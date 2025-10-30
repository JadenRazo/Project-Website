package middleware

import (
	"net/http"
	"strings"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/common/response"
	"github.com/gin-gonic/gin"
)

// RequestValidation validates incoming requests
func RequestValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("RequestID")

		// Validate Content-Type for POST/PUT/PATCH requests
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodPatch {
			contentType := c.GetHeader("Content-Type")
			if contentType == "" {
				logger.Warn("Missing Content-Type header",
					"request_id", requestID,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
				)
				response.SendValidationError(c, "Content-Type header is required for this request", nil)
				c.Abort()
				return
			}

			// Check if JSON content type for API endpoints
			if strings.HasPrefix(c.Request.URL.Path, "/api/") &&
				!strings.Contains(contentType, "application/json") &&
				!strings.Contains(contentType, "multipart/form-data") {
				logger.Warn("Invalid Content-Type for API request",
					"request_id", requestID,
					"content_type", contentType,
					"path", c.Request.URL.Path,
				)
				response.SendValidationError(c, "Content-Type must be application/json for API requests", nil)
				c.Abort()
				return
			}
		}

		// Validate request size
		if c.Request.ContentLength > 10*1024*1024 { // 10MB limit
			logger.Warn("Request body too large",
				"request_id", requestID,
				"content_length", c.Request.ContentLength,
				"path", c.Request.URL.Path,
			)
			response.SendError(c, http.StatusRequestEntityTooLarge, "Request body too large. Maximum size is 10MB.", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// QueryParamValidation validates query parameters
func QueryParamValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("RequestID")

		// Validate pagination parameters
		if page := c.Query("page"); page != "" {
			if _, err := c.GetQuery("page"); !err {
				logger.Warn("Invalid page parameter",
					"request_id", requestID,
					"page", page,
					"path", c.Request.URL.Path,
				)
			}
		}

		if pageSize := c.Query("page_size"); pageSize != "" {
			if _, err := c.GetQuery("page_size"); !err {
				logger.Warn("Invalid page_size parameter",
					"request_id", requestID,
					"page_size", pageSize,
					"path", c.Request.URL.Path,
				)
			}
		}

		c.Next()
	}
}
