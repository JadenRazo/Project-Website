# Deployment Guide

This guide provides step-by-step instructions for deploying the Jaden Razo website backend services on an Ubuntu server.

## Prerequisites

- Ubuntu 20.04 LTS or newer
- Root access
- Domain pointed to your server's IP address
- Open ports: 80, 443

## Step 1: Install Required Software

```bash
# Update system
apt update && apt upgrade -y

# Install required packages
apt install -y nginx postgresql postgresql-contrib certbot python3-certbot-nginx golang-go git build-essential

# Install Node.js and npm for frontend
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt install -y nodejs
```

## Step 2: Clone the Repository

```bash
# Create project directory
mkdir -p /root/Project-Website
cd /root/Project-Website

# Clone repository
git clone https://github.com/jadenrazo/website.git .
```

## Step 3: Build Backend Services

```bash
# Install Go dependencies
go mod tidy

# Build all services
mkdir -p bin
go build -o bin/devpanel cmd/devpanel/main.go
go build -o bin/urlshortener cmd/urlshortener/main.go
go build -o bin/messaging cmd/messaging/main.go

# Make startup script executable
chmod +x deploy/start-backend.sh
```

## Step 4: Set Up Databases

### PostgreSQL for URL Shortener

```bash
# Create database and user
sudo -u postgres psql -c "CREATE DATABASE urlshortener;"
sudo -u postgres psql -c "CREATE USER postgres WITH ENCRYPTED PASSWORD 'postgres';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE urlshortener TO postgres;"

# Enable PostgreSQL to start on boot
systemctl enable postgresql
systemctl start postgresql
```

## Step 5: Configure Backend Services

Copy and edit configuration files:

```bash
# Create directory
mkdir -p config

# Copy configuration files
cp deploy/config-examples/devpanel.yaml config/
cp deploy/config-examples/urlshortener.yaml config/
cp deploy/config-examples/messaging.yaml config/

# Edit configuration files as needed
# nano config/devpanel.yaml
# nano config/urlshortener.yaml
# nano config/messaging.yaml
```

## Step 6: Build Frontend

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Build for production
npm run build

# Return to project root
cd ..
```

## Step 7: Configure Nginx

```bash
# Copy nginx configuration
cp deploy/nginx/jadenrazo.conf /etc/nginx/sites-available/

# Create symlink
ln -s /etc/nginx/sites-available/jadenrazo.conf /etc/nginx/sites-enabled/

# Remove default configuration if it exists
rm -f /etc/nginx/sites-enabled/default

# Test configuration
nginx -t

# Reload nginx
systemctl reload nginx
```

## Step 8: Set Up SSL with Let's Encrypt

```bash
# Obtain SSL certificate
certbot --nginx -d jadenrazo.dev -d www.jadenrazo.dev

# Ensure certbot renewal is enabled
systemctl enable certbot.timer
systemctl start certbot.timer
```

## Step 9: Set Up Systemd Service

```bash
# Copy systemd service file
cp deploy/systemd/backend.service /etc/systemd/system/jadenrazo-backend.service

# Enable and start the service
systemctl daemon-reload
systemctl enable jadenrazo-backend
systemctl start jadenrazo-backend
```

## Step 10: Verify Deployment

1. Check if the backend services are running:
   ```bash
   systemctl status jadenrazo-backend
   ```

2. Check the logs:
   ```bash
   tail -f /root/Project-Website/logs/*.log
   ```

3. Test the frontend by visiting your domain in a browser:
   ```
   https://jadenrazo.dev
   ```

## Maintenance

### Restarting Services

```bash
# Restart all backend services
systemctl restart jadenrazo-backend

# Restart just nginx
systemctl restart nginx
```

### Updating the Application

```bash
# Pull latest changes
cd /root/Project-Website
git pull

# Rebuild backend
go build -o bin/devpanel cmd/devpanel/main.go
go build -o bin/urlshortener cmd/urlshortener/main.go
go build -o bin/messaging cmd/messaging/main.go

# Rebuild frontend
cd frontend
npm install
npm run build
cd ..

# Restart services
systemctl restart jadenrazo-backend
```

### Backup Strategy

```bash
# Backup databases
pg_dump -U postgres urlshortener > /root/backups/urlshortener_$(date +%Y%m%d).sql

# Backup SQLite databases
cp /root/Project-Website/data/*.db /root/backups/

# Backup configuration
cp -r /root/Project-Website/config /root/backups/config_$(date +%Y%m%d)
```

## Troubleshooting

### Common Issues

1. **Services won't start:**
   - Check logs: `tail -f /root/Project-Website/logs/*.log`
   - Verify permissions: `chmod +x /root/Project-Website/deploy/start-backend.sh`

2. **Database connection issues:**
   - Verify PostgreSQL is running: `systemctl status postgresql`
   - Check database credentials in config files

3. **Nginx 502 Bad Gateway:**
   - Ensure backend services are running: `systemctl status jadenrazo-backend`
   - Check Nginx configuration: `nginx -t`

4. **WebSocket connection failures:**
   - Verify WebSocket proxy settings in Nginx configuration
   - Check if firewall is blocking WebSocket connections

### Support Resources

For further assistance, contact:
- Email: support@jadenrazo.dev
- GitHub Issues: https://github.com/jadenrazo/website/issues 