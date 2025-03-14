package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/core/repository"
	"github.com/JadenRazo/Project-Website/backend/internal/domain"
)

// MessageGormRepository implements MessageRepository using Gorm
type MessageGormRepository struct {
	db       *gorm.DB
	baseRepo *repository.BaseRepository[domain.Message]
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) *MessageGormRepository {
	baseRepo := repository.NewBaseRepository(db, domain.Message{})
	return &MessageGormRepository{
		db: db,
		baseRepo: baseRepo.WithPreloads(
			func(db *gorm.DB) *gorm.DB { return db.Preload("Attachments") },
			func(db *gorm.DB) *gorm.DB { return db.Preload("Reactions") },
			func(db *gorm.DB) *gorm.DB {
				return db.Preload("Reactions.User", func(db *gorm.DB) *gorm.DB {
					return db.Select("id, username, avatar")
				})
			},
			func(db *gorm.DB) *gorm.DB {
				return db.Preload("Sender", func(db *gorm.DB) *gorm.DB {
					return db.Select("id, username, first_name, last_name, avatar, status")
				})
			},
		),
	}
}

// CreateMessage stores a new message
func (r *MessageGormRepository) CreateMessage(ctx context.Context, message *domain.Message) error {
	return r.baseRepo.Create(ctx, message)
}

// GetMessage retrieves a message by ID
func (r *MessageGormRepository) GetMessage(ctx context.Context, messageID uint) (*domain.Message, error) {
	return r.baseRepo.FindByID(ctx, messageID)
}

// GetChannelMessages retrieves messages from a channel with pagination
func (r *MessageGormRepository) GetChannelMessages(ctx context.Context, channelID uint, lastID uint, limit int) ([]domain.Message, error) {
	if channelID == 0 {
		return nil, repository.ErrInvalidInput
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}

	var messages []domain.Message
	query := r.db.WithContext(ctx).
		Preload("Attachments").
		Preload("Reactions").
		Preload("Reactions.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, avatar")
		}).
		Preload("Sender", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, first_name, last_name, avatar, status")
		}).
		Where("channel_id = ?", channelID).
		Order("created_at DESC").
		Limit(limit)

	// If lastID is provided, paginate from there
	if lastID > 0 {
		query = query.Where("id < ?", lastID)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

// UpdateMessage updates a message
func (r *MessageGormRepository) UpdateMessage(ctx context.Context, message *domain.Message) error {
	if message == nil || message.ID == 0 || message.SenderID == 0 {
		return repository.ErrInvalidInput
	}

	return r.db.WithContext(ctx).Model(message).
		Select("content", "is_edited", "edited_at").
		Updates(message).Error
}

// DeleteMessage soft-deletes a message
func (r *MessageGormRepository) DeleteMessage(ctx context.Context, messageID uint, userID uint) error {
	if messageID == 0 || userID == 0 {
		return repository.ErrInvalidInput
	}

	result := r.db.WithContext(ctx).
		Where("id = ? AND sender_id = ?", messageID, userID).
		Delete(&domain.Message{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repository.ErrUnauthorized
	}

	return nil
}

// AddReaction adds a reaction to a message
func (r *MessageGormRepository) AddReaction(ctx context.Context, reaction *domain.Reaction) error {
	if reaction == nil || reaction.MessageID == 0 || reaction.UserID == 0 || reaction.EmojiCode == "" {
		return repository.ErrInvalidInput
	}

	// Check if the same reaction already exists
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Reaction{}).
		Where("message_id = ? AND user_id = ? AND emoji_code = ?",
			reaction.MessageID, reaction.UserID, reaction.EmojiCode).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Reaction already exists
	}

	return r.db.WithContext(ctx).Create(reaction).Error
}

// RemoveReaction removes a reaction from a message
func (r *MessageGormRepository) RemoveReaction(ctx context.Context, reactionID uint, userID uint) error {
	if reactionID == 0 || userID == 0 {
		return repository.ErrInvalidInput
	}

	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", reactionID, userID).
		Delete(&domain.Reaction{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repository.ErrUnauthorized
	}

	return nil
}

// GetMessageReactions gets all reactions for a message
func (r *MessageGormRepository) GetMessageReactions(ctx context.Context, messageID uint) ([]domain.Reaction, error) {
	if messageID == 0 {
		return nil, repository.ErrInvalidInput
	}

	var reactions []domain.Reaction
	err := r.db.WithContext(ctx).
		Where("message_id = ?", messageID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, avatar")
		}).
		Find(&reactions).Error

	if err != nil {
		return nil, err
	}

	return reactions, nil
}

// MarkMessageAsRead marks a message as read by a user
func (r *MessageGormRepository) MarkMessageAsRead(ctx context.Context, messageID uint, userID uint) error {
	if messageID == 0 || userID == 0 {
		return repository.ErrInvalidInput
	}

	// Fetch current message
	var message domain.Message
	if err := r.db.WithContext(ctx).Select("id, read_by").First(&message, messageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrNotFound
		}
		return err
	}

	// Parse read_by JSON
	var readByIDs []uint
	if message.ReadBy != "" {
		if err := json.Unmarshal([]byte(message.ReadBy), &readByIDs); err != nil {
			// If invalid JSON, start fresh
			readByIDs = []uint{}
		}
	}

	// Check if user already marked as read
	for _, id := range readByIDs {
		if id == userID {
			return nil // Already marked as read
		}
	}

	// Add user to read_by
	readByIDs = append(readByIDs, userID)
	readByJSON, err := json.Marshal(readByIDs)
	if err != nil {
		return err
	}

	// Update message
	return r.db.WithContext(ctx).Model(&domain.Message{}).
		Where("id = ?", messageID).
		Update("read_by", string(readByJSON)).Error
}

// GetUnreadCount gets count of unread messages for a user in a channel
func (r *MessageGormRepository) GetUnreadCount(ctx context.Context, channelID uint, userID uint) (int64, error) {
	if channelID == 0 || userID == 0 {
		return 0, repository.ErrInvalidInput
	}

	var count int64

	// Count messages where read_by doesn't contain userID
	// This is a simplistic implementation - in a production system, you might
	// want to use a more efficient approach like a separate table for message read status
	err := r.db.WithContext(ctx).Model(&domain.Message{}).
		Where("channel_id = ?", channelID).
		Where("sender_id != ?", userID). // Don't count user's own messages
		Where("read_by NOT LIKE ?", "%"+strconv.FormatUint(uint64(userID), 10)+"%").
		Count(&count).Error

	return count, err
}

// ChannelGormRepository implements ChannelRepository using Gorm
type ChannelGormRepository struct {
	db       *gorm.DB
	baseRepo *repository.BaseRepository[domain.Channel]
}

// NewChannelRepository creates a new channel repository
func NewChannelRepository(db *gorm.DB) *ChannelGormRepository {
	baseRepo := repository.NewBaseRepository(db, domain.Channel{})
	return &ChannelGormRepository{
		db: db,
		baseRepo: baseRepo.WithPreloads(
			func(db *gorm.DB) *gorm.DB { return db.Preload("Owner", "id, username, avatar") },
			func(db *gorm.DB) *gorm.DB {
				return db.Preload("Members", "id, username, first_name, last_name, avatar, status")
			},
		),
	}
}

// CreateChannel creates a new channel
func (r *ChannelGormRepository) CreateChannel(ctx context.Context, channel *domain.Channel) error {
	if channel == nil || channel.Name == "" || channel.OwnerID == 0 {
		return repository.ErrInvalidInput
	}

	return r.baseRepo.Create(ctx, channel)
}

// GetChannel retrieves a channel by ID
func (r *ChannelGormRepository) GetChannel(ctx context.Context, channelID uint) (*domain.Channel, error) {
	return r.baseRepo.FindByID(ctx, channelID)
}

// GetUserChannels gets all channels a user is a member of
func (r *ChannelGormRepository) GetUserChannels(ctx context.Context, userID uint) ([]domain.Channel, error) {
	if userID == 0 {
		return nil, repository.ErrInvalidInput
	}

	var channels []domain.Channel
	err := r.db.WithContext(ctx).Where("owner_id = ? OR id IN (SELECT channel_id FROM user_channels WHERE user_id = ?)", userID, userID).Find(&channels).Error
	return channels, err
}
