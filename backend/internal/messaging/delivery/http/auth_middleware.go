package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MessageUserIDKey is the context key for user ID in message handlers
const MessageUserIDKey = "message_user_id"

// GetMessageUserID extracts the user ID from the context
func GetMessageUserID(ctx context.Context) (uint, bool) {
	value, exists := ctx.Value(MessageUserIDKey).(uint)
	return value, exists
}

// MessageAuthMiddleware is a middleware that extracts the user ID from the request context
// and sets it in the request context for handlers to use
func MessageAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user ID is in the context
		// This assumes that an authentication middleware has already set the user ID in the context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}

		// Check if user ID is a uint
		userIDValue, ok := userID.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}

		// Set user ID in context with our key
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), MessageUserIDKey, userIDValue))

		c.Next()
	}
}
