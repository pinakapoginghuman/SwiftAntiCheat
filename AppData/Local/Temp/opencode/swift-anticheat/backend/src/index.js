const express = require('express');
const cors = require('cors');
const path = require('path');
const Database = require('better-sqlite3');

const scansRouter = require('./routes/scans');
const playersRouter = require('./routes/players');

const app = express();
const PORT = process.env.PORT || 3000;

const db = new Database(path.join(__dirname, '..', 'data.db'));
db.pragma('journal_mode = WAL');

db.exec(`
  CREATE TABLE IF NOT EXISTS scans (
    id TEXT PRIMARY KEY,
    player_name TEXT,
    player_uuid TEXT,
    staff_name TEXT,
    status TEXT DEFAULT 'pending',
    created_at INTEGER DEFAULT (unixepoch()),
    completed_at INTEGER,
    results TEXT
  );
  CREATE TABLE IF NOT EXISTS bans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_uuid TEXT,
    player_name TEXT,
    reason TEXT,
    duration TEXT,
    banned_by TEXT,
    created_at INTEGER DEFAULT (unixepoch()),
    expires_at INTEGER
  );
  CREATE TABLE IF NOT EXISTS hwids (
    player_uuid TEXT PRIMARY KEY,
    hwid_hash TEXT,
    last_seen INTEGER DEFAULT (unixepoch())
  );
`);

app.use(cors());
app.use(express.json({ limit: '50mb' }));

app.use('/api/scans', scansRouter(db));
app.use('/api/players', playersRouter(db));

app.get('/api/health', (req, res) => {
  res.json({ status: 'ok', timestamp: Date.now() });
});

app.listen(PORT, () => {
  console.log(`SwiftAntiCheat API running on port ${PORT}`);
});
