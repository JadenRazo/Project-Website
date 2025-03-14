package events

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging/websocket"
)

// WebSocketHandler broadcasts events to WebSocket clients
type WebSocketHandler struct {
	hub *websocket.Hub
}

// NewWebSocketHandler creates a new WebSocket event handler
func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleMessageCreated handles message created events
func (h *WebSocketHandler) HandleMessageCreated(ctx context.Context, event *Event) error {
	data, ok := event.Data.(*MessageEventData)
	if !ok {
		return ErrInvalidEventData
	}

	// Format data for WebSocket clients
	wsData := map[string]interface{}{
		"id":        data.Message.ID,
		"content":   data.Message.Content,
		"senderId":  data.Message.SenderID,
		"channelId": data.Message.ChannelID,
		"createdAt": data.Message.CreatedAt.Format(time.RFC3339),
	}

	// Add attachments if present
	if len(data.Message.Attachments) > 0 {
		attachments := make([]map[string]interface{}, len(data.Message.Attachments))
		for i, attachment := range data.Message.Attachments {
			attachments[i] = map[string]interface{}{
				"id":       attachment.ID,
				"fileName": attachment.FileName,
				"fileType": attachment.FileType,
				"fileSize": attachment.FileSize,
				"fileUrl":  attachment.FileURL,
				"isImage":  attachment.IsImage,
			}
		}
		wsData["attachments"] = attachments
	}

	// Broadcast to the channel
	h.hub.BroadcastToChannel(data.ChannelID, websocket.EventTypeMessage, wsData)
	return nil
}

// HandleMessageUpdated handles message updated events
func (h *WebSocketHandler) HandleMessageUpdated(ctx context.Context, event *Event) error {
	data, ok := event.Data.(*MessageEventData)
	if !ok {
		return ErrInvalidEventData
	}

	wsData := map[string]interface{}{
		"id":        data.Message.ID,
		"content":   data.Message.Content,
		"isEdited":  data.Message.IsEdited,
		"editedAt":  data.Message.EditedAt.Format(time.RFC3339),
		"channelId": data.Message.ChannelID,
	}

	h.hub.BroadcastToChannel(data.ChannelID, websocket.EventTypeMessageEdit, wsData)
	return nil
}

// HandleMessageDeleted handles message deleted events
func (h *WebSocketHandler) HandleMessageDeleted(ctx context.Context, event *Event) error {
	data, ok := event.Data.(*MessageEventData)
	if !ok {
		return ErrInvalidEventData
	}

	wsData := map[string]interface{}{
		"id":        data.Message.ID,
		"channelId": data.Message.ChannelID,
	}

	h.hub.BroadcastToChannel(data.ChannelID, websocket.EventTypeMessageDelete, wsData)
	return nil
}

// HandleReactionAdded handles reaction added events
func (h *WebSocketHandler) HandleReactionAdded(ctx context.Context, event *Event) error {
	data, ok := event.Data.(*ReactionEventData)
	if !ok {
		return ErrInvalidEventData
	}

	// Get channel ID from the message (simplified - in a real app, query the DB)
	// Assuming we have the channel ID from the message
	channelID := uint(0) // This should come from the message

	wsData := map[string]interface{}{
		"id":        data.Reaction.ID,
		"messageId": data.Reaction.MessageID,
		"userId":    data.Reaction.UserID,
		"emojiCode": data.Reaction.EmojiCode,
		"emojiName": data.Reaction.EmojiName,
	}

	h.hub.BroadcastToChannel(channelID, websocket.EventTypeReaction, wsData)
	return nil
}

// HandleChannelCreated handles channel created events
func (h *WebSocketHandler) HandleChannelCreated(ctx context.Context, event *Event) error {
	data, ok := event.Data.(*ChannelEventData)
	if !ok {
		return ErrInvalidEventData
	}

	wsData := map[string]interface{}{
		"id":          data.Channel.ID,
		"name":        data.Channel.Name,
		"description": data.Channel.Description,
		"type":        data.Channel.Type,
		"ownerId":     data.Channel.OwnerID,
		"icon":        data.Channel.Icon,
		"createdAt":   data.Channel.CreatedAt.Format(time.RFC3339),
	}

	// In a real app, broadcast to all users who should know about this channel
	// For simplicity, we'll just broadcast to the owner
	h.hub.BroadcastToUser(data.Channel.OwnerID, websocket.EventTypeChannelCreate, wsData)
	return nil
}

// RegisterHandlers registers all event handlers with the dispatcher
func RegisterHandlers(dispatcher *EventDispatcher, hub *websocket.Hub) {
	wsHandler := NewWebSocketHandler(hub)

	// Register message handlers
	dispatcher.Subscribe(EventMessageCreated, wsHandler.HandleMessageCreated)
	dispatcher.Subscribe(EventMessageUpdated, wsHandler.HandleMessageUpdated)
	dispatcher.Subscribe(EventMessageDeleted, wsHandler.HandleMessageDeleted)
	dispatcher.Subscribe(EventReactionAdded, wsHandler.HandleReactionAdded)

	// Register channel handlers
	dispatcher.Subscribe(EventChannelCreated, wsHandler.HandleChannelCreated)

	// Add more handlers as needed
}

// Common errors
var (
	ErrInvalidEventData = NewEventError("invalid_event_data", "The event data has an invalid format")
)

// EventError represents an error that occurred during event handling
type EventError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewEventError creates a new event error
func NewEventError(code, message string) *EventError {
	return &EventError{
		Code:    code,
		Message: message,
	}
}

func (e *EventError) Error() string {
	return e.Message
}
