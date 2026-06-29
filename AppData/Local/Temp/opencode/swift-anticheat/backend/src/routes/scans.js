const { Router } = require('express');
const { v4: uuidv4 } = require('uuid');
const authMiddleware = require('../middleware/auth');

module.exports = function scansRouter(db) {
  const router = Router();

  router.post('/create', authMiddleware, (req, res) => {
    const { playerName, playerUUID, staffName } = req.body;
    if (!playerName || !playerUUID || !staffName) {
      return res.status(400).json({ error: 'Missing required fields' });
    }

    const id = uuidv4();
    const stmt = db.prepare(
      'INSERT INTO scans (id, player_name, player_uuid, staff_name) VALUES (?, ?, ?, ?)'
    );
    stmt.run(id, playerName, playerUUID, staffName);

    res.json({ id, playerName, playerUUID, staffName, status: 'pending' });
  });

  router.get('/:id', (req, res) => {
    const scan = db.prepare('SELECT * FROM scans WHERE id = ?').get(req.params.id);
    if (!scan) return res.status(404).json({ error: 'Scan not found' });

    let results = null;
    if (scan.results) {
      try { results = JSON.parse(scan.results); } catch (e) {}
    }

    res.json({
      id: scan.id,
      playerName: scan.player_name,
      playerUUID: scan.player_uuid,
      staffName: scan.staff_name,
      status: scan.status,
      createdAt: scan.created_at,
      completedAt: scan.completed_at,
      results
    });
  });

  router.post('/:id/results', (req, res) => {
    const { results, hwidHash } = req.body;
    if (!results) return res.status(400).json({ error: 'Missing results' });

    const scan = db.prepare('SELECT * FROM scans WHERE id = ?').get(req.params.id);
    if (!scan) return res.status(404).json({ error: 'Scan not found' });

    db.prepare(
      'UPDATE scans SET results = ?, status = ?, completed_at = unixepoch() WHERE id = ?'
    ).run(JSON.stringify(results), 'completed', req.params.id);

    if (hwidHash && scan.player_uuid) {
      db.prepare(
        'INSERT INTO hwids (player_uuid, hwid_hash, last_seen) VALUES (?, ?, unixepoch()) ON CONFLICT(player_uuid) DO UPDATE SET hwid_hash = ?, last_seen = unixepoch()'
      ).run(scan.player_uuid, hwidHash, hwidHash);
    }

    res.json({ status: 'ok' });
  });

  router.get('/', (req, res) => {
    const { player, status, limit = 50, offset = 0 } = req.query;
    let query = 'SELECT * FROM scans WHERE 1=1';
    const params = [];

    if (player) { query += ' AND (player_name = ? OR player_uuid = ?)'; params.push(player, player); }
    if (status) { query += ' AND status = ?'; params.push(status); }
    query += ' ORDER BY created_at DESC LIMIT ? OFFSET ?';
    params.push(parseInt(limit), parseInt(offset));

    const scans = db.prepare(query).all(...params).map(s => ({
      ...s,
      results: s.results ? JSON.parse(s.results) : null
    }));

    res.json(scans);
  });

  return router;
};
