# Admin Setup Guide

This guide explains how to securely set up admin accounts for your portfolio website.

## Security Overview

The admin panel has multiple security layers:
- **IP Whitelisting**: Only allowed IPs can access `/devpanel`
- **Localhost-only by default**: Binds to 127.0.0.1
- **Rate limiting**: Prevents brute force attempts
- **Domain restrictions**: Only configured domains can be admins
- **Audit logging**: All actions are logged

## Configuration

### 1. Configure Admin Domains

Edit your `.env` file:

```bash
# Set your domain(s) for admin emails
ADMIN_DOMAINS=yourdomain.com

# For multiple domains
ADMIN_DOMAINS=yourdomain.com,yourcompany.dev

# Optionally specify exact allowed emails
ADMIN_EMAILS=admin@yourdomain.com,cto@yourcompany.dev
```

### 2. Configure IP Whitelist (Important!)

By default, only localhost (127.0.0.1) can access the admin panel.

To allow specific IPs, create or edit `backend/config/devpanel-security.yaml`:

```yaml
security:
  ipWhitelist:
    enabled: true
    localhostOnly: false  # Set to false to use custom IPs
    allowedIPs:
      - "127.0.0.1"       # Localhost
      - "::1"             # IPv6 localhost
      - "YOUR.IP.HERE"    # Your static IP
      - "10.0.0.0/24"     # Or CIDR range
```

Or use environment variable:
```bash
IP_WHITELIST=127.0.0.1,YOUR.IP.HERE
```

## Setup Methods

### Method 1: CLI Tool (Recommended for Development)

Most secure - no passwords stored in files:

```bash
# Set your domain
export ADMIN_DOMAINS=yourdomain.com

# Run the interactive admin creation tool
cd backend
go run cmd/admin-create/main.go

# You'll be prompted for email and password
```

### Method 2: SSH Tunnel (Recommended for Production)

Since the admin panel binds to localhost, use SSH tunneling:

```bash
# From your local machine
ssh -L 8081:localhost:8081 your-server.com

# Now access http://localhost:8081/devpanel locally
```

### Method 3: Web Setup Token

1. SSH into your server
2. Run: `curl http://localhost:8081/api/v1/auth/admin/setup/request`
3. Check logs for the setup token
4. Visit the setup page and use the token

### Method 4: Environment Variables (Development Only)

⚠️ **WARNING**: Only for local development. Never commit these!

Instead of storing passwords in .env, use a separate file:

```bash
# Create .env.local (add to .gitignore!)
DEV_ADMIN_EMAIL=admin@yourdomain.com
DEV_ADMIN_PASSWORD=temp-dev-password

# Load it when starting
source .env.local && ./start-dev.sh
```

## Security Best Practices

### For Development
1. Use the CLI tool - don't store passwords
2. If using env vars, use `.env.local` and add to `.gitignore`
3. Use strong passwords even in development
4. Keep IP whitelist enabled

### For Production
1. **Never** use DEV_ADMIN_* variables
2. Always use SSH tunneling or VPN
3. Configure specific IP whitelist
4. Use strong, unique passwords
5. Enable TLS/HTTPS
6. Monitor audit logs: `/var/log/devpanel-audit.log`

## IP Whitelist Examples

### Personal Setup (Home + VPN)
```yaml
security:
  ipWhitelist:
    enabled: true
    localhostOnly: false
    allowedIPs:
      - "127.0.0.1"
      - "YOUR.HOME.IP"
      - "YOUR.VPN.IP"
```

### Team Setup (Office Network)
```yaml
security:
  ipWhitelist:
    enabled: true
    localhostOnly: false
    allowedIPs:
      - "127.0.0.1"
      - "203.0.113.0/24"  # Office subnet
      - "198.51.100.50"   # Specific workstation
```

### Cloud/Docker Setup
```yaml
security:
  ipWhitelist:
    enabled: true
    localhostOnly: false
    allowedIPs:
      - "127.0.0.1"
      - "172.17.0.0/16"   # Docker network
      - "10.0.0.0/8"      # Private network
```

## Troubleshooting

### "403 Forbidden" Error
- Your IP is not whitelisted
- Check your IP: `curl ifconfig.me`
- Add it to the whitelist configuration

### Can't Access from External IP
- By default, binds to localhost only
- Use SSH tunneling or configure proper reverse proxy
- Check `DEVPANEL_BIND` environment variable

### Rate Limit Errors
- Auth endpoints: 5 requests/minute (production)
- Wait before retrying
- Check audit logs for details

## Monitoring

View security events:
```bash
# Audit log
tail -f /var/log/devpanel-audit.log

# Failed login attempts
grep "LOGIN_FAILED" /var/log/devpanel-audit.log

# Blocked IPs
grep "IP_BLOCKED" /var/log/devpanel-audit.log
```