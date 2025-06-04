package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/JadenRazo/Project-Website/backend/internal/common/database"
	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
)

func main() {
	// Parse command line flags
	action := flag.String("action", "up", "Migration action: up, down, reset")
	envFlag := flag.String("env", "development", "Environment to use (development, staging, production)")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*envFlag)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Running migrations for environment: %s\n", *envFlag)
	fmt.Printf("Action: %s\n", *action)

	// Initialize database connection
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Execute migration based on action
	switch *action {
	case "up":
		err = runMigrationsUp(db)
	case "down":
		err = runMigrationsDown(db)
	case "reset":
		err = resetDatabase(db)
	default:
		log.Fatalf("Unknown action: %s. Use 'up', 'down', or 'reset'", *action)
	}

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migrations completed successfully")
}

// runMigrationsUp applies the schema.sql file
func runMigrationsUp(db database.Database) error {
	fmt.Println("Applying database schema...")
	
	// Get the project root directory
	execDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	
	// Look for schema.sql in the backend directory
	schemaPath := filepath.Join(execDir, "schema.sql")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		// Try looking in the parent directory (if running from cmd/migration)
		schemaPath = filepath.Join(execDir, "..", "..", "schema.sql")
		if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
			return fmt.Errorf("schema.sql not found. Please ensure it exists in the backend directory")
		}
	}
	
	// Read the schema file
	sqlContent, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %v", err)
	}
	
	// Execute the SQL
	gormDB := db.GetDB()
	err = gormDB.Exec(string(sqlContent)).Error
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}
	
	fmt.Println("Schema applied successfully")
	return nil
}

// runMigrationsDown drops all tables (dangerous!)
func runMigrationsDown(db database.Database) error {
	fmt.Println("WARNING: This will drop all tables and data!")
	fmt.Print("Are you sure? (type 'yes' to confirm): ")
	
	var confirmation string
	fmt.Scanln(&confirmation)
	
	if confirmation != "yes" {
		fmt.Println("Migration cancelled")
		return nil
	}
	
	gormDB := db.GetDB()
	
	// Drop all tables in reverse dependency order
	tables := []string{
		"rate_limits",
		"system_logs",
		"projects",
		"read_receipts",
		"message_reactions",
		"message_attachments",
		"messages",
		"channel_members",
		"channels",
		"url_clicks",
		"shortened_urls",
		"users",
	}
	
	for _, table := range tables {
		err := gormDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error
		if err != nil {
			fmt.Printf("Warning: Failed to drop table %s: %v\n", table, err)
		}
	}
	
	// Drop functions
	err := gormDB.Exec("DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE").Error
	if err != nil {
		fmt.Printf("Warning: Failed to drop function: %v\n", err)
	}
	
	fmt.Println("All tables dropped successfully")
	return nil
}

// resetDatabase drops all tables and reapplies the schema
func resetDatabase(db database.Database) error {
	fmt.Println("Resetting database (drop all tables and reapply schema)...")
	
	err := runMigrationsDown(db)
	if err != nil {
		return fmt.Errorf("failed to drop tables: %v", err)
	}
	
	err = runMigrationsUp(db)
	if err != nil {
		return fmt.Errorf("failed to apply schema: %v", err)
	}
	
	fmt.Println("Database reset completed successfully")
	return nil
}
