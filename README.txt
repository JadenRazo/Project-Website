# Personal Portfolio & Project Showcase

## Overview

This repository contains my personal website showcasing various software development projects along with their live implementations. The website serves both as a portfolio and as a platform hosting functional applications, starting with a robust URL shortener service.

```
Project-Website/
├─ backend/
│  ├─ config/
│  │  ├─ config.go
│  ├─ data/
│  ├─ db/
│  │  ├─ analytics.go
│  │  ├─ database.go
│  │  ├─ migrations.go
│  │  ├─ models.go
│  │  ├─ repository.go
│  │  ├─ url_model.go
│  ├─ handlers/
│  │  ├─ analytics_handlers.go
│  │  ├─ url_handlers.go
│  │  ├─ user_handlers.go
│  ├─ middleware/
│  │  ├─ auth.go
│  │  ├─ rate_limiter.go
│  ├─ static/
│  │  ├─ css/
│  │  │  ├─ style.css
│  ├─ templates/
│  │  ├─ 404.html
│  │  ├─ expired.html
│  │  ├─ index.html
│  ├─ utils/
│  │  ├─ token.go
│  │  ├─ url_validator.go
│  ├─ go.mod
│  ├─ main.go
├─ frontend/
│  ├─ build/
│  ├─ node_modules/
│  ├─ public/
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
│  │  ├─ contexts/
│  │  │  ├─ Themecontext.tsx
│  │  ├─ hooks/
│  │  │  ├─ useClickOutside.ts
│  │  ├─ styles/
│  │  │  ├─ GlobalStyles.ts
│  │  │  ├─ theme.types.ts
│  │  │  ├─ themes.ts
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
