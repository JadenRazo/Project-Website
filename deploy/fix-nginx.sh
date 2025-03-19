#!/bin/bash

# Script to fix Nginx configuration for Jadenrazo.dev
# This script will:
# 1. Create required directories
# 2. Install the new Nginx configuration
# 3. Test and reload Nginx
# 4. Check backend services

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

echo -e "${GREEN}=== Fixing Nginx Configuration for Jadenrazo.dev ===${NC}\n"

# Create Let's Encrypt verification directory
echo -e "${YELLOW}Creating Let's Encrypt verification directory...${NC}"
mkdir -p /var/www/letsencrypt/.well-known/acme-challenge
chmod -R 755 /var/www/letsencrypt

# Check for existing SSL certificates
echo -e "${YELLOW}Checking for SSL certificates...${NC}"
SSL_DIR="/etc/letsencrypt/live/jadenrazo.dev"
if [ ! -d "$SSL_DIR" ]; then
  echo -e "${YELLOW}SSL certificates not found at $SSL_DIR${NC}"
  echo -e "${YELLOW}You'll need to either:${NC}"
  echo "1. Generate Let's Encrypt certificates using certbot"
  echo "2. Create a self-signed certificate for testing"
  echo "3. Modify the Nginx config to use your existing certificates"
  
  # Ask if we should generate a self-signed certificate
  read -p "Generate a self-signed certificate for testing? (y/n): " -n 1 -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Generating self-signed certificate...${NC}"
    mkdir -p /etc/nginx/ssl
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
      -keyout /etc/nginx/ssl/jadenrazo.key \
      -out /etc/nginx/ssl/jadenrazo.crt \
      -subj "/CN=jadenrazo.dev/O=Jaden Razo/C=US"
    
    # Update the Nginx configuration to use the self-signed certificate
    sed -i 's|ssl_certificate .*|ssl_certificate /etc/nginx/ssl/jadenrazo.crt;|' deploy/nginx/jadenrazo.conf
    sed -i 's|ssl_certificate_key .*|ssl_certificate_key /etc/nginx/ssl/jadenrazo.key;|' deploy/nginx/jadenrazo.conf
    
    echo -e "${GREEN}Self-signed certificate generated and Nginx config updated${NC}"
  fi
fi

# Install the Nginx configuration
echo -e "${YELLOW}Installing Nginx configuration...${NC}"
cp deploy/nginx/jadenrazo.conf /etc/nginx/sites-available/

# Create symlink if it doesn't exist
if [ ! -f /etc/nginx/sites-enabled/jadenrazo.conf ]; then
  echo -e "${YELLOW}Creating symlink to enable the site...${NC}"
  ln -s /etc/nginx/sites-available/jadenrazo.conf /etc/nginx/sites-enabled/
fi

# Remove default site if it exists and user confirms
if [ -f /etc/nginx/sites-enabled/default ]; then
  read -p "Remove default Nginx site? (y/n): " -n 1 -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Removing default Nginx site...${NC}"
    rm -f /etc/nginx/sites-enabled/default
  fi
fi

# Test Nginx configuration
echo -e "${YELLOW}Testing Nginx configuration...${NC}"
nginx -t
if [ $? -ne 0 ]; then
  echo -e "${RED}Nginx configuration test failed. Please check the errors above.${NC}"
  exit 1
fi

# Reload Nginx
echo -e "${YELLOW}Reloading Nginx...${NC}"
systemctl reload nginx
if [ $? -ne 0 ]; then
  echo -e "${RED}Failed to reload Nginx. Trying to restart...${NC}"
  systemctl restart nginx
  if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to restart Nginx. Please check the Nginx error logs.${NC}"
    echo -e "${YELLOW}Nginx error logs: /var/log/nginx/error.log${NC}"
    exit 1
  fi
fi

# Check backend services
echo -e "${YELLOW}Checking backend services...${NC}"
ps aux | grep bin/ | grep -v grep
if [ $? -ne 0 ]; then
  echo -e "${RED}No backend services found running.${NC}"
  echo -e "${YELLOW}Starting backend services...${NC}"
  
  if [ -f "/root/Project-Website/deploy/start-services.sh" ]; then
    # Try to run the start services script
    cd /root/Project-Website
    ./deploy/start-services.sh &
    sleep 5
    echo -e "${GREEN}Backend services started in the background${NC}"
  else
    echo -e "${RED}Could not find start-services.sh script.${NC}"
    echo -e "${YELLOW}Please start your backend services manually.${NC}"
  fi
fi

# Test backend connections
echo -e "${YELLOW}Testing backend connections...${NC}"
if curl -s http://127.0.0.1:8080/devpanel/health > /dev/null; then
  echo -e "${GREEN}Developer Panel service is responding${NC}"
else
  echo -e "${RED}Developer Panel service is not responding${NC}"
fi

if curl -s http://127.0.0.1:8081/health > /dev/null; then
  echo -e "${GREEN}URL Shortener service is responding${NC}"
else
  echo -e "${RED}URL Shortener service is not responding${NC}"
fi

if curl -s http://127.0.0.1:8082/health > /dev/null; then
  echo -e "${GREEN}Messaging service is responding${NC}"
else
  echo -e "${RED}Messaging service is not responding${NC}"
fi

# Print Cloudflare DNS configuration information
echo -e "\n${GREEN}=== Cloudflare DNS Configuration ===${NC}"
echo -e "${YELLOW}To properly set up Cloudflare:${NC}"
echo "1. Make sure your domain has an A record pointing to your server's IP address"
echo "2. Set the 'Proxied' status to ON for the A record"
echo "3. Set SSL/TLS encryption mode to 'Full' in Cloudflare"
echo "4. If using a self-signed certificate, set SSL/TLS encryption mode to 'Flexible'"

# Print next steps
echo -e "\n${GREEN}=== Next Steps ===${NC}"
echo -e "${YELLOW}1. Test your website locally:${NC}"
echo "   curl -H 'Host: jadenrazo.dev' https://localhost -k"
echo -e "${YELLOW}2. Check for any errors in Nginx logs:${NC}"
echo "   tail -f /var/log/nginx/error.log"
echo -e "${YELLOW}3. If you're using a self-signed certificate with Cloudflare:${NC}"
echo "   Set SSL/TLS encryption mode to 'Flexible' in Cloudflare"
echo -e "${YELLOW}4. To test whether Cloudflare can reach your origin:${NC}"
echo "   Turn on Development Mode in Cloudflare temporarily"

echo -e "\n${GREEN}=== Configuration Complete ===${NC}" 