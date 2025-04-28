package project

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Status represents the status of a project
type Status string

const (
	// StatusActive indicates an active project
	StatusActive Status = "active"

	// StatusArchived indicates an archived project
	StatusArchived Status = "archived"

	// StatusDraft indicates a draft project
	StatusDraft Status = "draft"
)

// Project represents a project in the developer panel
type Project struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	OwnerID     uuid.UUID  `json:"owner_id" db:"owner_id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Status      Status     `json:"status" db:"status"`
	RepoURL     string     `json:"repo_url" db:"repo_url"`
	LiveURL     string     `json:"live_url,omitempty" db:"live_url"`
	Tags        []string   `json:"tags" db:"tags"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Repository defines the interface for project data access
type Repository interface {
	Create(ctx context.Context, project *Project) error
	Get(ctx context.Context, id uuid.UUID) (*Project, error)
	GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]*Project, error)
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *Filter, pagination *Pagination) ([]*Project, int, error)
}

// Filter contains filtering options for projects
type Filter struct {
	Status   *Status   `json:"status"`
	OwnerID  uuid.UUID `json:"owner_id"`
	Tag      string    `json:"tag"`
	Search   string    `json:"search"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

// Pagination contains pagination information
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// Service defines the project business logic
type Service struct {
	repo Repository
}

// NewService creates a new project service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// ErrProjectNotFound is returned when a project is not found
var ErrProjectNotFound = errors.New("project not found")

// ErrInvalidProject is returned when project data is invalid
var ErrInvalidProject = errors.New("invalid project data")

// Create creates a new project
func (s *Service) Create(ctx context.Context, project *Project) error {
	if project == nil {
		return ErrInvalidProject
	}

	if project.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidProject)
	}

	if project.OwnerID == uuid.Nil {
		return fmt.Errorf("%w: owner ID is required", ErrInvalidProject)
	}

	// Set defaults
	now := time.Now()
	project.ID = uuid.New()
	project.CreatedAt = now
	project.UpdatedAt = now

	if project.Status == "" {
		project.Status = StatusDraft
	}

	return s.repo.Create(ctx, project)
}

// Get retrieves a project by ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Project, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid project ID", ErrInvalidProject)
	}

	project, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, ErrProjectNotFound
	}

	return project, nil
}

// GetByOwner retrieves all projects for an owner
func (s *Service) GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]*Project, error) {
	if ownerID == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid owner ID", ErrInvalidProject)
	}

	return s.repo.GetByOwner(ctx, ownerID)
}

// Update updates an existing project
func (s *Service) Update(ctx context.Context, project *Project) error {
	if project == nil {
		return ErrInvalidProject
	}

	if project.ID == uuid.Nil {
		return fmt.Errorf("%w: project ID is required", ErrInvalidProject)
	}

	// Check if project exists
	existing, err := s.repo.Get(ctx, project.ID)
	if err != nil {
		return err
	}

	if existing == nil {
		return ErrProjectNotFound
	}

	// Update timestamps
	project.CreatedAt = existing.CreatedAt
	project.UpdatedAt = time.Now()

	return s.repo.Update(ctx, project)
}

// Delete soft-deletes a project
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid project ID", ErrInvalidProject)
	}

	// Check if project exists
	project, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if project == nil {
		return ErrProjectNotFound
	}

	// Set deletion timestamp
	now := time.Now()
	project.DeletedAt = &now
	project.UpdatedAt = now
	project.Status = StatusArchived

	return s.repo.Update(ctx, project)
}

// List retrieves a list of projects with filtering and pagination
func (s *Service) List(ctx context.Context, filter *Filter, pagination *Pagination) ([]*Project, int, error) {
	if pagination == nil {
		pagination = &Pagination{
			Page:     1,
			PageSize: 10,
		}
	}

	if pagination.Page < 1 {
		pagination.Page = 1
	}

	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 10
	}

	if filter == nil {
		filter = &Filter{}
	}

	return s.repo.List(ctx, filter, pagination)
}

// Archive changes a project's status to archived
func (s *Service) Archive(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid project ID", ErrInvalidProject)
	}

	// Check if project exists
	project, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if project == nil {
		return ErrProjectNotFound
	}

	// Update status
	project.Status = StatusArchived
	project.UpdatedAt = time.Now()

	return s.repo.Update(ctx, project)
}

// Activate changes a project's status to active
func (s *Service) Activate(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid project ID", ErrInvalidProject)
	}

	// Check if project exists
	project, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if project == nil {
		return ErrProjectNotFound
	}

	// Update status
	project.Status = StatusActive
	project.UpdatedAt = time.Now()

	return s.repo.Update(ctx, project)
}
