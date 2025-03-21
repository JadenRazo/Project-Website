.PHONY: all dev build clean test lint migrate seed backup restore

# Default target
all: dev

# Development
dev:
	@echo "Starting development environment..."
	@npm run dev

# Build all services
build:
	@echo "Building all services..."
	@npm run build

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf frontend/build
	@rm -rf backend/bin
	@rm -rf node_modules
	@rm -rf frontend/node_modules

# Run tests
test:
	@echo "Running tests..."
	@cd frontend && npm test
	@cd backend && go test ./...

# Run linters
lint:
	@echo "Running linters..."
	@cd frontend && npm run lint
	@cd backend && go vet ./...
	@cd backend && golangci-lint run

# Database migrations
migrate:
	@echo "Running database migrations..."
	@cd backend && go run cmd/migration/main.go up

# Seed database
seed:
	@echo "Seeding database..."
	@./scripts/seed.sh

# Backup database
backup:
	@echo "Backing up database..."
	@./scripts/backup_db.sh

# Restore database
restore:
	@echo "Restoring database..."
	@./scripts/restore_db.sh

# Install dependencies
install:
	@echo "Installing dependencies..."
	@npm install
	@cd frontend && npm install
	@cd backend && go mod download

# Start all services
start:
	@echo "Starting all services..."
	@./scripts/start-services.sh

# Stop all services
stop:
	@echo "Stopping all services..."
	@killall -9 api devpanel messaging urlshortener || true

# Development with hot reload
dev-reload:
	@echo "Starting development with hot reload..."
	@npm run dev

# Production build and start
prod:
	@echo "Building and starting production environment..."
	@npm run build
	@npm run start

# Help
help:
	@echo "Available targets:"
	@echo "  all        - Start development environment (default)"
	@echo "  dev        - Start development environment"
	@echo "  build      - Build all services"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linters"
	@echo "  migrate    - Run database migrations"
	@echo "  seed       - Seed database"
	@echo "  backup     - Backup database"
	@echo "  restore    - Restore database"
	@echo "  install    - Install dependencies"
	@echo "  start      - Start all services"
	@echo "  stop       - Stop all services"
	@echo "  dev-reload - Start development with hot reload"
	@echo "  prod       - Build and start production environment" 