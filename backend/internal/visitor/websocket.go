package visitor

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var visitorAllowedOrigins []string

func init() {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins != "" {
		visitorAllowedOrigins = strings.Split(origins, ",")
	} else if os.Getenv("APP_ENV") == "production" {
		visitorAllowedOrigins = []string{"https://jadenrazo.dev", "https://www.jadenrazo.dev"}
	} else {
		visitorAllowedOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	}
}

func isVisitorOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	for _, allowed := range visitorAllowedOrigins {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}
	return false
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return isVisitorOriginAllowed(origin)
	},
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.Mutex
	done       chan struct{}
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
		done:       make(chan struct{}),
	}
}

// Run starts the hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			h.mutex.Lock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					go func(c *websocket.Conn) {
						h.unregister <- c
					}(client)
				}
			}
			h.mutex.Unlock()
		case <-h.done:
			h.mutex.Lock()
			for client := range h.clients {
				client.Close()
			}
			h.mutex.Unlock()
			return
		}
	}
}

// Stop gracefully stops the hub and closes all connections
func (h *Hub) Stop() {
	close(h.done)
}

// ServeWs handles websocket requests from the peer.
func (s *Service) ServeWs(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	hub.register <- conn

	// Unregister the client when the connection is closed
	defer func() {
		hub.unregister <- conn
	}()

	for {
		// Read message from browser
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
	}
}
