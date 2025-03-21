package service

import (
	"context"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// MessagingService defines the interface for messaging operations
type MessagingService interface {
	// Message operations
	CreateMessage(ctx context.Context, message *domain.MessagingMessage) error
	GetMessage(ctx context.Context, id uint) (*domain.MessagingMessage, error)
	UpdateMessage(ctx context.Context, message *domain.MessagingMessage) error
	DeleteMessage(ctx context.Context, id uint) error
	GetChannelMessages(ctx context.Context, channelID uint, limit, offset int) ([]*domain.MessagingMessage, error)
	GetThreadMessages(ctx context.Context, threadID uint, limit, offset int) ([]*domain.MessagingMessage, error)
	SearchMessages(ctx context.Context, filters repository.MessageSearchFilters) ([]*domain.MessagingMessage, error)

	// Channel operations
	CreateChannel(ctx context.Context, channel *domain.MessagingChannel) error
	GetChannel(ctx context.Context, id uint) (*domain.MessagingChannel, error)
	UpdateChannel(ctx context.Context, channel *domain.MessagingChannel) error
	DeleteChannel(ctx context.Context, id uint) error
	GetUserChannels(ctx context.Context, userID uint) ([]*domain.MessagingChannel, error)
	AddChannelMember(ctx context.Context, member *domain.MessagingChannelMember) error
	RemoveChannelMember(ctx context.Context, channelID, userID uint) error
	GetChannelMembers(ctx context.Context, channelID uint) ([]*domain.MessagingChannelMember, error)

	// Category operations
	CreateCategory(ctx context.Context, category *domain.MessagingCategory) error
	GetCategory(ctx context.Context, id uint) (*domain.MessagingCategory, error)
	UpdateCategory(ctx context.Context, category *domain.MessagingCategory) error
	DeleteCategory(ctx context.Context, id uint) error
	GetChannelCategories(ctx context.Context, channelID uint) ([]*domain.MessagingCategory, error)

	// User operations
	GetUser(ctx context.Context, id uint) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint) error

	// Role operations
	CreateRole(ctx context.Context, role *domain.MessagingRole) error
	GetRole(ctx context.Context, id uint) (*domain.MessagingRole, error)
	UpdateRole(ctx context.Context, role *domain.MessagingRole) error
	DeleteRole(ctx context.Context, id uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]*domain.MessagingRole, error)
	AssignRole(ctx context.Context, userID, roleID uint) error
	RemoveRole(ctx context.Context, userID, roleID uint) error

	// Reaction operations
	AddReaction(ctx context.Context, reaction *domain.MessagingReaction) error
	RemoveReaction(ctx context.Context, messageID, userID uint, emoji string) error
	GetMessageReactions(ctx context.Context, messageID uint) ([]*domain.MessagingReaction, error)

	// Read receipt operations
	MarkMessageAsRead(ctx context.Context, receipt *domain.MessagingReadReceipt) error
	GetUnreadCount(ctx context.Context, userID, channelID uint) (int64, error)
	GetReadReceipts(ctx context.Context, messageID uint) ([]*domain.MessagingReadReceipt, error)

	// Permission operations
	HasPermission(ctx context.Context, userID, channelID uint, permission domain.Permission) (bool, error)
	GetUserPermissions(ctx context.Context, userID, channelID uint) ([]domain.Permission, error)
}
