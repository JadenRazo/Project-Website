package blog

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Post struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Title           string         `json:"title" gorm:"type:varchar(255);not null"`
	Slug            string         `json:"slug" gorm:"type:varchar(255);uniqueIndex;not null"`
	Content         string         `json:"content" gorm:"type:text"`
	Excerpt         string         `json:"excerpt" gorm:"type:text"`
	FeaturedImage   string         `json:"featured_image" gorm:"type:varchar(500)"`
	AuthorID        *uuid.UUID     `json:"author_id" gorm:"type:uuid"`
	Status          string         `json:"status" gorm:"type:varchar(50);default:'draft'"`
	PublishedAt     *time.Time     `json:"published_at"`
	Tags            pq.StringArray `json:"tags" gorm:"type:text[]"`
	ViewCount       int            `json:"view_count" gorm:"default:0"`
	ReadTimeMinutes int            `json:"read_time_minutes" gorm:"default:1"`
	IsFeatured      bool           `json:"is_featured" gorm:"default:false"`
	IsVisible       bool           `json:"is_visible" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (Post) TableName() string {
	return "posts"
}

type ListParams struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Tag      string `form:"tag"`
	Search   string `form:"search"`
	Status   string `form:"status"`
}

type ListResult struct {
	Posts      []Post `json:"posts"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

type Repository interface {
	Create(post *Post) error
	GetByID(id uuid.UUID) (*Post, error)
	GetBySlug(slug string) (*Post, error)
	List(params ListParams) (*ListResult, error)
	ListPublished(params ListParams) (*ListResult, error)
	GetFeatured(limit int) ([]Post, error)
	Update(post *Post) error
	Delete(id uuid.UUID) error
	IncrementViewCount(id uuid.UUID) error
}

type Service interface {
	CreatePost(post *Post) error
	GetByID(id uuid.UUID) (*Post, error)
	GetBySlug(slug string) (*Post, error)
	List(params ListParams) (*ListResult, error)
	ListPublished(params ListParams) (*ListResult, error)
	GetFeatured(limit int) ([]Post, error)
	UpdatePost(id uuid.UUID, updates map[string]interface{}) error
	DeletePost(id uuid.UUID) error
	IncrementViewCount(slug string, r *http.Request) error
}
