package metrics

import (
	"net/http"
	"strconv"
	"time"
)

// LatencyMiddleware creates middleware that tracks request latency
func (m *Manager) LatencyMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !m.config.Enabled || m.latencyTracker == nil {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			
			// Create a response writer wrapper to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			next.ServeHTTP(wrapped, r)
			
			// Calculate latency in milliseconds
			latency := float64(time.Since(start).Nanoseconds()) / 1e6
			
			// Record latency for successful requests only
			if wrapped.statusCode < 400 {
				m.RecordLatency(latency, r.URL.Path)
			}
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// HealthCheckMiddleware specifically tracks health check latency
func (m *Manager) HealthCheckMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			next.ServeHTTP(w, r)
			
			// Always record health check latency
			latency := float64(time.Since(start).Nanoseconds()) / 1e6
			if m.config.Enabled && m.latencyTracker != nil {
				m.RecordLatency(latency, "/health")
			}
		})
	}
}

// PeriodicLatencyCollector runs a background job to collect periodic latency measurements
func (m *Manager) StartPeriodicLatencyCollector(interval time.Duration, healthCheckURL string) {
	if !m.config.Enabled || m.latencyTracker == nil {
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		for range ticker.C {
			start := time.Now()
			
			resp, err := client.Get(healthCheckURL)
			if err != nil {
				continue
			}
			resp.Body.Close()
			
			if resp.StatusCode == http.StatusOK {
				latency := float64(time.Since(start).Nanoseconds()) / 1e6
				m.RecordLatency(latency, "/health-check")
			}
		}
	}()
}