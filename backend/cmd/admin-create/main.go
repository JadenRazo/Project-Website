package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
)

func main() {
	fmt.Println("=== Admin User Creation Tool ===")
	fmt.Println("WARNING: Only use this in development environments!")
	fmt.Println()

	// Check if we're in development
	env := os.Getenv("ENVIRONMENT")
	if env == "production" {
		fmt.Println("❌ This tool cannot be used in production environments.")
		fmt.Println("Use the web-based admin setup instead: /devpanel/setup")
		os.Exit(1)
	}

	// Connect to database
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Check if admin already exists
	adminAuth := auth.NewAdminAuth(db)
	hasAdmin, err := adminAuth.HasAdminAccount()
	if err != nil {
		log.Fatalf("Failed to check for existing admin: %v", err)
	}

	if hasAdmin {
		fmt.Println("⚠️  An admin account already exists.")
		fmt.Print("Do you want to create another admin? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			fmt.Println("Exiting...")
			return
		}
	}

	// Get admin details
	email, err := getAdminEmail(db)
	if err != nil {
		log.Fatalf("Error getting email: %v", err)
	}

	password, err := getSecurePassword()
	if err != nil {
		log.Fatalf("Error getting password: %v", err)
	}

	// Create admin user
	err = createAdminUser(db, email, password)
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Printf("\n✅ Admin user created successfully!\n")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("You can now login at: http://localhost:8081/devpanel\n")
}

func connectDB() (*gorm.DB, error) {
	// Try to use the same database as the main app
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./database.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate to ensure tables exist
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getAdminEmail(db *gorm.DB) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	// Show configured domains
	domains := os.Getenv("ADMIN_DOMAINS")
	if domains == "" {
		domains = "example.com"
	}
	fmt.Printf("Configured admin domains: %s\n", domains)

	for {
		fmt.Print("Enter admin email: ")
		email, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		email = strings.TrimSpace(email)

		// Validate email domain
		adminAuth := auth.NewAdminAuth(db)
		if !adminAuth.ValidateAdminEmail(email) {
			fmt.Printf("❌ Invalid email. Must use one of the configured admin domains: %s\n", domains)
			continue
		}

		return email, nil
	}
}

func getSecurePassword() (string, error) {
	for {
		fmt.Print("Enter admin password (min 8 chars): ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		fmt.Println() // New line after password input

		if len(password) < 8 {
			fmt.Println("❌ Password must be at least 8 characters long.")
			continue
		}

		fmt.Print("Confirm password: ")
		confirmPassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		fmt.Println()

		if string(password) != string(confirmPassword) {
			fmt.Println("❌ Passwords do not match.")
			continue
		}

		return string(password), nil
	}
}

func createAdminUser(db *gorm.DB, email, password string) error {
	// Check if user already exists
	var existingUser entity.User
	err := db.Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		return errors.New("user with this email already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Create new admin user
	username := strings.Split(email, "@")[0]
	user, err := entity.NewUser(email, username, password, "Administrator")
	if err != nil {
		return err
	}

	user.Role = entity.RoleAdmin
	user.IsVerified = true
	user.IsActive = true

	return db.Create(user).Error
}
