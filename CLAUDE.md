# Portfolio Website - Claude Development Guide

## Project Overview
Full-stack portfolio website: React/TypeScript frontend + Go backend microservices (URL shortener, messaging, devpanel).

## Essential Commands
**Start Development:** `./start-dev.sh --fresh` 
**Build:** `npm run build` (frontend) | `go build -o bin/[service] cmd/[service]/main.go` (backend)
**Quality Checks:** `npm run lint:fix && npm run type-check` (frontend) | `./scripts/lint.sh` (backend)
**Database:** `go run cmd/migration/main.go up` | `./scripts/seed.sh`

## Service Ports
- Frontend: :3000 | Main API: :8080 | DevPanel: :8081 | Messaging: :8082 | URL Shortener: :8083 | Worker: :8084

## File Maintenance Notice
**IMPORTANT:** When adding new features or refactoring, update this file with objectively important information that would help future development. Keep additions concise and focused on essential patterns, commands, or architectural decisions.

## Code Style Guidelines
**Frontend:** ES modules, functional components + hooks, TypeScript interfaces, PascalCase components, camelCase functions
**Backend:** Standard Go idioms, explicit error handling, dependency injection, domain in `internal/domain/`, handlers in `internal/*/delivery/http/`
**Universal:** NO comments unless explicitly requested

## Key Architecture
**Frontend:** `src/App.tsx` (main) | `src/pages/` | `src/components/` | `src/hooks/` | `src/contexts/`
**Assets:** `src/assets/icons/` (favicons) | `src/assets/images/` | `src/assets/videos/` | `src/assets/data/` (JSON) | `src/config/` (env-config.js)
**Backend:** `cmd/*/main.go` (entry) | `internal/*/` (services) | `internal/common/` (shared utils)
**Config:** `backend/config/*.yaml` | `frontend/.env*` | `craco.config.js`

## Development Workflow
1. `git pull` → `./start-dev.sh --fresh` → feature branch from `main`
2. Tmux sessions: `tmux attach -t portfolio-{frontend|backend}` (Ctrl+B + window#)
3. Before commit: run quality checks, write tests, atomic commits

## Critical Security
- NEVER commit secrets/keys - use env vars
- Validate all inputs, use prepared statements, JWT auth, bcrypt passwords

## SQL/Database Guidelines
- Add SQL data to existing .sql files when coherent - avoid creating new .sql files unless absolutely necessary
- Only create new .sql files when it's the best possible solution for organization

## Common Fixes
- Port conflicts: `lsof -i :PORT` and kill processes
- Cache issues: `npm run dev:fresh` or `window.__clearAllCaches()` in dev console
- WebSocket issues: check WS_URL in frontend env