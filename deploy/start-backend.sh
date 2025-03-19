#!/bin/bash

# Start backend services script for Jaden Razo Website
# This script starts all three backend components in parallel

LOG_DIR="/root/Project-Website/logs"
mkdir -p $LOG_DIR

# Function to start a service
start_service() {
    local service_name=$1
    local service_bin=$2
    local service_config=$3
    
    echo "Starting $service_name..."
    
    # Run the service and redirect output to log files
    nohup $service_bin -config $service_config > "$LOG_DIR/$service_name.log" 2>&1 &
    local pid=$!
    
    echo "$service_name started with PID $pid"
    echo $pid > "$LOG_DIR/$service_name.pid"
}

# Kill any existing processes
echo "Stopping any existing processes..."
if [ -f "$LOG_DIR/devpanel.pid" ]; then
    kill -15 $(cat "$LOG_DIR/devpanel.pid") 2>/dev/null || true
    rm "$LOG_DIR/devpanel.pid"
fi

if [ -f "$LOG_DIR/urlshortener.pid" ]; then
    kill -15 $(cat "$LOG_DIR/urlshortener.pid") 2>/dev/null || true
    rm "$LOG_DIR/urlshortener.pid"
fi

if [ -f "$LOG_DIR/messaging.pid" ]; then
    kill -15 $(cat "$LOG_DIR/messaging.pid") 2>/dev/null || true
    rm "$LOG_DIR/messaging.pid"
fi

# Wait a moment for processes to clean up
sleep 2

# Start developer panel
start_service "devpanel" "/root/Project-Website/bin/devpanel" "/root/Project-Website/config/devpanel.yaml"

# Start URL shortener
start_service "urlshortener" "/root/Project-Website/bin/urlshortener" "/root/Project-Website/config/urlshortener.yaml"

# Start messaging platform
start_service "messaging" "/root/Project-Website/bin/messaging" "/root/Project-Website/config/messaging.yaml"

# Keep the script running to maintain the systemd service active
# This also allows us to monitor the logs
echo "All services started. Monitoring logs..."

tail -f "$LOG_DIR"/*.log 