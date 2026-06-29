const { Router } = require('express');
const authMiddleware = require('../middleware/auth');

module.exports = function playersRouter(db) {
  const router = Router();

  router.get('/:uuid', authMiddleware, (req, res) => {
    const { uuid } = req.params;

    const bans = db.prepare(
      'SELECT * FROM bans WHERE player_uuid = ? ORDER BY created_at DESC'
    ).all(uuid);

    const scans = db.prepare(
      'SELECT id, status, created_at, completed_at FROM scans WHERE player_uuid = ? ORDER BY created_at DESC'
    ).all(uuid);

    const hwid = db.prepare('SELECT * FROM hwids WHERE player_uuid = ?').get(uuid);

    res.json({ uuid, bans, scans, hwid });
  });

  router.post('/:uuid/ban', authMiddleware, (req, res) => {
    const { uuid } = req.params;
    const { playerName, reason, duration, bannedBy, expiresAt } = req.body;

    if (!playerName || !reason) {
      return res.status(400).json({ error: 'Missing required fields' });
    }

    const stmt = db.prepare(
      'INSERT INTO bans (player_uuid, player_name, reason, duration, banned_by, expires_at) VALUES (?, ?, ?, ?, ?, ?)'
    );
    const result = stmt.run(uuid, playerName, reason, duration || 'permanent', bannedBy || 'console', expiresAt || null);

    res.json({ id: result.lastInsertRowid, uuid, playerName, reason });
  });

  router.post('/:uuid/hwid', authMiddleware, (req, res) => {
    const { uuid } = req.params;
    const { hwidHash } = req.body;

    if (!hwidHash) return res.status(400).json({ error: 'Missing hwidHash' });

    db.prepare(
      'INSERT INTO hwids (player_uuid, hwid_hash, last_seen) VALUES (?, ?, unixepoch()) ON CONFLICT(player_uuid) DO UPDATE SET hwid_hash = ?, last_seen = unixepoch()'
    ).run(uuid, hwidHash, hwidHash);

    res.json({ status: 'ok' });
  });

  return router;
};
