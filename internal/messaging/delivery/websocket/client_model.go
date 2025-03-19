package websocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client
type Client struct {
	hub      *Hub                   // Reference to the hub
	conn     Connection             // WebSocket connection
	send     chan []byte            // Buffered channel of outbound messages
	UserID   uint                   // User ID of the connected client
	Channels []uint                 // Channels the client is subscribed to
	IP       string                 // IP address of the client for rate limiting
	Metadata map[string]interface{} // Additional metadata about the client
}

// ConnectionManager manages connection rate limiting
type ConnectionManager struct {
	// Maps IP addresses to connection counts and timestamps
	connections map[string]struct {
		count     int       // Number of connections
		timestamp time.Time // Last connection time
	}
	mu sync.Mutex // Mutex for thread safety

	// Rate limiting settings
	maxConnPerIP     int           // Maximum connections per IP
	connectionWindow time.Duration // Time window for rate limiting
}

// NewConnectionManager creates a new connection manager with default settings
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]struct {
			count     int
			timestamp time.Time
		}),
		maxConnPerIP:     50,              // 50 connections per IP
		connectionWindow: 5 * time.Minute, // 5-minute window
	}
}

// CanConnect checks if an IP can establish a new connection
func (m *ConnectionManager) CanConnect(ip string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	info, exists := m.connections[ip]

	// If no record or expired record, allow connection
	if !exists || now.Sub(info.timestamp) > m.connectionWindow {
		return true
	}

	// Check if under rate limit
	return info.count < m.maxConnPerIP
}

// AddConnection records a new connection from an IP
func (m *ConnectionManager) AddConnection(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	info, exists := m.connections[ip]

	if !exists || now.Sub(info.timestamp) > m.connectionWindow {
		// New or expired record
		m.connections[ip] = struct {
			count     int
			timestamp time.Time
		}{
			count:     1,
			timestamp: now,
		}
	} else {
		// Update existing record
		info.count++
		info.timestamp = now
		m.connections[ip] = info
	}
}

// RemoveConnection decreases the connection count for an IP
func (m *ConnectionManager) RemoveConnection(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if info, exists := m.connections[ip]; exists {
		if info.count > 1 {
			info.count--
			m.connections[ip] = info
		} else {
			delete(m.connections, ip)
		}
	}
}

// NewWebSocketConnection creates a new WebSocketConnection that implements Connection
func NewWebSocketConnection(conn *websocket.Conn) Connection {
	return &WebSocketConnection{
		conn: conn,
	}
}

// WebSocketConnection implements the Connection interface for gorilla/websocket
type WebSocketConnection struct {
	conn *websocket.Conn
}

// Close closes the underlying WebSocket connection
func (c *WebSocketConnection) Close() error {
	return c.conn.Close()
}

// ReadMessage reads a message from the WebSocket connection
func (c *WebSocketConnection) ReadMessage() (messageType int, p []byte, err error) {
	return c.conn.ReadMessage()
}

// WriteMessage writes a message to the WebSocket connection
func (c *WebSocketConnection) WriteMessage(messageType int, data []byte) error {
	return c.conn.WriteMessage(messageType, data)
}

// SetReadLimit sets the maximum size for a message read from the peer
func (c *WebSocketConnection) SetReadLimit(limit int64) {
	c.conn.SetReadLimit(limit)
}

// SetReadDeadline sets the read deadline on the underlying network connection
func (c *WebSocketConnection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the write deadline on the underlying network connection
func (c *WebSocketConnection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// SetPongHandler sets the handler for pong messages received from the peer
func (c *WebSocketConnection) SetPongHandler(h func(appData string) error) {
	c.conn.SetPongHandler(h)
}
