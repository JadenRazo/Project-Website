// Package db provides database access patterns and implementations
package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ===== Error Definitions =====

// Common repository errors that provide meaningful feedback
var (
	// ErrNotFound is returned when a requested record doesn't exist
	ErrNotFound = errors.New("record not found")
	
	// ErrInvalidInput is returned when the provided input is invalid
	ErrInvalidInput = errors.New("invalid input provided")
	
	// ErrDuplicateKey is returned when a unique constraint is violated
	ErrDuplicateKey = errors.New("unique constraint violation")
	
	// ErrTransactional is returned when a transaction operation fails
	ErrTransactional = errors.New("transaction operation failed")
	
	// ErrTimeout is returned when a database operation times out
	ErrTimeout = errors.New("database operation timed out")
	
	// ErrPermission is returned when the operation isn't permitted
	ErrPermission = errors.New("permission denied for operation")
)

// ===== Core Models =====

// Model is the base struct providing common fields for all database models
type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// ===== Repository Interface =====

// Repository defines a standard generic interface for data access operations
// Type parameter T represents the entity type managed by the repository
type Repository[T any] interface {
	// Basic CRUD operations
	Find(id uint) (*T, error)
	FindAll(options ...QueryOption) ([]T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id uint) error
	
	// Query variations
	FindBy(field string, value interface{}) (*T, error)
	FindAllBy(field string, value interface{}, options ...QueryOption) ([]T, error)
	FindWhere(condition string, args ...interface{}) (*T, error)
	FindAllWhere(condition string, args ...interface{}) ([]T, error)
	
	// Batch operations
	BatchCreate(entities []*T) error
	BatchUpdate(entities []*T) error
	BatchDelete(ids []uint) error
	
	// Aggregation operations
	Count(options ...QueryOption) (int64, error)
	
	// Transaction support
	WithTransaction(tx *gorm.DB) Repository[T]
	
	// Advanced features
	Exists(id uint) (bool, error)
	Paginate(page, pageSize int, options ...QueryOption) ([]T, *Pagination, error)
}

// ===== Pagination =====

// Pagination provides metadata for paginated results
type Pagination struct {
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	TotalItems  int64 `json:"totalItems"`
	TotalPages  int   `json:"totalPages"`
	HasNext     bool  `json:"hasNext"`
	HasPrev     bool  `json:"hasPrev"`
}

// ===== Query Options =====

// SortDirection defines the direction for sorting results
type SortDirection string

const (
	// SortAscending represents ascending order (A-Z, 0-9)
	SortAscending SortDirection = "ASC"
	
	// SortDescending represents descending order (Z-A, 9-0)
	SortDescending SortDirection = "DESC"
)

// QueryOption is a function type for modifying database queries
// This follows the functional options pattern for flexible API design
type QueryOption func(*gorm.DB) *gorm.DB

// WithPreload adds eager loading for a relationship
func WithPreload(relationship string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(relationship)
	}
}

// WithNestedPreload adds eager loading for nested relationships
func WithNestedPreload(path string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(clause.Associations).Preload(path)
	}
}

// WithOrder adds ordering to the query
func WithOrder(field string, direction SortDirection) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf("%s %s", field, direction))
	}
}

// WithLimit adds a limit to the query
func WithLimit(limit int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

// WithOffset adds an offset to the query
func WithOffset(offset int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

// WithTimeout adds a timeout context to the query
func WithTimeout(seconds int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
		return db.WithContext(ctx)
	}
}

// WithSelect specifies columns to select
func WithSelect(columns ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(columns)
	}
}

// WithJoin adds a join to the query
func WithJoin(join string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(join)
	}
}

// WithOmit excludes specific columns from operations
func WithOmit(columns ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Omit(columns...)
	}
}

// ===== Factory Type =====

// SpecializedRepositoryFactory creates a specialized repository from a base repository
type SpecializedRepositoryFactory[T any, S any] func(base Repository[T]) S

// ===== Base Repository Implementation =====

// BaseRepository implements the common repository operations
type BaseRepository[T any] struct {
	db         *gorm.DB
	entityType reflect.Type
	tableName  string
}

// NewRepository creates a new repository for the given entity type
func NewRepository[T any](entityType T) Repository[T] {
	// Get the type of T
	typ := reflect.TypeOf(entityType)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	
	// Create an instance to get the table name
	instance := reflect.New(typ).Interface()
	
	// Try to get table name if the type has a TableName method
	var tableName string
	if tabler, ok := instance.(interface{ TableName() string }); ok {
		tableName = tabler.TableName()
	} else {
		// Default to pluralized type name with snake case if TableName() not defined
		tableName = ToSnakeCase(typ.Name()) + "s"
	}
	
	return &BaseRepository[T]{
		db:         GetDB(),
		entityType: typ,
		tableName:  tableName,
	}
}

// Find retrieves an entity by ID
func (r *BaseRepository[T]) Find(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find entity: %w", err)
	}
	return &entity, nil
}

// FindAll retrieves all instances of the entity with optional query modifiers
func (r *BaseRepository[T]) FindAll(options ...QueryOption) ([]T, error) {
	var entities []T
	
	query := r.db
	for _, option := range options {
		query = option(query)
	}
	
	err := query.Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find entities: %w", err)
	}
	
	return entities, nil
}

// Create adds a new entity to the database
func (r *BaseRepository[T]) Create(entity *T) error {
	err := r.db.Create(entity).Error
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || 
		   strings.Contains(err.Error(), "Duplicate entry") {
			return ErrDuplicateKey
		}
		return fmt.Errorf("failed to create entity: %w", err)
	}
	return nil
}

// Update modifies an existing entity
func (r *BaseRepository[T]) Update(entity *T) error {
	result := r.db.Save(entity)
	
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") || 
		   strings.Contains(result.Error.Error(), "Duplicate entry") {
			return ErrDuplicateKey
		}
		return fmt.Errorf("failed to update entity: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}

// Delete removes an entity by ID
func (r *BaseRepository[T]) Delete(id uint) error {
	var entity T
	result := r.db.Delete(&entity, id)
	
	if result.Error != nil {
		return fmt.Errorf("failed to delete entity: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}

// FindBy retrieves an entity by a specific field value
func (r *BaseRepository[T]) FindBy(field string, value interface{}) (*T, error) {
	var entity T
	err := r.db.Where(fmt.Sprintf("%s = ?", field), value).First(&entity).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find entity by %s: %w", field, err)
	}
	
	return &entity, nil
}

// FindAllBy retrieves all entities that match a field value
func (r *BaseRepository[T]) FindAllBy(field string, value interface{}, options ...QueryOption) ([]T, error) {
	var entities []T
	
	query := r.db.Where(fmt.Sprintf("%s = ?", field), value)
	for _, option := range options {
		query = option(query)
	}
	
	err := query.Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find entities by %s: %w", field, err)
	}
	
	return entities, nil
}

// FindWhere retrieves an entity based on a custom condition
func (r *BaseRepository[T]) FindWhere(condition string, args ...interface{}) (*T, error) {
	var entity T
	err := r.db.Where(condition, args...).First(&entity).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find entity with condition: %w", err)
	}
	
	return &entity, nil
}

// FindAllWhere retrieves all entities that match a custom condition
func (r *BaseRepository[T]) FindAllWhere(condition string, args ...interface{}) ([]T, error) {
	var entities []T
	
	query := r.db.Where(condition, args...)
	
	// Check if any args are query options
	for i := 0; i < len(args); i++ {
		if option, ok := args[i].(QueryOption); ok {
			query = option(query)
			// Remove the option from args to avoid passing it to the WHERE clause
			args = append(args[:i], args[i+1:]...)
			i-- // Adjust index since we removed an element
		}
	}
	
	err := query.Find(&entities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find entities with condition: %w", err)
	}
	
	return entities, nil
}

// BatchCreate adds multiple entities in a single transaction
func (r *BaseRepository[T]) BatchCreate(entities []*T) error {
	if len(entities) == 0 {
		return nil
	}
	
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			if err := tx.Create(entity).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchUpdate modifies multiple entities in a single transaction
func (r *BaseRepository[T]) BatchUpdate(entities []*T) error {
	if len(entities) == 0 {
		return nil
	}
	
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			if err := tx.Save(entity).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchDelete removes multiple entities by their IDs in a single transaction
func (r *BaseRepository[T]) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	
	var entity T
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&entity, ids).Error
	})
}

// Count returns the number of entities matching the query options
func (r *BaseRepository[T]) Count(options ...QueryOption) (int64, error) {
	var count int64
	
	query := r.db.Model(new(T))
	for _, option := range options {
		query = option(query)
	}
	
	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count entities: %w", err)
	}
	
	return count, nil
}

// WithTransaction creates a new repository instance using a transaction
// This allows for composing multiple operations in a single atomic unit
func (r *BaseRepository[T]) WithTransaction(tx *gorm.DB) Repository[T] {
	return &BaseRepository[T]{
		db:         tx,
		entityType: r.entityType,
		tableName:  r.tableName,
	}
}

// Exists checks if an entity with the given ID exists
func (r *BaseRepository[T]) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(new(T)).Where("id = ?", id).Count(&count).Error
	
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	
	return count > 0, nil
}

// Paginate returns a paginated list of entities
func (r *BaseRepository[T]) Paginate(page, pageSize int, options ...QueryOption) ([]T, *Pagination, error) {
	// Normalize pagination parameters
	if page < 1 {
		page = 1
	}
	
	if pageSize < 1 {
		pageSize = 10
	}
	
	var entities []T
	var totalItems int64
	
	// Create separate query instances for counting and fetching
	countQuery := r.db.Model(new(T))
	fetchQuery := r.db.Model(new(T))
	
	// Apply all options to both queries except pagination ones
	for _, option := range options {
		// Skip pagination options for count query
		if !isPaginationOption(option) {
			countQuery = option(countQuery)
		}
		fetchQuery = option(fetchQuery)
	}
	
	// Get total count for pagination
	if err := countQuery.Count(&totalItems).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to count total items: %w", err)
	}
	
	// Apply pagination to fetch query
	fetchQuery = fetchQuery.Limit(pageSize).Offset((page - 1) * pageSize)
	
	// Execute the final query
	if err := fetchQuery.Find(&entities).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to find paginated entities: %w", err)
	}
	
	// Calculate total pages
	totalPages := int((totalItems + int64(pageSize) - 1) / int64(pageSize))
	
	// Build pagination metadata
	pagination := &Pagination{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}
	
	return entities, pagination, nil
}

// ===== Helper Functions =====

// isPaginationOption checks if a query option is related to pagination
// This is used internally to skip pagination options in count queries
func isPaginationOption(option QueryOption) bool {
	optionPtr := reflect.ValueOf(option).Pointer()
	limitPtr := reflect.ValueOf(WithLimit).Pointer()
	offsetPtr := reflect.ValueOf(WithOffset).Pointer()
	
	return optionPtr == limitPtr || optionPtr == offsetPtr
}

// NewTransaction creates a new database transaction with convenient commit and rollback functions
func NewTransaction() (*gorm.DB, func(), func(error) error) {
	tx := GetDB().Begin()
	
	// Return commit function
	commitFn := func() {
		tx.Commit()
	}
	
	// Return rollback function with error handling
	rollbackFn := func(err error) error {
		tx.Rollback()
		return fmt.Errorf("%w: %v", ErrTransactional, err)
	}
	
	return tx, commitFn, rollbackFn
}

// ToSnakeCase converts a string from camelCase or PascalCase to snake_case
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
