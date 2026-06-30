const { Router } = require('express');

module.exports = function reportsRouter(db) {
  const router = Router();

  router.post('/upload', (req, res) => {
    const { reportCode, results, hwidHash } = req.body;
    if (!reportCode || !results) {
      return res.status(400).json({ error: 'Missing required fields' });
    }

    const existing = db.prepare('SELECT * FROM reports WHERE report_code = ?').get(reportCode);
    if (existing) {
      return res.status(409).json({ error: 'Report code already exists' });
    }

    db.prepare(
      'INSERT INTO reports (report_code, results, hwid_hash) VALUES (?, ?, ?)'
    ).run(reportCode, JSON.stringify(results), hwidHash || '');

    res.json({ status: 'ok', reportCode });
  });

  router.get('/:code', (req, res) => {
    const report = db.prepare('SELECT * FROM reports WHERE report_code = ?').get(req.params.code);
    if (!report) return res.status(404).json({ error: 'Report not found' });

    res.json({
      reportCode: report.report_code,
      results: report.results ? JSON.parse(report.results) : null,
      hwidHash: report.hwid_hash,
      createdAt: report.created_at
    });
  });

  return router;
};
