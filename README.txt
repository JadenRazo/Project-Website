# Personal Portfolio & Project Showcase

## Overview

This repository contains my personal website showcasing various software development projects along with their live implementations. The website serves both as a portfolio and as a platform hosting functional applications, starting with a robust URL shortener service.

```
Project-Website/
├── cmd
│   ├── admin
│   │   └── main.go
│   ├── api
│   │   └── main.go
│   ├── migration
│   │   └── main.go
│   ├── server
│   │   └── main.go
│   └── worker
│       └── main.go
├── config
│   ├── app.yaml
│   ├── config.go
│   ├── development.yaml
│   ├── env.go
│   ├── production.yaml
│   └── testing.yaml
├── deployments
│   ├── docker
│   │   ├── api
│   │   │   └── Dockerfile
│   │   ├── docker-compose.dev.yml
│   │   ├── docker-compose.yml
│   │   └── worker
│   │       └── Dockerfile
│   ├── kubernetes
│   │   ├── base
│   │   └── overlays
│   │       ├── development
│   │       ├── production
│   │       └── staging
│   ├── nginx
│   │   ├── api.conf
│   │   └── ssl
│   └── systemd
│       ├── api.service
│       └── worker.service
├── docs
│   ├── api
│   │   ├── endpoints
│   │   │   ├── messaging.md
│   │   │   └── urlshortener.md
│   │   ├── openapi.json
│   │   └── swagger.yaml
│   └── architecture.md
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   ├── bootstrap.go
│   │   ├── config
│   │   │   └── config.go
│   │   ├── context
│   │   │   └── context.go
│   │   ├── health
│   │   └── server
│   │       ├── middleware
│   │       │   ├── auth.go
│   │       │   ├── cors.go
│   │       │   ├── logging.go
│   │       │   ├── recovery.go
│   │       │   └── spa.go
│   │       ├── router.go
│   │       └── server.go
│   ├── common
│   │   ├── auth
│   │   │   ├── jwt.go
│   │   │   └── password.go
│   │   ├── cache
│   │   │   └── redis.go
│   │   ├── database
│   │   │   ├── db.go
│   │   │   └── transaction.go
│   │   ├── errors
│   │   │   └── errors.go
│   │   ├── logger
│   │   │   └── logger.go
│   │   ├── metrics
│   │   │   ├── grafana.go
│   │   │   └── prometheus.go
│   │   ├── ratelimit
│   │   ├── resilience
│   │   ├── storage
│   │   │   └── storage.go
│   │   ├── tracing
│   │   │   ├── jaeger.go
│   │   │   └── opentelemetry.go
│   │   ├── utils
│   │   │   ├── token.go
│   │   │   └── url_validator.go
│   │   └── validator
│   │       └── validator.go
│   ├── core
│   ├── domain
│   │   ├── entity
│   │   │   ├── audit.go
│   │   │   └── user.go
│   │   ├── errors
│   │   │   └── errors.go
│   │   └── models.go
│   ├── messaging
│   │   ├── attachments
│   │   │   └── service.go
│   │   ├── delivery
│   │   │   ├── http
│   │   │   │   ├── handlers.go
│   │   │   │   ├── middleware.go
│   │   │   │   ├── read_receipt_handler.go
│   │   │   │   └── routes.go
│   │   │   └── websocket
│   │   │       ├── client.go
│   │   │       ├── connection_manager.go
│   │   │       ├── hub.go
│   │   │       ├── presence.go
│   │   │       └── types.go
│   │   ├── domain
│   │   │   ├── channel.go
│   │   │   ├── message.go
│   │   │   ├── reaction.go
│   │   │   ├── read_receipts.go
│   │   │   └── repository.go
│   │   ├── errors
│   │   │   └── errors.go
│   │   ├── events
│   │   │   ├── dispatcher.go
│   │   │   ├── event_types.go
│   │   │   └── handlers.go
│   │   ├── middleware
│   │   │   └── auth.go
│   │   ├── repository
│   │   │   ├── attachment_repository.go
│   │   │   ├── cache
│   │   │   │   ├── channel.go
│   │   │   │   └── message.go
│   │   │   ├── errors.go
│   │   │   ├── gorm_repository.go
│   │   │   ├── mock
│   │   │   │   └── repository.go
│   │   │   ├── postgres
│   │   │   │   ├── attachment_repository.go
│   │   │   │   ├── channel_repository.go
│   │   │   │   ├── message_repository.go
│   │   │   │   ├── reaction_repository.go
│   │   │   │   └── read_receipt_repository.go
│   │   │   └── read_receipt_repository.go
│   │   ├── service
│   │   │   ├── attachment_service.go
│   │   │   ├── channel_service.go
│   │   │   ├── messaging_service.go
│   │   │   ├── reaction_service.go
│   │   │   ├── read_receipt_service.go
│   │   │   └── service.go
│   │   ├── storage
│   │   └── usecase
│   │       ├── mark_as_read.go
│   │       ├── search_messages.go
│   │       ├── send_message.go
│   │       └── upload_attachment.go
│   ├── urlshortener
│   │   ├── delivery
│   │   │   └── http
│   │   │       ├── handlers.go
│   │   │       ├── middleware.go
│   │   │       └── routes.go
│   │   ├── domain
│   │   │   ├── repository.go
│   │   │   ├── stats.go
│   │   │   └── url.go
│   │   ├── repository
│   │   │   ├── cache
│   │   │   │   └── url.go
│   │   │   ├── gorm_repository.go
│   │   │   ├── mock
│   │   │   │   └── repository.go
│   │   │   ├── postgres
│   │   │   │   ├── stats.go
│   │   │   │   └── url.go
│   │   │   └── repository.go
│   │   ├── service
│   │   │   ├── service.go
│   │   │   ├── service_imp.go
│   │   │   ├── stats.go
│   │   │   └── url.go
│   │   ├── usecase
│   │   │   ├── resolve_url.go
│   │   │   ├── shorten_url.go
│   │   │   └── track_click.go
│   │   └── validator
│   └── worker
│       ├── queue
│       │   ├── kafka.go
│       │   └── rabbitmq.go
│       └── tasks
│           ├── messaging_tasks.go
│           └── scheduled.go
├── migrations
│   ├── common
│   │   ├── 000001_create_users_table.down.sql
│   │   └── 000001_create_users_table.up.sql
│   ├── messaging
│   │   ├── 000001_create_channels_table.down.sql
│   │   └── 000001_create_channels_table.up.sql
│   └── urlshortener
│       ├── 000001_create_urls_table.down.sql
│       └── 000001_create_urls_table.up.sql
├── pkg
│   ├── httputil
│   ├── pagination
│   └── validator
├── scripts
│   ├── backup_db.sh
│   ├── lint.sh
│   ├── restore_db.sh
│   ├── seed.sh
│   └── setup.sh
└─── web
│    ├── static
│    │   ├── css
│    │   ├── images
│    │   └── js
│    └── templates
│        ├── layouts
│        ├── pages
│        └── partials
├─ frontend/
│  ├─ build/
│  │  ├─ static
│  │  │  ├─ css
│  │  │  │  ├── main.e6c13ad2.css
│  │  │  │  └── main.e6c13ad2.css.map
│  │  │  ├─ js
│  │  │  │  ├── main.e896c9ee.js
│  │  │  │  ├── main.e896c9ee.js.LICENSE.txt
│  │  │  │  └── main.e896c9ee.js.map
│  │  ├── apple-touch-icon.png
│  │  ├── asset-manifest.json
│  │  ├── favicon-16x16.png
│  │  ├── favicon-32x32.png
│  │  ├── favicon.ico
│  │  ├── index.html
│  │  ├── manifest.json
│  │  └── robots.txt
│  ├─ node_modules/
│  ├─ public/
│  │  ├── apple-touch-icon.png
│  │  ├── favicon-16x16.png
│  │  ├── favicon-32x32.png
│  │  ├── favicon.ico
│  │  ├── index.html
│  │  ├── manifest.json
│  │  └── robots.txt
│  ├─ src/
│  │  ├─ assests/
│  │  │  ├─ images/
│  │  ├─ components/
│  │  │  ├─ animations/
│  │  │  │  ├─ FloatingElement.tsx
│  │  │  │  ├─ LoadingScreen.tsx
│  │  │  │  ├─ NetworkBackground.tsx
│  │  │  │  ├─ ParticleBackground.tsx
│  │  │  │  ├─ ScrollIndacator.tsx
│  │  │  ├─ layout/
│  │  │  │  ├─ Layout.tsx
│  │  │  │  ├─ NavigationBar.tsx
│  │  │  ├─ navigation/
│  │  │  │  ├─ BurgerMenu.tsx
│  │  │  ├─ sections/
│  │  │  │  ├─ Hero.tsx
│  │  │  │  ├─ ParallaxHero.tsx
│  │  │  │  ├─ Projects.tsx
│  │  │  │  ├─ Timeline.tsx
│  │  │  ├─ ui/
│  │  │  │  ├─ ProjectCard.tsx
│  │  │  │  ├─ SkillBar.tsx
│  │  │  │  ├─ Timeline.tsx
│  │  ├─ constants/
│  │  ├─ contexts/
│  │  │  ├─ Themecontext.tsx
│  │  ├─ docs/
│  │  │  ├─ ScrollTransform.md
│  │  ├─ types/
│  │  ├─ utils/
│  │  │  ├─ debugHelpers.ts 
│  │  │  ├─ MemoryManager.ts
│  │  │  └── performence.ts
│  │  ├─ hooks/
│  │  │  ├─ useClickOutside.ts
│  │  │  ├─ useAnimationController.ts
│  │  │  ├─ useDeviceCapabilities.ts
│  │  │  ├─ usePerformenceOptimizations.ts
│  │  │  ├─ useTouchInteractions.ts
│  │  │  └── useZIndex.ts
│  │  ├─ styles/
│  │  │  ├─ GlobalStyles.ts
│  │  │  ├─ theme.types.ts
│  │  │  └── themes.ts
│  │  ├─ app.css
│  │  ├─ app.tsx
│  │  ├─ custom.d.ts
│  │  ├─ index.css
│  │  ├─ index.html
│  │  ├─ index.tsx
│  │  ├─ logo.svg
│  │  ├─ styled.d.ts
│  ├─ package.json
├─ index.html
├─ README.md
```

## Architecture

The project is split into two main components:

### Backend (Go)
- Built with Go using Gin framework and GORM
- Implements a clean repository pattern architecture
- Provides RESTful APIs for the frontend
- Features robust database management and migrations
- Includes middlewares for authentication and rate limiting

### Frontend (React/TypeScript)
- Developed with React and TypeScript
- Showcases interactive UI elements and animations
- Implements responsive design with dark/light theme support
- Features component-based architecture for reusability
- Uses custom hooks for enhanced functionality

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

### Backend Setup
1. Navigate to backend: `cd Project-Website/backend`
2. Install dependencies: `go mod download`
3. Configure environment: Copy `.env.example` to `.env` and adjust values
4. Start server: `go run main.go`

### Frontend Setup
1. Navigate to frontend: `cd Project-Website/frontend`
2. Install dependencies: `npm install`
3. Start development server: `npm start`

## API Endpoints

### URL Shortener
- `POST /api/urls/shorten`: Create shortened URL
- `GET /api/urls`: List user's URLs
- `GET /:shortCode`: Redirect to original URL
- `GET /api/urls/:shortCode/analytics`: Get URL analytics

### Authentication
- `POST /api/auth/register`: Create account
- `POST /api/auth/login`: Authenticate user

## Technologies Used

### Backend
- Go (Golang)
- Gin Web Framework
- GORM (ORM)
- SQLite (Database)
- JWT (Authentication)

### Frontend
- React
- TypeScript
- Styled Components
- React Router
- Context API

## Development Approach

This project follows modern development practices:
- Type-safe programming with TypeScript and Go
- Component-based frontend architecture
- Repository pattern for database operations
- Clean separation between UI and business logic
- Comprehensive error handling

## Running in Production
For production deployment, additional considerations include:
- Database configuration for production environment
- Frontend build optimization
- Server configuration and deployment

---
