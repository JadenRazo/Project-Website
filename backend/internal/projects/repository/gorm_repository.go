package repository

import (
	"context"
	"strings"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/project"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type gormProjectRepository struct {
	db *gorm.DB
}

// ProjectModel represents the database model for projects
type ProjectModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OwnerID     uuid.UUID      `gorm:"type:uuid;not null" json:"owner_id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"type:varchar(50);not null;default:'draft'" json:"status"`
	RepoURL     string         `gorm:"type:text" json:"repo_url"`
	LiveURL     string         `gorm:"type:text" json:"live_url"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	CreatedAt   time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt   *time.Time     `gorm:"index" json:"deleted_at"`
}

// TableName sets the table name for GORM
func (ProjectModel) TableName() string {
	return "projects"
}

// NewGormRepository creates a new GORM-based project repository
func NewGormRepository(db *gorm.DB) project.Repository {
	return &gormProjectRepository{
		db: db,
	}
}

func (r *gormProjectRepository) Create(ctx context.Context, proj *project.Project) error {
	model := projectToModel(proj)
	result := r.db.WithContext(ctx).Create(&model)
	if result.Error != nil {
		return result.Error
	}

	*proj = *modelToProject(&model)
	return nil
}

func (r *gormProjectRepository) Get(ctx context.Context, id uuid.UUID) (*project.Project, error) {
	var model ProjectModel
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return modelToProject(&model), nil
}

func (r *gormProjectRepository) GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]*project.Project, error) {
	var models []ProjectModel
	result := r.db.WithContext(ctx).Where("owner_id = ? AND deleted_at IS NULL", ownerID).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	projects := make([]*project.Project, len(models))
	for i, model := range models {
		projects[i] = modelToProject(&model)
	}

	return projects, nil
}

func (r *gormProjectRepository) Update(ctx context.Context, proj *project.Project) error {
	model := projectToModel(proj)
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", proj.ID).Updates(&model)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *gormProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&ProjectModel{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", now)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *gormProjectRepository) List(ctx context.Context, filter *project.Filter, pagination *project.Pagination) ([]*project.Project, int, error) {
	query := r.db.WithContext(ctx).Model(&ProjectModel{}).Where("deleted_at IS NULL")

	// Apply filters
	if filter != nil {
		if filter.Status != nil {
			query = query.Where("status = ?", string(*filter.Status))
		}

		if filter.OwnerID != uuid.Nil {
			query = query.Where("owner_id = ?", filter.OwnerID)
		}

		if filter.Tag != "" {
			query = query.Where("? = ANY(tags)", filter.Tag)
		}

		if filter.Search != "" {
			searchTerm := "%" + strings.ToLower(filter.Search) + "%"
			query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
		}

		if !filter.FromDate.IsZero() {
			query = query.Where("created_at >= ?", filter.FromDate)
		}

		if !filter.ToDate.IsZero() {
			query = query.Where("created_at <= ?", filter.ToDate)
		}
	}

	// Count total records
	var total int64
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (pagination.Page - 1) * pagination.PageSize
	query = query.Order("created_at DESC").Limit(pagination.PageSize).Offset(offset)

	var models []ProjectModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	projects := make([]*project.Project, len(models))
	for i, model := range models {
		projects[i] = modelToProject(&model)
	}

	return projects, int(total), nil
}

// Helper functions to convert between domain models and database models
func projectToModel(proj *project.Project) ProjectModel {
	return ProjectModel{
		ID:          proj.ID,
		OwnerID:     proj.OwnerID,
		Name:        proj.Name,
		Description: proj.Description,
		Status:      string(proj.Status),
		RepoURL:     proj.RepoURL,
		LiveURL:     proj.LiveURL,
		Tags:        proj.Tags,
		CreatedAt:   proj.CreatedAt,
		UpdatedAt:   proj.UpdatedAt,
		DeletedAt:   proj.DeletedAt,
	}
}

func modelToProject(model *ProjectModel) *project.Project {
	return &project.Project{
		ID:          model.ID,
		OwnerID:     model.OwnerID,
		Name:        model.Name,
		Description: model.Description,
		Status:      project.Status(model.Status),
		RepoURL:     model.RepoURL,
		LiveURL:     model.LiveURL,
		Tags:        model.Tags,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   model.DeletedAt,
	}
}
