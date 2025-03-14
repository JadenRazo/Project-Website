package repository

import (
	"context"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
)

// AttachmentRepository defines the interface for attachment storage operations
type AttachmentRepository interface {
	// CreateAttachment creates a new attachment record
	CreateAttachment(ctx context.Context, attachment *domain.Attachment) error

	// GetAttachment retrieves an attachment by ID
	GetAttachment(ctx context.Context, attachmentID uint) (*domain.Attachment, error)

	// UpdateAttachment updates an attachment record
	UpdateAttachment(ctx context.Context, attachment *domain.Attachment) error

	// DeleteAttachment deletes an attachment record
	DeleteAttachment(ctx context.Context, attachmentID uint) error

	// GetMessageAttachments retrieves all attachments for a message
	GetMessageAttachments(ctx context.Context, messageID uint) ([]domain.Attachment, error)
}
