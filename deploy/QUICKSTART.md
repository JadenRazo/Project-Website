# Quick Start Guide for Jaden Razo Website Backend

This guide provides quick instructions for getting the backend services up and running.

## Development Environment

### 1. Build the Backend Services

```bash
# Build all backend services at once
go build -o bin/devpanel cmd/devpanel/main.go
go build -o bin/urlshortener cmd/urlshortener/main.go

# For messaging, use the simplified implementation if you encounter build errors:
go build -o bin/messaging cmd/messaging/simple_main.go
# Or use the full implementation if it's working:
# go build -o bin/messaging cmd/messaging/main.go
```

### 2. Start Services (Development Mode)

The easiest way to start all services in development mode is to use the provided script:

```bash
# Start all services with the integrated script
./deploy/start-services.sh
```

This will:
- Start the Developer Panel on port 8080
- Start the URL Shortener on port 8081
- Start the Messaging Platform on port 8082 (if available)
- Log all output to the `logs` directory
- Allow you to stop all services with Ctrl+C

### 3. Test the Services

Once running, you can test the services:

- Developer Panel: http://localhost:8080/devpanel
- URL Shortener: http://localhost:8081/s
- Messaging Platform: http://localhost:8082 (if available)

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

If you see "Error: Port XXXX is already in use", find and kill the process:

```bash
# Find the process using the port
sudo lsof -i :8080   # For Developer Panel
sudo lsof -i :8081   # For URL Shortener
sudo lsof -i :8082   # For Messaging Platform

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
1. Check if you're using the simplified implementation (recommended for now)
2. Make sure your frontend is connecting to the correct WebSocket URL
3. Check the logs for any connection errors

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

3. **Configure Backups**: Set up regular database backups using the scripts in `scripts/backup_db.sh` 

4. **Resolve WebSocket Implementation**: If you're using the simplified messaging implementation, plan to resolve the duplicate declarations in the full implementation when time permits 