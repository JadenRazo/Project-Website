// internal/messaging/service/service.go
package service

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
)

// MessagingService defines the interface for messaging operations
type MessagingService interface {
	// Message operations
	CreateMessage(ctx context.Context, content string, senderID, channelID uint, attachments []domain.MessagingAttachment, embeds []domain.MessagingEmbed) (*domain.MessagingMessage, error)
	GetMessage(ctx context.Context, messageID uint) (*domain.MessagingMessage, error)
	UpdateMessage(ctx context.Context, messageID uint, content string, userID uint) (*domain.MessagingMessage, error)
	DeleteMessage(ctx context.Context, messageID uint, userID uint) error
	GetChannelMessages(ctx context.Context, channelID uint, lastMessageID uint, limit int) ([]domain.MessagingMessage, error)
	GetThreadMessages(ctx context.Context, threadID uint, lastMessageID uint, limit int) ([]domain.MessagingMessage, error)
	PinMessage(ctx context.Context, messageID uint, userID uint) error
	UnpinMessage(ctx context.Context, messageID uint, userID uint) error
	SearchMessages(ctx context.Context, channelID uint, query string, from, to time.Time) ([]domain.MessagingMessage, error)

	// Channel operations
	CreateChannel(ctx context.Context, name, description, channelType string, ownerID uint, isPrivate bool, categoryID *uint) (*domain.MessagingChannel, error)
	GetChannel(ctx context.Context, channelID uint) (*domain.MessagingChannel, error)
	UpdateChannel(ctx context.Context, channelID uint, name, description string, userID uint) (*domain.MessagingChannel, error)
	DeleteChannel(ctx context.Context, channelID uint, userID uint) error
	ListUserChannels(ctx context.Context, userID uint) ([]domain.MessagingChannel, error)
	AddChannelMember(ctx context.Context, channelID, userID, addedByID uint) error
	RemoveChannelMember(ctx context.Context, channelID, userID, removedByID uint) error
	GetChannelMembers(ctx context.Context, channelID uint) ([]domain.MessagingUser, error)
	SetChannelSlowMode(ctx context.Context, channelID uint, seconds int, userID uint) error
	SetChannelNSFW(ctx context.Context, channelID uint, isNSFW bool, userID uint) error

	// Category operations
	CreateCategory(ctx context.Context, name, description string, serverID uint, position int) (*domain.MessagingCategory, error)
	GetCategory(ctx context.Context, categoryID uint) (*domain.MessagingCategory, error)
	UpdateCategory(ctx context.Context, categoryID uint, name, description string, userID uint) (*domain.MessagingCategory, error)
	DeleteCategory(ctx context.Context, categoryID uint, userID uint) error
	ListServerCategories(ctx context.Context, serverID uint) ([]domain.MessagingCategory, error)
	UpdateCategoryPosition(ctx context.Context, categoryID uint, position int, userID uint) error

	// User operations
	UpdateUserStatus(ctx context.Context, userID uint, status, statusMessage string) error
	GetUserStatus(ctx context.Context, userID uint) (*domain.MessagingUser, error)
	BlockUser(ctx context.Context, userID, blockedID uint) error
	UnblockUser(ctx context.Context, userID, blockedID uint) error
	MuteUser(ctx context.Context, userID, mutedID uint, duration time.Duration) error
	UnmuteUser(ctx context.Context, userID, mutedID uint) error
	GetBlockedUsers(ctx context.Context, userID uint) ([]domain.MessagingUser, error)
	GetMutedUsers(ctx context.Context, userID uint) ([]domain.MessagingUser, error)

	// Role operations
	CreateRole(ctx context.Context, name, color string, permissions uint64, userID uint) (*domain.MessagingRole, error)
	GetRole(ctx context.Context, roleID uint) (*domain.MessagingRole, error)
	UpdateRole(ctx context.Context, roleID uint, name, color string, permissions uint64, userID uint) (*domain.MessagingRole, error)
	DeleteRole(ctx context.Context, roleID uint, userID uint) error
	AddUserRole(ctx context.Context, userID, roleID uint, addedByID uint) error
	RemoveUserRole(ctx context.Context, userID, roleID uint, removedByID uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]domain.MessagingRole, error)

	// Reaction operations
	AddReaction(ctx context.Context, messageID, userID uint, emojiCode, emojiName string) error
	RemoveReaction(ctx context.Context, messageID, userID uint, emojiCode string) error
	GetMessageReactions(ctx context.Context, messageID uint) ([]domain.MessagingReaction, error)
	GetUserReactions(ctx context.Context, userID uint) ([]domain.MessagingReaction, error)

	// Read receipt operations
	MarkMessageAsRead(ctx context.Context, messageID, userID uint) error
	MarkChannelAsRead(ctx context.Context, channelID, userID uint) error
	GetMessageReadReceipts(ctx context.Context, messageID uint) ([]domain.MessagingReadReceipt, error)
	GetUserReadReceipts(ctx context.Context, userID uint) ([]domain.MessagingReadReceipt, error)

	// Permission operations
	HasPermission(ctx context.Context, userID uint, channelID uint, permission uint64) (bool, error)
	GetUserPermissions(ctx context.Context, userID uint, channelID uint) (uint64, error)
}
