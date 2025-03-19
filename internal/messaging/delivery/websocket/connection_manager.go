package websocket

import (
	"sync"
	"time"
)

// ConnectionManager manages WebSocket connections with rate limiting
type ConnectionManager struct {
	// Maximum number of connections allowed per IP
	maxConnections int

	// Map of IP addresses to connection counts
	connections map[string]int

	// Map of IP addresses to last connection time
	lastConnection map[string]time.Time

	// Mutex to protect the maps
	mu sync.RWMutex
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(maxConnections int) *ConnectionManager {
	return &ConnectionManager{
		maxConnections: maxConnections,
		connections:    make(map[string]int),
		lastConnection: make(map[string]time.Time),
		mu:             sync.RWMutex{},
	}
}

// CanConnect checks if an IP is allowed to connect
func (m *ConnectionManager) CanConnect(ip string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := m.connections[ip]
	lastConn, exists := m.lastConnection[ip]

	// Rate limiting: allow connection if it's been more than 5 seconds since last connection
	// or if the client hasn't hit the connection limit
	return !exists || time.Since(lastConn) > 5*time.Second || count < m.maxConnections
}

// AddConnection records a new connection from an IP
func (m *ConnectionManager) AddConnection(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.connections[ip]++
	m.lastConnection[ip] = time.Now()
}

// RemoveConnection removes a connection for an IP
func (m *ConnectionManager) RemoveConnection(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connections[ip] > 0 {
		m.connections[ip]--
	}

	// If no more connections from this IP, clean up
	if m.connections[ip] == 0 {
		delete(m.connections, ip)
		delete(m.lastConnection, ip)
	}
}

// GetConnectionCount gets the current connection count for an IP
func (m *ConnectionManager) GetConnectionCount(ip string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.connections[ip]
}
