
package repository

import (
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/prompt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) prompt.Repository {
	return &gormRepository{db: db}
}

// Prompts

func (r *gormRepository) CreatePrompt(p *prompt.Prompt) error {
	return r.db.Create(p).Error
}

func (r *gormRepository) GetPromptByID(id uuid.UUID) (*prompt.Prompt, error) {
	var p prompt.Prompt
	err := r.db.Preload("Category").First(&p, "id = ?", id).Error
	return &p, err
}

func (r *gormRepository) GetAllPrompts(includeHidden bool) ([]prompt.Prompt, error) {
	var prompts []prompt.Prompt
	query := r.db.Preload("Category").Order("sort_order asc")
	if !includeHidden {
		query = query.Where("is_visible = ?", true)
	}
	err := query.Find(&prompts).Error
	return prompts, err
}

func (r *gormRepository) GetVisiblePrompts() ([]prompt.Prompt, error) {
	var prompts []prompt.Prompt
	err := r.db.Preload("Category").Order("sort_order asc").Where("is_visible = ?", true).Find(&prompts).Error
	return prompts, err
}

func (r *gormRepository) UpdatePrompt(p *prompt.Prompt) error {
	return r.db.Save(p).Error
}

func (r *gormRepository) DeletePrompt(id uuid.UUID) error {
	return r.db.Delete(&prompt.Prompt{}, "id = ?", id).Error
}

// Categories

func (r *gormRepository) CreateCategory(cat *prompt.Category) error {
	return r.db.Create(cat).Error
}

func (r *gormRepository) GetCategoryByID(id uuid.UUID) (*prompt.Category, error) {
	var cat prompt.Category
	err := r.db.First(&cat, "id = ?", id).Error
	return &cat, err
}

func (r *gormRepository) GetAllCategories(includeHidden bool) ([]prompt.Category, error) {
	var cats []prompt.Category
	query := r.db.Model(&prompt.Category{}).Order("sort_order asc")
	if !includeHidden {
		query = query.Where("is_visible = ?", true)
	}
	err := query.Find(&cats).Error
	return cats, err
}

func (r *gormRepository) GetVisibleCategories() ([]prompt.Category, error) {
	var cats []prompt.Category
	err := r.db.Order("sort_order asc").Where("is_visible = ?", true).Find(&cats).Error
	return cats, err
}

func (r *gormRepository) UpdateCategory(cat *prompt.Category) error {
	return r.db.Save(cat).Error
}

func (r *gormRepository) DeleteCategory(id uuid.UUID) error {
	return r.db.Delete(&prompt.Category{}, "id = ?", id).Error
}
