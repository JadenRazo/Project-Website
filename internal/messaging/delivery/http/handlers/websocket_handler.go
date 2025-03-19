package handlers

import (
	"log"
	"net/http"

	"jadenrazo.dev/internal/messaging/delivery/websocket"
	"jadenrazo.dev/internal/messaging/service"
)

// WebSocketHandler handles WebSocket connections for the messaging platform
type WebSocketHandler struct {
	hub              *websocket.Hub
	messagingService service.MessagingService
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(messagingService service.MessagingService) *WebSocketHandler {
	hub := websocket.NewHub()

	// Start the WebSocket hub in a goroutine
	go hub.Run()

	return &WebSocketHandler{
		hub:              hub,
		messagingService: messagingService,
	}
}

// RegisterRoutes registers the WebSocket routes to the router
func (h *WebSocketHandler) RegisterRoutes(router http.Handler) {
	// Use type assertion to access the required methods
	if mux, ok := router.(MuxRouter); ok {
		mux.HandleFunc("/ws", h.HandleWebSocket)
		mux.HandleFunc("/ws/channel/{channelId}", h.HandleChannelWebSocket)

		log.Println("WebSocket routes registered")
	} else {
		log.Println("Router does not implement MuxRouter interface")
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	websocket.ServeWs(h.hub, w, r)
}

// HandleChannelWebSocket handles WebSocket connections for a specific channel
func (h *WebSocketHandler) HandleChannelWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract channel ID from URL path
	channelID := r.PathValue("channelId")

	// Add channel ID to query parameters to pass to the WebSocket handler
	q := r.URL.Query()
	q.Add("channelId", channelID)
	r.URL.RawQuery = q.Encode()

	// Serve WebSocket connection
	websocket.ServeWs(h.hub, w, r)
}

// MuxRouter interface for router compatibility
type MuxRouter interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// SendSystemMessage sends a system message to a channel
func (h *WebSocketHandler) SendSystemMessage(channelID uint, content string) {
	message := &websocket.Message{
		Type:      websocket.EventTypeMessage,
		ChannelID: channelID,
		Data: map[string]interface{}{
			"content": content,
			"system":  true,
		},
	}

	h.hub.Broadcast <- message
}

// SendUserOfflineNotification sends a notification when a user goes offline
func (h *WebSocketHandler) SendUserOfflineNotification(userID uint) {
	update := &websocket.PresenceUpdate{
		UserID: userID,
		Status: "offline",
	}

	h.hub.Presence <- update
}

// GetOnlineUsers returns the list of currently online users
func (h *WebSocketHandler) GetOnlineUsers() []uint {
	userIDs := make([]uint, 0)

	// Use the hub's state to get online users
	// This will need to be extended based on the hub implementation

	return userIDs
}
