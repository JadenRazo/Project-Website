package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the client
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from client
	maxMessageSize = 512 * 1024 // 512KB
)

// ServeWs handles WebSocket requests from clients
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Extract user ID from query parameter or authentication token
	// In a production environment, this would come from authentication middleware
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var userID uint
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get client IP for rate limiting
	ip := r.RemoteAddr
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ip = forwardedFor
	}

	// Check rate limiting
	if !hub.connManager.CanConnect(ip) {
		http.Error(w, "Connection rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Create client
	client := &Client{
		hub:      hub,
		conn:     NewWebSocketConnection(conn),
		send:     make(chan []byte, 256),
		UserID:   userID,
		Channels: []uint{},
		IP:       ip,
		Metadata: make(map[string]interface{}),
	}

	// Register client with hub
	client.hub.register <- client

	// Update connection manager
	hub.connManager.AddConnection(ip)

	// Start client read/write goroutines
	go client.readPump()
	go client.writePump()
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Set read parameters
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Continuously read messages
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading from WebSocket: %v", err)
			}
			break
		}

		// Parse message
		var clientMsg ClientMessage
		if err := json.Unmarshal(message, &clientMsg); err != nil {
			log.Printf("Error parsing client message: %v", err)
			c.sendError("invalid_message", "Invalid message format")
			continue
		}

		// Process message based on type
		c.processMessage(&clientMsg)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Write message to WebSocket
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// Send ping
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// processMessage processes incoming client messages based on their type
func (c *Client) processMessage(msg *ClientMessage) {
	switch msg.Type {
	case EventTypeChannelSubscribe:
		c.handleChannelSubscribe(msg)
	case EventTypeChannelUnsubscribe:
		c.handleChannelUnsubscribe(msg)
	case EventTypeMessage:
		c.handleNewMessage(msg)
	case EventTypePresenceUpdate:
		c.handlePresenceUpdate(msg)
	case EventTypeTyping:
		c.handleTyping(msg)
	case EventTypeReadReceipt:
		c.handleReadReceipt(msg)
	default:
		c.sendError("unknown_message_type", "Unknown message type")
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(code string, message string) {
	errorResponse := ErrorResponse{
		Code:    code,
		Message: message,
	}

	response := map[string]interface{}{
		"type": EventTypeError,
		"data": errorResponse,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling error response: %v", err)
		return
	}

	c.send <- payload
}

// sendAck sends an acknowledgment message to the client
func (c *Client) sendAck(eventType string, data interface{}) {
	response := map[string]interface{}{
		"type": EventTypeAck,
		"data": map[string]interface{}{
			"type": eventType,
			"data": data,
		},
	}

	payload, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling ack response: %v", err)
		return
	}

	c.send <- payload
}

// extractData extracts data from an interface{} into a struct
func extractData(data interface{}, target interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

// handleChannelSubscribe processes channel subscription requests
func (c *Client) handleChannelSubscribe(msg *ClientMessage) {
	var data struct {
		ChannelID uint `json:"channelId"`
	}

	if err := extractData(msg.Data, &data); err != nil {
		c.sendError("invalid_data", "Invalid channel subscription data")
		return
	}

	// Add client to channel
	c.hub.SubscribeToChannel(c, data.ChannelID)

	// Acknowledge subscription
	c.sendAck(EventTypeChannelSubscribe, map[string]interface{}{
		"channelId": data.ChannelID,
		"success":   true,
	})
}

// handleChannelUnsubscribe processes channel unsubscription requests
func (c *Client) handleChannelUnsubscribe(msg *ClientMessage) {
	var data struct {
		ChannelID uint `json:"channelId"`
	}

	if err := extractData(msg.Data, &data); err != nil {
		c.sendError("invalid_data", "Invalid channel unsubscription data")
		return
	}

	// Remove client from channel
	c.hub.UnsubscribeFromChannel(c, data.ChannelID)

	// Acknowledge unsubscription
	c.sendAck(EventTypeChannelUnsubscribe, map[string]interface{}{
		"channelId": data.ChannelID,
		"success":   true,
	})
}

// handleNewMessage processes new message requests from clients
func (c *Client) handleNewMessage(msg *ClientMessage) {
	var data struct {
		Content   string `json:"content"`
		ChannelID uint   `json:"channelId"`
	}

	if err := extractData(msg.Data, &data); err != nil {
		c.sendError("invalid_data", "Invalid message data")
		return
	}

	// Forward to hub as a message event
	message := &Message{
		Type:      EventTypeMessage,
		ChannelID: data.ChannelID,
		Data: map[string]interface{}{
			"content":  data.Content,
			"senderId": c.UserID,
		},
		Timestamp: time.Now().Unix(),
	}

	c.hub.broadcast <- message
}

// handlePresenceUpdate processes presence update requests
func (c *Client) handlePresenceUpdate(msg *ClientMessage) {
	var data struct {
		Status string `json:"status"`
	}

	if err := extractData(msg.Data, &data); err != nil {
		c.sendError("invalid_data", "Invalid presence data")
		return
	}

	// Forward to hub as a presence update
	update := &PresenceUpdate{
		UserID:    c.UserID,
		Status:    data.Status,
		Timestamp: time.Now().Unix(),
	}

	c.hub.presence <- update
}

// handleTyping processes typing indicator requests
func (c *Client) handleTyping(msg *ClientMessage) {
	var data struct {
		ChannelID uint `json:"channelId"`
		IsTyping  bool `json:"isTyping"`
	}

	if err := extractData(msg.Data, &data); err != nil {
		c.sendError("invalid_data", "Invalid typing data")
		return
	}

	// Forward to hub as a typing event
	typing := &TypingEvent{
		ChannelID: data.ChannelID,
		UserID:    c.UserID,
		IsTyping:  data.IsTyping,
		Timestamp: time.Now().Unix(),
	}

	c.hub.typing <- typing
}

// handleReadReceipt processes read receipt requests
func (c *Client) handleReadReceipt(msg *ClientMessage) {
	var data struct {
		MessageID uint `json:"messageId"`
		ChannelID uint `json:"channelId"`
	}

	if err := extractData(msg.Data, &data); err != nil {
		c.sendError("invalid_data", "Invalid read receipt data")
		return
	}

	// Forward to hub as a read receipt
	receipt := &ReadReceipt{
		MessageID: data.MessageID,
		ChannelID: data.ChannelID,
		UserID:    c.UserID,
		Timestamp: time.Now().Unix(),
	}

	c.hub.readReceipts <- receipt
}
