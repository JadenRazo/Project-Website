package postgres

import (
	"context"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
	"gorm.io/gorm"
)

// ReadReceiptRepo implements the ReadReceiptRepository interface
type ReadReceiptRepo struct {
	db *gorm.DB
}

// NewReadReceiptRepository creates a new PostgreSQL read receipt repository
func NewReadReceiptRepository(db *gorm.DB) repository.ReadReceiptRepository {
	return &ReadReceiptRepo{
		db: db,
	}
}

// CreateReadReceipt creates a new read receipt
func (r *ReadReceiptRepo) CreateReadReceipt(ctx context.Context, receipt *repository.ReadReceipt) error {
	return r.db.WithContext(ctx).Create(receipt).Error
}

// CreateBulkReadReceipts creates multiple read receipts in a single operation
func (r *ReadReceiptRepo) CreateBulkReadReceipts(ctx context.Context, receipts []repository.ReadReceipt) error {
	// Use a transaction for batch insert
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(&receipts).Error
	})
}

// ReadReceiptExists checks if a read receipt already exists for a message and user
func (r *ReadReceiptRepo) ReadReceiptExists(ctx context.Context, messageID, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&repository.ReadReceipt{}).
		Where("message_id = ? AND user_id = ?", messageID, userID).
		Count(&count).
		Error

	return count > 0, err
}

// GetMessageReadReceipts gets all read receipts for a message
func (r *ReadReceiptRepo) GetMessageReadReceipts(ctx context.Context, messageID uint) ([]repository.ReadReceipt, error) {
	var receipts []repository.ReadReceipt
	err := r.db.WithContext(ctx).
		Where("message_id = ?", messageID).
		Find(&receipts).
		Error

	return receipts, err
}

// GetUnreadCount gets the count of unread messages for a user in a channel
func (r *ReadReceiptRepo) GetUnreadCount(ctx context.Context, channelID, userID uint) (int, error) {
	// This query counts messages that:
	// 1. Are in the specified channel
	// 2. Weren't sent by the current user
	// 3. Don't have a read receipt from the current user
	var count int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM messages m
		WHERE m.channel_id = ? 
		AND m.sender_id != ?
		AND NOT EXISTS (
			SELECT 1 FROM read_receipts rr 
			WHERE rr.message_id = m.id 
			AND rr.user_id = ?
		)
	`, channelID, userID, userID).Count(&count).Error

	return int(count), err
}

// DeleteReadReceipts deletes all read receipts for a message
func (r *ReadReceiptRepo) DeleteReadReceipts(ctx context.Context, messageID uint) error {
	return r.db.WithContext(ctx).
		Where("message_id = ?", messageID).
		Delete(&repository.ReadReceipt{}).
		Error
}
