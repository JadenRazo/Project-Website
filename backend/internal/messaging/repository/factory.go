package repository

import (
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository/postgres"
	"gorm.io/gorm"
)

// RepositoryFactory creates repository instances
type RepositoryFactory struct {
	db *gorm.DB
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(db *gorm.DB) *RepositoryFactory {
	return &RepositoryFactory{
		db: db,
	}
}

// MessageRepository creates a message repository instance
func (f *RepositoryFactory) MessageRepository() MessageRepository {
	return postgres.NewMessageRepository(f.db)
}

// ChannelRepository creates a channel repository instance
func (f *RepositoryFactory) ChannelRepository() ChannelRepository {
	return postgres.NewChannelRepository(f.db)
}

// AttachmentRepository creates an attachment repository instance
func (f *RepositoryFactory) AttachmentRepository() AttachmentRepository {
	return postgres.NewAttachmentRepository(f.db)
}

// EmbedRepository creates an embed repository instance
func (f *RepositoryFactory) EmbedRepository() EmbedRepository {
	return postgres.NewEmbedRepository(f.db)
}

// ReadReceiptRepository creates a read receipt repository instance
func (f *RepositoryFactory) ReadReceiptRepository() ReadReceiptRepository {
	return postgres.NewReadReceiptRepository(f.db)
}

// ReactionRepository creates a reaction repository instance
func (f *RepositoryFactory) ReactionRepository() ReactionRepository {
	return postgres.NewReactionRepository(f.db)
}
