#!/bin/bash

# Script to start all backend services for the Jaden Razo website
# This script manages the developer panel, URL shortener, and messaging platform

# Set default configuration paths
DEVPANEL_CONFIG="config/devpanel.yaml"
URLSHORTENER_CONFIG="config/urlshortener.yaml"
MESSAGING_CONFIG="config/messaging.yaml"

# Set default ports
DEVPANEL_PORT=8080
URLSHORTENER_PORT=8081
MESSAGING_PORT=8082

# Create logs directory if it doesn't exist
LOGS_DIR="logs"
mkdir -p $LOGS_DIR

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Print header
echo -e "${GREEN}=== Jaden Razo Website Backend Services ===${NC}"
echo -e "${GREEN}Starting all backend services...${NC}"

# Flag to track if we're exiting from the trap handler
EXITING_FROM_TRAP=0

# Function to check if a port is already in use
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null ; then
        return 0 # Port is in use
    else
        return 1 # Port is free
    fi
}

# Function to start a service
start_service() {
    local name=$1
    local binary=$2
    local config=$3
    local port=$4
    
    echo -e "${YELLOW}Starting $name on port $port...${NC}"
    
    # Check if the binary exists
    if [ ! -f "bin/$binary" ]; then
        echo -e "${RED}Error: Binary 'bin/$binary' not found.${NC}"
        echo -e "${YELLOW}Attempting to build it...${NC}"
        
        # Determine the main.go file path based on service name
        local main_path="cmd/$binary/main.go"
        
        # For messaging service, check if simple_main.go exists and use it
        if [ "$binary" = "messaging" ] && [ -f "cmd/$binary/simple_main.go" ]; then
            echo -e "${YELLOW}Using simplified messaging implementation...${NC}"
            go build -o "bin/$binary" "cmd/$binary/simple_main.go"
        elif [ -f "$main_path" ]; then
            go build -o "bin/$binary" "$main_path"
        else
            echo -e "${RED}Error: Main file at $main_path not found.${NC}"
            return 1
        fi
        
        if [ $? -ne 0 ]; then
            echo -e "${RED}Failed to build $binary. Check the Go code for errors.${NC}"
            return 1
        fi
        echo -e "${GREEN}Successfully built $binary${NC}"
    fi
    
    # Check if the config file exists
    if [ ! -f "$config" ]; then
        echo -e "${RED}Warning: Config file '$config' not found. Using default configuration.${NC}"
    fi
    
    # Check if the port is already in use
    if check_port $port; then
        echo -e "${RED}Error: Port $port is already in use. Cannot start $name.${NC}"
        return 1
    fi
    
    # Start the service and redirect output to log file
    nohup bin/$binary -config $config > "$LOGS_DIR/$binary.log" 2>&1 &
    local pid=$!
    
    # Check if service started successfully by waiting for a moment and checking if process exists
    sleep 2
    if kill -0 $pid 2>/dev/null; then
        echo -e "${GREEN}$name started successfully with PID $pid${NC}"
        echo $pid > "$LOGS_DIR/$binary.pid"
        return 0
    else
        echo -e "${RED}Failed to start $name. Check the logs at $LOGS_DIR/$binary.log${NC}"
        return 1
    fi
}

# Function to stop services
stop_services() {
    echo -e "${YELLOW}Stopping all services...${NC}"
    
    # Check for PID files and kill processes
    for service in devpanel urlshortener messaging; do
        if [ -f "$LOGS_DIR/$service.pid" ]; then
            pid=$(cat "$LOGS_DIR/$service.pid")
            echo -e "${YELLOW}Stopping $service (PID: $pid)...${NC}"
            kill -15 $pid 2>/dev/null
            rm "$LOGS_DIR/$service.pid"
            echo -e "${GREEN}$service stopped${NC}"
        fi
    done
    
    echo -e "${GREEN}All services stopped${NC}"
    
    # Only exit if called from the trap handler
    if [ $EXITING_FROM_TRAP -eq 1 ]; then
        exit 0
    fi
}

# Trap handler for Ctrl+C
trap_handler() {
    EXITING_FROM_TRAP=1
    stop_services
}

# Set up trap for graceful shutdown on CTRL+C
trap trap_handler INT TERM

# Stop existing services first
stop_services

# Start Developer Panel
start_service "Developer Panel" "devpanel" "$DEVPANEL_CONFIG" $DEVPANEL_PORT

# Start URL Shortener
start_service "URL Shortener" "urlshortener" "$URLSHORTENER_CONFIG" $URLSHORTENER_PORT

# Start Messaging Platform (if it exists)
if [ -f "cmd/messaging/simple_main.go" ] || [ -f "cmd/messaging/main.go" ]; then
    start_service "Messaging Platform" "messaging" "$MESSAGING_CONFIG" $MESSAGING_PORT
else
    echo -e "${YELLOW}Messaging Platform main.go not found. Skipping.${NC}"
fi

# Display service status
echo -e "\n${GREEN}=== Service Status ===${NC}"
echo -e "${GREEN}Developer Panel:     http://localhost:$DEVPANEL_PORT${NC}"
echo -e "${GREEN}URL Shortener:       http://localhost:$URLSHORTENER_PORT${NC}"
if [ -f "$LOGS_DIR/messaging.pid" ]; then
    echo -e "${GREEN}Messaging Platform:  http://localhost:$MESSAGING_PORT${NC}"
fi

echo -e "\n${YELLOW}Logs are available in the $LOGS_DIR directory${NC}"
echo -e "${YELLOW}To stop all services, press Ctrl+C${NC}"

# Keep the script running to allow for Ctrl+C to work
echo -e "\n${GREEN}Monitoring services. Press Ctrl+C to stop all services.${NC}"
while true; do
    sleep 1
done 