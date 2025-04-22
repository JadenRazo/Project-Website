package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/auth"
)

// Middleware provides JWT authentication middleware
type Middleware struct {
	authService  *auth.Auth
	maxAttempts  int           // Maximum auth attempts per IP
	windowTime   time.Duration // Time window for rate limiting
	ipCache      map[string]*rateLimitEntry
	cacheMutex   sync.RWMutex
	cleanupTimer *time.Ticker
}

// rateLimitEntry tracks authentication attempts
type rateLimitEntry struct {
	attempts    int
	lastAttempt time.Time
	blocked     bool
	blockUntil  time.Time
}

// New creates a new auth middleware
func New(authService *auth.Auth) *Middleware {
	m := &Middleware{
		authService: authService,
		maxAttempts: 5,               // 5 failed attempts
		windowTime:  5 * time.Minute, // Within 5 minutes
		ipCache:     make(map[string]*rateLimitEntry),
	}

	// Start cache cleanup routine
	m.cleanupTimer = time.NewTicker(10 * time.Minute)
	go m.cleanupCache()

	return m
}

// cleanupCache periodically removes expired entries
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

// Stop stops the middleware cleanup
func (m *Middleware) Stop() {
	if m.cleanupTimer != nil {
		m.cleanupTimer.Stop()
	}
}

// Authenticate is the main authentication middleware
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for preflight requests and public paths
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// Get client IP for rate limiting
		clientIP := getClientIP(r)

		// Check rate limiting
		if m.isRateLimited(clientIP) {
			w.Header().Set("Retry-After", "300") // 5 minutes
			http.Error(w, "Too many failed authentication attempts", http.StatusTooManyRequests)
			return
		}

		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.recordFailedAttempt(clientIP)
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// Check that it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			m.recordFailedAttempt(clientIP)
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token with context
		claims, err := m.authService.ValidateToken(r.Context(), tokenString)
		if err != nil {
			m.recordFailedAttempt(clientIP)
			if errors.Is(err, auth.ErrTokenExpired) {
				http.Error(w, "Unauthorized: Token expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			}
			return
		}

		// Reset failed attempts on successful auth
		m.resetAttempts(clientIP)

		// Add user info to context with strong typing
		userKey := auth.ContextKey("user")
		ctx := context.WithValue(r.Context(), userKey, claims)

		// Call the next handler with the enriched context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware checks if user has required role
func (m *Middleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get claims from context using typed key
			userKey := auth.ContextKey("user")
			claims, ok := r.Context().Value(userKey).(*auth.Claims)
			if !ok {
				http.Error(w, "Unauthorized: No user in context", http.StatusUnauthorized)
				return
			}

			// Check role
			if claims.Role != role {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isRateLimited checks if a client is rate limited
func (m *Middleware) isRateLimited(ip string) bool {
	m.cacheMutex.RLock()
	entry, exists := m.ipCache[ip]
	m.cacheMutex.RUnlock()

	if !exists {
		return false
	}

	return entry.blocked && time.Now().Before(entry.blockUntil)
}

// recordFailedAttempt records a failed auth attempt
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
		return
	}

	// Increment attempts
	entry.attempts++
	entry.lastAttempt = now

	// Block if exceeded attempts
	if entry.attempts >= m.maxAttempts {
		entry.blocked = true
		entry.blockUntil = now.Add(m.windowTime)
	}
}

// resetAttempts resets failed attempts for an IP
func (m *Middleware) resetAttempts(ip string) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	delete(m.ipCache, ip)
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	// Try X-Forwarded-For header first
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// Take the first IP if there are multiple
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Try X-Real-IP
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Fall back to RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}
