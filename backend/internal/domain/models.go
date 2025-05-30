package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
	"gorm.io/gorm"
)

// BaseModel provides common fields for all models
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User is an alias for entity.User to maintain backward compatibility
type User = entity.User

// Channel represents a messaging channel
type Channel struct {
	BaseModel
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	Type        string `gorm:"default:public" json:"type"` // public, private, direct
	OwnerID     uint   `json:"ownerId"`
	Icon        string `json:"icon,omitempty"`
	Banner      string `json:"banner,omitempty"`
	IsArchived  bool   `gorm:"default:false" json:"isArchived"`

	// Relations
	Owner    User      `gorm:"foreignKey:OwnerID" json:"-"`
	Members  []User    `gorm:"many2many:user_channels" json:"-"`
	Messages []Message `gorm:"foreignKey:ChannelID" json:"-"`
}

// Message represents a single message in a channel
type Message struct {
	BaseModel
	Content   string     `gorm:"not null" json:"content"`
	SenderID  uint       `json:"senderId"`
	ChannelID uint       `json:"channelId"`
	IsEdited  bool       `gorm:"default:false" json:"isEdited"`
	IsPinned  bool       `gorm:"default:false" json:"isPinned"`
	EditedAt  *time.Time `json:"editedAt,omitempty"`
	ReadBy    string     `gorm:"type:json" json:"-"` // JSON string of user IDs who have read the message

	// Relations
	Sender      User         `gorm:"foreignKey:SenderID" json:"-"`
	Channel     Channel      `gorm:"foreignKey:ChannelID" json:"-"`
	Attachments []Attachment `gorm:"foreignKey:MessageID" json:"attachments,omitempty"`
	Reactions   []Reaction   `gorm:"foreignKey:MessageID" json:"reactions,omitempty"`
	Mentions    []User       `gorm:"many2many:message_mentions" json:"-"`
}

// Attachment represents a file attached to a message
type Attachment struct {
	BaseModel
	FileName  string `gorm:"not null" json:"fileName"`
	FileType  string `gorm:"not null" json:"fileType"`
	FileSize  int64  `gorm:"not null" json:"fileSize"`
	FileURL   string `gorm:"not null" json:"fileUrl"`
	MessageID uint   `json:"messageId"`
	IsImage   bool   `gorm:"default:false" json:"isImage"`
	Thumbnail string `json:"thumbnail,omitempty"`

	// Relations
	Message Message `gorm:"foreignKey:MessageID" json:"-"`
}

// Reaction represents an emoji reaction to a message
type Reaction struct {
	BaseModel
	EmojiCode string `gorm:"not null" json:"emojiCode"`
	EmojiName string `json:"emojiName"`
	UserID    uint   `json:"userId"`
	MessageID uint   `json:"messageId"`

	// Relations
	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Message Message `gorm:"foreignKey:MessageID" json:"-"`
}

// ShortURL represents a shortened URL
type ShortURL struct {
	BaseModel
	OriginalURL    string    `gorm:"not null" json:"originalUrl"`
	ShortCode      string    `gorm:"uniqueIndex;not null" json:"shortCode"`
	CreatorID      uint      `json:"creatorId"`
	ExpiresAt      time.Time `json:"expiresAt,omitempty"`
	ClickCount     int64     `gorm:"default:0" json:"clickCount"`
	LastAccessedAt time.Time `json:"lastAccessedAt,omitempty"`
	Title          string    `json:"title,omitempty"`
	Description    string    `json:"description,omitempty"`
	IsCustomCode   bool      `gorm:"default:false" json:"isCustomCode"`

	// Relations
	Creator User       `gorm:"foreignKey:CreatorID" json:"-"`
	Clicks  []URLClick `gorm:"foreignKey:ShortURLID" json:"-"`
}

// URLClick represents a click on a shortened URL
type URLClick struct {
	BaseModel
	ShortURLID uint   `json:"shortUrlId"`
	IP         string `json:"-"` // Stored hashed or anonymized
	UserAgent  string `json:"userAgent,omitempty"`
	Referrer   string `json:"referrer,omitempty"`
	Country    string `json:"country,omitempty"`
	City       string `json:"city,omitempty"`
	DeviceType string `json:"deviceType,omitempty"`
}

// IsExpired checks if the shortened URL has expired
func (s *ShortURL) IsExpired() bool {
	return !s.ExpiresAt.IsZero() && time.Now().After(s.ExpiresAt)
}

// IncrementClicks increments the click counter and updates last accessed time
func (s *ShortURL) IncrementClicks() {
	s.ClickCount++
	s.LastAccessedAt = time.Now()
}

// GetFullURL returns the complete shortened URL using the base URL
func (s *ShortURL) GetFullURL(baseURL string) string {
	return baseURL + "/" + s.ShortCode
}

// CustomDomain represents a custom domain for shortened URLs
type CustomDomain struct {
	BaseModel
	Domain      string    `gorm:"uniqueIndex;not null" json:"domain"`
	UserID      uint      `json:"userId"`
	IsVerified  bool      `gorm:"default:false" json:"isVerified"`
	VerifiedAt  time.Time `json:"verifiedAt,omitempty"`
	VerifyToken string    `json:"-"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// CodeStats represents code statistics for the project
type CodeStats struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	Languages    LanguageStatsArr `gorm:"type:jsonb" json:"languages"`
	TotalLines   int64            `json:"total_lines"`
	TotalFiles   int              `json:"total_files"`
	TotalBlanks  int64            `json:"total_blanks"`
	TotalCode    int64            `json:"total_code"`
	TotalComment int64            `json:"total_comment"`
	UpdatedAt    time.Time        `json:"updated_at"`
	CreatedAt    time.Time        `json:"created_at"`
}

// LanguageStats represents statistics for a single language
type LanguageStats struct {
	Name     string `json:"name"`
	Files    int    `json:"files"`
	Lines    int64  `json:"lines"`
	Code     int64  `json:"code"`
	Comments int64  `json:"comments"`
	Blanks   int64  `json:"blanks"`
}

// LanguageStatsArr is a custom type for JSONB storage
type LanguageStatsArr []LanguageStats

// Value implements the driver.Valuer interface
func (l LanguageStatsArr) Value() (driver.Value, error) {
	return json.Marshal(l)
}

// Scan implements the sql.Scanner interface
func (l *LanguageStatsArr) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, l)
}
