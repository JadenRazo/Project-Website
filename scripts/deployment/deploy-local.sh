#!/bin/bash

# Local deployment script for testing
echo "=== Portfolio Website Local Deployment ==="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}1. Building backend services...${NC}"
cd /main/Project-Website/backend
go build -o bin/api cmd/api/main.go
go build -o bin/devpanel cmd/devpanel/main.go

echo -e "${GREEN}2. Frontend already built${NC}"

echo -e "${GREEN}3. Starting services...${NC}"
# Kill any existing services
pkill -f "bin/api" 2>/dev/null || true
pkill -f "bin/devpanel" 2>/dev/null || true
sleep 2

# Start services
./bin/api > /tmp/api.log 2>&1 &
./bin/devpanel > /tmp/devpanel.log 2>&1 &

sleep 3

# Check if services are running
if pgrep -f "bin/api" > /dev/null; then
    echo -e "${GREEN}✓ API service is running${NC}"
else
    echo -e "${RED}✗ API service failed to start${NC}"
    tail -n 20 /tmp/api.log
fi

if pgrep -f "bin/devpanel" > /dev/null; then
    echo -e "${GREEN}✓ DevPanel service is running${NC}"
else
    echo -e "${RED}✗ DevPanel service failed to start${NC}"
    tail -n 20 /tmp/devpanel.log
fi

echo ""
echo -e "${GREEN}Local deployment complete!${NC}"
echo ""
echo "Services are running on:"
echo "- API: http://localhost:8080"
echo "- DevPanel API: http://localhost:8081"
echo ""
echo "Admin credentials:"
echo "Email: admin@jadenrazo.dev"
echo "Password: JadenRazoAdmin5795%"
echo ""
echo "Note: For production deployment, copy the built files to your server and run deploy.sh"