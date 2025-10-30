package main

import (
	"fmt"
	"log"

	"github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Username     string    `gorm:"type:citext;uniqueIndex;not null"`
	Email        string    `gorm:"type:citext;uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	FullName     string
	Role         string `gorm:"default:'user'"`
	IsActive     bool   `gorm:"default:true"`
	IsVerified   bool   `gorm:"default:false"`
}

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

	// Admin credentials
	email := "admin@jadenrazo.dev"
	password := "JadenRazoAdmin5795%"

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Check if admin already exists
	var existingUser User
	result := database.Where("email = ?", email).First(&existingUser)

	if result.Error == nil {
		// Update existing user
		existingUser.PasswordHash = string(hashedPassword)
		existingUser.IsActive = true
		existingUser.IsVerified = true
		existingUser.Role = "admin"

		if err := database.Save(&existingUser).Error; err != nil {
			log.Fatalf("Failed to update admin user: %v", err)
		}

		fmt.Println("Admin user updated successfully!")
	} else if result.Error == gorm.ErrRecordNotFound {
		// Create new admin user
		adminUser := User{
			ID:           uuid.New(),
			Email:        email,
			Username:     "jadenrazo_admin",
			PasswordHash: string(hashedPassword),
			FullName:     "Administrator",
			IsActive:     true,
			IsVerified:   true,
			Role:         "admin",
		}

		if err := database.Create(&adminUser).Error; err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}

		fmt.Println("Admin user created successfully!")
	} else {
		log.Fatalf("Database error: %v", result.Error)
	}

	fmt.Printf("\nAdmin credentials:\n")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Password: %s\n", password)
	fmt.Println("\nYou can now log in to the DevPanel at https://jadenrazo.dev/devpanel")

	// Set environment variable hint
	fmt.Println("\nMake sure ADMIN_DOMAINS environment variable includes 'jadenrazo.dev'")
}
