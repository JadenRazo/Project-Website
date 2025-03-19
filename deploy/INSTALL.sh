#!/bin/bash

# Installation script for Jaden Razo's Website Backend
# This script automates the deployment process for Ubuntu servers

set -e  # Exit on error

# Print colored output
print_green() {
    echo -e "\e[32m$1\e[0m"
}

print_yellow() {
    echo -e "\e[33m$1\e[0m"
}

print_red() {
    echo -e "\e[31m$1\e[0m"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    print_red "Please run as root"
    exit 1
fi

# Confirm installation
print_yellow "This script will install the Jaden Razo Website Backend on this server."
print_yellow "It will install: Nginx, PostgreSQL, Go, Node.js, and all required dependencies."
read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_red "Installation aborted."
    exit 1
fi

# Create a backup directory
print_green "Creating backup directory..."
mkdir -p /root/backups

# Step 1: Update system and install dependencies
print_green "Step 1: Updating system and installing dependencies..."
apt update && apt upgrade -y
apt install -y nginx postgresql postgresql-contrib certbot python3-certbot-nginx golang-go git build-essential curl

# Install Node.js
print_green "Installing Node.js..."
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt install -y nodejs

# Step 2: Set up project directory
print_green "Step 2: Setting up project directory..."
mkdir -p /root/Project-Website
cd /root/Project-Website

# Check if the directory is empty, if not, offer to backup
if [ "$(ls -A /root/Project-Website)" ]; then
    print_yellow "Project directory is not empty."
    read -p "Backup existing files before proceeding? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        backup_dir="/root/backups/project_$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$backup_dir"
        cp -r /root/Project-Website/* "$backup_dir"
        print_green "Backup created at $backup_dir"
    fi
    
    read -p "Clear existing files? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -rf /root/Project-Website/*
        print_green "Existing files cleared."
    else
        print_red "Cannot continue without clearing the directory."
        exit 1
    fi
fi

# Ask for the repository URL
print_yellow "Please enter the git repository URL:"
read -p "URL [https://github.com/jadenrazo/website.git]: " repo_url
repo_url=${repo_url:-https://github.com/jadenrazo/website.git}

# Clone the repository
print_green "Cloning repository..."
git clone "$repo_url" .

# Step 3: Set up PostgreSQL
print_green "Step 3: Setting up PostgreSQL..."
systemctl enable postgresql
systemctl start postgresql

# Create database and user
print_yellow "Setting up PostgreSQL database..."
sudo -u postgres psql -c "CREATE DATABASE urlshortener;" || true
sudo -u postgres psql -c "CREATE USER postgres WITH ENCRYPTED PASSWORD 'postgres';" || true
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE urlshortener TO postgres;" || true

# Step 4: Build backend services
print_green "Step 4: Building backend services..."
go mod tidy
mkdir -p bin
go build -o bin/devpanel cmd/devpanel/main.go
go build -o bin/urlshortener cmd/urlshortener/main.go
go build -o bin/messaging cmd/messaging/main.go

# Make startup script executable
chmod +x deploy/start-backend.sh

# Step 5: Configure backend services
print_green "Step 5: Configuring backend services..."
mkdir -p config
if [ -d "deploy/config-examples" ]; then
    cp deploy/config-examples/*.yaml config/ 2>/dev/null || true
else
    print_yellow "Config examples not found. You'll need to create config files manually."
fi

# Step 6: Build frontend
print_green "Step 6: Building frontend..."
if [ -d "frontend" ]; then
    cd frontend
    npm install
    npm run build
    cd ..
    print_green "Frontend built successfully."
else
    print_yellow "Frontend directory not found. Skipping frontend build."
fi

# Step 7: Configure Nginx
print_green "Step 7: Configuring Nginx..."
if [ -f "deploy/nginx/jadenrazo.conf" ]; then
    cp deploy/nginx/jadenrazo.conf /etc/nginx/sites-available/
    ln -sf /etc/nginx/sites-available/jadenrazo.conf /etc/nginx/sites-enabled/
    
    # Ask for domain name
    print_yellow "Please enter your domain name:"
    read -p "Domain [jadenrazo.dev]: " domain_name
    domain_name=${domain_name:-jadenrazo.dev}
    
    # Update domain in nginx config
    sed -i "s/jadenrazo.dev/$domain_name/g" /etc/nginx/sites-available/jadenrazo.conf
    
    # Disable default site if enabled
    if [ -f "/etc/nginx/sites-enabled/default" ]; then
        rm -f /etc/nginx/sites-enabled/default
    fi
    
    # Test nginx config
    nginx -t
    systemctl reload nginx
    print_green "Nginx configured successfully."
else
    print_yellow "Nginx config not found. You'll need to configure Nginx manually."
fi

# Step 8: Set up systemd service
print_green "Step 8: Setting up systemd service..."
if [ -f "deploy/systemd/backend.service" ]; then
    cp deploy/systemd/backend.service /etc/systemd/system/jadenrazo-backend.service
    systemctl daemon-reload
    systemctl enable jadenrazo-backend
    systemctl start jadenrazo-backend
    print_green "Systemd service configured and started."
else
    print_yellow "Systemd service file not found. You'll need to configure the service manually."
fi

# Step 9: Set up SSL with Let's Encrypt
print_green "Step 9: Setting up SSL with Let's Encrypt..."
print_yellow "Do you want to set up SSL with Let's Encrypt now?"
read -p "This requires your domain to be pointed to this server already (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [ -z "$domain_name" ]; then
        read -p "Domain name: " domain_name
    fi
    
    certbot --nginx -d $domain_name -d www.$domain_name
    systemctl enable certbot.timer
    systemctl start certbot.timer
    print_green "SSL certificates installed."
else
    print_yellow "SSL setup skipped. You can run 'certbot --nginx' later to set up SSL."
fi

# Step 10: Create logs directory
print_green "Step 10: Creating logs directory..."
mkdir -p /root/Project-Website/logs

# Step 11: Summary
print_green "============================================="
print_green "Installation completed!"
print_green "============================================="
print_green "Backend services: http://localhost:8080 (Developer Panel)"
print_green "               http://localhost:8081 (URL Shortener)"
print_green "               http://localhost:8082 (Messaging Platform)"
print_green ""
print_green "Frontend: http://$domain_name (if configured)"
print_green ""
print_green "Check service status: systemctl status jadenrazo-backend"
print_green "View logs: tail -f /root/Project-Website/logs/*.log"
print_green "=============================================" 