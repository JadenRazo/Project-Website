package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"

	"jadenrazo.dev/internal/devpanel/project"
)

// ProjectRepository defines the interface for project data operations
type ProjectRepository interface {
	// GetAll returns all projects
	GetAll() ([]project.Project, error)

	// GetByID returns a project by its ID
	GetByID(id string) (*project.Project, error)

	// Create creates a new project
	Create(p project.Project) error

	// Update updates an existing project
	Update(p project.Project) error

	// Delete deletes a project by its ID
	Delete(id string) error

	// SaveConfig saves the current configuration to disk
	SaveConfig() error
}

// FileProjectRepository implements ProjectRepository using a YAML file
type FileProjectRepository struct {
	configPath string
	projects   map[string]*project.Project
	mutex      sync.RWMutex
}

// NewFileProjectRepository creates a new file-based project repository
func NewFileProjectRepository(configPath string) (*FileProjectRepository, error) {
	repo := &FileProjectRepository{
		configPath: configPath,
		projects:   make(map[string]*project.Project),
		mutex:      sync.RWMutex{},
	}

	// Load initial configuration
	if err := repo.loadConfig(); err != nil {
		return nil, err
	}

	return repo, nil
}

// loadConfig loads the project configuration from the YAML file
func (r *FileProjectRepository) loadConfig() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if the file exists
	if _, err := os.Stat(r.configPath); os.IsNotExist(err) {
		// If the file doesn't exist, create an empty configuration
		return r.SaveConfig()
	}

	// Read the file
	data, err := ioutil.ReadFile(r.configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Parse the YAML
	var projectList project.ProjectList
	if err := yaml.Unmarshal(data, &projectList); err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	// Initialize the projects map
	r.projects = make(map[string]*project.Project)
	for i := range projectList.Projects {
		p := &projectList.Projects[i]
		r.projects[p.ID] = p

		// Set initial status
		if p.Enabled {
			p.Status = "Ready"
		} else {
			p.Status = "Disabled"
		}
	}

	return nil
}

// GetAll returns all projects
func (r *FileProjectRepository) GetAll() ([]project.Project, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	projects := make([]project.Project, 0, len(r.projects))
	for _, p := range r.projects {
		projects = append(projects, *p)
	}

	return projects, nil
}

// GetByID returns a project by its ID
func (r *FileProjectRepository) GetByID(id string) (*project.Project, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	p, exists := r.projects[id]
	if !exists {
		return nil, fmt.Errorf("project not found: %s", id)
	}

	return p, nil
}

// Create creates a new project
func (r *FileProjectRepository) Create(p project.Project) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if a project with the same ID already exists
	if _, exists := r.projects[p.ID]; exists {
		return fmt.Errorf("project with ID %s already exists", p.ID)
	}

	// Initialize status
	if p.Enabled {
		p.Status = "Ready"
	} else {
		p.Status = "Disabled"
	}

	// Add to map
	r.projects[p.ID] = &p

	// Save to disk
	return r.SaveConfig()
}

// Update updates an existing project
func (r *FileProjectRepository) Update(p project.Project) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if the project exists
	if _, exists := r.projects[p.ID]; !exists {
		return fmt.Errorf("project not found: %s", p.ID)
	}

	// Update in map
	r.projects[p.ID] = &p

	// Save to disk
	return r.SaveConfig()
}

// Delete deletes a project by its ID
func (r *FileProjectRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if the project exists
	if _, exists := r.projects[id]; !exists {
		return fmt.Errorf("project not found: %s", id)
	}

	// Delete from map
	delete(r.projects, id)

	// Save to disk
	return r.SaveConfig()
}

// SaveConfig saves the current configuration to disk
func (r *FileProjectRepository) SaveConfig() error {
	// Create a project list from the map
	projects := make([]project.Project, 0, len(r.projects))
	for _, p := range r.projects {
		projects = append(projects, *p)
	}

	projectList := project.ProjectList{
		Projects: projects,
	}

	// Marshal to YAML
	data, err := yaml.Marshal(projectList)
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(r.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	// Write to file
	if err := ioutil.WriteFile(r.configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}
