package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
