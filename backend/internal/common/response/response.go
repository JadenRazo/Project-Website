package response

import (
	"encoding/json"
	"net/http"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a structured API error response
type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// SuccessResponse represents a structured API success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ValidationErrorResponse represents a validation error response
type ValidationErrorResponse struct {
	Error            string      `json:"error"`
	ValidationErrors interface{} `json:"validation_errors,omitempty"`
}

// SendError sends a generic error response with the specified status code
func SendError(c *gin.Context, statusCode int, message string, details interface{}) {
	requestID := c.GetString("RequestID")
	
	// Log the error for debugging
	logger.Error("API Error Response",
		"request_id", requestID,
		"status_code", statusCode,
		"message", message,
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
	)

	response := ErrorResponse{
		Error:   message,
		Details: details,
	}

	c.JSON(statusCode, response)
}

// SendValidationError sends a validation error response (400 Bad Request)
func SendValidationError(c *gin.Context, message string, validationErrors interface{}) {
	requestID := c.GetString("RequestID")
	
	// Log the validation error
	logger.Warn("Validation Error",
		"request_id", requestID,
		"message", message,
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
	)

	response := ValidationErrorResponse{
		Error:            message,
		ValidationErrors: validationErrors,
	}

	c.JSON(http.StatusBadRequest, response)
}

// SendInternalError sends an internal server error response (500 Internal Server Error)
func SendInternalError(c *gin.Context, message string, err error) {
	requestID := c.GetString("RequestID")
	
	// Log the internal error with more details
	errorDetails := "unknown error"
	if err != nil {
		errorDetails = err.Error()
	}
	
	logger.Error("Internal Server Error",
		"request_id", requestID,
		"message", message,
		"error_details", errorDetails,
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
	)

	response := ErrorResponse{
		Error: message,
		// Don't expose internal error details in production
		Details: nil,
	}

	c.JSON(http.StatusInternalServerError, response)
}

// SendSuccess sends a successful response with data
func SendSuccess(c *gin.Context, data interface{}) {
	response := SuccessResponse{
		Success: true,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// SendSuccessWithMessage sends a successful response with data and message
func SendSuccessWithMessage(c *gin.Context, data interface{}, message string) {
	response := SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	c.JSON(http.StatusOK, response)
}

// SendCreated sends a 201 Created response with data
func SendCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// SendNoContent sends a 204 No Content response
func SendNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Standard HTTP Response Utilities (for non-Gin handlers)

// writeJSONResponse writes a JSON response to standard http.ResponseWriter
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// If encoding fails, try to send a basic error response
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error - failed to encode response"}`))  
		}
	}
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	response := ErrorResponse{
		Error: message,
	}
	writeJSONResponse(w, http.StatusUnauthorized, response)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Forbidden"
	}
	response := ErrorResponse{
		Error: message,
	}
	writeJSONResponse(w, http.StatusForbidden, response)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Bad Request"
	}
	response := ErrorResponse{
		Error: message,
	}
	writeJSONResponse(w, http.StatusBadRequest, response)
}

// NotFound sends a 404 Not Found response
func NotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Not Found"
	}
	response := ErrorResponse{
		Error: message,
	}
	writeJSONResponse(w, http.StatusNotFound, response)
}

// MethodNotAllowed sends a 405 Method Not Allowed response
func MethodNotAllowed(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Method Not Allowed"
	}
	response := ErrorResponse{
		Error: message,
	}
	writeJSONResponse(w, http.StatusMethodNotAllowed, response)
}

// Conflict sends a 409 Conflict response
func Conflict(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Conflict"
	}
	response := ErrorResponse{
		Error: message,
	}
	writeJSONResponse(w, http.StatusConflict, response)
}

// TooManyRequests sends a 429 Too Many Requests response
func TooManyRequests(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Too Many Requests"
	}
	response := ErrorResponse{
		Error: message,
	}
	writeJSONResponse(w, http.StatusTooManyRequests, response)
}

// InternalError sends a 500 Internal Server Error response
func InternalError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Internal Server Error"
	}
	response := ErrorResponse{
		Error: message,
	}
	writeJSONResponse(w, http.StatusInternalServerError, response)
}

// ServiceUnavailable sends a 503 Service Unavailable response
func ServiceUnavailable(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Service Unavailable"
	}
	response := ErrorResponse{
		Error: message,
	}
	writeJSONResponse(w, http.StatusServiceUnavailable, response)
}

// Success sends a 200 OK response with data
func Success(w http.ResponseWriter, data interface{}) {
	writeJSONResponse(w, http.StatusOK, data)
}

// Created sends a 201 Created response with data
func Created(w http.ResponseWriter, data interface{}) {
	writeJSONResponse(w, http.StatusCreated, data)
}

// NoContent sends a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// JSON sends a custom status code with JSON data
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	writeJSONResponse(w, statusCode, data)
}

// Error sends a custom error response with status code
func Error(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{
		Error: message,
	}
	writeJSONResponse(w, statusCode, response)
}