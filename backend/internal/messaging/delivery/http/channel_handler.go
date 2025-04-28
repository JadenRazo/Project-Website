package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/JadenRazo/Project-Website/backend/internal/contexts"
	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ChannelHandler represents a channel handler
type ChannelHandler struct {
	channelUseCase domain.ChannelUseCase
	validator      *validator.Validate
}

// NewChannelHandler creates a new channel handler
func NewChannelHandler(channelUseCase domain.ChannelUseCase, validator *validator.Validate) *ChannelHandler {
	if channelUseCase == nil {
		panic("channelUseCase cannot be nil")
	}
	if validator == nil {
		panic("validator cannot be nil")
	}

	return &ChannelHandler{
		channelUseCase: channelUseCase,
		validator:      validator,
	}
}

// RegisterRoutes registers the channel routes
func (h *ChannelHandler) RegisterRoutes(r chi.Router) {
	if r == nil {
		panic("router cannot be nil")
	}

	r.Post("/channels", h.handleCreateChannel)
	r.Get("/channels", h.handleListChannels)
	r.Get("/channels/{channelID}", h.handleGetChannel)
	r.Put("/channels/{channelID}", h.handleUpdateChannel)
	r.Delete("/channels/{channelID}", h.handleDeleteChannel)

	r.Post("/channels/{channelID}/members", h.handleAddMember)
	r.Get("/channels/{channelID}/members", h.handleListMembers)
	r.Delete("/channels/{channelID}/members/{memberID}", h.handleRemoveMember)
	r.Put("/channels/{channelID}/members/{memberID}", h.handleUpdateMember)
}

// ChannelCreateRequest represents a request to create a channel
type ChannelCreateRequest struct {
	Name        string   `json:"name" validate:"required,min=3,max=100"`
	Description string   `json:"description,omitempty" validate:"max=500"`
	IsPrivate   bool     `json:"is_private"`
	MemberIDs   []string `json:"member_ids,omitempty" validate:"dive,uuid"`
}

// ChannelUpdateRequest represents a request to update a channel
type ChannelUpdateRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description string `json:"description,omitempty" validate:"max=500"`
	IsPrivate   *bool  `json:"is_private,omitempty"`
}

// MemberRequest represents a request to add or update a member
type MemberRequest struct {
	Role string `json:"role" validate:"required,oneof=member admin"`
}

// Helper response functions for this handler
// ChannelResponse converts a channel to a map for response
func channelToResponse(channel *domain.MessagingChannel) map[string]interface{} {
	return map[string]interface{}{
		"id":          channel.ID,
		"name":        channel.Name,
		"description": channel.Description,
		"is_private":  channel.IsPrivate,
		"created_at":  channel.CreatedAt,
		"updated_at":  channel.UpdatedAt,
	}
}

// channelsToResponse converts a slice of channels to a slice of response maps
func channelsToResponseList(channels []domain.MessagingChannel) []map[string]interface{} {
	responseList := make([]map[string]interface{}, len(channels))
	for i, channel := range channels {
		channelCopy := channel // Create a copy to avoid pointer issues
		responseList[i] = channelToResponse(&channelCopy)
	}
	return responseList
}

// memberToResponse converts a member to a map for response
func memberToResponse(member *domain.MessagingMember) map[string]interface{} {
	return map[string]interface{}{
		"id":         member.ID,
		"channel_id": member.ChannelID,
		"user_id":    member.UserID,
		"role":       member.Role,
		"joined_at":  member.JoinedAt,
		"updated_at": member.UpdatedAt,
	}
}

// membersToResponseList converts a slice of members to a slice of response maps
func membersToResponseList(members []domain.MessagingMember) []map[string]interface{} {
	responseList := make([]map[string]interface{}, len(members))
	for i, member := range members {
		memberCopy := member // Create a copy to avoid pointer issues
		responseList[i] = memberToResponse(&memberCopy)
	}
	return responseList
}

// formatValidationErrors formats validation errors into a readable string
func formatValidationErrors(err error) string {
	if ve, ok := err.(validator.ValidationErrors); ok {
		errors := make([]string, len(ve))
		for i, fe := range ve {
			errors[i] = fmt.Sprintf("%s: %s", fe.Field(), fe.Tag())
		}
		return fmt.Sprintf("%v", errors)
	}
	return err.Error()
}

// handleCreateChannel handles the request to create a channel
func (h *ChannelHandler) handleCreateChannel(w http.ResponseWriter, r *http.Request) {
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	var req ChannelCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(req); err != nil {
		validationErrors := ProcessValidationErrors(err)
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+validationErrors)
		return
	}

	channel := &domain.MessagingChannel{
		Name:        req.Name,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
	}

	// If member IDs are provided, create initial members
	members := make([]*domain.MessagingMember, 0, len(req.MemberIDs)+1)

	// Always add the creator as an admin
	creatorMember := &domain.MessagingMember{
		UserID: userID,
		Role:   "admin",
	}
	members = append(members, creatorMember)

	// Add other members
	for _, memberID := range req.MemberIDs {
		// Skip if it's the creator
		if memberID == userID {
			continue
		}

		// Validate UUID format
		if _, err := uuid.Parse(memberID); err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid member ID format: %s", memberID))
			return
		}

		member := &domain.MessagingMember{
			UserID: memberID,
			Role:   "member",
		}
		members = append(members, member)
	}

	createdChannel, err := h.channelUseCase.CreateChannel(r.Context(), channel, members)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, CreateChannelResponse(createdChannel))
}

// handleUpdateChannel handles the request to update a channel
func (h *ChannelHandler) handleUpdateChannel(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if _, err := uuid.Parse(channelID); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid channel ID format")
		return
	}

	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	var req ChannelUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(req); err != nil {
		validationErrors := ProcessValidationErrors(err)
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+validationErrors)
		return
	}

	// Check if user is an admin of the channel
	isAdmin, err := h.channelUseCase.IsChannelAdmin(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	if !isAdmin {
		RespondWithError(w, http.StatusForbidden, "Only channel admins can update the channel")
		return
	}

	// Get existing channel
	channel, err := h.channelUseCase.GetChannel(r.Context(), channelID)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Update only provided fields
	if req.Name != "" {
		channel.Name = req.Name
	}
	if req.Description != "" {
		channel.Description = req.Description
	}
	if req.IsPrivate != nil {
		channel.IsPrivate = *req.IsPrivate
	}

	updatedChannel, err := h.channelUseCase.UpdateChannel(r.Context(), channel)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, CreateChannelResponse(updatedChannel))
}

// handleGetChannel handles the request to get a channel
func (h *ChannelHandler) handleGetChannel(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if _, err := uuid.Parse(channelID); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid channel ID format")
		return
	}

	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	channel, err := h.channelUseCase.GetChannel(r.Context(), channelID)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Check if user has access to the channel
	if channel.IsPrivate {
		isAdmin, err := h.channelUseCase.IsChannelAdmin(r.Context(), channelID, userID)
		if err != nil {
			HandleError(w, err)
			return
		}

		if !isAdmin {
			// Check if user is a member
			members, _, err := h.channelUseCase.ListChannelMembers(r.Context(), channelID, &domain.Pagination{Page: 1, PageSize: 1000})
			if err != nil {
				HandleError(w, err)
				return
			}

			isMember := false
			for _, member := range members {
				if member.UserID == userID {
					isMember = true
					break
				}
			}

			if !isMember {
				RespondWithError(w, http.StatusForbidden, "You don't have access to this channel")
				return
			}
		}
	}

	RespondWithJSON(w, http.StatusOK, CreateChannelResponse(channel))
}

// handleDeleteChannel handles the request to delete a channel
func (h *ChannelHandler) handleDeleteChannel(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if _, err := uuid.Parse(channelID); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid channel ID format")
		return
	}

	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Check if the channel exists
	_, err := h.channelUseCase.GetChannel(r.Context(), channelID)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Check if user is an admin
	isAdmin, err := h.channelUseCase.IsChannelAdmin(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	if !isAdmin {
		RespondWithError(w, http.StatusForbidden, "Only channel admins can delete the channel")
		return
	}

	if err := h.channelUseCase.DeleteChannel(r.Context(), channelID); err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Channel deleted successfully",
	})
}

// handleListChannels handles the request to list channels
func (h *ChannelHandler) handleListChannels(w http.ResponseWriter, r *http.Request) {
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Parse pagination parameters
	pagination := &domain.Pagination{
		Page:     1,
		PageSize: 10,
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil && page > 0 {
			pagination.Page = page
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err == nil && pageSize > 0 && pageSize <= 100 {
			pagination.PageSize = pageSize
			pagination.PerPage = pageSize
		}
	} else if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		perPage, err := strconv.Atoi(perPageStr)
		if err == nil && perPage > 0 && perPage <= 100 {
			pagination.PageSize = perPage
			pagination.PerPage = perPage
		}
	}

	channels, total, err := h.channelUseCase.ListChannels(r.Context(), userID, pagination)
	if err != nil {
		HandleError(w, err)
		return
	}

	pagination.Total = total

	// Create response channels list
	channelResponses := make([]ChannelResponse, len(channels))
	for i, ch := range channels {
		channelCopy := ch
		channelResponses[i] = CreateChannelResponse(&channelCopy)
	}

	RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"channels":   channelResponses,
		"pagination": pagination,
	})
}

// handleAddMember handles the request to add a member to a channel
func (h *ChannelHandler) handleAddMember(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if _, err := uuid.Parse(channelID); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid channel ID format")
		return
	}

	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Check if user is admin of the channel
	isAdmin, err := h.channelUseCase.IsChannelAdmin(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	if !isAdmin {
		RespondWithError(w, http.StatusForbidden, "Only channel admins can add members")
		return
	}

	var req MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(req); err != nil {
		validationErrors := ProcessValidationErrors(err)
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+validationErrors)
		return
	}

	// Get the member ID from the request
	memberID := r.URL.Query().Get("user_id")
	if memberID == "" {
		RespondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	if _, err := uuid.Parse(memberID); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	member, err := h.channelUseCase.AddChannelMember(r.Context(), channelID, memberID, req.Role)
	if err != nil {
		if err.Error() == domain.ErrDuplicateMember {
			RespondWithError(w, http.StatusConflict, "Member already exists in the channel")
			return
		}
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, CreateMemberResponse(member))
}

// handleRemoveMember handles the request to remove a member from a channel
func (h *ChannelHandler) handleRemoveMember(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if _, err := uuid.Parse(channelID); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid channel ID format")
		return
	}

	memberID := chi.URLParam(r, "memberID")
	if _, err := uuid.Parse(memberID); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid member ID format")
		return
	}

	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Check if the channel exists
	_, err := h.channelUseCase.GetChannel(r.Context(), channelID)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Check if user is an admin
	isAdmin, err := h.channelUseCase.IsChannelAdmin(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Users can remove themselves
	if memberID != userID && !isAdmin {
		RespondWithError(w, http.StatusForbidden, "Only channel admins can remove other members")
		return
	}

	if err := h.channelUseCase.RemoveChannelMember(r.Context(), channelID, memberID); err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Member removed successfully",
	})
}

// handleUpdateMember handles the request to update a member's role
func (h *ChannelHandler) handleUpdateMember(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if _, err := uuid.Parse(channelID); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid channel ID format")
		return
	}

	memberID := chi.URLParam(r, "memberID")
	if memberID == "" {
		RespondWithError(w, http.StatusBadRequest, "Member ID is required")
		return
	}

	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Check if user is an admin
	isAdmin, err := h.channelUseCase.IsChannelAdmin(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	if !isAdmin {
		RespondWithError(w, http.StatusForbidden, "Only channel admins can update member roles")
		return
	}

	var req MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(req); err != nil {
		validationErrors := ProcessValidationErrors(err)
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+validationErrors)
		return
	}

	// Check if the channel exists
	_, err = h.channelUseCase.GetChannel(r.Context(), channelID)
	if err != nil {
		HandleError(w, err)
		return
	}

	member, err := h.channelUseCase.UpdateChannelMember(r.Context(), channelID, memberID, req.Role)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, CreateMemberResponse(member))
}

// handleListMembers handles the request to list members of a channel
func (h *ChannelHandler) handleListMembers(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if _, err := uuid.Parse(channelID); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid channel ID format")
		return
	}

	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	channel, err := h.channelUseCase.GetChannel(r.Context(), channelID)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Check if user has access to the channel
	isAdmin, err := h.channelUseCase.IsChannelAdmin(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	if channel.IsPrivate && !isAdmin {
		// Check if user is a member
		members, _, err := h.channelUseCase.ListChannelMembers(r.Context(), channelID, &domain.Pagination{Page: 1, PageSize: 1000})
		if err != nil {
			HandleError(w, err)
			return
		}

		isMember := false
		for _, member := range members {
			if member.UserID == userID {
				isMember = true
				break
			}
		}

		if !isMember {
			RespondWithError(w, http.StatusForbidden, "You don't have access to this channel")
			return
		}
	}

	// Parse pagination parameters
	pagination := &domain.Pagination{
		Page:     1,
		PageSize: 10,
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil && page > 0 {
			pagination.Page = page
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err == nil && pageSize > 0 && pageSize <= 100 {
			pagination.PageSize = pageSize
		}
	}

	members, total, err := h.channelUseCase.ListChannelMembers(r.Context(), channelID, pagination)
	if err != nil {
		HandleError(w, err)
		return
	}

	pagination.Total = total

	// Create response members list
	memberResponses := make([]MemberResponse, len(members))
	for i, m := range members {
		memberCopy := m
		memberResponses[i] = CreateMemberResponse(&memberCopy)
	}

	RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"members":    memberResponses,
		"pagination": pagination,
	})
}

// ListChannels handles listing all channels
func (h *ChannelHandler) ListChannels(w http.ResponseWriter, r *http.Request) {
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	pagination := getPaginationFromRequest(r)
	channels, total, err := h.channelUseCase.ListChannels(r.Context(), userID, pagination.Offset, pagination.PerPage)
	if err != nil {
		HandleError(w, err)
		return
	}

	response := createPaginatedResponse(channelsToResponseList(channels), pagination, total)
	RespondWithJSON(w, http.StatusOK, response)
}

// CreateChannel handles creating a new channel
func (h *ChannelHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	var req ChannelCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(req); err != nil {
		validationErrors := ProcessValidationErrors(err)
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+validationErrors)
		return
	}

	channel := &domain.MessagingChannel{
		Name:        req.Name,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
	}

	createdChannel, err := h.channelUseCase.CreateChannel(r.Context(), channel, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Add members if specified
	if len(req.MemberIDs) > 0 {
		for _, memberID := range req.MemberIDs {
			_, err := h.channelUseCase.AddMember(r.Context(), createdChannel.ID, memberID, "member")
			if err != nil {
				// Log the error but continue
				fmt.Printf("Failed to add member %s: %v\n", memberID, err)
			}
		}
	}

	RespondWithJSON(w, http.StatusCreated, CreateChannelResponse(createdChannel))
}

// ListDirectChannels handles listing direct message channels
func (h *ChannelHandler) ListDirectChannels(w http.ResponseWriter, r *http.Request) {
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	pagination := getPaginationFromRequest(r)
	channels, total, err := h.channelUseCase.ListDirectChannels(r.Context(), userID, pagination.Offset, pagination.PerPage)
	if err != nil {
		HandleError(w, err)
		return
	}

	response := createPaginatedResponse(channelsToResponseList(channels), pagination, total)
	RespondWithJSON(w, http.StatusOK, response)
}

// CreateDirectChannel handles creating a direct message channel
func (h *ChannelHandler) CreateDirectChannel(w http.ResponseWriter, r *http.Request) {
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	var req struct {
		TargetUserID string `json:"target_user_id" validate:"required,uuid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(req); err != nil {
		validationErrors := ProcessValidationErrors(err)
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+validationErrors)
		return
	}

	channel, err := h.channelUseCase.CreateDirectChannel(r.Context(), userID, req.TargetUserID)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, CreateChannelResponse(channel))
}

// GetChannel handles retrieving a channel by ID
func (h *ChannelHandler) GetChannel(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if channelID == "" {
		RespondWithError(w, http.StatusBadRequest, "Channel ID is required")
		return
	}

	channel, err := h.channelUseCase.GetChannel(r.Context(), channelID)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, CreateChannelResponse(channel))
}

// DeleteChannel handles deleting a channel
func (h *ChannelHandler) DeleteChannel(w http.ResponseWriter, r *http.Request) {
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

	// Check if user is an admin of the channel
	isAdmin, err := h.channelUseCase.IsChannelAdmin(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	if !isAdmin {
		RespondWithError(w, http.StatusForbidden, "Only channel admins can delete the channel")
		return
	}

	err = h.channelUseCase.DeleteChannel(r.Context(), channelID)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Channel deleted successfully"})
}

// ListChannelMembers handles listing members of a channel
func (h *ChannelHandler) ListChannelMembers(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if channelID == "" {
		RespondWithError(w, http.StatusBadRequest, "Channel ID is required")
		return
	}

	pagination := getPaginationFromRequest(r)
	members, total, err := h.channelUseCase.ListMembers(r.Context(), channelID, pagination.Offset, pagination.PerPage)
	if err != nil {
		HandleError(w, err)
		return
	}

	response := createPaginatedResponse(membersToResponseList(members), pagination, total)
	RespondWithJSON(w, http.StatusOK, response)
}

// AddChannelMember handles adding a member to a channel
func (h *ChannelHandler) AddChannelMember(w http.ResponseWriter, r *http.Request) {
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

	// Check if user is an admin of the channel
	isAdmin, err := h.channelUseCase.IsChannelAdmin(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	if !isAdmin {
		RespondWithError(w, http.StatusForbidden, "Only channel admins can add members")
		return
	}

	var req struct {
		UserID string `json:"user_id" validate:"required,uuid"`
		Role   string `json:"role" validate:"required,oneof=member admin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(req); err != nil {
		validationErrors := ProcessValidationErrors(err)
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+validationErrors)
		return
	}

	member, err := h.channelUseCase.AddMember(r.Context(), channelID, req.UserID, req.Role)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, memberToResponse(member))
}

// RemoveChannelMember handles removing a member from a channel
func (h *ChannelHandler) RemoveChannelMember(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if channelID == "" {
		RespondWithError(w, http.StatusBadRequest, "Channel ID is required")
		return
	}

	targetUserID := chi.URLParam(r, "userID")
	if targetUserID == "" {
		RespondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Check if user is an admin of the channel
	isAdmin, err := h.channelUseCase.IsChannelAdmin(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Users can remove themselves
	if targetUserID != userID && !isAdmin {
		RespondWithError(w, http.StatusForbidden, "Only channel admins can remove other members")
		return
	}

	err = h.channelUseCase.RemoveMember(r.Context(), channelID, targetUserID)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Member removed successfully"})
}

// LeaveChannel handles a user leaving a channel
func (h *ChannelHandler) LeaveChannel(w http.ResponseWriter, r *http.Request) {
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

	err := h.channelUseCase.RemoveMember(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Left channel successfully"})
}

// JoinChannel handles a user joining a channel
func (h *ChannelHandler) JoinChannel(w http.ResponseWriter, r *http.Request) {
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

	// Check if channel is public
	channel, err := h.channelUseCase.GetChannel(r.Context(), channelID)
	if err != nil {
		HandleError(w, err)
		return
	}

	if channel.IsPrivate {
		RespondWithError(w, http.StatusForbidden, "Cannot join private channels")
		return
	}

	member, err := h.channelUseCase.AddMember(r.Context(), channelID, userID, "member")
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, memberToResponse(member))
}

// MarkChannelAsRead marks all messages in a channel as read
func (h *ChannelHandler) MarkChannelAsRead(w http.ResponseWriter, r *http.Request) {
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

	err := h.channelUseCase.MarkChannelAsRead(r.Context(), channelID, userID)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Channel marked as read"})
}

// ArchiveChannel archives a channel
func (h *ChannelHandler) ArchiveChannel(w http.ResponseWriter, r *http.Request) {
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

	err := h.channelUseCase.ArchiveChannel(r.Context(), channelID)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Channel archived successfully"})
}

// UnarchiveChannel unarchives a channel
func (h *ChannelHandler) UnarchiveChannel(w http.ResponseWriter, r *http.Request) {
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

	err := h.channelUseCase.UnarchiveChannel(r.Context(), channelID)
	if err != nil {
		HandleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Channel unarchived successfully"})
}

// GetUserChannels returns all channels a user is a member of
func (h *ChannelHandler) GetUserChannels(w http.ResponseWriter, r *http.Request) {
	userID := contexts.GetUserID(r.Context())
	if userID == "" {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	pagination := getPaginationFromRequest(r)
	channels, total, err := h.channelUseCase.GetUserChannels(r.Context(), userID, pagination.Offset, pagination.PerPage)
	if err != nil {
		HandleError(w, err)
		return
	}

	response := createPaginatedResponse(channelsToResponseList(channels), pagination, total)
	RespondWithJSON(w, http.StatusOK, response)
}
