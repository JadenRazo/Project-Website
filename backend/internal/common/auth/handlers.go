package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/common/auth/totp"
)

type AdminAuthHandlers struct {
	adminAuth *AdminAuth
}

func NewAdminAuthHandlers(db *gorm.DB) *AdminAuthHandlers {
	return &AdminAuthHandlers{
		adminAuth: NewAdminAuth(db),
	}
}

func (h *AdminAuthHandlers) Login(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := h.adminAuth.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AdminAuthHandlers) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		return
	}

	claims, err := h.adminAuth.ValidateToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       claims.UserID,
		"email":    claims.Email,
		"isAdmin":  claims.IsAdmin,
		"username": claims.Email, // Use email as username for admin
	})
}

func (h *AdminAuthHandlers) RequestSetup(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if !h.adminAuth.ValidateAdminEmail(req.Email) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Email not authorized for admin access"})
		return
	}

	// Check if admin already exists
	hasAdmin, err := h.adminAuth.HasAdminAccount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service temporarily unavailable"})
		return
	}
	if hasAdmin {
		c.JSON(http.StatusConflict, gin.H{"error": "Admin account already exists"})
		return
	}

	// Generate and send setup token
	setupToken, err := h.adminAuth.GenerateSetupToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate setup token"})
		return
	}

	err = h.adminAuth.SendSetupEmail(req.Email, setupToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send setup email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setup email sent successfully"})
}

func (h *AdminAuthHandlers) CompleteSetup(c *gin.Context) {
	var req SetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.adminAuth.CompleteSetup(&req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "Admin account already exists"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Setup failed"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin account created successfully"})
}

func (h *AdminAuthHandlers) CheckSetupStatus(c *gin.Context) {
	hasAdmin, err := h.adminAuth.HasAdminAccount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service temporarily unavailable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hasAdmin": hasAdmin})
}

func (h *AdminAuthHandlers) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		claims, err := h.adminAuth.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !claims.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func (h *AdminAuthHandlers) VerifyMFA(c *gin.Context) {
	var req MFAVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	encryptionKey := os.Getenv("JWT_SECRET")
	if encryptionKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	totpSvc := totp.NewTOTPService("Portfolio-DevPanel")
	encryptor := totp.NewEncryptor(encryptionKey)

	resp, err := h.adminAuth.VerifyMFAAndLogin(&req, totpSvc, encryptor)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
