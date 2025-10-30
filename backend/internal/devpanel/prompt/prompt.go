
package prompt

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Prompt struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Prompt      string         `gorm:"type:text;not null" json:"prompt"`
	CategoryID  uuid.UUID      `gorm:"type:uuid" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category"`
	IsFeatured  bool           `gorm:"default:false" json:"is_featured"`
	IsVisible   bool           `gorm:"default:true" json:"is_visible"`
	SortOrder   int            `gorm:"default:100" json:"sort_order"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Category struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name        string         `gorm:"type:varchar(255);unique;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	IsVisible   bool           `gorm:"default:true" json:"is_visible"`
	SortOrder   int            `gorm:"default:100" json:"sort_order"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Repository interface {
	// Prompts
	CreatePrompt(prompt *Prompt) error
	GetPromptByID(id uuid.UUID) (*Prompt, error)
	GetAllPrompts(includeHidden bool) ([]Prompt, error)
	GetVisiblePrompts() ([]Prompt, error)
	UpdatePrompt(prompt *Prompt) error
	DeletePrompt(id uuid.UUID) error

	// Categories
	CreateCategory(cat *Category) error
	GetCategoryByID(id uuid.UUID) (*Category, error)
	GetAllCategories(includeHidden bool) ([]Category, error)
	GetVisibleCategories() ([]Category, error)
	UpdateCategory(cat *Category) error
	DeleteCategory(id uuid.UUID) error
}

type Service interface {
	// Prompts
	CreatePrompt(prompt *Prompt) error
	GetPromptByID(id uuid.UUID) (*Prompt, error)
	GetAllPrompts(includeHidden bool) ([]Prompt, error)
	GetVisiblePrompts() ([]Prompt, error)
	UpdatePrompt(id uuid.UUID, updates map[string]interface{}) error
	DeletePrompt(id uuid.UUID) error

	// Categories
	CreateCategory(cat *Category) error
	GetCategoryByID(id uuid.UUID) (*Category, error)
	GetAllCategories(includeHidden bool) ([]Category, error)
	GetVisibleCategories() ([]Category, error)
	UpdateCategory(id uuid.UUID, updates map[string]interface{}) error
	DeleteCategory(id uuid.UUID) error
}

func (Category) TableName() string {
	return "prompt_categories"
}

func (Prompt) TableName() string {
	return "prompts"
}
