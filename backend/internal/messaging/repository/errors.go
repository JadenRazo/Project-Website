package repository

import "errors"

// Common repository errors
var (
	// ErrNotFound is returned when a record is not found
	ErrNotFound = errors.New("record not found")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized is returned when the user is not authorized to perform an action
	ErrUnauthorized = errors.New("unauthorized")

	// ErrDuplicateEntry is returned when trying to create a record that already exists
	ErrDuplicateEntry = errors.New("duplicate entry")

	// ErrInvalidOperation is returned when trying to perform an invalid operation
	ErrInvalidOperation = errors.New("invalid operation")

	// ErrDatabaseOperation is returned when a database operation fails
	ErrDatabaseOperation = errors.New("database operation failed")

	// ErrMessageNotFound is returned when a message is not found
	ErrMessageNotFound = errors.New("message not found")

	// ErrChannelNotFound is returned when a channel is not found
	ErrChannelNotFound = errors.New("channel not found")

	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrAttachmentNotFound is returned when an attachment is not found
	ErrAttachmentNotFound = errors.New("attachment not found")

	// ErrEmbedNotFound is returned when an embed is not found
	ErrEmbedNotFound = errors.New("embed not found")

	// ErrReactionNotFound is returned when a reaction is not found
	ErrReactionNotFound = errors.New("reaction not found")

	// ErrReadReceiptNotFound is returned when a read receipt is not found
	ErrReadReceiptNotFound = errors.New("read receipt not found")

	// ErrDuplicateReaction is returned when trying to add a duplicate reaction
	ErrDuplicateReaction = errors.New("duplicate reaction")

	// ErrDuplicateReadReceipt is returned when trying to add a duplicate read receipt
	ErrDuplicateReadReceipt = errors.New("duplicate read receipt")

	// ErrUserAlreadyMember is returned when trying to add a user that is already a member
	ErrUserAlreadyMember = errors.New("user already member")

	// ErrUserNotMember is returned when trying to remove a user that is not a member
	ErrUserNotMember = errors.New("user not member")

	// ErrCannotRemoveOwner is returned when trying to remove a channel owner
	ErrCannotRemoveOwner = errors.New("cannot remove owner")

	// ErrPinLimitExceeded is returned when trying to pin more messages than allowed
	ErrPinLimitExceeded = errors.New("pin limit exceeded")

	// ErrFileTooLarge is returned when trying to upload a file that exceeds the size limit
	ErrFileTooLarge = errors.New("file too large")

	// ErrInvalidFileType is returned when trying to upload a file with an invalid type
	ErrInvalidFileType = errors.New("invalid file type")

	// ErrInvalidURL is returned when trying to create an embed with an invalid URL
	ErrInvalidURL = errors.New("invalid url")

	// ErrInvalidEmojiCode is returned when trying to add a reaction with an invalid emoji code
	ErrInvalidEmojiCode = errors.New("invalid emoji code")
)
