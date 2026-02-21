package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AdminAuth struct {
	db        *gorm.DB
	jwtSecret string
}

type AdminClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	IsAdmin bool   `json:"is_admin"`
	jwt.StandardClaims
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token       string     `json:"token,omitempty"`
	User        *AdminUser `json:"user,omitempty"`
	ExpiresIn   int64      `json:"expires_in,omitempty"`
	RequiresMFA bool       `json:"requires_mfa"`
	MFAType     string     `json:"mfa_type,omitempty"`
	TempToken   string     `json:"temp_token,omitempty"`
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

type MFAVerificationRequest struct {
	TempToken    string `json:"temp_token" binding:"required"`
	MFACode      string `json:"mfa_code" binding:"required"`
	IsBackupCode bool   `json:"is_backup_code"`
}

func NewAdminAuth(db *gorm.DB) *AdminAuth {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is required but not set")
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

	// Check if TOTP MFA is enabled
	if user.TOTPEnabled {
		tempToken, err := a.GenerateTempToken(&user)
		if err != nil {
			return nil, fmt.Errorf("failed to generate temp token: %v", err)
		}

		return &AdminLoginResponse{
			RequiresMFA: true,
			MFAType:     "totp",
			TempToken:   tempToken,
		}, nil
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
		RequiresMFA: false,
		Token:       token,
		ExpiresIn:   expiresIn,
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

	expectedToken := os.Getenv("ADMIN_SETUP_TOKEN")
	if expectedToken == "" {
		return errors.New("admin setup is disabled")
	}

	if req.SetupToken != expectedToken {
		return errors.New("invalid setup token")
	}

	var existingUser entity.User
	err := a.db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return errors.New("admin account already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("database error: %v", err)
	}

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

func (a *AdminAuth) GenerateTempToken(user *entity.User) (string, error) {
	expiresIn := 5 * time.Minute
	claims := &AdminClaims{
		UserID:  user.ID.String(),
		Email:   user.Email,
		Role:    string(user.Role),
		IsAdmin: false,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   "mfa-pending",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AdminAuth) VerifyMFAAndLogin(req *MFAVerificationRequest, totpService interface {
	ValidateToken(secret, token string) bool
	ValidateBackupCode(hashedCodesJSON, inputCode string) (int, bool)
	MarkBackupCodeUsed(hashedCodesJSON string, index int) (string, error)
}, encryptor interface {
	Decrypt(encoded string) (string, error)
}) (*AdminLoginResponse, error) {
	claims, err := a.ValidateToken(req.TempToken)
	if err != nil {
		return nil, errors.New("invalid or expired temp token")
	}

	if claims.Subject != "mfa-pending" {
		return nil, errors.New("invalid temp token type")
	}

	var user entity.User
	err = a.db.Where("id = ?", claims.UserID).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.TOTPEnabled || user.TOTPSecret == nil {
		return nil, errors.New("MFA not enabled for this account")
	}

	decryptedSecret, err := encryptor.Decrypt(*user.TOTPSecret)
	if err != nil {
		return nil, errors.New("failed to decrypt TOTP secret")
	}

	var verified bool
	if req.IsBackupCode {
		if user.TOTPBackupCodes == nil {
			return nil, errors.New("no backup codes available")
		}

		index, valid := totpService.ValidateBackupCode(*user.TOTPBackupCodes, req.MFACode)
		if !valid {
			a.logMFAEvent(user.ID.String(), "recovery_failed", true)
			return nil, errors.New("invalid backup code")
		}

		updatedCodes, err := totpService.MarkBackupCodeUsed(*user.TOTPBackupCodes, index)
		if err != nil {
			return nil, errors.New("failed to mark backup code as used")
		}

		user.TOTPBackupCodes = &updatedCodes
		user.TOTPRecoveryUsed++
		a.db.Save(&user)

		a.logMFAEvent(user.ID.String(), "recovery_used", true)
		verified = true
	} else {
		verified = totpService.ValidateToken(decryptedSecret, req.MFACode)
		if !verified {
			a.logMFAEvent(user.ID.String(), "failed", false)
			return nil, errors.New("invalid TOTP code")
		}
		a.logMFAEvent(user.ID.String(), "verified", true)
	}

	if !verified {
		return nil, errors.New("MFA verification failed")
	}

	token, expiresIn, err := a.GenerateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	user.RecordLogin()
	a.db.Save(&user)

	return &AdminLoginResponse{
		RequiresMFA: false,
		Token:       token,
		ExpiresIn:   expiresIn,
		User: &AdminUser{
			ID:       user.ID.String(),
			Email:    user.Email,
			IsAdmin:  user.IsAdmin(),
			Username: user.Username,
		},
	}, nil
}

func (a *AdminAuth) logMFAEvent(userID, eventType string, success bool) {
	metadataJSON, _ := json.Marshal(map[string]interface{}{})

	event := map[string]interface{}{
		"user_id":    userID,
		"event_type": eventType,
		"success":    success,
		"metadata":   metadataJSON,
		"created_at": time.Now(),
	}

	a.db.Table("mfa_events").Create(event)
}
