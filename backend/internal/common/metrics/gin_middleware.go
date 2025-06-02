package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
)

// GinLatencyMiddleware creates Gin-compatible middleware that tracks request latency
func (m *Manager) GinLatencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.config.Enabled || m.latencyTracker == nil {
			c.Next()
			return
		}

		start := time.Now()
		
		c.Next()
		
		// Calculate latency in milliseconds
		latency := float64(time.Since(start).Nanoseconds()) / 1e6
		
		// Record latency for successful requests only
		if c.Writer.Status() < 400 {
			m.RecordLatency(latency, c.Request.URL.Path)
		}
	}
}