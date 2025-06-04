# Personal Portfolio & Project Showcase

## Overview

This repository contains a full-stack portfolio website with integrated microservices. The system includes a modern React frontend and a Go-based backend with multiple microservices for URL shortening, messaging, and developer panel functionality.

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
│   │   │   │   ├── CreativeShaderBackground.tsx
│   │   │   │   ├── FloatingElement.tsx  # Floating animations
│   │   │   │   ├── LoadingScreen.tsx  # Loading UI
│   │   │   │   ├── PixelGridAnimation.tsx  # Grid animations
│   │   │   │   ├── ScrollIndicator.tsx  # Scroll indicator
│   │   │   │   └── SpaceAnimation.tsx  # Space animations
│   │   │   ├── Footer/     # Footer components
│   │   │   │   └── Footer.tsx
│   │   │   ├── layout/     # Layout components
│   │   │   │   └── Layout.tsx  # Main layout wrapper
│   │   │   ├── NavigationBar/  # Navigation
│   │   │   │   └── NavigationBar.tsx  # Nav bar
│   │   │   ├── navigation/  # Navigation utilities
│   │   │   │   └── ScrollToTop.tsx  # Auto-scroll
│   │   │   ├── sections/    # Page sections
│   │   │   │   ├── About.tsx  # About section
│   │   │   │   ├── Hero.tsx  # Hero section
│   │   │   │   ├── ParallaxHero.tsx  # Parallax effects
│   │   │   │   ├── Projects.tsx  # Projects section
│   │   │   │   ├── Skills.tsx  # Skills section
│   │   │   │   └── SkillsSection.tsx  # Skills wrapper
│   │   │   └── ui/          # UI elements
│   │   │       ├── LanguageFilter.tsx  # Language filter
│   │   │       ├── OptimizedImage.tsx  # Optimized images
│   │   │       ├── ProjectCard.tsx  # Project cards
│   │   │       ├── SkillBar.tsx  # Skill visualization
│   │   │       ├── Timeline.tsx  # Timeline component
│   │   │       └── VirtualizedList.tsx  # Virtual lists
│   │   ├── contexts/        # React contexts
│   │   │   ├── ThemeContext.tsx  # Theme management
│   │   │   └── ZIndexContext.tsx  # Z-index management
│   │   ├── hooks/           # Custom hooks
│   │   │   ├── useAnimationController.ts  # Animation control
│   │   │   ├── useAuth.tsx  # Authentication
│   │   │   ├── useClickOutside.ts  # Click detection
│   │   │   ├── useDeviceCapabilities.ts  # Device features
│   │   │   ├── usePerformanceOptimizations.ts  # Performance
│   │   │   ├── useThemeToggle.ts  # Theme switching
│   │   │   ├── useTouchInteractions.ts  # Touch gestures
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
│   │   │   └── urlshortener/ # URL Shortener
│   │   │       └── UrlShortener.tsx
│   │   ├── styles/          # Styling
│   │   │   ├── GlobalStyles.ts  # Global styles
│   │   │   ├── theme.types.ts  # Theme types
│   │   │   └── themes.ts  # Theme configs
│   │   ├── utils/           # Utilities
│   │   │   ├── debugHelpers.ts  # Debugging tools
│   │   │   ├── MemoryManager.tsx  # Memory management
│   │   │   └── performance.ts  # Performance tools
│   │   ├── App.tsx          # Main App component
│   │   └── index.tsx        # Entry point
│   └── package.json         # Dependencies
│
├── backend/                  # Go backend
│   ├── cmd/                  # Entry points
│   │   ├── api/              # API service
│   │   │   └── main.go
│   │   ├── devpanel/         # Developer panel
│   │   │   └── main.go
│   │   ├── migration/        # DB migrations
│   │   │   └── main.go
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
│   │   │   ├── bootstrap.go
│   │   │   └── server/       # HTTP server
│   │   │       ├── middleware/
│   │   │       │   ├── auth.go
│   │   │       │   └── cors.go
│   │   │       └── server.go
│   │   ├── common/           # Shared utils
│   │   │   ├── auth/         # Authentication
│   │   │   │   ├── jwt.go
│   │   │   │   └── password.go
│   │   │   ├── cache/        # Caching
│   │   │   │   └── redis.go
│   │   │   ├── database/     # DB access
│   │   │   │   └── db.go
│   │   │   └── utils/        # Utilities
│   │   │       └── url_validator.go
│   │   ├── domain/           # Business domain
│   │   │   ├── entity/       # Core entities
│   │   │   │   ├── user.go
│   │   │   │   └── audit.go
│   │   │   └── errors/       # Domain errors
│   │   │       └── errors.go
│   │   ├── messaging/        # Messaging
│   │   │   ├── delivery/     # HTTP/WS delivery
│   │   │   │   ├── http/     # HTTP handlers
│   │   │   │   │   └── handlers.go
│   │   │   │   └── websocket/ # WS handlers
│   │   │   │       ├── client.go
│   │   │   │       └── hub.go
│   │   │   ├── domain/       # Models
│   │   │   │   ├── channel.go
│   │   │   │   └── message.go
│   │   │   ├── repository/   # Data access
│   │   │   │   └── postgres/
│   │   │   │       ├── channel_repository.go
│   │   │   │       └── message_repository.go
│   │   │   └── service/      # Business logic
│   │   │       ├── messaging_service.go
│   │   │       └── service.go
│   │   └── urlshortener/     # URL shortener
│   │       ├── delivery/     # HTTP delivery
│   │       │   └── http/
│   │       │       ├── handlers.go
│   │       │       └── routes.go
│   │       ├── domain/       # Models
│   │       │   ├── url.go
│   │       │   └── stats.go
│   │       ├── repository/   # Data access
│   │       │   └── postgres/
│   │       │       ├── url.go
│   │       │       └── stats.go
│   │       └── service/      # Business logic
│   │           ├── url.go
│   │           └── stats.go
│   ├── migrations/           # DB migrations
│   │   ├── common/           # Shared migrations
│   │   │   └── 000001_create_users_table.up.sql
│   │   ├── messaging/        # Messaging migrations
│   │   │   └── 000001_create_channels_table.up.sql
│   │   └── urlshortener/     # URL shortener migrations
│   │       └── 000001_create_urls_table.up.sql
│   ├── scripts/              # Scripts
│   │   ├── run.sh            # Run app
│   │   └── setup.sh          # Setup env
│   ├── go.mod                # Go dependencies
│   └── go.sum                # Go checksums
└── README.md                # Documentation
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
- Check all required dependencies
- Set up PostgreSQL database connections
- Run database migrations
- Start all backend microservices
- Start the frontend development server
- Create organized tmux sessions for easy monitoring

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
- **Code Statistics**: Automatic code line counting with tokei (updates hourly)
- **URL Shortener**: Create and track short URLs
- **Messaging System**: Real-time messaging with WebSocket support
- **Developer Panel**: Admin interface for system management

## Available Scripts

### Development
```bash
# RECOMMENDED: Start all services with fresh cache clearing
./start-dev.sh --fresh

# Start all services (normal mode)
./start-dev.sh

# Start with additional options
./start-dev.sh --kill-existing  # Kill processes on required ports
./start-dev.sh --skip-deps       # Skip dependency installation
./start-dev.sh --verbose         # Enable verbose output
./start-dev.sh --help           # Show help message

# Start with fresh cache clearing (recommended for cache issues)
./start-dev.sh --fresh

# Stop all services
./stop-dev.sh

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
./scripts/start-services.sh
```

## Database Setup

### Development
```bash
# Run migrations
cd backend
go run cmd/migration/main.go up

# Seed development data
./scripts/seed.sh
```

### Production
```bash
# Create production database
createdb project_website

# Run migrations
go run cmd/migration/main.go up -env production
```

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
   ```

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
- API: http://localhost:8080/health
- DevPanel: http://localhost:8081/health
- Messaging: http://localhost:8082/health
- URL Shortener: http://localhost:8083/health

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
./scripts/backup_db.sh

# Restore database
./scripts/restore_db.sh
```

## Troubleshooting

### Common Issues

1. Port Conflicts
   ```bash
   # Check for processes using ports
   sudo lsof -i :8080
   sudo lsof -i :8081
   sudo lsof -i :8082
   sudo lsof -i :8083
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

