# Personal Portfolio & Project Showcase

## Overview

This repository contains a full-stack portfolio website with integrated microservices. The system includes a modern React frontend and a Go-based backend with multiple microservices for URL shortening, messaging, and developer panel functionality.

## System Requirements

### Development Environment
- Go 1.21 or higher
- Node.js 18.x or higher
- npm 9.x or higher
- Git
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
├── frontend/                 # React/TypeScript frontend
│   ├── src/
│   │   ├── components/      # Reusable UI components
│   │   ├── contexts/        # React contexts (theme, auth)
│   │   ├── hooks/          # Custom React hooks
│   │   ├── pages/          # Page components
│   │   ├── styles/         # Global styles and themes
│   │   └── utils/          # Utility functions
│   └── public/             # Static assets
│
├── backend/                 # Go backend services
│   ├── cmd/                # Service entry points
│   │   ├── api/           # Main API service
│   │   ├── devpanel/      # Developer panel service
│   │   ├── messaging/     # Real-time messaging service
│   │   └── urlshortener/  # URL shortening service
│   │
│   ├── internal/          # Internal packages
│   │   ├── common/        # Shared utilities
│   │   ├── domain/        # Business logic
│   │   ├── messaging/     # Messaging service
│   │   └── urlshortener/  # URL shortener service
│   │
│   └── migrations/        # Database migrations
│
├── scripts/               # Utility scripts
└── deployments/          # Deployment configurations
```

## Local Development Setup

### Windows Setup

1. Install Prerequisites:
   ```powershell
   # Install Chocolatey (if not installed)
   Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

   # Install Go
   choco install golang

   # Install Node.js
   choco install nodejs

   # Install Git
   choco install git
   ```

2. Clone the repository:
   ```powershell
   git clone https://github.com/yourusername/Project-Website.git
   cd Project-Website
   ```

3. Install dependencies:
   ```powershell
   # Install frontend dependencies
   cd frontend
   npm install

   # Install backend dependencies
   cd ../backend
   go mod download
   ```

4. Start development servers:
   ```powershell
   # From project root
   npm run dev
   ```

### Ubuntu Setup

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

   # Install Git
   sudo apt install -y git
   ```

2. Clone and setup:
   ```bash
   git clone https://github.com/yourusername/Project-Website.git
   cd Project-Website

   # Install frontend dependencies
   cd frontend
   npm install

   # Install backend dependencies
   cd ../backend
   go mod download
   ```

3. Start development servers:
   ```bash
   # From project root
   npm run dev
   ```

## Service Ports

- Frontend Development: http://localhost:3000
- Main API: http://localhost:8080
- DevPanel: http://localhost:8081
- Messaging Service: http://localhost:8082
- URL Shortener: http://localhost:8083

## Available Scripts

### Development
```bash
# Start all services in development mode
npm run dev

# Start only frontend
npm run dev:frontend

# Start only backend
npm run dev:backend
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

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 