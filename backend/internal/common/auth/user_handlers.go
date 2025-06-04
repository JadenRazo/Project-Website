package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoginRequest represents a login request payload
type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// RegisterRequest represents a registration request payload
type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Username        string `json:"username" binding:"required,min=3,max=30"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	FullName        string `json:"full_name" binding:"required,min=2,max=100"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse represents the response for authentication operations
type AuthResponse struct {
	User    interface{} `json:"user"`
	Tokens  *TokenPair  `json:"tokens"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// RegisterHandler handles user registration
func (a *Auth) RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Validate password confirmation
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "password_mismatch",
			Message: "Password and confirm password do not match",
		})
		return
	}

	// Validate password strength
	if err := ValidatePasswordStrength(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "weak_password",
			Message: err.Error(),
		})
		return
	}

	// Create user
	user, err := a.CreateUser(req.Email, req.Username, req.Password, req.FullName)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "user_exists",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "registration_failed",
			Message: "Failed to register user",
		})
		return
	}

	// Generate tokens
	tokens, err := a.GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication tokens",
		})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		User:    user.Sanitize(),
		Tokens:  tokens,
		Message: "Registration successful",
	})
}

// LoginHandler handles user login
func (a *Auth) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Authenticate user
	user, err := a.AuthenticateUser(req.EmailOrUsername, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "authentication_failed",
			Message: "Invalid credentials",
		})
		return
	}

	// Generate tokens
	tokens, err := a.GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication tokens",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		User:    user.Sanitize(),
		Tokens:  tokens,
		Message: "Login successful",
	})
}

// RefreshTokenHandler handles token refresh
func (a *Auth) RefreshTokenHandler(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Refresh tokens
	tokens, err := a.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "token_refresh_failed",
			Message: "Invalid or expired refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
		"message": "Token refreshed successfully",
	})
}

// LogoutHandler handles user logout
func (a *Auth) LogoutHandler(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Revoke refresh token
	err := a.RevokeToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "logout_failed",
			Message: "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// ValidateTokenHandler validates a JWT token
func (a *Auth) ValidateTokenHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "missing_token",
			Message: "Authorization header is required",
		})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "invalid_token_format",
			Message: "Invalid authorization header format",
		})
		return
	}

	claims, err := a.ValidateToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "invalid_token",
			Message: "Invalid or expired token",
		})
		return
	}

	// Get fresh user data
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID in token",
		})
		return
	}

	user, err := a.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "user_not_found",
			Message: "User associated with token not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user":  user.Sanitize(),
		"claims": claims,
	})
}

// ChangePasswordHandler handles password changes
func (a *Auth) ChangePasswordHandler(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Validate password confirmation
	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "password_mismatch",
			Message: "New password and confirm password do not match",
		})
		return
	}

	// Validate new password strength
	if err := ValidatePasswordStrength(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "weak_password",
			Message: err.Error(),
		})
		return
	}

	// Parse user ID
	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID",
		})
		return
	}

	// Update password
	err = a.UpdateUserPassword(uid, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if strings.Contains(err.Error(), "incorrect") {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "incorrect_password",
				Message: "Current password is incorrect",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "password_update_failed",
			Message: "Failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// GetProfileHandler returns the current user's profile
func (a *Auth) GetProfileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID",
		})
		return
	}

	user, err := a.GetUserByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "user_not_found",
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.Sanitize(),
	})
}

// RegisterUserRoutes registers user authentication routes
func (a *Auth) RegisterUserRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", a.RegisterHandler)
		auth.POST("/login", a.LoginHandler)
		auth.POST("/refresh", a.RefreshTokenHandler)
		auth.POST("/logout", a.LogoutHandler)
		auth.GET("/validate", a.ValidateTokenHandler)
		auth.GET("/profile", a.AuthMiddleware(), a.GetProfileHandler)
		auth.PUT("/password", a.AuthMiddleware(), a.ChangePasswordHandler)
	}
}

// AuthMiddleware returns a Gin middleware function that validates JWT tokens
func (a *Auth) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for certain paths
		if a.isPublicPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "missing_token",
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "invalid_token_format",
				Message: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := a.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "invalid_token",
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}