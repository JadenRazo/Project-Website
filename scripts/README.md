# Scripts Directory Organization

This directory contains all scripts for the Portfolio Website project, organized by functionality.

## Directory Structure

```
scripts/
├── backend/          # Backend-specific management scripts
├── database/         # Database management and migrations
├── deployment/       # Production deployment scripts
├── development/      # Development environment scripts
├── monitoring/       # Monitoring and metrics scripts
├── ssl/             # SSL certificate management
├── testing/         # Testing and validation scripts
└── utilities/       # General utility scripts
```

## Quick Access

Main scripts are symlinked in the project root for convenience:
- `./start-dev.sh` → Development environment startup
- `./start-prod.sh` → Production environment startup
- `./stop-dev.sh` → Stop development services
- `./stop-prod.sh` → Stop production services

## Script Categories

### Development (`development/`)
- **start-dev.sh** - Start development environment with hot reload
- **start-prod.sh** - Start production-optimized environment locally
- **stop-dev.sh** - Stop development tmux sessions
- **stop-prod.sh** - Stop production tmux sessions
- **lint.sh** - Run linting checks on codebase
- **optimize-images.sh** - Optimize image assets

### Database (`database/`)
- **seed.sh** - Seed database with initial data
- **backup_db.sh** - Backup PostgreSQL database
- **restore_db.sh** - Restore database from backup
- **check-db-status.sh** - Check database connectivity
- **update-database.sh** - Apply database migrations

### Deployment (`deployment/`)
- **deploy-local.sh** - Deploy to local environment
- **start-backend.sh** - Start backend services only
- **start-services.sh** - Start all services in production
- **INSTALL.sh** - Initial server setup and installation
- **check_system.sh** - Verify system requirements
- **install-systemd-service.sh** - Install systemd services

### Backend (`backend/`)
- **run.sh** - Run backend services
- **status.sh** - Check backend service status
- **stop.sh** - Stop backend services
- **autostart.sh** - Configure auto-start on boot

### Monitoring (`monitoring/`)
- **update_code_stats.sh** - Update code statistics
- **setup_code_stats_cron.sh** - Setup cron job for stats

### SSL (`ssl/`)
- **setup-ssl-automation.sh** - Configure SSL auto-renewal
- **ssl-auto-renew.sh** - Renew SSL certificates
- **ssl-monitor.sh** - Monitor SSL certificate status
- **fix-ssl-now.sh** - Quick SSL fix script

### Testing (`testing/`)
- **test_db_connection.sh** - Test database connectivity

### Utilities (`utilities/`)
- **backend-lint.sh** - Lint backend Go code
- **backend-setup.sh** - Setup backend dependencies
- **initial-setup.sh** - Initial project setup
- **setup-admin.sh** - Configure admin user
- **setup-redis-security.sh** - Configure Redis security

## Usage Examples

### Development Workflow
```bash
# Start development environment (automatic session cleanup included)
./start-dev.sh

# Start with fresh build and cache clear
./start-dev.sh --fresh

# Stop development environment
./stop-dev.sh
```

### Production Deployment
```bash
# Start production environment locally
./start-prod.sh

# Deploy to production server
cd scripts/deployment
./INSTALL.sh          # First-time setup
./start-services.sh   # Start all services
```

### Database Management
```bash
# Seed database
./scripts/database/seed.sh

# Backup database
./scripts/database/backup_db.sh

# Check database status
./scripts/database/check-db-status.sh
```

### Monitoring
```bash
# Update code statistics
./scripts/monitoring/update_code_stats.sh

# Setup automatic code stats updates
./scripts/monitoring/setup_code_stats_cron.sh
```

## Script Conventions

1. **Naming**: Use kebab-case for script names
2. **Extension**: All scripts should have `.sh` extension
3. **Shebang**: Start with `#!/bin/bash`
4. **Error Handling**: Use `set -e` for error handling
5. **Documentation**: Include usage comments at the top
6. **Permissions**: Ensure scripts are executable (`chmod +x`)

## Environment Variables

Scripts may require these environment variables:
- `NODE_ENV` - Node environment (development/production)
- `ENVIRONMENT` - Go environment (development/production)
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` - Database config
- `API_PORT`, `DEVPANEL_PORT`, `MESSAGING_PORT`, etc. - Service ports

## Troubleshooting

### Scripts not executing
```bash
# Make script executable
chmod +x scripts/path/to/script.sh
```

### Port conflicts
```bash
# The start scripts now automatically detect and clean up existing tmux sessions
./start-dev.sh
```

### Database connection issues
```bash
# Test database connection
./scripts/testing/test_db_connection.sh
```

## Adding New Scripts

1. Place script in appropriate category directory
2. Make it executable: `chmod +x script.sh`
3. Update this README with description
4. Create symlink in root if it's a primary script
5. Test thoroughly before committing

## Maintenance

Regular maintenance tasks:
- Review and remove unused scripts
- Update script documentation
- Ensure all scripts follow conventions
- Test scripts after major changes