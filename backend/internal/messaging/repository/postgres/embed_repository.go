package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	msgerrors "github.com/JadenRazo/Project-Website/backend/internal/messaging/errors"
	"gorm.io/gorm"
)

// EmbedRepo implements repository.EmbedRepository using PostgreSQL
type EmbedRepo struct {
	db *gorm.DB
}

// NewEmbedRepository creates a new PostgreSQL embed repository
func NewEmbedRepository(db *gorm.DB) *EmbedRepo {
	return &EmbedRepo{
		db: db,
	}
}

// Create implements repository.BaseRepository
func (r *EmbedRepo) Create(ctx context.Context, embed *domain.MessagingEmbed) error {
	if embed.CreatedAt.IsZero() {
		embed.CreatedAt = time.Now()
	}
	if embed.UpdatedAt.IsZero() {
		embed.UpdatedAt = time.Now()
	}
	return r.db.WithContext(ctx).Create(embed).Error
}

// FindByID implements repository.BaseRepository
func (r *EmbedRepo) FindByID(ctx context.Context, id uint) (*domain.MessagingEmbed, error) {
	var embed domain.MessagingEmbed
	err := r.db.WithContext(ctx).
		Preload("Message").
		First(&embed, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, msgerrors.ErrNotFound
		}
		return nil, err
	}
	return &embed, nil
}

// Update implements repository.BaseRepository
func (r *EmbedRepo) Update(ctx context.Context, embed *domain.MessagingEmbed) error {
	embed.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(embed).Error
}

// Delete implements repository.BaseRepository
func (r *EmbedRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.MessagingEmbed{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return msgerrors.ErrNotFound
	}
	return nil
}

// FindAll implements repository.BaseRepository
func (r *EmbedRepo) FindAll(ctx context.Context) ([]domain.MessagingEmbed, error) {
	var embeds []domain.MessagingEmbed
	err := r.db.WithContext(ctx).
		Preload("Message").
		Find(&embeds).Error
	return embeds, err
}

// GetMessageEmbeds retrieves all embeds for a message
func (r *EmbedRepo) GetMessageEmbeds(ctx context.Context, messageID uint) ([]domain.MessagingEmbed, error) {
	var embeds []domain.MessagingEmbed
	err := r.db.WithContext(ctx).
		Where("message_id = ?", messageID).
		Find(&embeds).Error
	return embeds, err
}

// GetEmbedByURL retrieves an embed by its URL
func (r *EmbedRepo) GetEmbedByURL(ctx context.Context, url string) (*domain.MessagingEmbed, error) {
	var embed domain.MessagingEmbed
	err := r.db.WithContext(ctx).
		Where("url = ?", url).
		First(&embed).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, msgerrors.ErrNotFound
		}
		return nil, err
	}
	return &embed, nil
}

// GetEmbedCount gets the count of embeds for a message
func (r *EmbedRepo) GetEmbedCount(ctx context.Context, messageID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.MessagingEmbed{}).
		Where("message_id = ?", messageID).
		Count(&count).Error
	return int(count), err
}
