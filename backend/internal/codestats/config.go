package codestats

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ConfigFile represents the structure of the codestats.yaml file
type ConfigFile struct {
	Projects       []ProjectConfig `yaml:"projects"`
	UpdateInterval string          `yaml:"updateInterval"`
	Output         OutputConfig    `yaml:"output"`
}

type OutputConfig struct {
	IncludeLanguageBreakdown bool `yaml:"includeLanguageBreakdown"`
	IncludeFileCount         bool `yaml:"includeFileCount"`
}

// LoadConfig loads the code stats configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var configFile ConfigFile
	if err := yaml.Unmarshal(data, &configFile); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Parse update interval
	updateInterval, err := time.ParseDuration(configFile.UpdateInterval)
	if err != nil {
		updateInterval = 1 * time.Hour // Default to 1 hour
	}

	return &Config{
		Projects:       configFile.Projects,
		UpdateInterval: updateInterval,
	}, nil
}

// DefaultConfig returns a default configuration if no config file is found
func DefaultConfig() *Config {
	return &Config{
		Projects: []ProjectConfig{
			{
				Name: "Quiz Bot",
				Path: "/quiz_bot",
			},
			{
				Name: "Project Website",
				Path: "/main/Project-Website",
				Exclude: []string{"build", "node_modules", "logs", "bin"},
			},
		},
		UpdateInterval: 1 * time.Hour,
	}
}