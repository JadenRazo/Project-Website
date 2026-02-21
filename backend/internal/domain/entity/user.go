package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Role represents a user role
type Role string

const (
	// RoleUser is the standard user role
	RoleUser Role = "user"

	// RoleAdmin is the administrator role
	RoleAdmin Role = "admin"

	// RoleModerator is the moderator role
	RoleModerator Role = "moderator"
)

// User represents a user in the system
type User struct {
	ID                  uuid.UUID  `json:"id" db:"id" gorm:"column:id;type:uuid;primaryKey"`
	Email               string     `json:"email" db:"email" gorm:"column:email;type:citext;uniqueIndex"`
	Username            string     `json:"username" db:"username" gorm:"column:username;type:citext;uniqueIndex"`
	HashedPassword      string     `json:"-" db:"password_hash" gorm:"column:password_hash"`
	FullName            string     `json:"full_name" db:"full_name" gorm:"column:full_name"`
	AvatarURL           string     `json:"avatar_url,omitempty" db:"avatar_url" gorm:"column:avatar_url"`
	Role                Role       `json:"role" db:"role" gorm:"column:role"`
	IsActive            bool       `json:"is_active" db:"is_active" gorm:"column:is_active"`
	IsVerified          bool       `json:"is_verified" db:"is_verified" gorm:"column:is_verified"`
	TwoFactorEnabled    bool       `json:"two_factor_enabled" db:"two_factor_enabled" gorm:"column:two_factor_enabled"`
	TwoFactorProvider   *string    `json:"two_factor_provider,omitempty" db:"two_factor_provider" gorm:"column:two_factor_provider"`
	TwoFactorProviderID *string    `json:"two_factor_provider_id,omitempty" db:"two_factor_provider_id" gorm:"column:two_factor_provider_id"`
	TOTPSecret          *string    `json:"-" db:"totp_secret" gorm:"column:totp_secret"`
	TOTPEnabled         bool       `json:"totp_enabled" db:"totp_enabled" gorm:"column:totp_enabled"`
	TOTPVerified        bool       `json:"totp_verified" db:"totp_verified" gorm:"column:totp_verified"`
	TOTPBackupCodes     *string    `json:"-" db:"totp_backup_codes" gorm:"column:totp_backup_codes"`
	TOTPRecoveryUsed    int        `json:"totp_recovery_used" db:"totp_recovery_used" gorm:"column:totp_recovery_used"`
	TOTPEnabledAt       *time.Time `json:"totp_enabled_at,omitempty" db:"totp_enabled_at" gorm:"column:totp_enabled_at"`
	LastLogin           *time.Time `json:"last_login,omitempty" db:"last_login" gorm:"column:last_login_at"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at" gorm:"column:created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at" gorm:"column:updated_at"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// NewUser creates a new user with the given information
func NewUser(email, username, password, fullName string) (*User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &User{
		ID:             uuid.New(),
		Email:          email,
		Username:       username,
		HashedPassword: string(hashedPassword),
		FullName:       fullName,
		Role:           RoleUser,
		IsActive:       true,
		IsVerified:     false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// CheckPassword validates if the provided password matches the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}

// UpdatePassword updates the user's password
func (u *User) UpdatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.HashedPassword = string(hashedPassword)
	u.UpdatedAt = time.Now()

	return nil
}

// IsAdmin checks if the user has administrator privileges
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsModerator checks if the user has moderator privileges
func (u *User) IsModerator() bool {
	return u.Role == RoleModerator || u.Role == RoleAdmin
}

// HasPermission checks if the user has the specified role or higher
func (u *User) HasPermission(requiredRole Role) bool {
	switch requiredRole {
	case RoleUser:
		return true
	case RoleModerator:
		return u.Role == RoleModerator || u.Role == RoleAdmin
	case RoleAdmin:
		return u.Role == RoleAdmin
	default:
		return false
	}
}

// RecordLogin updates the user's last login time
func (u *User) RecordLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

// Sanitize removes sensitive information from the user object
func (u *User) Sanitize() *User {
	// Create a copy to avoid modifying the original
	sanitized := *u
	sanitized.HashedPassword = ""
	return &sanitized
}
