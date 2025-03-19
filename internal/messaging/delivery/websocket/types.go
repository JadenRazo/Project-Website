package websocket

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Connection defines the interface for WebSocket connections
type Connection interface {
	// ReadMessage reads a message from the connection
	ReadMessage() (messageType int, p []byte, err error)

	// WriteMessage writes a message to the connection
	WriteMessage(messageType int, data []byte) error

	// Close closes the connection
	Close() error

	// SetReadDeadline sets the read deadline on the connection
	SetReadDeadline(t time.Time) error

	// SetWriteDeadline sets the write deadline on the connection
	SetWriteDeadline(t time.Time) error

	// SetPongHandler sets the handler for pong messages
	SetPongHandler(h func(appData string) error)

	// SetPingHandler sets the handler for ping messages
	SetPingHandler(h func(appData string) error)
}

// ClientList is a map of connected clients with thread-safe access
type ClientList struct {
	sync.RWMutex
	clients map[*Client]bool
}

// Upgrader configures WebSocket upgrade parameters
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for now, in production this should be restricted
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketConnection is an implementation of the Connection interface using gorilla/websocket
type WebSocketConnection struct {
	*websocket.Conn
}

// NewWebSocketConnection creates a new WebSocketConnection
func NewWebSocketConnection(conn *websocket.Conn) Connection {
	return &WebSocketConnection{Conn: conn}
}

// EventTypes for WebSocket messages
const (
	EventTypeMessage          = "message"
	EventTypeChannelSubscribe = "channel_subscribe"
	EventTypePresenceUpdate   = "presence_update"
	EventTypeTyping           = "typing"
	EventTypeReadReceipt      = "read_receipt"
	EventTypeError            = "error"
	EventTypeAck              = "ack"
)

// ClientMessage represents a message from a client
type ClientMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	ID   string      `json:"id,omitempty"`
}

// ErrorResponse represents an error message sent to clients
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
