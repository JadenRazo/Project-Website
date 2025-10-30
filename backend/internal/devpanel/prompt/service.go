
package prompt

import (
	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Prompts

func (s *service) CreatePrompt(p *Prompt) error {
	// Validation can be added here
	return s.repo.CreatePrompt(p)
}

func (s *service) GetPromptByID(id uuid.UUID) (*Prompt, error) {
	return s.repo.GetPromptByID(id)
}

func (s *service) GetAllPrompts(includeHidden bool) ([]Prompt, error) {
	return s.repo.GetAllPrompts(includeHidden)
}

func (s *service) GetVisiblePrompts() ([]Prompt, error) {
	return s.repo.GetVisiblePrompts()
}

func (s *service) UpdatePrompt(id uuid.UUID, updates map[string]interface{}) error {
	p, err := s.repo.GetPromptByID(id)
	if err != nil {
		return err
	}

	// This is a simple way to update. A more robust solution would use a dedicated update struct and validation.
	if name, ok := updates["name"].(string); ok {
		p.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		p.Description = description
	}
	if prompt, ok := updates["prompt"].(string); ok {
		p.Prompt = prompt
	}
	if categoryID, ok := updates["category_id"].(string); ok {
		parsedID, err := uuid.Parse(categoryID)
		if err == nil {
			p.CategoryID = parsedID
		}
	}
	if isFeatured, ok := updates["is_featured"].(bool); ok {
		p.IsFeatured = isFeatured
	}
	if isVisible, ok := updates["is_visible"].(bool); ok {
		p.IsVisible = isVisible
	}
	if sortOrder, ok := updates["sort_order"].(float64); ok {
		p.SortOrder = int(sortOrder)
	}

	return s.repo.UpdatePrompt(p)
}

func (s *service) DeletePrompt(id uuid.UUID) error {
	return s.repo.DeletePrompt(id)
}

// Categories

func (s *service) CreateCategory(cat *Category) error {
	return s.repo.CreateCategory(cat)
}

func (s *service) GetCategoryByID(id uuid.UUID) (*Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *service) GetAllCategories(includeHidden bool) ([]Category, error) {
	return s.repo.GetAllCategories(includeHidden)
}

func (s *service) GetVisibleCategories() ([]Category, error) {
	return s.repo.GetVisibleCategories()
}

func (s *service) UpdateCategory(id uuid.UUID, updates map[string]interface{}) error {
	cat, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return err
	}

	if name, ok := updates["name"].(string); ok {
		cat.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		cat.Description = description
	}
	if isVisible, ok := updates["is_visible"].(bool); ok {
		cat.IsVisible = isVisible
	}
	if sortOrder, ok := updates["sort_order"].(float64); ok {
		cat.SortOrder = int(sortOrder)
	}

	return s.repo.UpdateCategory(cat)
}

func (s *service) DeleteCategory(id uuid.UUID) error {
	return s.repo.DeleteCategory(id)
}
