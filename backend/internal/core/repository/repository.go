package repository

import (
	"context"

	"gorm.io/gorm"
)

// Repository defines core repository functionality
type Repository interface {
	// GetDB returns the database connection
	GetDB() *gorm.DB

	// WithContext returns a repository with context
	WithContext(ctx context.Context) Repository

	// Transaction executes operations in a transaction
	Transaction(fn func(tx *gorm.DB) error) error
}

// BaseRepository provides common repository functionality
type BaseRepository struct {
	db  *gorm.DB
	ctx context.Context
}

// NewBaseRepository creates a new repository instance
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{
		db:  db,
		ctx: context.Background(),
	}
}

// GetDB returns the database connection
func (r *BaseRepository) GetDB() *gorm.DB {
	if r.ctx != nil {
		return r.db.WithContext(r.ctx)
	}
	return r.db
}

// WithContext returns a repository with context
func (r *BaseRepository) WithContext(ctx context.Context) Repository {
	if ctx == nil {
		ctx = context.Background()
	}
	return &BaseRepository{
		db:  r.db,
		ctx: ctx,
	}
}

// Transaction executes operations in a transaction
func (r *BaseRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
