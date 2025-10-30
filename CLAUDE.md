# Portfolio Website - Claude Development Guide

## Project Overview
Full-stack portfolio website: React/TypeScript frontend + Go backend with 5 independent microservices running as separate processes (URL shortener, messaging, devpanel, visitor analytics, worker).

## Quick Start
**Development:** `./start-dev.sh --fresh` (from project root) - Starts all services in tmux sessions with automatic cleanup of existing sessions
**Production:** `./start-prod.sh` - Start production services with automatic session cleanup | `docker-compose up -d` - Deploy with Docker
**Database Setup:** Apply `backend/schema.sql` then run `./scripts/database/seed.sh`

## Essential Commands

### Development
- **Start:** `./start-dev.sh --fresh` (automatic session cleanup included)
- **Frontend Build:** `npm run build` | `npm run build:prod` (optimized)
- **Backend Build:** `go build -o bin/[service] backend/cmd/[service]/main.go`
- **Quality Checks:** `npm run lint:fix && npm run type-check` (frontend) | `./scripts/development/lint.sh` (backend)
- **Cache Clear:** `npm run dev:fresh` or `window.__clearAllCaches()` in console

### Testing
- **Frontend:** `npm run test` - React Testing Library tests
- **Backend:** `go test ./...` - Unit tests with coverage
- **Performance:** `npm run build:analyze` - Bundle size analysis

### Database & Admin
- **Initialize DB:** `psql -U user -d db < backend/schema.sql` (required before first run)
  - Creates 48 tables including 7 visitor analytics tables
  - Enables extensions: uuid-ossp, pgcrypto, citext
  - Sets up indexes for performance optimization
- **Seed Data:** `./scripts/database/seed.sh` - Populate with sample projects/data
- **Create Admin:** `go run backend/cmd/init-admin/main.go` - Interactive admin creation
- **Verify Admin:** `go run backend/cmd/verify-admin/main.go` - Check admin exists and credentials work
- **Admin Panel:** `go run backend/cmd/admin-create/main.go` - Alternative admin creation tool
- **Verify Visitor Tables:** API service runs `verifyVisitorTables()` on startup, logs error if any of 7 tables missing

## Service Ports
- Frontend: :3000 | Main API: :8080 | DevPanel: :8081 | Messaging: :8082
- URL Shortener: :8083 | Worker: :8084 | Prometheus: :9090 | Grafana: :3001

## Recent Major Features (All Production-Ready)
All features listed below are fully implemented, tested, and operational in production:

- **Service Management System:** 5 microservices registered with ServiceManager (urlshortener, messaging, devpanel, visitor, worker), full lifecycle control via DevPanel UI (Start/Stop/Restart), real-time status monitoring, health checks
- **Service Analytics:** Real-time metrics collection with in-memory storage (7-day retention), 30-second collection interval, CPU/memory tracking via gopsutil, atomic request/error counters in BaseService, historical data API with time-range queries (1h/6h/24h/7d)
- **Visitor Analytics:** Complete tracking system with middleware integration (line 429), 7 database tables (sessions/metrics/locations/realtime/summaries), automated aggregation via 5 worker cron jobs (hourly/daily/weekly), SHA256 session hashing, IP geolocation, real-time active counts, GDPR/CCPA/LGPD compliance with consent tracking
- **Worker Service:** Full implementation of core.Service interface, wraps ScheduledTasks with robfig/cron v3, 5 automated cron jobs for visitor data aggregation (hourly metrics, daily summaries, session cleanup, location aggregates, old data purge), managed lifecycle via DevPanel API
- **Secure Redis:** TLS 1.3 encryption, ACL authentication, secure connection handling, setup automation via `./scripts/utilities/setup-redis-security.sh`
- **DevPanel Expansion:** CRUD for certifications, skills, prompts, project paths, visitor analytics dashboard with real-time stats/graphs, service management UI with time-series visualization
- **Circuit Breaker:** Fault tolerance pattern implemented for external service dependencies, automatic retry logic
- **Enhanced Security:** Rate limiting (per-IP and per-user), JWT authentication with refresh tokens, bcrypt password hashing, secure middleware stack, CORS with origin validation

## State Management & Performance
- **State:** Zustand stores (authStore, performanceStore, themeStore)
- **Optimization:** Lazy loading, code splitting, memoization
- **Caching:** Redis with secure TLS, in-memory caching
- **WebSocket:** Connection manager with presence system
- **Monitoring:** Real-time metrics via Prometheus/Grafana

## Key Architecture

### Frontend Structure
- **Main:** `src/App.tsx` | `src/index.tsx`
- **Pages:** `src/pages/` (Home, Projects, About, Status, DevPanel, Messaging, UrlShortener, Contact, NotFound)
- **Components:** `src/components/`
  - `animations/` - Animation components
  - `auth/` - Authentication UI components
  - `common/` - Shared reusable components
  - `devpanel/` - AdminLogin, ProjectManager, CertificationsManager, SkillsManager, PromptsManager, ProjectPathsManager, VisitorAnalytics, DevPanelLoadingState, DevPanelSectionNav
  - `Footer/` - Footer component
  - `layout/` - Layout wrappers
  - `metrics/` - SystemMetrics component
  - `NavigationBar/` - Main navigation
  - `navigation/` - ScrollToTop, PageTop, PageTransition
  - `notifications/` - Notification system
  - `sections/` - Hero, Projects sections
  - `skeletons/` - Loading skeleton components
  - `ThemeToggle/` - Theme switching UI
  - `ui/` - UI primitives
- **Hooks:** `src/hooks/`
  - useAuth, useAnimationController, useClickOutside
  - useCrudOperations, useDebounce, useDeviceCapabilities
  - useInlineFormScroll, useLazyLoad, useLoadingState
  - useMemoryManager, useMobileOptimizations, useModalScroll
  - useOptimizedScrollHandler, usePerformanceOptimizations, usePreloader
  - useScrollTo, useScrollToForm, useTheme, useThemeToggle
  - useTouchInteractions, useWebSocketConnection, useZIndex
- **State:** `src/stores/` (authStore, performanceStore, themeStore)
- **Utils:** `src/utils/`
  - apiConfig, debugHelpers, devCacheManager, devPanelApi
  - errorHandler, lazyWithPreload, MemoryManager, performance
  - performanceConfig, performanceMonitor, preloader
  - promptApi, scrollConfig, scrollTestUtils, validation
- **Assets:** `src/assets/` (data, icons, images, videos)
- **Config:** `src/config/env-config.js` | `public/env-config.js`

### Backend Structure
- **Entry Points:** `backend/cmd/*/main.go`
  - `api/` - Main API service with gateway, auth, and service registry
  - `devpanel/` - DevPanel standalone service
  - `messaging/` - Messaging standalone service
  - `urlshortener/` - URL shortener standalone service
  - `worker/` - Background worker standalone service
  - `admin/`, `admin-create/`, `init-admin/`, `verify-admin/` - Admin CLI tools
  - `check-db/` - Database health check utility
  - `simple-api/`, `server/` - Alternative server implementations for testing
- **Services:** `backend/internal/*/` (domain-driven design)
  - `codestats/` - Code statistics with projectpath management
  - `devpanel/` - certification/, skill/, prompt/ managers + service
  - `messaging/` - Full messaging system with WebSocket, events, use cases
  - `projects/` - Project management CRUD
  - `status/` - Service health monitoring
  - `urlshortener/` - URL shortening with entities
  - `visitor/` - Visitor analytics and tracking
  - `worker/` - Background task processing
- **Common:** `backend/internal/common/`
  - auth/ (JWT, admin, password), cache/ (Redis with TLS/ACL)
  - circuitbreaker/, compression/, metrics/ (Prometheus integration, Manager)
  - middleware/, ratelimit/, repository/, response/
  - security/, testutil/, utils/, validator/
- **Core:** `backend/internal/core/` (config, db, service_manager, base_service with atomic counters)
  - ServiceManager: Registers 5 microservices, provides lifecycle control (Start/Stop/Restart), health monitoring
  - Service interface: Start(), Stop(), Restart(), Status(), Name(), HealthCheck()
  - BaseService: Embedded in all microservices, tracks uptime/requests/errors with atomic counters
- **Domain:** `backend/internal/domain/` (entity, errors)
- **DevPanel Metrics:** `backend/internal/devpanel/metrics.go`
  - MetricsCollector: In-memory time-series storage, configurable collection interval (default 30 seconds)
  - Started in `backend/cmd/api/main.go` line 348 via `metricsCollector.StartCollecting(serviceManager)`
  - Collects CPU/memory via gopsutil (process-level metrics), request/error counts from BaseService atomic counters
  - 7-day retention (168 hours) with automatic trimming of old data points
  - Runs in background goroutine with ticker, no persistence across restarts
- **Service Architecture:** 5 independent microservices with separate binaries and processes
  - Each service has its own `backend/cmd/[service]/main.go` entry point
  - Independent processes: api (8080), devpanel (8081), messaging (8082), urlshortener (8083), worker (8084)
  - ServiceManager in API service provides centralized monitoring and control via DevPanel UI
  - Start/Stop controls in DevPanel affect service registration state, not actual OS processes
- **Config:** `backend/config/*.yaml` (app, development, production, devpanel) | Environment variables

## Scripts Organization
- **Project Root:**
  - `./start-dev.sh` - Start all services in tmux (development mode, auto-detects and cleans existing sessions)
    - Frontend session: `portfolio-frontend` (1 window for React dev server)
    - Backend session: `portfolio-backend` (5 windows: api, devpanel, messaging, urlshortener, worker)
  - `./start-prod.sh` - Start production services in tmux (auto-detects and cleans existing sessions)
    - Frontend session: `portfolio-frontend-prod` (1 window for static file server)
    - Backend session: `portfolio-backend-prod` (5 windows: api, devpanel, messaging, urlshortener, worker)
  - `./docker-compose.yml` - Production deployment configuration
- **Subdirectories:**
  - `scripts/database/` - Migration management, seeding, backups (backup_db.sh, restore_db.sh, seed.sh, update-db.sh)
  - `scripts/development/` - Linting, building, local development tools (lint.sh, start-dev.sh, stop-dev.sh, start-prod.sh, stop-prod.sh)
  - `scripts/monitoring/` - Metrics collection, health checks
  - `scripts/utilities/` - One-time setup scripts (setup-admin.sh, setup-redis-security.sh, backend-setup.sh, backend-lint.sh)

## Docker Deployment
```bash
docker-compose up -d              # Start all services
docker-compose logs -f [service]  # View logs
docker-compose down               # Stop services
docker-compose build --no-cache   # Rebuild images
```
**Services:** frontend, backend, postgres, redis, prometheus, grafana
**Networks:** Internal network with service discovery
**Volumes:** Persistent data for DB, Redis, monitoring

## Code Style Guidelines
### Frontend
- ES modules, functional components with hooks
- TypeScript strict mode, interfaces over types
- PascalCase components, camelCase functions
- Tailwind utility classes, styled-components for complex styles
- NO comments unless explicitly requested

### Backend
- Standard Go idioms, explicit error handling
- Dependency injection, interface-based design
- Domain-driven structure, repository pattern
- Context propagation, graceful shutdowns
- NO comments unless explicitly requested
- **Middleware Order Critical:** Register tracking middleware BEFORE route registration to ensure execution on all requests
  - Example: `router.Use(visitor.TrackingMiddleware(visitorService))` must come before `apiGateway.RegisterService()` calls
  - Current implementation: Line 429 (middleware) comes before line 433 (first route registration)

### Implementing Service Metrics
Services should track request and error counts for analytics:
```go
// In service struct, embed BaseService
type MyService struct {
    *core.BaseService
    // other fields
}

// In handler functions
func (s *MyService) MyHandler(c *gin.Context) {
    s.IncrementRequests()  // Track incoming request

    result, err := s.doSomething()
    if err != nil {
        s.IncrementErrors()  // Track errors
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, result)
}
```
Metrics are automatically collected every 30 seconds by MetricsCollector.

## Development Workflow
1. `git pull origin main` - Get latest changes
2. `./start-dev.sh --fresh` - Start development environment (automatically handles existing sessions)
3. Create feature branch: `git checkout -b feature/description`
4. Tmux navigation:
   - Frontend: `tmux attach -t portfolio-frontend`
   - Backend: `tmux attach -t portfolio-backend` (Ctrl+B + 0-4 to switch between api/devpanel/messaging/urlshortener/worker windows)
5. Make changes, run quality checks
6. Commit atomically with clear messages
7. Push and create PR to main branch

**Note:** The start scripts automatically detect and clean up existing tmux sessions, preventing port conflicts. Backend session has 5 windows for independent microservice processes.

## Testing
### Frontend Testing
- **Unit Tests:** `npm run test` - Component and hook tests
- **Coverage:** `npm run test -- --coverage`
- **E2E:** Cypress tests in `cypress/`

### Backend Testing
- **Unit Tests:** `go test ./...`
- **Coverage:** `go test -cover ./...`
- **Integration:** `go test -tags=integration ./...`
- **Benchmarks:** `go test -bench=. ./internal/performance/`

## Monitoring & Observability
- **Metrics:** Prometheus at :9090, custom metrics via `/metrics` endpoints on each service
- **Service Analytics (DevPanel):**
  - MetricsCollector: Runs on 30-second ticker, stores time-series data in-memory for 7 days
  - **Data Collection:** CPU/memory via gopsutil (process-level), request/error counts from BaseService atomic counters
  - **API Endpoints:**
    - Service list: `GET /api/v1/devpanel/services` - All registered services with current stats
    - Service detail: `GET /api/v1/devpanel/services/:name` - Specific service status
    - Historical metrics: `GET /api/v1/devpanel/metrics/:service/history?duration=1h|6h|24h|7d`
    - Lifecycle control: `POST /api/v1/devpanel/services/:name/start|stop|restart`
  - **UI Dashboard:** Real-time graphs showing CPU, memory, requests, errors per service with time-series visualization
- **Visitor Analytics (DevPanel):**
  - Real-time active visitors count with 5-minute activity window
  - Daily/weekly/monthly aggregate statistics
  - Geographic distribution from IP geolocation
  - Top pages by view count
  - Session tracking with SHA256 privacy hashing
- **Service Management:** All 5 microservices visible in DevPanel UI (urlshortener, messaging, devpanel, visitor, worker) with lifecycle controls
- **Dashboards:** Grafana at :3001 with pre-configured panels (when enabled)
- **Logging:** Structured logging with log levels, log streaming in DevPanel, centralized error tracking
- **Tracing:** OpenTelemetry integration (optional, when configured)
- **Alerts:** Configured in `deploy/prometheus/alerts.yml` (production deployment)

## Service-Specific Notes

### Messaging Service
- WebSocket real-time communication
- Presence system for online status
- Message persistence with PostgreSQL
- File attachments with S3 integration
- Rate limiting per connection

### URL Shortener
- Analytics tracking (clicks, referrers, geographic data)
- Custom slugs support
- QR code generation
- Bulk operations API
- Cache-first architecture

### DevPanel
- Admin authentication required
- Project management CRUD
- **Service Management UI:** Monitor all 5 independent microservices (urlshortener, messaging, devpanel, visitor, worker), real-time status/uptime/metrics, Start/Stop/Restart controls affect registration state only
- **Visitor Analytics Dashboard:** Real-time active visitors, daily/weekly/monthly stats, geographic distribution, top pages, session tracking
- Real-time logs streaming
- System metrics monitoring
- Service analytics with real-time graphs (CPU, memory, requests, errors)
- MetricsCollector started in main.go, collects every 30 seconds
- Services track metrics via BaseService.IncrementRequests() and IncrementErrors()
- Visitor data endpoints: `/api/v1/devpanel/visitors/stats`, `/api/v1/devpanel/visitors/realtime`, `/api/v1/devpanel/visitors/timeline`, `/api/v1/devpanel/visitors/locations`
- Service control endpoints: `/api/v1/devpanel/services`, `/api/v1/devpanel/services/:name`, `/api/v1/devpanel/services/:name/start|stop|restart`

### Worker Service
- Implements core.Service interface, registered with ServiceManager as "worker"
- Wraps ScheduledTasks (cron jobs) with lifecycle management
- Start() launches all cron jobs in goroutine, Stop() cancels context and halts all tasks
- Visible in DevPanel service management UI with real-time status/uptime/metrics
- Background job processing with robfig/cron v3 (supports seconds precision)
- **Visitor Metrics Aggregation:** 5 automated cron jobs for visitor analytics
  - **Hourly metrics aggregation:** `0 * * * *` (every hour at :00) - Aggregates visitor data into `visitor_metrics` table
  - **Daily summary generation:** `5 0 * * *` (daily at 12:05 AM) - Generates daily summaries in `visitor_daily_summary`
  - **Expired session cleanup:** `0 */6 * * *` (every 6 hours) - Removes expired sessions from `visitor_sessions`
  - **Location aggregates update:** `0 2 * * *` (daily at 2:00 AM) - Updates `visitor_locations` aggregate data
  - **Old metrics cleanup:** `0 3 * * 0` (weekly Sunday 3:00 AM) - Purges metrics older than 90 days
- Service control: Can be started/stopped/restarted via DevPanel API (`/api/v1/devpanel/services/worker/start|stop|restart`)
- Logs success/failure for each task execution with structured logging

## Security Best Practices
- **Secrets:** Never commit - use environment variables or `.env` files
- **Input Validation:** Sanitize all user inputs, use prepared statements
- **Authentication:** JWT with refresh tokens, bcrypt for passwords
- **Rate Limiting:** Per-IP and per-user limits on all endpoints
- **CORS:** Strict origin checking, credentials support
- **TLS:** Enforce HTTPS in production, secure WebSocket connections
- **Visitor Privacy:** No IP storage, SHA256 session hashing, consent tracking (GDPR/CCPA/LGPD), automatic session expiration (24 hours), bot detection

## SQL/Database Guidelines
- Consolidate migrations into `backend/schema.sql` after applying
- Use transactions for multi-table operations
- Index foreign keys and frequently queried columns
- Avoid N+1 queries, use joins or batch loading
- Regular backups via `scripts/database/backup_db.sh`

## Common Issues & Solutions

### Port Conflicts & Tmux Sessions
**Automatic Resolution:** The `./start-dev.sh` and `./start-prod.sh` scripts now automatically detect and clean up existing tmux sessions that could cause port conflicts.

**Manual Resolution (if needed):**
```bash
lsof -i :PORT        # Find process
kill -9 PID          # Kill process
```

### Cache Issues
```bash
npm run dev:fresh              # Frontend cache clear
redis-cli FLUSHALL            # Redis cache clear (dev only)
```

### WebSocket Connection
- Check `WS_URL` in frontend `.env`
- Verify CORS settings in backend
- Check nginx WebSocket upgrade headers

### Build Failures
- Frontend: Clear `node_modules` and reinstall
- Backend: Update dependencies with `go mod tidy`

### Database Connection
- Check `DATABASE_URL` format
- Verify PostgreSQL is running
- Check connection pool settings

### Service Analytics Issues
- **Graphs are blank:** Ensure MetricsCollector started at line 348 in `backend/cmd/api/main.go` (`metricsCollector.StartCollecting(serviceManager)`)
- **Old data missing after restart:** By design - metrics stored in-memory only (no persistence), data lost on service restart for performance reasons
- **Service not showing stats:** Verify service embeds `*core.BaseService` and calls `IncrementRequests()` in handlers and `IncrementErrors()` on errors
- **No historical data initially:** Wait 30 seconds for first collection cycle to complete, MetricsCollector runs on 30-second ticker
- **Only 5 services visible:** By design - only lifecycle-managed microservices registered with ServiceManager (urlshortener, messaging, devpanel, visitor, worker); utility services (projects, codestats, status, auth) excluded intentionally
- **CPU/Memory always same:** Metrics are process-level (gopsutil tracks entire API process), not per-service isolation
- **Request counts not increasing:** Check handlers call `s.IncrementRequests()` where `s` is the service instance with embedded BaseService
- **Worker service not aggregating:** Verify running state via `/api/v1/devpanel/services/worker`, check logs for cron task registration and execution, ensure not stopped via Stop API

### Visitor Analytics Issues
- **No visitor data appearing:**
  - Check all 7 database tables exist: `psql $DATABASE_URL < backend/schema.sql`
  - Verify middleware registered BEFORE routes at line 429: `router.Use(visitor.TrackingMiddleware(visitorService))`
  - Check API startup logs for `verifyVisitorTables()` validation results
- **Active visitors always 0:**
  - Query `visitor_realtime` table for recent entries: `SELECT COUNT(*) FROM visitor_realtime WHERE last_activity > NOW() - INTERVAL '5 minutes';`
  - Look for `[VISITOR TRACKING ERROR]` prefix in logs indicating middleware failures
  - Verify middleware is executing (add debug logging if needed)
- **Tracking errors in console:**
  - `[VISITOR TRACKING ERROR]` prefix indicates database issues or missing tables
  - Check database connection string and credentials
  - Run `verifyVisitorTables()` validation (automatic on API startup)
- **No aggregated metrics/summaries:**
  - Ensure worker service running: `curl http://localhost:8080/api/v1/devpanel/services/worker`
  - Check logs for success messages: "Hourly visitor metrics aggregated successfully", "Daily visitor summary generated successfully"
  - Verify cron jobs registered: Look for "Starting scheduled tasks" in worker logs
- **Tables missing warning on startup:**
  - Apply complete schema including all 7 visitor tables: `psql $DATABASE_URL < backend/schema.sql`
  - Required tables: `visitor_sessions`, `page_views`, `visitor_realtime`, `visitor_metrics`, `visitor_daily_summary`, `privacy_consents`, `visitor_locations`
  - Startup validation logs missing tables if any are absent

## Architecture Decisions & Tradeoffs

### Service Analytics Design
- **In-Memory Storage:** Metrics stored in RAM (not database) for performance
  - **Benefit:** Zero database overhead, instant queries, no write amplification
  - **Tradeoff:** Data lost on service restart, no long-term historical analysis
  - **Mitigation:** 7-day retention sufficient for operational monitoring, Prometheus available for long-term storage
- **Process-Level Metrics:** CPU/memory tracked per microservice process
  - **Implementation:** Each service (api, devpanel, messaging, urlshortener, worker) runs as independent process
  - **Limitation:** Current metrics collection only tracks API process (gopsutil limitation), not all 5 processes
  - **Implication:** CPU/memory values reflect API process only; request/error counts are per-service

### Visitor Analytics Design
- **Middleware Placement:** Registered at line 429, before all route registrations
  - **Critical:** Ensures tracking executes for every request, including 404s
  - **Performance:** Minimal overhead (async database writes, no blocking)
- **Privacy-First:** SHA256 session hashing, no raw IP storage, consent tracking
  - **Compliance:** GDPR/CCPA/LGPD ready, user can revoke consent
  - **Limitation:** Cannot reconstruct original IPs from hashes
- **Aggregation Strategy:** Worker cron jobs run hourly/daily, not real-time
  - **Benefit:** Reduces database load, optimizes query performance
  - **Tradeoff:** Aggregate data may lag by up to 1 hour

### Service Management Design
- **State vs Process Control:** Start/Stop in DevPanel updates registration state, not OS processes
  - **Architecture:** 5 independent microservices run as separate processes with own binaries
  - **ServiceManager Role:** Centralized registry and monitoring in API service, not process orchestrator
  - **Process Control:** Use tmux sessions or systemd to start/stop actual OS processes
  - **DevPanel Controls:** Affect service visibility and registration state, useful for feature flags/debugging
  - **Implication:** "Stopped" services in DevPanel are still running as OS processes, just marked inactive in registry

## Performance Optimization Tips
- Use React.memo for expensive components, implement virtual scrolling for long lists
- Enable gzip compression in production (already configured in api/main.go line 366)
- Use CDN for static assets in production deployments
- Database query optimization: Use EXPLAIN ANALYZE, add indexes to foreign keys
- **In-Memory Metrics:** Service analytics optimized for speed over persistence (7-day retention, automatic cleanup)
- **Batch Operations:** Worker aggregates visitor data in batches to minimize database queries
- **Connection Pooling:** PostgreSQL connection pool configured for optimal throughput

## File Maintenance Notice
**IMPORTANT:** When adding new features or refactoring, update this file with objectively important information that would help future development. Keep additions concise and focused on essential patterns, commands, or architectural decisions. Remember to maintain 100% accuracy in this documentation.
- The .gitignore file should have .gitignore in it because I do not desire for the file to be saved to github specifically, so do not EVER remove .gitignore from the gitignore file.
- credentials should generally be gathered from .env or similar files.