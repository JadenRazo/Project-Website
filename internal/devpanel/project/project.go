package project

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Route represents an API route for a project
type Route struct {
	Path        string   `json:"path" yaml:"path"`
	Description string   `json:"description" yaml:"description"`
	Methods     []string `json:"methods" yaml:"methods"`
}

// Binary represents a binary executable for a project
type Binary struct {
	Path         string            `json:"path" yaml:"path"`
	BuildCommand string            `json:"buildCommand" yaml:"buildCommand"`
	RunCommand   string            `json:"runCommand" yaml:"runCommand"`
	Environment  map[string]string `json:"environment" yaml:"environment"`
	Process      *os.Process       `json:"-" yaml:"-"`
}

// Project represents a deployable project
type Project struct {
	ID           string   `json:"id" yaml:"id"`
	Name         string   `json:"name" yaml:"name"`
	Description  string   `json:"description" yaml:"description"`
	Enabled      bool     `json:"enabled" yaml:"enabled"`
	Routes       []Route  `json:"routes" yaml:"routes"`
	Binaries     []Binary `json:"binaries" yaml:"binaries"`
	Dependencies []string `json:"dependencies" yaml:"dependencies"`
	Status       string   `json:"status" yaml:"-"`
	LastUpdated  string   `json:"lastUpdated" yaml:"-"`
}

// ProjectList represents a list of projects
type ProjectList struct {
	Projects []Project `json:"projects" yaml:"projects"`
}

// Enable enables a project
func (p *Project) Enable() error {
	p.Enabled = true
	p.Status = "Starting"
	p.LastUpdated = time.Now().Format(time.RFC3339)
	return nil
}

// Disable disables a project
func (p *Project) Disable() error {
	// Stop any running processes first
	err := p.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop project: %v", err)
	}

	p.Enabled = false
	p.Status = "Stopped"
	p.LastUpdated = time.Now().Format(time.RFC3339)
	return nil
}

// Start starts a project
func (p *Project) Start() error {
	if !p.Enabled {
		return fmt.Errorf("cannot start disabled project: %s", p.ID)
	}

	// First build the project
	if err := p.Build(); err != nil {
		p.Status = "Build Failed"
		return fmt.Errorf("failed to build project: %v", err)
	}

	// Then start the binaries
	for i, binary := range p.Binaries {
		// Split the run command into command and arguments
		parts := strings.Fields(binary.RunCommand)
		if len(parts) == 0 {
			p.Status = "Invalid Run Command"
			return fmt.Errorf("invalid run command for binary: %s", binary.Path)
		}

		cmd := exec.Command(parts[0], parts[1:]...)

		// Set environment variables
		env := os.Environ()
		for k, v := range binary.Environment {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env

		// Start the process
		if err := cmd.Start(); err != nil {
			p.Status = "Start Failed"
			return fmt.Errorf("failed to start binary: %v", err)
		}

		// Store the process reference
		p.Binaries[i].Process = cmd.Process
	}

	p.Status = "Running"
	p.LastUpdated = time.Now().Format(time.RFC3339)
	return nil
}

// Stop stops a project
func (p *Project) Stop() error {
	var lastErr error

	for i, binary := range p.Binaries {
		if binary.Process != nil {
			// Kill the process
			if err := binary.Process.Kill(); err != nil {
				lastErr = fmt.Errorf("failed to kill process: %v", err)
			}

			// Clear the process reference
			p.Binaries[i].Process = nil
		}
	}

	if lastErr != nil {
		p.Status = "Stop Failed"
		return lastErr
	}

	p.Status = "Stopped"
	p.LastUpdated = time.Now().Format(time.RFC3339)
	return nil
}

// Build builds a project's binaries
func (p *Project) Build() error {
	for _, binary := range p.Binaries {
		// Split the build command into command and arguments
		parts := strings.Fields(binary.BuildCommand)
		if len(parts) == 0 {
			return fmt.Errorf("invalid build command for binary: %s", binary.Path)
		}

		cmd := exec.Command(parts[0], parts[1:]...)

		// Execute the build command
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("build failed: %v, output: %s", err, string(output))
		}
	}

	return nil
}

// Restart restarts a project
func (p *Project) Restart() error {
	if err := p.Stop(); err != nil {
		return fmt.Errorf("failed to stop project: %v", err)
	}

	return p.Start()
}
