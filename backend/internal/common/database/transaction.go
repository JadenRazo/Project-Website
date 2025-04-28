package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"gorm.io/gorm"
)

// TransactionOptions contains options for database transactions
type TransactionOptions struct {
	// Timeout is the maximum duration a transaction can run before being rolled back
	Timeout time.Duration

	// IsolationLevel sets the transaction isolation level
	IsolationLevel sql.IsolationLevel

	// ReadOnly sets whether the transaction is read-only
	ReadOnly bool
}

// DefaultTransactionOptions returns the default transaction options
func DefaultTransactionOptions() TransactionOptions {
	return TransactionOptions{
		Timeout:        30 * time.Second,
		IsolationLevel: sql.LevelDefault,
		ReadOnly:       false,
	}
}

// WithContext executes a function within a transaction with a context
// If the function returns an error, the transaction is rolled back
// Otherwise, the transaction is committed
func WithContext(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return WithContextAndOptions(ctx, db, DefaultTransactionOptions(), fn)
}

// WithContextAndOptions executes a function within a transaction with a context and options
func WithContextAndOptions(ctx context.Context, db *gorm.DB, opts TransactionOptions, fn func(tx *gorm.DB) error) error {
	// Create a context with timeout if needed
	var (
		txCtx    context.Context
		cancelFn context.CancelFunc
	)

	if opts.Timeout > 0 {
		txCtx, cancelFn = context.WithTimeout(ctx, opts.Timeout)
		defer cancelFn()
	} else {
		txCtx = ctx
	}

	// Start transaction with options
	tx := db.WithContext(txCtx).Session(&gorm.Session{
		SkipDefaultTransaction: true,
	})

	if opts.ReadOnly {
		tx = tx.Session(&gorm.Session{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			QueryFields:            true,
		})
	}

	if opts.IsolationLevel != sql.LevelDefault {
		tx = tx.Begin(&sql.TxOptions{
			Isolation: opts.IsolationLevel,
			ReadOnly:  opts.ReadOnly,
		})
	} else {
		tx = tx.Begin()
	}

	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Start time for metrics
	startTime := time.Now()

	// Set up recovery for panics
	defer func() {
		// Handle panic
		if r := recover(); r != nil {
			// Log the panic
			logger.Error("Panic in database transaction",
				"panic", fmt.Sprintf("%v", r),
				"duration_ms", time.Since(startTime).Milliseconds(),
			)

			// Rollback the transaction
			if err := tx.Rollback().Error; err != nil {
				logger.Error("Failed to rollback transaction after panic",
					"error", err.Error(),
				)
			}

			// Re-panic
			panic(r)
		}
	}()

	// Execute the function
	err := fn(tx)

	// Duration metrics
	duration := time.Since(startTime)

	// Rollback or commit based on error
	if err != nil {
		// Rollback transaction
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			// Log rollback error but return original error
			logger.Error("Failed to rollback transaction",
				"original_error", err.Error(),
				"rollback_error", rollbackErr.Error(),
				"duration_ms", duration.Milliseconds(),
			)
			return fmt.Errorf("transaction failed: %w (rollback failed: %v)", err, rollbackErr)
		}

		logger.Debug("Transaction rolled back",
			"error", err.Error(),
			"duration_ms", duration.Milliseconds(),
		)

		return err
	}

	// Commit transaction
	if commitErr := tx.Commit().Error; commitErr != nil {
		logger.Error("Failed to commit transaction",
			"error", commitErr.Error(),
			"duration_ms", duration.Milliseconds(),
		)
		return fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	logger.Debug("Transaction completed successfully",
		"duration_ms", duration.Milliseconds(),
	)

	return nil
}

// WithRetry executes a function within a transaction with retry logic
// It will retry the transaction up to maxAttempts times if it encounters
// a conflict or deadlock error
func WithRetry(ctx context.Context, db *gorm.DB, maxAttempts int, fn func(tx *gorm.DB) error) error {
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = WithContext(ctx, db, fn)

		// If no error or error is not retryable, break
		if err == nil || !isRetryableError(err) {
			break
		}

		// Log retry attempt
		logger.Warn("Retrying database transaction",
			"attempt", attempt,
			"max_attempts", maxAttempts,
			"error", err.Error(),
		)

		// Exponential backoff
		if attempt < maxAttempts {
			backoff := time.Duration(1<<uint(attempt-1)) * 10 * time.Millisecond
			time.Sleep(backoff)
		}
	}

	return err
}

// isRetryableError checks if an error should trigger a transaction retry
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common retryable errors like deadlocks and conflicts
	// This will depend on your database system
	var retryableErrors = []string{
		"deadlock",
		"conflict",
		"serialize",
		"lock timeout",
		"try again",
		"try restarting transaction",
		"could not serialize access",
	}

	errMsg := err.Error()
	for _, msg := range retryableErrors {
		if strings.Contains(strings.ToLower(errMsg), msg) {
			return true
		}
	}

	return false
}
