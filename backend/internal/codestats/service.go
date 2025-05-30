package codestats

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"gorm.io/gorm"
)

type Service struct {
	db       *gorm.DB
	projects []ProjectConfig
}

type Config struct {
	Projects       []ProjectConfig
	UpdateInterval time.Duration
}

type ProjectConfig struct {
	Name    string   `yaml:"name"`
	Path    string   `yaml:"path"`
	Exclude []string `yaml:"exclude,omitempty"`
}

type CodeStats struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Languages    []Language `gorm:"serializer:json" json:"languages"`
	TotalLines   int64     `json:"total_lines"`
	TotalFiles   int       `json:"total_files"`
	TotalBlanks  int64     `json:"total_blanks"`
	TotalCode    int64     `json:"total_code"`
	TotalComment int64     `json:"total_comment"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Language struct {
	Name     string `json:"name"`
	Files    int    `json:"files"`
	Lines    int64  `json:"lines"`
	Code     int64  `json:"code"`
	Comments int64  `json:"comments"`
	Blanks   int64  `json:"blanks"`
}

func NewService(db *gorm.DB, cfg Config) *Service {
	return &Service{
		db:       db,
		projects: cfg.Projects,
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
	logger.Info("Running tokei to collect code statistics")
	
	// Build tokei command with all project paths and exclusions
	args := []string{"/root/.cargo/bin/tokei", "--output", "json"}
	
	// Add exclusion patterns
	excludePatterns := make(map[string]bool)
	for _, project := range s.projects {
		for _, exclude := range project.Exclude {
			excludePatterns[exclude] = true
		}
	}
	
	for pattern := range excludePatterns {
		args = append(args, "--exclude", pattern)
	}
	
	// Add all project paths
	for _, project := range s.projects {
		args = append(args, project.Path)
		logger.Info("Adding project to scan", "name", project.Name, "path", project.Path)
	}
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			logger.Error("Tokei failed", "stderr", string(exitErr.Stderr))
		}
		return fmt.Errorf("failed to run tokei: %w", err)
	}

	var tokeiOutput map[string]interface{}
	if err := json.Unmarshal(output, &tokeiOutput); err != nil {
		return fmt.Errorf("failed to parse tokei output: %w", err)
	}

	stats := s.parseTokeniOutput(tokeiOutput)
	stats.UpdatedAt = time.Now()
	stats.CreatedAt = time.Now()

	if err := s.db.Create(&stats).Error; err != nil {
		return fmt.Errorf("failed to save stats: %w", err)
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
	ticker := time.NewTicker(interval)
	go func() {
		// Run immediately on start
		if err := s.UpdateStats(); err != nil {
			logger.Error("Failed to update code stats", "error", err)
		}

		for range ticker.C {
			if err := s.UpdateStats(); err != nil {
				logger.Error("Failed to update code stats", "error", err)
			}
		}
	}()
}