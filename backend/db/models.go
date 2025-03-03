package db

import (
	"time"

	"gorm.io/gorm"
)

// Model is a base struct for all models providing common fields
type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Repository defines a standard interface for data access operations
type Repository interface {
	Find(id uint, result interface{}) error
	FindAll(results interface{}) error
	Create(value interface{}) error
	Update(value interface{}) error
	Delete(value interface{}) error
}

// BaseRepository implements common repository operations
type BaseRepository struct {
	DB *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository() *BaseRepository {
	return &BaseRepository{
		DB: GetDB(),
	}
}

// Find retrieves a record by ID
func (r *BaseRepository) Find(id uint, result interface{}) error {
	return r.DB.First(result, id).Error
}

// FindAll retrieves all records
func (r *BaseRepository) FindAll(results interface{}) error {
	return r.DB.Find(results).Error
}

// Create inserts a new record
func (r *BaseRepository) Create(value interface{}) error {
	return r.DB.Create(value).Error
}

// Update updates an existing record
func (r *BaseRepository) Update(value interface{}) error {
	return r.DB.Save(value).Error
}

// Delete removes a record
func (r *BaseRepository) Delete(value interface{}) error {
	return r.DB.Delete(value).Error
}
