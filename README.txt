# Personal Portfolio & Project Showcase

## Overview

This repository contains my personal website showcasing various software development projects along with their live implementations. The website serves both as a portfolio and as a platform hosting functional applications, starting with a robust URL shortener service.

```
Project-Website/
├── frontend/                 # React/TypeScript frontend
│   ├── build/                # Production build
│   ├── public/               # Static assets
│   │   ├── index.html        # HTML entry point
│   │   ├── favicon.ico       # Website favicon
│   │   └── manifest.json     # PWA manifest
│   ├── src/
│   │   ├── assets/           # Images and static resources
│   │   ├── components/       # Reusable UI components
│   │   │   ├── animations/   # Animation components
│   │   │   │   ├── CreativeShaderBackground.tsx  # Advanced shader-based background
│   │   │   │   ├── FloatingElement.tsx           # Floating animation component
│   │   │   │   ├── LoadingScreen.tsx             # Loading state UI
│   │   │   │   ├── PixelGridAnimation.tsx        # Grid-based animation effect
│   │   │   │   ├── ScrollIndicator.tsx           # Scroll direction indicator
│   │   │   │   └── SpaceAnimation.tsx            # Space-themed animation
│   │   │   ├── Footer/      # Footer component with site links
│   │   │   │   └── Footer.tsx
│   │   │   ├── layout/      # Layout components
│   │   │   │   └── Layout.tsx                   # Main app layout wrapper
│   │   │   ├── NavigationBar/ # Main navigation
│   │   │   │   └── NavigationBar.tsx            # Responsive navigation bar
│   │   │   ├── navigation/  # Navigation utilities
│   │   │   │   └── ScrollToTop.tsx              # Auto-scroll to top on navigation
│   │   │   ├── sections/    # Homepage sections
│   │   │   │   ├── About.tsx                    # About section (916 lines)
│   │   │   │   ├── Hero.tsx                     # Hero banner section (999 lines)
│   │   │   │   ├── ParallaxHero.tsx             # Parallax scrolling hero
│   │   │   │   ├── Projects.tsx                 # Projects showcase section
│   │   │   │   ├── Skills.tsx                   # Skills display section
│   │   │   │   └── SkillsSection.tsx            # Skills section wrapper
│   │   │   └── ui/          # Reusable UI elements
│   │   │       ├── LanguageFilter.tsx           # Programming language filter
│   │   │       ├── OptimizedImage.tsx           # Performance-optimized image component
│   │   │       ├── ProjectCard.tsx              # Project display card
│   │   │       ├── SkillBar.tsx                 # Skill level visualization
│   │   │       ├── Timeline.tsx                 # Timeline visualization
│   │   │       └── VirtualizedList.tsx          # Virtualized scrolling list
│   │   │
│   │   ├── contexts/        # React contexts
│   │   │   ├── ThemeContext.tsx                 # Theme management context
│   │   │   └── ZIndexContext.tsx                # Z-index layer management
│   │   │
│   │   ├── hooks/           # Custom React hooks
│   │   │   ├── useAnimationController.ts        # Animation control hook
│   │   │   ├── useAuth.tsx                      # Authentication state hook
│   │   │   ├── useClickOutside.ts               # Click outside element detection
│   │   │   ├── useDeviceCapabilities.ts         # Device capability detection
│   │   │   ├── usePerformanceOptimizations.ts   # Performance settings management
│   │   │   ├── useThemeToggle.ts                # Theme switching hook
│   │   │   ├── useTouchInteractions.ts          # Touch gesture handling
│   │   │   └── useZIndex.ts                     # Z-index management
│   │   │
│   │   ├── pages/           # Application pages
│   │   │   ├── About/       # About page
│   │   │   │   └── About.tsx                   # Detailed about page
│   │   │   ├── Contact/     # Contact page
│   │   │   │   └── Contact.tsx                 # Contact form page
│   │   │   ├── devpanel/    # Developer panel
│   │   │   │   └── DevPanel.tsx                # Admin tools interface
│   │   │   ├── Home/        # Homepage
│   │   │   │   └── Home.tsx                    # Main landing page
│   │   │   ├── messaging/   # Messaging app
│   │   │   │   └── Messaging.tsx               # Real-time messaging interface
│   │   │   ├── NotFound/    # 404 page
│   │   │   │   └── NotFound.tsx                # Error page for invalid routes
│   │   │   └── urlshortener/ # URL Shortener app
│   │   │       └── UrlShortener.tsx            # URL shortening interface
│   │   │
│   │   ├── styles/          # Styling
│   │   │   ├── GlobalStyles.ts                  # Global CSS styles
│   │   │   ├── theme.types.ts                   # TypeScript theme definitions
│   │   │   └── themes.ts                        # Theme configurations
│   │   │
│   │   ├── utils/           # Utility functions
│   │   │   ├── debugHelpers.ts                  # Development debugging utilities
│   │   │   ├── MemoryManager.tsx                # Memory optimization utilities
│   │   │   └── performance.ts                   # Performance measurement tools
│   │   │
│   │   ├── App.tsx          # Main App component with routing (109 lines)
│   │   └── index.tsx        # React entry point
│   │
│   └── package.json         # Dependencies and scripts
│
├── backend/                  # Go backend services
│   ├── cmd/                  # Service entry points
│   │   ├── api/              # Main API service
│   │   │   └── main.go
│   │   ├── devpanel/         # Developer panel service
│   │   │   └── main.go
│   │   ├── migration/        # Database migration tool
│   │   │   └── main.go
│   │   └── worker/           # Background worker
│   │       └── main.go
│   │
│   ├── config/               # Configuration
│   │   ├── app.yaml          # Main app config
│   │   ├── development.yaml  # Dev environment config
│   │   ├── production.yaml   # Production config
│   │   └── config.go         # Config loader
│   │
│   ├── deployments/          # Deployment configurations
│   │   ├── docker/           # Docker setup
│   │   │   └── docker-compose.yml
│   │   ├── nginx/            # Web server config
│   │   │   └── api.conf
│   │   └── systemd/          # Service definitions
│   │       └── api.service
│   │
│   ├── internal/             # Internal packages
│   │   ├── app/              # Application bootstrap
│   │   │   ├── bootstrap.go
│   │   │   └── server/       # HTTP server
│   │   │       ├── middleware/
│   │   │       │   ├── auth.go
│   │   │       │   └── cors.go
│   │   │       └── server.go
│   │   │
│   │   ├── common/           # Shared utilities
│   │   │   ├── auth/         # Authentication
│   │   │   │   ├── jwt.go
│   │   │   │   └── password.go
│   │   │   ├── cache/        # Caching
│   │   │   │   └── redis.go
│   │   │   ├── database/     # Database access
│   │   │   │   └── db.go
│   │   │   └── utils/        # Utility functions
│   │   │       └── url_validator.go
│   │   │
│   │   ├── domain/           # Business domain
│   │   │   ├── entity/       # Core entities
│   │   │   │   ├── user.go
│   │   │   │   └── audit.go
│   │   │   └── errors/       # Domain errors
│   │   │       └── errors.go
│   │   │
│   │   ├── messaging/        # Messaging service
│   │   │   ├── delivery/     # HTTP/WebSocket delivery
│   │   │   │   ├── http/     # HTTP handlers
│   │   │   │   │   └── handlers.go
│   │   │   │   └── websocket/# WebSocket handlers
│   │   │   │       ├── client.go
│   │   │   │       └── hub.go
│   │   │   ├── domain/       # Domain models
│   │   │   │   ├── channel.go
│   │   │   │   └── message.go
│   │   │   ├── repository/   # Data access
│   │   │   │   └── postgres/
│   │   │   │       ├── channel_repository.go
│   │   │   │       └── message_repository.go
│   │   │   └── service/      # Business logic
│   │   │       ├── messaging_service.go
│   │   │       └── service.go
│   │   │
│   │   └── urlshortener/     # URL shortener service
│   │       ├── delivery/     # HTTP delivery
│   │       │   └── http/
│   │       │       ├── handlers.go
│   │       │       └── routes.go
│   │       ├── domain/       # Domain models
│   │       │   ├── url.go
│   │       │   └── stats.go
│   │       ├── repository/   # Data access
│   │       │   └── postgres/
│   │       │       ├── url.go
│   │       │       └── stats.go
│   │       └── service/      # Business logic
│   │           ├── url.go
│   │           └── stats.go
│   │
│   ├── migrations/           # Database migrations
│   │   ├── common/           # Shared migrations
│   │   │   └── 000001_create_users_table.up.sql
│   │   ├── messaging/        # Messaging migrations
│   │   │   └── 000001_create_channels_table.up.sql
│   │   └── urlshortener/     # URL shortener migrations
│   │       └── 000001_create_urls_table.up.sql
│   │
│   ├── scripts/              # Utility scripts
│   │   ├── run.sh            # Run the application
│   │   └── setup.sh          # Setup environment
│   │
│   ├── go.mod                # Go dependencies
│   └── go.sum                # Go dependencies checksums
│
└── README.md                # Project documentation
```

## Architecture

The project is split into two main components:

### Backend (Go)
- Built with Go using Gin framework and GORM
- Implements a clean repository pattern architecture
- Provides RESTful APIs for the frontend
- Features robust database management and migrations
- Includes middlewares for authentication and rate limiting

### Frontend Features
- **Interactive UI**: Parallax effects, particle backgrounds, and smooth animations
- **Responsive Design**: Mobile-first approach ensuring compatibility across devices
- **Theme System**: Context-based theming with dark/light mode support
- **Custom Components**: Reusable UI elements like project cards, skill bars, and timelines

## Featured Projects

### URL Shortener
The first fully implemented project offering:
- Short URL creation with custom code options
- User account management and authentication
- Advanced analytics for URL performance tracking
- Custom domain support
- Private and expirable links

### Coming Soon
Additional planned projects to be integrated into the portfolio.

## Technical Highlights

### Backend Implementation
- **Database Layer**: Uses repository pattern with generic Go implementations
- **Authentication**: JWT-based auth with role-based permissions
- **Analytics**: Comprehensive click tracking and visualization
- **API Security**: Input validation, rate limiting, and proper error handling

### Frontend Features
- **Interactive UI**: Parallax effects, particle backgrounds, and smooth animations
- **Responsive Design**: Mobile-first approach ensuring compatibility across devices
- **Theme System**: Context-based theming with dark/light mode support
- **Custom Components**: Reusable UI elements like project cards, skill bars, and timelines

## Getting Started

### Setup and Installation

1. Clone the repository
   ```
   git clone https://github.com/JadenRazo/Project-Website.git
   cd Project-Website
   ```

2. Run the setup script to install dependencies and prepare the environment
   ```
   ./backend/scripts/setup.sh
   ```
   
   This script:
   - Installs backend dependencies (Go modules)
   - Installs frontend dependencies (npm)
   - Creates initial configuration files
   - Sets up the database with migrations

3. Run the application with the run script
   ```
   ./backend/scripts/run.sh
   ```
   
   This script:
   - Starts both the backend and frontend in tmux sessions
   - Provides a development environment with live reloading
   
   Available options:
   - First use chmod -x run.sh in /Project-Website/backend/scripts
   - Use bash or sh run.sh with any of the following below:
   - `-e, --env ENV` - Run in specific environment (development, staging, production)
   - `-w, --watch` - Run backend with file watching (hot reload)
   - `-d, --debug` - Run in debug mode
   - `-i, --install` - Install/update dependencies
   - `-s, --setup` - Setup configuration files
   - `-h, --help` - Print help message

4. Access the application
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

### Development Workflow

The application uses tmux to manage the backend and frontend processes. To view logs:

```
# View backend logs
tmux attach-session -t backend

# View frontend logs
tmux attach-session -t frontend
```

To detach from a tmux session without stopping it, press `Ctrl+b` then `d`.

To stop all services:
```
tmux kill-session -t backend
tmux kill-session -t frontend
```

Or simply run the script again to restart everything.

## API Endpoints

### URL Shortener
- `