package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/codestats/projectpath"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ProjectPathModel struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name            string         `gorm:"type:varchar(255);not null"`
	Path            string         `gorm:"type:text;not null"`
	Description     *string        `gorm:"type:text"`
	ExcludePatterns pq.StringArray `gorm:"type:text[]"`
	IsActive        bool           `gorm:"default:true"`
	CreatedAt       time.Time      `gorm:"not null"`
	UpdatedAt       time.Time      `gorm:"not null"`
	DeletedAt       *time.Time     `gorm:"index"`
}

func (ProjectPathModel) TableName() string {
	return "project_paths"
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, path *projectpath.ProjectPath) error {
	model := &ProjectPathModel{
		ID:              path.ID,
		Name:            path.Name,
		Path:            path.Path,
		Description:     path.Description,
		ExcludePatterns: path.ExcludePatterns,
		IsActive:        path.IsActive,
		CreatedAt:       path.CreatedAt,
		UpdatedAt:       path.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create project path: %w", err)
	}

	*path = *r.modelToDomain(model)
	return nil
}

func (r *GormRepository) Get(ctx context.Context, id uuid.UUID) (*projectpath.ProjectPath, error) {
	var model ProjectPathModel
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project path not found")
		}
		return nil, fmt.Errorf("failed to get project path: %w", err)
	}

	return r.modelToDomain(&model), nil
}

func (r *GormRepository) Update(ctx context.Context, path *projectpath.ProjectPath) error {
	model := &ProjectPathModel{
		ID:              path.ID,
		Name:            path.Name,
		Path:            path.Path,
		Description:     path.Description,
		ExcludePatterns: path.ExcludePatterns,
		IsActive:        path.IsActive,
		UpdatedAt:       time.Now(),
	}

	result := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", path.ID).
		Updates(model)

	if result.Error != nil {
		return fmt.Errorf("failed to update project path: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("project path not found")
	}

	*path = *r.modelToDomain(model)
	return nil
}

func (r *GormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&ProjectPathModel{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("failed to delete project path: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("project path not found")
	}

	return nil
}

func (r *GormRepository) List(ctx context.Context, filter *projectpath.Filter, pagination *projectpath.Pagination) ([]*projectpath.ProjectPath, int, error) {
	var models []ProjectPathModel
	var total int64

	query := r.db.WithContext(ctx).Model(&ProjectPathModel{}).Where("deleted_at IS NULL")

	if filter != nil {
		if filter.Name != "" {
			query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}
		if filter.Search != "" {
			search := "%" + filter.Search + "%"
			query = query.Where("name ILIKE ? OR path ILIKE ? OR description ILIKE ?", search, search, search)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count project paths: %w", err)
	}

	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.PageSize
		query = query.Offset(offset).Limit(pagination.PageSize)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list project paths: %w", err)
	}

	paths := make([]*projectpath.ProjectPath, len(models))
	for i, model := range models {
		paths[i] = r.modelToDomain(&model)
	}

	return paths, int(total), nil
}

func (r *GormRepository) GetActive(ctx context.Context) ([]*projectpath.ProjectPath, error) {
	var models []ProjectPathModel
	if err := r.db.WithContext(ctx).
		Where("is_active = ? AND deleted_at IS NULL", true).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get active project paths: %w", err)
	}

	paths := make([]*projectpath.ProjectPath, len(models))
	for i, model := range models {
		paths[i] = r.modelToDomain(&model)
	}

	return paths, nil
}

func (r *GormRepository) modelToDomain(model *ProjectPathModel) *projectpath.ProjectPath {
	return &projectpath.ProjectPath{
		ID:              model.ID,
		Name:            model.Name,
		Path:            model.Path,
		Description:     model.Description,
		ExcludePatterns: model.ExcludePatterns,
		IsActive:        model.IsActive,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
		DeletedAt:       model.DeletedAt,
	}
}