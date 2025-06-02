#!/bin/bash

# Quick DevPanel Demo Startup Script
echo "🚀 Starting DevPanel Demo..."

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Kill any existing processes on our ports
echo -e "${YELLOW}Cleaning up existing processes...${NC}"
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
lsof -ti:3000 | xargs kill -9 2>/dev/null || true

# Start the simple API
echo -e "${BLUE}Starting Simple API Server...${NC}"
cd backend
go run cmd/simple-api/main.go > /tmp/devpanel-api.log 2>&1 &
API_PID=$!
cd ..

# Wait for API to start
sleep 3

# Test API
if curl -s http://localhost:8080/health > /dev/null; then
    echo -e "${GREEN}✅ API Server started successfully on port 8080${NC}"
else
    echo -e "${YELLOW}⚠️  API might still be starting...${NC}"
fi

# Start the frontend
echo -e "${BLUE}Starting Frontend Development Server...${NC}"
cd frontend
npm start > /tmp/devpanel-frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..

# Wait for frontend to start
sleep 5

echo -e "${GREEN}🎉 DevPanel Demo Environment Started!${NC}"
echo ""
echo -e "${BLUE}📋 Service URLs:${NC}"
echo -e "  Frontend:  ${GREEN}http://localhost:3000${NC}"
echo -e "  API:       ${GREEN}http://localhost:8080${NC}"
echo -e "  DevPanel:  ${GREEN}http://localhost:3000/devpanel${NC}"
echo ""
echo -e "${BLUE}📊 API Endpoints:${NC}"
echo -e "  Health:    ${GREEN}http://localhost:8080/health${NC}"
echo -e "  Projects:  ${GREEN}http://localhost:8080/api/v1/projects${NC}"
echo -e "  Featured:  ${GREEN}http://localhost:8080/api/v1/projects/featured${NC}"
echo ""
echo -e "${YELLOW}🛑 To stop all services:${NC}"
echo -e "  kill $API_PID $FRONTEND_PID"
echo ""
echo -e "${YELLOW}💡 Log files:${NC}"
echo -e "  API:       tail -f /tmp/devpanel-api.log"
echo -e "  Frontend:  tail -f /tmp/devpanel-frontend.log"
echo ""

# Keep script running
wait