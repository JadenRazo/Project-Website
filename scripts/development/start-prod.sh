#!/bin/bash

# =============================================================================
# Portfolio Website Production Startup Script
# =============================================================================
# This script builds the frontend for production and starts both frontend 
# and backend in separate tmux sessions optimized for production deployment
# =============================================================================

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Configuration
# Get the real path of the script, resolving symlinks
SCRIPT_PATH="$(readlink -f "${BASH_SOURCE[0]}")"
PROJECT_ROOT="$(cd "$(dirname "$SCRIPT_PATH")/../.." && pwd)"
FRONTEND_DIR="$PROJECT_ROOT/frontend"
BACKEND_DIR="$PROJECT_ROOT/backend"
SESSION_PREFIX="portfolio"
FRONTEND_SESSION="${SESSION_PREFIX}-frontend-prod"
BACKEND_SESSION="${SESSION_PREFIX}-backend-prod"

# Default ports
FRONTEND_PORT=3000
API_PORT=8080
DEVPANEL_PORT=8081
MESSAGING_PORT=8082
URLSHORTENER_PORT=8083
WORKER_PORT=8084

# Script options
FORCE_FRESH=false
SKIP_DEPS=false
VERBOSE=false
KILL_EXISTING=false

# Load environment variables early
if [ -f "$BACKEND_DIR/.env.production" ]; then
    set -a
    source "$BACKEND_DIR/.env.production"
    set +a
elif [ -f "$BACKEND_DIR/.env" ]; then
    set -a
    source "$BACKEND_DIR/.env"
    set +a
fi

# =============================================================================
# Helper Functions
# =============================================================================

print_banner() {
    echo -e "${CYAN}"
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                  Portfolio Website PRODUCTION Launcher                      ║"
    echo "║                                                                              ║"
    echo "║  🚀 Frontend: Optimized React Build (port $FRONTEND_PORT)                           ║"
    echo "║  ⚙️  Backend: Go microservices (ports $API_PORT-$WORKER_PORT)                               ║"
    echo "║  🗄️  Database: PostgreSQL                                                    ║"
    echo "║  📺 Sessions: tmux ($FRONTEND_SESSION, $BACKEND_SESSION)               ║"
    echo "║  🎯 Mode: PRODUCTION (optimized, minified, no dev tools)                    ║"
    echo "╚══════════════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_updated_banner() {
    echo -e "${CYAN}"
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                  Portfolio Website PRODUCTION Launcher                      ║"
    echo "║                                                                              ║"
    printf "║  🚀 Frontend: Optimized React Build (port %-4s)                             ║\\n" "$FRONTEND_PORT"
    printf "║  ⚙️  Backend: API:%-4s DevPanel:%-4s Messaging:%-4s URLShort:%-4s        ║\\n" "$API_PORT" "$DEVPANEL_PORT" "$MESSAGING_PORT" "$URLSHORTENER_PORT"
    printf "║  🔧 Worker: %-4s (Code Stats, Scheduled Tasks)                              ║\\n" "$WORKER_PORT"
    echo "║  🗄️  Database: PostgreSQL                                                    ║"
    printf "║  📺 Sessions: tmux (%-20s, %-20s)    ║\\n" "$FRONTEND_SESSION" "$BACKEND_SESSION"
    echo "║  🎯 Mode: PRODUCTION (optimized, minified, no dev tools)                    ║"
    echo "╚══════════════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

log() {
    echo -e "${BLUE}[$(date +'%H:%M:%S')] $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

verbose() {
    if [ "$VERBOSE" = true ]; then
        echo -e "${PURPLE}[DEBUG] $1${NC}"
    fi
}

# =============================================================================
# Dependency Checks
# =============================================================================

check_dependencies() {
    log "Checking system dependencies..."
    
    # Check tmux
    if ! command -v tmux &> /dev/null; then
        error "tmux is required but not installed. Please install it first."
    fi
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        error "Node.js is required but not installed."
    fi
    
    # Check npm
    if ! command -v npm &> /dev/null; then
        error "npm is required but not installed."
    fi
    
    # Check Go
    if ! command -v go &> /dev/null; then
        error "Go is required but not installed."
    fi
    
    # Check tokei for code stats
    if ! command -v tokei &> /dev/null; then
        warning "tokei is not installed - code statistics will not be available"
        warning "Install with: cargo install tokei"
    else
        verbose "tokei found - code statistics available"
    fi
    
    # Check PostgreSQL connection (optional)
    if command -v psql &> /dev/null; then
        verbose "PostgreSQL client found"
    else
        warning "PostgreSQL client not found - database connectivity check will be skipped"
    fi
    
    # Check if pg_isready is available for better PostgreSQL checks
    if command -v pg_isready &> /dev/null; then
        verbose "pg_isready found - enhanced PostgreSQL checks available"
    fi
    
    success "All required dependencies found"
}

# =============================================================================
# Port Management
# =============================================================================

check_port() {
    local port=$1
    if lsof -i :$port >/dev/null 2>&1; then
        return 0  # Port is in use
    else
        return 1  # Port is free
    fi
}

find_available_port() {
    local start_port=$1
    local max_attempts=20
    local current_port=$start_port
    
    for ((i=0; i<max_attempts; i++)); do
        if ! check_port $current_port; then
            echo $current_port
            return 0
        fi
        ((current_port++))
    done
    
    error "Could not find available port starting from $start_port"
}

kill_port() {
    local port=$1
    local process_name=$2
    
    if check_port $port; then
        warning "Port $port is in use by $process_name"
        if [ "$KILL_EXISTING" = true ]; then
            log "Killing process on port $port..."
            lsof -ti:$port | xargs kill -9 2>/dev/null || true
            sleep 2
            success "Process on port $port killed"
        else
            return 1  # Port is in use, don't kill
        fi
    fi
    return 0  # Port is free or was killed
}

check_and_assign_ports() {
    log "Checking and assigning available ports..."
    
    # Check frontend port
    if check_port $FRONTEND_PORT; then
        warning "Port $FRONTEND_PORT is in use - will prompt to use another port"
        # Don't change the port, let serve handle it
    fi
    
    # Check backend ports
    if ! kill_port $API_PORT "API"; then
        local new_port=$(find_available_port $API_PORT)
        warning "API port $API_PORT in use, using port $new_port instead"
        API_PORT=$new_port
    fi
    
    if ! kill_port $DEVPANEL_PORT "DevPanel"; then
        local new_port=$(find_available_port $DEVPANEL_PORT)
        warning "DevPanel port $DEVPANEL_PORT in use, using port $new_port instead"
        DEVPANEL_PORT=$new_port
    fi
    
    if ! kill_port $MESSAGING_PORT "Messaging"; then
        local new_port=$(find_available_port $MESSAGING_PORT)
        warning "Messaging port $MESSAGING_PORT in use, using port $new_port instead"
        MESSAGING_PORT=$new_port
    fi
    
    if ! kill_port $URLSHORTENER_PORT "URL Shortener"; then
        local new_port=$(find_available_port $URLSHORTENER_PORT)
        warning "URL Shortener port $URLSHORTENER_PORT in use, using port $new_port instead"
        URLSHORTENER_PORT=$new_port
    fi
    
    if ! kill_port $WORKER_PORT "Worker"; then
        local new_port=$(find_available_port $WORKER_PORT)
        warning "Worker port $WORKER_PORT in use, using port $new_port instead"
        WORKER_PORT=$new_port
    fi
    
    success "All ports assigned successfully"
    log "Frontend: $FRONTEND_PORT, API: $API_PORT, DevPanel: $DEVPANEL_PORT, Messaging: $MESSAGING_PORT, URL Shortener: $URLSHORTENER_PORT, Worker: $WORKER_PORT"
}

# =============================================================================
# tmux Session Management
# =============================================================================

auto_detect_existing_sessions() {
    local sessions_found=false

    if tmux has-session -t "$FRONTEND_SESSION" 2>/dev/null; then
        warning "Detected existing tmux session: $FRONTEND_SESSION"
        sessions_found=true
    fi

    if tmux has-session -t "$BACKEND_SESSION" 2>/dev/null; then
        warning "Detected existing tmux session: $BACKEND_SESSION"
        sessions_found=true
    fi

    if [ "$sessions_found" = true ]; then
        log "Auto-enabling session cleanup to prevent conflicts..."
        KILL_EXISTING=true
    fi
}

kill_existing_sessions() {
    log "Cleaning up existing tmux sessions..."
    
    if tmux has-session -t "$FRONTEND_SESSION" 2>/dev/null; then
        tmux kill-session -t "$FRONTEND_SESSION"
        verbose "Killed existing frontend session"
    fi
    
    if tmux has-session -t "$BACKEND_SESSION" 2>/dev/null; then
        tmux kill-session -t "$BACKEND_SESSION"
        verbose "Killed existing backend session"
    fi
    
    success "tmux sessions cleaned up"
}

# =============================================================================
# Database Connectivity
# =============================================================================

test_database_connection() {
    log "Testing database connectivity..."

    # Environment variables are already loaded at script start
    if [ -f "$BACKEND_DIR/.env.production" ]; then
        verbose "Using production environment variables from backend/.env.production"
    elif [ -f "$BACKEND_DIR/.env" ]; then
        verbose "Using environment variables from backend/.env"
    fi

    # Try to connect to the database using environment variables or defaults
    local db_host="${DB_HOST:-127.0.0.1}"
    local db_port="${DB_PORT:-5432}"
    local db_name="${DB_NAME:-project_website}"
    local db_user="${DB_USER:-postgres}"
    local db_password="${DB_PASSWORD:-postgres}"

    # Show connection parameters for debugging
    verbose "Database connection parameters:"
    verbose "  Host: $db_host"
    verbose "  Port: $db_port"
    verbose "  Database: $db_name"
    verbose "  User: $db_user"
    verbose "  Password: $db_password"

    # Wait for PostgreSQL to be ready
    local max_attempts=3
    local attempt=0
    log "Waiting for database to become available..."
    while true; do
        psql_error=$(PGPASSWORD="$db_password" timeout 2 psql -h "$db_host" -p "$db_port" -U "$db_user" -d "$db_name" -c '\q' 2>&1)
        if [ $? -eq 0 ]; then
            break
        fi

        if [ $attempt -ge $max_attempts ]; then
            warning "Database connection failed after $max_attempts attempts."
            warning "Last error: $psql_error"
            return 1
        fi

        verbose "Attempting to connect to database (attempt $((attempt+1)) of $max_attempts)..."
        verbose "psql error: $psql_error"
        sleep 1
        attempt=$((attempt+1))
    done

    success "Database connection successful"
    return 0
}

run_database_migrations() {
    log "Running database migrations..."
    
    cd "$BACKEND_DIR"
    
    # Temporarily skip automatic migrations due to schema conflicts
    warning "Skipping automatic migrations - run manually if needed"
    
    # Ensure code_stats table exists
    if command -v psql &> /dev/null; then
        local db_host="${DB_HOST:-localhost}"
        local db_port="${DB_PORT:-5432}"
        local db_name="${DB_NAME:-project_website}"
        local db_user="${DB_USER:-postgres}"
        local db_password="${DB_PASSWORD:-postgres}"
        
        verbose "Ensuring code_stats table exists..."
        PGPASSWORD="$db_password" psql -h "$db_host" -p "$db_port" -d "$db_name" -U "$db_user" <<EOF 2>/dev/null || true
CREATE TABLE IF NOT EXISTS code_stats (
    id SERIAL PRIMARY KEY,
    lines_of_code INTEGER NOT NULL,
    file_count INTEGER NOT NULL,
    language_stats JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
EOF
    fi
    
    cd "$PROJECT_ROOT"
}

# =============================================================================
# Frontend Setup
# =============================================================================

setup_frontend() {
    log "Setting up frontend..."
    
    cd "$FRONTEND_DIR"
    
    # Install dependencies if needed
    if [ "$SKIP_DEPS" = false ]; then
        if [ ! -d "node_modules" ] || [ "$FORCE_FRESH" = true ]; then
            log "Installing frontend dependencies..."
            npm install
            success "Frontend dependencies installed"
        else
            verbose "Skipping frontend dependency installation"
        fi
    fi
    
    # Clear caches if requested
    if [ "$FORCE_FRESH" = true ]; then
        log "Clearing frontend caches..."
        npm run clear-cache 2>/dev/null || {
            log "Manual cache clearing..."
            rm -rf node_modules/.cache 2>/dev/null || true
            rm -rf build 2>/dev/null || true
            rm -f .eslintcache 2>/dev/null || true
            rm -f tsconfig.tsbuildinfo 2>/dev/null || true
        }
        success "Frontend caches cleared"
    fi
    
    # Update code statistics before building
    log "Updating code statistics with tokei..."
    if [ -x "$PROJECT_ROOT/scripts/monitoring/update_code_stats.sh" ]; then
        "$PROJECT_ROOT/scripts/monitoring/update_code_stats.sh" >/dev/null 2>&1
        success "Code statistics updated"
    else
        verbose "Code stats script not found or not executable"
    fi
    
    # Run production build with better memory management
    log "Creating optimized production build..."
    export NODE_OPTIONS="--max-old-space-size=4096"
    export NODE_ENV=production
    export GENERATE_SOURCEMAP=false
    npm run build:prod || npm run build
    
    if [ ! -d "build" ]; then
        error "Production build failed - build directory not created"
    fi
    
    # Verify build quality
    local build_size=$(du -sh build/ | cut -f1)
    success "Production build completed successfully (size: $build_size)"
    
    # Copy env-config.js to build directory
    if [ -f "public/env-config.js" ]; then
        cp public/env-config.js build/env-config.js
        verbose "Copied env-config.js to build directory"
    fi
    
    cd "$PROJECT_ROOT"
}

start_frontend() {
    log "Starting frontend production server in tmux session..."
    
    tmux new-session -d -s "$FRONTEND_SESSION" -c "$FRONTEND_DIR"
    
    # Set up the frontend session
    tmux send-keys -t "$FRONTEND_SESSION" "clear" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo 'Frontend Production Server Starting...'" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo 'Session: $FRONTEND_SESSION'" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo 'Port: $FRONTEND_PORT'" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo 'Mode: PRODUCTION (serving build/)'" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo ''" Enter
    
    # Set the port environment variable for the server
    tmux send-keys -t "$FRONTEND_SESSION" "export PORT=$FRONTEND_PORT" Enter
    
    # Start the production server
    tmux send-keys -t "$FRONTEND_SESSION" "npx serve -s build -l $FRONTEND_PORT --single" Enter
    
    success "Frontend session started: $FRONTEND_SESSION on port $FRONTEND_PORT"
}

# =============================================================================
# Production Configuration
# =============================================================================

create_production_config() {
    log "Creating production configuration..."
    
    cat > "$BACKEND_DIR/config/production.yaml" << EOF
# Production configuration
port: "8080"
jwtSecret: "${JWT_SECRET}"
environment: "production"
baseURL: "https://www.jadenrazo.dev"
logLevel: "info"
metricsEnabled: true
tracingEnabled: false
enablePprof: false

# Server timeouts (can be tighter or adjusted based on needs)
readTimeout: "10s"
writeTimeout: "15s"
idleTimeout: "60s"
shutdownTimeout: "15s"

# TLS settings (disabled since nginx handles TLS termination)
tlsEnabled: false
tlsCert: "${TLS_CERT_PATH:-/etc/letsencrypt/live/jadenrazo.dev/fullchain.pem}" # Path to your TLS certificate
tlsKey: "${TLS_KEY_PATH:-/etc/letsencrypt/live/jadenrazo.dev/privkey.pem}"    # Path to your TLS private key

# CORS allowed origins
allowedOrigins:
  - "https://jadenrazo.dev"
  - "https://www.jadenrazo.dev"

# Rate limiting (adjust based on expected traffic and capacity)
apiRateLimit: 200       # requests per minute
redirectRateLimit: 1000 # requests per minute

# Database configuration for PostgreSQL
database:
  driver: "postgres"
  host: "${DB_HOST}"    # Database host - set in .env
  port: 5432
  user: "${DB_USER}"    # Database username - set in .env
  password: "${DB_PASSWORD}" # Database password - set in .env
  dbName: "project_website"   # Production database name
  sslMode: "require" # Or "verify-full", "verify-ca". CRITICAL for production security

# URL Shortener Configuration
urlShortener:
  baseURL: "https://www.jadenrazo.dev"
  shortCodeLength: 7

# Messaging Configuration
messaging:
  maxMessageSize: 8192 # 8KB, adjust based on production needs

# DevPanel Configuration
devPanel:
  metricsInterval: "60s"
  maxLogLines: 5000
  logRetention: "720h" # 30 days
EOF
    
    success "Production configuration created"
}

# =============================================================================
# Backend Setup
# =============================================================================

setup_backend() {
    log "Setting up backend..."
    
    cd "$BACKEND_DIR"
    
    # Download Go modules if needed
    if [ "$SKIP_DEPS" = false ]; then
        if [ ! -d "vendor" ] || [ "$FORCE_FRESH" = true ]; then
            log "Downloading Go modules..."
            go mod download
            go mod tidy
            success "Go modules updated"
        else
            verbose "Skipping Go module download"
        fi
    fi
    
    cd "$PROJECT_ROOT"
}

start_backend() {
    log "Starting backend services in tmux session..."

    # Create a bin directory if it doesn't exist
    mkdir -p "$BACKEND_DIR/bin"

    # Build each service for production with optimizations
    log "Building backend services for production..."
    cd "$BACKEND_DIR"

    # Build API service
    if ! go build -ldflags="-s -w" -o "bin/api" "./cmd/api"; then
        error "Failed to build API service"
        return 1
    fi

    # Build DevPanel service
    if ! go build -ldflags="-s -w" -o "bin/devpanel" "./cmd/devpanel"; then
        error "Failed to build DevPanel service"
        return 1
    fi

    # Build Messaging service
    if ! go build -ldflags="-s -w" -o "bin/messaging" "./cmd/messaging"; then
        error "Failed to build Messaging service"
        return 1
    fi

    # Build URL Shortener service
    if ! go build -ldflags="-s -w" -o "bin/urlshortener" "./cmd/urlshortener"; then
        error "Failed to build URL Shortener service"
        return 1
    fi

    # Build Worker service
    if ! go build -ldflags="-s -w" -o "bin/worker" "./cmd/worker"; then
        error "Failed to build Worker service"
        return 1
    fi

    cd "$PROJECT_ROOT"
    success "Backend services built successfully"

    tmux new-session -d -s "$BACKEND_SESSION" -c "$PROJECT_ROOT"

    # Create windows for each service
    tmux rename-window -t "$BACKEND_SESSION:0" "api"
    tmux new-window -t "$BACKEND_SESSION" -n "devpanel" -c "$PROJECT_ROOT"
    tmux new-window -t "$BACKEND_SESSION" -n "messaging" -c "$PROJECT_ROOT"
    tmux new-window -t "$BACKEND_SESSION" -n "urlshortener" -c "$PROJECT_ROOT"
    tmux new-window -t "$BACKEND_SESSION" -n "worker" -c "$PROJECT_ROOT"

    # Start API service
    tmux send-keys -t "$BACKEND_SESSION:api" "clear" Enter
    tmux send-keys -t "$BACKEND_SESSION:api" "echo 'Starting API Service on port $API_PORT...'" Enter
    tmux send-keys -t "$BACKEND_SESSION:api" "cd $PROJECT_ROOT" Enter
    tmux send-keys -t "$BACKEND_SESSION:api" "set -a && source $BACKEND_DIR/.env.production && set +a" Enter
    tmux send-keys -t "$BACKEND_SESSION:api" "export API_PORT=$API_PORT" Enter
    tmux send-keys -t "$BACKEND_SESSION:api" "export NODE_ENV=production && export ENVIRONMENT=production" Enter
    tmux send-keys -t "$BACKEND_SESSION:api" "$BACKEND_DIR/bin/api" Enter

    # Start DevPanel service
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "clear" Enter
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "echo 'Starting DevPanel Service on port $DEVPANEL_PORT...'" Enter
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "cd $PROJECT_ROOT" Enter
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "set -a && source $BACKEND_DIR/.env.production && set +a" Enter
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "export DEVPANEL_PORT=$DEVPANEL_PORT" Enter
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "export NODE_ENV=production && export ENVIRONMENT=production" Enter
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "$BACKEND_DIR/bin/devpanel" Enter

    # Start Messaging service
    tmux send-keys -t "$BACKEND_SESSION:messaging" "clear" Enter
    tmux send-keys -t "$BACKEND_SESSION:messaging" "echo 'Starting Messaging Service on port $MESSAGING_PORT...'" Enter
    tmux send-keys -t "$BACKEND_SESSION:messaging" "cd $PROJECT_ROOT" Enter
    tmux send-keys -t "$BACKEND_SESSION:messaging" "set -a && source $BACKEND_DIR/.env.production && set +a" Enter
    tmux send-keys -t "$BACKEND_SESSION:messaging" "export MESSAGING_PORT=$MESSAGING_PORT" Enter
    tmux send-keys -t "$BACKEND_SESSION:messaging" "export NODE_ENV=production && export ENVIRONMENT=production" Enter
    tmux send-keys -t "$BACKEND_SESSION:messaging" "$BACKEND_DIR/bin/messaging" Enter

    # Start URL Shortener service
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "clear" Enter
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "echo 'Starting URL Shortener Service on port $URLSHORTENER_PORT...'" Enter
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "cd $PROJECT_ROOT" Enter
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "set -a && source $BACKEND_DIR/.env.production && set +a" Enter
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "export URLSHORTENER_PORT=$URLSHORTENER_PORT" Enter
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "export NODE_ENV=production && export ENVIRONMENT=production" Enter
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "$BACKEND_DIR/bin/urlshortener" Enter

    # Start Worker service
    tmux send-keys -t "$BACKEND_SESSION:worker" "clear" Enter
    tmux send-keys -t "$BACKEND_SESSION:worker" "echo 'Starting Worker Service on port $WORKER_PORT...'" Enter
    tmux send-keys -t "$BACKEND_SESSION:worker" "cd $PROJECT_ROOT" Enter
    tmux send-keys -t "$BACKEND_SESSION:worker" "set -a && source $BACKEND_DIR/.env.production && set +a" Enter
    tmux send-keys -t "$BACKEND_SESSION:worker" "export WORKER_PORT=$WORKER_PORT" Enter
    tmux send-keys -t "$BACKEND_SESSION:worker" "export NODE_ENV=production && export ENVIRONMENT=production" Enter
    tmux send-keys -t "$BACKEND_SESSION:worker" "$BACKEND_DIR/bin/worker" Enter

    # Focus on the API window
    tmux select-window -t "$BACKEND_SESSION:api"

    success "Backend session started: $BACKEND_SESSION with 5 microservices"
}

# =============================================================================
# Health Checks
# =============================================================================

wait_for_services() {
    log "Waiting for services to start..."
    
    local max_attempts=30
    local attempt=0
    
    # Wait for frontend
    while [ $attempt -lt $max_attempts ]; do
        if check_port $FRONTEND_PORT; then
            success "Frontend is running on http://localhost:$FRONTEND_PORT"
            break
        fi
        sleep 2
        ((attempt++))
    done
    
    # Wait for backend services
    sleep 5  # Give backend a bit more time
    
    local services=("$API_PORT:API" "$DEVPANEL_PORT:DevPanel" "$MESSAGING_PORT:Messaging" "$URLSHORTENER_PORT:URL Shortener" "$WORKER_PORT:Worker")
    
    for service in "${services[@]}"; do
        local port="${service%:*}"
        local name="${service#*:}"
        
        if check_port $port; then
            success "$name is running on http://localhost:$port"
        else
            warning "$name may not be running on port $port"
        fi
    done
}

# =============================================================================
# Usage and Help
# =============================================================================

show_help() {
    echo -e "${WHITE}Portfolio Website Production Launcher${NC}"
    echo ""
    echo -e "${CYAN}USAGE:${NC}"
    echo "  $0 [OPTIONS]"
    echo ""
    echo -e "${CYAN}OPTIONS:${NC}"
    echo "  -f, --fresh         Force fresh start with cache clearing and rebuild"
    echo "  -k, --kill-existing Kill existing processes on required ports (optional - auto-detected)"
    echo "  -s, --skip-deps     Skip dependency installation/updates"
    echo "  -v, --verbose       Enable verbose output"
    echo "  -h, --help          Show this help message"
    echo ""
    echo -e "${CYAN}AUTO-CLEANUP:${NC}"
    echo "  The script automatically detects existing tmux sessions and cleans them up"
    echo "  to prevent port conflicts. The --kill-existing flag is now optional."
    echo ""
    echo -e "${CYAN}EXAMPLES:${NC}"
    echo "  $0                  # Normal startup (auto-cleans existing sessions)"
    echo "  $0 --fresh          # Fresh start with cache clearing and rebuild"
    echo "  $0 -f -k            # Fresh start with explicit session cleanup"
    echo "  $0 --verbose        # Verbose output showing auto-detection"
    echo ""
    echo -e "${CYAN}TMUX SESSIONS:${NC}"
    echo "  Frontend: tmux attach -t $FRONTEND_SESSION"
    echo "  Backend:  tmux attach -t $BACKEND_SESSION"
    echo ""
    echo -e "${CYAN}STOPPING:${NC}"
    echo "  tmux kill-session -t $FRONTEND_SESSION"
    echo "  tmux kill-session -t $BACKEND_SESSION"
    echo ""
    echo -e "${CYAN}DIFFERENCE FROM DEV MODE:${NC}"
    echo "  • Uses optimized, minified React build"
    echo "  • No hot reload or development WebSocket"
    echo "  • Smaller bundle sizes and faster loading"
    echo "  • Production environment variables"
    echo "  • No source maps or dev tools"
    echo ""
}

show_completion_summary() {
    echo ""
    echo -e "${GREEN}🎉 Production environment started successfully!${NC}"
    echo ""
    echo -e "${CYAN}📋 Service URLs:${NC}"
    echo -e "  Main App:      ${WHITE}http://localhost:$API_PORT${NC} (Frontend + API)"
    echo -e "  Frontend Only: ${WHITE}http://localhost:$FRONTEND_PORT${NC} (Static files only, no API proxy)"
    echo -e "  DevPanel:      ${WHITE}http://localhost:$DEVPANEL_PORT${NC}"
    echo -e "  Messaging:     ${WHITE}http://localhost:$MESSAGING_PORT${NC}"
    echo -e "  URL Shortener: ${WHITE}http://localhost:$URLSHORTENER_PORT${NC}"
    echo -e "  Worker:        ${WHITE}http://localhost:$WORKER_PORT${NC} (Background tasks)"
    echo ""
    echo -e "${CYAN}📊 Key Features:${NC}"
    echo -e "  Status Page:   ${WHITE}http://localhost:$API_PORT/status${NC}"
    echo -e "  Code Stats:    ${WHITE}http://localhost:$API_PORT/api/v1/code/stats${NC}"
    echo -e "  Health Check:  ${WHITE}http://localhost:$API_PORT/health${NC}"
    echo ""
    echo -e "${YELLOW}⚠️  IMPORTANT:${NC}"
    echo -e "  For full functionality, access the app at ${WHITE}http://localhost:$API_PORT${NC}"
    echo -e "  The frontend-only server on port $FRONTEND_PORT does not proxy API requests."
    echo ""
    echo -e "${CYAN}🖥️  Production Sessions:${NC}"
    echo -e "  Frontend: ${WHITE}tmux attach -t $FRONTEND_SESSION${NC}"
    echo -e "  Backend:  ${WHITE}tmux attach -t $BACKEND_SESSION${NC}"
    echo ""
    echo -e "${CYAN}🛑 To Stop Production Services:${NC}"
    echo -e "  ${WHITE}tmux kill-session -t $FRONTEND_SESSION${NC}"
    echo -e "  ${WHITE}tmux kill-session -t $BACKEND_SESSION${NC}"
    echo ""
    echo -e "${YELLOW}💡 Tip: Use Ctrl+B then D to detach from tmux sessions${NC}"
    echo ""
    echo -e "${GREEN}✨ This is a production-optimized deployment!${NC}"
    echo ""
}

# =============================================================================
# Main Execution
# =============================================================================

main() {
    # Auto-detect existing tmux sessions and enable cleanup if found
    auto_detect_existing_sessions

    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -f|--fresh)
                FORCE_FRESH=true
                shift
                ;;
            -k|--kill-existing)
                KILL_EXISTING=true
                shift
                ;;
            -s|--skip-deps)
                SKIP_DEPS=true
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                error "Unknown option: $1"
                ;;
        esac
    done
    
    # Main execution flow
    print_banner
    
    check_dependencies
    check_and_assign_ports
    kill_existing_sessions
    # Skip database check if not available
    test_database_connection || warning "Database connection failed, continuing anyway..."
    # run_database_migrations
    
    # Show updated banner with final ports
    echo ""
    print_updated_banner
    echo ""
    
    setup_frontend
    create_production_config
    setup_backend
    
    start_backend
    sleep 2
    start_frontend
    
    wait_for_services
    show_completion_summary
}

# Handle script interruption
trap 'echo -e "\n${YELLOW}Production startup interrupted. Cleaning up...${NC}"; kill_existing_sessions; exit 1' INT

# Run main function
main "$@"