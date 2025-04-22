#!/bin/bash
set -e

# Print colorful messages
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Project-Website Database Seed Script${NC}"
echo -e "${BLUE}========================================${NC}"

# Navigate to backend directory
cd "$(dirname "${BASH_SOURCE[0]}")/.."
BACKEND_DIR="$(pwd)"
echo -e "${GREEN}Backend directory:${NC} ${YELLOW}${BACKEND_DIR}${NC}"

# Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go first.${NC}"
    echo -e "  Visit ${YELLOW}https://go.dev/doc/install${NC} for installation instructions."
    exit 1
fi

# Validate seed file location
SEED_DIR="${BACKEND_DIR}/internal/common/database/seeds"
if [ ! -d "$SEED_DIR" ]; then
    SEED_DIR="${BACKEND_DIR}/migrations/seeds"
    if [ ! -d "$SEED_DIR" ]; then
        echo -e "${YELLOW}Creating seeds directory...${NC}"
        mkdir -p "${BACKEND_DIR}/internal/common/database/seeds"
        SEED_DIR="${BACKEND_DIR}/internal/common/database/seeds"
    fi
fi

# Check for seed configuration
echo -e "\n${GREEN}Checking for seed configuration...${NC}"
SEED_CONFIG="${BACKEND_DIR}/config/seeds.yaml"
if [ ! -f "$SEED_CONFIG" ]; then
    echo -e "${YELLOW}Seed configuration not found. Creating default configuration...${NC}"
    
    # Create seeds directory if it doesn't exist
    mkdir -p $(dirname "$SEED_CONFIG")
    
    # Create default seed configuration
    cat > "$SEED_CONFIG" << EOF
# Seed configuration
environment: development
seeds:
  users: true
  projects: true
  messaging_channels: true
  messaging_messages: true
  urls: true
  # Add other seed categories here
  
# Seed data counts
counts:
  users: 10
  projects: 5
  messaging_channels: 3
  messaging_messages: 50
  urls: 15
  # Add other seed count configurations here
EOF

    echo -e "${GREEN}Default seed configuration created at:${NC} ${YELLOW}${SEED_CONFIG}${NC}"
fi

echo -e "\n${GREEN}Seed information:${NC}"
echo -e "  Seed directory: ${YELLOW}${SEED_DIR}${NC}"
echo -e "  Seed config: ${YELLOW}${SEED_CONFIG}${NC}"

# Confirm before proceeding
echo -e "\n${YELLOW}Warning: This will add seed data to your database. In some cases, this might conflict with existing data.${NC}"
read -p "Do you want to continue? (y/N): " CONFIRM
if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}Seed operation canceled.${NC}"
    exit 0
fi

# Check if seed command exists
SEED_CMD="${BACKEND_DIR}/cmd/seed/main.go"
if [ ! -f "$SEED_CMD" ]; then
    echo -e "${YELLOW}Seed command not found. Creating default seed command...${NC}"
    
    # Create seed directory if it doesn't exist
    mkdir -p $(dirname "$SEED_CMD")
    
    # Create default seed command
    cat > "$SEED_CMD" << EOF
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/JadenRazo/Project-Website/backend/internal/common/database"
	"gorm.io/gorm"
)

// Main seeder function
func main() {
	fmt.Println("Starting database seeder...")

	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = filepath.Join("config", "development.yaml")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run seeders
	if err := runSeeders(db.GetDB()); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	fmt.Println("Database seeding completed successfully!")
}

// runSeeders runs all database seeders
func runSeeders(db *gorm.DB) error {
	// Add your seeders here
	seeders := []func(*gorm.DB) error{
		seedUsers,
		seedProjects,
		// Add more seeders as needed
	}

	for _, seeder := range seeders {
		if err := seeder(db); err != nil {
			return err
		}
	}

	return nil
}

// Example seeders
func seedUsers(db *gorm.DB) error {
	fmt.Println("Seeding users...")
	// Implementation depends on your database models
	// Create an admin user
	// db.Create(&models.User{...})
	return nil
}

func seedProjects(db *gorm.DB) error {
	fmt.Println("Seeding projects...")
	// Implementation depends on your database models
	// db.Create(&models.Project{...})
	return nil
}
EOF

    echo -e "${GREEN}Default seed command created at:${NC} ${YELLOW}${SEED_CMD}${NC}"
    echo -e "${YELLOW}Please customize the seed command to match your database models before running it.${NC}"
fi

# Run seed command
echo -e "\n${GREEN}Running database seeder...${NC}"
CONFIG_PATH="${SEED_CONFIG}" go run "${SEED_CMD}"

echo -e "\n${GREEN}Database seeding completed!${NC}"
