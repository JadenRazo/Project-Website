#!/bin/bash
set -e

# Print colorful messages
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Project-Website Setup Script${NC}"
echo -e "${BLUE}========================================${NC}"

# Define project directories
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BACKEND_DIR="${ROOT_DIR}/backend"
FRONTEND_DIR="${ROOT_DIR}/frontend"

# Navigate to root directory
cd "${ROOT_DIR}"
echo -e "${GREEN}Project root directory:${NC} ${YELLOW}${ROOT_DIR}${NC}"

# Check if directories exist
if [ ! -d "${BACKEND_DIR}" ]; then
    echo -e "${RED}Error: Backend directory not found: ${BACKEND_DIR}${NC}"
    exit 1
fi

if [ ! -d "${FRONTEND_DIR}" ]; then
    echo -e "${RED}Error: Frontend directory not found: ${FRONTEND_DIR}${NC}"
    exit 1
fi

echo -e "\n${GREEN}This script will set up your development environment.${NC}"
echo -e "${YELLOW}For running the application, please use run.sh after setup is complete.${NC}"

# Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go first.${NC}"
    echo -e "  Visit ${YELLOW}https://go.dev/doc/install${NC} for installation instructions."
    exit 1
fi

# Check for Node.js and npm
if ! command -v node &> /dev/null || ! command -v npm &> /dev/null; then
    echo -e "${RED}Error: Node.js or npm is not installed. Please install them first.${NC}"
    echo -e "  Visit ${YELLOW}https://nodejs.org/${NC} for installation instructions."
    exit 1
fi

# Check for tmux
if ! command -v tmux &> /dev/null; then
    echo -e "${YELLOW}tmux is recommended for running the application.${NC}"
    echo -e "  Ubuntu/Debian: ${YELLOW}sudo apt-get install tmux${NC}"
    echo -e "  CentOS/RHEL: ${YELLOW}sudo yum install tmux${NC}"
    echo -e "  MacOS: ${YELLOW}brew install tmux${NC}"
    echo -e "${YELLOW}The run.sh script will attempt to install it for you if missing.${NC}"
fi

echo -e "\n${GREEN}Step 1:${NC} Setting up backend dependencies..."
cd "${BACKEND_DIR}"
go mod tidy
go mod download

echo -e "\n${GREEN}Step 2:${NC} Setting up frontend dependencies..."
cd "${FRONTEND_DIR}"
npm install

echo -e "\n${GREEN}Step 3:${NC} Creating configuration files..."
cd "${BACKEND_DIR}"
if [ ! -f "config/development.yaml" ]; then
    echo -e "${YELLOW}Creating default development.yaml configuration...${NC}"
    cp config/development.yaml.example config/development.yaml 2>/dev/null || echo "development.yaml.example not found, skipping..."
fi

echo -e "\n${GREEN}Step 4:${NC} Setting up the database..."
cd "${BACKEND_DIR}"
echo -e "${YELLOW}Running database migrations...${NC}"
go run cmd/migration/main.go

echo -e "\n${GREEN}Setup complete!${NC}"
echo -e "\n${YELLOW}To run the application, use:${NC}"
echo -e "  ${BLUE}${BACKEND_DIR}/scripts/run.sh${NC}"
echo -e "\n${YELLOW}Available options for run.sh:${NC}"
echo -e "  ${BLUE}-e, --env ENV${NC}        Run in specific environment (development, staging, production)"
echo -e "  ${BLUE}-w, --watch${NC}          Run backend with file watching (requires air)"
echo -e "  ${BLUE}-d, --debug${NC}          Run in debug mode"
echo -e "  ${BLUE}-i, --install${NC}        Install/update dependencies"
echo -e "  ${BLUE}-s, --setup${NC}          Setup configuration files"
echo -e "  ${BLUE}-h, --help${NC}           Print help message"
