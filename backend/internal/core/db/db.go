package db

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global DB instance
var dbInstance *gorm.DB

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDB returns the database connection
func GetDB() (*gorm.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	// Build PostgreSQL DSN
	var dsn string
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		// Use DATABASE_URL if provided
		dsn = dbURL
	} else {
		// Build from individual environment variables with proper URL encoding for password
		password := url.QueryEscape(getEnvWithDefault("DB_PASSWORD", "postgres"))
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			getEnvWithDefault("DB_HOST", "195.201.136.53"),
			getEnvWithDefault("DB_PORT", "5432"),
			getEnvWithDefault("DB_USER", "postgres"),
			password,
			getEnvWithDefault("DB_NAME", "project_website"),
			getEnvWithDefault("DB_SSL_MODE", "disable"),
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
