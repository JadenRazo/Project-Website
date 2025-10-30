package security

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type SecurityManager struct {
	config      SecurityConfig
	ipWhitelist *IPWhitelist
	rateLimiter *RateLimiter
	auditLogger *AuditLogger
}

func NewSecurityManager(config SecurityConfig) (*SecurityManager, error) {
	// Initialize IP whitelist
	ipWhitelist, err := NewIPWhitelist(config.DevPanel.IPWhitelist)
	if err != nil {
		return nil, err
	}

	// Initialize rate limiter
	rateLimiter := NewRateLimiter(config.RateLimit)

	// Initialize audit logger
	auditLogger, err := NewAuditLogger(config.AuditLog)
	if err != nil {
		return nil, err
	}

	return &SecurityManager{
		config:      config,
		ipWhitelist: ipWhitelist,
		rateLimiter: rateLimiter,
		auditLogger: auditLogger,
	}, nil
}

func (sm *SecurityManager) DevPanelMiddleware() []gin.HandlerFunc {
	middlewares := []gin.HandlerFunc{}

	// Add security headers
	middlewares = append(middlewares, sm.securityHeadersMiddleware())

	// Add audit logging (should be first to capture all requests)
	if sm.auditLogger.IsEnabled() {
		middlewares = append(middlewares, sm.auditLogger.Middleware())
	}

	// Add IP whitelist (should be early to block unauthorized IPs)
	if sm.ipWhitelist.IsEnabled() {
		middlewares = append(middlewares, sm.ipWhitelist.Middleware())
	}

	// Add rate limiting
	middlewares = append(middlewares, sm.rateLimiter.Middleware())

	return middlewares
}

func (sm *SecurityManager) AuthMiddleware() []gin.HandlerFunc {
	middlewares := []gin.HandlerFunc{}

	// Add stricter rate limiting for auth endpoints
	middlewares = append(middlewares, sm.rateLimiter.AuthMiddleware())

	// Add HTTPS redirect for production
	if sm.config.IsProduction() && sm.config.DevPanel.RequireTLS {
		middlewares = append(middlewares, sm.httpsRedirectMiddleware())
	}

	return middlewares
}

func (sm *SecurityManager) GetCORSConfig() cors.Config {
	config := cors.DefaultConfig()

	config.AllowOrigins = sm.config.CORS.AllowedOrigins
	config.AllowMethods = sm.config.CORS.AllowedMethods
	config.AllowHeaders = sm.config.CORS.AllowedHeaders
	config.AllowCredentials = sm.config.CORS.AllowCredentials
	config.MaxAge = time.Duration(sm.config.CORS.MaxAge) * time.Second

	// In localhost-only mode, restrict to local origins
	if sm.config.DevPanel.LocalhostOnly {
		localOrigins := []string{}
		for _, origin := range config.AllowOrigins {
			if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
				localOrigins = append(localOrigins, origin)
			}
		}
		config.AllowOrigins = localOrigins
	}

	return config
}

func (sm *SecurityManager) securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy for admin panel
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self'; " +
			"connect-src 'self'"
		c.Header("Content-Security-Policy", csp)

		// HSTS for production
		if sm.config.IsProduction() && sm.config.DevPanel.RequireTLS {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}

func (sm *SecurityManager) httpsRedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Forwarded-Proto") != "https" && c.Request.TLS == nil {
			httpsURL := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(http.StatusMovedPermanently, httpsURL)
			c.Abort()
			return
		}
		c.Next()
	}
}

func (sm *SecurityManager) LogAdminAction(c *gin.Context, action, resource string, details map[string]interface{}) {
	sm.auditLogger.LogAdminAction(c, action, resource, details)
}

func (sm *SecurityManager) LogSecurityEvent(clientIP, action string, details map[string]interface{}) {
	sm.auditLogger.LogSecurityEvent(clientIP, action, details)
}

func (sm *SecurityManager) GetBindAddress(port int) string {
	return sm.config.GetBindAddress(port)
}

func (sm *SecurityManager) IsLocalhostOnly() bool {
	return sm.config.DevPanel.LocalhostOnly
}

func (sm *SecurityManager) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"ip_whitelist": map[string]interface{}{
			"enabled":        sm.ipWhitelist.IsEnabled(),
			"localhost_only": sm.ipWhitelist.IsLocalhostOnly(),
			"allowed_ips":    sm.ipWhitelist.GetAllowedIPs(),
			"allowed_cidrs":  sm.ipWhitelist.GetAllowedCIDRs(),
		},
		"rate_limiter": sm.rateLimiter.GetStats(),
		"audit_log": map[string]interface{}{
			"enabled": sm.auditLogger.IsEnabled(),
		},
		"environment": sm.config.Environment,
	}
}

func (sm *SecurityManager) Cleanup() {
	if sm.rateLimiter != nil {
		sm.rateLimiter.Stop()
	}
}

// Implement AuditLogger interface
func (sm *SecurityManager) LogAuthEvent(clientIP, action, userID string, success bool, details map[string]interface{}) {
	sm.auditLogger.LogAuthEvent(clientIP, action, userID, success, details)
}

// Implement SessionStore interface (NoOp implementation for now)
func (sm *SecurityManager) CreateSession(userID, sessionID string, expiresAt time.Time) error {
	return nil
}

func (sm *SecurityManager) ValidateSession(sessionID string) (bool, error) {
	return true, nil
}

func (sm *SecurityManager) InvalidateSession(sessionID string) error {
	return nil
}

func (sm *SecurityManager) InvalidateUserSessions(userID string) error {
	return nil
}
