package certification

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Certifications

func (s *service) CreateCertification(cert *Certification) error {
	if cert.Name == "" || cert.Issuer == "" {
		return errors.New("name and issuer are required")
	}

	if cert.IssueDate.IsZero() {
		return errors.New("issue date is required")
	}

	if cert.ExpiryDate != nil && cert.ExpiryDate.Before(cert.IssueDate) {
		return errors.New("expiry date must be after issue date")
	}

	cert.ID = uuid.New()
	cert.CreatedAt = time.Now()
	cert.UpdatedAt = time.Now()

	return s.repo.CreateCertification(cert)
}

func (s *service) GetCertificationByID(id uuid.UUID) (*Certification, error) {
	return s.repo.GetCertificationByID(id)
}

func (s *service) GetAllCertifications(includeHidden bool) ([]Certification, error) {
	return s.repo.GetAllCertifications(includeHidden)
}

func (s *service) GetVisibleCertifications() ([]Certification, error) {
	return s.repo.GetVisibleCertifications()
}

func (s *service) UpdateCertification(id uuid.UUID, updates map[string]interface{}) error {
	cert, err := s.repo.GetCertificationByID(id)
	if err != nil {
		return err
	}

	// Update fields
	if name, ok := updates["name"].(string); ok && name != "" {
		cert.Name = name
	}

	if issuer, ok := updates["issuer"].(string); ok && issuer != "" {
		cert.Issuer = issuer
	}

	if credentialID, ok := updates["credential_id"].(string); ok {
		cert.CredentialID = &credentialID
	}

	if issueDate, ok := updates["issue_date"].(time.Time); ok {
		cert.IssueDate = issueDate
	}

	if expiryDate, ok := updates["expiry_date"].(time.Time); ok {
		cert.ExpiryDate = &expiryDate
	}

	if verificationURL, ok := updates["verification_url"].(string); ok {
		cert.VerificationURL = &verificationURL
	}

	if verificationText, ok := updates["verification_text"].(string); ok {
		cert.VerificationText = &verificationText
	}

	if badgeURL, ok := updates["badge_url"].(string); ok {
		cert.BadgeURL = &badgeURL
	}

	if description, ok := updates["description"].(string); ok {
		cert.Description = &description
	}

	if categoryID, ok := updates["category_id"].(string); ok {
		if categoryID == "" {
			cert.CategoryID = nil
		} else {
			catID, err := uuid.Parse(categoryID)
			if err == nil {
				cert.CategoryID = &catID
			}
		}
	}

	if isFeatured, ok := updates["is_featured"].(bool); ok {
		cert.IsFeatured = isFeatured
	}

	if isVisible, ok := updates["is_visible"].(bool); ok {
		cert.IsVisible = isVisible
	}

	if sortOrder, ok := updates["sort_order"].(float64); ok {
		cert.SortOrder = int(sortOrder)
	}

	cert.UpdatedAt = time.Now()

	return s.repo.UpdateCertification(cert)
}

func (s *service) DeleteCertification(id uuid.UUID) error {
	return s.repo.DeleteCertification(id)
}

// Categories

func (s *service) CreateCategory(cat *Category) error {
	if cat.Name == "" {
		return errors.New("category name is required")
	}

	cat.ID = uuid.New()
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()

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

	// Update fields
	if name, ok := updates["name"].(string); ok && name != "" {
		cat.Name = name
	}

	if description, ok := updates["description"].(string); ok {
		cat.Description = description
	}

	if sortOrder, ok := updates["sort_order"].(float64); ok {
		cat.SortOrder = int(sortOrder)
	}

	if isVisible, ok := updates["is_visible"].(bool); ok {
		cat.IsVisible = isVisible
	}

	cat.UpdatedAt = time.Now()

	return s.repo.UpdateCategory(cat)
}

func (s *service) DeleteCategory(id uuid.UUID) error {
	return s.repo.DeleteCategory(id)
}
