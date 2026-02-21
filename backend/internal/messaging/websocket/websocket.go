package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/response"
	"github.com/gorilla/websocket"
)

var allowedOrigins []string

func init() {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	} else if os.Getenv("APP_ENV") == "production" {
		allowedOrigins = []string{"https://jadenrazo.dev", "https://www.jadenrazo.dev"}
	} else {
		allowedOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	}
}

func isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	for _, allowed := range allowedOrigins {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}
	return false
}

// Common errors
var (
	ErrInvalidMessage   = errors.New("invalid message format")
	ErrConnectionClosed = errors.New("connection closed")
	ErrRateLimited      = errors.New("rate limited")
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for concurrent access to clients map
	mutex sync.RWMutex

	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Hub reference
	hub *Hub

	// User ID
	userID string

	// Rate limiting
	messageCount   int
	lastMessageTs  time.Time
	maxMessageRate int // messages per minute

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

// Message represents a websocket message
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	ID      string          `json:"id"`
	Time    int64           `json:"time"`
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client, 100),
		unregister: make(chan *Client, 100),
		clients:    make(map[*Client]bool),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Run starts the hub with context support for graceful shutdown
func (h *Hub) Run() {
	defer h.cancel()

	for {
		select {
		case <-h.ctx.Done():
			// Shutdown all clients
			h.mutex.Lock()
			for client := range h.clients {
				close(client.send)
				client.cancel()
			}
			h.clients = make(map[*Client]bool)
			h.mutex.Unlock()
			return

		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				client.cancel()
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			if !isValidMessage(message) {
				log.Printf("Invalid message format, skipping broadcast")
				continue
			}

			var toRemove []*Client
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					toRemove = append(toRemove, client)
				}
			}
			h.mutex.RUnlock()

			if len(toRemove) > 0 {
				h.mutex.Lock()
				for _, client := range toRemove {
					if _, ok := h.clients[client]; ok {
						delete(h.clients, client)
						close(client.send)
						client.cancel()
					}
				}
				h.mutex.Unlock()
			}
		}
	}
}

// Stop gracefully shuts down the hub
func (h *Hub) Stop() {
	h.cancel()
}

// isValidMessage validates message format
func isValidMessage(message []byte) bool {
	// Check minimum length
	if len(message) < 2 || len(message) > 65536 {
		return false
	}

	// Validate JSON structure
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		return false
	}

	// Validate required fields
	if msg.Type == "" {
		return false
	}

	return true
}

// ServeWs handles websocket requests from clients
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, userID string) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return isOriginAllowed(origin)
		},
		HandshakeTimeout: 10 * time.Second,
	}

	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		response.BadRequest(w, "Could not open websocket connection")
		return
	}

	// Create client context
	ctx, cancel := context.WithCancel(r.Context())

	// Create a client with rate limiting
	client := &Client{
		hub:            hub,
		conn:           conn,
		send:           make(chan []byte, 256),
		userID:         userID,
		lastMessageTs:  time.Now(),
		maxMessageRate: 60, // 60 messages per minute
		ctx:            ctx,
		cancel:         cancel,
	}

	// Register the client
	client.hub.register <- client

	// Start client goroutines
	go client.readPump()
	go client.writePump()
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Set read deadline
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		select {
		case <-c.ctx.Done():
			return

		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			// Rate limiting
			now := time.Now()
			if now.Sub(c.lastMessageTs) < time.Minute {
				c.messageCount++
				if c.messageCount > c.maxMessageRate {
					log.Printf("Rate limit exceeded for user %s", c.userID)
					// Send rate limit notification
					c.send <- []byte(`{"type":"error","payload":{"code":"rate_limited","message":"Message rate limit exceeded"}}`)
					continue
				}
			} else {
				// Reset counter for new minute
				c.messageCount = 1
				c.lastMessageTs = now
			}

			// Validate message format
			if !isValidMessage(message) {
				// Send error back to client
				c.send <- []byte(`{"type":"error","payload":{"code":"invalid_format","message":"Invalid message format"}}`)
				continue
			}

			// Process the message and broadcast
			c.hub.broadcast <- message
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			// Receive shutdown signal
			c.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutting down"))
			return

		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}

			// Send the message with proper error handling
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping: %v", err)
				return
			}
		}
	}
}
