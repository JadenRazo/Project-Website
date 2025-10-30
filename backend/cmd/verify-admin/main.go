package main

import (
	"fmt"
	"log"

	"github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatal("Failed to load configuration")
	}

	// Initialize database
	database, err := db.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Query for the admin user
	var result struct {
		ID           string
		Email        string
		Username     string
		PasswordHash string `gorm:"column:password_hash"`
		Role         string
		IsActive     bool `gorm:"column:is_active"`
		IsVerified   bool `gorm:"column:is_verified"`
	}

	email := "admin@jadenrazo.dev"
	err = database.Table("users").Where("email = ?", email).First(&result).Error
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}

	fmt.Printf("Found user:\n")
	fmt.Printf("  ID: %s\n", result.ID)
	fmt.Printf("  Email: %s\n", result.Email)
	fmt.Printf("  Username: %s\n", result.Username)
	fmt.Printf("  Role: %s\n", result.Role)
	fmt.Printf("  IsActive: %v\n", result.IsActive)
	fmt.Printf("  IsVerified: %v\n", result.IsVerified)
	fmt.Printf("  Password Hash: %s...\n", result.PasswordHash[:20])

	// Test password
	testPassword := "JadenRazoAdmin5795%"
	err = bcrypt.CompareHashAndPassword([]byte(result.PasswordHash), []byte(testPassword))
	if err == nil {
		fmt.Printf("\n✓ Password verification SUCCESSFUL\n")
	} else {
		fmt.Printf("\n✗ Password verification FAILED: %v\n", err)
	}

	// Check if it's actually an admin
	if result.Role != "admin" {
		fmt.Printf("\nWARNING: User role is '%s', not 'admin'\n", result.Role)
	}
}
