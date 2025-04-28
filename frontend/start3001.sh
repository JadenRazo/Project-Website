#!/bin/bash
# Simple script to start React on port 3001 and bind to all interfaces
echo "Starting React frontend on port 3001 (accessible externally)..."

# More aggressively kill any process on port 3001
echo "Checking for processes on port 3001..."
PORT_PID=$(lsof -ti:3001)
if [ ! -z "$PORT_PID" ]; then
  echo "Port 3001 is already in use by PID $PORT_PID. Killing process..."
  kill -9 $PORT_PID
  sleep 2
  echo "Process killed, starting new instance."
fi

# Kill any React scripts that might be running
echo "Killing any existing React processes..."
pkill -f "react-scripts" || true
sleep 1

# Create a .env file to force the port and avoid the prompt
echo "Setting up React environment..."
cat > .env << EOF
PORT=3001
HOST=0.0.0.0
BROWSER=none
FAST_REFRESH=false
FORCE_COLOR=true
EOF

# Start React directly, bypassing the prompt mechanism
echo "Starting React app on port 3001..."
npm run start 