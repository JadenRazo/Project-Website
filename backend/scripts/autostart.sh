#!/bin/bash
set -e

# Define variables for the source script
SKIP_FRONTEND_START=true

# Source the original script for common functions but don't start the frontend
source "$(dirname "$0")/run.sh" --env development

# Define a simpler frontend auto-start function
start_frontend_auto() {
    print_header "Starting frontend on port 3001"
    
    # Kill existing frontend processes
    echo -e "${CYAN}Killing any existing React and frontend processes...${NC}"
    pkill -f "react-scripts" || true
    pkill -f "node.*start.js" || true
    tmux kill-session -t frontend 2>/dev/null || true
    sleep 2
    
    # Kill any processes on port 3001 directly
    echo -e "${CYAN}Killing any processes on port 3001...${NC}"
    fuser -k 3001/tcp || true
    sleep 2
    
    # Set up environment variables file
    echo -e "${CYAN}Creating .env file for React...${NC}"
    cat > "${FRONTEND_DIR}/.env" << EOF
PORT=3001
HOST=0.0.0.0
BROWSER=none
FAST_REFRESH=false
DANGEROUSLY_DISABLE_HOST_CHECK=true
FORCE_COLOR=true
EOF

    # Start React directly
    echo -e "${CYAN}Starting React on port 3001...${NC}"
    cd "${FRONTEND_DIR}"
    tmux new-session -d -s frontend "cd ${FRONTEND_DIR} && PORT=3001 HOST=0.0.0.0 BROWSER=none npm start; echo 'Process exited with code $?'; read"
    
    # Wait and verify
    echo -e "${CYAN}Waiting for React to start on port 3001 (this may take 30 seconds)...${NC}"
    for i in {1..15}; do
        echo -e "${YELLOW}Checking attempt $i/15...${NC}"
        if nc -z localhost 3001; then
            echo -e "${GREEN}✓ Frontend is now running on port 3001!${NC}"
            break
        fi
        sleep 2
    done
    
    # Final check
    if nc -z localhost 3001; then
        echo -e "${GREEN}✓ Frontend successfully started on port 3001${NC}"
    else
        echo -e "${RED}✗ Failed to start frontend service${NC}"
        echo -e "${YELLOW}Debugging information:${NC}"
        echo -e "${YELLOW}1. Check tmux logs: tmux attach-session -t frontend${NC}"
        echo -e "${YELLOW}2. Try starting manually: cd ${FRONTEND_DIR} && PORT=3001 npm start${NC}"
        echo -e "${YELLOW}3. Check if port 3001 is free: lsof -i :3001${NC}"
    fi
}

# Override the frontend start with our auto-confirm version
start_frontend_auto

# Restart nginx to apply configuration changes
if command -v nginx &> /dev/null; then
    echo -e "${CYAN}Restarting Nginx...${NC}"
    sudo systemctl restart nginx
    echo -e "${GREEN}✓ Nginx restarted${NC}"
else
    echo -e "${YELLOW}⚠ Nginx command not found. Please restart Nginx manually.${NC}"
fi

# Final status report
print_header "Service Status Report"

# Check services again
all_services_running=true
echo -e "${CYAN}Backend Services:${NC}"
for service in "${!SERVICES[@]}"; do
    port="${SERVICES[$service]}"
    if nc -z localhost "$port" >/dev/null 2>&1; then
        echo -e "  ${GREEN}✓ ${service} (port ${port})${NC}: Running"
    else
        echo -e "  ${RED}✗ ${service} (port ${port})${NC}: Not running"
        all_services_running=false
    fi
done

echo -e "\n${CYAN}Frontend:${NC}"
if nc -z localhost 3001 >/dev/null 2>&1; then
    echo -e "  ${GREEN}✓ Development server (port 3001)${NC}: Running"
else
    echo -e "  ${RED}✗ Frontend server${NC}: Not running"
    all_services_running=false
fi

echo -e "\n${MAGENTA}Website should now be accessible at https://jadenrazo.dev 🚀${NC}" 