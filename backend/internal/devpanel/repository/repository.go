package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// VCSType represents the version control system type
type VCSType string

const (
	// VCSTypeGit represents Git repositories
	VCSTypeGit VCSType = "git"

	// VCSTypeSVN represents SVN repositories
	VCSTypeSVN VCSType = "svn"

	// VCSTypeHg represents Mercurial repositories
	VCSTypeHg VCSType = "hg"
)

// VisibilityType represents repository visibility
type VisibilityType string

const (
	// VisibilityPublic represents public repositories
	VisibilityPublic VisibilityType = "public"

	// VisibilityPrivate represents private repositories
	VisibilityPrivate VisibilityType = "private"

	// VisibilityInternal represents internal repositories
	VisibilityInternal VisibilityType = "internal"
)

// Repository represents a code repository
type Repository struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	ProjectID   uuid.UUID      `json:"project_id" db:"project_id"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	URL         string         `json:"url" db:"url"`
	VCSType     VCSType        `json:"vcs_type" db:"vcs_type"`
	Visibility  VisibilityType `json:"visibility" db:"visibility"`
	Branch      string         `json:"branch" db:"branch"`
	LastCommit  string         `json:"last_commit,omitempty" db:"last_commit"`
	LastSync    *time.Time     `json:"last_sync,omitempty" db:"last_sync"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// Storage defines the interface for repository data access
type Storage interface {
	Create(ctx context.Context, repo *Repository) error
	Get(ctx context.Context, id uuid.UUID) (*Repository, error)
	GetByProject(ctx context.Context, projectID uuid.UUID) ([]*Repository, error)
	Update(ctx context.Context, repo *Repository) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *Filter, pagination *Pagination) ([]*Repository, int, error)
}

// Filter contains filtering options for repositories
type Filter struct {
	ProjectID  *uuid.UUID      `json:"project_id"`
	VCSType    *VCSType        `json:"vcs_type"`
	Visibility *VisibilityType `json:"visibility"`
	Search     string          `json:"search"`
	FromDate   *time.Time      `json:"from_date"`
	ToDate     *time.Time      `json:"to_date"`
}

// Pagination contains pagination information
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// Service defines the repository business logic
type Service struct {
	storage Storage
}

// NewService creates a new repository service
func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// ErrRepositoryNotFound is returned when a repository is not found
var ErrRepositoryNotFound = errors.New("repository not found")

// ErrInvalidRepository is returned when repository data is invalid
var ErrInvalidRepository = errors.New("invalid repository data")

// Create creates a new repository
func (s *Service) Create(ctx context.Context, repo *Repository) error {
	if repo == nil {
		return ErrInvalidRepository
	}

	if repo.ProjectID == uuid.Nil {
		return fmt.Errorf("%w: project ID is required", ErrInvalidRepository)
	}

	if repo.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidRepository)
	}

	if repo.URL == "" {
		return fmt.Errorf("%w: URL is required", ErrInvalidRepository)
	}

	// Set defaults
	now := time.Now()
	repo.ID = uuid.New()
	repo.CreatedAt = now
	repo.UpdatedAt = now

	if repo.VCSType == "" {
		repo.VCSType = VCSTypeGit
	}

	if repo.Visibility == "" {
		repo.Visibility = VisibilityPrivate
	}

	if repo.Branch == "" {
		repo.Branch = "main"
	}

	return s.storage.Create(ctx, repo)
}

// Get retrieves a repository by ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Repository, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid repository ID", ErrInvalidRepository)
	}

	repo, err := s.storage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if repo == nil {
		return nil, ErrRepositoryNotFound
	}

	return repo, nil
}

// GetByProject retrieves all repositories for a project
func (s *Service) GetByProject(ctx context.Context, projectID uuid.UUID) ([]*Repository, error) {
	if projectID == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid project ID", ErrInvalidRepository)
	}

	return s.storage.GetByProject(ctx, projectID)
}

// Update updates an existing repository
func (s *Service) Update(ctx context.Context, repo *Repository) error {
	if repo == nil {
		return ErrInvalidRepository
	}

	if repo.ID == uuid.Nil {
		return fmt.Errorf("%w: repository ID is required", ErrInvalidRepository)
	}

	// Check if repository exists
	existing, err := s.storage.Get(ctx, repo.ID)
	if err != nil {
		return err
	}

	if existing == nil {
		return ErrRepositoryNotFound
	}

	// Update timestamps
	repo.CreatedAt = existing.CreatedAt
	repo.UpdatedAt = time.Now()

	return s.storage.Update(ctx, repo)
}

// Delete deletes a repository
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid repository ID", ErrInvalidRepository)
	}

	// Check if repository exists
	repo, err := s.storage.Get(ctx, id)
	if err != nil {
		return err
	}

	if repo == nil {
		return ErrRepositoryNotFound
	}

	return s.storage.Delete(ctx, id)
}

// List retrieves a list of repositories with filtering and pagination
func (s *Service) List(ctx context.Context, filter *Filter, pagination *Pagination) ([]*Repository, int, error) {
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

	return s.storage.List(ctx, filter, pagination)
}

// SyncRepository updates the repository with the latest commit information
func (s *Service) SyncRepository(ctx context.Context, id uuid.UUID, lastCommit string) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid repository ID", ErrInvalidRepository)
	}

	// Check if repository exists
	repo, err := s.storage.Get(ctx, id)
	if err != nil {
		return err
	}

	if repo == nil {
		return ErrRepositoryNotFound
	}

	// Update sync information
	now := time.Now()
	repo.LastCommit = lastCommit
	repo.LastSync = &now
	repo.UpdatedAt = now

	return s.storage.Update(ctx, repo)
}
