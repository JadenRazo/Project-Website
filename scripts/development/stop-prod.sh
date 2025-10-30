#!/bin/bash

# =============================================================================
# Portfolio Website Production Stop Script
# =============================================================================
# This script stops all production services and cleans up tmux sessions
# =============================================================================

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Configuration
SESSION_PREFIX="portfolio"
FRONTEND_SESSION="${SESSION_PREFIX}-frontend-prod"
BACKEND_SESSION="${SESSION_PREFIX}-backend-prod"

# Ports to clean up (extended range to catch any port variations)
PORTS=(3000 3001 3002 3003 3004 3005 8080 8081 8082 8083 8084 8085 8086 8087 8088 8089)

print_banner() {
    echo -e "${CYAN}"
    echo "╔══════════════════════════════════════════════════════════════════════════════╗"
    echo "║                    Portfolio Website Production Stopper                     ║"
    echo "║                                                                              ║"
    echo "║  🛑 Stopping all production services...                                     ║"
    echo "║  🧹 Cleaning up tmux sessions...                                            ║"
    echo "║  🔌 Freeing up ports...                                                     ║"
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
}

stop_tmux_sessions() {
    log "Stopping tmux sessions..."
    
    if tmux has-session -t "$FRONTEND_SESSION" 2>/dev/null; then
        tmux kill-session -t "$FRONTEND_SESSION"
        success "Frontend session stopped: $FRONTEND_SESSION"
    else
        warning "Frontend session not found: $FRONTEND_SESSION"
    fi
    
    if tmux has-session -t "$BACKEND_SESSION" 2>/dev/null; then
        tmux kill-session -t "$BACKEND_SESSION"
        success "Backend session stopped: $BACKEND_SESSION"
    else
        warning "Backend session not found: $BACKEND_SESSION"
    fi
}

cleanup_ports() {
    log "Cleaning up ports..."
    
    for port in "${PORTS[@]}"; do
        if lsof -i :$port >/dev/null 2>&1; then
            log "Killing processes on port $port..."
            lsof -ti:$port | xargs kill -9 2>/dev/null || true
            success "Port $port cleaned up"
        else
            success "Port $port is already free"
        fi
    done
}

show_completion() {
    echo ""
    echo -e "${GREEN}🎉 All production services stopped successfully!${NC}"
    echo ""
    echo -e "${CYAN}📋 What was stopped:${NC}"
    echo -e "  • Frontend production server"
    echo -e "  • Backend API service"
    echo -e "  • DevPanel service"
    echo -e "  • Messaging service"
    echo -e "  • URL Shortener service"
    echo -e "  • tmux sessions: $FRONTEND_SESSION, $BACKEND_SESSION"
    echo ""
    echo -e "${YELLOW}💡 To start again, run: ./start-prod.sh${NC}"
    echo ""
}

main() {
    print_banner
    stop_tmux_sessions
    sleep 2
    cleanup_ports
    show_completion
}

# Run main function
main "$@"
