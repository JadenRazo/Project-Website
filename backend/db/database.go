package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB(dbPath string) error {
	var err error
	
	// Configure logger
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	
	// Open database connection
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	
	// Auto migrate schemas
	err = DB.AutoMigrate(&URL{}, &User{}, &ClickAnalytics{}, &CustomDomain{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	
	fmt.Println("Database migration completed")
	return nil
}
