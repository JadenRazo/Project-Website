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

// CreateMessage creates a new message
func (r *MessageRepo) CreateMessage(ctx context.Context, message *domain.Message) error {
	// Set creation time if not already set
	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now()
	}

	return r.db.WithContext(ctx).Create(message).Error
}

// GetMessage retrieves a message by ID
func (r *MessageRepo) GetMessage(ctx context.Context, messageID uint) (*domain.Message, error) {
	var message domain.Message
	result := r.db.WithContext(ctx).
		Preload("Attachments"). // Load associated attachments
		First(&message, messageID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrMessageNotFound
		}
		return nil, result.Error
	}
	return &message, nil
}

// GetChannelMessages retrieves messages from a channel with pagination
func (r *MessageRepo) GetChannelMessages(ctx context.Context, channelID uint, lastMessageID uint, limit int) ([]domain.Message, error) {
	var messages []domain.Message
	query := r.db.WithContext(ctx).
		Preload("Attachments").
		Where("channel_id = ?", channelID).
		Order("created_at DESC")

	// If lastMessageID is provided, paginate from that message
	if lastMessageID > 0 {
		var lastMessage domain.Message
		if err := r.db.First(&lastMessage, lastMessageID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, repository.ErrMessageNotFound
			}
			return nil, err
		}
		query = query.Where("created_at < ?", lastMessage.CreatedAt)
	}

	// Apply limit
	if limit <= 0 {
		limit = 50 // Default limit
	}
	query = query.Limit(limit)

	// Execute query
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

// UpdateMessage updates a message
func (r *MessageRepo) UpdateMessage(ctx context.Context, message *domain.Message) error {
	// Get existing message
	var existing domain.Message
	if err := r.db.WithContext(ctx).First(&existing, message.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrMessageNotFound
		}
		return err
	}

	// Ensure we don't update certain fields
	message.CreatedAt = existing.CreatedAt
	message.SenderID = existing.SenderID
	message.ChannelID = existing.ChannelID

	// Mark as edited
	message.IsEdited = true
	message.EditedAt = time.Now()

	// Update only content and edit metadata
	return r.db.WithContext(ctx).Model(message).
		Select("content", "is_edited", "edited_at").
		Updates(message).Error
}

// DeleteMessage deletes a message
func (r *MessageRepo) DeleteMessage(ctx context.Context, messageID uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Message{}, messageID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrMessageNotFound
	}
	return nil
}

// AddReaction adds a reaction to a message
func (r *MessageRepo) AddReaction(ctx context.Context, reaction *domain.Reaction) error {
	// Check if the message exists
	var count int64
	r.db.WithContext(ctx).Model(&domain.Message{}).Where("id = ?", reaction.MessageID).Count(&count)
	if count == 0 {
		return repository.ErrMessageNotFound
	}

	// Check for duplicate reaction
	r.db.WithContext(ctx).Model(&domain.Reaction{}).
		Where("message_id = ? AND user_id = ? AND emoji_code = ?",
			reaction.MessageID, reaction.UserID, reaction.EmojiCode).
		Count(&count)

	if count > 0 {
		return repository.ErrDuplicateReaction
	}

	// Set creation time
	if reaction.CreatedAt.IsZero() {
		reaction.CreatedAt = time.Now()
	}

	return r.db.WithContext(ctx).Create(reaction).Error
}

// RemoveReaction removes a reaction from a message
func (r *MessageRepo) RemoveReaction(ctx context.Context, reactionID uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Reaction{}, reactionID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrReactionNotFound
	}
	return nil
}

// GetMessageReactions retrieves all reactions for a message
func (r *MessageRepo) GetMessageReactions(ctx context.Context, messageID uint) ([]domain.Reaction, error) {
	var reactions []domain.Reaction
	if err := r.db.WithContext(ctx).
		Where("message_id = ?", messageID).
		Find(&reactions).Error; err != nil {
		return nil, err
	}
	return reactions, nil
}

// SearchMessages searches for messages in a channel by content
func (r *MessageRepo) SearchMessages(ctx context.Context, channelID uint, query string, limit int) ([]domain.Message, error) {
	var messages []domain.Message

	if limit <= 0 {
		limit = 50 // Default limit
	}

	// Full-text search - this assumes you've set up full-text search in PostgreSQL
	err := r.db.WithContext(ctx).
		Where("channel_id = ? AND to_tsvector('english', content) @@ plainto_tsquery('english', ?)",
			channelID, query).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}
	return messages, nil
}

// GetUnreadMessages retrieves unread messages for a user in a channel before the specified time
func (r *MessageRepo) GetUnreadMessages(
	ctx context.Context,
	channelID uint,
	userID uint,
	upToTime time.Time,
) ([]domain.Message, error) {
	var messages []domain.Message

	// Query for messages in the channel that:
	// 1. Were created before upToTime
	// 2. Weren't sent by the current user
	// 3. Don't have a read receipt from the current user
	err := r.db.WithContext(ctx).
		Preload("Attachments").
		Table("messages m").
		Where("m.channel_id = ?", channelID).
		Where("m.sender_id != ?", userID).
		Where("m.created_at <= ?", upToTime).
		Where("NOT EXISTS (SELECT 1 FROM read_receipts rr WHERE rr.message_id = m.id AND rr.user_id = ?)", userID).
		Order("m.created_at ASC").
		Find(&messages).
		Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}
