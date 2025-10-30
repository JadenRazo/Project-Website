package projectpath

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ProjectPath struct {
	ID              uuid.UUID      `json:"id"`
	Name            string         `json:"name"`
	Path            string         `json:"path"`
	Description     *string        `json:"description"`
	ExcludePatterns pq.StringArray `json:"exclude_patterns"`
	IsActive        bool           `json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       *time.Time     `json:"deleted_at,omitempty"`
}

type Filter struct {
	Name     string
	IsActive *bool
	Search   string
}

type Pagination struct {
	Page     int
	PageSize int
}

type Repository interface {
	Create(ctx context.Context, path *ProjectPath) error
	Get(ctx context.Context, id uuid.UUID) (*ProjectPath, error)
	Update(ctx context.Context, path *ProjectPath) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *Filter, pagination *Pagination) ([]*ProjectPath, int, error)
	GetActive(ctx context.Context) ([]*ProjectPath, error)
}

type Service interface {
	CreateProjectPath(ctx context.Context, path *ProjectPath) error
	GetProjectPath(ctx context.Context, id uuid.UUID) (*ProjectPath, error)
	UpdateProjectPath(ctx context.Context, path *ProjectPath) error
	DeleteProjectPath(ctx context.Context, id uuid.UUID) error
	ListProjectPaths(ctx context.Context, filter *Filter, pagination *Pagination) ([]*ProjectPath, int, error)
	GetActiveProjectPaths(ctx context.Context) ([]*ProjectPath, error)
}

type CreateRequest struct {
	Name            string   `json:"name" binding:"required,min=1,max=255"`
	Path            string   `json:"path" binding:"required"`
	Description     *string  `json:"description"`
	ExcludePatterns []string `json:"exclude_patterns"`
	IsActive        bool     `json:"is_active"`
}

type UpdateRequest struct {
	Name            string   `json:"name" binding:"required,min=1,max=255"`
	Path            string   `json:"path" binding:"required"`
	Description     *string  `json:"description"`
	ExcludePatterns []string `json:"exclude_patterns"`
	IsActive        bool     `json:"is_active"`
}

type ListRequest struct {
	Name     string `form:"name"`
	IsActive *bool  `form:"is_active"`
	Search   string `form:"search"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type ListResponse struct {
	Data        []*ProjectPath `json:"data"`
	TotalCount  int            `json:"total_count"`
	CurrentPage int            `json:"current_page"`
	PageSize    int            `json:"page_size"`
	HasNext     bool           `json:"has_next"`
}