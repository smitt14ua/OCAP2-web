# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

OCAP2-Web is a Go-based web server for the Operation Capture And Playback system. It archives and replays Arma 3 military simulation missions, rendering unit positions, movements, and combat events on interactive Leaflet maps.

## Build Commands

```bash
# Build (Linux)
go build -o ocap-webserver ./cmd

# Build (Windows)
go build -o ocap-webserver.exe ./cmd

# Run tests
go test ./...

# Run single test
go test -run TestName ./server

# Docker build
docker build -t ocap-webserver .
```

## Architecture

### Backend (Go)

Entry point: `cmd/main.go` initializes repositories and starts the Echo server.

**Core packages in `server/`:**
- `handler.go` - HTTP endpoints using Echo framework
- `operation.go` - SQLite repository for mission metadata (RepoOperation)
- `marker.go` - Dynamic marker image generation with color transforms (RepoMarker)
- `ammo.go` - Equipment/gear icon lookup (RepoAmmo)
- `setting.go` - Configuration via Viper (env vars, JSON/YAML files)

**API Endpoints (all under configurable `prefixURL`, default `/aar/`):**
- `GET /api/v1/operations` - Query missions with filters
- `POST /api/v1/operations/add` - Upload mission (requires `secret`)
- `GET /api/v1/customize` - UI customization settings
- `GET /api/version` - Build info
- `GET /data/:name` - Stream mission data (gzipped JSON)
- `GET /file/:name` - Download mission as attachment
- `GET /images/markers/:name/:color` - Generate colored marker
- `GET /images/markers/magicons/:name` - Equipment icons
- `GET /images/maps/*` - Map tiles

### Frontend (JavaScript)

Static SPA in `static/` using Leaflet for map rendering. No build step required.

**Key files in `static/scripts/`:**
- `ocap.js` - Main app, map init, playback engine
- `ocap.ui.js` - UI panels (unit list, events, timeline)
- `ocap.event.js` - Mission event handling
- `ocap.marker.js` - Unit/vehicle marker rendering
- `ocap.entity.js` - Base entity class
- `localizable.js` - i18n support

### Data Storage

- **SQLite database** - Mission metadata (`db/` or `OCAP_DB`)
- **Mission files** - Gzipped JSON in `data/` or `OCAP_DATA`
- **Map tiles** - Downloaded separately to `maps/` or `OCAP_MAPS`
- **Assets** - `markers/` and `ammo/` directories contain mod-specific icons

## Configuration

Settings loaded via Viper with priority: environment variables → config files → defaults.

**Config file:** `setting.json`
```json
{
  "listen": "127.0.0.1:5000",
  "prefixURL": "/aar/",
  "secret": "same-secret",
  "logger": true
}
```

**Environment variables** (prefix `OCAP_`):
- `OCAP_LISTEN`, `OCAP_SECRET`, `OCAP_PREFIXURL`
- `OCAP_DB`, `OCAP_DATA`, `OCAP_MAPS`, `OCAP_MARKERS`, `OCAP_AMMO`, `OCAP_STATIC`
- `OCAP_CUSTOMIZE_WEBSITEURL`, `OCAP_CUSTOMIZE_WEBSITELOGO`

## Key Implementation Details

- Path traversal protection via `paramPath()` using `filepath.Clean()` validation
- Marker colors: predefined names (blufor, opfor, ind, civ, dead, hit) or hex codes
- Cache control: 7 days for assets, no-cache for HTML
- Build metadata injected via ldflags: `BuildCommit`, `BuildDate`
