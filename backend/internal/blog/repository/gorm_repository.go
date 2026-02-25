package repository

import (
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/blog"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) blog.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(post *blog.Post) error {
	return r.db.Create(post).Error
}

func (r *gormRepository) GetByID(id uuid.UUID) (*blog.Post, error) {
	var post blog.Post
	err := r.db.First(&post, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *gormRepository) GetBySlug(slug string) (*blog.Post, error) {
	var post blog.Post
	err := r.db.Where("slug = ?", slug).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *gormRepository) List(params blog.ListParams) (*blog.ListResult, error) {
	var posts []blog.Post
	var total int64

	query := r.db.Model(&blog.Post{})

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Tag != "" {
		query = query.Where("? = ANY(tags)", params.Tag)
	}

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("title ILIKE ? OR excerpt ILIKE ?", search, search)
	}

	query.Count(&total)

	offset := (params.Page - 1) * params.PageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))

	return &blog.ListResult{
		Posts:      posts,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *gormRepository) ListPublished(params blog.ListParams) (*blog.ListResult, error) {
	var posts []blog.Post
	var total int64

	query := r.db.Model(&blog.Post{}).
		Where("status = ? AND is_visible = ?", "published", true)

	if params.Tag != "" {
		query = query.Where("? = ANY(tags)", params.Tag)
	}

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("title ILIKE ? OR excerpt ILIKE ?", search, search)
	}

	query.Count(&total)

	offset := (params.Page - 1) * params.PageSize
	err := query.Order("published_at DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))

	return &blog.ListResult{
		Posts:      posts,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *gormRepository) GetFeatured(limit int) ([]blog.Post, error) {
	var posts []blog.Post
	err := r.db.Where("status = ? AND is_visible = ? AND is_featured = ?", "published", true, true).
		Order("published_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

func (r *gormRepository) Update(post *blog.Post) error {
	return r.db.Save(post).Error
}

func (r *gormRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&blog.Post{}, "id = ?", id).Error
}

func (r *gormRepository) IncrementViewCount(id uuid.UUID) error {
	return r.db.Model(&blog.Post{}).
		Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}
