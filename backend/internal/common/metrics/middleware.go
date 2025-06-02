package metrics

import (
	"net/http"
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
			Timeout: 5 * time.Second,
		}

		for range ticker.C {
			start := time.Now()
			
			resp, err := client.Get(healthCheckURL)
			latency := float64(time.Since(start).Nanoseconds()) / 1e6
			
			if err != nil {
				// Still record the latency even if there's an error (shows system issues)
				m.RecordLatency(latency, "/api/status")
				continue
			}
			resp.Body.Close()
			
			// Record latency regardless of status code to show real response times
			m.RecordLatency(latency, "/api/status")
		}
	}()
}