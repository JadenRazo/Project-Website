package repository

import (
	"context"
	"time"
)

// ReadReceipt represents a user's read status for a message
type ReadReceipt struct {
	ID        uint      `json:"id"`
	MessageID uint      `json:"message_id"`
	UserID    uint      `json:"user_id"`
	ReadAt    time.Time `json:"read_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReadReceiptRepository defines the interface for read receipt storage operations
type ReadReceiptRepository interface {
	// CreateReadReceipt creates a new read receipt
	CreateReadReceipt(ctx context.Context, receipt *ReadReceipt) error

	// CreateBulkReadReceipts creates multiple read receipts in a single operation
	CreateBulkReadReceipts(ctx context.Context, receipts []ReadReceipt) error

	// ReadReceiptExists checks if a read receipt already exists for a message and user
	ReadReceiptExists(ctx context.Context, messageID, userID uint) (bool, error)

	// GetMessageReadReceipts gets all read receipts for a message
	GetMessageReadReceipts(ctx context.Context, messageID uint) ([]ReadReceipt, error)

	// GetUnreadCount gets the count of unread messages for a user in a channel
	GetUnreadCount(ctx context.Context, channelID, userID uint) (int, error)

	// DeleteReadReceipts deletes all read receipts for a message
	DeleteReadReceipts(ctx context.Context, messageID uint) error
}
