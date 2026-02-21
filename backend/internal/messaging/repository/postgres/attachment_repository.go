package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	msgerrors "github.com/JadenRazo/Project-Website/backend/internal/messaging/errors"
	"gorm.io/gorm"
)

// AttachmentRepo implements repository.AttachmentRepository using PostgreSQL
type AttachmentRepo struct {
	db *gorm.DB
}

// NewAttachmentRepository creates a new PostgreSQL attachment repository
func NewAttachmentRepository(db *gorm.DB) *AttachmentRepo {
	return &AttachmentRepo{
		db: db,
	}
}

// Create implements repository.BaseRepository
func (r *AttachmentRepo) Create(ctx context.Context, attachment *domain.MessagingAttachment) error {
	if attachment.CreatedAt.IsZero() {
		attachment.CreatedAt = time.Now()
	}
	if attachment.UpdatedAt.IsZero() {
		attachment.UpdatedAt = time.Now()
	}
	return r.db.WithContext(ctx).Create(attachment).Error
}

// FindByID implements repository.BaseRepository
func (r *AttachmentRepo) FindByID(ctx context.Context, id uint) (*domain.MessagingAttachment, error) {
	var attachment domain.MessagingAttachment
	err := r.db.WithContext(ctx).
		Preload("Message").
		First(&attachment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, msgerrors.ErrNotFound
		}
		return nil, err
	}
	return &attachment, nil
}

// Update implements repository.BaseRepository
func (r *AttachmentRepo) Update(ctx context.Context, attachment *domain.MessagingAttachment) error {
	attachment.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(attachment).Error
}

// Delete implements repository.BaseRepository
func (r *AttachmentRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.MessagingAttachment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return msgerrors.ErrNotFound
	}
	return nil
}

// FindAll implements repository.BaseRepository
func (r *AttachmentRepo) FindAll(ctx context.Context) ([]domain.MessagingAttachment, error) {
	var attachments []domain.MessagingAttachment
	err := r.db.WithContext(ctx).
		Preload("Message").
		Find(&attachments).Error
	return attachments, err
}

// GetMessageAttachments retrieves all attachments for a message
func (r *AttachmentRepo) GetMessageAttachments(ctx context.Context, messageID uint) ([]domain.MessagingAttachment, error) {
	var attachments []domain.MessagingAttachment
	err := r.db.WithContext(ctx).
		Where("message_id = ?", messageID).
		Find(&attachments).Error
	return attachments, err
}

// GetAttachmentByURL retrieves an attachment by its URL
func (r *AttachmentRepo) GetAttachmentByURL(ctx context.Context, url string) (*domain.MessagingAttachment, error) {
	var attachment domain.MessagingAttachment
	err := r.db.WithContext(ctx).
		Where("url = ?", url).
		First(&attachment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, msgerrors.ErrNotFound
		}
		return nil, err
	}
	return &attachment, nil
}

// GetAttachmentByHash retrieves an attachment by its hash
func (r *AttachmentRepo) GetAttachmentByHash(ctx context.Context, hash string) (*domain.MessagingAttachment, error) {
	var attachment domain.MessagingAttachment
	err := r.db.WithContext(ctx).
		Where("hash = ?", hash).
		First(&attachment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, msgerrors.ErrNotFound
		}
		return nil, err
	}
	return &attachment, nil
}

// GetAttachmentCount gets the count of attachments for a message
func (r *AttachmentRepo) GetAttachmentCount(ctx context.Context, messageID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.MessagingAttachment{}).
		Where("message_id = ?", messageID).
		Count(&count).Error
	return int(count), err
}
