package database

import (
	"context"
	"fmt"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database defines the interface for database operations
type Database interface {
	// Ping checks database connectivity
	Ping(ctx context.Context) error

	// GetDB returns the underlying database connection
	GetDB() *gorm.DB

	// Transaction runs operations in a transaction
	Transaction(fn func(tx *gorm.DB) error) error

	// Close closes the database connection
	Close() error
}

// GormDatabase implements Database interface using GORM
type GormDatabase struct {
	db *gorm.DB
}

// NewDatabase creates a new database connection with retry mechanism
func NewDatabase(cfg *config.DatabaseConfig) (Database, error) {
	var dialector gorm.Dialector

	// Maximum retry attempts
	maxRetries := 5
	var err error
	var db *gorm.DB

	switch cfg.Driver {
	case "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	case "postgres":
		dialector = postgres.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	// Configure GORM logger
	gormLogger := logger.New(
		newLogAdapter(),
		logger.Config{
			SlowThreshold:             time.Duration(cfg.SlowThresholdMs) * time.Millisecond,
			LogLevel:                  getLogLevel(cfg.LogLevel),
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Retry logic for database connection
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Connect to database
		db, err = gorm.Open(dialector, &gorm.Config{
			Logger: gormLogger,
			// Enable PrepareStmt for better performance
			PrepareStmt: true,
			// Disable foreign key constraint when migrating to avoid issues
			DisableForeignKeyConstraintWhenMigrating: true,
		})

		if err == nil {
			break
		}

		// Wait with exponential backoff before retrying
		waitTime := time.Duration(attempt*attempt) * time.Second
		if attempt < maxRetries {
			fmt.Printf("Failed to connect to database, retrying in %v (attempt %d/%d): %v\n",
				waitTime, attempt, maxRetries, err)
			time.Sleep(waitTime)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Configure connection pool based on environment
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMinutes) * time.Minute)

	// Health check
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize connection pool
	for i := 0; i < 5; i++ {
		go func() {
			_ = sqlDB.Ping()
		}()
	}

	return &GormDatabase{db: db}, nil
}

// Ping checks database connectivity
func (d *GormDatabase) Ping(ctx context.Context) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Use context for timeout management
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// GetDB returns the underlying database connection
func (d *GormDatabase) GetDB() *gorm.DB {
	return d.db
}

// Transaction runs operations in a transaction
func (d *GormDatabase) Transaction(fn func(tx *gorm.DB) error) error {
	return d.db.Transaction(fn)
}

// Close closes the database connection
func (d *GormDatabase) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	return sqlDB.Close()
}

// Helper functions for GORM logger
func getLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}

// logAdapter adapts GORM logger to our logger
type logAdapter struct{}

func newLogAdapter() *logAdapter {
	return &logAdapter{}
}

func (l *logAdapter) Printf(format string, args ...interface{}) {
	// Implementation can be updated to use your custom logger
	fmt.Printf(format, args...)
}
