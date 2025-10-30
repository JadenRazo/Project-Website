package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/response"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// Context keys for user information
	UserKey      ContextKey = "user"
	UserIDKey    ContextKey = "user_id"
	UserRoleKey  ContextKey = "user_role"
	UserEmailKey ContextKey = "user_email"
)

// Middleware provides comprehensive authentication middleware with rate limiting
type Middleware struct {
	jwtManager   *JWTManager
	maxAttempts  int
	windowTime   time.Duration
	blockTime    time.Duration
	ipCache      map[string]*rateLimitEntry
	cacheMutex   sync.RWMutex
	cleanupTimer *time.Ticker
	publicPaths  []string
}

// rateLimitEntry tracks authentication attempts for rate limiting
type rateLimitEntry struct {
	attempts    int
	lastAttempt time.Time
	blocked     bool
	blockUntil  time.Time
}

// MiddlewareConfig holds configuration for the auth middleware
type MiddlewareConfig struct {
	MaxAttempts int           // Maximum failed attempts before blocking
	WindowTime  time.Duration // Time window for counting attempts
	BlockTime   time.Duration // How long to block after max attempts
	PublicPaths []string      // Paths that don't require authentication
}

// NewMiddleware creates a comprehensive authentication middleware
func NewMiddleware(jwtManager *JWTManager, config *MiddlewareConfig) *Middleware {
	if config == nil {
		config = &MiddlewareConfig{
			MaxAttempts: 5,
			WindowTime:  5 * time.Minute,
			BlockTime:   15 * time.Minute,
			PublicPaths: []string{
				"/api/auth/login",
				"/api/auth/register",
				"/api/auth/refresh",
				"/api/health",
				"/metrics",
				"/swagger",
			},
		}
	}

	m := &Middleware{
		jwtManager:  jwtManager,
		maxAttempts: config.MaxAttempts,
		windowTime:  config.WindowTime,
		blockTime:   config.BlockTime,
		ipCache:     make(map[string]*rateLimitEntry),
		publicPaths: config.PublicPaths,
	}

	// Start cleanup routine
	m.cleanupTimer = time.NewTicker(10 * time.Minute)
	go m.cleanupCache()

	return m
}

// RequireAuth is the main authentication middleware
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// Check if path is public
		if m.isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Get client IP for rate limiting
		clientIP := getClientIP(r)

		// Check rate limiting
		if m.isRateLimited(clientIP) {
			w.Header().Set("Retry-After", "900") // 15 minutes
			response.TooManyRequests(w, "Too many failed authentication attempts. Please try again later.")
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		token, err := ExtractTokenFromHeader(authHeader)
		if err != nil {
			m.recordFailedAttempt(clientIP)
			m.handleAuthError(w, err)
			return
		}

		// Validate token
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			m.recordFailedAttempt(clientIP)
			m.handleAuthError(w, err)
			return
		}

		// Reset failed attempts on successful auth
		m.resetAttempts(clientIP)

		// Add user info to context
		ctx := m.addUserToContext(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth middleware adds user info to context if token is present but doesn't require it
func (m *Middleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// Try to extract and validate token
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			if token, err := ExtractTokenFromHeader(authHeader); err == nil {
				if claims, err := m.jwtManager.ValidateToken(token); err == nil {
					// Add user info to context if token is valid
					ctx := m.addUserToContext(r.Context(), claims)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RequireRole middleware checks if the authenticated user has the required role
func (m *Middleware) RequireRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, ok := r.Context().Value(UserKey).(*Claims)
			if !ok {
				response.Unauthorized(w, "Unauthorized: No user in context")
				return
			}

			// Check if user has required role
			hasRole := false
			for _, role := range requiredRoles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				response.Forbidden(w, "Forbidden: Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin middleware ensures the user has admin role
func (m *Middleware) RequireAdmin(next http.Handler) http.Handler {
	return m.RequireRole("admin")(next)
}

// addUserToContext adds user claims to request context with multiple keys for compatibility
func (m *Middleware) addUserToContext(ctx context.Context, claims *Claims) context.Context {
	ctx = context.WithValue(ctx, UserKey, claims)
	ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
	ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
	ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
	return ctx
}

// handleAuthError provides consistent error responses for authentication failures
func (m *Middleware) handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrMissingToken):
		response.Unauthorized(w, "Unauthorized: Authorization token required")
	case errors.Is(err, ErrInvalidFormat):
		response.Unauthorized(w, "Unauthorized: Invalid authorization format")
	case errors.Is(err, ErrTokenExpired):
		http.Error(w, "Unauthorized: Token expired", http.StatusUnauthorized)
	case errors.Is(err, ErrInvalidToken):
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
	default:
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// isPublicPath checks if a path should be publicly accessible
func (m *Middleware) isPublicPath(path string) bool {
	for _, publicPath := range m.publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}
	return false
}

// Rate limiting methods
func (m *Middleware) isRateLimited(ip string) bool {
	m.cacheMutex.RLock()
	entry, exists := m.ipCache[ip]
	m.cacheMutex.RUnlock()

	if !exists {
		return false
	}

	return entry.blocked && time.Now().Before(entry.blockUntil)
}

func (m *Middleware) recordFailedAttempt(ip string) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	entry, exists := m.ipCache[ip]
	now := time.Now()

	if !exists {
		m.ipCache[ip] = &rateLimitEntry{
			attempts:    1,
			lastAttempt: now,
		}
		return
	}

	// Reset if outside window
	if now.Sub(entry.lastAttempt) > m.windowTime {
		entry.attempts = 1
		entry.lastAttempt = now
		entry.blocked = false
		return
	}

	// Increment attempts
	entry.attempts++
	entry.lastAttempt = now

	// Block if exceeded attempts
	if entry.attempts >= m.maxAttempts {
		entry.blocked = true
		entry.blockUntil = now.Add(m.blockTime)
	}
}

func (m *Middleware) resetAttempts(ip string) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	delete(m.ipCache, ip)
}

func (m *Middleware) cleanupCache() {
	for range m.cleanupTimer.C {
		m.cacheMutex.Lock()
		now := time.Now()

		for ip, entry := range m.ipCache {
			// Remove entries older than window time if not blocked
			if !entry.blocked && now.Sub(entry.lastAttempt) > m.windowTime {
				delete(m.ipCache, ip)
			}

			// Remove block if block time expired
			if entry.blocked && now.After(entry.blockUntil) {
				entry.blocked = false
				entry.attempts = 0
			}
		}

		m.cacheMutex.Unlock()
	}
}

// Stop stops the middleware cleanup routine
func (m *Middleware) Stop() {
	if m.cleanupTimer != nil {
		m.cleanupTimer.Stop()
	}
}

// Helper function to extract client IP
func getClientIP(r *http.Request) string {
	// Try X-Forwarded-For header first (for load balancers/proxies)
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// Take the first IP if there are multiple
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Try X-Real-IP (nginx)
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if colonPos := strings.LastIndex(ip, ":"); colonPos != -1 {
		ip = ip[:colonPos]
	}
	return ip
}

// GetUserFromContext extracts user claims from context
func GetUserFromContext(ctx context.Context) (*Claims, bool) {
	user, ok := ctx.Value(UserKey).(*Claims)
	return user, ok
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetUserRoleFromContext extracts user role from context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}
