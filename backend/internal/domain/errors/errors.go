package errors

import "errors"

// Common domain errors
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
	ErrDuplicate    = errors.New("duplicate entry")
)
