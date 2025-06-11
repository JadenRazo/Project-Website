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
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="$PROJECT_ROOT/frontend"
BACKEND_DIR="$PROJECT_ROOT/backend"
SESSION_PREFIX="portfolio"
FRONTEND_SESSION="${SESSION_PREFIX}-frontend-prod"
BACKEND_SESSION="${SESSION_PREFIX}-backend-prod"

# Production ports
FRONTEND_PORT=3000
API_PORT=8080
DEVPANEL_PORT=8081
MESSAGING_PORT=8082
URLSHORTENER_PORT=8083
WORKER_PORT=8084

# Script options
FORCE_REBUILD=false
SKIP_DEPS=false
VERBOSE=false
KILL_EXISTING=false

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
# Dependency and Environment Checks
# =============================================================================

check_dependencies() {
    log "Checking system dependencies..."
    
    # Check Node.js version
    if ! command -v node >/dev/null 2>&1; then
        error "Node.js not found. Please install Node.js 18+ and npm."
    fi
    
    local node_version=$(node -v | sed 's/v//')
    local node_major=$(echo $node_version | cut -d. -f1)
    
    if [ "$node_major" -lt 18 ]; then
        error "Node.js version $node_version is too old. Please upgrade to Node.js 18+."
    fi
    
    verbose "Node.js version: $node_version ✓"
    
    # Check npm
    if ! command -v npm >/dev/null 2>&1; then
        error "npm not found. Please install npm."
    fi
    
    verbose "npm version: $(npm -v) ✓"
    
    # Check Go
    if ! command -v go >/dev/null 2>&1; then
        error "Go not found. Please install Go 1.19+ for backend services."
    fi
    
    local go_version=$(go version | awk '{print $3}' | sed 's/go//')
    verbose "Go version: $go_version ✓"
    
    # Check tmux
    if ! command -v tmux >/dev/null 2>&1; then
        error "tmux not found. Please install tmux for session management."
    fi
    
    verbose "tmux version: $(tmux -V) ✓"
    
    success "All system dependencies verified"
}

check_old_packages() {
    log "Checking for problematic package versions..."
    
    cd "$FRONTEND_DIR"
    
    # Check for known problematic packages
    if npm ls react-scripts@5.0.1 >/dev/null 2>&1; then
        warning "react-scripts 5.0.1 detected - this may cause build issues"
        log "Consider upgrading to a more recent version"
    fi
    
    # Check TypeScript version compatibility
    local ts_version=$(npm ls typescript --depth=0 2>/dev/null | grep typescript@ | sed 's/.*typescript@//' | cut -d' ' -f1)
    if [[ "$ts_version" =~ ^4\. ]]; then
        warning "TypeScript 4.x detected. Some dependencies expect TypeScript 5+"
    fi
    
    cd "$PROJECT_ROOT"
    success "Package version check completed"
}

# =============================================================================
# Production Build Functions
# =============================================================================

build_frontend() {
    log "Building frontend for production..."
    
    cd "$FRONTEND_DIR"
    
    # Install/update dependencies if needed
    if [ "$SKIP_DEPS" = false ]; then
        if [ ! -d "node_modules" ] || [ "$FORCE_REBUILD" = true ]; then
            log "Cleaning npm cache and installing dependencies..."
            npm cache clean --force
            rm -rf node_modules package-lock.json
            npm install
            success "Frontend dependencies installed"
        fi
    fi
    
    # Clear previous build
    if [ "$FORCE_REBUILD" = true ]; then
        log "Clearing previous build and caches..."
        rm -rf build/
        rm -rf node_modules/.cache/ || true
        npm cache clean --force
        success "Previous build and caches cleared"
    fi
    
    # Ensure we have the latest package-lock.json
    if [ ! -f "package-lock.json" ]; then
        log "Generating package-lock.json..."
        npm install --package-lock-only
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

check_production_readiness() {
    log "Checking production readiness..."
    
    # Check if build directory exists
    if [ ! -d "$FRONTEND_DIR/build" ]; then
        error "No production build found. Run with --rebuild or create build first."
    fi
    
    # Check if build is recent (within last hour for safety)
    local build_age=$(find "$FRONTEND_DIR/build" -maxdepth 0 -mmin -60 2>/dev/null | wc -l)
    if [ "$build_age" -eq 0 ] && [ "$FORCE_REBUILD" = false ]; then
        warning "Build is older than 1 hour. Consider rebuilding with --rebuild"
    fi
    
    # Check essential files
    local essential_files=("$FRONTEND_DIR/build/index.html" "$FRONTEND_DIR/build/static")
    for file in "${essential_files[@]}"; do
        if [ ! -e "$file" ]; then
            error "Essential build file missing: $file"
        fi
    done
    
    success "Production build verification passed"
}

# =============================================================================
# Session Management
# =============================================================================

kill_existing_sessions() {
    log "Cleaning up existing production sessions..."
    
    if tmux has-session -t "$FRONTEND_SESSION" 2>/dev/null; then
        tmux kill-session -t "$FRONTEND_SESSION"
        verbose "Killed existing frontend production session"
    fi
    
    if tmux has-session -t "$BACKEND_SESSION" 2>/dev/null; then
        tmux kill-session -t "$BACKEND_SESSION"
        verbose "Killed existing backend production session"
    fi
    
    success "Production sessions cleaned up"
}

start_frontend_production() {
    log "Starting frontend production server..."
    
    cd "$FRONTEND_DIR"
    
    # Verify serve is available
    if ! command -v npx >/dev/null 2>&1; then
        error "npx not found. Please ensure Node.js is properly installed."
    fi
    
    tmux new-session -d -s "$FRONTEND_SESSION" -c "$FRONTEND_DIR"
    
    # Set up the frontend production session
    tmux send-keys -t "$FRONTEND_SESSION" "clear" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo 'Frontend Production Server Starting...'" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo 'Session: $FRONTEND_SESSION'" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo 'Port: $FRONTEND_PORT'" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo 'Mode: PRODUCTION (serving build/)'" Enter
    tmux send-keys -t "$FRONTEND_SESSION" "echo ''" Enter
    
    # Use the project's serve installation and add fallback
    tmux send-keys -t "$FRONTEND_SESSION" "npx serve -s build -l $FRONTEND_PORT --single || npm run serve" Enter
    
    success "Frontend production session started: $FRONTEND_SESSION on port $FRONTEND_PORT"
    
    cd "$PROJECT_ROOT"
}

start_backend_production() {
    log "Starting backend services for production..."
    
    cd "$BACKEND_DIR"
    
    # Set production environment
    export NODE_ENV=production
    export GO_ENV=production
    
    tmux new-session -d -s "$BACKEND_SESSION" -c "$BACKEND_DIR"
    
    # Create windows for each service
    tmux rename-window -t "$BACKEND_SESSION:0" "api"
    tmux new-window -t "$BACKEND_SESSION" -n "devpanel" -c "$BACKEND_DIR"
    tmux new-window -t "$BACKEND_SESSION" -n "messaging" -c "$BACKEND_DIR"
    tmux new-window -t "$BACKEND_SESSION" -n "urlshortener" -c "$BACKEND_DIR"
    tmux new-window -t "$BACKEND_SESSION" -n "worker" -c "$BACKEND_DIR"
    
    # Start services with production settings
    tmux send-keys -t "$BACKEND_SESSION:api" "export NODE_ENV=production && export GO_ENV=production" Enter
    tmux send-keys -t "$BACKEND_SESSION:api" "go run cmd/api/main.go" Enter
    
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "export NODE_ENV=production && export GO_ENV=production" Enter
    tmux send-keys -t "$BACKEND_SESSION:devpanel" "go run cmd/devpanel/main.go" Enter
    
    tmux send-keys -t "$BACKEND_SESSION:messaging" "export NODE_ENV=production && export GO_ENV=production" Enter
    tmux send-keys -t "$BACKEND_SESSION:messaging" "go run cmd/messaging/main.go" Enter
    
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "export NODE_ENV=production && export GO_ENV=production" Enter
    tmux send-keys -t "$BACKEND_SESSION:urlshortener" "go run cmd/urlshortener/main.go" Enter
    
    tmux send-keys -t "$BACKEND_SESSION:worker" "export NODE_ENV=production && export GO_ENV=production" Enter
    tmux send-keys -t "$BACKEND_SESSION:worker" "go run cmd/worker/main.go" Enter
    
    tmux select-window -t "$BACKEND_SESSION:api"
    
    success "Backend production session started: $BACKEND_SESSION"
    
    cd "$PROJECT_ROOT"
}

# =============================================================================
# Main Functions
# =============================================================================

show_help() {
    echo -e "${WHITE}Portfolio Website Production Launcher${NC}"
    echo ""
    echo -e "${CYAN}USAGE:${NC}"
    echo "  $0 [OPTIONS]"
    echo ""
    echo -e "${CYAN}OPTIONS:${NC}"
    echo "  -r, --rebuild       Force rebuild of frontend production build"
    echo "  -k, --kill-existing Kill existing production sessions"
    echo "  -s, --skip-deps     Skip dependency installation/updates"
    echo "  -v, --verbose       Enable verbose output"
    echo "  -h, --help          Show this help message"
    echo ""
    echo -e "${CYAN}EXAMPLES:${NC}"
    echo "  $0                  # Start with existing build"
    echo "  $0 --rebuild        # Force rebuild and start"
    echo "  $0 -r -k            # Rebuild and kill existing sessions"
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
    echo -e "${CYAN}📋 Production URLs:${NC}"
    echo -e "  Frontend:      ${WHITE}http://localhost:$FRONTEND_PORT${NC} (production build)"
    echo -e "  API:           ${WHITE}http://localhost:$API_PORT${NC}"
    echo -e "  DevPanel:      ${WHITE}http://localhost:$DEVPANEL_PORT${NC}"
    echo -e "  Status:        ${WHITE}http://localhost:$FRONTEND_PORT/status${NC}"
    echo ""
    echo -e "${CYAN}🖥️  Production Sessions:${NC}"
    echo -e "  Frontend: ${WHITE}tmux attach -t $FRONTEND_SESSION${NC}"
    echo -e "  Backend:  ${WHITE}tmux attach -t $BACKEND_SESSION${NC}"
    echo ""
    echo -e "${CYAN}🛑 To Stop Production Services:${NC}"
    echo -e "  ${WHITE}tmux kill-session -t $FRONTEND_SESSION${NC}"
    echo -e "  ${WHITE}tmux kill-session -t $BACKEND_SESSION${NC}"
    echo ""
    echo -e "${GREEN}✨ This is a production-optimized deployment suitable for recruiters!${NC}"
    echo ""
}

main() {
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -r|--rebuild)
                FORCE_REBUILD=true
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
    
    print_banner
    
    # Run dependency and environment checks
    check_dependencies
    check_old_packages
    
    if [ "$KILL_EXISTING" = true ]; then
        kill_existing_sessions
    fi
    
    if [ "$FORCE_REBUILD" = true ]; then
        build_frontend
    else
        check_production_readiness
    fi
    
    start_backend_production
    sleep 3
    start_frontend_production
    
    sleep 5
    show_completion_summary
}

# Handle script interruption
trap 'echo -e "\n${YELLOW}Production startup interrupted. Cleaning up...${NC}"; kill_existing_sessions; exit 1' INT

# Run main function
main "$@"