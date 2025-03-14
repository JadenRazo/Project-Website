package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
	"gorm.io/gorm"
)

// ChannelRepo implements repository.ChannelRepository using PostgreSQL
type ChannelRepo struct {
	db *gorm.DB
}

// NewChannelRepository creates a new PostgreSQL channel repository
func NewChannelRepository(db *gorm.DB) repository.ChannelRepository {
	return &ChannelRepo{
		db: db,
	}
}

// CreateChannel creates a new channel
func (r *ChannelRepo) CreateChannel(ctx context.Context, channel *domain.Channel) error {
	// Set creation time if not already set
	if channel.CreatedAt.IsZero() {
		channel.CreatedAt = time.Now()
	}

	result := r.db.WithContext(ctx).Create(channel)
	if result.Error != nil {
		return result.Error
	}

	// Create channel membership for owner
	membership := domain.ChannelMember{
		ChannelID: channel.ID,
		UserID:    channel.OwnerID,
		Role:      "owner",
		JoinedAt:  time.Now(),
	}

	return r.db.WithContext(ctx).Create(&membership).Error
}

// GetChannel retrieves a channel by ID
func (r *ChannelRepo) GetChannel(ctx context.Context, channelID uint) (*domain.Channel, error) {
	var channel domain.Channel
	result := r.db.WithContext(ctx).First(&channel, channelID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrChannelNotFound
		}
		return nil, result.Error
	}
	return &channel, nil
}

// ListUserChannels lists all channels a user is a member of
func (r *ChannelRepo) ListUserChannels(ctx context.Context, userID uint) ([]domain.Channel, error) {
	var channels []domain.Channel
	err := r.db.WithContext(ctx).
		Joins("JOIN channel_members ON channel_members.channel_id = channels.id").
		Where("channel_members.user_id = ?", userID).
		Find(&channels).Error

	if err != nil {
		return nil, err
	}
	return channels, nil
}

// AddUserToChannel adds a user to a channel
func (r *ChannelRepo) AddUserToChannel(ctx context.Context, channelID uint, userID uint) error {
	// Check if channel exists
	var channel domain.Channel
	if err := r.db.WithContext(ctx).First(&channel, channelID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrChannelNotFound
		}
		return err
	}

	// Check if user already a member
	var count int64
	r.db.WithContext(ctx).
		Model(&domain.ChannelMember{}).
		Where("channel_id = ? AND user_id = ?", channelID, userID).
		Count(&count)

	if count > 0 {
		return repository.ErrUserAlreadyMember
	}

	// Add user to channel
	membership := domain.ChannelMember{
		ChannelID: channelID,
		UserID:    userID,
		Role:      "member",
		JoinedAt:  time.Now(),
	}

	return r.db.WithContext(ctx).Create(&membership).Error
}

// RemoveUserFromChannel removes a user from a channel
func (r *ChannelRepo) RemoveUserFromChannel(ctx context.Context, channelID uint, userID uint) error {
	// Check if channel exists
	var channel domain.Channel
	if err := r.db.WithContext(ctx).First(&channel, channelID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrChannelNotFound
		}
		return err
	}

	// Don't allow removal of channel owner
	if channel.OwnerID == userID {
		return repository.ErrCannotRemoveOwner
	}

	// Delete membership
	result := r.db.WithContext(ctx).
		Where("channel_id = ? AND user_id = ?", channelID, userID).
		Delete(&domain.ChannelMember{})

	if result.RowsAffected == 0 {
		return repository.ErrUserNotMember
	}

	return result.Error
}

// GetChannelMembers retrieves all members of a channel
func (r *ChannelRepo) GetChannelMembers(ctx context.Context, channelID uint) ([]domain.User, error) {
	var users []domain.User
	err := r.db.WithContext(ctx).
		Joins("JOIN channel_members ON channel_members.user_id = users.id").
		Where("channel_members.channel_id = ?", channelID).
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateChannel updates channel information
func (r *ChannelRepo) UpdateChannel(ctx context.Context, channel *domain.Channel) error {
	// Fetch existing to verify it exists
	var existing domain.Channel
	if err := r.db.WithContext(ctx).First(&existing, channel.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrChannelNotFound
		}
		return err
	}

	// Set some fields that shouldn't be updated
	channel.CreatedAt = existing.CreatedAt
	channel.OwnerID = existing.OwnerID // Don't allow changing the owner

	// Update only selected fields
	return r.db.WithContext(ctx).Model(channel).
		Select("name", "description", "type", "icon", "topic", "updated_at").
		Updates(channel).Error
}

// DeleteChannel deletes a channel
func (r *ChannelRepo) DeleteChannel(ctx context.Context, channelID uint) error {
	// Verify channel exists
	var channel domain.Channel
	if err := r.db.WithContext(ctx).First(&channel, channelID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrChannelNotFound
		}
		return err
	}

	// Start transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete channel members
	if err := tx.Where("channel_id = ?", channelID).Delete(&domain.ChannelMember{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete channel
	if err := tx.Delete(&domain.Channel{}, channelID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
