// backend/internal/messaging/websocket/hub.go
package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// ClientList is a map of connected clients with thread-safe access
type ClientList struct {
	sync.RWMutex
	clients map[*Client]bool
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients by user ID
	userClients map[uint]*ClientList

	// Channel mapping for quicker lookups
	channelClients map[uint]*ClientList

	// Registration/unregistration channels
	register   chan *Client
	unregister chan *Client

	// Broadcast message channel
	broadcast chan *Message

	// Presence updates for status changes
	presence chan *PresenceUpdate

	// Typing indicator updates
	typing chan *TypingEvent

	// Read receipts
	readReceipts chan *ReadReceipt

	// Mutex for thread-safe access to maps
	mu sync.RWMutex

	// User presence map (userID -> presence data)
	userPresence map[uint]*PresenceData

	// Connection manager for rate limiting and connection count
	connManager *ConnectionManager
}

// NewHub creates a new hub instance
func NewHub(connManager *ConnectionManager) *Hub {
	return &Hub{
		userClients:    make(map[uint]*ClientList),
		channelClients: make(map[uint]*ClientList),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan *Message),
		presence:       make(chan *PresenceUpdate),
		typing:         make(chan *TypingEvent),
		readReceipts:   make(chan *ReadReceipt),
		userPresence:   make(map[uint]*PresenceData),
		connManager:    connManager,
	}
}

// Run starts the hub's main loop for handling client events
func (h *Hub) Run() {
	// Periodic cleanup for inactive clients
	ticker := time.NewTicker(2 * time.Minute)

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)

		case presence := <-h.presence:
			h.updatePresence(presence)

		case typing := <-h.typing:
			h.broadcastTypingEvent(typing)

		case receipt := <-h.readReceipts:
			h.handleReadReceipt(receipt)

		case <-ticker.C:
			h.cleanupInactiveClients()
		}
	}
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
	h.connManager.AddConnection()
}

// registerClient adds a new client connection
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Add to user client list
	if _, exists := h.userClients[client.UserID]; !exists {
		h.userClients[client.UserID] = &ClientList{
			clients: make(map[*Client]bool),
		}
	}

	h.userClients[client.UserID].Lock()
	h.userClients[client.UserID].clients[client] = true
	h.userClients[client.UserID].Unlock()

	// Add to channel client lists
	for _, channelID := range client.Channels {
		if _, exists := h.channelClients[channelID]; !exists {
			h.channelClients[channelID] = &ClientList{
				clients: make(map[*Client]bool),
			}
		}

		h.channelClients[channelID].Lock()
		h.channelClients[channelID].clients[client] = true
		h.channelClients[channelID].Unlock()
	}

	// Update user presence to online
	h.updateUserPresence(client.UserID, StatusOnline)

	// Send initial presence data to the client
	h.sendPresenceData(client)
}

// unregisterClient removes a client connection
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Decrement connection count
	h.connManager.RemoveConnection()

	// Remove from user client list
	if clients, exists := h.userClients[client.UserID]; exists {
		clients.Lock()
		delete(clients.clients, client)
		numClients := len(clients.clients)
		clients.Unlock()

		// If no more clients for this user, clean up and update presence
		if numClients == 0 {
			delete(h.userClients, client.UserID)
			h.updateUserPresence(client.UserID, StatusOffline)
		}
	}

	// Remove from channel client lists
	for _, channelID := range client.Channels {
		if clients, exists := h.channelClients[channelID]; exists {
			clients.Lock()
			delete(clients.clients, client)
			numClients := len(clients.clients)
			clients.Unlock()

			// If no more clients for this channel, clean up
			if numClients == 0 {
				delete(h.channelClients, channelID)
			}
		}
	}

	// Close client connection
	client.Close()
}

// broadcastMessage sends a message to all relevant clients
func (h *Hub) broadcastMessage(msg *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Get all clients in the target channel
	channelClients, ok := h.channelClients[msg.ChannelID]
	if !ok {
		return
	}

	// Convert message to JSON
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	// Send to all clients in the channel
	channelClients.RLock()
	for client := range channelClients.clients {
		select {
		case client.send <- data:
			// Message sent successfully
		default:
			// Buffer full, client might be slow or disconnected
			h.unregister <- client
		}
	}
	channelClients.RUnlock()
}

// broadcastTypingEvent broadcasts typing event to all clients in a channel
func (h *Hub) broadcastTypingEvent(event *TypingEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Get all clients in the target channel
	channelClients, ok := h.channelClients[event.ChannelID]
	if !ok {
		return
	}

	// Convert event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal typing event: %v", err)
		return
	}

	// Send to all clients in the channel
	channelClients.RLock()
	for client := range channelClients.clients {
		// Don't send typing events back to the user who is typing
		if client.UserID == event.UserID {
			continue
		}

		select {
		case client.send <- data:
			// Event sent successfully
		default:
			// Buffer full, client might be slow or disconnected
			h.unregister <- client
		}
	}
	channelClients.RUnlock()
}

// handleReadReceipt updates read status and notifies relevant clients
func (h *Hub) handleReadReceipt(receipt *ReadReceipt) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Only notify the sender of the original message
	userClients, ok := h.userClients[receipt.MessageSenderID]
	if !ok {
		return
	}

	// Convert receipt to JSON
	data, err := json.Marshal(receipt)
	if err != nil {
		log.Printf("Failed to marshal read receipt: %v", err)
		return
	}

	// Send to the original message sender's clients
	userClients.RLock()
	for client := range userClients.clients {
		select {
		case client.send <- data:
			// Receipt sent successfully
		default:
			// Buffer full
			h.unregister <- client
		}
	}
	userClients.RUnlock()
}

// updatePresence broadcasts a presence update to relevant clients
func (h *Hub) updatePresence(update *PresenceUpdate) {
	h.mu.Lock()

	// Update presence in our map
	h.userPresence[update.UserID] = &PresenceData{
		Status:    update.Status,
		LastSeen:  time.Now(),
		StatusMsg: update.StatusMsg,
	}

	h.mu.Unlock()

	// Broadcast the update to all relevant clients
	data, err := json.Marshal(update)
	if err != nil {
		log.Printf("Failed to marshal presence update: %v", err)
		return
	}

	// For each channel this user is in, notify all clients
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Notify all connected clients
	for _, clients := range h.channelClients {
		clients.RLock()
		for client := range clients.clients {
			select {
			case client.send <- data:
				// Update sent
			default:
				// Buffer full
				go func(c *Client) {
					h.unregister <- c
				}(client)
			}
		}
		clients.RUnlock()
	}
}

// updateUserPresence updates a user's presence status
func (h *Hub) updateUserPresence(userID uint, status string) {
	update := &PresenceUpdate{
		Type:      EventTypePresence,
		UserID:    userID,
		Status:    status,
		Timestamp: time.Now().Unix(),
	}

	h.presence <- update
}

// sendPresenceData sends current presence data to a newly connected client
func (h *Hub) sendPresenceData(client *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	presenceList := make([]PresenceUpdate, 0, len(h.userPresence))

	// Create a bulk presence update
	for userID, data := range h.userPresence {
		presenceList = append(presenceList, PresenceUpdate{
			Type:      EventTypePresence,
			UserID:    userID,
			Status:    data.Status,
			StatusMsg: data.StatusMsg,
			Timestamp: time.Now().Unix(),
		})
	}

	// Send bulk presence update
	bulkUpdate := BulkPresenceUpdate{
		Type:      EventTypeBulkPresence,
		Presences: presenceList,
	}

	data, err := json.Marshal(bulkUpdate)
	if err != nil {
		log.Printf("Failed to marshal bulk presence: %v", err)
		return
	}

	client.send <- data
}

// cleanupInactiveClients removes clients that haven't been active
func (h *Hub) cleanupInactiveClients() {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, clientList := range h.userClients {
		clientList.RLock()
		for client := range clientList.clients {
			if time.Since(client.lastActivity) > 10*time.Minute {
				h.unregister <- client
			}
		}
		clientList.RUnlock()
	}

	// Clean up stale presence data
	timeout := 24 * time.Hour
	now := time.Now()

	for userID, presence := range h.userPresence {
		if presence.Status == StatusOffline && now.Sub(presence.LastSeen) > timeout {
			delete(h.userPresence, userID)
		}
	}
}

// SubscribeToChannel adds a client to a channel's broadcast list
func (h *Hub) SubscribeToChannel(client *Client, channelID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Add channel to client's subscription list if not already there
	found := false
	for _, id := range client.Channels {
		if id == channelID {
			found = true
			break
		}
	}

	if !found {
		client.Channels = append(client.Channels, channelID)
	}

	// Add client to channel distribution list
	if _, exists := h.channelClients[channelID]; !exists {
		h.channelClients[channelID] = &ClientList{
			clients: make(map[*Client]bool),
		}
	}

	h.channelClients[channelID].Lock()
	h.channelClients[channelID].clients[client] = true
	h.channelClients[channelID].Unlock()
}

// UnsubscribeFromChannel removes a client from a channel's broadcast list
func (h *Hub) UnsubscribeFromChannel(client *Client, channelID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Remove channel from client's subscription list
	for i, id := range client.Channels {
		if id == channelID {
			client.Channels = append(client.Channels[:i], client.Channels[i+1:]...)
			break
		}
	}

	// Remove client from channel distribution list
	if clients, exists := h.channelClients[channelID]; exists {
		clients.Lock()
		delete(clients.clients, client)
		numClients := len(clients.clients)
		clients.Unlock()

		// If no more clients for this channel, clean up
		if numClients == 0 {
			delete(h.channelClients, channelID)
		}
	}
}

// BroadcastToChannel sends a message to a specific channel
func (h *Hub) BroadcastToChannel(channelID uint, msgType string, data interface{}) {
	message := &Message{
		Type:      msgType,
		ChannelID: channelID,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	h.broadcast <- message
}

// BroadcastTyping sends a typing indicator to a channel
func (h *Hub) BroadcastTyping(userID, channelID uint, isTyping bool) {
	event := &TypingEvent{
		Type:      EventTypeTyping,
		UserID:    userID,
		ChannelID: channelID,
		IsTyping:  isTyping,
		Timestamp: time.Now().Unix(),
	}

	h.typing <- event
}

// SendReadReceipt notifies that a user has read a message
func (h *Hub) SendReadReceipt(messageID, channelID, userID, senderID uint) {
	receipt := &ReadReceipt{
		Type:            EventTypeRead,
		MessageID:       messageID,
		ChannelID:       channelID,
		UserID:          userID,   // Who read the message
		MessageSenderID: senderID, // Who sent the original message
		Timestamp:       time.Now().Unix(),
	}

	h.readReceipts <- receipt
}

// BroadcastToUser sends a message to all clients of a specific user
func (h *Hub) BroadcastToUser(userID uint, msgType string, data interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.userClients[userID]
	if !ok {
		return
	}

	message := map[string]interface{}{
		"type":      msgType,
		"data":      data,
		"timestamp": time.Now().Unix(),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal user message: %v", err)
		return
	}

	clients.RLock()
	defer clients.RUnlock()

	for client := range clients.clients {
		select {
		case client.send <- jsonData:
			// Message sent
		default:
			// Buffer full
			go func(c *Client) {
				h.unregister <- c
			}(client)
		}
	}
}
