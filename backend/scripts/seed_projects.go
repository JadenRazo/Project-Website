package main

import (
	"context"
	"log"
	"os"

	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/project"
	projectRepo "github.com/JadenRazo/Project-Website/backend/internal/projects/repository"
	"github.com/google/uuid"
)

func main() {
	// Change to backend directory
	if err := os.Chdir(".."); err != nil {
		log.Fatalf("Failed to change directory: %v", err)
	}

	// Initialize database connection
	database, err := db.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the projects table
	if err := database.AutoMigrate(&projectRepo.ProjectModel{}); err != nil {
		log.Fatalf("Failed to migrate projects table: %v", err)
	}

	// Initialize repository and service
	repository := projectRepo.NewGormRepository(database)
	service := project.NewService(repository)

	// Create a dummy owner ID (in real app, this would be the actual user ID)
	ownerID := uuid.New()

	// Define the projects to seed
	projects := []*project.Project{
		{
			Name:        "Portfolio Website",
			Description: "A modern, responsive portfolio website built with React, TypeScript, and styled-components featuring real-time messaging, URL shortener, and developer panel.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Project-Website",
			LiveURL:     "https://jadenrazo.dev",
			Tags:        []string{"React", "TypeScript", "Go", "PostgreSQL", "WebSocket", "REST API"},
			OwnerID:     ownerID,
		},
		{
			Name:        "Educational Quiz Discord Bot",
			Description: "An advanced Discord bot that leverages LLMs to create educational quizzes with multi-guild support, achievement system, and real-time leaderboards.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Discord-Bot-Python",
			LiveURL:     "https://discord.gg/your-bot-invite",
			Tags:        []string{"Python", "Discord.py", "PostgreSQL", "OpenAI API", "Anthropic Claude", "Google Gemini"},
			OwnerID:     ownerID,
		},
		{
			Name:        "DevPanel",
			Description: "A development environment management system with real-time monitoring, service control, and comprehensive project management capabilities.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Project-Website/tree/main/backend/internal/devpanel",
			LiveURL:     "https://jadenrazo.dev/devpanel",
			Tags:        []string{"React", "Go", "WebSocket", "TypeScript", "Real-time Monitoring"},
			OwnerID:     ownerID,
		},
		{
			Name:        "Messaging Platform",
			Description: "A real-time messaging platform with WebSocket integration, file attachments, reactions, and modern UI similar to Discord.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Project-Website/tree/main/backend/internal/messaging",
			LiveURL:     "https://jadenrazo.dev/messaging",
			Tags:        []string{"React", "WebSocket", "Go", "TypeScript", "Real-time Chat"},
			OwnerID:     ownerID,
		},
		{
			Name:        "URL Shortener Service",
			Description: "A high-performance URL shortening service with analytics, custom short codes, and comprehensive statistics tracking.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Project-Website/tree/main/backend/internal/urlshortener",
			LiveURL:     "https://jadenrazo.dev/s/",
			Tags:        []string{"Go", "PostgreSQL", "Analytics", "REST API", "Microservice"},
			OwnerID:     ownerID,
		},
		{
			Name:        "Code Statistics Tracker",
			Description: "Automated system for tracking lines of code across projects with scheduled updates and API integration.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Project-Website/tree/main/scripts",
			LiveURL:     "https://jadenrazo.dev/api/v1/code/stats",
			Tags:        []string{"Go", "Automation", "CLI", "Statistics", "CRON"},
			OwnerID:     ownerID,
		},
	}

	ctx := context.Background()

	// Seed the projects
	for _, proj := range projects {
		if err := service.Create(ctx, proj); err != nil {
			log.Printf("Failed to create project %s: %v", proj.Name, err)
		} else {
			log.Printf("Successfully created project: %s", proj.Name)
		}
	}

	log.Println("Project seeding completed!")
}
