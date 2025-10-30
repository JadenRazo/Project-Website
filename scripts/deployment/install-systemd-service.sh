#!/bin/bash

# Script to install the Jaden Razo Website Backend Service as a systemd service
# This script must be run as root

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Error: This script must be run as root${NC}"
  echo "Please run with sudo or as the root user"
  exit 1
fi

# Get the current directory (should be project root)
CURRENT_DIR=$(pwd)

# Print header
echo -e "${GREEN}=== Installing Jaden Razo Website Backend Service ===${NC}\n"

# Copy the systemd service file
echo -e "${YELLOW}Copying systemd service file...${NC}"
cp "${CURRENT_DIR}/deploy/systemd/jadenrazo-backend.service" /etc/systemd/system/

# Make sure start-services.sh is executable
echo -e "${YELLOW}Making start-services.sh executable...${NC}"
chmod +x "${CURRENT_DIR}/deploy/start-services.sh"

# Create logs directory if it doesn't exist
echo -e "${YELLOW}Creating logs directory...${NC}"
mkdir -p "${CURRENT_DIR}/logs"

# Reload systemd daemon
echo -e "${YELLOW}Reloading systemd daemon...${NC}"
systemctl daemon-reload

# Enable the service to start on boot
echo -e "${YELLOW}Enabling service to start on boot...${NC}"
systemctl enable jadenrazo-backend

# Start the service
echo -e "${YELLOW}Starting service...${NC}"
systemctl start jadenrazo-backend

# Check status
echo -e "${YELLOW}Checking service status...${NC}"
systemctl status jadenrazo-backend

echo -e "\n${GREEN}=== Installation Complete ===${NC}"
echo -e "The backend service is now installed and running."
echo -e "You can manage it with the following commands:"
echo -e "  ${YELLOW}systemctl status jadenrazo-backend${NC} - Check service status"
echo -e "  ${YELLOW}systemctl start jadenrazo-backend${NC} - Start the service"
echo -e "  ${YELLOW}systemctl stop jadenrazo-backend${NC} - Stop the service"
echo -e "  ${YELLOW}systemctl restart jadenrazo-backend${NC} - Restart the service"
echo -e "  ${YELLOW}journalctl -u jadenrazo-backend${NC} - View service logs" 