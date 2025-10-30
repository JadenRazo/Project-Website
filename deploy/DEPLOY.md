# Deployment Guide

I've created these step-by-step instructions for deploying my website backend services on an Ubuntu server.

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
# Navigate to backend directory
cd backend

# Install Go dependencies
go mod tidy

# Build all services
mkdir -p bin
go build -o bin/api cmd/api/main.go
go build -o bin/devpanel cmd/devpanel/main.go
go build -o bin/urlshortener cmd/urlshortener/main.go
go build -o bin/messaging cmd/messaging/main.go
go build -o bin/worker cmd/worker/main.go

# Make startup script executable
cd ..
chmod +x deploy/start-backend.sh
chmod +x deploy/start-services.sh
```

## Step 4: Set Up Databases

### PostgreSQL Database

```bash
# Create database and user
sudo -u postgres psql -c "CREATE DATABASE portfolio;"
sudo -u postgres psql -c "CREATE USER portfolio WITH ENCRYPTED PASSWORD 'your_secure_password';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE portfolio TO portfolio;"

# Apply schema
psql -U portfolio -d portfolio -f backend/schema.sql

# Enable PostgreSQL to start on boot
systemctl enable postgresql
systemctl start postgresql
```

## Step 5: Configure Backend Services

Set up environment variables and configuration:

```bash
# Create .env file in backend directory
cd backend
cat > .env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_USER=portfolio
DB_PASSWORD=your_secure_password
DB_NAME=portfolio
JWT_SECRET=your_jwt_secret_here
ADMIN_TOKEN=your_admin_token_here
REDIS_PASSWORD=your_redis_password
ENV=production
EOF

# Ensure configuration files exist
ls config/*.yaml
cd ..
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
cd backend
go build -o bin/api cmd/api/main.go
go build -o bin/devpanel cmd/devpanel/main.go
go build -o bin/urlshortener cmd/urlshortener/main.go
go build -o bin/messaging cmd/messaging/main.go
go build -o bin/worker cmd/worker/main.go
cd ..

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
# Create backup directory
mkdir -p /root/backups

# Backup PostgreSQL database
pg_dump -U portfolio portfolio > /root/backups/portfolio_$(date +%Y%m%d).sql

# Or use the backup script
./scripts/database/backup_db.sh

# Backup configuration
cp -r /root/Project-Website/backend/config /root/backups/config_$(date +%Y%m%d)
cp /root/Project-Website/backend/.env /root/backups/.env_$(date +%Y%m%d)
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