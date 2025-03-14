package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connected websocket client
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// User information
	UserID uint

	// Channels the client is subscribed to
	Channels []uint

	// Hub for broadcasting messages
	hub *Hub

	// Buffered channel of outbound messages
	send chan []byte

	// Time of last activity for timeout management
	lastActivity time.Time
}

// NewClient creates a new websocket client
func NewClient(conn *websocket.Conn, userID uint, hub *Hub) *Client {
	return &Client{
		conn:         conn,
		UserID:       userID,
		hub:          hub,
		send:         make(chan []byte, 256), // Buffer size for outgoing messages
		lastActivity: time.Now(),
		Channels:     make([]uint, 0),
	}
}

// ReadPump handles messages from the websocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Set read limits to prevent memory issues
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			}
			break
		}

		// Update last activity timestamp
		c.lastActivity = time.Now()

		// Parse the incoming message
		var clientMsg ClientMessage
		if err := json.Unmarshal(message, &clientMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			c.sendError("invalid_format", "Invalid message format")
			continue
		}

		// Handle different message types
		c.handleMessage(&clientMsg)
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
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

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming client messages based on their type
func (c *Client) handleMessage(msg *ClientMessage) {
	switch msg.Type {
	case EventTypeTyping:
		c.handleTypingIndicator(msg)
	case EventTypePresence:
		c.handlePresenceUpdate(msg)
	case EventTypeChannelSubscribe:
		c.handleChannelSubscribe(msg)
	case EventTypeChannelUnsubscribe:
		c.handleChannelUnsubscribe(msg)
	case EventTypeMessage:
		c.handleNewMessage(msg)
	default:
		c.sendError("unknown_message_type", "Unknown message type")
	}
}

// handleTypingIndicator processes typing events
func (c *Client) handleTypingIndicator(msg *ClientMessage) {
	var data struct {
		ChannelID uint `json:"channelId"`
		IsTyping  bool `json:"isTyping"`
	}

	if err := extractData(msg.Data, &data); err != nil {
		c.sendError("invalid_data", "Invalid typing data format")
		return
	}

	// Create typing event
	typingEvent := Message{
		Type:      EventTypeTyping,
		ChannelID: data.ChannelID,
		Data: map[string]interface{}{
			"userId":   c.UserID,
			"isTyping": data.IsTyping,
		},
		Timestamp: time.Now().Unix(),
	}

	// Broadcast to channel
	c.hub.broadcast <- &typingEvent
}

// handlePresenceUpdate processes presence status updates
func (c *Client) handlePresenceUpdate(msg *ClientMessage) {
	var data struct {
		Status    string `json:"status"`
		StatusMsg string `json:"statusMsg,omitempty"`
	}

	if err := extractData(msg.Data, &data); err != nil {
		c.sendError("invalid_data", "Invalid presence data format")
		return
	}

	// Validate status
	validStatus := false
	for _, status := range []string{StatusOnline, StatusIdle, StatusDND, StatusOffline} {
		if data.Status == status {
			validStatus = true
			break
		}
	}

	if !validStatus {
		c.sendError("invalid_status", "Invalid status value")
		return
	}

	// Create presence update
	presenceUpdate := &PresenceUpdate{
		Type:      EventTypePresence,
		UserID:    c.UserID,
		Status:    data.Status,
		StatusMsg: data.StatusMsg,
		Timestamp: time.Now().Unix(),
	}

	c.hub.presence <- presenceUpdate
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

// sendError sends an error message to the client
func (c *Client) sendError(code, message string) {
	errorMsg := ErrorMessage{
		Type:    EventTypeError,
		Code:    code,
		Message: message,
	}

	data, err := json.Marshal(errorMsg)
	if err != nil {
		log.Printf("Error marshaling error message: %v", err)
		return
	}

	c.send <- data
}

// sendAck sends an acknowledgment message
func (c *Client) sendAck(eventType string, data interface{}) {
	ack := map[string]interface{}{
		"type":      eventType + "_ack",
		"data":      data,
		"timestamp": time.Now().Unix(),
	}

	jsonData, err := json.Marshal(ack)
	if err != nil {
		log.Printf("Error marshaling ack message: %v", err)
		return
	}

	c.send <- jsonData
}

// Close closes the client connection and cleans up resources
func (c *Client) Close() {
	close(c.send)
}

// Helper function to extract data from interface{} to struct
func extractData(raw interface{}, target interface{}) error {
	data, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// WebSocket constants
const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512KB
)

var (
	newline = []byte{'\n'}
)
