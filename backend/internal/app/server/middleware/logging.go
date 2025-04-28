package middleware

import (
	"net/http"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/google/uuid"
)

// RequestLogger is middleware that logs the start and end of each request
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a request ID if not already present
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			r.Header.Set("X-Request-ID", requestID)
		}

		// Create a response wrapper to capture status code
		ww := NewResponseWriter(w)

		// Log request start
		logger.Info("Request started",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		// Call the next handler
		next.ServeHTTP(ww, r)

		// Calculate request duration
		duration := time.Since(start)

		// Log request completion
		logger.Info("Request completed",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration_ms", duration.Milliseconds(),
			"size", ww.BytesWritten(),
		)
	})
}

// ResponseWriter is a wrapper for http.ResponseWriter that captures status code and response size
type ResponseWriter struct {
	http.ResponseWriter
	status       int
	bytesWritten int64
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		status:         http.StatusOK, // Default status code is 200
	}
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the response size
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += int64(size)
	return size, err
}

// Status returns the status code
func (rw *ResponseWriter) Status() int {
	return rw.status
}

// BytesWritten returns the number of bytes written
func (rw *ResponseWriter) BytesWritten() int64 {
	return rw.bytesWritten
}

// Flush implements the http.Flusher interface
func (rw *ResponseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements the http.Hijacker interface
func (rw *ResponseWriter) Hijack() (interface{}, interface{}, error) {
	if h, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}
