package repository

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/domain/errors"
)

// BaseRepository defines core CRUD operations for all repositories
type BaseRepository[T any] interface {
	// Core CRUD operations
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id interface{}) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id interface{}) error
	FindAll(ctx context.Context, options ...QueryOption) ([]*T, error)
	
	// Advanced operations
	FindWhere(ctx context.Context, conditions map[string]interface{}, options ...QueryOption) ([]*T, error)
	Count(ctx context.Context, conditions map[string]interface{}) (int64, error)
	Exists(ctx context.Context, id interface{}) (bool, error)
	BulkCreate(ctx context.Context, entities []*T) error
	BulkUpdate(ctx context.Context, entities []*T, fields ...string) error
	
	// Transaction support
	WithTransaction(tx *gorm.DB) BaseRepository[T]
	Transaction(fn func(repo BaseRepository[T]) error) error
	
	// Advanced querying
	Paginate(ctx context.Context, page, pageSize int, conditions map[string]interface{}, options ...QueryOption) ([]*T, int64, error)
	
	// Database utilities
	GetDB() *gorm.DB
	WithContext(ctx context.Context) BaseRepository[T]
}

// QueryOption defines options for database queries
type QueryOption func(*gorm.DB) *gorm.DB

// Common query options
func WithPreload(associations ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		for _, assoc := range associations {
			db = db.Preload(assoc)
		}
		return db
	}
}

func WithSelect(fields ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fields)
	}
}

func WithOrder(order string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	}
}

func WithLimit(limit int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func WithOffset(offset int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func WithSoftDelete() QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at IS NULL")
	}
}

// GormBaseRepository implements BaseRepository using GORM
type GormBaseRepository[T any] struct {
	db    *gorm.DB
	model T
}

// NewBaseRepository creates new base repository instance
func NewBaseRepository[T any](db *gorm.DB, model T) BaseRepository[T] {
	return &GormBaseRepository[T]{
		db:    db,
		model: model,
	}
}

// Create implements BaseRepository.Create
func (r *GormBaseRepository[T]) Create(ctx context.Context, entity *T) error {
	if entity == nil {
		return errors.ErrInvalidInput
	}
	
	if err := ValidateUUIDs(entity); err != nil {
		return fmt.Errorf("UUID validation failed: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return handleDBError(err)
	}
	
	return nil
}

// FindByID implements BaseRepository.FindByID
func (r *GormBaseRepository[T]) FindByID(ctx context.Context, id interface{}) (*T, error) {
	if id == nil {
		return nil, errors.ErrInvalidInput
	}
	
	if err := ValidateID(id); err != nil {
		return nil, fmt.Errorf("ID validation failed: %w", err)
	}
	
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, handleDBError(err)
	}
	
	return &entity, nil
}

// Update implements BaseRepository.Update
func (r *GormBaseRepository[T]) Update(ctx context.Context, entity *T) error {
	if entity == nil {
		return errors.ErrInvalidInput
	}
	
	if err := ValidateUUIDs(entity); err != nil {
		return fmt.Errorf("UUID validation failed: %w", err)
	}
	
	result := r.db.WithContext(ctx).Save(entity)
	if result.Error != nil {
		return handleDBError(result.Error)
	}
	
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	
	return nil
}

// Delete implements BaseRepository.Delete
func (r *GormBaseRepository[T]) Delete(ctx context.Context, id interface{}) error {
	if id == nil {
		return errors.ErrInvalidInput
	}
	
	if err := ValidateID(id); err != nil {
		return fmt.Errorf("ID validation failed: %w", err)
	}
	
	result := r.db.WithContext(ctx).Delete(&r.model, id)
	if result.Error != nil {
		return handleDBError(result.Error)
	}
	
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	
	return nil
}

// FindAll implements BaseRepository.FindAll
func (r *GormBaseRepository[T]) FindAll(ctx context.Context, options ...QueryOption) ([]*T, error) {
	var entities []*T
	
	db := r.db.WithContext(ctx)
	for _, option := range options {
		db = option(db)
	}
	
	if err := db.Find(&entities).Error; err != nil {
		return nil, handleDBError(err)
	}
	
	return entities, nil
}

// FindWhere implements BaseRepository.FindWhere
func (r *GormBaseRepository[T]) FindWhere(ctx context.Context, conditions map[string]interface{}, options ...QueryOption) ([]*T, error) {
	var entities []*T
	
	db := r.db.WithContext(ctx)
	
	// Apply conditions
	for key, value := range conditions {
		db = db.Where(key, value)
	}
	
	// Apply options
	for _, option := range options {
		db = option(db)
	}
	
	if err := db.Find(&entities).Error; err != nil {
		return nil, handleDBError(err)
	}
	
	return entities, nil
}

// Count implements BaseRepository.Count
func (r *GormBaseRepository[T]) Count(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	var count int64
	
	db := r.db.WithContext(ctx).Model(&r.model)
	
	// Apply conditions
	for key, value := range conditions {
		db = db.Where(key, value)
	}
	
	if err := db.Count(&count).Error; err != nil {
		return 0, handleDBError(err)
	}
	
	return count, nil
}

// Exists implements BaseRepository.Exists
func (r *GormBaseRepository[T]) Exists(ctx context.Context, id interface{}) (bool, error) {
	if id == nil {
		return false, errors.ErrInvalidInput
	}
	
	if err := ValidateID(id); err != nil {
		return false, fmt.Errorf("ID validation failed: %w", err)
	}
	
	var count int64
	if err := r.db.WithContext(ctx).Model(&r.model).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, handleDBError(err)
	}
	
	return count > 0, nil
}

// BulkCreate implements BaseRepository.BulkCreate
func (r *GormBaseRepository[T]) BulkCreate(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return errors.ErrInvalidInput
	}
	
	// Validate all entities
	for _, entity := range entities {
		if err := ValidateUUIDs(entity); err != nil {
			return fmt.Errorf("UUID validation failed: %w", err)
		}
	}
	
	if err := r.db.WithContext(ctx).CreateInBatches(entities, 100).Error; err != nil {
		return handleDBError(err)
	}
	
	return nil
}

// BulkUpdate implements BaseRepository.BulkUpdate
func (r *GormBaseRepository[T]) BulkUpdate(ctx context.Context, entities []*T, fields ...string) error {
	if len(entities) == 0 {
		return errors.ErrInvalidInput
	}
	
	// Use transaction for bulk updates
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			if err := ValidateUUIDs(entity); err != nil {
				return fmt.Errorf("UUID validation failed: %w", err)
			}
			
			query := tx.Model(entity)
			if len(fields) > 0 {
				query = query.Select(fields)
			}
			
			if err := query.Updates(entity).Error; err != nil {
				return handleDBError(err)
			}
		}
		return nil
	})
}

// WithTransaction implements BaseRepository.WithTransaction
func (r *GormBaseRepository[T]) WithTransaction(tx *gorm.DB) BaseRepository[T] {
	return &GormBaseRepository[T]{
		db:    tx,
		model: r.model,
	}
}

// Transaction implements BaseRepository.Transaction
func (r *GormBaseRepository[T]) Transaction(fn func(repo BaseRepository[T]) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := r.WithTransaction(tx)
		return fn(txRepo)
	})
}

// Paginate implements BaseRepository.Paginate
func (r *GormBaseRepository[T]) Paginate(ctx context.Context, page, pageSize int, conditions map[string]interface{}, options ...QueryOption) ([]*T, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	
	// Count total records
	count, err := r.Count(ctx, conditions)
	if err != nil {
		return nil, 0, err
	}
	
	// Calculate offset
	offset := (page - 1) * pageSize
	
	// Add pagination options
	paginationOptions := append(options, WithLimit(pageSize), WithOffset(offset))
	
	// Find records
	entities, err := r.FindWhere(ctx, conditions, paginationOptions...)
	if err != nil {
		return nil, 0, err
	}
	
	return entities, count, nil
}

// GetDB implements BaseRepository.GetDB
func (r *GormBaseRepository[T]) GetDB() *gorm.DB {
	return r.db
}

// WithContext implements BaseRepository.WithContext
func (r *GormBaseRepository[T]) WithContext(ctx context.Context) BaseRepository[T] {
	return &GormBaseRepository[T]{
		db:    r.db.WithContext(ctx),
		model: r.model,
	}
}

// Database utility functions

// ValidateID validates that an ID is not nil/empty and is of correct type
func ValidateID(id interface{}) error {
	if id == nil {
		return errors.ErrInvalidInput
	}
	
	switch v := id.(type) {
	case uint:
		if v == 0 {
			return errors.ErrInvalidInput
		}
	case uint32:
		if v == 0 {
			return errors.ErrInvalidInput
		}
	case uint64:
		if v == 0 {
			return errors.ErrInvalidInput
		}
	case int:
		if v <= 0 {
			return errors.ErrInvalidInput
		}
	case int32:
		if v <= 0 {
			return errors.ErrInvalidInput
		}
	case int64:
		if v <= 0 {
			return errors.ErrInvalidInput
		}
	case string:
		if v == "" {
			return errors.ErrInvalidInput
		}
		// If it looks like a UUID, validate it
		if len(v) == 36 {
			if _, err := uuid.Parse(v); err != nil {
				return fmt.Errorf("invalid UUID format: %w", err)
			}
		}
	case uuid.UUID:
		if v == uuid.Nil {
			return errors.ErrInvalidInput
		}
	default:
		// For other types, check if it's a zero value
		if isZeroValue(v) {
			return errors.ErrInvalidInput
		}
	}
	
	return nil
}

// ValidateUUIDs validates UUID fields in an entity using reflection
func ValidateUUIDs(entity interface{}) error {
	if entity == nil {
		return errors.ErrInvalidInput
	}
	
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return nil
	}
	
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		
		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}
		
		// Check UUID fields
		if field.Type() == reflect.TypeOf(uuid.UUID{}) {
			uuidVal := field.Interface().(uuid.UUID)
			
			// Check if this is a required UUID field (non-zero UUID expected)
			// Skip validation for fields that can be nil (like optional foreign keys)
			if fieldType.Tag.Get("gorm") != "" {
				gormTag := fieldType.Tag.Get("gorm")
				// If it's marked as not null and it's nil UUID, that's invalid
				if contains(gormTag, "not null") && uuidVal == uuid.Nil {
					return fmt.Errorf("field %s: UUID cannot be nil", fieldType.Name)
				}
			}
		}
		
		// Check pointer to UUID fields
		if field.Type() == reflect.TypeOf((*uuid.UUID)(nil)) && !field.IsNil() {
			uuidPtr := field.Interface().(*uuid.UUID)
			if *uuidPtr == uuid.Nil {
				return fmt.Errorf("field %s: UUID cannot be nil", fieldType.Name)
			}
		}
	}
	
	return nil
}

// HandleDBError converts database errors to domain errors
func handleDBError(err error) error {
	if err == nil {
		return nil
	}
	
	switch err {
	case gorm.ErrRecordNotFound:
		return errors.ErrNotFound
	case gorm.ErrInvalidTransaction:
		return fmt.Errorf("transaction error: %w", err)
	case gorm.ErrNotImplemented:
		return fmt.Errorf("operation not supported: %w", err)
	case gorm.ErrMissingWhereClause:
		return fmt.Errorf("unsafe operation: missing where clause: %w", err)
	case gorm.ErrUnsupportedRelation:
		return fmt.Errorf("unsupported relation: %w", err)
	case gorm.ErrPrimaryKeyRequired:
		return fmt.Errorf("primary key required: %w", err)
	case gorm.ErrModelValueRequired:
		return fmt.Errorf("model value required: %w", err)
	case gorm.ErrInvalidData:
		return errors.ErrInvalidInput
	default:
		// Check for common database constraint violations
		errStr := err.Error()
		if contains(errStr, "duplicate key") || contains(errStr, "UNIQUE constraint") {
			return errors.ErrDuplicate
		}
		if contains(errStr, "foreign key constraint") {
			return fmt.Errorf("foreign key constraint violation: %w", err)
		}
		if contains(errStr, "check constraint") {
			return fmt.Errorf("check constraint violation: %w", err)
		}
		
		// Return the original error wrapped
		return fmt.Errorf("database error: %w", err)
	}
}

// Context utilities

// WithTimeout creates a context with timeout
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = 30 * time.Second // Default timeout
	}
	return context.WithTimeout(parent, timeout)
}

// WithDeadline creates a context with deadline
func WithDeadline(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, deadline)
}

// Helper functions

// isZeroValue checks if a value is the zero value for its type
func isZeroValue(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 len(s) > len(substr)+1 && 
		 func() bool {
			 for i := 1; i <= len(s)-len(substr); i++ {
				 if s[i:i+len(substr)] == substr {
					 return true
				 }
			 }
			 return false
		 }()))
}

// Transaction management utilities

// TransactionManager provides advanced transaction utilities
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// ExecuteInTransaction executes multiple repository operations in a single transaction
func (tm *TransactionManager) ExecuteInTransaction(ctx context.Context, operations func(tx *gorm.DB) error) error {
	return tm.db.WithContext(ctx).Transaction(operations)
}

// ExecuteWithRetry executes an operation with retry logic for transient failures
func (tm *TransactionManager) ExecuteWithRetry(ctx context.Context, maxRetries int, operation func() error) error {
	var lastErr error
	
	for i := 0; i <= maxRetries; i++ {
		if err := operation(); err != nil {
			lastErr = err
			
			// Check if this is a retryable error
			if !isRetryableError(err) {
				return err
			}
			
			// Wait before retry (exponential backoff)
			if i < maxRetries {
				backoff := time.Duration(i+1) * 100 * time.Millisecond
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(backoff):
					continue
				}
			}
		} else {
			return nil
		}
	}
	
	return fmt.Errorf("operation failed after %d retries: %w", maxRetries, lastErr)
}

// isRetryableError determines if an error is worth retrying
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := err.Error()
	retryablePatterns := []string{
		"connection refused",
		"connection reset",
		"timeout",
		"temporary failure",
		"deadlock",
		"lock wait timeout",
	}
	
	for _, pattern := range retryablePatterns {
		if contains(errStr, pattern) {
			return true
		}
	}
	
	return false
}