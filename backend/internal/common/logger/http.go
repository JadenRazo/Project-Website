package logger

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// responseWriter is a wrapper for http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the response size
func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// HTTPLogger is a middleware that logs HTTP requests
func HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a wrapper for the response writer
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process the request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Determine log level based on status code
		var level logrus.Level
		switch {
		case rw.statusCode >= 500:
			level = logrus.ErrorLevel
		case rw.statusCode >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}

		// Log the request
		entry := log.WithFields(logrus.Fields{
			"remote_addr":  r.RemoteAddr,
			"method":       r.Method,
			"path":         r.URL.Path,
			"query":        r.URL.RawQuery,
			"status":       rw.statusCode,
			"size":         rw.size,
			"duration_ms":  duration.Milliseconds(),
			"user_agent":   r.UserAgent(),
			"referer":      r.Referer(),
			"request_id":   r.Header.Get("X-Request-ID"),
			"content_type": r.Header.Get("Content-Type"),
		})

		msg := "HTTP Request"
		switch level {
		case logrus.ErrorLevel:
			entry.Error(msg)
		case logrus.WarnLevel:
			entry.Warn(msg)
		default:
			entry.Info(msg)
		}
	})
}

// RequestLogger returns a logger with request context
func RequestLogger(r *http.Request) *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"remote_addr": r.RemoteAddr,
		"method":      r.Method,
		"path":        r.URL.Path,
		"request_id":  r.Header.Get("X-Request-ID"),
	})
}
