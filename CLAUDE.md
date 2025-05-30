# Portfolio Website - Claude Development Guide

## Project Overview
This is a full-stack portfolio website with React/TypeScript frontend and Go backend microservices architecture. The system includes URL shortening, messaging, and developer panel functionality.

## Common Commands

### Development
- `./start-dev.sh` - **RECOMMENDED** Start all services in tmux sessions
- `./start-dev.sh --fresh` - Start with cache clearing and fresh setup
- `./start-dev.sh --kill-existing` - Kill existing processes on required ports
- `./start-dev.sh --skip-deps` - Skip dependency installation
- `./start-dev.sh --verbose` - Enable verbose output
- `./start-dev.sh --help` - Show help message
- `npm run dev` - Alias for ./start-dev.sh
- `npm run dev:fresh` - Alias for ./start-dev.sh --fresh
- `npm run dev:kill` - Alias for ./start-dev.sh --kill-existing

### Building & Production
- `npm run build` - Build frontend for production
- `go build -o bin/api cmd/api/main.go` - Build API service
- `go build -o bin/messaging cmd/messaging/main.go` - Build messaging service
- `go build -o bin/urlshortener cmd/urlshortener/main.go` - Build URL shortener
- `go build -o bin/devpanel cmd/devpanel/main.go` - Build developer panel
- `go build -o bin/worker cmd/worker/main.go` - Build worker service (scheduled tasks)

### Testing & Quality
- `npm run test` - Run frontend tests
- `npm run type-check` - TypeScript type checking
- `npm run lint` - Run ESLint
- `npm run lint:fix` - Auto-fix linting issues
- `npm run format` - Format code with Prettier
- `go test ./...` - Run Go tests
- `./scripts/lint.sh` - Run Go linting

### Database
- `go run cmd/migration/main.go up` - Run database migrations
- `go run cmd/migration/main.go down` - Rollback migrations
- `./scripts/seed.sh` - Seed development data
- `./scripts/backup_db.sh` - Backup database
- `./scripts/restore_db.sh` - Restore database

### Git Workflow
- `git status` - Check current changes
- `git diff` - View uncommitted changes
- `git log --oneline -10` - View recent commits
- Always run `npm run lint:fix` and `npm run type-check` before committing frontend changes
- Always run `./scripts/lint.sh` before committing backend changes

## Code Style Guidelines

### TypeScript/React (Frontend)
- Use ES modules syntax (`import/export`), NOT CommonJS (`require`)
- Destructure imports when possible: `import { useState, useEffect } from 'react'`
- Use functional components with hooks, NOT class components
- Use TypeScript interfaces for props and state
- Place component files in appropriate directories under `src/components/`
- Use styled-components for styling when possible
- Follow existing naming conventions: PascalCase for components, camelCase for functions
- IMPORTANT: DO NOT add comments unless explicitly requested

### Go (Backend)
- Follow standard Go conventions and idioms
- Use meaningful variable and function names
- Keep functions small and focused
- Handle errors explicitly - never ignore error returns
- Use dependency injection for testability
- Place domain logic in `internal/domain/`
- Place HTTP handlers in `internal/*/delivery/http/`
- Place data access in `internal/*/repository/`
- IMPORTANT: DO NOT add comments unless explicitly requested

## Project Structure & Key Files

### Frontend Core Files
- `frontend/src/App.tsx` - Main application component
- `frontend/src/index.tsx` - Application entry point
- `frontend/src/pages/` - Page components
- `frontend/src/components/` - Reusable UI components
- `frontend/src/hooks/` - Custom React hooks
- `frontend/src/contexts/` - React contexts (Theme, ZIndex)
- `frontend/src/styles/` - Global styles and themes

### Backend Core Files
- `backend/cmd/*/main.go` - Service entry points
- `backend/internal/app/bootstrap.go` - Application initialization
- `backend/internal/common/` - Shared utilities (auth, cache, database)
- `backend/internal/messaging/` - Messaging service implementation
- `backend/internal/urlshortener/` - URL shortener implementation
- `backend/config/` - Configuration files

### Configuration Files
- `backend/config/development.yaml` - Development config
- `backend/config/production.yaml` - Production config
- `frontend/.env` - Frontend environment variables
- `frontend/.env.development` - Frontend development-specific variables
- `frontend/craco.config.js` - CRACO configuration for webpack overrides
- `backend/.env` - Backend environment variables

## Service Ports
- Frontend: http://localhost:3000
- Main API: http://localhost:8080
- DevPanel: http://localhost:8081
- Messaging: http://localhost:8082
- URL Shortener: http://localhost:8083
- Worker: http://localhost:8084 (Background tasks)

## Key Features
- Status Page: http://localhost:3000/status (Service health monitoring)
- Code Stats API: http://localhost:8080/api/code-stats (Lines of code tracking)
- Health Check: http://localhost:8080/api/status (Backend status)

## Testing Guidelines
- Write unit tests for new features
- Test files should be colocated with source files
- Frontend tests use Jest and React Testing Library
- Backend tests use Go's built-in testing package
- Run tests before committing changes
- Ensure all tests pass in CI/CD pipeline

## Development Workflow
1. Always pull latest changes before starting work
2. **Start development environment**: `./start-dev.sh --fresh`
3. Create feature branches from `main`
4. **Work in tmux sessions**: 
   - Frontend: `tmux attach -t portfolio-frontend`
   - Backend: `tmux attach -t portfolio-backend`
   - Navigate windows: `Ctrl+B` then window number (0-4)
5. Run linting and type checking before commits
6. Write descriptive commit messages
7. Keep commits focused and atomic
8. Update tests when changing functionality
9. **Stop services when done**: 
   - `tmux kill-session -t portfolio-frontend`
   - `tmux kill-session -t portfolio-backend`

## Performance Considerations
- Frontend uses React.memo and useMemo for optimization
- Implement lazy loading for heavy components
- Use virtual scrolling for long lists (VirtualizedList component)
- Backend uses Redis for caching frequently accessed data
- Database queries should use appropriate indexes
- Monitor memory usage with MemoryManager component

## Security Best Practices
- NEVER commit secrets or API keys
- Use environment variables for sensitive configuration
- Validate all user inputs on both frontend and backend
- Use prepared statements for database queries
- Implement proper CORS configuration
- Use JWT tokens for authentication
- Hash passwords with bcrypt

## Debugging Tips
- Frontend: Use React Developer Tools and Redux DevTools
- Backend: Use delve debugger (`dlv debug`)
- Check service health endpoints: `/health`
- Monitor logs in `logs/` directory
- Use `debugHelpers.ts` utilities for frontend debugging
- Enable verbose logging in development mode

## Common Issues & Solutions
- Port conflicts: Check `lsof -i :PORT` and kill conflicting processes
- Database connection errors: Verify PostgreSQL is running and credentials are correct
- CORS issues: Check backend CORS middleware configuration
- Build failures: Clear caches (`npm cache clean --force`, `go clean -modcache`)
- WebSocket connection issues: Ensure correct WS_URL in frontend env
- **Browser cache issues during development**:
  - Run `npm run dev:fresh` to start with cleared caches
  - Use `npm run clear-cache` to manually clear all caches
  - In browser console: run `window.__clearAllCaches()` (dev mode only)
  - Hard refresh: Ctrl+Shift+R (Cmd+Shift+R on Mac)
  - Check for lingering service workers in DevTools > Application > Service Workers

## Deployment Checklist
- [ ] Run all tests and ensure they pass
- [ ] Build frontend production bundle
- [ ] Build all backend services
- [ ] Run database migrations
- [ ] Update environment variables
- [ ] Configure Nginx reverse proxy
- [ ] Set up SSL certificates
- [ ] Configure systemd services
- [ ] Set up monitoring and logging

## Repository Conventions
- Branch naming: `feature/description`, `fix/description`, `chore/description`
- PR titles should be descriptive and include ticket numbers if applicable
- Squash commits when merging to keep history clean
- Update README.md when adding new features or changing setup
- Document API changes in `backend/docs/api/`
- Keep dependencies up to date but test thoroughly

## Important Reminders
- This project uses both TypeScript and Go - ensure you're using the right toolchain
- Frontend and backend are separate but interdependent - test integration points
- The messaging service uses WebSockets - test real-time functionality
- URL shortener tracks analytics - verify stats are being recorded
- DevPanel is for internal use only - ensure proper access controls
- Always check git status before committing to avoid accidental file additions
- **Frontend uses CRACO** for webpack configuration without ejecting from Create React App
- **Cache clearing is automatic** in development mode on startup
- **DevCacheManager** provides utilities for managing browser caches in development