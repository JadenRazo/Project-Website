# Quick Start Guide for My Website Backend

I've written this guide to provide quick instructions for getting the backend services up and running.

## Development Environment

### 1. Build the Backend Services

```bash
# Navigate to backend directory
cd backend

# Build all backend services
go build -o bin/api cmd/api/main.go
go build -o bin/devpanel cmd/devpanel/main.go
go build -o bin/urlshortener cmd/urlshortener/main.go
go build -o bin/messaging cmd/messaging/main.go
go build -o bin/worker cmd/worker/main.go

cd ..
```

### 2. Start Services (Development Mode)

The easiest way to start all services in development mode is to use the main development script:

```bash
# Start all services with the integrated script
./start-dev.sh --fresh
```

This will:
- Automatically detect and clean up any existing tmux sessions
- Start the Main API on port 8080
- Start the Developer Panel on port 8081
- Start the Messaging Service on port 8082
- Start the URL Shortener on port 8083
- Start the Worker service on port 8084
- Start the Frontend on port 3000
- Create tmux sessions for easy monitoring

**Note:** The script now includes automatic session cleanup, so you don't need to manually kill existing processes.

### 3. Test the Services

Once running, you can test the services:

- Frontend: http://localhost:3000
- Main API: http://localhost:8080
- Developer Panel: http://localhost:8081
- Messaging: http://localhost:8082
- URL Shortener: http://localhost:8083
- Status Page: http://localhost:3000/status

## Production Deployment

### 1. Set Up as a System Service

To ensure the services run continuously:

```bash
# Copy the systemd service file
sudo cp deploy/systemd/jadenrazo-backend.service /etc/systemd/system/

# Reload systemd
sudo systemctl daemon-reload

# Enable the service to start on boot
sudo systemctl enable jadenrazo-backend

# Start the service
sudo systemctl start jadenrazo-backend

# Check status
sudo systemctl status jadenrazo-backend
```

### 2. Set Up Nginx

To expose the services to the web with proper routing:

```bash
# Copy the nginx configuration
sudo cp deploy/nginx/jadenrazo.conf /etc/nginx/sites-available/

# Create a symlink to enable it
sudo ln -s /etc/nginx/sites-available/jadenrazo.conf /etc/nginx/sites-enabled/

# Remove default site if it exists
sudo rm -f /etc/nginx/sites-enabled/default

# Test configuration
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

### 3. Set Up SSL Certificates

For HTTPS support:

```bash
# Install Certbot (if not already installed)
sudo apt install certbot python3-certbot-nginx

# Obtain SSL certificates
sudo certbot --nginx -d jadenrazo.dev -d www.jadenrazo.dev

# Ensure auto-renewal is enabled
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer
```

## Common Tasks

### View Logs

```bash
# View logs for a specific service
tail -f logs/devpanel.log
tail -f logs/urlshortener.log
tail -f logs/messaging.log

# View all logs
tail -f logs/*.log
```

### Restart Services

```bash
# Restart all services through systemd
sudo systemctl restart jadenrazo-backend

# OR manually
./deploy/start-services.sh
```

### Check Status

```bash
# Check status through systemd
sudo systemctl status jadenrazo-backend

# Check running processes
ps aux | grep bin/
```

## Troubleshooting

### Port Already in Use

**Note:** The `./start-dev.sh` and `./start-prod.sh` scripts now automatically detect and clean up existing tmux sessions, preventing most port conflicts. You typically won't need manual intervention.

If you still encounter "Error: Port XXXX is already in use":

```bash
# Find the process using the port
sudo lsof -i :3000   # For Frontend
sudo lsof -i :8080   # For Main API
sudo lsof -i :8081   # For Developer Panel
sudo lsof -i :8082   # For Messaging
sudo lsof -i :8083   # For URL Shortener
sudo lsof -i :8084   # For Worker

# Kill the process
sudo kill -9 <PID>
```

### Services Not Starting

Check the logs for specific error messages:

```bash
# Check the logs
cat logs/devpanel.log
cat logs/urlshortener.log
cat logs/messaging.log
```

### WebSocket Issues

If you encounter issues with the messaging platform's WebSocket functionality:
1. Make sure your frontend is connecting to the correct WebSocket URL (ws://localhost:8082)
2. Check the logs for any connection errors
3. Verify CORS settings in the backend allow WebSocket upgrades

### Database Connection Issues

Make sure PostgreSQL is running:

```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Start PostgreSQL if it's not running
sudo systemctl start postgresql
```

## Next Steps

1. **Configure Frontend**: Make sure your frontend code is properly built and accessible at `/root/Project-Website/frontend/build`

2. **Set Up Monitoring**: Consider setting up monitoring tools like Prometheus and Grafana

3. **Configure Backups**: Set up regular database backups using the scripts in `scripts/database/backup_db.sh` 

4. **Resolve WebSocket Implementation**: If you're using the simplified messaging implementation, plan to resolve the duplicate declarations in the full implementation when time permits 