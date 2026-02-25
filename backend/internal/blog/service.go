package blog

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func slugify(title string) string {
	s := strings.ToLower(title)
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' {
			return r
		}
		return -1
	}, s)
	reg := regexp.MustCompile(`\s+`)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func calculateReadTime(content string) int {
	words := len(strings.Fields(content))
	minutes := words / 200
	if minutes < 1 {
		return 1
	}
	return minutes
}

func (s *service) CreatePost(post *Post) error {
	if post.Title == "" {
		return errors.New("title is required")
	}

	post.ID = uuid.New()
	post.Slug = slugify(post.Title)
	post.ReadTimeMinutes = calculateReadTime(post.Content)
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	if post.Status == "published" && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}

	return s.repo.Create(post)
}

func (s *service) GetByID(id uuid.UUID) (*Post, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetBySlug(slug string) (*Post, error) {
	return s.repo.GetBySlug(slug)
}

func (s *service) List(params ListParams) (*ListResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 50 {
		params.PageSize = 10
	}
	return s.repo.List(params)
}

func (s *service) ListPublished(params ListParams) (*ListResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 50 {
		params.PageSize = 10
	}
	return s.repo.ListPublished(params)
}

func (s *service) GetFeatured(limit int) ([]Post, error) {
	if limit < 1 || limit > 10 {
		limit = 3
	}
	return s.repo.GetFeatured(limit)
}

func (s *service) UpdatePost(id uuid.UUID, updates map[string]interface{}) error {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if title, ok := updates["title"].(string); ok && title != "" {
		post.Title = title
		post.Slug = slugify(title)
	}

	if slug, ok := updates["slug"].(string); ok && slug != "" {
		post.Slug = slug
	}

	if content, ok := updates["content"].(string); ok {
		post.Content = content
		post.ReadTimeMinutes = calculateReadTime(content)
	}

	if excerpt, ok := updates["excerpt"].(string); ok {
		post.Excerpt = excerpt
	}

	if featuredImage, ok := updates["featured_image"].(string); ok {
		post.FeaturedImage = featuredImage
	}

	if tags, ok := updates["tags"].([]interface{}); ok {
		strTags := make([]string, 0, len(tags))
		for _, t := range tags {
			if str, ok := t.(string); ok {
				strTags = append(strTags, str)
			}
		}
		post.Tags = strTags
	}

	if status, ok := updates["status"].(string); ok {
		if status == "published" && post.Status != "published" && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
		post.Status = status
	}

	if isFeatured, ok := updates["is_featured"].(bool); ok {
		post.IsFeatured = isFeatured
	}

	if isVisible, ok := updates["is_visible"].(bool); ok {
		post.IsVisible = isVisible
	}

	post.UpdatedAt = time.Now()

	return s.repo.Update(post)
}

func (s *service) DeletePost(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *service) IncrementViewCount(slug string) error {
	post, err := s.repo.GetBySlug(slug)
	if err != nil {
		return err
	}
	return s.repo.IncrementViewCount(post.ID)
}
