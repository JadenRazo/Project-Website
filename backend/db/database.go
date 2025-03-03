// Package db provides database access for the application
package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBManager handles the database connection and operations
type DBManager struct {
	DB        *gorm.DB
	config    *gorm.Config
	models    []interface{}
	dsn       string
	driverName string
}

// NewDBManager creates a new database manager instance
func NewDBManager() *DBManager {
	// Configure logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	return &DBManager{
		config: &gorm.Config{
			Logger: newLogger,
			NowFunc: func() time.Time {
				return time.Now().UTC() // Standardize on UTC for all timestamps
			},
			TranslateError: true, // Translate database errors to GORM errors
		},
		driverName: "sqlite", // Default to SQLite
		models:     make([]interface{}, 0),
	}
}

// RegisterModels adds models to auto-migration list
func (m *DBManager) RegisterModels(models ...interface{}) {
	m.models = append(m.models, models...)
}

// SetSQLite configures the manager to use SQLite with the given path
func (m *DBManager) SetSQLite(dbPath string) *DBManager {
	m.driverName = "sqlite"
	m.dsn = dbPath
	return m
}

// Connect establishes a database connection
func (m *DBManager) Connect() error {
	var err error

	// Log connection attempt
	log.Printf("Connecting to %s database at %s", m.driverName, m.dsn)

	switch m.driverName {
	case "sqlite":
		m.DB, err = gorm.Open(sqlite.Open(m.dsn), m.config)
	default:
		return fmt.Errorf("unsupported database driver: %s", m.driverName)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the underlying SQL DB to set connection pool parameters
	sqlDB, err := m.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Set connection pool parameters for better performance
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")
	return nil
}

// AutoMigrate runs migrations for all registered models
func (m *DBManager) AutoMigrate() error {
	if len(m.models) == 0 {
		log.Println("Warning: No models registered for migration")
		return nil
	}

	log.Printf("Running auto-migration for %d models", len(m.models))
	err := m.DB.AutoMigrate(m.models...)

	if err != nil {
		return fmt.Errorf("failed to run auto-migration: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// Close closes the database connection
func (m *DBManager) Close() error {
	if m.DB != nil {
		sqlDB, err := m.DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get database connection: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}

// Global instance for ease of use in simple applications
var DefaultManager *DBManager

// Initialize creates the default database manager
func Initialize() {
	DefaultManager = NewDBManager()
}

// GetDB returns the global DB instance from the default manager
func GetDB() *gorm.DB {
	if DefaultManager != nil && DefaultManager.DB != nil {
		return DefaultManager.DB
	}
	panic("Database not initialized. Call Initialize() and connect first")
}
