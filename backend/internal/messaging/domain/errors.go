package domain

import "errors"

// Error types for the moderation system
var (
	ErrValidation     = errors.New("validation error")
	ErrNotFound       = errors.New("not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInternal       = errors.New("internal server error")
	ErrDuplicateEntry = errors.New("duplicate entry")
)

// NewValidationError creates a new validation error
func NewValidationError(message string) error {
	return errors.Join(ErrValidation, errors.New(message))
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidation)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) error {
	return errors.Join(ErrNotFound, errors.New(message))
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) error {
	return errors.Join(ErrUnauthorized, errors.New(message))
}

// IsUnauthorizedError checks if an error is an unauthorized error
func IsUnauthorizedError(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// NewInternalError creates a new internal server error
func NewInternalError(message string) error {
	return errors.Join(ErrInternal, errors.New(message))
}

// IsInternalError checks if an error is an internal server error
func IsInternalError(err error) bool {
	return errors.Is(err, ErrInternal)
}

// NewDuplicateEntryError creates a new duplicate entry error
func NewDuplicateEntryError(message string) error {
	return errors.Join(ErrDuplicateEntry, errors.New(message))
}

// IsDuplicateEntryError checks if an error is a duplicate entry error
func IsDuplicateEntryError(err error) bool {
	return errors.Is(err, ErrDuplicateEntry)
}
