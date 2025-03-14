package websocket

import "time"

// WebSocket event types
const (
	EventTypeMessage            = "message"
	EventTypeMessageEdit        = "message_edit"
	EventTypeMessageDelete      = "message_delete"
	EventTypeReaction           = "reaction"
	EventTypeTyping             = "typing"
	EventTypePresence           = "presence"
	EventTypeBulkPresence       = "bulk_presence"
	EventTypeChannelCreate      = "channel_create"
	EventTypeChannelUpdate      = "channel_update"
	EventTypeChannelDelete      = "channel_delete"
	EventTypeChannelSubscribe   = "channel_subscribe"
	EventTypeChannelUnsubscribe = "channel_unsubscribe"
	EventTypeError              = "error"
	EventTypeRead               = "read_receipt"
	EventTypeAttachment         = "attachment"
)

// User status constants
const (
	StatusOnline  = "online"
	StatusIdle    = "idle"
	StatusDND     = "dnd" // Do Not Disturb
	StatusOffline = "offline"
)

// ClientMessage represents a message received from a client
type ClientMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp,omitempty"`
}

// Message represents a message to be broadcast to clients
type Message struct {
	Type      string      `json:"type"`
	ChannelID uint        `json:"channelId"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// PresenceData contains a user's presence information
type PresenceData struct {
	Status    string    `json:"status"`
	LastSeen  time.Time `json:"-"`
	StatusMsg string    `json:"statusMsg,omitempty"`
}

// PresenceUpdate represents a change in a user's presence
type PresenceUpdate struct {
	Type      string `json:"type"`
	UserID    uint   `json:"userId"`
	Status    string `json:"status"`
	StatusMsg string `json:"statusMsg,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// BulkPresenceUpdate contains presence data for multiple users
type BulkPresenceUpdate struct {
	Type      string           `json:"type"`
	Presences []PresenceUpdate `json:"presences"`
}

// TypingEvent represents a typing indicator from a user
type TypingEvent struct {
	Type      string `json:"type"`
	UserID    uint   `json:"userId"`
	ChannelID uint   `json:"channelId"`
	IsTyping  bool   `json:"isTyping"`
	Timestamp int64  `json:"timestamp"`
}

// ReadReceipt represents a read receipt for a message
type ReadReceipt struct {
	Type            string `json:"type"`
	MessageID       uint   `json:"messageId"`
	ChannelID       uint   `json:"channelId"`
	UserID          uint   `json:"userId"`          // User who read the message
	MessageSenderID uint   `json:"messageSenderId"` // User who sent the original message
	Timestamp       int64  `json:"timestamp"`
}

// AttachmentEvent represents an attachment being uploaded or processed
type AttachmentEvent struct {
	Type         string `json:"type"`
	AttachmentID uint   `json:"attachmentId"`
	MessageID    uint   `json:"messageId"`
	ChannelID    uint   `json:"channelId"`
	Status       string `json:"status"`             // "uploading", "processing", "complete", "error"
	Progress     int    `json:"progress,omitempty"` // 0-100 for uploading/processing
	Error        string `json:"error,omitempty"`
	Timestamp    int64  `json:"timestamp"`
}

// ErrorMessage represents an error to be sent to clients
type ErrorMessage struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
