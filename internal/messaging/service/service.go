package service

import (
	"net/http"
	"strconv"
	"time"

	"jadenrazo.dev/internal/messaging/delivery/websocket"
	"jadenrazo.dev/internal/messaging/repository/sqlite"
)

// MessagingService handles the business logic for messaging
type MessagingService struct {
	repo *sqlite.Repository
	hub  *websocket.Hub
}

// NewMessagingService creates a new messaging service
func NewMessagingService(repo *sqlite.Repository, hub *websocket.Hub) *MessagingService {
	return &MessagingService{
		repo: repo,
		hub:  hub,
	}
}

// GetUserIDFromContext extracts the user ID from the context
func GetUserIDFromContext(r *http.Request) uint {
	// In a real app, this would extract user ID from JWT or session
	// For now, use a query parameter for testing
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		return 0
	}

	userID, _ := strconv.ParseUint(userIDStr, 10, 64)
	return uint(userID)
}

// respondWithJSON responds with a JSON payload
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// In a real app, use proper JSON encoding
	w.Write([]byte(`{"status":"ok"}`))
}

// ListChannels handles listing all channels
func (s *MessagingService) ListChannels(w http.ResponseWriter, r *http.Request) {
	// Get all public channels for now
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Channels listed"})
}

// CreateChannel handles creating a new channel
func (s *MessagingService) CreateChannel(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	isPrivateStr := r.FormValue("isPrivate")
	isPrivate := isPrivateStr == "true"

	channel := &sqlite.Channel{
		Name:        name,
		Description: description,
		OwnerID:     userID,
		IsPrivate:   isPrivate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateChannel(channel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, channel)
}

// GetChannel handles retrieving a channel by ID
func (s *MessagingService) GetChannel(w http.ResponseWriter, r *http.Request) {
	// Extract channel ID from URL
	channelIDStr := r.PathValue("channelID")
	if channelIDStr == "" {
		http.Error(w, "Missing channel ID", http.StatusBadRequest)
		return
	}

	channelID, err := strconv.ParseUint(channelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	channel, err := s.repo.GetChannel(uint(channelID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, channel)
}

// UpdateChannel handles updating a channel
func (s *MessagingService) UpdateChannel(w http.ResponseWriter, r *http.Request) {
	// Not implemented yet
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// DeleteChannel handles deleting a channel
func (s *MessagingService) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	// Not implemented yet
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// GetMessages handles retrieving messages from a channel
func (s *MessagingService) GetMessages(w http.ResponseWriter, r *http.Request) {
	// Extract channel ID from URL
	channelIDStr := r.PathValue("channelID")
	if channelIDStr == "" {
		http.Error(w, "Missing channel ID", http.StatusBadRequest)
		return
	}

	channelID, err := strconv.ParseUint(channelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	// Parse pagination params
	lastMessageIDStr := r.URL.Query().Get("lastMessageId")
	lastMessageID := uint(0)
	if lastMessageIDStr != "" {
		id, err := strconv.ParseUint(lastMessageIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid lastMessageId", http.StatusBadRequest)
			return
		}
		lastMessageID = uint(id)
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	messages, err := s.repo.GetChannelMessages(uint(channelID), lastMessageID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, messages)
}

// SendMessage handles sending a new message to a channel
func (s *MessagingService) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract channel ID from URL
	channelIDStr := r.PathValue("channelID")
	if channelIDStr == "" {
		http.Error(w, "Missing channel ID", http.StatusBadRequest)
		return
	}

	channelID, err := strconv.ParseUint(channelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Missing message content", http.StatusBadRequest)
		return
	}

	message := &sqlite.Message{
		Content:   content,
		SenderID:  userID,
		ChannelID: uint(channelID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateMessage(message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast to WebSocket clients
	s.hub.BroadcastToChannel(uint(channelID), "new_message", message)

	respondWithJSON(w, http.StatusCreated, message)
}

// GetChannelMembers handles retrieving members of a channel
func (s *MessagingService) GetChannelMembers(w http.ResponseWriter, r *http.Request) {
	// Extract channel ID from URL
	channelIDStr := r.PathValue("channelID")
	if channelIDStr == "" {
		http.Error(w, "Missing channel ID", http.StatusBadRequest)
		return
	}

	channelID, err := strconv.ParseUint(channelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	members, err := s.repo.GetChannelMembers(uint(channelID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, members)
}

// AddChannelMember handles adding a user to a channel
func (s *MessagingService) AddChannelMember(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract channel ID from URL
	channelIDStr := r.PathValue("channelID")
	if channelIDStr == "" {
		http.Error(w, "Missing channel ID", http.StatusBadRequest)
		return
	}

	channelID, err := strconv.ParseUint(channelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	memberIDStr := r.FormValue("userId")
	if memberIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	memberID, err := strconv.ParseUint(memberIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if user has permission to add members (owner or admin)
	channel, err := s.repo.GetChannel(uint(channelID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if channel.OwnerID != userID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	role := r.FormValue("role")
	if role == "" {
		role = "member"
	}

	err = s.repo.AddChannelMember(uint(memberID), uint(channelID), role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify through WebSocket
	s.hub.BroadcastToChannel(uint(channelID), "member_added", map[string]interface{}{
		"channelId": channelID,
		"userId":    memberID,
		"role":      role,
	})

	w.WriteHeader(http.StatusNoContent)
}

// RemoveChannelMember handles removing a user from a channel
func (s *MessagingService) RemoveChannelMember(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract channel ID from URL
	channelIDStr := r.PathValue("channelID")
	if channelIDStr == "" {
		http.Error(w, "Missing channel ID", http.StatusBadRequest)
		return
	}

	channelID, err := strconv.ParseUint(channelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	// Extract user ID to remove from URL
	memberIDStr := r.PathValue("userID")
	if memberIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	memberID, err := strconv.ParseUint(memberIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if user has permission to remove members (owner or admin)
	channel, err := s.repo.GetChannel(uint(channelID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if channel.OwnerID != userID && uint(memberID) != userID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	err = s.repo.RemoveChannelMember(uint(memberID), uint(channelID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify through WebSocket
	s.hub.BroadcastToChannel(uint(channelID), "member_removed", map[string]interface{}{
		"channelId": channelID,
		"userId":    memberID,
	})

	w.WriteHeader(http.StatusNoContent)
}

// GetUserChannels handles retrieving channels for a user
func (s *MessagingService) GetUserChannels(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	channels, err := s.repo.GetUserChannels(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, channels)
}

// GetUserStatus handles retrieving a user's status
func (s *MessagingService) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	// Not implemented yet
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// UpdateUserStatus handles updating a user's status
func (s *MessagingService) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	// Not implemented yet
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
