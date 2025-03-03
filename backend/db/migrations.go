package db

import (
	"fmt"
	"log"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	db      *DBManager
	models  []interface{}
	version string
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *DBManager) *MigrationManager {
	return &MigrationManager{
		db:      db,
		version: "1.0",
	}
}

// RegisterModels registers models for migration
func (m *MigrationManager) RegisterModels(models ...interface{}) {
	m.models = append(m.models, models...)
	m.db.RegisterModels(models...)
}

// Migrate performs database migration
func (m *MigrationManager) Migrate() error {
	log.Printf("Running migrations (version %s)", m.version)
	
	// Register all models
	m.RegisterModels(
		&URL{},
		&User{},
		&ClickAnalytics{},
		&CustomDomain{},
	)
	
	// Run auto-migration
	if err := m.db.AutoMigrate(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	
	// Additional custom migrations could be run here
	
	return nil
}
