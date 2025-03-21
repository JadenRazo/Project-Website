package cache

import (
	"context"
	"fmt"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
	"github.com/redis/go-redis/v9"
)

// ChannelCache implements repository.ChannelRepository using Redis
type ChannelCache struct {
	*BaseCache
	repo repository.ChannelRepository
}

// NewChannelCache creates a new Redis channel cache repository
func NewChannelCache(client redis.UniversalClient, repo repository.ChannelRepository) repository.ChannelRepository {
	return &ChannelCache{
		BaseCache: NewBaseCache(client),
		repo:      repo,
	}
}

// Create implements repository.BaseRepository
func (c *ChannelCache) Create(ctx context.Context, channel *domain.MessagingChannel) error {
	if err := c.repo.Create(ctx, channel); err != nil {
		return err
	}
	return c.invalidateChannelCache(ctx, channel)
}

// FindByID implements repository.BaseRepository
func (c *ChannelCache) FindByID(ctx context.Context, id uint) (*domain.MessagingChannel, error) {
	key := fmt.Sprintf("%s%d", channelKeyPrefix, id)

	var channel domain.MessagingChannel
	if err := c.getFromCache(ctx, key, &channel); err == nil {
		return &channel, nil
	}

	channelPtr, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := c.setInCache(ctx, key, channelPtr, channelCacheTTL); err != nil {
		// Log cache error but don't fail the request
	}

	return channelPtr, nil
}

// Update implements repository.BaseRepository
func (c *ChannelCache) Update(ctx context.Context, channel *domain.MessagingChannel) error {
	if err := c.repo.Update(ctx, channel); err != nil {
		return err
	}
	return c.invalidateChannelCache(ctx, channel)
}

// Delete implements repository.BaseRepository
func (c *ChannelCache) Delete(ctx context.Context, id uint) error {
	if err := c.repo.Delete(ctx, id); err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", channelKeyPrefix, id)
	return c.deleteFromCache(ctx, key)
}

// FindAll implements repository.BaseRepository
func (c *ChannelCache) FindAll(ctx context.Context) ([]domain.MessagingChannel, error) {
	return c.repo.FindAll(ctx)
}

// AddMember implements repository.ChannelRepository
func (c *ChannelCache) AddMember(ctx context.Context, member *domain.MessagingChannelMember) error {
	if err := c.repo.AddMember(ctx, member); err != nil {
		return err
	}
	return c.invalidateMembersCache(ctx, member.ChannelID)
}

// RemoveMember implements repository.ChannelRepository
func (c *ChannelCache) RemoveMember(ctx context.Context, channelID uint, userID uint) error {
	if err := c.repo.RemoveMember(ctx, channelID, userID); err != nil {
		return err
	}
	return c.invalidateMembersCache(ctx, channelID)
}

// GetChannelMembers implements repository.ChannelRepository
func (c *ChannelCache) GetChannelMembers(ctx context.Context, channelID uint) ([]domain.MessagingChannelMember, error) {
	key := fmt.Sprintf(channelMembersKey, channelID)

	var members []domain.MessagingChannelMember
	if err := c.getFromCache(ctx, key, &members); err == nil {
		return members, nil
	}

	members, err := c.repo.GetChannelMembers(ctx, channelID)
	if err != nil {
		return nil, err
	}

	if err := c.setInCache(ctx, key, members, membersCacheTTL); err != nil {
		// Log cache error but don't fail the request
	}

	return members, nil
}

// GetUserChannels implements repository.ChannelRepository
func (c *ChannelCache) GetUserChannels(ctx context.Context, userID uint) ([]domain.MessagingChannel, error) {
	key := fmt.Sprintf(userChannelsKey, userID)

	var channels []domain.MessagingChannel
	if err := c.getFromCache(ctx, key, &channels); err == nil {
		return channels, nil
	}

	channels, err := c.repo.GetUserChannels(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := c.setInCache(ctx, key, channels, channelCacheTTL); err != nil {
		// Log cache error but don't fail the request
	}

	return channels, nil
}

// PinMessage implements repository.ChannelRepository
func (c *ChannelCache) PinMessage(ctx context.Context, pinnedMessage *domain.MessagingPinnedMessage) error {
	if err := c.repo.PinMessage(ctx, pinnedMessage); err != nil {
		return err
	}
	return c.invalidatePinnedMessagesCache(ctx, pinnedMessage.ChannelID)
}

// UnpinMessage implements repository.ChannelRepository
func (c *ChannelCache) UnpinMessage(ctx context.Context, messageID uint) error {
	if err := c.repo.UnpinMessage(ctx, messageID); err != nil {
		return err
	}
	return c.invalidatePinnedMessagesCache(ctx, messageID)
}

// GetPinnedMessages implements repository.ChannelRepository
func (c *ChannelCache) GetPinnedMessages(ctx context.Context, channelID uint) ([]domain.MessagingPinnedMessage, error) {
	key := fmt.Sprintf(pinnedMessagesKey, channelID)

	var pinnedMessages []domain.MessagingPinnedMessage
	if err := c.getFromCache(ctx, key, &pinnedMessages); err == nil {
		return pinnedMessages, nil
	}

	pinnedMessages, err := c.repo.GetPinnedMessages(ctx, channelID)
	if err != nil {
		return nil, err
	}

	if err := c.setInCache(ctx, key, pinnedMessages, pinnedCacheTTL); err != nil {
		// Log cache error but don't fail the request
	}

	return pinnedMessages, nil
}

// Helper functions for cache invalidation
func (c *ChannelCache) invalidateChannelCache(ctx context.Context, channel *domain.MessagingChannel) error {
	keys := []string{
		fmt.Sprintf("%s%d", channelKeyPrefix, channel.ID),
	}
	return c.deleteFromCacheMulti(ctx, keys)
}

func (c *ChannelCache) invalidateMembersCache(ctx context.Context, channelID uint) error {
	key := fmt.Sprintf(channelMembersKey, channelID)
	return c.deleteFromCache(ctx, key)
}

func (c *ChannelCache) invalidatePinnedMessagesCache(ctx context.Context, channelID uint) error {
	key := fmt.Sprintf(pinnedMessagesKey, channelID)
	return c.deleteFromCache(ctx, key)
}
