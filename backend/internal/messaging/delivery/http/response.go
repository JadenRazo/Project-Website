package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/go-playground/validator/v10"
)

// Common error constants
var (
	ErrMemberNotFoundMsg  = "member not found"
	ErrChannelNotFoundMsg = "channel not found"
	ErrMessageNotFoundMsg = "message not found"
	ErrNotFound           = errors.New("resource not found")
	ErrForbidden          = errors.New("forbidden access")
	ErrBadRequest         = errors.New("bad request")
	ErrDuplicateEntity    = errors.New("entity already exists")
	ErrDuplicateMember    = errors.New("member already exists in channel")
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationResponse represents pagination information in a response
type PaginationResponse struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

// RespondWithJSON writes a JSON response
func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Error encoding response"})
		}
	}
}

// RespondWithError writes an error response
func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, map[string]string{"error": message})
}

// ProcessValidationErrors formats validation errors into a readable string
func ProcessValidationErrors(err error) string {
	if ve, ok := err.(validator.ValidationErrors); ok {
		errors := make([]string, len(ve))
		for i, fe := range ve {
			errors[i] = fmt.Sprintf("%s: %s", fe.Field(), fe.Tag())
		}
		return fmt.Sprintf("%v", errors)
	}
	return err.Error()
}

// HandleError handles domain-specific errors and returns appropriate HTTP responses
func HandleError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "Internal server error"

	errMsg := err.Error()
	switch errMsg {
	case ErrMemberNotFoundMsg:
		status = http.StatusNotFound
		message = "Member not found"
	case ErrNotFound.Error(), ErrChannelNotFoundMsg, ErrMessageNotFoundMsg:
		status = http.StatusNotFound
		message = "Resource not found"
	case ErrForbidden.Error():
		status = http.StatusForbidden
		message = "Forbidden action"
	case ErrBadRequest.Error():
		status = http.StatusBadRequest
		message = "Invalid request"
	case ErrDuplicateEntity.Error(), ErrDuplicateMember.Error():
		status = http.StatusConflict
		message = "Resource already exists"
	default:
		// Log the error for debugging
		fmt.Printf("Unexpected error: %v\n", err)
	}

	RespondWithError(w, status, message)
}

// BuildPaginationResponse creates a pagination response from domain pagination
func BuildPaginationResponse(total int, pagination domain.Pagination) PaginationResponse {
	totalPages := ComputeTotalPages(total, pagination.PageSize)

	return PaginationResponse{
		CurrentPage: pagination.Page,
		PerPage:     pagination.PageSize,
		TotalItems:  total,
		TotalPages:  totalPages,
	}
}

// ComputeTotalPages calculates the total number of pages based on item count and page size
func ComputeTotalPages(totalItems, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}

	totalPages := totalItems / pageSize
	if totalItems%pageSize > 0 {
		totalPages++
	}

	return totalPages
}

// Channel response types

// ChannelResponse represents a channel in API responses
type ChannelResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ChannelListResponse represents a list of channels in API responses
type ChannelListResponse struct {
	Channels   []ChannelResponse  `json:"channels"`
	Pagination PaginationResponse `json:"pagination"`
}

// CreateChannelResponse creates a ChannelResponse from a domain.MessagingChannel
func CreateChannelResponse(channel *domain.MessagingChannel) ChannelResponse {
	return ChannelResponse{
		ID:          channel.ID,
		Name:        channel.Name,
		Description: channel.Description,
		IsPrivate:   channel.IsPrivate,
		CreatedAt:   FormatTime(channel.CreatedAt),
		UpdatedAt:   FormatTime(channel.UpdatedAt),
	}
}

// MemberResponse represents a channel member in API responses
type MemberResponse struct {
	ID        uint   `json:"id"`
	ChannelID uint   `json:"channel_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	JoinedAt  string `json:"joined_at"`
	UpdatedAt string `json:"updated_at"`
}

// MemberListResponse represents a list of members in API responses
type MemberListResponse struct {
	Members    []MemberResponse   `json:"members"`
	Pagination PaginationResponse `json:"pagination"`
}

// CreateMemberResponse creates a MemberResponse from a domain.MessagingMember
func CreateMemberResponse(member *domain.MessagingMember) MemberResponse {
	return MemberResponse{
		ID:        member.ID,
		ChannelID: member.ChannelID,
		UserID:    member.UserID,
		Role:      member.Role,
		JoinedAt:  member.JoinedAt,
		UpdatedAt: member.UpdatedAt,
	}
}

// MessageResponse represents a message response
type MessageResponse struct {
	ID          string               `json:"id"`
	ChannelID   string               `json:"channel_id"`
	UserID      string               `json:"user_id"`
	Content     string               `json:"content"`
	Reactions   []ReactionResponse   `json:"reactions,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	User        *UserInfoResponse    `json:"user,omitempty"`
	Mentions    []MentionResponse    `json:"mentions,omitempty"`
	Metadata    map[string]string    `json:"metadata,omitempty"`
	Attachments []AttachmentResponse `json:"attachments,omitempty"`
}

// MessageListResponse represents a list of messages in API responses
type MessageListResponse struct {
	Messages   []MessageResponse  `json:"messages"`
	Pagination PaginationResponse `json:"pagination"`
}

// UserInfoResponse represents a simplified user response for embedding in other responses
type UserInfoResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// ReactionResponse represents a reaction in API responses
type ReactionResponse struct {
	ID        string    `json:"id"`
	MessageID string    `json:"message_id"`
	UserID    string    `json:"user_id"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
}

// MentionResponse represents a mention in API responses
type MentionResponse struct {
	ID        string `json:"id"`
	MessageID string `json:"message_id"`
	UserID    string `json:"user_id"`
	Type      string `json:"type"`
}

// AttachmentResponse represents an attachment in API responses
type AttachmentResponse struct {
	ID        string            `json:"id"`
	MessageID string            `json:"message_id"`
	Type      string            `json:"type"`
	URL       string            `json:"url"`
	Name      string            `json:"name"`
	Size      int64             `json:"size"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// FormatTime formats a time.Time into a string using RFC3339 format
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ErrorResponse creates an error response object
func ErrorResponse(message string) map[string]string {
	return map[string]string{"error": message}
}

// NewError creates a new error with the given message
func NewError(message string) error {
	return errors.New(message)
}
