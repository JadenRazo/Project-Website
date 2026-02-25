# Personal Portfolio & Project Showcase

## Overview

My full-stack portfolio website with 5 independent Go microservices. It includes a modern React frontend and a Go-based backend architecture with separate microservice processes for URL shortening, messaging, developer panel, visitor analytics, and background workers.

## Architecture

![Architecture](docs/architecture.svg)

## System Requirements

### Development Environment
- Go 1.21 or higher
- Node.js 18.x or higher
- npm 9.x or higher
- Git
- tmux (for development session management)
- PostgreSQL 14+ (for database)
- tokei (optional, for code statistics)
- Make (optional, for using Makefile commands)

### Production Server Requirements
- Ubuntu 22.04 LTS or higher (recommended)
- 2+ CPU cores
- 4GB+ RAM
- 50GB+ storage
- Nginx (for reverse proxy)
- PostgreSQL 14+ (for production database)
- Redis (for caching)

## Project Structure

```
Project-Website/
├── frontend/             # React/TypeScript frontend
│   ├── build/            # Production build
│   ├── public/           # Essential public assets only
│   │   ├── index.html    # HTML entry point
│   │   ├── manifest.json # PWA manifest
│   │   └── robots.txt    # SEO configuration
│   ├── src/
│   │   ├── assets/       # Organized assets
│   │   │   ├── data/     # JSON data files
│   │   │   │   └── code_stats.json
│   │   │   ├── icons/    # Favicons and icons
│   │   │   │   ├── favicon.ico
│   │   │   │   ├── favicon-16x16.png
│   │   │   │   ├── favicon-32x32.png
│   │   │   │   └── apple-touch-icon.png
│   │   │   ├── images/   # Image assets
│   │   │   │   └── headshot.jpg
│   │   │   └── videos/   # Video assets
│   │   │       └── web_ready_quizbot_example_video.mp4
│   │   ├── config/       # Configuration files
│   │   │   └── env-config.js
│   │   ├── components/   # UI components
│   │   │   ├── animations/  # Animation components
│   │   │   ├── auth/        # Authentication components
│   │   │   ├── common/      # Shared common components
│   │   │   ├── devpanel/    # Developer panel components
│   │   │   │   ├── AdminLogin.tsx
│   │   │   │   ├── ProjectManager.tsx
│   │   │   │   ├── CertificationsManager.tsx
│   │   │   │   ├── SkillsManager.tsx
│   │   │   │   ├── PromptsManager.tsx
│   │   │   │   ├── ProjectPathsManager.tsx
│   │   │   │   └── VisitorAnalytics.tsx
│   │   │   ├── Footer/      # Footer components
│   │   │   │   └── Footer.tsx
│   │   │   ├── layout/      # Layout components
│   │   │   │   └── Layout.tsx
│   │   │   ├── metrics/     # System metrics components
│   │   │   │   └── SystemMetrics.tsx
│   │   │   ├── NavigationBar/  # Navigation bar
│   │   │   │   └── NavigationBar.tsx
│   │   │   ├── navigation/  # Navigation utilities
│   │   │   │   ├── ScrollToTop.tsx
│   │   │   │   ├── PageTop.tsx
│   │   │   │   └── PageTransition.tsx
│   │   │   ├── notifications/  # Notification system
│   │   │   ├── sections/    # Page sections
│   │   │   │   ├── Hero.tsx
│   │   │   │   └── Projects.tsx
│   │   │   ├── skeletons/   # Loading skeletons
│   │   │   ├── ThemeToggle/ # Theme switching
│   │   │   └── ui/          # UI elements
│   │   ├── contexts/        # React contexts
│   │   │   ├── ThemeContext.tsx  # Theme management
│   │   │   └── ZIndexContext.tsx  # Z-index management
│   │   ├── hooks/           # Custom hooks
│   │   │   ├── useAnimationController.ts  # Animation control
│   │   │   ├── useAuth.tsx  # Authentication
│   │   │   ├── useClickOutside.ts  # Click detection
│   │   │   ├── useCrudOperations.ts  # CRUD operations
│   │   │   ├── useDebounce.ts  # Debouncing utility
│   │   │   ├── useDeviceCapabilities.ts  # Device features
│   │   │   ├── useInlineFormScroll.ts  # Form scroll management
│   │   │   ├── useLazyLoad.ts  # Lazy loading
│   │   │   ├── useLoadingState.ts  # Loading state management
│   │   │   ├── useMemoryManager.ts  # Memory optimization
│   │   │   ├── useMobileOptimizations.ts  # Mobile performance
│   │   │   ├── useModalScroll.ts  # Modal scroll handling
│   │   │   ├── useOptimizedScrollHandler.ts  # Optimized scrolling
│   │   │   ├── usePerformanceOptimizations.ts  # Performance
│   │   │   ├── usePreloader.ts  # Resource preloading
│   │   │   ├── useScrollTo.ts  # Scroll navigation
│   │   │   ├── useScrollToForm.ts  # Form scroll utilities
│   │   │   ├── useTheme.ts  # Theme access
│   │   │   ├── useThemeToggle.ts  # Theme switching
│   │   │   ├── useTouchInteractions.ts  # Touch gestures
│   │   │   ├── useWebSocketConnection.ts  # WebSocket management
│   │   │   └── useZIndex.ts  # Z-index utilities
│   │   ├── pages/           # App pages
│   │   │   ├── About/       # About page
│   │   │   │   └── About.tsx
│   │   │   ├── Contact/     # Contact page
│   │   │   │   └── Contact.tsx
│   │   │   ├── devpanel/    # Developer panel
│   │   │   │   └── DevPanel.tsx
│   │   │   ├── Home/        # Homepage
│   │   │   │   └── Home.tsx
│   │   │   ├── messaging/   # Messaging app
│   │   │   │   └── Messaging.tsx
│   │   │   ├── NotFound/    # 404 page
│   │   │   │   └── NotFound.tsx
│   │   │   ├── Projects/    # Projects page
│   │   │   │   └── index.tsx
│   │   │   ├── Status/      # Status page
│   │   │   │   └── Status.tsx
│   │   │   └── urlshortener/ # URL Shortener
│   │   │       └── UrlShortener.tsx
│   │   ├── stores/          # Zustand state stores
│   │   │   ├── authStore.ts  # Authentication state
│   │   │   ├── performanceStore.ts  # Performance metrics
│   │   │   └── themeStore.ts  # Theme state
│   │   ├── styles/          # Styling
│   │   │   ├── GlobalStyles.ts  # Global styles
│   │   │   ├── theme.types.ts  # Theme types
│   │   │   └── themes.ts  # Theme configs
│   │   ├── utils/           # Utilities
│   │   │   ├── apiConfig.ts  # API configuration
│   │   │   ├── debugHelpers.ts  # Debugging tools
│   │   │   ├── devCacheManager.ts  # Development cache
│   │   │   ├── devPanelApi.ts  # DevPanel API client
│   │   │   ├── errorHandler.ts  # Error handling
│   │   │   ├── lazyWithPreload.ts  # Lazy loading with preload
│   │   │   ├── MemoryManager.tsx  # Memory management
│   │   │   ├── performance.ts  # Performance tools
│   │   │   ├── performanceConfig.ts  # Performance configuration
│   │   │   ├── performanceMonitor.ts  # Performance monitoring
│   │   │   ├── preloader.ts  # Resource preloader
│   │   │   ├── promptApi.ts  # Prompts API client
│   │   │   ├── scrollConfig.ts  # Scroll configuration
│   │   │   ├── scrollTestUtils.ts  # Scroll testing utilities
│   │   │   └── validation.ts  # Input validation
│   │   ├── App.tsx          # Main App component
│   │   └── index.tsx        # Entry point
│   └── package.json         # Dependencies
│
├── backend/                  # Go backend
│   ├── cmd/                  # Entry points
│   │   ├── admin/            # Admin CLI tools
│   │   ├── admin-create/     # Create admin user
│   │   ├── api/              # Main API service
│   │   │   └── main.go
│   │   ├── check-db/         # Database health check
│   │   ├── devpanel/         # Developer panel
│   │   │   └── main.go
│   │   ├── init-admin/       # Initialize admin
│   │   ├── messaging/        # Messaging service
│   │   │   └── main.go
│   │   ├── server/           # Server utilities
│   │   ├── simple-api/       # Simple API variant
│   │   ├── urlshortener/     # URL shortener service
│   │   │   └── main.go
│   │   ├── verify-admin/     # Verify admin credentials
│   │   └── worker/           # Background worker
│   │       └── main.go
│   ├── config/               # Configuration
│   │   ├── app.yaml          # Main config
│   │   ├── development.yaml  # Dev config
│   │   ├── production.yaml   # Prod config
│   │   └── config.go         # Config loader
│   ├── deployments/          # Deployment
│   │   ├── docker/           # Docker setup
│   │   │   └── docker-compose.yml
│   │   ├── nginx/            # Web server
│   │   │   └── api.conf
│   │   └── systemd/          # Service defs
│   │       └── api.service
│   ├── internal/             # Packages
│   │   ├── app/              # Bootstrap
│   │   │   └── server/       # HTTP server
│   │   │       └── middleware/
│   │   ├── codestats/        # Code statistics
│   │   │   ├── delivery/http/
│   │   │   ├── projectpath/  # Project path management
│   │   │   └── service.go
│   │   ├── common/           # Shared utilities
│   │   │   ├── auth/         # Authentication (JWT, password, admin)
│   │   │   ├── cache/        # Caching (Redis with TLS)
│   │   │   ├── circuitbreaker/  # Circuit breaker pattern
│   │   │   ├── compression/  # Response compression
│   │   │   ├── metrics/      # Prometheus metrics
│   │   │   ├── middleware/   # HTTP middleware
│   │   │   ├── ratelimit/    # Rate limiting
│   │   │   ├── repository/   # Base repository patterns
│   │   │   ├── response/     # Standard API responses
│   │   │   ├── security/     # Security utilities
│   │   │   ├── testutil/     # Testing utilities
│   │   │   ├── utils/        # General utilities
│   │   │   └── validator/    # Input validation
│   │   ├── core/             # Core services
│   │   │   ├── config/       # Configuration management
│   │   │   ├── db/           # Database access
│   │   │   └── service_manager.go
│   │   ├── devpanel/         # Developer panel
│   │   │   ├── certification/  # Certifications management
│   │   │   ├── prompt/       # AI prompts management
│   │   │   ├── server/       # DevPanel server
│   │   │   ├── skill/        # Skills management
│   │   │   └── service.go
│   │   ├── domain/           # Business domain
│   │   │   ├── entity/       # Core entities
│   │   │   └── errors/       # Domain errors
│   │   ├── gateway/          # API gateway
│   │   ├── messaging/        # Messaging service
│   │   │   ├── api/          # API handlers
│   │   │   ├── delivery/     # HTTP/WebSocket delivery
│   │   │   ├── entity/       # Message entities
│   │   │   ├── events/       # Event types
│   │   │   ├── service/      # Business logic
│   │   │   ├── usecase/      # Use cases
│   │   │   └── websocket/    # WebSocket management
│   │   ├── projects/         # Projects service
│   │   │   ├── delivery/http/
│   │   │   ├── repository/
│   │   │   └── service/
│   │   ├── status/           # Status monitoring
│   │   │   └── service.go
│   │   ├── urlshortener/     # URL shortener
│   │   │   ├── entity/       # URL entities
│   │   │   └── service.go
│   │   ├── visitor/          # Visitor analytics
│   │   │   └── service.go
│   │   └── worker/           # Background workers
│   │       └── tasks/        # Task definitions
│   ├── schema.sql            # Complete database schema
│   ├── scripts/              # Scripts
│   │   ├── run.sh            # Run app
│   │   ├── seed.sh           # Seed database
│   │   ├── seed_projects.go  # Seed projects data
│   │   └── setup.sh          # Setup environment
│   ├── go.mod                # Go dependencies
│   └── go.sum                # Go checksums
├── scripts/                  # Root scripts directory
│   ├── database/             # Database management
│   │   ├── backup_db.sh      # Database backup
│   │   ├── check-db-status.sh  # Check DB status
│   │   ├── restore_db.sh     # Database restore
│   │   ├── seed.sh           # Seed data
│   │   ├── update-database.sh  # Update schema
│   │   └── update-db.sh      # Safe schema updates
│   ├── development/          # Development tools
│   │   ├── lint.sh           # Backend linting
│   │   ├── optimize-images.sh  # Image optimization
│   │   ├── setup.sh          # Development setup
│   │   ├── start-prod.sh     # Production start
│   │   ├── stop-dev.sh       # Stop development
│   │   └── stop-prod.sh      # Stop production
│   ├── monitoring/           # Monitoring tools
│   │   ├── setup_code_stats_cron.sh  # Code stats cron
│   │   └── update_code_stats.sh  # Update code stats
│   ├── ssl/                  # SSL management
│   │   ├── fix-ssl-now.sh    # SSL fix
│   │   ├── setup-ssl-automation.sh  # SSL automation
│   │   ├── ssl-auto-renew.sh  # Auto renewal
│   │   └── ssl-monitor.sh    # SSL monitoring
│   ├── setup-admin.sh        # Initialize admin user
│   └── setup-redis-security.sh  # Configure Redis TLS/ACL
├── deploy/                   # Deployment files
│   ├── grafana/              # Grafana dashboards
│   ├── nginx/                # Nginx configurations
│   ├── prometheus/           # Prometheus configs
│   ├── redis/                # Redis secure configs
│   ├── start-backend.sh      # Start backend services
│   └── start-services.sh     # Start all services
├── docker-compose.yml        # Docker deployment
├── Makefile                  # Development commands
└── README.md                 # Documentation
```

## Local Development Setup

### Quick Start

```bash
# Clone the repository
git clone https://github.com/JadenRazo/Project-Website.git
cd Project-Website

# Start development environment with all services
./start-dev.sh --fresh

# Alternative: Start without clearing caches
./start-dev.sh
```

The `start-dev.sh` script will:
- Automatically detect and clean up any existing tmux sessions
- Check all required dependencies
- Set up PostgreSQL database connections
- Initialize the database schema
- Start all 5 backend microservices in separate tmux windows (api, devpanel, messaging, urlshortener, worker)
- Start the frontend development server in a separate tmux session
- Create organized tmux sessions for easy monitoring (portfolio-frontend, portfolio-backend)

**Note:** The script now includes automatic session cleanup, eliminating the need to manually kill existing processes or use the `--kill-existing` flag.

### Windows Setup

1. Install Prerequisites:
   ```powershell
   # Install Chocolatey (if not installed)
   Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

   # Install Go
   choco install golang

   # Install Node.js
   choco install nodejs

   # Install Git and tmux
   choco install git
   choco install tmux

   # Install PostgreSQL
   choco install postgresql

   # Optional: Install tokei for code statistics
   cargo install tokei
   ```

2. Clone and start:
   ```powershell
   git clone https://github.com/JadenRazo/Project-Website.git
   cd Project-Website
   ./start-dev.sh --fresh
   ```

### Ubuntu/Linux Setup

1. Install Prerequisites:
   ```bash
   # Update system
   sudo apt update && sudo apt upgrade -y

   # Install Go
   wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc

   # Install Node.js
   curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
   sudo apt install -y nodejs

   # Install Git, tmux, and PostgreSQL
   sudo apt install -y git tmux postgresql postgresql-contrib

   # Optional: Install tokei for code statistics
   cargo install tokei
   ```

2. Setup PostgreSQL:
   ```bash
   # Start PostgreSQL
   sudo systemctl start postgresql
   sudo systemctl enable postgresql

   # Create database user
   sudo -u postgres createuser --interactive
   ```

3. Clone and start:
   ```bash
   git clone https://github.com/JadenRazo/Project-Website.git
   cd Project-Website
   ./start-dev.sh --fresh
   ```

## Service Ports

- Frontend Development: http://localhost:3000
- Main API: http://localhost:8080
- DevPanel: http://localhost:8081
- Messaging Service: http://localhost:8082
- URL Shortener: http://localhost:8083
- Worker Service: http://localhost:8084 (Background tasks)

## Key Features

- **Status Monitoring**: Real-time service health monitoring at http://localhost:3000/status
- **Visitor Analytics**: Real-time visitor tracking with geographic data and behavior metrics
- **Code Statistics**: Automatic code line counting with tokei (updates hourly)
- **URL Shortener**: Create and track short URLs with analytics
- **Messaging System**: Real-time messaging with WebSocket support and presence tracking
- **Developer Panel**: Comprehensive admin interface with project management, certifications, skills, prompts, and visitor analytics
- **Performance Monitoring**: Prometheus metrics and Grafana dashboards
- **Secure Infrastructure**: Redis with TLS/ACL, rate limiting, circuit breakers

## Available Scripts

### Development
```bash
# RECOMMENDED: Start all services with fresh cache clearing
./start-dev.sh --fresh

# Start all services (normal mode)
./start-dev.sh

# Start with additional options
./start-dev.sh --skip-deps       # Skip dependency installation
./start-dev.sh --verbose         # Enable verbose output
./start-dev.sh --help           # Show help message

# Stop all services
./scripts/development/stop-dev.sh

# Production mode
./start-prod.sh                  # Start production services with automatic cleanup

# Legacy commands (still available)
npm run dev                              # Start all services
npm run dev:frontend                     # Start only frontend
npm run dev:backend                      # Start only backend
cd frontend && npm run dev:fresh         # Start frontend with cleared caches
cd frontend && npm run clear-cache       # Clear all frontend development caches
```

### Production
```bash
# Build all services
npm run build

# Start all services in production mode
npm run start
```

### Individual Service Management
```bash
# Using the start-services script
./deploy/start-services.sh
```

## Database Setup

The project uses PostgreSQL with a comprehensive schema including URL shortener, messaging, and monitoring tables.

### Initial Setup
```bash
# For new installations - applies complete schema
./scripts/database/update-db.sh

# For existing databases - safe migration that preserves data
./scripts/database/update-db.sh --safe

# Check database status
./scripts/database/check-db-status.sh
```

### Manual Setup
```bash
# Create database if it doesn't exist
createdb project_website

# Apply schema
psql -d project_website -f backend/schema.sql
```

### What's Included
The database schema includes tables for:
- URL shortener (shortened_urls, url_clicks, url_tags)
- Status monitoring (incidents, status_history)
- Messaging (channels, messages, word_filters)
- User management and authentication
- Audit logging and metrics

## Environment Configuration

1. Frontend (.env):
   ```
   REACT_APP_API_URL=http://localhost:8080
   REACT_APP_WS_URL=ws://localhost:8082
   ```

2. Backend (.env):
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=project_website
   JWT_SECRET=your_jwt_secret
   ```

## Deployment

### Manual Deployment

1. Build frontend:
   ```bash
   cd frontend
   npm run build
   ```

2. Build backend services:
   ```bash
   cd backend
   go build -o bin/api cmd/api/main.go
   go build -o bin/devpanel cmd/devpanel/main.go
   go build -o bin/messaging cmd/messaging/main.go
   go build -o bin/urlshortener cmd/urlshortener/main.go
   go build -o bin/worker cmd/worker/main.go
   ```

   Each service is an independent binary that runs as a separate process.

3. Configure Nginx:
   ```nginx
   server {
       listen 80;
       server_name yourdomain.com;

       location / {
           root /path/to/frontend/build;
           try_files $uri $uri/ /index.html;
       }

       location /api {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

## Docker Deployment

The application runs as 10 containers: 5 backend microservices, frontend, PostgreSQL, Redis, Prometheus, and Grafana.

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- 4GB RAM minimum
- 10GB disk space

Verify prerequisites:
```bash
docker --version          # Should show 20.10 or higher
docker-compose --version  # Should show 2.0 or higher
```

### Container Architecture

| Container | Port | Purpose |
|-----------|------|---------|
| portfolio_api | 8080 | Main API gateway, authentication, visitor tracking |
| portfolio_devpanel | 8081 | Admin panel for service management and analytics |
| portfolio_messaging | 8082 | WebSocket real-time messaging and presence |
| portfolio_urlshortener | 8083 | URL shortening with click analytics |
| portfolio_worker | 8084 | Background cron jobs for visitor analytics aggregation |
| portfolio_frontend | 3000 | React SPA served via nginx |
| postgres | 5432 | Primary database |
| redis_secure | 6379 | Cache with TLS 1.3 encryption |
| prometheus | 9090 | Metrics collection |
| grafana | 3001 | Metrics visualization |

### Setup

**Step 1: Generate secure secrets**

Run the automated secret generation script:
```bash
./scripts/utilities/generate-secrets.sh
```

This creates `.env` from `.env.example` and generates secure passwords for all services.

**Step 2: Generate Redis TLS certificates**

```bash
./scripts/utilities/generate-redis-certs.sh
```

Creates self-signed TLS certificates in `deploy/redis/certs/` (ca.crt, ca.key, redis.crt, redis.key). Replace with CA-signed certificates for production.

**Step 3: Configure environment**

Edit `.env` and verify these required variables are set:
```bash
DB_PASSWORD=your_secure_password
JWT_SECRET=your_jwt_secret           # 32+ characters
ADMIN_TOKEN=your_admin_token         # openssl rand -base64 32
REDIS_APP_PASSWORD=your_redis_app_password
REDIS_MONITOR_PASSWORD=your_redis_monitor_password
REDIS_WORKER_PASSWORD=your_redis_worker_password
REDIS_HEALTH_PASSWORD=your_redis_health_password
GRAFANA_PASSWORD=your_grafana_password
ADMIN_ALLOWED_EMAILS=your-email@example.com,@your-domain.com
FRONTEND_URL=http://localhost:3000
```

**Step 4: Build and start**

```bash
docker-compose build --no-cache
docker-compose up -d
```

**Step 5: Verify deployment**

```bash
docker-compose ps
```

All services should show "healthy" or "Up" status. First deployment may take 5-10 minutes for database initialization.

### Health Checks

```bash
curl http://localhost:8080/api/health   # API
curl http://localhost:8081/api/health   # DevPanel
curl http://localhost:8082/api/health   # Messaging
curl http://localhost:8083/api/health   # URL Shortener
curl http://localhost:8084/api/health   # Worker
curl http://localhost:3000/health       # Frontend

# PostgreSQL
docker-compose exec postgres pg_isready -U portfolio

# Redis (with TLS)
docker-compose exec redis redis-cli --tls --cacert /etc/redis/certs/ca.crt \
  --user health -a "$REDIS_HEALTH_PASSWORD" ping

# Prometheus
curl http://localhost:9090/-/healthy
```

### Database Verification

```bash
# Check all 48 tables were created
docker exec -it portfolio_postgres psql -U postgres -d project_website -c "\dt" | wc -l

# Verify visitor analytics tables
docker exec -it portfolio_postgres psql -U postgres -d project_website \
  -c "SELECT table_name FROM information_schema.tables WHERE table_name LIKE 'visitor%';"
```

Expected: 7 visitor tables (visitor_sessions, page_views, visitor_realtime, visitor_metrics, visitor_daily_summary, privacy_consents, visitor_locations).

### Redis Security

Redis is configured with TLS 1.3 encryption and ACL authentication with four user roles:

- `app` - Full read/write access (api, messaging, urlshortener)
- `monitor` - Read-only access (devpanel)
- `worker` - Background job access (worker)
- `health` - Health check ping only

Dangerous commands (FLUSHDB, KEYS, CONFIG) are disabled.

Connection string format:
```
rediss://username:password@redis:6379?tls=true&ca_cert=/path/to/ca.crt
```

### Docker Networks

Three isolated networks:

1. **internal_network** - PostgreSQL and Redis (no external access)
2. **external_network** - Frontend, API services, and nginx
3. **monitoring_network** - Prometheus and Grafana (isolated)

### Volumes

- `postgres_data` - PostgreSQL database
- `redis_data` - Redis persistence
- `prometheus_data` - Metrics history
- `grafana_data` - Dashboards and configuration

### Resource Limits

| Service | CPU Limit | Memory Limit | CPU Reserved | Memory Reserved |
|---------|-----------|--------------|--------------|-----------------|
| PostgreSQL | 1.0 | 1GB | 0.5 | 512MB |
| Redis | 0.5 | 512MB | 0.25 | 256MB |
| API | 1.0 | 512MB | 0.5 | 256MB |
| DevPanel | 0.5 | 256MB | 0.25 | 128MB |
| Messaging | 0.5 | 384MB | 0.25 | 192MB |
| URL Shortener | 0.5 | 256MB | 0.25 | 128MB |
| Worker | 0.5 | 256MB | 0.25 | 128MB |
| Frontend | 0.5 | 256MB | 0.25 | 128MB |
| Prometheus | 0.5 | 512MB | 0.25 | 256MB |
| Grafana | 0.5 | 512MB | 0.25 | 256MB |

For production under expected traffic, increase limits by 50-100%.

### Common Docker Operations

```bash
# View logs
docker-compose logs -f                     # All services
docker-compose logs -f api                 # Specific service
docker-compose logs --tail=100 api         # Last 100 lines

# Restart services
docker-compose restart                     # All services
docker-compose restart api                 # Specific service

# Update application
git pull origin main
docker-compose build --no-cache
docker-compose up -d

# Clean up
docker-compose down                        # Stop services
docker-compose down -v                     # Stop and remove volumes (WARNING: deletes data)
docker image prune -a                      # Remove unused images

# Rebuild single service
docker-compose build api
docker-compose up -d api
```

### Docker Troubleshooting

**Container fails to start:**
```bash
docker-compose logs <service-name>
```
Common causes: missing environment variables, port conflicts (`lsof -i :<port>`), database not healthy.

**Redis TLS connection failed:**
```bash
ls -la deploy/redis/certs/   # Verify ca.crt, ca.key, redis.crt, redis.key exist
./scripts/utilities/generate-redis-certs.sh  # Regenerate if missing
docker-compose restart redis
```

**Database schema not initialized:**
```bash
docker-compose down
docker volume rm project-website_postgres_data
docker-compose up -d postgres
# Wait for PostgreSQL to be healthy, then:
docker-compose up -d
```

**Services show "unhealthy" after startup:**
Wait 60 seconds - health checks run every 30 seconds and may require 2-3 cycles. If still unhealthy after 2 minutes, check logs and restart the service.

**Worker cron jobs not running:**
```bash
docker-compose logs worker | grep "Starting scheduled tasks"
docker-compose ps worker
```

**Rollback if deployment fails:**
```bash
docker-compose down                 # Stop services (volumes preserved)
# Fix configuration or restore backup
docker-compose up -d                # Restart
```

## OAuth Setup

The DevPanel supports OAuth authentication via Google, GitHub, and Microsoft. OAuth is optional - username/password authentication works without it.

### Callback URLs

Configure these redirect URIs in each provider's developer console:

| Provider | Development | Production |
|----------|-------------|------------|
| Google | `http://localhost:8080/api/v1/auth/admin/oauth/callback/google` | `https://your-domain.com/api/v1/auth/admin/oauth/callback/google` |
| GitHub | `http://localhost:8080/api/v1/auth/admin/oauth/callback/github` | `https://your-domain.com/api/v1/auth/admin/oauth/callback/github` |
| Microsoft | `http://localhost:8080/api/v1/auth/admin/oauth/callback/microsoft` | `https://your-domain.com/api/v1/auth/admin/oauth/callback/microsoft` |

### Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/) and create a project
2. Navigate to **APIs & Services** → **Library**, search for and enable "Google+ API"
3. Go to **APIs & Services** → **OAuth consent screen**, select External, fill in app details, and add scopes `.../auth/userinfo.email` and `.../auth/userinfo.profile`
4. Go to **APIs & Services** → **Credentials**, click **Create Credentials** → **OAuth client ID**
5. Select Web application, add authorized JavaScript origins (`http://localhost:3000`) and authorized redirect URIs
6. Copy the Client ID and Client Secret

```bash
# .env
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/admin/oauth/callback/google
```

### GitHub OAuth

1. Go to [GitHub Developer Settings](https://github.com/settings/developers) → **OAuth Apps** → **New OAuth App**
2. Set Homepage URL to `http://localhost:3000` and Authorization callback URL to the GitHub callback above
3. Click Register application, then generate a client secret
4. Copy the Client ID and Client Secret immediately (the secret cannot be viewed again)

```bash
# .env
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://localhost:8080/api/v1/auth/admin/oauth/callback/github
```

### Microsoft OAuth

1. Go to [Azure Portal](https://portal.azure.com/) → **Azure Active Directory** → **App registrations** → **New registration**
2. Set supported account types to "Accounts in any organizational directory and personal Microsoft accounts"
3. Add the Microsoft redirect URI under Redirect URI → Web
4. Go to **Certificates & secrets** → **New client secret**, copy the Value immediately
5. Go to **API permissions**, add Microsoft Graph delegated permissions: `User.Read`, `email`, `profile`, `openid`
6. Copy the Application (client) ID from Overview

```bash
# .env
MICROSOFT_CLIENT_ID=your-application-client-id
MICROSOFT_CLIENT_SECRET=your-client-secret-value
MICROSOFT_REDIRECT_URL=http://localhost:8080/api/v1/auth/admin/oauth/callback/microsoft
```

### OAuth .env Configuration

```bash
FRONTEND_URL=http://localhost:3000
ADMIN_ALLOWED_EMAILS=your-email@gmail.com,@your-domain.com
OAUTH_ENCRYPTION_KEY=<generated by generate-secrets.sh>

GOOGLE_CLIENT_ID=xxxxxx.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-xxxxxxxxxx
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/admin/oauth/callback/google

GITHUB_CLIENT_ID=Iv1.xxxxxxxxxxxxxx
GITHUB_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxx
GITHUB_REDIRECT_URL=http://localhost:8080/api/v1/auth/admin/oauth/callback/github

MICROSOFT_CLIENT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
MICROSOFT_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxx
MICROSOFT_REDIRECT_URL=http://localhost:8080/api/v1/auth/admin/oauth/callback/microsoft
```

`ADMIN_ALLOWED_EMAILS` accepts specific addresses and domain wildcards (e.g., `@company.com` allows anyone from that domain). Unconfigured OAuth providers will not show buttons in the UI.

### OAuth Troubleshooting

- **Redirect URI mismatch**: The URI in your OAuth app must match exactly what is in `.env`, including protocol, port, and path
- **Unauthorized / Access Denied**: Verify your email is listed in `ADMIN_ALLOWED_EMAILS` and the OAuth encryption key is set
- **Invalid client**: Double-check that `CLIENT_ID` and `CLIENT_SECRET` have no extra spaces or newlines
- **Google "This app isn't verified"**: Expected in development - click Advanced and proceed. Submit for verification for production

### OAuth Security Practices

- Never commit `.env` - it is already in `.gitignore`
- Use separate OAuth apps for development and production
- Enable 2FA on all OAuth provider accounts
- Rotate secrets periodically using `./scripts/utilities/generate-secrets.sh`
- Only add trusted addresses to `ADMIN_ALLOWED_EMAILS`

## Production Deployment

### Environment Updates

```bash
ENVIRONMENT=production
NODE_ENV=production
GO_ENV=production

FRONTEND_URL=https://your-domain.com
REACT_APP_API_URL=https://your-domain.com
REACT_APP_WS_URL=wss://your-domain.com

# Update OAuth redirect URLs to https
GOOGLE_REDIRECT_URL=https://your-domain.com/api/v1/auth/admin/oauth/callback/google
GITHUB_REDIRECT_URL=https://your-domain.com/api/v1/auth/admin/oauth/callback/github
MICROSOFT_REDIRECT_URL=https://your-domain.com/api/v1/auth/admin/oauth/callback/microsoft
```

### Production Checklist

- Replace self-signed Redis TLS certificates in `deploy/redis/certs/` with CA-signed certificates
- Configure firewall: open ports 80 and 443, block direct access to 3000, 5432, 6379, 8080-8084, 9090, 3001
- Set up nginx for SSL/TLS termination and routing to backend services
- Create admin user: `docker exec -it portfolio_api go run /app/cmd/init-admin/main.go`
- Configure automated database backups:
  ```bash
  # Daily backup cron
  0 2 * * * /path/to/scripts/database/backup_db.sh
  ```
- Set up log rotation to prevent disk space issues
- Configure Prometheus alert rules for service downtime, high memory/CPU, database failures, and disk warnings
- Update all OAuth provider consoles with production redirect URLs
- Submit Google OAuth app for verification if it will have external users
- Set resource limits in `docker-compose.yml` appropriate for your server capacity
- Test disaster recovery using backup and restore scripts

### Admin User Creation

```bash
# Docker deployment
docker exec -it portfolio_api go run /app/cmd/init-admin/main.go

# Non-Docker deployment
go run backend/cmd/init-admin/main.go
```

## Monitoring and Maintenance

### Health Checks
Each microservice runs independently on its own port:
- API: http://localhost:8080/health
- DevPanel: http://localhost:8081/health
- Messaging: http://localhost:8082/health
- URL Shortener: http://localhost:8083/health
- Worker: http://localhost:8084/health (background tasks)

### Logs
```bash
# View all service logs
tail -f logs/*.log

# View specific service logs
tail -f logs/api.log
```

### Backup
```bash
# Backup database
./scripts/database/backup_db.sh

# Restore database
./scripts/database/restore_db.sh
```

## Troubleshooting

### Common Issues

1. Port Conflicts and Tmux Sessions
   **Note:** The `./start-dev.sh` and `./start-prod.sh` scripts now automatically detect and clean up existing tmux sessions, preventing most port conflicts.

   If you still encounter issues:
   ```bash
   # Check for processes using ports (5 microservices + frontend)
   sudo lsof -i :3000  # Frontend
   sudo lsof -i :8080  # API
   sudo lsof -i :8081  # DevPanel
   sudo lsof -i :8082  # Messaging
   sudo lsof -i :8083  # URL Shortener
   sudo lsof -i :8084  # Worker
   ```

2. Database Connection
   ```bash
   # Test database connection
   psql -h localhost -U postgres -d project_website
   ```

3. Frontend Build Issues
   ```bash
   # Clear npm cache
   npm cache clean --force
   
   # Remove node_modules and reinstall
   rm -rf node_modules
   npm install
   ```

4. Browser Cache Issues in Development
   ```bash
   # Use fresh start command
   cd frontend && npm run dev:fresh
   
   # Or manually clear caches
   cd frontend && npm run clear-cache
   
   # In browser console (dev mode only)
   window.__clearAllCaches()
   ```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

