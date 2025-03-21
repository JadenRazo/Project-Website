package repository

import (
	"context"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/errors"
)

// BaseRepository defines common CRUD operations
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	FindAll(ctx context.Context) ([]*T, error)
}

// MessageSearchFilters defines criteria for searching messages
type MessageSearchFilters struct {
	ChannelID      *uint
	UserID         *uint
	ThreadID       *uint
	Query          string
	StartTime      *int64
	EndTime        *int64
	HasAttachments bool
	HasMentions    bool
	IsPinned       bool
	IsNSFW         bool
	IsSpoiler      bool
	Limit          int
	Offset         int
}

// MessageRepository defines operations for message management
type MessageRepository interface {
	BaseRepository[domain.MessagingMessage]
	GetChannelMessages(ctx context.Context, channelID uint, limit, offset int) ([]*domain.MessagingMessage, error)
	GetThreadMessages(ctx context.Context, threadID uint, limit, offset int) ([]*domain.MessagingMessage, error)
	SearchMessages(ctx context.Context, filters MessageSearchFilters) ([]*domain.MessagingMessage, error)
	GetMessage(ctx context.Context, id uint) (*domain.MessagingMessage, error)
	AddReaction(ctx context.Context, reaction *domain.MessagingReaction) error
	RemoveReaction(ctx context.Context, reactionID uint) error
	GetMessageReactions(ctx context.Context, messageID uint) ([]domain.MessagingReaction, error)
	MarkAsRead(ctx context.Context, messageID uint, userID uint) error
	GetUnreadCount(ctx context.Context, channelID uint, userID uint) (int, error)
}

// ChannelRepository defines operations for channel management
type ChannelRepository interface {
	BaseRepository[domain.MessagingChannel]
	AddMember(ctx context.Context, member *domain.MessagingChannelMember) error
	RemoveMember(ctx context.Context, channelID, userID uint) error
	GetMembers(ctx context.Context, channelID uint) ([]*domain.MessagingChannelMember, error)
	GetUserChannels(ctx context.Context, userID uint) ([]*domain.MessagingChannel, error)
	GetChannel(ctx context.Context, id uint) (*domain.MessagingChannel, error)
	UpdateChannel(ctx context.Context, channel *domain.MessagingChannel) error
	PinMessage(ctx context.Context, pinnedMessage *domain.MessagingPinnedMessage) error
	UnpinMessage(ctx context.Context, messageID uint) error
	GetPinnedMessages(ctx context.Context, channelID uint) ([]domain.MessagingPinnedMessage, error)
}

// CategoryRepository defines category-specific operations
type CategoryRepository interface {
	BaseRepository[domain.MessagingCategory]
	GetChannelCategories(ctx context.Context, channelID uint) ([]*domain.MessagingCategory, error)
}

// UserRepository defines user-specific operations
type UserRepository interface {
	BaseRepository[domain.User]
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	BlockUser(ctx context.Context, userID uint, blockedID uint) error
	UnblockUser(ctx context.Context, userID uint, blockedID uint) error
	MuteUser(ctx context.Context, userID uint, mutedID uint) error
	UnmuteUser(ctx context.Context, userID uint, mutedID uint) error
	GetBlockedUsers(ctx context.Context, userID uint) ([]domain.MessagingUser, error)
	GetMutedUsers(ctx context.Context, userID uint) ([]domain.MessagingUser, error)
}

// RoleRepository defines role-specific operations
type RoleRepository interface {
	BaseRepository[domain.MessagingRole]
	GetUserRoles(ctx context.Context, userID uint) ([]*domain.MessagingRole, error)
	AssignRole(ctx context.Context, userID, roleID uint) error
	RemoveRole(ctx context.Context, userID, roleID uint) error
}

// ReactionRepository defines operations for reaction management
type ReactionRepository interface {
	BaseRepository[domain.MessagingReaction]
	GetMessageReactions(ctx context.Context, messageID uint) ([]*domain.MessagingReaction, error)
	GetUserReactions(ctx context.Context, userID uint) ([]domain.MessagingReaction, error)
}

// ReadReceiptRepository defines operations for read receipt management
type ReadReceiptRepository interface {
	BaseRepository[domain.MessagingReadReceipt]
	GetUnreadCount(ctx context.Context, userID, channelID uint) (int64, error)
	GetMessageReceipts(ctx context.Context, messageID uint) ([]*domain.MessagingReadReceipt, error)
}

// AttachmentRepository defines operations for attachment management
type AttachmentRepository interface {
	BaseRepository[domain.MessagingAttachment]
	GetMessageAttachments(ctx context.Context, messageID uint) ([]domain.MessagingAttachment, error)
}

// EmbedRepository defines operations for embed management
type EmbedRepository interface {
	BaseRepository[domain.MessagingEmbed]
	GetMessageEmbeds(ctx context.Context, messageID uint) ([]domain.MessagingEmbed, error)
}

// Common repository errors
var (
	ErrNotFound          = errors.New("entity not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrDuplicateEntry    = errors.New("duplicate entry")
	ErrInvalidOperation  = errors.New("invalid operation")
	ErrDatabaseOperation = errors.New("database operation failed")
)
