package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
	"gorm.io/gorm"
)

// ReadReceiptRepo implements repository.ReadReceiptRepository using PostgreSQL
type ReadReceiptRepo struct {
	db *gorm.DB
}

// NewReadReceiptRepository creates a new PostgreSQL read receipt repository
func NewReadReceiptRepository(db *gorm.DB) repository.ReadReceiptRepository {
	return &ReadReceiptRepo{
		db: db,
	}
}

// Create implements repository.BaseRepository
func (r *ReadReceiptRepo) Create(ctx context.Context, receipt *domain.MessagingReadReceipt) error {
	if receipt.CreatedAt.IsZero() {
		receipt.CreatedAt = time.Now()
	}
	if receipt.UpdatedAt.IsZero() {
		receipt.UpdatedAt = time.Now()
	}
	return r.db.WithContext(ctx).Create(receipt).Error
}

// FindByID implements repository.BaseRepository
func (r *ReadReceiptRepo) FindByID(ctx context.Context, id uint) (*domain.MessagingReadReceipt, error) {
	var receipt domain.MessagingReadReceipt
	err := r.db.WithContext(ctx).
		Preload("Message").
		Preload("User").
		First(&receipt, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &receipt, nil
}

// Update implements repository.BaseRepository
func (r *ReadReceiptRepo) Update(ctx context.Context, receipt *domain.MessagingReadReceipt) error {
	receipt.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(receipt).Error
}

// Delete implements repository.BaseRepository
func (r *ReadReceiptRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.MessagingReadReceipt{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// FindAll implements repository.BaseRepository
func (r *ReadReceiptRepo) FindAll(ctx context.Context) ([]domain.MessagingReadReceipt, error) {
	var receipts []domain.MessagingReadReceipt
	err := r.db.WithContext(ctx).
		Preload("Message").
		Preload("User").
		Find(&receipts).Error
	return receipts, err
}

// GetMessageReadReceipts retrieves all read receipts for a message
func (r *ReadReceiptRepo) GetMessageReadReceipts(ctx context.Context, messageID uint) ([]domain.MessagingReadReceipt, error) {
	var receipts []domain.MessagingReadReceipt
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("message_id = ?", messageID).
		Find(&receipts).Error
	return receipts, err
}

// GetUserReadReceipts retrieves all read receipts for a user
func (r *ReadReceiptRepo) GetUserReadReceipts(ctx context.Context, userID uint) ([]domain.MessagingReadReceipt, error) {
	var receipts []domain.MessagingReadReceipt
	err := r.db.WithContext(ctx).
		Preload("Message").
		Where("user_id = ?", userID).
		Find(&receipts).Error
	return receipts, err
}

// GetUnreadCount gets the count of unread messages for a user in a channel
func (r *ReadReceiptRepo) GetUnreadCount(ctx context.Context, channelID uint, userID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.MessagingMessage{}).
		Where("channel_id = ? AND sender_id != ?", channelID, userID).
		Where("NOT EXISTS (SELECT 1 FROM messaging_read_receipts WHERE message_id = messaging_messages.id AND user_id = ?)", userID).
		Count(&count).Error
	return int(count), err
}

// MarkAsRead marks a message as read by a user
func (r *ReadReceiptRepo) MarkAsRead(ctx context.Context, messageID uint, userID uint) error {
	receipt := &domain.MessagingReadReceipt{
		MessageID: messageID,
		UserID:    userID,
		ReadAt:    time.Now(),
	}
	return r.db.WithContext(ctx).Create(receipt).Error
}

// GetLastReadTime gets the last time a user read messages in a channel
func (r *ReadReceiptRepo) GetLastReadTime(ctx context.Context, channelID uint, userID uint) (*time.Time, error) {
	var receipt domain.MessagingReadReceipt
	err := r.db.WithContext(ctx).
		Joins("JOIN messaging_messages ON messaging_messages.id = messaging_read_receipts.message_id").
		Where("messaging_messages.channel_id = ? AND messaging_read_receipts.user_id = ?", channelID, userID).
		Order("messaging_read_receipts.read_at DESC").
		First(&receipt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &receipt.ReadAt, nil
}
