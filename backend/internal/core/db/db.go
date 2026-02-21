package db

import (
	"context"
	"fmt"
	"os"
	"regexp"
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

// GetDB returns the database connection with retry logic for startup scenarios
func GetDB() (*gorm.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	// Build PostgreSQL DSN
	var dsn string
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		re := regexp.MustCompile(`^postgres://(?:[^:@/]+(?::[^@/]*)?@)?(?:[^:/]+|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::\d+)?/(?:[^?]+)(?:\?.*)?$`)
		if !re.MatchString(dbURL) {
			fmt.Printf("Warning: Invalid DATABASE_URL format, falling back to individual DB variables.\n")
		} else {
			dsn = dbURL
		}
	}

	if dsn == "" {
		password := getEnvWithDefault("DB_PASSWORD", "postgres")
		dsn = fmt.Sprintf("host=%s port=%s user=%s password='%s' dbname=%s sslmode=%s",
			getEnvWithDefault("DB_HOST", "localhost"),
			getEnvWithDefault("DB_PORT", "5432"),
			getEnvWithDefault("DB_USER", "postgres"),
			password,
			getEnvWithDefault("DB_NAME", "project_website"),
			getEnvWithDefault("DB_SSL_MODE", "disable"),
		)
	}

	maxRetries := 10
	retryDelay := 3 * time.Second

	var db *gorm.DB
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})

		if err == nil {
			sqlDB, sqlErr := db.DB()
			if sqlErr == nil {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				pingErr := sqlDB.PingContext(ctx)
				cancel()

				if pingErr == nil {
					sqlDB.SetMaxIdleConns(5)
					sqlDB.SetMaxOpenConns(20)
					sqlDB.SetConnMaxLifetime(time.Hour)
					dbInstance = db
					if attempt > 1 {
						fmt.Printf("Database connection established after %d attempts\n", attempt)
					}
					return db, nil
				}
				err = pingErr
			} else {
				err = sqlErr
			}
		}

		if attempt < maxRetries {
			fmt.Printf("Database connection attempt %d/%d failed: %v. Retrying in %v...\n", attempt, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// GetDBWithContext returns a DB connection with context
func GetDBWithContext(ctx context.Context) (*gorm.DB, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	return db.WithContext(ctx), nil
}

// RunMigrations is deprecated - use backend/schema.sql to initialize database
// This function is kept for backward compatibility but should not be used
func RunMigrations(models ...interface{}) error {
	return fmt.Errorf("automatic migrations are disabled - please use backend/schema.sql to initialize the database")
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
