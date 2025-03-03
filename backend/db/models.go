package db

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// URL represents a shortened URL
type URL struct {
	Model
	ShortCode    string    `gorm:"uniqueIndex;not null" json:"shortCode"`
	OriginalURL  string    `gorm:"not null" json:"originalUrl"`
	ExpiresAt    *time.Time `json:"expiresAt,omitempty"`
	Clicks       uint       `gorm:"default:0" json:"clicks"`
	UserID       uint       `json:"userId,omitempty"`
	User         *User      `gorm:"foreignKey:UserID" json:"-"`
	IsPrivate    bool       `gorm:"default:false" json:"isPrivate"`
	CustomDomain string    `json:"customDomain,omitempty"`
	Tags         string    `json:"tags,omitempty"`
	Analytics    []ClickAnalytics `gorm:"foreignKey:URLID" json:"-"`
}

// User represents a registered user
type User struct {
	Model
	Username    string `gorm:"uniqueIndex;not null" json:"username"`
	Email       string `gorm:"uniqueIndex;not null" json:"email"`
	Password    string `gorm:"not null" json:"-"`
	Role        string `gorm:"default:user" json:"role"` // user or admin
	URLs        []URL  `gorm:"foreignKey:UserID" json:"-"`
	IsActive    bool   `gorm:"default:true" json:"isActive"`
	LastLoginAt time.Time `json:"lastLoginAt"`
	Domains     []CustomDomain `gorm:"foreignKey:UserID" json:"-"`
}

// CustomDomain stores users' custom domains for URL shortening
type CustomDomain struct {
	Model
	Domain      string `gorm:"uniqueIndex;not null" json:"domain"`
	UserID      uint   `json:"userId"`
	User        *User  `gorm:"foreignKey:UserID" json:"-"`
	IsVerified  bool   `gorm:"default:false" json:"isVerified"`
	VerifiedAt  *time.Time `json:"verifiedAt,omitempty"`
}

// ClickAnalytics tracks detailed information about URL clicks
type ClickAnalytics struct {
	Model
	URLID       uint   `json:"urlId"`
	URL         *URL   `gorm:"foreignKey:URLID" json:"-"`
	ClickedAt   time.Time `gorm:"index" json:"clickedAt"`
	IPAddress   string `json:"ipAddress,omitempty"`
	UserAgent   string `json:"userAgent,omitempty"`
	Referrer    string `json:"referrer,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	City        string `json:"city,omitempty"`
	Device      string `json:"device,omitempty"` // mobile, desktop, tablet, etc.
	Browser     string `json:"browser,omitempty"`
	OS          string `json:"os,omitempty"`
}

// HashPassword securely hashes a user's password
func (u *User) HashPassword() error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPass)
	return nil
}

// CheckPassword verifies the provided password against the stored hash
func (u *User) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	return err == nil
}

// BeforeSave handles password hashing before saving
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.ID == 0 || (u.Password != "" && !isHashedPassword(u.Password)) {
		return u.HashPassword()
	}
	return nil
}

// isHashedPassword checks if a password is already hashed
func isHashedPassword(password string) bool {
	// Bcrypt hashed passwords start with $2a$, $2b$, etc.
	return len(password) == 60 && password[0:4] == "$2a$" || password[0:4] == "$2b$"
}

// TableName specifies the table name for URL
func (URL) TableName() string {
	return "urls"
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// TableName specifies the table name for CustomDomain
func (CustomDomain) TableName() string {
	return "custom_domains"
}

// TableName specifies the table name for ClickAnalytics
func (ClickAnalytics) TableName() string {
	return "click_analytics"
}
