package repository

import (
	"context"
	"strings"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/skill"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type gormSkillRepository struct {
	db *gorm.DB
}

// SkillModel represents the database model for skills
type SkillModel struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name             string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Description      string         `gorm:"type:text" json:"description"`
	Category         string         `gorm:"type:varchar(50);not null" json:"category"`
	ProficiencyLevel string         `gorm:"type:varchar(50);not null" json:"proficiency_level"`
	ProficiencyValue int            `gorm:"type:integer;not null;check:proficiency_value >= 0 AND proficiency_value <= 100" json:"proficiency_value"`
	IsFeatured       bool           `gorm:"type:boolean;not null;default:false" json:"is_featured"`
	SortOrder        int            `gorm:"type:integer;not null;default:1000" json:"sort_order"`
	IconURL          string         `gorm:"type:text" json:"icon_url"`
	Color            string         `gorm:"type:varchar(7)" json:"color"`
	Tags             pq.StringArray `gorm:"type:text[]" json:"tags"`
	CreatedAt        time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt        *time.Time     `gorm:"index" json:"deleted_at"`
}

// TableName sets the table name for GORM
func (SkillModel) TableName() string {
	return "skills"
}

// NewGormRepository creates a new GORM-based skill repository
func NewGormRepository(db *gorm.DB) skill.Repository {
	return &gormSkillRepository{
		db: db,
	}
}

func (r *gormSkillRepository) Create(ctx context.Context, sk *skill.Skill) error {
	model := skillToModel(sk)
	result := r.db.WithContext(ctx).Create(&model)
	if result.Error != nil {
		return result.Error
	}

	*sk = *modelToSkill(&model)
	return nil
}

func (r *gormSkillRepository) Get(ctx context.Context, id uuid.UUID) (*skill.Skill, error) {
	var model SkillModel
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return modelToSkill(&model), nil
}

func (r *gormSkillRepository) GetByCategory(ctx context.Context, category skill.Category) ([]*skill.Skill, error) {
	var models []SkillModel
	result := r.db.WithContext(ctx).Where("category = ? AND deleted_at IS NULL", string(category)).Order("sort_order ASC, name ASC").Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	skills := make([]*skill.Skill, len(models))
	for i, model := range models {
		skills[i] = modelToSkill(&model)
	}

	return skills, nil
}

func (r *gormSkillRepository) GetFeatured(ctx context.Context) ([]*skill.Skill, error) {
	var models []SkillModel
	result := r.db.WithContext(ctx).Where("is_featured = true AND deleted_at IS NULL").Order("sort_order ASC, name ASC").Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	skills := make([]*skill.Skill, len(models))
	for i, model := range models {
		skills[i] = modelToSkill(&model)
	}

	return skills, nil
}

func (r *gormSkillRepository) Update(ctx context.Context, sk *skill.Skill) error {
	model := skillToModel(sk)
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", sk.ID).Updates(&model)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *gormSkillRepository) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&SkillModel{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", now)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *gormSkillRepository) List(ctx context.Context, filter *skill.Filter, pagination *skill.Pagination) ([]*skill.Skill, int, error) {
	query := r.db.WithContext(ctx).Model(&SkillModel{}).Where("deleted_at IS NULL")

	// Apply filters
	if filter != nil {
		if filter.Category != nil {
			query = query.Where("category = ?", string(*filter.Category))
		}

		if filter.ProficiencyLevel != nil {
			query = query.Where("proficiency_level = ?", string(*filter.ProficiencyLevel))
		}

		if filter.IsFeatured != nil {
			query = query.Where("is_featured = ?", *filter.IsFeatured)
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
	query = query.Order("sort_order ASC, name ASC").Limit(pagination.PageSize).Offset(offset)

	var models []SkillModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	skills := make([]*skill.Skill, len(models))
	for i, model := range models {
		skills[i] = modelToSkill(&model)
	}

	return skills, int(total), nil
}

func (r *gormSkillRepository) UpdateSortOrder(ctx context.Context, skillID uuid.UUID, sortOrder int) error {
	result := r.db.WithContext(ctx).Model(&SkillModel{}).Where("id = ? AND deleted_at IS NULL", skillID).Updates(map[string]interface{}{
		"sort_order": sortOrder,
		"updated_at": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Helper functions to convert between domain models and database models
func skillToModel(sk *skill.Skill) SkillModel {
	return SkillModel{
		ID:               sk.ID,
		Name:             sk.Name,
		Description:      sk.Description,
		Category:         string(sk.Category),
		ProficiencyLevel: string(sk.ProficiencyLevel),
		ProficiencyValue: sk.ProficiencyValue,
		IsFeatured:       sk.IsFeatured,
		SortOrder:        sk.SortOrder,
		IconURL:          sk.IconURL,
		Color:            sk.Color,
		Tags:             sk.Tags,
		CreatedAt:        sk.CreatedAt,
		UpdatedAt:        sk.UpdatedAt,
		DeletedAt:        sk.DeletedAt,
	}
}

func modelToSkill(model *SkillModel) *skill.Skill {
	return &skill.Skill{
		ID:               model.ID,
		Name:             model.Name,
		Description:      model.Description,
		Category:         skill.Category(model.Category),
		ProficiencyLevel: skill.ProficiencyLevel(model.ProficiencyLevel),
		ProficiencyValue: model.ProficiencyValue,
		IsFeatured:       model.IsFeatured,
		SortOrder:        model.SortOrder,
		IconURL:          model.IconURL,
		Color:            model.Color,
		Tags:             model.Tags,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
		DeletedAt:        model.DeletedAt,
	}
}
