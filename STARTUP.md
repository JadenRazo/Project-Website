# Portfolio Website Startup Guide

## Quick Start

### Prerequisites
- PostgreSQL running with the database schema loaded
- tmux installed (`sudo apt install tmux` or `brew install tmux`)
- Node.js 18+ and npm
- Go 1.21+

### One-Command Startup
```bash
# Fresh start with cache clearing (recommended)
./start-dev.sh --fresh

# Normal start
./start-dev.sh

# Stop everything
./stop-dev.sh
```

## Detailed Setup

### 1. Environment Configuration
```bash
# Copy and configure environment variables
cp .env.example .env
# Edit .env with your database credentials and other settings
```

### 2. Database Setup
- Ensure PostgreSQL is running
- Create database: `createdb project_website`
- Load schema using pgAdmin or: `psql -d project_website -f database/schema.sql`

### 3. Start Development Environment
```bash
./start-dev.sh --fresh
```

This will:
- ✅ Automatically detect and clean up any existing tmux sessions
- ✅ Check all dependencies (tmux, Node.js, Go, PostgreSQL)
- ✅ Test database connectivity
- ✅ Clear frontend caches
- ✅ Install/update dependencies
- ✅ Start backend services in tmux session (`portfolio-backend`)
- ✅ Start frontend in tmux session (`portfolio-frontend`)
- ✅ Verify all services are running

## Working with tmux Sessions

### Attach to Sessions
```bash
# Frontend session (React development server)
tmux attach -t portfolio-frontend

# Backend session (Go microservices)
tmux attach -t portfolio-backend
```

### Navigate Backend Services
The backend session has multiple windows:
- `api` - Main API service (port 8080)
- `devpanel` - Developer panel (port 8081)
- `messaging` - Messaging service (port 8082)
- `urlshortener` - URL shortener service (port 8083)

Switch between windows:
```bash
# While in tmux session
Ctrl+B, then press:
- 0 for api
- 1 for devpanel
- 2 for messaging
- 3 for urlshortener
```

### tmux Commands
```bash
# Detach from session (keeps it running)
Ctrl+B, then D

# List all sessions
tmux list-sessions

# Kill a specific session
tmux kill-session -t portfolio-frontend
tmux kill-session -t portfolio-backend
```

## Service URLs

Once started, access your services at:
- **Frontend**: http://localhost:3000
- **API**: http://localhost:8080
- **DevPanel**: http://localhost:8081
- **Messaging**: http://localhost:8082
- **URL Shortener**: http://localhost:8083

## Script Options

### start-dev.sh Options
```bash
./start-dev.sh [OPTIONS]

Options:
  -f, --fresh         Force fresh start with cache clearing
  -s, --skip-deps     Skip dependency installation/updates
  -v, --verbose       Enable verbose output
  -h, --help          Show help message

Note: The script now automatically detects and cleans up existing tmux sessions,
eliminating the need for the --kill-existing flag.
```

### Examples
```bash
# Fresh start with cache clearing
./start-dev.sh --fresh

# Skip dependency updates (faster startup)
./start-dev.sh --skip-deps

# Verbose output for debugging
./start-dev.sh --verbose
```

## Troubleshooting

### Port Conflicts
The `./start-dev.sh` script now automatically detects and cleans up existing tmux sessions, preventing most port conflicts. If you still encounter issues, the script will handle them automatically when you run:
```bash
./start-dev.sh
```

### Database Connection Issues
1. Ensure PostgreSQL is running
2. Check credentials in `.env` file
3. Test connection: `psql -h localhost -d project_website -U postgres`

### Cache Issues
```bash
# Force fresh start
./start-dev.sh --fresh

# Manual cache clearing
cd frontend && npm run clear-cache
```

### tmux Not Found
```bash
# Ubuntu/Debian
sudo apt install tmux

# macOS
brew install tmux

# CentOS/RHEL
sudo yum install tmux
```

### Services Not Starting
1. Check logs in tmux sessions
2. Verify all dependencies are installed
3. Check for port conflicts
4. Ensure database is accessible

## Development Tips

### Monitoring Logs
```bash
# Watch all tmux sessions
tmux list-sessions

# Attach and monitor specific service
tmux attach -t portfolio-backend
# Then switch to specific window (Ctrl+B, 0-3)
```

### Quick Restart
```bash
# Stop everything
./stop-dev.sh

# Start fresh
./start-dev.sh --fresh
```

### Working with the Frontend
- Hot module replacement is enabled
- Changes auto-refresh the browser
- React DevTools work normally
- TypeScript errors show in console

### Working with the Backend
- File watching enabled for Go files
- Services restart automatically on changes
- Each service runs in its own tmux window
- Database connections are tested on startup

## Performance Notes

- First startup takes longer (dependency installation)
- Subsequent startups are faster with `--skip-deps`
- Fresh cache clearing adds ~10-30 seconds
- All services start in parallel for speed

## Security Notes

- `.env` file contains sensitive data - never commit it
- Default admin user is created (change password immediately!)
- Database connection is tested but not required for startup
- Services bind to localhost only by default