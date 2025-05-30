package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// ErrorType defines the type of error
type ErrorType string

// Error types
const (
	ErrorTypeValidation  ErrorType = "validation"
	ErrorTypeDatabase    ErrorType = "database"
	ErrorTypeAuth        ErrorType = "auth"
	ErrorTypeNotFound    ErrorType = "not_found"
	ErrorTypeRateLimit   ErrorType = "rate_limit"
	ErrorTypeInternal    ErrorType = "internal"
	ErrorTypeWebSocket   ErrorType = "websocket"
	ErrorTypeUnsupported ErrorType = "unsupported"
)

// AppError represents an application-specific error
type AppError struct {
	Type    ErrorType `json:"type"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Cause   error     `json:"-"` // Not serialized
	Status  int       `json:"-"` // HTTP status code
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// JSON returns the error as a JSON string
func (e *AppError) JSON() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf(`{"type":"internal","code":"json_marshal_error","message":"Failed to marshal error"}`)
	}
	return string(data)
}

// HTTPStatusCode returns the appropriate HTTP status code for this error
func (e *AppError) HTTPStatusCode() int {
	if e.Status != 0 {
		return e.Status
	}

	switch e.Type {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeAuth:
		return http.StatusUnauthorized
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeRateLimit:
		return http.StatusTooManyRequests
	case ErrorTypeUnsupported:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

// Error constructors

// NewValidationError creates a new validation error
func NewValidationError(code, message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Code:    code,
		Message: message,
		Cause:   cause,
		Status:  http.StatusBadRequest,
	}
}

// NewDatabaseError creates a new database error
func NewDatabaseError(code, message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeDatabase,
		Code:    code,
		Message: message,
		Cause:   cause,
		Status:  http.StatusInternalServerError,
	}
}

// NewAuthError creates a new authentication error
func NewAuthError(code, message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeAuth,
		Code:    code,
		Message: message,
		Cause:   cause,
		Status:  http.StatusUnauthorized,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(code, message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Code:    code,
		Message: message,
		Cause:   cause,
		Status:  http.StatusNotFound,
	}
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(code, message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeRateLimit,
		Code:    code,
		Message: message,
		Cause:   cause,
		Status:  http.StatusTooManyRequests,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(code, message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Code:    code,
		Message: message,
		Cause:   cause,
		Status:  http.StatusInternalServerError,
	}
}

// NewWebSocketError creates a new WebSocket error
func NewWebSocketError(code, message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeWebSocket,
		Code:    code,
		Message: message,
		Cause:   cause,
		Status:  http.StatusInternalServerError,
	}
}

// Helper functions

// AsAppError converts an error to an AppError if possible
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// WrapError wraps a standard error as an AppError
func WrapError(err error, errType ErrorType, code, message string) *AppError {
	return &AppError{
		Type:    errType,
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// Common messaging errors

// Channel errors
var (
	ErrChannelNotFound = NewNotFoundError(
		"channel_not_found",
		"The requested channel does not exist or you don't have access to it",
		nil,
	)

	ErrUserNotMember = NewAuthError(
		"user_not_member",
		"You are not a member of this channel",
		nil,
	)

	ErrUserAlreadyMember = NewValidationError(
		"user_already_member",
		"User is already a member of this channel",
		nil,
	)

	ErrCannotRemoveOwner = NewValidationError(
		"cannot_remove_owner",
		"Cannot remove the channel owner from the channel",
		nil,
	)
)

// Message errors
var (
	ErrMessageNotFound = NewNotFoundError(
		"message_not_found",
		"The requested message does not exist",
		nil,
	)

	ErrMessageTooLong = NewValidationError(
		"message_too_long",
		"Message content exceeds the maximum allowed length",
		nil,
	)

	ErrEmptyMessage = NewValidationError(
		"empty_message",
		"Message content cannot be empty",
		nil,
	)
)

// Attachment errors
var (
	ErrAttachmentNotFound = NewNotFoundError(
		"attachment_not_found",
		"The requested attachment does not exist",
		nil,
	)

	ErrAttachmentTooLarge = NewValidationError(
		"attachment_too_large",
		"Attachment size exceeds the maximum allowed limit",
		nil,
	)

	ErrInvalidAttachmentType = NewValidationError(
		"invalid_attachment_type",
		"The attachment type is not supported",
		nil,
	)
)

// Reaction errors
var (
	ErrReactionNotFound = NewNotFoundError(
		"reaction_not_found",
		"The requested reaction does not exist",
		nil,
	)

	ErrDuplicateReaction = NewValidationError(
		"duplicate_reaction",
		"User has already used this reaction on this message",
		nil,
	)
)

// WebSocket errors
var (
	ErrWebSocketConnLimit = NewRateLimitError(
		"ws_connection_limit",
		"Maximum number of WebSocket connections reached",
		nil,
	)

	ErrWebSocketRateLimit = NewRateLimitError(
		"ws_rate_limit",
		"Too many connection attempts. Please try again later",
		nil,
	)
)

// HandleError converts an error to an appropriate HTTP status code and response
func HandleError(err error) (int, interface{}) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.HTTPStatusCode(), appErr.JSON()
	}

	// Default to internal server error for unknown errors
	return http.StatusInternalServerError, map[string]interface{}{
		"error": "Internal server error",
		"type":  string(ErrorTypeInternal),
	}
}

var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrFileTooLarge     = errors.New("file size exceeds limit")
	ErrInvalidFileType  = errors.New("file type not allowed")
	ErrPinLimitExceeded = errors.New("channel pin limit exceeded")
)

// Common repository errors
var (
	// ErrNotFound is returned when an entity is not found
	ErrNotFound = errors.New("entity not found")

	// ErrInvalidInput is returned when the input is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized is returned when the user is not authorized
	ErrUnauthorized = errors.New("unauthorized")

	// ErrDuplicateEntry is returned when a duplicate entry is found
	ErrDuplicateEntry = errors.New("duplicate entry")

	// ErrInvalidOperation is returned when an invalid operation is attempted
	ErrInvalidOperation = errors.New("invalid operation")

	// ErrDatabaseOperation is returned when a database operation fails
	ErrDatabaseOperation = errors.New("database operation failed")
)

// ErrorCode represents a specific error code
type ErrorCode string

const (
	// Message errors
	ErrMessageNotFound     ErrorCode = "message_not_found"
	ErrMessageInvalid      ErrorCode = "message_invalid"
	ErrMessageUnauthorized ErrorCode = "message_unauthorized"
	ErrMessageDeleted      ErrorCode = "message_deleted"
	ErrMessageNotEditable  ErrorCode = "message_not_editable"
	ErrMessageNotDeletable ErrorCode = "message_not_deletable"

	// Channel errors
	ErrChannelNotFound     ErrorCode = "channel_not_found"
	ErrChannelInvalid      ErrorCode = "channel_invalid"
	ErrChannelUnauthorized ErrorCode = "channel_unauthorized"
	ErrChannelFull         ErrorCode = "channel_full"
	ErrChannelArchived     ErrorCode = "channel_archived"

	// User errors
	ErrUserNotFound     ErrorCode = "user_not_found"
	ErrUserInvalid      ErrorCode = "user_invalid"
	ErrUserUnauthorized ErrorCode = "user_unauthorized"
	ErrUserBlocked      ErrorCode = "user_blocked"
	ErrUserMuted        ErrorCode = "user_muted"

	// Attachment errors
	ErrAttachmentNotFound ErrorCode = "attachment_not_found"
	ErrAttachmentInvalid  ErrorCode = "attachment_invalid"
	ErrAttachmentTooLarge ErrorCode = "attachment_too_large"
	ErrAttachmentType     ErrorCode = "attachment_type_invalid"

	// Moderation errors
	ErrModerationRuleNotFound ErrorCode = "moderation_rule_not_found"
	ErrModerationRuleInvalid  ErrorCode = "moderation_rule_invalid"
	ErrModerationActionDenied ErrorCode = "moderation_action_denied"
	ErrModerationLevelInvalid ErrorCode = "moderation_level_invalid"

	// Word filter errors
	ErrWordFilterNotFound ErrorCode = "word_filter_not_found"
	ErrWordFilterInvalid  ErrorCode = "word_filter_invalid"
	ErrWordFilterExists   ErrorCode = "word_filter_exists"

	// System errors
	ErrInternalServer     ErrorCode = "internal_server_error"
	ErrDatabase           ErrorCode = "database_error"
	ErrStorage            ErrorCode = "storage_error"
	ErrValidation         ErrorCode = "validation_error"
	ErrRateLimit          ErrorCode = "rate_limit_exceeded"
	ErrServiceUnavailable ErrorCode = "service_unavailable"
)

// MessagingError represents a custom error in the messaging system
type MessagingError struct {
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	EntityID  uuid.UUID `json:"entity_id,omitempty"`
	Timestamp int64     `json:"timestamp"`
}

// NewMessagingError creates a new MessagingError
func NewMessagingError(code ErrorCode, message string, details string, entityID uuid.UUID) *MessagingError {
	return &MessagingError{
		Code:      code,
		Message:   message,
		Details:   details,
		EntityID:  entityID,
		Timestamp: time.Now().Unix(),
	}
}

// Error implements the error interface
func (e *MessagingError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (Details: %s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// HTTPStatus returns the appropriate HTTP status code for the error
func (e *MessagingError) HTTPStatus() int {
	switch e.Code {
	case ErrMessageNotFound, ErrChannelNotFound, ErrUserNotFound, ErrAttachmentNotFound,
		ErrModerationRuleNotFound, ErrWordFilterNotFound:
		return http.StatusNotFound
	case ErrMessageInvalid, ErrChannelInvalid, ErrUserInvalid, ErrAttachmentInvalid,
		ErrModerationRuleInvalid, ErrWordFilterInvalid, ErrValidation:
		return http.StatusBadRequest
	case ErrMessageUnauthorized, ErrChannelUnauthorized, ErrUserUnauthorized,
		ErrModerationActionDenied:
		return http.StatusForbidden
	case ErrMessageDeleted, ErrChannelArchived:
		return http.StatusGone
	case ErrAttachmentTooLarge:
		return http.StatusRequestEntityTooLarge
	case ErrRateLimit:
		return http.StatusTooManyRequests
	case ErrServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// IsNotFound checks if the error is a not found error
func (e *MessagingError) IsNotFound() bool {
	return e.HTTPStatus() == http.StatusNotFound
}

// IsInvalid checks if the error is an invalid request error
func (e *MessagingError) IsInvalid() bool {
	return e.HTTPStatus() == http.StatusBadRequest
}

// IsUnauthorized checks if the error is an unauthorized error
func (e *MessagingError) IsUnauthorized() bool {
	return e.HTTPStatus() == http.StatusForbidden
}

// Common error constructors
func NewMessageNotFoundError(messageID uuid.UUID) *MessagingError {
	return NewMessagingError(
		ErrMessageNotFound,
		"Message not found",
		"",
		messageID,
	)
}

func NewChannelNotFoundError(channelID uuid.UUID) *MessagingError {
	return NewMessagingError(
		ErrChannelNotFound,
		"Channel not found",
		"",
		channelID,
	)
}

func NewUserNotFoundError(userID uuid.UUID) *MessagingError {
	return NewMessagingError(
		ErrUserNotFound,
		"User not found",
		"",
		userID,
	)
}

func NewAttachmentNotFoundError(attachmentID uuid.UUID) *MessagingError {
	return NewMessagingError(
		ErrAttachmentNotFound,
		"Attachment not found",
		"",
		attachmentID,
	)
}

func NewModerationRuleNotFoundError(ruleID uuid.UUID) *MessagingError {
	return NewMessagingError(
		ErrModerationRuleNotFound,
		"Moderation rule not found",
		"",
		ruleID,
	)
}

func NewWordFilterNotFoundError(filterID uuid.UUID) *MessagingError {
	return NewMessagingError(
		ErrWordFilterNotFound,
		"Word filter not found",
		"",
		filterID,
	)
}

func NewValidationError(details string) *MessagingError {
	return NewMessagingError(
		ErrValidation,
		"Validation failed",
		details,
		uuid.Nil,
	)
}

func NewInternalServerError(details string) *MessagingError {
	return NewMessagingError(
		ErrInternalServer,
		"Internal server error",
		details,
		uuid.Nil,
	)
}
