#!/bin/bash

# Script to check system dependencies and database setup
# This helps ensure all required components are installed and configured

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Print header
echo -e "${GREEN}=== Jaden Razo Website Backend System Check ===${NC}\n"

# Check Go installation
echo -e "${YELLOW}Checking Go installation...${NC}"
if command -v go >/dev/null 2>&1; then
    go_version=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓ Go is installed ($go_version)${NC}"
else
    echo -e "${RED}✗ Go is not installed. Please install Go 1.21 or later.${NC}"
    echo -e "   Visit https://golang.org/doc/install for installation instructions."
fi

# Check Node.js installation (for frontend)
echo -e "\n${YELLOW}Checking Node.js installation...${NC}"
if command -v node >/dev/null 2>&1; then
    node_version=$(node --version)
    echo -e "${GREEN}✓ Node.js is installed ($node_version)${NC}"
else
    echo -e "${RED}✗ Node.js is not installed. It's required for the frontend.${NC}"
    echo -e "   Run: curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -"
    echo -e "   Then: sudo apt-get install -y nodejs"
fi

# Check PostgreSQL installation
echo -e "\n${YELLOW}Checking PostgreSQL installation...${NC}"
if command -v psql >/dev/null 2>&1; then
    pg_version=$(psql --version | awk '{print $3}')
    echo -e "${GREEN}✓ PostgreSQL is installed ($pg_version)${NC}"
    
    # Check if PostgreSQL service is running
    if systemctl is-active --quiet postgresql; then
        echo -e "${GREEN}✓ PostgreSQL service is running${NC}"
    else
        echo -e "${RED}✗ PostgreSQL service is not running${NC}"
        echo -e "   Run: sudo systemctl start postgresql"
    fi
    
    # Try to connect to PostgreSQL
    echo -e "\n${YELLOW}Testing PostgreSQL connection...${NC}"
    if sudo -u postgres psql -c '\l' >/dev/null 2>&1; then
        echo -e "${GREEN}✓ PostgreSQL connection successful${NC}"
        
        # Check if urlshortener database exists
        if sudo -u postgres psql -lqt | cut -d \| -f 1 | grep -qw urlshortener; then
            echo -e "${GREEN}✓ 'urlshortener' database exists${NC}"
        else
            echo -e "${RED}✗ 'urlshortener' database does not exist${NC}"
            echo -e "   Run: sudo -u postgres psql -c 'CREATE DATABASE urlshortener;'"
        fi
    else
        echo -e "${RED}✗ Failed to connect to PostgreSQL${NC}"
    fi
else
    echo -e "${RED}✗ PostgreSQL is not installed${NC}"
    echo -e "   Run: sudo apt-get install postgresql postgresql-contrib"
fi

# Check Nginx installation
echo -e "\n${YELLOW}Checking Nginx installation...${NC}"
if command -v nginx >/dev/null 2>&1; then
    nginx_version=$(nginx -v 2>&1 | awk -F/ '{print $2}')
    echo -e "${GREEN}✓ Nginx is installed ($nginx_version)${NC}"
    
    # Check if Nginx service is running
    if systemctl is-active --quiet nginx; then
        echo -e "${GREEN}✓ Nginx service is running${NC}"
    else
        echo -e "${RED}✗ Nginx service is not running${NC}"
        echo -e "   Run: sudo systemctl start nginx"
    fi
    
    # Check if our configuration exists
    if [ -f "/etc/nginx/sites-available/jadenrazo.conf" ]; then
        echo -e "${GREEN}✓ Nginx site configuration exists${NC}"
        
        if [ -f "/etc/nginx/sites-enabled/jadenrazo.conf" ]; then
            echo -e "${GREEN}✓ Nginx site is enabled${NC}"
        else
            echo -e "${RED}✗ Nginx site is not enabled${NC}"
            echo -e "   Run: sudo ln -s /etc/nginx/sites-available/jadenrazo.conf /etc/nginx/sites-enabled/"
        fi
    else
        echo -e "${RED}✗ Nginx site configuration does not exist${NC}"
        echo -e "   Run: sudo cp deploy/nginx/jadenrazo.conf /etc/nginx/sites-available/"
    fi
else
    echo -e "${RED}✗ Nginx is not installed${NC}"
    echo -e "   Run: sudo apt-get install nginx"
fi

# Check required directories
echo -e "\n${YELLOW}Checking required directories...${NC}"
directories=("bin" "config" "logs")
for dir in "${directories[@]}"; do
    if [ -d "$dir" ]; then
        echo -e "${GREEN}✓ Directory '$dir' exists${NC}"
    else
        echo -e "${RED}✗ Directory '$dir' does not exist${NC}"
        echo -e "   Run: mkdir -p $dir"
    fi
done

# Check config files
echo -e "\n${YELLOW}Checking configuration files...${NC}"
config_files=("config/devpanel.yaml" "config/urlshortener.yaml" "config/messaging.yaml")
for config in "${config_files[@]}"; do
    if [ -f "$config" ]; then
        echo -e "${GREEN}✓ Configuration file '$config' exists${NC}"
    else
        echo -e "${RED}✗ Configuration file '$config' does not exist${NC}"
        echo -e "   You need to create this configuration file."
    fi
done

# Check if binaries are built
echo -e "\n${YELLOW}Checking built binaries...${NC}"
binaries=("bin/devpanel" "bin/urlshortener" "bin/messaging")
for binary in "${binaries[@]}"; do
    if [ -f "$binary" ]; then
        echo -e "${GREEN}✓ Binary '$binary' exists${NC}"
    else
        base_name=$(basename "$binary")
        echo -e "${RED}✗ Binary '$binary' does not exist${NC}"
        
        # Special handling for messaging binary
        if [ "$base_name" = "messaging" ] && [ -f "cmd/messaging/simple_main.go" ]; then
            echo -e "   Run: go build -o $binary cmd/$base_name/simple_main.go"
        else
            echo -e "   Run: go build -o $binary cmd/$base_name/main.go"
        fi
    fi
done

# Check messaging implementation
echo -e "\n${YELLOW}Checking messaging implementation...${NC}"
if [ -f "cmd/messaging/simple_main.go" ]; then
    echo -e "${GREEN}✓ Simplified messaging implementation is available${NC}"
    
    if [ -f "cmd/messaging/main.go" ]; then
        echo -e "${YELLOW}⚠ Both simplified and full messaging implementations exist${NC}"
        echo -e "   Note: The start-services.sh script will prioritize the simplified version"
    fi
else
    if [ -f "cmd/messaging/main.go" ]; then
        echo -e "${YELLOW}⚠ Only full messaging implementation exists${NC}"
        echo -e "   Note: You may encounter build errors due to duplicate declarations"
    else
        echo -e "${RED}✗ No messaging implementation found${NC}"
        echo -e "   Create either cmd/messaging/simple_main.go or cmd/messaging/main.go"
    fi
fi

# Print summary
echo -e "\n${GREEN}=== System Check Summary ===${NC}"
echo -e "${YELLOW}To start all services in development mode:${NC}"
echo -e "   ./deploy/start-services.sh"
echo -e "\n${YELLOW}To deploy as a service:${NC}"
echo -e "   Follow the instructions in deploy/QUICKSTART.md"
echo -e "\n${YELLOW}For more detailed deployment information:${NC}"
echo -e "   See deploy/DEPLOY.md" 