package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/common/response"
	"github.com/gin-gonic/gin"
)

// ErrorRecovery recovers from panics and returns a proper error response
func ErrorRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("RequestID")

				// Get stack trace
				stack := debug.Stack()

				// Log the panic
				logger.Error("Panic recovered",
					"request_id", requestID,
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"stack", string(stack),
				)

				// Return error response
				response.SendInternalError(c,
					"An unexpected error occurred. Please try again later.",
					fmt.Errorf("panic: %v", err),
				)

				// Abort the request
				c.Abort()
			}
		}()

		c.Next()
	}
}
