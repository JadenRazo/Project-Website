# DevPanel Security Guide

This document outlines the comprehensive security features implemented for the DevPanel to ensure localhost-only access and configurable security policies.

## Overview

The DevPanel now includes multiple layers of security:

1. **IP Whitelist Control** - Restrict access to specific IP addresses/ranges
2. **Localhost-Only Binding** - Bind service only to localhost interface
3. **Rate Limiting** - Prevent brute force and DoS attacks
4. **Audit Logging** - Track all admin actions and security events
5. **Enhanced Authentication** - Improved JWT validation with session tracking
6. **Security Headers** - XSS, clickjacking, and other attack prevention

## Security Configuration

### Default Configuration (Development)

```yaml
environment: development

devpanel:
  ip_whitelist:
    enabled: true
    localhost_only: true  # Only allow localhost connections
    allowed_ips: []       # Additional specific IPs (if needed)
    allowed_cidrs: []     # Additional CIDR ranges (if needed)
  bind_address: "127.0.0.1"  # Bind only to localhost
  localhost_only: true
  require_tls: false      # TLS not required in development
  session_timeout: "24h"
  max_sessions: 5

rate_limit:
  enabled: true
  requests_per_minute: 100
  burst_size: 10
  auth_endpoints:
    requests_per_minute: 10  # Stricter limits for auth
    burst_size: 3
```

### Production Configuration

```yaml
environment: production

devpanel:
  ip_whitelist:
    enabled: true
    localhost_only: true
  require_tls: true       # TLS required in production
  session_timeout: "8h"   # Shorter session timeout
  max_sessions: 2         # Fewer concurrent sessions

rate_limit:
  requests_per_minute: 50  # Lower rate limits
  auth_endpoints:
    requests_per_minute: 5
    burst_size: 2
```

## Security Features

### 1. IP Whitelist Middleware

- **Purpose**: Restrict access to authorized IP addresses only
- **Default**: Localhost-only mode enabled
- **Configuration**: Supports individual IPs and CIDR ranges
- **Behavior**: Blocks requests from non-whitelisted IPs with HTTP 403

### 2. Localhost-Only Binding

- **Server Binding**: Service binds to `127.0.0.1:8081` by default
- **Network Isolation**: Prevents external network access to DevPanel
- **Override**: Can be configured for specific deployment scenarios

### 3. Rate Limiting

- **Global Limits**: 100 requests/minute (dev), 50 requests/minute (prod)
- **Auth Limits**: 10 login attempts/minute (dev), 5 attempts/minute (prod)
- **Protection**: Prevents brute force attacks and resource exhaustion
- **Cleanup**: Automatic cleanup of rate limiter state

### 4. Audit Logging

- **Authentication Events**: Login attempts, token validation failures
- **Security Events**: Blocked IPs, rate limit violations
- **Admin Actions**: All management operations (create/edit/delete)
- **Log Rotation**: Automatic log rotation with size/age limits

### 5. Enhanced Authentication

- **JWT Validation**: Improved token validation with IP tracking
- **Session Context**: User and session information in request context
- **Security Headers**: CSRF, XSS, and clickjacking protection
- **Audit Trail**: All auth events logged with client IP and user agent

### 6. CORS Security

- **Localhost Origins**: Restricted to localhost in localhost-only mode
- **Credential Handling**: Secure credential handling
- **Method Restrictions**: Only necessary HTTP methods allowed

## Environment Variables

```bash
# Security Configuration
ENVIRONMENT=development|production

# JWT Secret (required)
JWT_SECRET=your-secure-jwt-secret

# Audit Log Path (optional)
AUDIT_LOG_PATH=/var/log/devpanel-audit.log

# Additional IP Whitelist (comma-separated)
DEVPANEL_ALLOWED_IPS=127.0.0.1,::1

# Additional CIDR Ranges (comma-separated)
DEVPANEL_ALLOWED_CIDRS=192.168.1.0/24,10.0.0.0/8
```

## Usage Examples

### Allow Specific IP Address

To allow a specific IP address to access the DevPanel:

1. **Environment Variable**:
```bash
DEVPANEL_ALLOWED_IPS=127.0.0.1,192.168.1.100
```

2. **Configuration File**:
```yaml
devpanel:
  ip_whitelist:
    enabled: true
    localhost_only: false
    allowed_ips:
      - "192.168.1.100"
      - "10.0.1.50"
```

### Allow CIDR Range

To allow an entire network range:

```yaml
devpanel:
  ip_whitelist:
    enabled: true
    localhost_only: false
    allowed_cidrs:
      - "192.168.1.0/24"  # Entire local network
      - "10.0.0.0/8"      # Private network range
```

### Production Deployment

For production deployment with TLS and strict security:

```yaml
environment: production
devpanel:
  localhost_only: true      # Always localhost-only in production
  require_tls: true        # Require HTTPS
  session_timeout: "4h"    # Short session timeout
  max_sessions: 1          # Single session per user
```

## Security Monitoring

### Audit Log Analysis

Monitor these security events in audit logs:

- `login_attempt` - Failed login attempts
- `auth_invalid_token` - Invalid JWT tokens
- `auth_insufficient_privileges` - Privilege escalation attempts
- `ip_not_whitelisted` - Blocked IP addresses
- `rate_limit_exceeded` - Rate limit violations

### Security Status Endpoint

Check security status via: `GET /api/v1/devpanel/security/status`

Returns:
```json
{
  "ip_whitelist": {
    "enabled": true,
    "localhost_only": true,
    "allowed_ips": [],
    "allowed_cidrs": []
  },
  "rate_limiter": {
    "active_limiters": 5,
    "enabled": true,
    "requests_per_min": 100,
    "burst_size": 10
  },
  "audit_log": {
    "enabled": true
  },
  "environment": "development"
}
```

## Best Practices

1. **Always use localhost-only mode** in production
2. **Enable TLS** for production deployments
3. **Monitor audit logs** for security events
4. **Rotate JWT secrets** regularly
5. **Use strong passwords** for admin accounts
6. **Limit session duration** in production
7. **Review IP whitelist** regularly

## Troubleshooting

### Access Denied Issues

1. Check if IP is whitelisted
2. Verify localhost-only mode settings
3. Review audit logs for blocked requests
4. Check rate limiting status

### Authentication Problems

1. Verify JWT_SECRET is set correctly
2. Check token expiration
3. Review audit logs for auth failures
4. Ensure proper headers are sent

### Configuration Issues

1. Validate YAML syntax in config files
2. Check environment variable values
3. Verify file permissions for audit logs
4. Test with development configuration first

## Upgrade Guide

To upgrade existing DevPanel installations:

1. **Backup existing configuration**
2. **Update main.go** to use new security features
3. **Add security configuration files**
4. **Set required environment variables**
5. **Test in development environment**
6. **Deploy to production with localhost-only mode**

The security enhancements are backward compatible and will default to secure settings if not configured explicitly.