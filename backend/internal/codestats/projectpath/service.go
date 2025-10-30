package projectpath

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/google/uuid"
)

// CodeStatsUpdater interface to avoid circular dependency
type CodeStatsUpdater interface {
	UpdateStats() error
}

type ProjectPathService struct {
	repo        Repository
	statsUpdater CodeStatsUpdater
}

func NewService(repo Repository) *ProjectPathService {
	return &ProjectPathService{repo: repo}
}

func NewServiceWithStatsUpdater(repo Repository, statsUpdater CodeStatsUpdater) *ProjectPathService {
	return &ProjectPathService{
		repo:        repo,
		statsUpdater: statsUpdater,
	}
}

func (s *ProjectPathService) CreateProjectPath(ctx context.Context, path *ProjectPath) error {
	path.ID = uuid.New()
	path.CreatedAt = time.Now()
	path.UpdatedAt = time.Now()

	if err := s.validateProjectPath(path); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.repo.Create(ctx, path); err != nil {
		return err
	}

	// Trigger async code stats update
	s.triggerStatsUpdate("created project path", path.Name)
	
	return nil
}

func (s *ProjectPathService) GetProjectPath(ctx context.Context, id uuid.UUID) (*ProjectPath, error) {
	return s.repo.Get(ctx, id)
}

func (s *ProjectPathService) UpdateProjectPath(ctx context.Context, path *ProjectPath) error {
	existing, err := s.repo.Get(ctx, path.ID)
	if err != nil {
		return err
	}

	path.CreatedAt = existing.CreatedAt
	path.UpdatedAt = time.Now()

	if err := s.validateProjectPath(path); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.repo.Update(ctx, path); err != nil {
		return err
	}

	// Trigger async code stats update if active status changed or path details changed
	shouldUpdate := existing.IsActive != path.IsActive || 
		existing.Path != path.Path || 
		fmt.Sprintf("%v", existing.ExcludePatterns) != fmt.Sprintf("%v", path.ExcludePatterns)
	
	if shouldUpdate {
		s.triggerStatsUpdate("updated project path", path.Name)
	}
	
	return nil
}

func (s *ProjectPathService) DeleteProjectPath(ctx context.Context, id uuid.UUID) error {
	// Get the path name before deletion for logging
	path, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Trigger async code stats update
	s.triggerStatsUpdate("deleted project path", path.Name)
	
	return nil
}

func (s *ProjectPathService) ListProjectPaths(ctx context.Context, filter *Filter, pagination *Pagination) ([]*ProjectPath, int, error) {
	if pagination == nil {
		pagination = &Pagination{Page: 1, PageSize: 10}
	}

	if pagination.Page < 1 {
		pagination.Page = 1
	}

	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 10
	}

	return s.repo.List(ctx, filter, pagination)
}

func (s *ProjectPathService) GetActiveProjectPaths(ctx context.Context) ([]*ProjectPath, error) {
	return s.repo.GetActive(ctx)
}

func (s *ProjectPathService) validateProjectPath(path *ProjectPath) error {
	if path.Name == "" {
		return fmt.Errorf("name is required")
	}

	if path.Path == "" {
		return fmt.Errorf("path is required")
	}

	if _, err := os.Stat(path.Path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path.Path)
	}

	if path.ExcludePatterns == nil {
		path.ExcludePatterns = make([]string, 0)
	}

	return nil
}

// triggerStatsUpdate asynchronously updates code statistics
func (s *ProjectPathService) triggerStatsUpdate(action, pathName string) {
	if s.statsUpdater == nil {
		logger.Warn("Code stats updater not configured, skipping stats update", 
			"action", action, "path", pathName)
		return
	}

	// Run stats update in a goroutine to avoid blocking the API response
	go func() {
		logger.Info("Triggering code stats update", 
			"action", action, "path", pathName)
		
		if err := s.statsUpdater.UpdateStats(); err != nil {
			logger.Error("Failed to update code statistics after project path change",
				"action", action, "path", pathName, "error", err)
		} else {
			logger.Info("Successfully updated code statistics",
				"action", action, "path", pathName)
		}
	}()
}