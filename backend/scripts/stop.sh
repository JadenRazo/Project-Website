#!/bin/bash

# Print colorful messages
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print a header
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${GREEN}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_header "Stopping Project-Website Services"

# Define all services
SERVICES=("api" "devpanel" "messaging" "urlshortener" "frontend")

# Check if any tmux sessions are running
if ! command -v tmux &> /dev/null; then
    echo -e "${RED}Error: tmux is not installed. Cannot stop services.${NC}"
    exit 1
fi

# Function to gracefully stop a service
stop_service() {
    local service=$1
    
    if tmux has-session -t "$service" 2>/dev/null; then
        echo -e "${YELLOW}Stopping ${service} service...${NC}"
        
        # Try to send SIGTERM to the tmux session
        tmux send-keys -t "$service" C-c
        
        # Wait a moment for the service to shut down gracefully
        sleep 2
        
        # Kill the tmux session
        tmux kill-session -t "$service" 2>/dev/null
        
        echo -e "${GREEN}✓ ${service} service stopped${NC}"
    else
        echo -e "${CYAN}ℹ ${service} service is not running${NC}"
    fi
}

# Stop each service
for service in "${SERVICES[@]}"; do
    stop_service "$service"
done

# Double-check no sessions are left
remaining_sessions=$(tmux list-sessions -F "#{session_name}" 2>/dev/null | grep -E "^(api|devpanel|messaging|urlshortener|frontend)$" || echo "")

if [ -n "$remaining_sessions" ]; then
    echo -e "${YELLOW}Some sessions are still running. Forcing termination...${NC}"
    echo "$remaining_sessions" | while read -r session; do
        echo -e "${RED}Killing session: ${session}${NC}"
        tmux kill-session -t "$session" 2>/dev/null
    done
fi

# Check for any remaining server processes
echo -e "${CYAN}Checking for remaining processes...${NC}"

# Check for Node.js frontend processes
node_procs=$(ps aux | grep -E "node.*start" | grep -v grep || echo "")
if [ -n "$node_procs" ]; then
    echo -e "${YELLOW}Found Node.js processes that might be related to the frontend:${NC}"
    echo "$node_procs"
    
    read -p "Do you want to kill these processes? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        ps aux | grep -E "node.*start" | grep -v grep | awk '{print $2}' | xargs -r kill
        echo -e "${GREEN}Node.js processes killed.${NC}"
    fi
fi

# Check for Go processes
go_procs=$(ps aux | grep -E "go run.*main.go" | grep -v grep || echo "")
if [ -n "$go_procs" ]; then
    echo -e "${YELLOW}Found Go processes that might be related to the backend:${NC}"
    echo "$go_procs"
    
    read -p "Do you want to kill these processes? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        ps aux | grep -E "go run.*main.go" | grep -v grep | awk '{print $2}' | xargs -r kill
        echo -e "${GREEN}Go processes killed.${NC}"
    fi
fi

print_header "All services have been stopped! 🛑"
echo -e "${BLUE}To restart services, run:${NC} ${YELLOW}./run.sh${NC}" 