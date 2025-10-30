package websocket

import (
	"sync"
	"time"
)

// WebSocketConfig holds websocket configuration
type WebSocketConfig struct {
	MaxConnections     int
	RateLimitPerMinute int
	ConnectionTimeout  time.Duration
}

// ConnectionManager handles WebSocket connection limits and rate limiting
type ConnectionManager struct {
	// Configuration settings
	config *WebSocketConfig

	// Active connections count
	activeConnections int

	// Map of client IPs to rate limit data
	rateLimits map[string]*RateLimitData

	// Mutex for thread-safe access to connections count
	connMutex sync.RWMutex

	// Mutex for thread-safe access to rate limits
	rateMutex sync.RWMutex
}

// RateLimitData tracks connection attempts for rate limiting
type RateLimitData struct {
	LastAttempt  time.Time
	AttemptCount int
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(cfg *WebSocketConfig) *ConnectionManager {
	return &ConnectionManager{
		config:            cfg,
		activeConnections: 0,
		rateLimits:        make(map[string]*RateLimitData),
	}
}

// CanConnect checks if a new connection is allowed based on connection limits
func (m *ConnectionManager) CanConnect() bool {
	m.connMutex.RLock()
	defer m.connMutex.RUnlock()

	return m.activeConnections < m.config.MaxConnections
}

// AddConnection increments the active connection count
func (m *ConnectionManager) AddConnection() {
	m.connMutex.Lock()
	defer m.connMutex.Unlock()

	m.activeConnections++
}

// RemoveConnection decrements the active connection count
func (m *ConnectionManager) RemoveConnection() {
	m.connMutex.Lock()
	defer m.connMutex.Unlock()

	if m.activeConnections > 0 {
		m.activeConnections--
	}
}

// GetConnectionCount returns the current number of active connections
func (m *ConnectionManager) GetConnectionCount() int {
	m.connMutex.RLock()
	defer m.connMutex.RUnlock()

	return m.activeConnections
}

// CheckRateLimit checks if an IP is rate limited
// Returns true if allowed, false if rate limited
func (m *ConnectionManager) CheckRateLimit(ip string) bool {
	m.rateMutex.Lock()
	defer m.rateMutex.Unlock()

	now := time.Now()

	// Get existing data or create new
	data, exists := m.rateLimits[ip]
	if !exists {
		data = &RateLimitData{
			LastAttempt:  now,
			AttemptCount: 1,
		}
		m.rateLimits[ip] = data
		return true
	}

	// Reset counter if window has passed (1 minute)
	if now.Sub(data.LastAttempt) > time.Minute {
		data.AttemptCount = 1
		data.LastAttempt = now
		return true
	}

	// Increment counter and check limit
	data.AttemptCount++
	data.LastAttempt = now

	// Limit to 10 connection attempts per minute
	return data.AttemptCount <= 10
}

// CleanupStaleLimits removes stale rate limit data
// Should be called periodically from a goroutine
func (m *ConnectionManager) CleanupStaleLimits() {
	m.rateMutex.Lock()
	defer m.rateMutex.Unlock()

	now := time.Now()

	for ip, data := range m.rateLimits {
		// Remove entries older than 5 minutes
		if now.Sub(data.LastAttempt) > 5*time.Minute {
			delete(m.rateLimits, ip)
		}
	}
}
