package cache

import (
	"context"
	"fmt"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
	"github.com/redis/go-redis/v9"
)

// MessageCache implements repository.MessageRepository using Redis
type MessageCache struct {
	*BaseCache
	repo repository.MessageRepository
}

// NewMessageCache creates a new Redis message cache repository
func NewMessageCache(client redis.UniversalClient, repo repository.MessageRepository) repository.MessageRepository {
	return &MessageCache{
		BaseCache: NewBaseCache(client),
		repo:      repo,
	}
}

// Create implements repository.BaseRepository
func (c *MessageCache) Create(ctx context.Context, message *domain.MessagingMessage) error {
	if err := c.repo.Create(ctx, message); err != nil {
		return err
	}
	return c.invalidateMessageCache(ctx, message)
}

// FindByID implements repository.BaseRepository
func (c *MessageCache) FindByID(ctx context.Context, id uint) (*domain.MessagingMessage, error) {
	key := fmt.Sprintf("%s%d", messageKeyPrefix, id)

	var message domain.MessagingMessage
	if err := c.getFromCache(ctx, key, &message); err == nil {
		return &message, nil
	}

	messagePtr, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := c.setInCache(ctx, key, messagePtr, messageCacheTTL); err != nil {
		// Log cache error but don't fail the request
	}

	return messagePtr, nil
}

// Update implements repository.BaseRepository
func (c *MessageCache) Update(ctx context.Context, message *domain.MessagingMessage) error {
	if err := c.repo.Update(ctx, message); err != nil {
		return err
	}
	return c.invalidateMessageCache(ctx, message)
}

// Delete implements repository.BaseRepository
func (c *MessageCache) Delete(ctx context.Context, id uint) error {
	if err := c.repo.Delete(ctx, id); err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", messageKeyPrefix, id)
	return c.deleteFromCache(ctx, key)
}

// FindAll implements repository.BaseRepository
func (c *MessageCache) FindAll(ctx context.Context) ([]domain.MessagingMessage, error) {
	return c.repo.FindAll(ctx)
}

// GetChannelMessages implements repository.MessageRepository
func (c *MessageCache) GetChannelMessages(ctx context.Context, channelID uint, lastMessageID uint, limit int) ([]domain.MessagingMessage, error) {
	key := fmt.Sprintf(channelMessagesKey, channelID)

	var messages []domain.MessagingMessage
	if err := c.getFromCache(ctx, key, &messages); err == nil {
		return messages, nil
	}

	messages, err := c.repo.GetChannelMessages(ctx, channelID, lastMessageID, limit)
	if err != nil {
		return nil, err
	}

	if err := c.setInCache(ctx, key, messages, channelCacheTTL); err != nil {
		// Log cache error but don't fail the request
	}

	return messages, nil
}

// GetThreadMessages implements repository.MessageRepository
func (c *MessageCache) GetThreadMessages(ctx context.Context, threadID uint, lastMessageID uint, limit int) ([]domain.MessagingMessage, error) {
	key := fmt.Sprintf(threadMessagesKey, threadID)

	var messages []domain.MessagingMessage
	if err := c.getFromCache(ctx, key, &messages); err == nil {
		return messages, nil
	}

	messages, err := c.repo.GetThreadMessages(ctx, threadID, lastMessageID, limit)
	if err != nil {
		return nil, err
	}

	if err := c.setInCache(ctx, key, messages, threadCacheTTL); err != nil {
		// Log cache error but don't fail the request
	}

	return messages, nil
}

// SearchMessages implements repository.MessageRepository
func (c *MessageCache) SearchMessages(ctx context.Context, filters repository.MessageSearchFilters) ([]domain.MessagingMessage, int, error) {
	return c.repo.SearchMessages(ctx, filters)
}

// AddReaction implements repository.MessageRepository
func (c *MessageCache) AddReaction(ctx context.Context, reaction *domain.MessagingReaction) error {
	if err := c.repo.AddReaction(ctx, reaction); err != nil {
		return err
	}
	return c.invalidateReactionsCache(ctx, reaction.MessageID)
}

// RemoveReaction implements repository.MessageRepository
func (c *MessageCache) RemoveReaction(ctx context.Context, reactionID uint) error {
	if err := c.repo.RemoveReaction(ctx, reactionID); err != nil {
		return err
	}
	return c.invalidateReactionsCache(ctx, reactionID)
}

// GetMessageReactions implements repository.MessageRepository
func (c *MessageCache) GetMessageReactions(ctx context.Context, messageID uint) ([]domain.MessagingReaction, error) {
	key := fmt.Sprintf(messageReactionsKey, messageID)

	var reactions []domain.MessagingReaction
	if err := c.getFromCache(ctx, key, &reactions); err == nil {
		return reactions, nil
	}

	reactions, err := c.repo.GetMessageReactions(ctx, messageID)
	if err != nil {
		return nil, err
	}

	if err := c.setInCache(ctx, key, reactions, reactionsCacheTTL); err != nil {
		// Log cache error but don't fail the request
	}

	return reactions, nil
}

// MarkAsRead implements repository.MessageRepository
func (c *MessageCache) MarkAsRead(ctx context.Context, messageID uint, userID uint) error {
	return c.repo.MarkAsRead(ctx, messageID, userID)
}

// GetUnreadCount implements repository.MessageRepository
func (c *MessageCache) GetUnreadCount(ctx context.Context, channelID uint, userID uint) (int, error) {
	return c.repo.GetUnreadCount(ctx, channelID, userID)
}

// Helper functions for cache invalidation
func (c *MessageCache) invalidateMessageCache(ctx context.Context, message *domain.MessagingMessage) error {
	keys := []string{
		fmt.Sprintf("%s%d", messageKeyPrefix, message.ID),
		fmt.Sprintf(channelMessagesKey, message.ChannelID),
	}
	if message.ThreadID != nil {
		keys = append(keys, fmt.Sprintf(threadMessagesKey, *message.ThreadID))
	}
	return c.deleteFromCacheMulti(ctx, keys)
}

func (c *MessageCache) invalidateReactionsCache(ctx context.Context, messageID uint) error {
	key := fmt.Sprintf(messageReactionsKey, messageID)
	return c.deleteFromCache(ctx, key)
}
