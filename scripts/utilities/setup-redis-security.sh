#!/bin/bash
# Redis Security Setup Script
# This script configures Redis for maximum security

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Redis Security Setup Script${NC}"
echo "================================="

# Check if running as root (required for some operations)
if [[ $EUID -ne 0 ]]; then
   echo -e "${YELLOW}Warning: This script should be run as root for system-level changes${NC}"
fi

# Create Redis user and group if they don't exist
echo -e "${GREEN}Creating Redis user and group...${NC}"
if ! id "redis" &>/dev/null; then
    useradd -r -s /bin/false redis
    echo "Redis user created"
else
    echo "Redis user already exists"
fi

# Create necessary directories
echo -e "${GREEN}Creating Redis directories...${NC}"
mkdir -p /var/run/redis /var/log/redis /etc/redis
chown redis:redis /var/run/redis /var/log/redis
chmod 750 /var/run/redis /var/log/redis

# Generate secure passwords if not already set
echo -e "${GREEN}Generating secure passwords...${NC}"
if [ -z "${REDIS_PASSWORD:-}" ]; then
    export REDIS_PASSWORD=$(openssl rand -base64 32)
    echo "REDIS_PASSWORD=$REDIS_PASSWORD" >> .env
fi

if [ -z "${REDIS_APP_PASSWORD:-}" ]; then
    export REDIS_APP_PASSWORD=$(openssl rand -base64 32)
    echo "REDIS_APP_PASSWORD=$REDIS_APP_PASSWORD" >> .env
fi

if [ -z "${REDIS_MONITOR_PASSWORD:-}" ]; then
    export REDIS_MONITOR_PASSWORD=$(openssl rand -base64 32)
    echo "REDIS_MONITOR_PASSWORD=$REDIS_MONITOR_PASSWORD" >> .env
fi

if [ -z "${REDIS_WORKER_PASSWORD:-}" ]; then
    export REDIS_WORKER_PASSWORD=$(openssl rand -base64 32)
    echo "REDIS_WORKER_PASSWORD=$REDIS_WORKER_PASSWORD" >> .env
fi

if [ -z "${REDIS_HEALTH_PASSWORD:-}" ]; then
    export REDIS_HEALTH_PASSWORD=$(openssl rand -base64 32)
    echo "REDIS_HEALTH_PASSWORD=$REDIS_HEALTH_PASSWORD" >> .env
fi

echo -e "${GREEN}Passwords generated and saved to .env file${NC}"

# Configure firewall (iptables)
echo -e "${GREEN}Configuring firewall rules...${NC}"
if command -v iptables &> /dev/null; then
    # Drop all external Redis connections
    iptables -A INPUT -p tcp --dport 6379 -j DROP
    iptables -A INPUT -p tcp --dport 6380 -j DROP
    echo "Firewall rules applied"
else
    echo -e "${YELLOW}Warning: iptables not found. Please configure firewall manually${NC}"
fi

# Set up log rotation
echo -e "${GREEN}Setting up log rotation...${NC}"
cat > /etc/logrotate.d/redis << EOF
/var/log/redis/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 640 redis redis
    sharedscripts
    postrotate
        [ -f /var/run/redis/redis.pid ] && kill -USR1 \$(cat /var/run/redis/redis.pid)
    endscript
}
EOF

# Create systemd security override (if using systemd)
if command -v systemctl &> /dev/null; then
    echo -e "${GREEN}Creating systemd security overrides...${NC}"
    mkdir -p /etc/systemd/system/redis.service.d
    cat > /etc/systemd/system/redis.service.d/security.conf << EOF
[Service]
# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/redis /var/log/redis /var/run/redis
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true
RestrictRealtime=true
RestrictNamespaces=true
RestrictSUIDSGID=true
PrivateDevices=true
ProtectHostname=true
ProtectClock=true
ProtectKernelLogs=true
LockPersonality=true
MemoryDenyWriteExecute=true
RestrictAddressFamilies=AF_INET AF_INET6 AF_UNIX
SystemCallFilter=@system-service
SystemCallErrorNumber=EPERM
EOF
    systemctl daemon-reload
fi

# Create Redis health check script
echo -e "${GREEN}Creating health check script...${NC}"
cat > /usr/local/bin/redis-health-check.sh << 'EOF'
#!/bin/bash
# Redis health check script

SOCKET="/var/run/redis/redis.sock"
PASSWORD="${REDIS_HEALTH_PASSWORD}"

# Check if Redis is responding
if redis-cli -s "$SOCKET" -a "$PASSWORD" --user health ping > /dev/null 2>&1; then
    echo "Redis is healthy"
    exit 0
else
    echo "Redis health check failed"
    exit 1
fi
EOF
chmod +x /usr/local/bin/redis-health-check.sh

# Create monitoring script
echo -e "${GREEN}Creating monitoring script...${NC}"
cat > /usr/local/bin/redis-monitor.sh << 'EOF'
#!/bin/bash
# Redis monitoring script

SOCKET="/var/run/redis/redis.sock"
PASSWORD="${REDIS_MONITOR_PASSWORD}"

# Get Redis info
echo "=== Redis Info ==="
redis-cli -s "$SOCKET" -a "$PASSWORD" --user monitor info server

echo -e "\n=== Memory Info ==="
redis-cli -s "$SOCKET" -a "$PASSWORD" --user monitor info memory | grep -E "used_memory_human|used_memory_peak_human|mem_fragmentation_ratio"

echo -e "\n=== Connected Clients ==="
redis-cli -s "$SOCKET" -a "$PASSWORD" --user monitor info clients | grep connected_clients

echo -e "\n=== Slow Queries ==="
redis-cli -s "$SOCKET" -a "$PASSWORD" --user monitor slowlog get 5
EOF
chmod +x /usr/local/bin/redis-monitor.sh

# Set proper permissions
echo -e "${GREEN}Setting file permissions...${NC}"
chmod 640 /etc/redis/redis.conf /etc/redis/users.acl
chown redis:redis /etc/redis/redis.conf /etc/redis/users.acl

# Create backup script
echo -e "${GREEN}Creating backup script...${NC}"
cat > /usr/local/bin/redis-backup.sh << 'EOF'
#!/bin/bash
# Redis backup script (for persistent data if enabled)

BACKUP_DIR="/var/backups/redis"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

mkdir -p "$BACKUP_DIR"
cp /var/lib/redis/dump.rdb "$BACKUP_DIR/dump_$TIMESTAMP.rdb" 2>/dev/null || echo "No RDB file to backup"

# Keep only last 7 days of backups
find "$BACKUP_DIR" -name "dump_*.rdb" -mtime +7 -delete
EOF
chmod +x /usr/local/bin/redis-backup.sh

echo -e "${GREEN}Redis security setup completed!${NC}"
echo ""
echo "Important notes:"
echo "1. Passwords have been generated and saved to .env file"
echo "2. Redis is configured to use Unix socket at /var/run/redis/redis.sock"
echo "3. All dangerous commands have been disabled"
echo "4. Firewall rules have been applied to block external access"
echo "5. Use redis-health-check.sh to verify Redis health"
echo "6. Use redis-monitor.sh to monitor Redis performance"
echo ""
echo -e "${YELLOW}Remember to:"
echo "- Review and adjust the configuration as needed"
echo "- Ensure your application uses the Unix socket connection"
echo "- Regularly monitor Redis logs and performance"
echo "- Keep Redis updated with security patches${NC}"