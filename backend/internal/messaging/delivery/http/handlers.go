package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/service"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/websocket"
	"github.com/gorilla/mux"
)

// MessagingHandler handles HTTP requests for messaging functionality
type MessagingHandler struct {
	messagingService service.MessagingService
	hub              *websocket.Hub
}

// NewMessagingHandler creates a new messaging handler
func NewMessagingHandler(messagingService service.MessagingService, hub *websocket.Hub) *MessagingHandler {
	return &MessagingHandler{
		messagingService: messagingService,
		hub:              hub,
	}
}

// RegisterRoutes registers the handler routes on the given router
func (h *MessagingHandler) RegisterRoutes(router *mux.Router) {
	// Channel routes
	router.HandleFunc("/channels", h.CreateChannel).Methods("POST")
	router.HandleFunc("/channels/{channelID}", h.GetChannel).Methods("GET")
	router.HandleFunc("/channels/{channelID}", h.UpdateChannel).Methods("PUT")
	router.HandleFunc("/channels/{channelID}", h.DeleteChannel).Methods("DELETE")
	router.HandleFunc("/channels/{channelID}/messages", h.GetMessages).Methods("GET")
	router.HandleFunc("/channels/{channelID}/messages", h.SendMessage).Methods("POST")
	router.HandleFunc("/channels/{channelID}/members", h.AddChannelMember).Methods("POST")
	router.HandleFunc("/channels/{channelID}/members/{userID}", h.RemoveChannelMember).Methods("DELETE")

	// User routes
	router.HandleFunc("/users/{userID}/channels", h.GetUserChannels).Methods("GET")

	// Message routes
	router.HandleFunc("/messages/{messageID}", h.GetMessage).Methods("GET")
	router.HandleFunc("/messages/{messageID}", h.UpdateMessage).Methods("PUT")
	router.HandleFunc("/messages/{messageID}", h.DeleteMessage).Methods("DELETE")
	router.HandleFunc("/messages/{messageID}/reactions", h.AddReaction).Methods("POST")
	router.HandleFunc("/messages/{messageID}/reactions/{reactionID}", h.RemoveReaction).Methods("DELETE")

	// WebSocket endpoint
	router.HandleFunc("/ws", h.HandleWebSocket).Methods("GET")
}

// CreateChannel handles channel creation
func (h *MessagingHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsPrivate   bool   `json:"isPrivate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create channel
	channel, err := h.messagingService.CreateChannel(r.Context(), req.Name, req.Description, req.IsPrivate, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	respond(w, http.StatusCreated, channel)
}

// GetChannel handles retrieving a channel by ID
func (h *MessagingHandler) GetChannel(w http.ResponseWriter, r *http.Request) {
	// Parse channel ID from URL
	channelID, err := parseChannelID(r)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get channel from service
	channel, err := h.messagingService.GetChannel(r.Context(), channelID)
	if err != nil {
		if err == repository.ErrChannelNotFound {
			http.Error(w, "Channel not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	respond(w, http.StatusOK, channel)
}

// UpdateChannel handles updating a channel
func (h *MessagingHandler) UpdateChannel(w http.ResponseWriter, r *http.Request) {
	// Parse channel ID from URL
	channelID, err := parseChannelID(r)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	// Parse request
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Topic       string `json:"topic"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get existing channel
	channel, err := h.messagingService.GetChannel(r.Context(), channelID)
	if err != nil {
		if err == repository.ErrChannelNotFound {
			http.Error(w, "Channel not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is the owner
	if channel.OwnerID != userID {
		http.Error(w, "Only the channel owner can update a channel", http.StatusForbidden)
		return
	}

	// Update channel fields
	channel.Name = req.Name
	channel.Description = req.Description
	channel.Icon = req.Icon
	channel.Topic = req.Topic
	channel.UpdatedAt = time.Now()

	// Update in service
	if err := h.messagingService.UpdateChannel(r.Context(), channel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	respond(w, http.StatusOK, channel)
}

// DeleteChannel handles deleting a channel
func (h *MessagingHandler) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	// Parse channel ID from URL
	channelID, err := parseChannelID(r)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get channel to check ownership
	channel, err := h.messagingService.GetChannel(r.Context(), channelID)
	if err != nil {
		if err == repository.ErrChannelNotFound {
			http.Error(w, "Channel not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify user is the owner
	if channel.OwnerID != userID {
		http.Error(w, "Only the channel owner can delete a channel", http.StatusForbidden)
		return
	}

	// Delete channel
	if err := h.messagingService.DeleteChannel(r.Context(), channelID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}

// GetMessages handles retrieving messages from a channel
func (h *MessagingHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	// Parse channel ID from URL
	channelID, err := parseChannelID(r)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	// Parse pagination params
	lastMessageID := uint(0)
	if lastIDStr := r.URL.Query().Get("lastMessageId"); lastIDStr != "" {
		lastID, err := strconv.ParseUint(lastIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid lastMessageId parameter", http.StatusBadRequest)
			return
		}
		lastMessageID = uint(lastID)
	}

	limit := 50 // Default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		if parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	// Get user ID from context
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get messages
	messages, err := h.messagingService.GetChannelMessages(r.Context(), channelID, lastMessageID, limit)
	if err != nil {
		if err == repository.ErrChannelNotFound {
			http.Error(w, "Channel not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	respond(w, http.StatusOK, messages)
}

// SendMessage handles sending a new message to a channel
func (h *MessagingHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	// Parse channel ID from URL
	channelID, err := parseChannelID(r)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	// Parse request
	var req struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Send message
	message, err := h.messagingService.SendMessage(r.Context(), req.Content, userID, channelID)
	if err != nil {
		if err == repository.ErrChannelNotFound || err == repository.ErrUserNotMember {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	respond(w, http.StatusCreated, message)
}

// GetUserChannels handles retrieving a user's channels
func (h *MessagingHandler) GetUserChannels(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from URL
	params := mux.Vars(r)
	userIDStr := params["userID"]
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userID := uint(userIDUint)

	// Get current user from context
	currentUserID := getUserIDFromContext(r.Context())
	if currentUserID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if requesting own channels or admin
	if currentUserID != userID {
		// Here you could check for admin status
		http.Error(w, "Forbidden: can only access your own channels", http.StatusForbidden)
		return
	}

	// Get channels
	channels, err := h.messagingService.ListUserChannels(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	respond(w, http.StatusOK, channels)
}

// AddChannelMember handles adding a user to a channel
func (h *MessagingHandler) AddChannelMember(w http.ResponseWriter, r *http.Request) {
	// Parse channel ID from URL
	channelID, err := parseChannelID(r)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	// Parse request
	var req struct {
		UserID uint `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get current user from context
	currentUserID := getUserIDFromContext(r.Context())
	if currentUserID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Add user to channel
	if err := h.messagingService.AddUserToChannel(r.Context(), req.UserID, channelID); err != nil {
		if err == repository.ErrChannelNotFound {
			http.Error(w, "Channel not found", http.StatusNotFound)
			return
		}
		if err == repository.ErrUserAlreadyMember {
			http.Error(w, "User is already a member of this channel", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}

// RemoveChannelMember handles removing a user from a channel
func (h *MessagingHandler) RemoveChannelMember(w http.ResponseWriter, r *http.Request) {
	// Parse channel ID and user ID from URL
	params := mux.Vars(r)
	channelIDStr := params["channelID"]
	channelIDUint, err := strconv.ParseUint(channelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}
	channelID := uint(channelIDUint)

	userIDStr := params["userID"]
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userID := uint(userIDUint)

	// Get current user from context
	currentUserID := getUserIDFromContext(r.Context())
	if currentUserID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if current user is removing themselves or if they're the channel owner
	channel, err := h.messagingService.GetChannel(r.Context(), channelID)
	if err != nil {
		if err == repository.ErrChannelNotFound {
			http.Error(w, "Channel not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Only allow removal if:
	// 1. User is removing themselves, or
	// 2. User is the channel owner
	if currentUserID != userID && channel.OwnerID != currentUserID {
		http.Error(w, "Forbidden: only channel owner can remove other users", http.StatusForbidden)
		return
	}

	// Remove user from channel
	if err := h.messagingService.RemoveUserFromChannel(r.Context(), userID, channelID); err != nil {
		if err == repository.ErrChannelNotFound || err == repository.ErrUserNotMember {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err == repository.ErrCannotRemoveOwner {
			http.Error(w, "Cannot remove the channel owner", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}

// GetMessage handles retrieving a message by ID
func (h *MessagingHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	// Parse message ID from URL
	params := mux.Vars(r)
	messageIDStr := params["messageID"]
	messageIDUint, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	messageID := uint(messageIDUint)

	// Get user ID from context
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get message
	message, err := h.messagingService.GetMessage(r.Context(), messageID)
	if err != nil {
		if err == repository.ErrMessageNotFound {
			http.Error(w, "Message not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	respond(w, http.StatusOK, message)
}

// UpdateMessage handles updating a message
func (h *MessagingHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	// Parse message ID from URL
	params := mux.Vars(r)
	messageIDStr := params["messageID"]
	messageIDUint, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	messageID := uint(messageIDUint)

	// Parse request
	var req struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get existing message
	message, err := h.messagingService.GetMessage(r.Context(), messageID)
	if err != nil {
		if err == repository.ErrMessageNotFound {
			http.Error(w, "Message not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify user is the sender
	if message.SenderID != userID {
		http.Error(w, "Forbidden: can only edit your own messages", http.StatusForbidden)
		return
	}

	// Update message
	message.Content = req.Content
	if err := h.messagingService.UpdateMessage(r.Context(), message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	respond(w, http.StatusOK, message)
}

// DeleteMessage handles deleting a message
func (h *MessagingHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	// Parse message ID from URL
	params := mux.Vars(r)
	messageIDStr := params["messageID"]
	messageIDUint, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	messageID := uint(messageIDUint)

	// Get user ID from context
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get existing message to check ownership
	message, err := h.messagingService.GetMessage(r.Context(), messageID)
	if err != nil {
		if err == repository.ErrMessageNotFound {
			http.Error(w, "Message not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is the sender
	if message.SenderID != userID {
		// Check if user is channel owner (they can delete any message)
		channel, err := h.messagingService.GetChannel(r.Context(), message.ChannelID)
		if err != nil || channel.OwnerID != userID {
			http.Error(w, "Forbidden: can only delete your own messages", http.StatusForbidden)
			return
		}
	}

	// Delete message
	if err := h.messagingService.DeleteMessage(r.Context(), messageID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}

// AddReaction handles adding a reaction to a message
func (h *MessagingHandler) AddReaction(w http.ResponseWriter, r *http.Request) {
	// Parse message ID from URL
	params := mux.Vars(r)
	messageIDStr := params["messageID"]
	messageIDUint, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	messageID := uint(messageIDUint)

	// Parse request
	var req struct {
		EmojiCode string `json:"emojiCode"`
		EmojiName string `json:"emojiName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create and add reaction
	reaction := &domain.Reaction{
		MessageID: messageID,
		UserID:    userID,
		EmojiCode: req.EmojiCode,
		EmojiName: req.EmojiName,
	}

	if err := h.messagingService.AddReaction(r.Context(), reaction); err != nil {
		if err == repository.ErrMessageNotFound {
			http.Error(w, "Message not found", http.StatusNotFound)
			return
		}
		if err == repository.ErrDuplicateReaction {
			http.Error(w, "User has already used this reaction on this message", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	respond(w, http.StatusCreated, reaction)
}

// RemoveReaction handles removing a reaction from a message
func (h *MessagingHandler) RemoveReaction(w http.ResponseWriter, r *http.Request) {
	// Parse message ID and reaction ID from URL
	params := mux.Vars(r)

	reactionIDStr := params["reactionID"]
	reactionIDUint, err := strconv.ParseUint(reactionIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid reaction ID", http.StatusBadRequest)
		return
	}
	reactionID := uint(reactionIDUint)

	// Get user ID from context
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get reaction to check ownership
	reaction, err := h.messagingService.GetReaction(r.Context(), reactionID)
	if err != nil {
		if err == repository.ErrReactionNotFound {
			http.Error(w, "Reaction not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify user owns the reaction
	if reaction.UserID != userID {
		http.Error(w, "Forbidden: can only remove your own reactions", http.StatusForbidden)
		return
	}

	// Remove reaction
	if err := h.messagingService.RemoveReaction(r.Context(), reactionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}

// HandleWebSocket upgrades an HTTP connection to WebSocket
func (h *MessagingHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := getUserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	// Create client and register with hub
	client := websocket.NewClient(conn, userID, h.hub)
	h.hub.Register(client)

	// Start goroutines for reading and writing
	go client.ReadPump()
	go client.WritePump()
}

// Helper functions

// parseChannelID extracts the channel ID from the request URL params
func parseChannelID(r *http.Request) (uint, error) {
	params := mux.Vars(r)
	channelIDStr := params["channelID"]
	channelIDUint, err := strconv.ParseUint(channelIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(channelIDUint), nil
}

// getUserIDFromContext extracts the user ID from the request context
// In a real app, this would be set by your auth middleware
func getUserIDFromContext(ctx context.Context) uint {
	// This is a placeholder - implement based on your auth system
	userID, ok := ctx.Value("userID").(uint)
	if !ok {
		return 0
	}
	return userID
}

// respond sends a JSON response with the given status and data
func respond(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// Global WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins in development
		// In production, validate the origin
		return true
	},
}
