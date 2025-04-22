package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global DB instance
var dbInstance *gorm.DB

// GetDB returns the database connection
func GetDB() (*gorm.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	// Use SQLite for simplicity to make imports work
	// For production, use a persistent database path
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// Enable security features
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	dbInstance = db
	return db, nil
}

// GetDBWithContext returns a DB connection with context
func GetDBWithContext(ctx context.Context) (*gorm.DB, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	return db.WithContext(ctx), nil
}

// RunMigrations runs database migrations
func RunMigrations(models ...interface{}) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	return db.AutoMigrate(models...)
}

// CloseDB closes the database connection
func CloseDB() error {
	if dbInstance == nil {
		return nil
	}

	sqlDB, err := dbInstance.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	return sqlDB.Close()
}
