package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/project"
	"github.com/google/uuid"
)

// MemoryProjectService is an in-memory implementation of the project service
type MemoryProjectService struct {
	projects map[string]*project.Project
	mutex    sync.RWMutex
}

// NewMemoryProjectService creates a new in-memory project service
func NewMemoryProjectService() *MemoryProjectService {
	service := &MemoryProjectService{
		projects: make(map[string]*project.Project),
		mutex:    sync.RWMutex{},
	}

	// Seed with initial data
	service.seedData()
	return service
}

func (s *MemoryProjectService) seedData() {
	ownerID := uuid.New()

	seedProjects := []*project.Project{
		{
			ID:          uuid.New(),
			OwnerID:     ownerID,
			Name:        "Portfolio Website",
			Description: "A modern, responsive portfolio website built with React, TypeScript, and styled-components featuring real-time messaging, URL shortener, and developer panel.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Project-Website",
			LiveURL:     "https://jadenrazo.dev",
			Tags:        []string{"React", "TypeScript", "Go", "PostgreSQL", "WebSocket", "REST API"},
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			OwnerID:     ownerID,
			Name:        "Educational Quiz Discord Bot",
			Description: "An advanced Discord bot that leverages LLMs to create educational quizzes with multi-guild support, achievement system, and real-time leaderboards.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Discord-Bot-Python",
			LiveURL:     "https://discord.gg/your-bot-invite",
			Tags:        []string{"Python", "Discord.py", "PostgreSQL", "OpenAI API", "Anthropic Claude", "Google Gemini"},
			CreatedAt:   time.Now().Add(-20 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			OwnerID:     ownerID,
			Name:        "DevPanel",
			Description: "A development environment management system with real-time monitoring, service control, and comprehensive project management capabilities.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Project-Website/tree/main/backend/internal/devpanel",
			LiveURL:     "https://jadenrazo.dev/devpanel",
			Tags:        []string{"React", "Go", "WebSocket", "TypeScript", "Real-time Monitoring"},
			CreatedAt:   time.Now().Add(-15 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			OwnerID:     ownerID,
			Name:        "Messaging Platform",
			Description: "A real-time messaging platform with WebSocket integration, file attachments, reactions, and modern UI similar to Discord.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Project-Website/tree/main/backend/internal/messaging",
			LiveURL:     "https://jadenrazo.dev/messaging",
			Tags:        []string{"React", "WebSocket", "Go", "TypeScript", "Real-time Chat"},
			CreatedAt:   time.Now().Add(-10 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			OwnerID:     ownerID,
			Name:        "URL Shortener Service",
			Description: "A high-performance URL shortening service with analytics, custom short codes, and comprehensive statistics tracking.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Project-Website/tree/main/backend/internal/urlshortener",
			LiveURL:     "https://jadenrazo.dev/s/",
			Tags:        []string{"Go", "PostgreSQL", "Analytics", "REST API", "Microservice"},
			CreatedAt:   time.Now().Add(-5 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			OwnerID:     ownerID,
			Name:        "Code Statistics Tracker",
			Description: "Automated system for tracking lines of code across projects with scheduled updates and API integration.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Project-Website/tree/main/scripts",
			LiveURL:     "https://jadenrazo.dev/api/v1/code/stats",
			Tags:        []string{"Go", "Automation", "CLI", "Statistics", "CRON"},
			CreatedAt:   time.Now().Add(-3 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			OwnerID:     ownerID,
			Name:        "Showers Auto Detail",
			Description: "A mobile-first auto detailing booking platform with instant quotes, online booking, Square payment integration, before/after gallery, and admin dashboard with 2FA authentication.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/showersautodetail",
			LiveURL:     "https://showersautodetail.com",
			Tags:        []string{"Astro", "React", "Node.js", "PostgreSQL", "Tailwind CSS", "Square API", "Docker"},
			CreatedAt:   time.Now().Add(-1 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			OwnerID:     ownerID,
			Name:        "WeenieSMP - Full-Stack Minecraft Server Ecosystem",
			Description: "A production-grade Minecraft server ecosystem featuring a Vue 3 website with Tebex e-commerce, Go backend microservices for real-time statistics and leaderboards, MariaDB integration with custom Minecraft plugin data (playtime tracking, land claims, bounty systems), Redis caching layer (60s TTL), multi-domain Nginx routing with TLS 1.3, and comprehensive Docker deployment infrastructure serving 12,000+ registered players.",
			Status:      project.StatusActive,
			RepoURL:     "https://github.com/JadenRazo/Project-Website/tree/main/weeniesmp",
			LiveURL:     "https://weeniesmp.net",
			Tags:        []string{"Vue 3", "TypeScript", "Go", "MariaDB", "Redis", "Pinia", "Tebex API", "Docker", "Nginx", "Microservices", "TLS", "Tailwind CSS"},
			CreatedAt:   time.Now().Add(-2 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}

	for _, proj := range seedProjects {
		s.projects[proj.ID.String()] = proj
	}
}

func (s *MemoryProjectService) Create(ctx context.Context, proj *project.Project) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if proj == nil {
		return project.ErrInvalidProject
	}

	if proj.Name == "" {
		return fmt.Errorf("%w: name is required", project.ErrInvalidProject)
	}

	// Set defaults
	now := time.Now()
	proj.ID = uuid.New()
	proj.CreatedAt = now
	proj.UpdatedAt = now

	if proj.Status == "" {
		proj.Status = project.StatusDraft
	}

	if proj.OwnerID == uuid.Nil {
		proj.OwnerID = uuid.New() // Default owner for testing
	}

	s.projects[proj.ID.String()] = proj
	return nil
}

func (s *MemoryProjectService) Get(ctx context.Context, id uuid.UUID) (*project.Project, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid project ID", project.ErrInvalidProject)
	}

	proj, exists := s.projects[id.String()]
	if !exists {
		return nil, project.ErrProjectNotFound
	}

	return proj, nil
}

func (s *MemoryProjectService) GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]*project.Project, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var result []*project.Project
	for _, proj := range s.projects {
		if proj.OwnerID == ownerID {
			result = append(result, proj)
		}
	}

	return result, nil
}

func (s *MemoryProjectService) Update(ctx context.Context, proj *project.Project) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if proj == nil {
		return project.ErrInvalidProject
	}

	if proj.ID == uuid.Nil {
		return fmt.Errorf("%w: project ID is required", project.ErrInvalidProject)
	}

	existing, exists := s.projects[proj.ID.String()]
	if !exists {
		return project.ErrProjectNotFound
	}

	// Preserve creation date
	proj.CreatedAt = existing.CreatedAt
	proj.UpdatedAt = time.Now()

	s.projects[proj.ID.String()] = proj
	return nil
}

func (s *MemoryProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid project ID", project.ErrInvalidProject)
	}

	if _, exists := s.projects[id.String()]; !exists {
		return project.ErrProjectNotFound
	}

	delete(s.projects, id.String())
	return nil
}

func (s *MemoryProjectService) List(ctx context.Context, filter *project.Filter, pagination *project.Pagination) ([]*project.Project, int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if pagination == nil {
		pagination = &project.Pagination{
			Page:     1,
			PageSize: 10,
		}
	}

	if pagination.Page < 1 {
		pagination.Page = 1
	}

	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 10
	}

	if filter == nil {
		filter = &project.Filter{}
	}

	var filteredProjects []*project.Project

	for _, proj := range s.projects {
		// Apply filters
		if filter.Status != nil && proj.Status != *filter.Status {
			continue
		}

		if filter.OwnerID != uuid.Nil && proj.OwnerID != filter.OwnerID {
			continue
		}

		if filter.Tag != "" {
			hasTag := false
			for _, tag := range proj.Tags {
				if tag == filter.Tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		if filter.Search != "" {
			// Case-insensitive search in name and description
			if !containsIgnoreCase(proj.Name, filter.Search) && !containsIgnoreCase(proj.Description, filter.Search) {
				continue
			}
		}

		filteredProjects = append(filteredProjects, proj)
	}

	total := len(filteredProjects)

	// Apply pagination
	start := (pagination.Page - 1) * pagination.PageSize
	end := start + pagination.PageSize

	if start > total {
		return []*project.Project{}, total, nil
	}

	if end > total {
		end = total
	}

	return filteredProjects[start:end], total, nil
}

func containsIgnoreCase(str, substr string) bool {
	// Convert both strings to lowercase for case-insensitive comparison
	strLower := strings.ToLower(str)
	substrLower := strings.ToLower(substr)
	return strings.Contains(strLower, substrLower)
}
