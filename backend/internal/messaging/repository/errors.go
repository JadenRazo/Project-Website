package repository

import "errors"

// Common repository errors
var (
	// Channel errors
	ErrChannelNotFound   = errors.New("channel not found")
	ErrUserAlreadyMember = errors.New("user is already a member of this channel")
	ErrUserNotMember     = errors.New("user is not a member of this channel")
	ErrCannotRemoveOwner = errors.New("cannot remove channel owner from channel")

	// Message errors
	ErrMessageNotFound    = errors.New("message not found")
	ErrDuplicateReaction  = errors.New("user has already used this reaction on this message")
	ErrReactionNotFound   = errors.New("reaction not found")
	ErrUnauthorizedAccess = errors.New("unauthorized access to this resource")
)
