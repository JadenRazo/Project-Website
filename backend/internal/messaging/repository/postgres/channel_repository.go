package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	msgerrors "github.com/JadenRazo/Project-Website/backend/internal/messaging/errors"
	"gorm.io/gorm"
)

// ChannelRepo implements repository.ChannelRepository using PostgreSQL
type ChannelRepo struct {
	db *gorm.DB
}

// NewChannelRepository creates a new PostgreSQL channel repository
func NewChannelRepository(db *gorm.DB) *ChannelRepo {
	return &ChannelRepo{
		db: db,
	}
}

// Create implements repository.BaseRepository
func (r *ChannelRepo) Create(ctx context.Context, channel *domain.MessagingChannel) error {
	if channel.CreatedAt.IsZero() {
		channel.CreatedAt = time.Now()
	}
	if channel.UpdatedAt.IsZero() {
		channel.UpdatedAt = time.Now()
	}
	return r.db.WithContext(ctx).Create(channel).Error
}

// FindByID implements repository.BaseRepository
func (r *ChannelRepo) FindByID(ctx context.Context, id uint) (*domain.MessagingChannel, error) {
	var channel domain.MessagingChannel
	err := r.db.WithContext(ctx).
		Preload("Members").
		Preload("PinnedMessages").
		First(&channel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, msgerrors.ErrNotFound
		}
		return nil, err
	}
	return &channel, nil
}

// Update implements repository.BaseRepository
func (r *ChannelRepo) Update(ctx context.Context, channel *domain.MessagingChannel) error {
	channel.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(channel).Error
}

// Delete implements repository.BaseRepository
func (r *ChannelRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.MessagingChannel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return msgerrors.ErrNotFound
	}
	return nil
}

// FindAll implements repository.BaseRepository
func (r *ChannelRepo) FindAll(ctx context.Context) ([]domain.MessagingChannel, error) {
	var channels []domain.MessagingChannel
	err := r.db.WithContext(ctx).
		Preload("Members").
		Preload("PinnedMessages").
		Find(&channels).Error
	return channels, err
}

// AddMember adds a user to a channel
func (r *ChannelRepo) AddMember(ctx context.Context, member *domain.MessagingChannelMember) error {
	member.CreatedAt = time.Now()
	member.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(member).Error
}

// RemoveMember removes a user from a channel
func (r *ChannelRepo) RemoveMember(ctx context.Context, channelID uint, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("channel_id = ? AND user_id = ?", channelID, userID).
		Delete(&domain.MessagingChannelMember{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return msgerrors.ErrNotFound
	}
	return nil
}

// GetChannelMembers retrieves all members of a channel
func (r *ChannelRepo) GetChannelMembers(ctx context.Context, channelID uint) ([]domain.MessagingChannelMember, error) {
	var members []domain.MessagingChannelMember
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("channel_id = ?", channelID).
		Find(&members).Error
	return members, err
}

// GetUserChannels retrieves all channels a user is a member of
func (r *ChannelRepo) GetUserChannels(ctx context.Context, userID uint) ([]domain.MessagingChannel, error) {
	var channels []domain.MessagingChannel
	err := r.db.WithContext(ctx).
		Joins("JOIN messaging_channel_members ON messaging_channel_members.channel_id = messaging_channels.id").
		Where("messaging_channel_members.user_id = ?", userID).
		Find(&channels).Error
	return channels, err
}

// PinMessage pins a message in a channel
func (r *ChannelRepo) PinMessage(ctx context.Context, pin *domain.MessagingPinnedMessage) error {
	pin.CreatedAt = time.Now()
	pin.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(pin).Error
}

// UnpinMessage unpins a message from a channel
func (r *ChannelRepo) UnpinMessage(ctx context.Context, channelID uint, messageID uint) error {
	result := r.db.WithContext(ctx).
		Where("channel_id = ? AND message_id = ?", channelID, messageID).
		Delete(&domain.MessagingPinnedMessage{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return msgerrors.ErrNotFound
	}
	return nil
}

// GetPinnedMessages retrieves all pinned messages in a channel
func (r *ChannelRepo) GetPinnedMessages(ctx context.Context, channelID uint) ([]domain.MessagingPinnedMessage, error) {
	var pins []domain.MessagingPinnedMessage
	err := r.db.WithContext(ctx).
		Preload("Message").
		Preload("PinnedBy").
		Where("channel_id = ?", channelID).
		Order("created_at DESC").
		Find(&pins).Error
	return pins, err
}

// GetPinnedMessageCount gets the count of pinned messages in a channel
func (r *ChannelRepo) GetPinnedMessageCount(ctx context.Context, channelID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.MessagingPinnedMessage{}).
		Where("channel_id = ?", channelID).
		Count(&count).Error
	return int(count), err
}
