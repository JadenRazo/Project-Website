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

### Docker Deployment
```bash
# Build and run with Docker Compose
docker-compose up --build
```

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

   Note: Each service is an independent binary that runs as a separate process.

3. Configure Nginx:
   ```nginx
   # Example Nginx configuration
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

       # Add similar blocks for other services
   }
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

