package codestats

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/codestats/projectpath"
	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"gorm.io/gorm"
)

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Service struct {
	db              *gorm.DB
	projectPathRepo projectpath.Repository
	updateMutex     sync.Mutex
	updating        bool
	ticker          *time.Ticker
	stopChan        chan struct{}
	stopped         bool
}

type CodeStats struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Languages    []Language `gorm:"serializer:json" json:"languages"`
	TotalLines   int64      `json:"total_lines"`
	TotalFiles   int        `json:"total_files"`
	TotalBlanks  int64      `json:"total_blanks"`
	TotalCode    int64      `json:"total_code"`
	TotalComment int64      `json:"total_comment"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type Language struct {
	Name     string `json:"name"`
	Files    int    `json:"files"`
	Lines    int64  `json:"lines"`
	Code     int64  `json:"code"`
	Comments int64  `json:"comments"`
	Blanks   int64  `json:"blanks"`
}

func NewService(db *gorm.DB, projectPathRepo projectpath.Repository) *Service {
	return &Service{
		db:              db,
		projectPathRepo: projectPathRepo,
		stopChan:        make(chan struct{}),
	}
}

func (s *Service) GetLatestStats() (*CodeStats, error) {
	var stats CodeStats
	err := s.db.Order("created_at DESC").First(&stats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &stats, nil
}

func (s *Service) UpdateStats() error {
	// Check if update is already in progress
	s.updateMutex.Lock()
	if s.updating {
		s.updateMutex.Unlock()
		logger.Info("Code stats update already in progress, skipping")
		return nil
	}
	s.updating = true
	s.updateMutex.Unlock()

	// Ensure we reset the updating flag when done
	defer func() {
		s.updateMutex.Lock()
		s.updating = false
		s.updateMutex.Unlock()
	}()

	logger.Info("Running tokei to collect code statistics")

	// Check if tokei binary exists before proceeding
	tokeiPath := "/root/.cargo/bin/tokei"
	if _, err := os.Stat(tokeiPath); os.IsNotExist(err) {
		return fmt.Errorf("tokei binary not found at %s. Please install tokei to enable code statistics", tokeiPath)
	}

	// Get active project paths from database
	ctx := context.Background()
	projectPaths, err := s.projectPathRepo.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("database error while retrieving project paths: %w. Check database connectivity and configure project paths in the DevPanel", err)
	}

	if len(projectPaths) == 0 {
		return fmt.Errorf("no active project paths configured. Please add project paths in the DevPanel")
	}

	// Build tokei command with all project paths and exclusions
	args := []string{"/root/.cargo/bin/tokei", "--output", "json"}

	// Add exclusion patterns from database
	excludePatterns := make(map[string]bool)
	for _, project := range projectPaths {
		for _, exclude := range project.ExcludePatterns {
			excludePatterns[exclude] = true
		}
	}

	for pattern := range excludePatterns {
		args = append(args, "--exclude", pattern)
	}

	// Add all project paths from database
	for _, project := range projectPaths {
		args = append(args, project.Path)
		logger.Info("Adding project to scan", "name", project.Name, "path", project.Path)
	}

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderrMsg := string(exitErr.Stderr)
			logger.Error("Tokei execution failed", "stderr", stderrMsg, "exit_code", exitErr.ExitCode())
			return fmt.Errorf("tokei command failed with exit code %d: %s. Check if all project paths are accessible", exitErr.ExitCode(), stderrMsg)
		}
		return fmt.Errorf("failed to execute tokei command: %w. Verify tokei installation and project path permissions", err)
	}

	var tokeiOutput map[string]interface{}
	if err := json.Unmarshal(output, &tokeiOutput); err != nil {
		logger.Error("Failed to parse tokei JSON output", "output_preview", string(output[:min(len(output), 500)]))
		return fmt.Errorf("failed to parse tokei JSON output: %w. The output may be corrupted or in unexpected format", err)
	}

	stats := s.parseTokeniOutput(tokeiOutput)
	stats.UpdatedAt = time.Now()
	stats.CreatedAt = time.Now()

	if err := s.db.Create(&stats).Error; err != nil {
		logger.Error("Database error while saving code statistics", "total_lines", stats.TotalLines, "db_error", err)
		return fmt.Errorf("database error while saving code statistics: %w. Check database connectivity and permissions", err)
	}

	// Also save to file for backward compatibility
	if err := s.saveToFile(stats); err != nil {
		logger.Error("Failed to save stats to file", "error", err)
	}

	logger.Info("Code statistics updated successfully",
		"total_files", stats.TotalFiles,
		"total_lines", stats.TotalLines,
		"total_code", stats.TotalCode)

	return nil
}

func (s *Service) parseTokeniOutput(output map[string]interface{}) *CodeStats {
	stats := &CodeStats{
		Languages: []Language{},
	}

	for lang, data := range output {
		if lang == "Total" {
			continue
		}

		if langData, ok := data.(map[string]interface{}); ok {
			language := Language{
				Name: lang,
			}

			// Count files from reports array
			if reports, ok := langData["reports"].([]interface{}); ok {
				language.Files = len(reports)
			}

			if code, ok := langData["code"].(float64); ok {
				language.Code = int64(code)
			}
			if comments, ok := langData["comments"].(float64); ok {
				language.Comments = int64(comments)
			}
			if blanks, ok := langData["blanks"].(float64); ok {
				language.Blanks = int64(blanks)
			}

			// Calculate total lines as code + comments + blanks
			language.Lines = language.Code + language.Comments + language.Blanks

			stats.Languages = append(stats.Languages, language)
			stats.TotalFiles += language.Files
			stats.TotalLines += language.Lines
			stats.TotalCode += language.Code
			stats.TotalComment += language.Comments
			stats.TotalBlanks += language.Blanks
		}
	}

	return stats
}

func (s *Service) saveToFile(stats *CodeStats) error {
	// Always save to the main project's frontend directory
	publicPath := "/main/Project-Website/frontend/public/code_stats.json"

	// Save in the format expected by frontend (totalLines field)
	fileData := map[string]interface{}{
		"totalLines": stats.TotalLines,
	}

	data, err := json.MarshalIndent(fileData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(publicPath, data, 0644)
}


func (s *Service) StartPeriodicUpdate(interval time.Duration) {
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Prevent multiple periodic updates
	if s.ticker != nil {
		logger.Warn("Periodic update already running, stopping previous instance")
		s.stopPeriodicUpdate()
	}

	s.ticker = time.NewTicker(interval)
	s.stopped = false

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic in periodic update goroutine", "panic", r)
			}
		}()

		// Run immediately on start
		if err := s.UpdateStats(); err != nil {
			logger.Error("Failed to update code stats on startup", "error", err)
		}

		for {
			select {
			case <-s.ticker.C:
				if err := s.UpdateStats(); err != nil {
					logger.Error("Failed to update code stats during periodic update", "error", err)
				}
			case <-s.stopChan:
				logger.Info("Stopping periodic code stats updates")
				return
			}
		}
	}()

	logger.Info("Started periodic code stats updates", "interval", interval)
}

// StopPeriodicUpdate stops the periodic update goroutine
func (s *Service) StopPeriodicUpdate() {
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()
	s.stopPeriodicUpdate()
}

// stopPeriodicUpdate is the internal method that must be called with mutex held
func (s *Service) stopPeriodicUpdate() {
	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}

	if !s.stopped {
		close(s.stopChan)
		s.stopped = true
		s.stopChan = make(chan struct{}) // Create new channel for potential restart
	}
}
