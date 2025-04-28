package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"github.com/JadenRazo/Project-Website/backend/internal/domain"
)

func main() {
	// Parse command line flags
	envFlag := flag.String("env", "development", "Environment to use (development, staging, production)")
	flag.Parse()

	// Load configuration
	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatalf("Failed to load configuration")
	}

	fmt.Printf("Running migrations for environment: %s\n", *envFlag)

	// Initialize database connection
	_, err := db.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	fmt.Println("Running database migrations...")
	err = db.RunMigrations(
		&domain.User{},
		&domain.ShortURL{},
		&domain.Message{},
		&domain.Channel{},
		&domain.Attachment{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations completed successfully")
}
