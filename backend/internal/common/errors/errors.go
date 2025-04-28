package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Standard error types
var (
	ErrNotFound           = errors.New("resource not found")
	ErrBadRequest         = errors.New("bad request")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInternalServer     = errors.New("internal server error")
	ErrConflict           = errors.New("resource conflict")
	ErrValidation         = errors.New("validation error")
	ErrTooManyRequests    = errors.New("too many requests")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrTimeout            = errors.New("request timeout")
)

// AppError represents an application error with context
type AppError struct {
	// Original is the underlying error
	Original error

	// Message is a user-friendly message
	Message string

	// Code is an application-specific error code
	Code string

	// StatusCode is the HTTP status code
	StatusCode int

	// Data contains additional error context
	Data map[string]interface{}

	// Internal indicates if this error should be shown to users
	Internal bool
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Original != nil {
		return e.Original.Error()
	}
	return "an error occurred"
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Original
}

// Is implements errors.Is interface
func (e *AppError) Is(target error) bool {
	if target == nil || e == nil {
		return target == e
	}

	// Check if target is also an AppError
	targetAppErr, ok := target.(*AppError)
	if ok {
		return e.Original == targetAppErr.Original
	}

	// Check if target matches the original error
	return errors.Is(e.Original, target)
}

// As implements errors.As interface
func (e *AppError) As(target interface{}) bool {
	if target == nil {
		return false
	}

	// Try to cast to AppError
	targetPtr, ok := target.(**AppError)
	if ok {
		*targetPtr = e
		return true
	}

	// If target isn't *AppError, try the original error
	return errors.As(e.Original, target)
}

// WithMessage adds a message to the error
func (e *AppError) WithMessage(message string) *AppError {
	e.Message = message
	return e
}

// WithCode adds a code to the error
func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

// WithData adds context data to the error
func (e *AppError) WithData(key string, value interface{}) *AppError {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}
	e.Data[key] = value
	return e
}

// WithStatusCode sets the HTTP status code
func (e *AppError) WithStatusCode(statusCode int) *AppError {
	e.StatusCode = statusCode
	return e
}

// IsInternal marks the error as internal (not to be exposed to users)
func (e *AppError) IsInternal() *AppError {
	e.Internal = true
	return e
}

// IsPublic marks the error as public (can be exposed to users)
func (e *AppError) IsPublic() *AppError {
	e.Internal = false
	return e
}

// NewAppError creates a new AppError
func NewAppError(err error, message string) *AppError {
	return &AppError{
		Original:   err,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Internal:   true,
	}
}

// New creates a new error
func New(message string) error {
	return errors.New(message)
}

// Wrap wraps an error with a message
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf wraps an error with a formatted message
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// NotFound creates a not found error
func NotFound(message string) *AppError {
	return &AppError{
		Original:   ErrNotFound,
		Message:    message,
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Internal:   false,
	}
}

// BadRequest creates a bad request error
func BadRequest(message string) *AppError {
	return &AppError{
		Original:   ErrBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Internal:   false,
	}
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *AppError {
	return &AppError{
		Original:   ErrUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Internal:   false,
	}
}

// Forbidden creates a forbidden error
func Forbidden(message string) *AppError {
	return &AppError{
		Original:   ErrForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Internal:   false,
	}
}

// InternalServer creates an internal server error
func InternalServer(err error) *AppError {
	return &AppError{
		Original:   err,
		Message:    "An internal server error occurred",
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_SERVER_ERROR",
		Internal:   true,
	}
}

// Conflict creates a conflict error
func Conflict(message string) *AppError {
	return &AppError{
		Original:   ErrConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
		Code:       "CONFLICT",
		Internal:   false,
	}
}

// Validation creates a validation error
func Validation(message string) *AppError {
	return &AppError{
		Original:   ErrValidation,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Internal:   false,
	}
}

// ValidationWithData creates a validation error with data
func ValidationWithData(fieldErrors map[string]string) *AppError {
	data := make(map[string]interface{})
	data["fields"] = fieldErrors

	return &AppError{
		Original:   ErrValidation,
		Message:    "Validation failed",
		StatusCode: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Internal:   false,
		Data:       data,
	}
}

// Is checks if an error is of a specific type
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// GetStatusCode returns the HTTP status code for an error
func GetStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode
	}

	// Map standard errors to status codes
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrTooManyRequests):
		return http.StatusTooManyRequests
	case errors.Is(err, ErrServiceUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, ErrTimeout):
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}
