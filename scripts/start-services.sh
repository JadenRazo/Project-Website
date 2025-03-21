#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting all services...${NC}"

# Function to check if a port is in use
check_port() {
    lsof -i :$1 >/dev/null 2>&1
    return $?
}

# Function to start a service
start_service() {
    local name=$1
    local port=$2
    local command=$3
    
    if check_port $port; then
        echo -e "${GREEN}$name is already running on port $port${NC}"
    else
        echo -e "${BLUE}Starting $name...${NC}"
        $command &
        sleep 2
        if check_port $port; then
            echo -e "${GREEN}$name started successfully on port $port${NC}"
        else
            echo -e "${RED}Failed to start $name${NC}"
            exit 1
        fi
    fi
}

# Start API service
start_service "API" 8080 "go run cmd/api/main.go"

# Start DevPanel service
start_service "DevPanel" 8081 "go run cmd/devpanel/main.go"

# Start Messaging service
start_service "Messaging" 8082 "go run cmd/messaging/main.go"

# Start URL Shortener service
start_service "URL Shortener" 8083 "go run cmd/urlshortener/main.go"

# Start frontend development server
echo -e "${BLUE}Starting frontend development server...${NC}"
cd frontend && npm run dev &

# Wait for all background processes
wait 