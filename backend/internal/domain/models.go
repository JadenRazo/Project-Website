package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// BaseModel provides common fields for all models
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User represents a user account in the system
type User struct {
	BaseModel
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"`
	FirstName    string    `json:"firstName,omitempty"`
	LastName     string    `json:"lastName,omitempty"`
	Avatar       string    `json:"avatar,omitempty"`
	Bio          string    `json:"bio,omitempty"`
	Role         string    `gorm:"default:user" json:"role"`
	Status       string    `gorm:"default:offline" json:"status"`
	LastLogin    time.Time `json:"lastLogin,omitempty"`
	IsActive     bool      `gorm:"default:true" json:"isActive"`

	// Relations
	Channels  []Channel  `gorm:"many2many:user_channels" json:"-"`
	Messages  []Message  `gorm:"foreignKey:SenderID" json:"-"`
	Reactions []Reaction `gorm:"foreignKey:UserID" json:"-"`
	ShortURLs []ShortURL `gorm:"foreignKey:CreatorID" json:"-"`
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifies if the provided password matches the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return u.FirstName + " " + u.LastName
}

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
