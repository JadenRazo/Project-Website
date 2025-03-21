package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
	"gorm.io/gorm"
)

// MessageRepo implements repository.MessageRepository using PostgreSQL
type MessageRepo struct {
	db *gorm.DB
}

// NewMessageRepository creates a new PostgreSQL message repository
func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &MessageRepo{
		db: db,
	}
}

// Create implements repository.BaseRepository
func (r *MessageRepo) Create(ctx context.Context, message *domain.MessagingMessage) error {
	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now()
	}
	if message.UpdatedAt.IsZero() {
		message.UpdatedAt = time.Now()
	}
	return r.db.WithContext(ctx).Create(message).Error
}

// FindByID implements repository.BaseRepository
func (r *MessageRepo) FindByID(ctx context.Context, id uint) (*domain.MessagingMessage, error) {
	var message domain.MessagingMessage
	err := r.db.WithContext(ctx).
		Preload("Attachments").
		Preload("Embeds").
		Preload("Reactions").
		Preload("Sender").
		First(&message, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &message, nil
}

// Update implements repository.BaseRepository
func (r *MessageRepo) Update(ctx context.Context, message *domain.MessagingMessage) error {
	message.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(message).Error
}

// Delete implements repository.BaseRepository
func (r *MessageRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.MessagingMessage{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// FindAll implements repository.BaseRepository
func (r *MessageRepo) FindAll(ctx context.Context) ([]domain.MessagingMessage, error) {
	var messages []domain.MessagingMessage
	err := r.db.WithContext(ctx).
		Preload("Attachments").
		Preload("Embeds").
		Preload("Reactions").
		Preload("Sender").
		Find(&messages).Error
	return messages, err
}

// GetChannelMessages retrieves messages from a channel with pagination
func (r *MessageRepo) GetChannelMessages(ctx context.Context, channelID uint, lastMessageID uint, limit int) ([]domain.MessagingMessage, error) {
	var messages []domain.MessagingMessage
	query := r.db.WithContext(ctx).
		Preload("Attachments").
		Preload("Embeds").
		Preload("Reactions").
		Preload("Sender").
		Where("channel_id = ?", channelID).
		Order("created_at DESC")

	if lastMessageID > 0 {
		query = query.Where("id < ?", lastMessageID)
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}
	query = query.Limit(limit)

	err := query.Find(&messages).Error
	return messages, err
}

// GetThreadMessages retrieves messages in a thread
func (r *MessageRepo) GetThreadMessages(ctx context.Context, threadID uint, lastMessageID uint, limit int) ([]domain.MessagingMessage, error) {
	var messages []domain.MessagingMessage
	query := r.db.WithContext(ctx).
		Preload("Attachments").
		Preload("Embeds").
		Preload("Reactions").
		Preload("Sender").
		Where("thread_id = ?", threadID).
		Order("created_at ASC")

	if lastMessageID > 0 {
		query = query.Where("id > ?", lastMessageID)
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}
	query = query.Limit(limit)

	err := query.Find(&messages).Error
	return messages, err
}

// SearchMessages searches for messages with given criteria
func (r *MessageRepo) SearchMessages(ctx context.Context, filters repository.MessageSearchFilters) ([]domain.MessagingMessage, int, error) {
	var messages []domain.MessagingMessage
	var total int64

	query := r.db.WithContext(ctx).
		Preload("Attachments").
		Preload("Embeds").
		Preload("Reactions").
		Preload("Sender")

	// Apply filters
	if filters.ChannelID != nil {
		query = query.Where("channel_id = ?", *filters.ChannelID)
	}
	if filters.UserID != nil {
		query = query.Where("sender_id = ?", *filters.UserID)
	}
	if filters.ThreadID != nil {
		query = query.Where("thread_id = ?", *filters.ThreadID)
	}
	if filters.Query != "" {
		query = query.Where("content ILIKE ?", "%"+filters.Query+"%")
	}
	if filters.StartTime != nil {
		query = query.Where("created_at >= ?", *filters.StartTime)
	}
	if filters.EndTime != nil {
		query = query.Where("created_at <= ?", *filters.EndTime)
	}
	if filters.HasAttachments {
		query = query.Joins("JOIN messaging_attachments ON messaging_attachments.message_id = messaging_messages.id")
	}
	if filters.HasMentions {
		query = query.Where("mentions IS NOT NULL AND array_length(mentions, 1) > 0")
	}
	if filters.IsPinned {
		query = query.Where("is_pinned = ?", true)
	}
	if filters.IsNSFW {
		query = query.Where("is_nsfw = ?", true)
	}
	if filters.IsSpoiler {
		query = query.Where("is_spoiler = ?", true)
	}

	// Get total count
	err := query.Model(&domain.MessagingMessage{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filters.Limit <= 0 {
		filters.Limit = 50 // Default limit
	}
	query = query.Offset(filters.Offset).Limit(filters.Limit)

	// Execute query
	err = query.Find(&messages).Error
	return messages, int(total), err
}

// AddReaction adds a reaction to a message
func (r *MessageRepo) AddReaction(ctx context.Context, reaction *domain.MessagingReaction) error {
	reaction.CreatedAt = time.Now()
	reaction.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(reaction).Error
}

// RemoveReaction removes a reaction from a message
func (r *MessageRepo) RemoveReaction(ctx context.Context, reactionID uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.MessagingReaction{}, reactionID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// GetMessageReactions retrieves all reactions on a message
func (r *MessageRepo) GetMessageReactions(ctx context.Context, messageID uint) ([]domain.MessagingReaction, error) {
	var reactions []domain.MessagingReaction
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("message_id = ?", messageID).
		Find(&reactions).Error
	return reactions, err
}

// MarkAsRead marks a message as read by a user
func (r *MessageRepo) MarkAsRead(ctx context.Context, messageID uint, userID uint) error {
	receipt := &domain.MessagingReadReceipt{
		MessageID: messageID,
		UserID:    userID,
		ReadAt:    time.Now(),
	}
	return r.db.WithContext(ctx).Create(receipt).Error
}

// GetUnreadCount gets the count of unread messages for a user in a channel
func (r *MessageRepo) GetUnreadCount(ctx context.Context, channelID uint, userID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.MessagingMessage{}).
		Where("channel_id = ? AND sender_id != ?", channelID, userID).
		Where("NOT EXISTS (SELECT 1 FROM messaging_read_receipts WHERE message_id = messaging_messages.id AND user_id = ?)", userID).
		Count(&count).Error
	return int(count), err
}
