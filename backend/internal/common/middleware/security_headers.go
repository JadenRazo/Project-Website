package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders returns a middleware that sets security-related HTTP headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		env := os.Getenv("APP_ENV")
		isProduction := env == "production"

		c.Header("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' https: data:; "+
				"font-src 'self' data:; "+
				"connect-src 'self' https://www.googleapis.com https://api.github.com https://login.microsoftonline.com https://graph.microsoft.com; "+
				"frame-ancestors 'none'; "+
				"base-uri 'self'; "+
				"form-action 'self'")

		c.Header("X-Content-Type-Options", "nosniff")

		c.Header("X-Frame-Options", "DENY")

		c.Header("X-XSS-Protection", "1; mode=block")

		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		if isProduction {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		c.Header("Cross-Origin-Embedder-Policy", "require-corp")
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Cross-Origin-Resource-Policy", "same-origin")

		c.Next()
	}
}
