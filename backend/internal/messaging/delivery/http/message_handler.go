package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JadenRazo/Project-Website/backend/internal/contexts"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// MessageUseCase defines the interface for message operations
type MessageUseCase interface {
	CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error)
	GetMessage(ctx context.Context, id string) (*domain.Message, error)
	UpdateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error)
	DeleteMessage(ctx context.Context, id string) error
	ListMessages(ctx context.Context, channelID string, before string, after string, limit int) ([]*domain.Message, error)
	AddReaction(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error)
	RemoveReaction(ctx context.Context, messageID string, userID string, emoji string) error
	IsChannelMember(ctx context.Context, channelID string, userID string) (bool, error)
	IsChannelAdmin(ctx context.Context, channelID string, userID string) (bool, error)
	MarkChannelAsRead(ctx context.Context, channelID string, userID string, lastReadMessageID string) error
}

// MessageHandler represents the HTTP handler for message-related requests
type MessageHandler struct {
	messageUseCase domain.MessageUseCase
	validator      *validator.Validate
}

// NewMessageHandler creates a new MessageHandler
func NewMessageHandler(messageUseCase domain.MessageUseCase, validator *validator.Validate) *MessageHandler {
	if messageUseCase == nil {
		panic("messageUseCase cannot be nil")
	}
	if validator == nil {
		panic("validator cannot be nil")
	}

	return &MessageHandler{
		messageUseCase: messageUseCase,
		validator:      validator,
	}
}

// RegisterRoutes registers the message routes
func (h *MessageHandler) RegisterRoutes(r chi.Router) {
	if r == nil {
		panic("router cannot be nil")
	}

	r.Post("/channels/{channelID}/messages", h.CreateMessage)
	r.Get("/channels/{channelID}/messages", h.ListChannelMessages)
	r.Get("/messages/{messageID}", h.GetMessage)
	r.Put("/messages/{messageID}", h.UpdateMessage)
	r.Delete("/messages/{messageID}", h.DeleteMessage)
	r.Post("/messages/{messageID}/reactions", h.AddReaction)
	r.Delete("/messages/{messageID}/reactions/{emoji}", h.RemoveReaction)
	r.Post("/channels/{channelID}/read", h.handleMarkChannelRead)
}

// CreateMessageRequest represents a request to create a message
type CreateMessageRequest struct {
	Content     string   `json:"content" validate:"required_without=Attachments"`
	Attachments []string `json:"attachments,omitempty"`
	Mentions    []string `json:"mentions,omitempty" validate:"dive,uuid"`
	ReplyTo     string   `json:"reply_to,omitempty" validate:"omitempty,uuid"`
}

// UpdateMessageRequest represents a request to update a message
type UpdateMessageRequest struct {
	Content  string   `json:"content" validate:"required_without=Attachments"`
	Mentions []string `json:"mentions,omitempty" validate:"dive,uuid"`
}

// AddReactionRequest represents a request to add a reaction
type AddReactionRequest struct {
	Emoji string `json:"emoji" validate:"required"`
}

// ToMessageResponse converts a domain message to response format
func ToMessageResponse(message *domain.Message) MessageResponse {
	response := MessageResponse{
		ID:        message.ID.String(),
		ChannelID: message.ChannelID.String(),
		UserID:    message.UserID.String(),
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}

	// Handle reply to ID if present
	// ReplyToID is handled in the MessageResponse type, we don't need to map it here

	// Convert attachments
	var attachments []AttachmentResponse
	for _, attachment := range message.Attachments {
		attachments = append(attachments, AttachmentResponse{
			ID:        attachment.ID.String(),
			MessageID: attachment.MessageID.String(),
			URL:       attachment.FileURL,
			Name:      attachment.FileName,
			Size:      attachment.FileSize,
			Type:      attachment.FileType,
		})
	}
	response.Attachments = attachments

	// Convert reactions
	var reactions []ReactionResponse
	for _, reaction := range message.Reactions {
		reactions = append(reactions, ReactionResponse{
			ID:        reaction.ID.String(),
			MessageID: reaction.MessageID.String(),
			UserID:    reaction.UserID.String(),
			Emoji:     reaction.Emoji,
			CreatedAt: reaction.CreatedAt,
		})
	}
	response.Reactions = reactions

	return response
}

// handleCreateMessage handles the request to create a message
func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if channelID == "" {
		RespondWithError(w, http.StatusBadRequest, "Channel ID is required")
		return
	}

	// Get user ID from context
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Parse request body
	var req struct {
		Content   string            `json:"content" validate:"required"`
		Metadata  map[string]string `json:"metadata,omitempty"`
		ParentID  string            `json:"parent_id,omitempty" validate:"omitempty,uuid"`
		Mentions  []string          `json:"mentions,omitempty" validate:"omitempty,dive,uuid"`
		Reactions []string          `json:"reactions,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		validationErrors := ProcessValidationErrors(err)
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+validationErrors)
		return
	}

	// Create message
	message, err := h.messageUseCase.CreateMessage(r.Context(), channelID, userID, req.Content, req.ParentID, req.Metadata)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Process mentions if any
	if len(req.Mentions) > 0 {
		err = h.messageUseCase.AddMentions(r.Context(), message.ID, req.Mentions)
		if err != nil {
			// Log error but continue
			// In a production system, this should be handled more gracefully
		}
	}

	// Return the created message
	RespondWithJSON(w, http.StatusCreated, message)
}

// handleUpdateMessage handles the request to update a message
func (h *MessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "messageID")
	if messageID == "" {
		RespondWithError(w, http.StatusBadRequest, "Message ID is required")
		return
	}

	// Get user ID from context
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Parse request body
	var req struct {
		Content  string            `json:"content" validate:"required"`
		Metadata map[string]string `json:"metadata,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		validationErrors := ProcessValidationErrors(err)
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+validationErrors)
		return
	}

	// Update message
	message, err := h.messageUseCase.UpdateMessage(r.Context(), messageID, userID, req.Content, req.Metadata)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Return updated message
	RespondWithJSON(w, http.StatusOK, message)
}

// handleGetMessage handles the request to get a message
func (h *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "messageID")
	if messageID == "" {
		RespondWithError(w, http.StatusBadRequest, "Message ID is required")
		return
	}

	// Get message
	message, err := h.messageUseCase.GetMessage(r.Context(), messageID)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Return message
	RespondWithJSON(w, http.StatusOK, message)
}

// handleDeleteMessage handles the request to delete a message
func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "messageID")
	if messageID == "" {
		RespondWithError(w, http.StatusBadRequest, "Message ID is required")
		return
	}

	// Get user ID from context
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Delete message
	err := h.messageUseCase.DeleteMessage(r.Context(), messageID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Return success
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Message deleted successfully"})
}

// handleListMessages handles the request to list messages
func (h *MessageHandler) ListChannelMessages(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if channelID == "" {
		RespondWithError(w, http.StatusBadRequest, "Channel ID is required")
		return
	}

	// Get pagination parameters
	pagination := getPaginationFromRequest(r)

	// Get messages from the use case
	messages, total, err := h.messageUseCase.GetChannelMessages(r.Context(), channelID, pagination.Offset, pagination.PerPage)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Create response
	response := createPaginatedResponse(messages, pagination, total)
	RespondWithJSON(w, http.StatusOK, response)
}

// handleAddReaction handles the request to add a reaction to a message
func (h *MessageHandler) AddReaction(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "messageID")
	if messageID == "" {
		RespondWithError(w, http.StatusBadRequest, "Message ID is required")
		return
	}

	// Get user ID from context
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Parse request body
	var req struct {
		Emoji string `json:"emoji" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		validationErrors := ProcessValidationErrors(err)
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+validationErrors)
		return
	}

	// Add reaction
	reaction, err := h.messageUseCase.AddReaction(r.Context(), messageID, userID, req.Emoji)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Return reaction
	RespondWithJSON(w, http.StatusCreated, reaction)
}

// handleRemoveReaction handles the request to remove a reaction from a message
func (h *MessageHandler) RemoveReaction(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "messageID")
	if messageID == "" {
		RespondWithError(w, http.StatusBadRequest, "Message ID is required")
		return
	}

	emoji := chi.URLParam(r, "emoji")
	if emoji == "" {
		RespondWithError(w, http.StatusBadRequest, "Emoji is required")
		return
	}

	// Get user ID from context
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Remove reaction
	err := h.messageUseCase.RemoveReaction(r.Context(), messageID, userID, emoji)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Return success
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Reaction removed successfully"})
}

// handleMarkChannelRead handles the request to mark a channel as read
func (h *MessageHandler) handleMarkChannelRead(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if channelID == "" {
		RespondWithError(w, http.StatusBadRequest, "Channel ID is required")
		return
	}

	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Check if user is a member of the channel
	isMember, err := h.messageUseCase.IsChannelMember(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	if !isMember {
		RespondWithError(w, http.StatusForbidden, "You are not a member of this channel")
		return
	}

	// Parse last read message ID if provided
	var lastReadMessageID string
	type markReadRequest struct {
		LastReadMessageID string `json:"last_read_message_id,omitempty"`
	}

	var req markReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
		lastReadMessageID = req.LastReadMessageID
	}

	// Mark channel as read
	if err := h.messageUseCase.MarkChannelAsRead(r.Context(), channelID, userID, lastReadMessageID); err != nil {
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// processValidationErrors formats validation errors into a readable string
func processValidationErrors(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorMessages := make([]string, 0, len(validationErrors))
		for _, fieldError := range validationErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", fieldError.Field(), fieldError.Tag()))
		}
		return fmt.Sprintf("%v", errorMessages)
	}
	return err.Error()
}

// handleMessageErrors maps domain errors to HTTP status codes and responses
func handleMessageErrors(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch err.Error() {
	case "resource not found", "message not found", "channel not found":
		RespondWithError(w, http.StatusNotFound, err.Error())
	case "forbidden access", "not a channel member":
		RespondWithError(w, http.StatusForbidden, err.Error())
	case "bad request":
		RespondWithError(w, http.StatusBadRequest, err.Error())
	case "entity already exists", "duplicate reaction":
		RespondWithError(w, http.StatusConflict, err.Error())
	default:
		// Log unexpected errors
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}
