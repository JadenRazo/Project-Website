package certification

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order" gorm:"default:1000"`
	IsVisible   bool      `json:"is_visible" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Certification struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name             string     `json:"name" gorm:"not null"`
	Issuer           string     `json:"issuer" gorm:"not null"`
	CredentialID     *string    `json:"credential_id"`
	IssueDate        time.Time  `json:"issue_date" gorm:"type:date;not null"`
	ExpiryDate       *time.Time `json:"expiry_date" gorm:"type:date"`
	VerificationURL  *string    `json:"verification_url"`
	VerificationText *string    `json:"verification_text"`
	BadgeURL         *string    `json:"badge_url"`
	Description      *string    `json:"description"`
	CategoryID       *uuid.UUID `json:"category_id"`
	Category         *Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	IsFeatured       bool       `json:"is_featured" gorm:"default:false"`
	IsVisible        bool       `json:"is_visible" gorm:"default:true"`
	SortOrder        int        `json:"sort_order" gorm:"default:1000"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type Repository interface {
	// Certifications
	CreateCertification(cert *Certification) error
	GetCertificationByID(id uuid.UUID) (*Certification, error)
	GetAllCertifications(includeHidden bool) ([]Certification, error)
	GetVisibleCertifications() ([]Certification, error)
	UpdateCertification(cert *Certification) error
	DeleteCertification(id uuid.UUID) error

	// Categories
	CreateCategory(cat *Category) error
	GetCategoryByID(id uuid.UUID) (*Category, error)
	GetAllCategories(includeHidden bool) ([]Category, error)
	GetVisibleCategories() ([]Category, error)
	UpdateCategory(cat *Category) error
	DeleteCategory(id uuid.UUID) error
}

type Service interface {
	// Certifications
	CreateCertification(cert *Certification) error
	GetCertificationByID(id uuid.UUID) (*Certification, error)
	GetAllCertifications(includeHidden bool) ([]Certification, error)
	GetVisibleCertifications() ([]Certification, error)
	UpdateCertification(id uuid.UUID, updates map[string]interface{}) error
	DeleteCertification(id uuid.UUID) error

	// Categories
	CreateCategory(cat *Category) error
	GetCategoryByID(id uuid.UUID) (*Category, error)
	GetAllCategories(includeHidden bool) ([]Category, error)
	GetVisibleCategories() ([]Category, error)
	UpdateCategory(id uuid.UUID, updates map[string]interface{}) error
	DeleteCategory(id uuid.UUID) error
}

func (Category) TableName() string {
	return "certification_categories"
}

func (Certification) TableName() string {
	return "certifications"
}
