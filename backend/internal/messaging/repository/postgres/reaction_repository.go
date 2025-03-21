package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
	"gorm.io/gorm"
)

// ReactionRepo implements repository.ReactionRepository using PostgreSQL
type ReactionRepo struct {
	db *gorm.DB
}

// NewReactionRepository creates a new PostgreSQL reaction repository
func NewReactionRepository(db *gorm.DB) repository.ReactionRepository {
	return &ReactionRepo{
		db: db,
	}
}

// Create implements repository.BaseRepository
func (r *ReactionRepo) Create(ctx context.Context, reaction *domain.MessagingReaction) error {
	if reaction.CreatedAt.IsZero() {
		reaction.CreatedAt = time.Now()
	}
	if reaction.UpdatedAt.IsZero() {
		reaction.UpdatedAt = time.Now()
	}
	return r.db.WithContext(ctx).Create(reaction).Error
}

// FindByID implements repository.BaseRepository
func (r *ReactionRepo) FindByID(ctx context.Context, id uint) (*domain.MessagingReaction, error) {
	var reaction domain.MessagingReaction
	err := r.db.WithContext(ctx).
		Preload("Message").
		Preload("User").
		First(&reaction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &reaction, nil
}

// Update implements repository.BaseRepository
func (r *ReactionRepo) Update(ctx context.Context, reaction *domain.MessagingReaction) error {
	reaction.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(reaction).Error
}

// Delete implements repository.BaseRepository
func (r *ReactionRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.MessagingReaction{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// FindAll implements repository.BaseRepository
func (r *ReactionRepo) FindAll(ctx context.Context) ([]domain.MessagingReaction, error) {
	var reactions []domain.MessagingReaction
	err := r.db.WithContext(ctx).
		Preload("Message").
		Preload("User").
		Find(&reactions).Error
	return reactions, err
}

// GetMessageReactions retrieves all reactions for a message
func (r *ReactionRepo) GetMessageReactions(ctx context.Context, messageID uint) ([]domain.MessagingReaction, error) {
	var reactions []domain.MessagingReaction
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("message_id = ?", messageID).
		Find(&reactions).Error
	return reactions, err
}

// GetUserReactions retrieves all reactions by a user
func (r *ReactionRepo) GetUserReactions(ctx context.Context, userID uint) ([]domain.MessagingReaction, error) {
	var reactions []domain.MessagingReaction
	err := r.db.WithContext(ctx).
		Preload("Message").
		Where("user_id = ?", userID).
		Find(&reactions).Error
	return reactions, err
}

// GetReactionCount gets the count of reactions for a message
func (r *ReactionRepo) GetReactionCount(ctx context.Context, messageID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.MessagingReaction{}).
		Where("message_id = ?", messageID).
		Count(&count).Error
	return int(count), err
}

// GetReactionByEmoji gets a specific reaction by emoji code for a message
func (r *ReactionRepo) GetReactionByEmoji(ctx context.Context, messageID uint, emojiCode string) (*domain.MessagingReaction, error) {
	var reaction domain.MessagingReaction
	err := r.db.WithContext(ctx).
		Where("message_id = ? AND emoji_code = ?", messageID, emojiCode).
		First(&reaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &reaction, nil
}

// GetReactionByUser gets a specific reaction by user for a message
func (r *ReactionRepo) GetReactionByUser(ctx context.Context, messageID uint, userID uint) (*domain.MessagingReaction, error) {
	var reaction domain.MessagingReaction
	err := r.db.WithContext(ctx).
		Where("message_id = ? AND user_id = ?", messageID, userID).
		First(&reaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &reaction, nil
}
