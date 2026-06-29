# SwiftAntiCheat

Full screenshare & anti-cheat system for Minecraft servers — plugin + PC scanner + web dashboard.

## System Components

| Component | Tech | What it does |
|-----------|------|-------------|
| **Plugin** | Java (Spigot/Paper) | `/ss`, freeze, chat mute, dodge-ban, `/swiftanticheat` link gen |
| **Backend API** | Node.js + SQLite | Link generation, scan storage, player history, HWID tracking |
| **Scanner** | Go (Windows) | PC binary the suspect runs — scans processes/files/registry/mods |
| **Dashboard** | Next.js | Web UI to view scans, check flags, search players, HWID tracking |

## How the Flow Works

1. **Staff suspects a player** → runs `/ss <player>` or `/swiftanticheat <player>`
2. **Player gets frozen** in-game (can't move, can chat/look around)
3. **Scan link generated** — unique URL sent to the player's chat
4. **Player clicks link** → downloads the scanner EXE, runs it
5. **Scanner scans their PC** (processes, files, Windows artifacts, Minecraft mods)
6. **Results upload** to the backend API
7. **Staff views results** on the web dashboard
8. **If player leaves** while frozen → automatic ban (configurable duration)
9. **Chat is silenced** only for the suspect during screenshare

## Quick Start

### 1. Backend API (Deploy First)

```bash
cd backend
npm install
cp .env.example .env
# Edit .env — set a random API_KEY
npm start
```

The API runs on `http://localhost:3000` (or your Render/Railway URL).

### 2. Plugin

Build with Maven:
```bash
cd plugin
mvn clean package
```

Copy `target/swift-anticheat-1.0.0.jar` to your server's `plugins/` folder.

Edit `plugins/SwiftAntiCheat/config.yml`:
```yaml
api:
  base-url: "https://your-api.onrender.com"  # your deployed backend URL
  api-key: "your-secret-key"
```

### 3. Dashboard

```bash
cd dashboard
npm install
npm run build
npm start
```

Or deploy to Vercel (free):
```bash
npx vercel --prod
```

### 4. Scanner (for players)

```bash
cd scanner
go build -o swiftac-scanner.exe
```

Players download this EXE from the scan link and run it with their scan ID:
```bash
swiftac-scanner.exe -id YOUR_SCAN_ID
```

## Free Hosting Guide

### Backend (Render.com — Free)

1. Create account at https://render.com
2. New → Web Service → connect your repo
3. Set:
   - Build Command: `cd backend && npm install`
   - Start Command: `cd backend && npm start`
   - Environment Variable: `API_KEY=your-random-secret-key`
4. Deploy — get URL like `https://swiftac-api.onrender.com`

### Database (SQLite)

The backend uses SQLite by default (stored as `data.db`). For free hosted SQL:
- **Supabase**: https://supabase.com (free 500MB PostgreSQL)
  - Change `src/index.js` to use `pg` instead of `better-sqlite3`

### Dashboard (Vercel — Free)

1. Create account at https://vercel.com
2. New Project → import `dashboard/` folder
3. Deploy — get URL like `https://swiftac-scan.vercel.app`
4. Staff visit this URL — enter the API URL in the top-right input

### Scanner Distribution

Upload the compiled `swiftac-scanner.exe` to:
- **GitHub Releases**: free, unlimited downloads
- **Cloudflare R2**: 10GB free storage

The scan link page (dashboard) should serve as a download page for the EXE.

## Commands

| Command | Permission | Description |
|---------|-----------|-------------|
| `/ss <player>` | `swiftac.staff` | Freeze a player for screenshare |
| `/uss <player>` | `swiftac.staff` | Unfreeze a player |
| `/swiftanticheat <player>` | `swiftac.staff` | Generate scan link + freeze |
| `/discord` | everyone | Get Discord invite link |

## Configuration

See `plugin/src/main/resources/config.yml` for all options:
- Ban duration for dodge
- Custom freeze messages
- Discord invite link
- API base URL and key
- Toggle AI freeze vs basic freeze

## Scanner Detection Methods

All user-mode (no kernel driver — safe):

- **Process scanning** — checks running processes against known cheat names
- **Filesystem scanning** — scans Desktop, Downloads, TEMP for cheat files
- **Windows Prefetch** — checks execution history for cheat executables
- **Registry scanning** — checks autorun entries for cheat persistence
- **PowerShell history** — checks for suspicious commands
- **Event log check** — detects cleared event logs
- **Minecraft mod scanning** — scans `.minecraft/mods/` for cheat mods
- **HWID fingerprinting** — hardware ID for ban evasion tracking
