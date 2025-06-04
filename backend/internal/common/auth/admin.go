package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminAuth struct {
	db        *gorm.DB
	jwtSecret string
}

type AdminClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.StandardClaims
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token     string `json:"token"`
	User      *AdminUser `json:"user"`
	ExpiresIn int64  `json:"expires_in"`
}

type AdminUser struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
	Username string `json:"username"`
}

type SetupRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	SetupToken      string `json:"setup_token" binding:"required"`
}

func NewAdminAuth(db *gorm.DB) *AdminAuth {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Fallback to a default secret (not recommended for production)
		jwtSecret = "default-secret-change-in-production"
	}
	return &AdminAuth{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (a *AdminAuth) ValidateAdminEmail(email string) bool {
	allowedDomains := []string{"jadenrazo.dev"}
	allowedEmails := []string{"support@jadenrazo.dev", "admin@jadenrazo.dev", "dev@jadenrazo.dev"}
	
	// Check exact matches first
	for _, allowed := range allowedEmails {
		if strings.ToLower(email) == strings.ToLower(allowed) {
			return true
		}
	}
	
	// Check domain matches
	for _, domain := range allowedDomains {
		if strings.HasSuffix(strings.ToLower(email), "@"+domain) {
			return true
		}
	}
	
	return false
}

func (a *AdminAuth) Login(req *AdminLoginRequest) (*AdminLoginResponse, error) {
	if !a.ValidateAdminEmail(req.Email) {
		return nil, errors.New("email not authorized for admin access")
	}
	
	var user entity.User
	err := a.db.Where("email = ? AND role = ?", req.Email, entity.RoleAdmin).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}
	
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}
	
	// Generate JWT token
	token, expiresIn, err := a.GenerateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}
	
	// Update last login
	user.RecordLogin()
	a.db.Save(&user)
	
	return &AdminLoginResponse{
		Token:     token,
		ExpiresIn: expiresIn,
		User: &AdminUser{
			ID:       user.ID.String(),
			Email:    user.Email,
			IsAdmin:  user.IsAdmin(),
			Username: user.Username,
		},
	}, nil
}

func (a *AdminAuth) GenerateToken(user *entity.User) (string, int64, error) {
	expiresIn := 24 * time.Hour
	claims := &AdminClaims{
		UserID:  user.ID.String(),
		Email:   user.Email,
		Role:    string(user.Role),
		IsAdmin: user.IsAdmin(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.ID.String(),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", 0, err
	}
	
	return tokenString, int64(expiresIn.Seconds()), nil
}

func (a *AdminAuth) ValidateToken(tokenString string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.jwtSecret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

func (a *AdminAuth) GenerateSetupToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (a *AdminAuth) SendSetupEmail(email, setupToken string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	
	if smtpHost == "" || smtpPort == "" {
		// Log the setup token for development
		fmt.Printf("\n=== ADMIN SETUP TOKEN ===\n")
		fmt.Printf("Email: %s\n", email)
		fmt.Printf("Setup Token: %s\n", setupToken)
		fmt.Printf("Use this token to complete admin setup at: /devpanel/setup\n")
		fmt.Printf("========================\n\n")
		return nil
	}
	
	subject := "Admin Account Setup - Project Website"
	body := fmt.Sprintf(`
You have been invited to set up an administrator account for the Project Website.

Setup Token: %s

Please visit the setup page and use this token to complete your account setup.

This token will expire in 24 hours.

If you did not request this, please ignore this email.
`, setupToken)
	
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)
	
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{email}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	
	return nil
}

func (a *AdminAuth) CompleteSetup(req *SetupRequest) error {
	if !a.ValidateAdminEmail(req.Email) {
		return errors.New("email not authorized for admin access")
	}
	
	if req.Password != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	
	// Verify setup token (in production, store tokens in database with expiration)
	// For now, we'll use a simple validation
	if len(req.SetupToken) < 20 {
		return errors.New("invalid setup token")
	}
	
	// Check if admin already exists
	var existingUser entity.User
	err := a.db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return errors.New("admin account already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("database error: %v", err)
	}
	
	// Create admin user
	username := strings.Split(req.Email, "@")[0]
	user, err := entity.NewUser(req.Email, username, req.Password, "Administrator")
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	
	user.Role = entity.RoleAdmin
	user.IsVerified = true
	
	err = a.db.Create(user).Error
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}
	
	return nil
}

func (a *AdminAuth) HasAdminAccount() (bool, error) {
	var count int64
	err := a.db.Model(&entity.User{}).Where("role = ?", entity.RoleAdmin).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
