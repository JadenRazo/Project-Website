package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/certification"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) certification.Repository {
	return &gormRepository{db: db}
}

// Certifications

func (r *gormRepository) CreateCertification(cert *certification.Certification) error {
	return r.db.Create(cert).Error
}

func (r *gormRepository) GetCertificationByID(id uuid.UUID) (*certification.Certification, error) {
	var cert certification.Certification
	err := r.db.Preload("Category").First(&cert, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *gormRepository) GetAllCertifications(includeHidden bool) ([]certification.Certification, error) {
	var certs []certification.Certification
	query := r.db.Preload("Category")

	if !includeHidden {
		query = query.Where("is_visible = ?", true)
	}

	err := query.Order("sort_order ASC, created_at DESC").Find(&certs).Error
	return certs, err
}

func (r *gormRepository) GetVisibleCertifications() ([]certification.Certification, error) {
	var certs []certification.Certification
	err := r.db.Preload("Category").
		Where("is_visible = ?", true).
		Where("deleted_at IS NULL").
		Order("is_featured DESC, sort_order ASC, issue_date DESC").
		Find(&certs).Error
	return certs, err
}

func (r *gormRepository) UpdateCertification(cert *certification.Certification) error {
	return r.db.Save(cert).Error
}

func (r *gormRepository) DeleteCertification(id uuid.UUID) error {
	return r.db.Delete(&certification.Certification{}, "id = ?", id).Error
}

// Categories

func (r *gormRepository) CreateCategory(cat *certification.Category) error {
	return r.db.Create(cat).Error
}

func (r *gormRepository) GetCategoryByID(id uuid.UUID) (*certification.Category, error) {
	var cat certification.Category
	err := r.db.First(&cat, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *gormRepository) GetAllCategories(includeHidden bool) ([]certification.Category, error) {
	var cats []certification.Category
	query := r.db.Model(&certification.Category{})

	if !includeHidden {
		query = query.Where("is_visible = ?", true)
	}

	err := query.Order("sort_order ASC, name ASC").Find(&cats).Error
	return cats, err
}

func (r *gormRepository) GetVisibleCategories() ([]certification.Category, error) {
	var cats []certification.Category
	err := r.db.Where("is_visible = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&cats).Error
	return cats, err
}

func (r *gormRepository) UpdateCategory(cat *certification.Category) error {
	return r.db.Save(cat).Error
}

func (r *gormRepository) DeleteCategory(id uuid.UUID) error {
	return r.db.Delete(&certification.Category{}, "id = ?", id).Error
}
