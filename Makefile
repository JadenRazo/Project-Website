.PHONY: all dev dev-fresh dev-kill dev-fresh-kill build clean test lint migrate seed backup restore install start stop setup status update-code-stats fresh kill-sessions

# Default target
all: dev-fresh-kill

# Development commands
dev:
	@echo "Starting development environment..."
	@./start-dev.sh

dev-fresh:
	@echo "Starting development environment with fresh cache..."
	@./start-dev.sh --fresh

dev-kill:
	@echo "Starting development environment with port cleanup..."
	@./start-dev.sh --kill-existing

# Development with fresh cache and port cleanup
dev-fresh-kill:
	@echo "Starting development environment with fresh cache and port cleanup..."
	@./start-dev.sh --fresh --kill-existing

# Alternative npm-based commands
npm-dev:
	@echo "Starting development via npm..."
	@npm run dev

npm-dev-fresh:
	@echo "Starting development with fresh cache via npm..."
	@npm run dev:fresh

npm-dev-kill:
	@echo "Starting development with port cleanup via npm..."
	@npm run dev:kill

# Build all services
build:
	@echo "Building all services..."
	@npm run build

build-frontend:
	@echo "Building frontend..."
	@cd frontend && npm run build

build-backend:
	@echo "Building backend services..."
	@cd backend && go build -o bin/api cmd/api/main.go
	@cd backend && go build -o bin/devpanel cmd/devpanel/main.go
	@cd backend && go build -o bin/messaging cmd/messaging/main.go
	@cd backend && go build -o bin/urlshortener cmd/urlshortener/main.go
	@cd backend && go build -o bin/worker cmd/worker/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf frontend/build
	@rm -rf backend/bin
	@rm -rf node_modules
	@rm -rf frontend/node_modules
	@rm -rf frontend/.eslintcache
	@rm -rf frontend/tsconfig.tsbuildinfo

# Clean development caches
clean-cache:
	@echo "Cleaning development caches..."
	@cd frontend && npm run clear-cache 2>/dev/null || echo "Manual cache clearing..."
	@rm -rf frontend/node_modules/.cache 2>/dev/null || true
	@rm -rf frontend/.eslintcache 2>/dev/null || true
	@rm -rf frontend/tsconfig.tsbuildinfo 2>/dev/null || true

# Run tests
test:
	@echo "Running tests..."
	@cd frontend && npm test
	@cd backend && go test ./...

test-frontend:
	@echo "Running frontend tests..."
	@cd frontend && npm test

test-backend:
	@echo "Running backend tests..."
	@cd backend && go test ./...

# Type checking
type-check:
	@echo "Running TypeScript type checking..."
	@cd frontend && npm run type-check

# Run linters
lint:
	@echo "Running linters..."
	@cd frontend && npm run lint
	@./scripts/lint.sh

lint-fix:
	@echo "Running linters with auto-fix..."
	@cd frontend && npm run lint:fix
	@./scripts/lint.sh

# Database operations
migrate:
	@echo "Running database migrations..."
	@cd backend && go run cmd/migration/main.go up

migrate-down:
	@echo "Rolling back database migrations..."
	@cd backend && go run cmd/migration/main.go down

seed:
	@echo "Seeding database..."
	@./scripts/seed.sh

backup:
	@echo "Backing up database..."
	@./scripts/backup_db.sh

restore:
	@echo "Restoring database..."
	@./scripts/restore_db.sh

# Install dependencies
install:
	@echo "Installing dependencies..."
	@npm install
	@cd frontend && npm install
	@cd backend && go mod download

install-frontend:
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install

install-backend:
	@echo "Installing backend dependencies..."
	@cd backend && go mod download && go mod tidy

# System setup
setup:
	@echo "Setting up development environment..."
	@./setup.sh

# Service management
start:
	@echo "Starting all services..."
	@./deploy/start-services.sh

stop:
	@echo "Stopping all services..."
	@./stop-dev.sh 2>/dev/null || ./scripts/stop.sh 2>/dev/null || echo "Killing tmux sessions..."
	@tmux kill-session -t portfolio-frontend 2>/dev/null || true
	@tmux kill-session -t portfolio-backend 2>/dev/null || true

kill-sessions:
	@echo "Killing tmux sessions..."
	@tmux kill-session -t portfolio-frontend 2>/dev/null || true
	@tmux kill-session -t portfolio-backend 2>/dev/null || true

# Production
prod:
	@echo "Building and starting production environment..."
	@npm run build
	@npm run start

# Status and monitoring
status:
	@echo "Checking service status..."
	@echo "Frontend (port 3000):"
	@lsof -i :3000 || echo "  Not running"
	@echo "API (port 8080):"
	@lsof -i :8080 || echo "  Not running"
	@echo "DevPanel (port 8081):"
	@lsof -i :8081 || echo "  Not running"
	@echo "Messaging (port 8082):"
	@lsof -i :8082 || echo "  Not running"
	@echo "URL Shortener (port 8083):"
	@lsof -i :8083 || echo "  Not running"
	@echo "Worker (port 8084):"
	@lsof -i :8084 || echo "  Not running"

# Code statistics
update-code-stats:
	@echo "Updating code statistics..."
	@./scripts/update_code_stats.sh

# Development workflow shortcuts
fresh: clean-cache dev-fresh-kill

# Pre-commit workflow
pre-commit: lint-fix type-check test

# Help
help:
	@echo "Portfolio Website Development Makefile"
	@echo ""
	@echo "MAIN DEVELOPMENT COMMANDS:"
	@echo "  dev-fresh-kill  - Development with fresh cache and port cleanup"
	@echo "  dev-fresh       - Start with fresh cache clearing"
	@echo "  dev-kill        - Start with port cleanup"
	@echo "  dev             - Normal development start"
	@echo "  fresh           - Clean cache and start development"
	@echo ""
	@echo "BUILD & TEST:"
	@echo "  build           - Build all services"
	@echo "  build-frontend  - Build frontend only"
	@echo "  build-backend   - Build backend services only"
	@echo "  test            - Run all tests"
	@echo "  test-frontend   - Run frontend tests"
	@echo "  test-backend    - Run backend tests"
	@echo "  type-check      - Run TypeScript type checking"
	@echo ""
	@echo "CODE QUALITY:"
	@echo "  lint            - Run linters"
	@echo "  lint-fix        - Run linters with auto-fix"
	@echo "  pre-commit      - Run full pre-commit checks"
	@echo ""
	@echo "DATABASE:"
	@echo "  migrate         - Run database migrations"
	@echo "  migrate-down    - Rollback migrations"
	@echo "  seed            - Seed database"
	@echo "  backup          - Backup database"
	@echo "  restore         - Restore database"
	@echo ""
	@echo "SETUP & MAINTENANCE:"
	@echo "  install         - Install all dependencies"
	@echo "  setup           - Full development environment setup"
	@echo "  clean           - Clean build artifacts"
	@echo "  clean-cache     - Clean development caches"
	@echo "  status          - Check service status"
	@echo "  update-code-stats - Update code statistics"
	@echo ""
	@echo "SERVICE CONTROL:"
	@echo "  start           - Start all services"
	@echo "  stop            - Stop all services"
	@echo "  kill-sessions   - Kill tmux sessions"
	@echo "  prod            - Production build and start"
	@echo ""
	@echo "TMUX SESSIONS:"
	@echo "  Frontend: tmux attach -t portfolio-frontend"
	@echo "  Backend:  tmux attach -t portfolio-backend" 