package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Message represents a message to be broadcast to clients
type Message struct {
	Type      string                 `json:"type"`
	ChannelID uint                   `json:"channelId"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

// PresenceUpdate represents a user presence status change
type PresenceUpdate struct {
	UserID    uint   `json:"userId"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

// TypingEvent represents a typing indicator event
type TypingEvent struct {
	ChannelID uint  `json:"channelId"`
	UserID    uint  `json:"userId"`
	IsTyping  bool  `json:"isTyping"`
	Timestamp int64 `json:"timestamp"`
}

// ReadReceipt represents a message read receipt
type ReadReceipt struct {
	MessageID uint  `json:"messageId"`
	ChannelID uint  `json:"channelId"`
	UserID    uint  `json:"userId"`
	Timestamp int64 `json:"timestamp"`
}

// PresenceData stores user presence information
type PresenceData struct {
	Status   string    `json:"status"`
	LastSeen time.Time `json:"lastSeen"`
	IsOnline bool      `json:"isOnline"`
}

// Client represents a connected client
type Client struct {
	hub      *Hub
	conn     Connection
	send     chan []byte
	UserID   uint
	Channels []uint
	IP       string
	Metadata map[string]interface{}
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients ClientList

	// Channels mapped to clients
	channels map[uint]*ClientList

	// Inbound messages from the clients
	broadcast chan *Message

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Presence updates from clients
	presence chan *PresenceUpdate

	// Typing indicators from clients
	typing chan *TypingEvent

	// Read receipts from clients
	readReceipts chan *ReadReceipt

	// Connection manager for rate limiting
	connManager *ConnectionManager

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// User presence information
	userPresence map[uint]*PresenceData
}

// NewHub creates a new hub instance
func NewHub() *Hub {
	return &Hub{
		clients:      ClientList{clients: make(map[*Client]bool)},
		channels:     make(map[uint]*ClientList),
		broadcast:    make(chan *Message),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		presence:     make(chan *PresenceUpdate),
		typing:       make(chan *TypingEvent),
		readReceipts: make(chan *ReadReceipt),
		connManager:  NewConnectionManager(),
		userPresence: make(map[uint]*PresenceData),
	}
}

// Run starts the hub, processing messages and client connections
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		case update := <-h.presence:
			h.broadcastPresence(update)
		case typing := <-h.typing:
			h.broadcastTyping(typing)
		case receipt := <-h.readReceipts:
			h.broadcastReadReceipt(receipt)
		}
	}
}

// registerClient adds a client to the hub
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients.Lock()
	h.clients.clients[client] = true
	h.clients.Unlock()

	// Update presence information
	h.userPresence[client.UserID] = &PresenceData{
		Status:   "online",
		LastSeen: time.Now(),
		IsOnline: true,
	}

	log.Printf("Client registered: User ID %d from IP %s", client.UserID, client.IP)
}

// unregisterClient removes a client from the hub and all channels
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients.Lock()
	if _, ok := h.clients.clients[client]; ok {
		delete(h.clients.clients, client)
		close(client.send)
	}
	h.clients.Unlock()

	// Remove from all channels
	for _, channelID := range client.Channels {
		if channelClients, ok := h.channels[channelID]; ok {
			channelClients.Lock()
			delete(channelClients.clients, client)
			channelClients.Unlock()
		}
	}

	// Update connection manager
	h.connManager.RemoveConnection(client.IP)

	// Update presence information
	if presence, ok := h.userPresence[client.UserID]; ok {
		presence.Status = "offline"
		presence.LastSeen = time.Now()
		presence.IsOnline = false
	}

	log.Printf("Client unregistered: User ID %d from IP %s", client.UserID, client.IP)
}

// SubscribeToChannel adds a client to a channel
func (h *Hub) SubscribeToChannel(client *Client, channelID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Initialize channel if it doesn't exist
	if _, ok := h.channels[channelID]; !ok {
		h.channels[channelID] = &ClientList{
			clients: make(map[*Client]bool),
		}
	}

	// Add client to channel
	h.channels[channelID].Lock()
	h.channels[channelID].clients[client] = true
	h.channels[channelID].Unlock()

	// Add channel to client's subscriptions if not already there
	channelExists := false
	for _, ch := range client.Channels {
		if ch == channelID {
			channelExists = true
			break
		}
	}

	if !channelExists {
		client.Channels = append(client.Channels, channelID)
	}

	log.Printf("Client subscribed: User ID %d to channel %d", client.UserID, channelID)
}

// UnsubscribeFromChannel removes a client from a channel
func (h *Hub) UnsubscribeFromChannel(client *Client, channelID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Remove client from channel
	if channelClients, ok := h.channels[channelID]; ok {
		channelClients.Lock()
		delete(channelClients.clients, client)
		channelClients.Unlock()
	}

	// Remove channel from client's subscriptions
	for i, ch := range client.Channels {
		if ch == channelID {
			client.Channels = append(client.Channels[:i], client.Channels[i+1:]...)
			break
		}
	}

	log.Printf("Client unsubscribed: User ID %d from channel %d", client.UserID, channelID)
}

// broadcastMessage sends a message to all clients in the specified channel
func (h *Hub) broadcastMessage(msg *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Marshal message to JSON
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	// Find clients in the target channel
	if clients, ok := h.channels[msg.ChannelID]; ok {
		clients.RLock()
		for client := range clients.clients {
			select {
			case client.send <- payload:
				// Message sent
			default:
				// Failed to send, client might be slow or disconnected
				go h.unregisterClient(client)
			}
		}
		clients.RUnlock()
	}

	log.Printf("Message broadcast to channel %d: %s", msg.ChannelID, msg.Type)
}

// broadcastPresence sends a presence update to all connected clients
func (h *Hub) broadcastPresence(update *PresenceUpdate) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Create message
	msg := map[string]interface{}{
		"type": EventTypePresenceUpdate,
		"data": update,
	}

	// Marshal message to JSON
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling presence update: %v", err)
		return
	}

	// Update presence information
	h.userPresence[update.UserID] = &PresenceData{
		Status:   update.Status,
		LastSeen: time.Now(),
		IsOnline: update.Status != "offline",
	}

	// Send to all clients
	h.clients.RLock()
	for client := range h.clients.clients {
		select {
		case client.send <- payload:
			// Message sent
		default:
			// Failed to send
			go h.unregisterClient(client)
		}
	}
	h.clients.RUnlock()

	log.Printf("Presence update broadcast: User ID %d status '%s'", update.UserID, update.Status)
}

// broadcastTyping sends a typing indicator to all clients in the specified channel
func (h *Hub) broadcastTyping(typing *TypingEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Create message
	msg := map[string]interface{}{
		"type": EventTypeTyping,
		"data": typing,
	}

	// Marshal message to JSON
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling typing event: %v", err)
		return
	}

	// Find clients in the target channel
	if clients, ok := h.channels[typing.ChannelID]; ok {
		clients.RLock()
		for client := range clients.clients {
			// Don't send typing indicators back to the sender
			if client.UserID == typing.UserID {
				continue
			}

			select {
			case client.send <- payload:
				// Message sent
			default:
				// Failed to send
				go h.unregisterClient(client)
			}
		}
		clients.RUnlock()
	}

	log.Printf("Typing event broadcast to channel %d: User ID %d is typing: %v",
		typing.ChannelID, typing.UserID, typing.IsTyping)
}

// broadcastReadReceipt sends a read receipt to all clients in the specified channel
func (h *Hub) broadcastReadReceipt(receipt *ReadReceipt) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Create message
	msg := map[string]interface{}{
		"type": EventTypeReadReceipt,
		"data": receipt,
	}

	// Marshal message to JSON
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling read receipt: %v", err)
		return
	}

	// Find clients in the target channel
	if clients, ok := h.channels[receipt.ChannelID]; ok {
		clients.RLock()
		for client := range clients.clients {
			select {
			case client.send <- payload:
				// Message sent
			default:
				// Failed to send
				go h.unregisterClient(client)
			}
		}
		clients.RUnlock()
	}

	log.Printf("Read receipt broadcast to channel %d: User ID %d read message %d",
		receipt.ChannelID, receipt.UserID, receipt.MessageID)
}
