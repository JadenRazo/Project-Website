# Portfolio Website Deployment Guide

## Production Deployment (Recommended)

### Quick Start
```bash
# Build and start production environment (automatic session cleanup included)
./scripts/development/start-prod.sh --rebuild
```

Note: The production script is located at `/scripts/development/start-prod.sh`

### What This Does
- ✅ **Builds optimized React bundle** (minified, tree-shaken, compressed)
- ✅ **Eliminates WebSocket errors** (no hot reload in production)
- ✅ **Faster performance** (smaller bundles, optimized assets)
- ✅ **Professional deployment** (production-ready setup)
- ✅ **Production environment variables** (HTTPS URLs, security headers)

### Key Differences from Development

| Feature | Development (`start-dev.sh`) | Production (`scripts/development/start-prod.sh`) |
|---------|------------------------------|------------------------------|
| Bundle Size | ~2-5MB (unminified) | ~500KB-1MB (minified) |
| Hot Reload | ✅ Enabled (WebSocket) | ❌ Disabled |
| Source Maps | ✅ Generated | ❌ Disabled |
| Dev Tools | ✅ React DevTools | ❌ Disabled |
| Performance | Slower (dev overhead) | ⚡ Optimized |
| Security | Dev headers | 🔒 Production headers |

## Commands

### Production Mode
```bash
# First time or after code changes
./scripts/development/start-prod.sh --rebuild

# Restart existing production build
./scripts/development/start-prod.sh

# Force clean rebuild (automatic session cleanup included)
./scripts/development/start-prod.sh --rebuild

# Stop production services
./scripts/development/stop-prod.sh
```

### Development Mode
```bash
# For active development only
./start-dev.sh --fresh

# Stop development services
./scripts/development/stop-dev.sh
```

## Deployment Checklist

### For Production Use
- [ ] Use `./scripts/development/start-prod.sh`
- [ ] Verify no WebSocket errors in console
- [ ] Check bundle size with `npm run build:analyze`
- [ ] Test HTTPS API endpoints work
- [ ] Verify proper favicon loading

### For Development
- [ ] Use `./start-dev.sh`
- [ ] Hot reload working
- [ ] Source maps available for debugging

### Docker Deployment
```bash
# Start all services with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f [service]

# Stop services
docker-compose down

# Rebuild images
docker-compose build --no-cache
```

## Architecture

```
Production Setup:
┌─────────────────┐    ┌──────────────────┐
│   Static Files  │    │   Go Backend     │
│  (React Build)  │    │  (Microservices) │
│   Port 3000     │    │   Ports 8080+    │
└─────────────────┘    └──────────────────┘
        │                       │
        └───────────────────────┘
                 │
        ┌─────────────────┐
        │   Your Domain   │
        │ jadenrazo.dev   │
        └─────────────────┘
```

## Troubleshooting

### WebSocket Errors
- **Cause**: Running dev server in production
- **Fix**: Use `./scripts/development/start-prod.sh` instead of `./start-dev.sh`

### Bundle Too Large
```bash
# Analyze bundle size
npm run build:analyze

# Check what's included
npx source-map-explorer 'build/static/js/*.js'
```

### Environment Issues
- Check `.env.production` for correct HTTPS URLs
- Verify `env-config.js` is copied to build directory

## Best Practices

1. **Always use production mode for public demos**
2. **Test production build before deploying**
3. **Monitor bundle size** - keep under 1MB when possible
4. **Use HTTPS in production** - never HTTP for public sites
5. **Enable caching headers** for static assets