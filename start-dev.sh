#!/bin/bash

# =============================================================================
# Project Website Development Startup Script
# =============================================================================
# This script starts both frontend and backend in separate tmux sessions
# with fresh cache clearing and database connectivity checks
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
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="$PROJECT_ROOT/frontend"
BACKEND_DIR="$PROJECT_ROOT/backend"
SESSION_PREFIX="portfolio"
FRONTEND_SESSION="${SESSION_PREFIX}-frontend"
BACKEND_SESSION="${SESSION_PREFIX}-backend"

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

# =============================================================================
# Helper Functions
# =============================================================================

print_banner() {
    echo -e "${CYAN}"
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                    Portfolio Website Development Launcher                   ║"
    echo "║                                                                              ║"
    echo "║  🚀 Frontend: React/TypeScript (port $FRONTEND_PORT)                                    ║"
    echo "║  ⚙️  Backend: Go microservices (ports $API_PORT-$WORKER_PORT)                               ║"
    echo "║  🗄️  Database: PostgreSQL                                                    ║"
    echo "║  📺 Sessions: tmux ($FRONTEND_SESSION, $BACKEND_SESSION)                      ║"
    echo "║  📊 Features: Status Monitoring, Code Stats (tokei), URL Shortener              ║"
    echo "╚══════════════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_updated_banner() {
    echo -e "${CYAN}"
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                    Portfolio Website Development Launcher                   ║"
    echo "║                                                                              ║"
    printf "║  🚀 Frontend: React/TypeScript (port %-4s)                                  ║\n" "$FRONTEND_PORT"
    printf "║  ⚙️  Backend: API:%-4s DevPanel:%-4s Messaging:%-4s URLShort:%-4s        ║\n" "$API_PORT" "$DEVPANEL_PORT" "$MESSAGING_PORT" "$URLSHORTENER_PORT"
    printf "║  🔧 Worker: %-4s (Code Stats, Scheduled Tasks)                              ║\n" "$WORKER_PORT"
    echo "║  🗄️  Database: PostgreSQL                                                    ║"
    printf "║  📺 Sessions: tmux (%-20s, %-20s)    ║\n" "$FRONTEND_SESSION" "$BACKEND_SESSION"
    echo "║  📊 Features: Status Monitoring, Code Stats (tokei), URL Shortener              ║"
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
        # Don't change the port, let React handle it
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
    
    # Load environment variables from backend/.env if it exists
    if [ -f "$BACKEND_DIR/.env" ]; then
        verbose "Loading environment variables from backend/.env"
        export $(grep -v '^#' "$BACKEND_DIR/.env" | xargs)
    fi
    
    # Try to connect to the database using environment variables or defaults
    local db_host="${DB_HOST:-195.201.136.53}"
    local db_port="${DB_PORT:-5432}"
    local db_name="${DB_NAME:-project_website}"
    local db_user="${DB_USER:-postgres}"
    local db_password="${DB_PASSWORD:-postgres}"
    
    # First check if PostgreSQL is running
    if command -v pg_isready &> /dev/null; then
        if ! pg_isready -h "$db_host" -p "$db_port" >/dev/null 2>&1; then
            warning "PostgreSQL is not running on $db_host:$db_port"
            warning "Please start PostgreSQL first:"
            warning "  • macOS: brew services start postgresql"
            warning "  • Ubuntu: sudo systemctl start postgresql"
            warning "  • Docker: docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres"
            return 1
        fi
        verbose "PostgreSQL is running on $db_host:$db_port"
    fi
    
    # Test the actual connection
    if command -v psql &> /dev/null; then
        # Try to connect to the database
        if PGPASSWORD="$db_password" psql -h "$db_host" -p "$db_port" -d "$db_name" -U "$db_user" -c "SELECT 1;" >/dev/null 2>&1; then
            success "Database connection successful"
            return 0
        else
            # Database might not exist, try to create it
            warning "Database '$db_name' might not exist, attempting to create it..."
            if PGPASSWORD="$db_password" createdb -h "$db_host" -p "$db_port" -U "$db_user" "$db_name" 2>/dev/null; then
                success "Database '$db_name' created successfully"
                return 0
            else
                warning "Could not create database - it may already exist or you may need different permissions"
            fi
        fi
    else
        warning "Cannot test database connection - psql not available"
    fi
    
    return 1
}

# =============================================================================
# Database Setup and Migrations
# =============================================================================

run_database_migrations() {
    log "Running database migrations..."
    
    cd "$BACKEND_DIR"
    
    # Temporarily skip automatic migrations due to schema conflicts
    warning "Skipping automatic migrations - run manually if needed"
    # TODO: Fix domain.User struct to match existing database schema
    
    # Check if migration command exists
    # if [ -f "cmd/migration/main.go" ]; then
    #     verbose "Running database migrations..."
    #     if go run cmd/migration/main.go up 2>/dev/null; then
    #         success "Database migrations completed"
    #     else
    #         warning "Database migration failed - this might be expected if migrations are already applied"
    #     fi
    # else
    #     verbose "No migration command found, skipping migrations"
    # fi
    
    # Ensure code_stats table exists
    if command -v psql &> /dev/null; then
        local db_host="${DB_HOST:-195.201.136.53}"
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
    
    # Update code statistics before starting
    log "Updating code statistics with tokei..."
    if [ -x "$PROJECT_ROOT/scripts/update_code_stats.sh" ]; then
        "$PROJECT_ROOT/scripts/update_code_stats.sh" >/dev/null 2>&1
        success "Code statistics updated"
    else
        verbose "Code stats script not found or not executable"
    fi
    
    cd "$PROJECT_ROOT"
}

start_frontend() {
    log "Starting frontend in tmux session..."
    
    tmux new-session -d -s "$FRONTEND_SESSION" -c "$FRONTEND_DIR"
    
    # Set up the frontend session
    tmux send-keys -t "$FRONTEND_SESSION" "clear" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo 'Frontend Development Server Starting...'" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo 'Session: $FRONTEND_SESSION'" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo 'Port: $FRONTEND_PORT'" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo ''" Enter
    
    # Set the port environment variable for React
    tmux send-keys -t "$FRONTEND_SESSION" "export PORT=$FRONTEND_PORT" Enter
    
    # Start the frontend server
    if [ "$FORCE_FRESH" = true ]; then
        tmux send-keys -t "$FRONTEND_SESSION" "npm run dev:fresh" Enter
    else
        tmux send-keys -t "$FRONTEND_SESSION" "npm start" Enter
    fi
    
    # Send 'y' after a short delay in case there's a port prompt
    # This handles the "Would you like to run the app on another port instead?" prompt
    (
        sleep 3
        tmux send-keys -t "$FRONTEND_SESSION" "y" Enter 2>/dev/null || true
    ) &
    
    success "Frontend session started: $FRONTEND_SESSION on port $FRONTEND_PORT"
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
    
    tmux new-session -d -s "$BACKEND_SESSION" -c "$BACKEND_DIR"
    
    # Create windows for each service
    tmux rename-window -t "$BACKEND_SESSION:0" "api"
    tmux new-window -t "$BACKEND_SESSION" -n "devpanel" -c "$BACKEND_DIR"
    tmux new-window -t "$BACKEND_SESSION" -n "messaging" -c "$BACKEND_DIR"
    tmux new-window -t "$BACKEND_SESSION" -n "urlshortener" -c "$BACKEND_DIR"
    tmux new-window -t "$BACKEND_SESSION" -n "worker" -c "$BACKEND_DIR"
    
    # Start API service (main service)
    tmux send-keys -t "$BACKEND_SESSION:api" "clear" Enter
    tmux send-keys -t "$BACKEND_SESSION:api" "echo 'Starting API Service on port $API_PORT...'" Enter
    tmux send-keys -t "$BACKEND_SESSION:api" "cd $BACKEND_DIR" Enter
    tmux send-keys -t "$BACKEND_SESSION:api" "export API_PORT=$API_PORT" Enter
    tmux send-keys -t "$BACKEND_SESSION:api" "go run cmd/api/main.go" Enter
    
    # Start DevPanel service
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "clear" Enter
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "echo 'Starting DevPanel Service on port $DEVPANEL_PORT...'" Enter
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "cd $BACKEND_DIR" Enter
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "export DEVPANEL_PORT=$DEVPANEL_PORT" Enter
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "go run cmd/devpanel/main.go" Enter
    
    # Start Messaging service
    tmux send-keys -t "$BACKEND_SESSION:messaging" "clear" Enter
    tmux send-keys -t "$BACKEND_SESSION:messaging" "echo 'Starting Messaging Service on port $MESSAGING_PORT...'" Enter
    tmux send-keys -t "$BACKEND_SESSION:messaging" "cd $BACKEND_DIR" Enter
    tmux send-keys -t "$BACKEND_SESSION:messaging" "export MESSAGING_PORT=$MESSAGING_PORT" Enter
    tmux send-keys -t "$BACKEND_SESSION:messaging" "go run cmd/messaging/main.go" Enter
    
    # Start URL Shortener service
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "clear" Enter
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "echo 'Starting URL Shortener Service on port $URLSHORTENER_PORT...'" Enter
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "cd $BACKEND_DIR" Enter
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "export URLSHORTENER_PORT=$URLSHORTENER_PORT" Enter
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "go run cmd/urlshortener/main.go" Enter
    
    # Start Worker service (includes scheduled tasks for code stats)
    tmux send-keys -t "$BACKEND_SESSION:worker" "clear" Enter
    tmux send-keys -t "$BACKEND_SESSION:worker" "echo 'Starting Worker Service (Code Stats, Scheduled Tasks)...'" Enter
    tmux send-keys -t "$BACKEND_SESSION:worker" "cd $BACKEND_DIR" Enter
    tmux send-keys -t "$BACKEND_SESSION:worker" "export WORKER_PORT=$WORKER_PORT" Enter
    tmux send-keys -t "$BACKEND_SESSION:worker" "go run cmd/worker/main.go" Enter
    
    # Focus on the API window
    tmux select-window -t "$BACKEND_SESSION:api"
    
    success "Backend session started: $BACKEND_SESSION"
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
    echo -e "${WHITE}Portfolio Website Development Launcher${NC}"
    echo ""
    echo -e "${CYAN}USAGE:${NC}"
    echo "  $0 [OPTIONS]"
    echo ""
    echo -e "${CYAN}OPTIONS:${NC}"
    echo "  -f, --fresh         Force fresh start with cache clearing"
    echo "  -k, --kill-existing Kill existing processes on required ports"
    echo "  -s, --skip-deps     Skip dependency installation/updates"
    echo "  -v, --verbose       Enable verbose output"
    echo "  -h, --help          Show this help message"
    echo ""
    echo -e "${CYAN}EXAMPLES:${NC}"
    echo "  $0                  # Normal startup"
    echo "  $0 --fresh          # Fresh start with cache clearing"
    echo "  $0 -f -k            # Fresh start and kill existing processes"
    echo "  $0 --verbose        # Verbose output"
    echo ""
    echo -e "${CYAN}TMUX SESSIONS:${NC}"
    echo "  Frontend: tmux attach -t $FRONTEND_SESSION"
    echo "  Backend:  tmux attach -t $BACKEND_SESSION"
    echo ""
    echo -e "${CYAN}STOPPING:${NC}"
    echo "  tmux kill-session -t $FRONTEND_SESSION"
    echo "  tmux kill-session -t $BACKEND_SESSION"
    echo ""
}

show_completion_summary() {
    echo ""
    echo -e "${GREEN}🎉 Development environment started successfully!${NC}"
    echo ""
    echo -e "${CYAN}📋 Service URLs:${NC}"
    echo -e "  Frontend:      ${WHITE}http://localhost:$FRONTEND_PORT${NC}"
    echo -e "  API:           ${WHITE}http://localhost:$API_PORT${NC}"
    echo -e "  DevPanel:      ${WHITE}http://localhost:$DEVPANEL_PORT${NC}"
    echo -e "  Messaging:     ${WHITE}http://localhost:$MESSAGING_PORT${NC}"
    echo -e "  URL Shortener: ${WHITE}http://localhost:$URLSHORTENER_PORT${NC}"
    echo -e "  Worker:        ${WHITE}http://localhost:$WORKER_PORT${NC} (Background tasks)"
    echo ""
    echo -e "${CYAN}📊 Key Features:${NC}"
    echo -e "  Status Page:   ${WHITE}http://localhost:$FRONTEND_PORT/status${NC}"
    echo -e "  Code Stats:    ${WHITE}http://localhost:$API_PORT/api/code-stats${NC}"
    echo -e "  Health Check:  ${WHITE}http://localhost:$API_PORT/api/status${NC}"
    echo ""
    echo -e "${CYAN}🖥️  tmux Sessions:${NC}"
    echo -e "  Frontend: ${WHITE}tmux attach -t $FRONTEND_SESSION${NC}"
    echo -e "  Backend:  ${WHITE}tmux attach -t $BACKEND_SESSION${NC}"
    echo ""
    echo -e "${CYAN}🛑 To Stop All Services:${NC}"
    echo -e "  ${WHITE}tmux kill-session -t $FRONTEND_SESSION${NC}"
    echo -e "  ${WHITE}tmux kill-session -t $BACKEND_SESSION${NC}"
    echo ""
    echo -e "${YELLOW}💡 Tip: Use Ctrl+B then D to detach from tmux sessions${NC}"
    echo ""
}

# =============================================================================
# Main Execution
# =============================================================================

main() {
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
    test_database_connection
    run_database_migrations
    
    # Show updated banner with final ports
    echo ""
    print_updated_banner
    echo ""
    
    setup_frontend
    setup_backend
    
    start_backend
    sleep 2
    start_frontend
    
    wait_for_services
    show_completion_summary
}

# Handle script interruption
trap 'echo -e "\n${YELLOW}Startup interrupted. Cleaning up...${NC}"; kill_existing_sessions; exit 1' INT

# Run main function
main "$@"