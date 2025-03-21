package interfaces

import (
	"context"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
)

// MessageSearchFilters defines filters for searching messages
type MessageSearchFilters struct {
	ChannelID uint
	UserID    uint
	Query     string
	FromDate  string
	ToDate    string
	Limit     int
	Offset    int
}

// BaseRepository defines common CRUD operations
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
}

// MessageRepository defines operations for message entities
type MessageRepository interface {
	BaseRepository[domain.MessagingMessage]
	GetChannelMessages(ctx context.Context, channelID uint, limit, offset int) ([]*domain.MessagingMessage, error)
	GetThreadMessages(ctx context.Context, threadID uint, limit, offset int) ([]*domain.MessagingMessage, error)
	SearchMessages(ctx context.Context, filters MessageSearchFilters) ([]*domain.MessagingMessage, error)
}

// ChannelRepository defines operations for channel entities
type ChannelRepository interface {
	BaseRepository[domain.MessagingChannel]
	GetUserChannels(ctx context.Context, userID uint) ([]*domain.MessagingChannel, error)
	AddMember(ctx context.Context, channelID, userID uint) error
	RemoveMember(ctx context.Context, channelID, userID uint) error
	GetMembers(ctx context.Context, channelID uint) ([]*domain.User, error)
}

// CategoryRepository defines operations for category entities
type CategoryRepository interface {
	BaseRepository[domain.MessagingCategory]
	GetChannelCategories(ctx context.Context, channelID uint) ([]*domain.MessagingCategory, error)
}

// UserRepository defines operations for user entities
type UserRepository interface {
	BaseRepository[domain.User]
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// RoleRepository defines operations for role entities
type RoleRepository interface {
	BaseRepository[domain.MessagingRole]
	GetUserRoles(ctx context.Context, userID uint) ([]*domain.MessagingRole, error)
	AssignRole(ctx context.Context, userID, roleID uint) error
	RemoveRole(ctx context.Context, userID, roleID uint) error
}

// ReactionRepository defines operations for reaction entities
type ReactionRepository interface {
	BaseRepository[domain.MessagingReaction]
	GetMessageReactions(ctx context.Context, messageID uint) ([]*domain.MessagingReaction, error)
}

// ReadReceiptRepository defines operations for read receipt entities
type ReadReceiptRepository interface {
	BaseRepository[domain.MessagingReadReceipt]
	GetMessageReceipts(ctx context.Context, messageID uint) ([]*domain.MessagingReadReceipt, error)
	GetUnreadCount(ctx context.Context, userID, channelID uint) (int64, error)
}

// AttachmentRepository defines operations for attachment entities
type AttachmentRepository interface {
	BaseRepository[domain.MessagingAttachment]
	GetMessageAttachments(ctx context.Context, messageID uint) ([]*domain.MessagingAttachment, error)
}

// EmbedRepository defines operations for embed entities
type EmbedRepository interface {
	BaseRepository[domain.MessagingEmbed]
	GetMessageEmbeds(ctx context.Context, messageID uint) ([]*domain.MessagingEmbed, error)
}
