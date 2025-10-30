package skill

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProficiencyLevel represents the skill proficiency level
type ProficiencyLevel string

const (
	// ProficiencyBeginner indicates beginner level
	ProficiencyBeginner ProficiencyLevel = "beginner"

	// ProficiencyIntermediate indicates intermediate level
	ProficiencyIntermediate ProficiencyLevel = "intermediate"

	// ProficiencyAdvanced indicates advanced level
	ProficiencyAdvanced ProficiencyLevel = "advanced"

	// ProficiencyExpert indicates expert level
	ProficiencyExpert ProficiencyLevel = "expert"
)

// Category represents the skill category
type Category string

const (
	// CategoryFrontend indicates frontend development skills
	CategoryFrontend Category = "frontend"

	// CategoryBackend indicates backend development skills
	CategoryBackend Category = "backend"

	// CategoryDesign indicates design and tools skills
	CategoryDesign Category = "design"

	// CategoryDatabase indicates database skills
	CategoryDatabase Category = "database"

	// CategoryDevOps indicates DevOps skills
	CategoryDevOps Category = "devops"

	// CategoryLanguage indicates programming language skills
	CategoryLanguage Category = "language"

	// CategoryFramework indicates framework skills
	CategoryFramework Category = "framework"

	// CategoryTool indicates tool skills
	CategoryTool Category = "tool"
)

// Skill represents a skill in the developer panel
type Skill struct {
	ID               uuid.UUID        `json:"id" db:"id"`
	Name             string           `json:"name" db:"name"`
	Description      string           `json:"description" db:"description"`
	Category         Category         `json:"category" db:"category"`
	ProficiencyLevel ProficiencyLevel `json:"proficiency_level" db:"proficiency_level"`
	ProficiencyValue int              `json:"proficiency_value" db:"proficiency_value"`
	IsFeatured       bool             `json:"is_featured" db:"is_featured"`
	SortOrder        int              `json:"sort_order" db:"sort_order"`
	IconURL          string           `json:"icon_url,omitempty" db:"icon_url"`
	Color            string           `json:"color,omitempty" db:"color"`
	Tags             []string         `json:"tags" db:"tags"`
	ProjectCount     int              `json:"project_count,omitempty" db:"-"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time       `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Repository defines the interface for skill data access
type Repository interface {
	Create(ctx context.Context, skill *Skill) error
	Get(ctx context.Context, id uuid.UUID) (*Skill, error)
	GetByCategory(ctx context.Context, category Category) ([]*Skill, error)
	GetFeatured(ctx context.Context) ([]*Skill, error)
	Update(ctx context.Context, skill *Skill) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *Filter, pagination *Pagination) ([]*Skill, int, error)
	UpdateSortOrder(ctx context.Context, skillID uuid.UUID, sortOrder int) error
}

// Filter contains filtering options for skills
type Filter struct {
	Category         *Category         `json:"category"`
	ProficiencyLevel *ProficiencyLevel `json:"proficiency_level"`
	IsFeatured       *bool             `json:"is_featured"`
	Tag              string            `json:"tag"`
	Search           string            `json:"search"`
	FromDate         time.Time         `json:"from_date"`
	ToDate           time.Time         `json:"to_date"`
}

// Pagination contains pagination information
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// Service defines the skill business logic
type Service struct {
	repo Repository
}

// NewService creates a new skill service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// ErrSkillNotFound is returned when a skill is not found
var ErrSkillNotFound = errors.New("skill not found")

// ErrInvalidSkill is returned when skill data is invalid
var ErrInvalidSkill = errors.New("invalid skill data")

// Create creates a new skill
func (s *Service) Create(ctx context.Context, skill *Skill) error {
	if skill == nil {
		return ErrInvalidSkill
	}

	if skill.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidSkill)
	}

	if skill.Category == "" {
		return fmt.Errorf("%w: category is required", ErrInvalidSkill)
	}

	if skill.ProficiencyValue < 0 || skill.ProficiencyValue > 100 {
		return fmt.Errorf("%w: proficiency value must be between 0 and 100", ErrInvalidSkill)
	}

	// Set defaults
	now := time.Now()
	skill.ID = uuid.New()
	skill.CreatedAt = now
	skill.UpdatedAt = now

	if skill.ProficiencyLevel == "" {
		skill.ProficiencyLevel = mapValueToProficiency(skill.ProficiencyValue)
	}

	if skill.SortOrder == 0 {
		skill.SortOrder = 1000 // Default to end of list
	}

	return s.repo.Create(ctx, skill)
}

// Get retrieves a skill by ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Skill, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid skill ID", ErrInvalidSkill)
	}

	skill, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if skill == nil {
		return nil, ErrSkillNotFound
	}

	return skill, nil
}

// GetByCategory retrieves all skills for a category
func (s *Service) GetByCategory(ctx context.Context, category Category) ([]*Skill, error) {
	if category == "" {
		return nil, fmt.Errorf("%w: invalid category", ErrInvalidSkill)
	}

	return s.repo.GetByCategory(ctx, category)
}

// GetFeatured retrieves all featured skills
func (s *Service) GetFeatured(ctx context.Context) ([]*Skill, error) {
	return s.repo.GetFeatured(ctx)
}

// Update updates an existing skill
func (s *Service) Update(ctx context.Context, skill *Skill) error {
	if skill == nil {
		return ErrInvalidSkill
	}

	if skill.ID == uuid.Nil {
		return fmt.Errorf("%w: skill ID is required", ErrInvalidSkill)
	}

	if skill.ProficiencyValue < 0 || skill.ProficiencyValue > 100 {
		return fmt.Errorf("%w: proficiency value must be between 0 and 100", ErrInvalidSkill)
	}

	// Check if skill exists
	existing, err := s.repo.Get(ctx, skill.ID)
	if err != nil {
		return err
	}

	if existing == nil {
		return ErrSkillNotFound
	}

	// Update timestamps and proficiency level
	skill.CreatedAt = existing.CreatedAt
	skill.UpdatedAt = time.Now()

	if skill.ProficiencyLevel == "" {
		skill.ProficiencyLevel = mapValueToProficiency(skill.ProficiencyValue)
	}

	return s.repo.Update(ctx, skill)
}

// Delete soft-deletes a skill
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid skill ID", ErrInvalidSkill)
	}

	// Check if skill exists
	skill, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if skill == nil {
		return ErrSkillNotFound
	}

	return s.repo.Delete(ctx, id)
}

// List retrieves a list of skills with filtering and pagination
func (s *Service) List(ctx context.Context, filter *Filter, pagination *Pagination) ([]*Skill, int, error) {
	if pagination == nil {
		pagination = &Pagination{
			Page:     1,
			PageSize: 50,
		}
	}

	if pagination.Page < 1 {
		pagination.Page = 1
	}

	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 50
	}

	if filter == nil {
		filter = &Filter{}
	}

	return s.repo.List(ctx, filter, pagination)
}

// UpdateSortOrder updates the sort order of a skill
func (s *Service) UpdateSortOrder(ctx context.Context, skillID uuid.UUID, sortOrder int) error {
	if skillID == uuid.Nil {
		return fmt.Errorf("%w: invalid skill ID", ErrInvalidSkill)
	}

	if sortOrder < 0 {
		return fmt.Errorf("%w: sort order must be non-negative", ErrInvalidSkill)
	}

	// Check if skill exists
	skill, err := s.repo.Get(ctx, skillID)
	if err != nil {
		return err
	}

	if skill == nil {
		return ErrSkillNotFound
	}

	return s.repo.UpdateSortOrder(ctx, skillID, sortOrder)
}

// ToggleFeatured toggles the featured status of a skill
func (s *Service) ToggleFeatured(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid skill ID", ErrInvalidSkill)
	}

	// Check if skill exists
	skill, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if skill == nil {
		return ErrSkillNotFound
	}

	// Toggle featured status
	skill.IsFeatured = !skill.IsFeatured
	skill.UpdatedAt = time.Now()

	return s.repo.Update(ctx, skill)
}

// mapValueToProficiency converts a numeric proficiency value to a proficiency level
func mapValueToProficiency(value int) ProficiencyLevel {
	switch {
	case value >= 90:
		return ProficiencyExpert
	case value >= 70:
		return ProficiencyAdvanced
	case value >= 40:
		return ProficiencyIntermediate
	default:
		return ProficiencyBeginner
	}
}

// GetAllCategories returns all available categories
func GetAllCategories() []Category {
	return []Category{
		CategoryFrontend,
		CategoryBackend,
		CategoryDesign,
		CategoryDatabase,
		CategoryDevOps,
		CategoryLanguage,
		CategoryFramework,
		CategoryTool,
	}
}

// GetAllProficiencyLevels returns all available proficiency levels
func GetAllProficiencyLevels() []ProficiencyLevel {
	return []ProficiencyLevel{
		ProficiencyBeginner,
		ProficiencyIntermediate,
		ProficiencyAdvanced,
		ProficiencyExpert,
	}
}
