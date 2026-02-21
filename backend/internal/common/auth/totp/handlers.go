package totp

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
)

type TOTPHandlers struct {
	db             *gorm.DB
	totpService    *TOTPService
	encryptor      *Encryptor
	issuer         string
}

type SetupRequest struct {
	Password string `json:"password" binding:"required"`
}

type VerifyRequest struct {
	Token string `json:"token" binding:"required,len=6"`
}

type DisableRequest struct {
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required,len=6"`
}

type MFAStatusResponse struct {
	TOTPEnabled         bool   `json:"totp_enabled"`
	TOTPVerified        bool   `json:"totp_verified"`
	BackupCodesRemaining int    `json:"backup_codes_remaining"`
	EnabledAt           *time.Time `json:"enabled_at,omitempty"`
}

func NewTOTPHandlers(db *gorm.DB) *TOTPHandlers {
	issuer := os.Getenv("TOTP_ISSUER")
	if issuer == "" {
		issuer = "Portfolio-DevPanel"
	}

	encryptionKey := os.Getenv("JWT_SECRET")
	if encryptionKey == "" {
		encryptionKey = "default-encryption-key-please-change"
	}

	return &TOTPHandlers{
		db:          db,
		totpService: NewTOTPService(issuer),
		encryptor:   NewEncryptor(encryptionKey),
		issuer:      issuer,
	}
}

func (h *TOTPHandlers) SetupTOTP(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req SetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user entity.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	if user.TOTPEnabled {
		c.JSON(http.StatusConflict, gin.H{"error": "TOTP is already enabled"})
		return
	}

	setup, err := h.totpService.GenerateSecret(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
		return
	}

	encryptedSecret, err := h.encryptor.Encrypt(setup.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt secret"})
		return
	}

	plainCodes, hashedCodes, err := h.totpService.GenerateBackupCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate backup codes"})
		return
	}

	hashedCodesJSON, err := json.Marshal(hashedCodes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process backup codes"})
		return
	}

	secret := encryptedSecret
	backupCodes := string(hashedCodesJSON)
	user.TOTPSecret = &secret
	user.TOTPBackupCodes = &backupCodes
	user.TOTPEnabled = false
	user.TOTPVerified = false

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save TOTP configuration"})
		return
	}

	h.logMFAEvent(user.ID, "setup_initiated", c.ClientIP(), c.Request.UserAgent(), true, nil)

	c.JSON(http.StatusOK, gin.H{
		"secret":       setup.Secret,
		"qr_code":      setup.QRCode,
		"url":          setup.URL,
		"backup_codes": plainCodes,
	})
}

func (h *TOTPHandlers) VerifyTOTP(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user entity.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.TOTPSecret == nil || *user.TOTPSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TOTP not set up"})
		return
	}

	if user.TOTPEnabled {
		c.JSON(http.StatusConflict, gin.H{"error": "TOTP already verified and enabled"})
		return
	}

	decryptedSecret, err := h.encryptor.Decrypt(*user.TOTPSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt TOTP secret"})
		return
	}

	if !h.totpService.ValidateToken(decryptedSecret, req.Token) {
		h.logMFAEvent(user.ID, "verification_failed", c.ClientIP(), c.Request.UserAgent(), false, nil)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP token"})
		return
	}

	now := time.Now()
	user.TOTPEnabled = true
	user.TOTPVerified = true
	user.TOTPEnabledAt = &now

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable TOTP"})
		return
	}

	h.logMFAEvent(user.ID, "enabled", c.ClientIP(), c.Request.UserAgent(), true, nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "TOTP successfully enabled",
		"enabled": true,
	})
}

func (h *TOTPHandlers) DisableTOTP(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req DisableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user entity.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	if !user.TOTPEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TOTP is not enabled"})
		return
	}

	if user.TOTPSecret != nil {
		decryptedSecret, err := h.encryptor.Decrypt(*user.TOTPSecret)
		if err == nil {
			if !h.totpService.ValidateToken(decryptedSecret, req.Token) {
				h.logMFAEvent(user.ID, "disable_failed", c.ClientIP(), c.Request.UserAgent(), false, nil)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP token"})
				return
			}
		}
	}

	user.TOTPSecret = nil
	user.TOTPEnabled = false
	user.TOTPVerified = false
	user.TOTPBackupCodes = nil
	user.TOTPRecoveryUsed = 0
	user.TOTPEnabledAt = nil

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable TOTP"})
		return
	}

	h.logMFAEvent(user.ID, "disabled", c.ClientIP(), c.Request.UserAgent(), true, nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "TOTP successfully disabled",
		"enabled": false,
	})
}

func (h *TOTPHandlers) GetMFAStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user entity.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	remaining := 0
	if user.TOTPBackupCodes != nil {
		remaining = h.totpService.CountRemainingBackupCodes(*user.TOTPBackupCodes)
	}

	c.JSON(http.StatusOK, MFAStatusResponse{
		TOTPEnabled:          user.TOTPEnabled,
		TOTPVerified:         user.TOTPVerified,
		BackupCodesRemaining: remaining,
		EnabledAt:            user.TOTPEnabledAt,
	})
}

func (h *TOTPHandlers) RegenerateBackupCodes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user entity.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	if !user.TOTPEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TOTP is not enabled"})
		return
	}

	plainCodes, hashedCodes, err := h.totpService.GenerateBackupCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate backup codes"})
		return
	}

	hashedCodesJSON, err := json.Marshal(hashedCodes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process backup codes"})
		return
	}

	backupCodes := string(hashedCodesJSON)
	user.TOTPBackupCodes = &backupCodes
	user.TOTPRecoveryUsed = 0

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save backup codes"})
		return
	}

	h.logMFAEvent(user.ID, "backup_regenerated", c.ClientIP(), c.Request.UserAgent(), true, nil)

	c.JSON(http.StatusOK, gin.H{
		"backup_codes": plainCodes,
		"message":      "Backup codes regenerated successfully",
	})
}

func (h *TOTPHandlers) logMFAEvent(userID uuid.UUID, eventType, ipAddress, userAgent string, success bool, metadata map[string]interface{}) {
	metadataJSON, _ := json.Marshal(metadata)

	event := map[string]interface{}{
		"user_id":    userID,
		"event_type": eventType,
		"ip_address": ipAddress,
		"user_agent": userAgent,
		"success":    success,
		"metadata":   metadataJSON,
		"created_at": time.Now(),
	}

	h.db.Table("mfa_events").Create(event)
}
