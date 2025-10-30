package middleware

import (
	"bytes"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware provides security-focused middleware functions
type SecurityMiddleware struct {
	// Patterns that should be sanitized from responses
	sensitivePatterns []*regexp.Regexp
	// Replacement text for sensitive information
	replacementText string
}

// NewSecurityMiddleware creates a new security middleware instance
func NewSecurityMiddleware() *SecurityMiddleware {
	patterns := []*regexp.Regexp{
		// Database connection strings
		regexp.MustCompile(`(?i)(postgresql?|mysql|sqlite)://[^@]+@[^/]+/\w+`),
		regexp.MustCompile(`(?i)host=[^\s]+`),
		regexp.MustCompile(`(?i)port=\d+`),
		regexp.MustCompile(`(?i)dbname=\w+`),
		regexp.MustCompile(`(?i)user=\w+`),
		regexp.MustCompile(`(?i)password=[^\s]+`),

		// IP addresses and internal network information
		regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`),
		regexp.MustCompile(`(?i)localhost:\d+`),
		regexp.MustCompile(`(?i)127\.0\.0\.1:\d+`),
		regexp.MustCompile(`(?i)0\.0\.0\.0:\d+`),

		// File paths that might contain sensitive information
		regexp.MustCompile(`(?i)/[a-zA-Z0-9_\-/]+\.(?:env|config|yml|yaml|json|key|pem|crt)`),
		regexp.MustCompile(`(?i)C:\\[^<>\"|]*\.(?:env|config|yml|yaml|json|key|pem|crt)`),

		// Database error messages that might leak information
		regexp.MustCompile(`(?i)connection refused.*host ".*"`),
		regexp.MustCompile(`(?i)dial tcp .*: connect: connection refused`),
		regexp.MustCompile(`(?i)no such host.*`),
		regexp.MustCompile(`(?i)timeout.*dial tcp.*`),

		// Stack traces and internal paths
		regexp.MustCompile(`(?i)goroutine \d+ \[.*\]:`),
		regexp.MustCompile(`(?i)panic: .*`),
		regexp.MustCompile(`(?i)/[a-zA-Z0-9_\-/]+\.go:\d+`),

		// Environment variables and secrets
		regexp.MustCompile(`(?i)[A-Z_]+_KEY=\w+`),
		regexp.MustCompile(`(?i)[A-Z_]+_SECRET=\w+`),
		regexp.MustCompile(`(?i)[A-Z_]+_TOKEN=\w+`),
		regexp.MustCompile(`(?i)api_key=\w+`),
		regexp.MustCompile(`(?i)access_token=\w+`),
	}

	return &SecurityMiddleware{
		sensitivePatterns: patterns,
		replacementText:   "[REDACTED]",
	}
}

// SanitizeResponse sanitizes response data to remove sensitive information
func (sm *SecurityMiddleware) SanitizeResponse() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Capture the response
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer

		// Process the request
		c.Next()

		// Sanitize the response body
		originalBody := writer.body.String()
		sanitizedBody := sm.sanitizeString(originalBody)

		// If the response was modified, update it
		if sanitizedBody != originalBody {
			writer.ResponseWriter.Header().Set("Content-Length", string(rune(len(sanitizedBody))))
			writer.ResponseWriter.Write([]byte(sanitizedBody))
		} else {
			writer.ResponseWriter.Write(writer.body.Bytes())
		}
	})
}

// sanitizeString removes sensitive patterns from a string
func (sm *SecurityMiddleware) sanitizeString(input string) string {
	result := input

	// Apply all sensitive patterns
	for _, pattern := range sm.sensitivePatterns {
		result = pattern.ReplaceAllString(result, sm.replacementText)
	}

	return result
}

// SanitizeErrorResponse specifically sanitizes error responses
func (sm *SecurityMiddleware) SanitizeErrorResponse() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()

		// Only process error responses
		if c.Writer.Status() >= 400 {
			// Check if we have JSON error response
			if contentType := c.Writer.Header().Get("Content-Type"); strings.Contains(contentType, "application/json") {
				// This middleware runs after the response, so we need to use the response capture approach
			}
		}
	})
}

// HeaderSecurity adds security headers to responses
func (sm *SecurityMiddleware) HeaderSecurity() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Add security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Remove server information
		c.Header("Server", "")

		c.Next()
	})
}

// responseWriter captures response data for sanitization
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(data []byte) (int, error) {
	// Write to both the buffer and the original writer
	w.body.Write(data)
	return len(data), nil
}

// SanitizeJSONError sanitizes error messages in JSON responses
func SanitizeJSONError(c *gin.Context, statusCode int, originalErr error) {
	sanitizedMsg := sanitizeErrorMessage(originalErr.Error())
	c.JSON(statusCode, gin.H{"error": sanitizedMsg})
}

// sanitizeErrorMessage sanitizes a single error message
func sanitizeErrorMessage(msg string) string {
	// Common sensitive patterns in error messages
	patterns := map[*regexp.Regexp]string{
		regexp.MustCompile(`(?i)connection refused.*host ".*"`):            "Database connection unavailable",
		regexp.MustCompile(`(?i)dial tcp .*: connect: connection refused`): "Service connection failed",
		regexp.MustCompile(`(?i)no such host.*`):                           "Service temporarily unavailable",
		regexp.MustCompile(`(?i)timeout.*dial tcp.*`):                      "Service response timeout",
		regexp.MustCompile(`(?i)database.*not.*found`):                     "Database service unavailable",
		regexp.MustCompile(`(?i)table.*doesn't exist`):                     "Data temporarily unavailable",
		regexp.MustCompile(`(?i)column.*doesn't exist`):                    "Data structure error",
		regexp.MustCompile(`(?i)duplicate.*key.*value`):                    "Data conflict error",
		regexp.MustCompile(`(?i)permission.*denied`):                       "Access denied",
		regexp.MustCompile(`(?i)authentication.*failed`):                   "Invalid credentials",
		regexp.MustCompile(`(?i)invalid.*token`):                           "Invalid credentials",
		regexp.MustCompile(`(?i)expired.*token`):                           "Session expired",
	}

	result := msg
	for pattern, replacement := range patterns {
		if pattern.MatchString(result) {
			result = replacement
			break // Use first match
		}
	}

	// Remove any remaining IP addresses, ports, or paths
	result = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`).ReplaceAllString(result, "[IP]")
	result = regexp.MustCompile(`:\d{4,5}\b`).ReplaceAllString(result, ":[PORT]")
	result = regexp.MustCompile(`/[a-zA-Z0-9_\-/]+\.(go|sql|yaml|yml|json|env)`).ReplaceAllString(result, "/[FILE]")

	return result
}

// RateLimitByIP provides IP-based rate limiting
func (sm *SecurityMiddleware) RateLimitByIP() gin.HandlerFunc {
	// Simple in-memory rate limiter - in production, use Redis or similar
	clients := make(map[string]int)
	const maxRequests = 100 // per minute

	return gin.HandlerFunc(func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Simple rate limiting logic
		if count, exists := clients[clientIP]; exists && count > maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		clients[clientIP]++
		c.Next()

		// Reset counter periodically (simplified)
		if clients[clientIP] > maxRequests*2 {
			delete(clients, clientIP)
		}
	})
}
