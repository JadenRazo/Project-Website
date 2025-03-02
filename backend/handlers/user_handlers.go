package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	
	"urlshortener/db"
	"urlshortener/utils"
)

// RegisterRequest represents the request for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the request for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the response for user data
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// Register handles user registration
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	// Check if username is already taken
	var existingUser db.User
	if result := db.DB.Where("username = ?", req.Username).First(&existingUser); result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}
	
	// Check if email is already taken
	if result := db.DB.Where("email = ?", req.Email).First(&existingUser); result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	
	// Create user
	user := db.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    string(hashedPassword),
		Role:        "user", // Default role
		IsActive:    true,
		LastLoginAt: time.Now(),
	}
	
	if result := db.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	
	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"token": token,
		"user": UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	// Find user by email
	var user db.User
	if result := db.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	
	// Check if account is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Your account has been deactivated"})
		return
	}
	
	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	
	// Update last login time
	db.DB.Model(&user).Update("last_login_at", time.Now())
	
	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token": token,
		"user": UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	})
}

// GetUserProfile returns the current user's profile
func GetUserProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var user db.User
	if result := db.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}
	
	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	})
}

// UpdateUserProfile updates the current user's profile
func UpdateUserProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var user db.User
	if result := db.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}
	
	var req struct {
		Username string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
		Email    string `json:"email,omitempty" binding:"omitempty,email"`
		Password string `json:"password,omitempty" binding:"omitempty,min=8"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	// Check if updating username
	if req.Username != "" && req.Username != user.Username {
		// Check if username is already taken
		var existingUser db.User
		if result := db.DB.Where("username = ? AND id <> ?", req.Username, userID).First(&existingUser); result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			return
		}
		user.Username = req.Username
	}
	
	// Check if updating email
	if req.Email != "" && req.Email != user.Email {
		// Check if email is already taken
		var existingUser db.User
		if result := db.DB.Where("email = ? AND id <> ?", req.Email, userID).First(&existingUser); result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		user.Email = req.Email
	}
	
	// Check if updating password
	if req.Password != "" {
		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}
	
	// Save changes
	if result := db.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	})
}

// AddCustomDomain adds a new custom domain for the user
func AddCustomDomain(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var req struct {
		Domain string `json:"domain" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	// Check if domain is already registered
	var existingDomain db.CustomDomain
	if result := db.DB.Where("domain = ?", req.Domain).First(&existingDomain); result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Domain already registered"})
		return
	}
	
	// Create domain record
	domain := db.CustomDomain{
		Domain:     req.Domain,
		UserID:     userID.(uint),
		IsVerified: false, // Domains need to be verified
	}
	
	if result := db.DB.Create(&domain); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add custom domain"})
		return
	}
	
	// Generate verification token
	verificationToken := "verify_" + uuid.New().String()
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Domain added successfully. Please verify ownership.",
		"domain": domain.Domain,
		"verificationMethod": "DNS TXT Record",
		"verificationToken": verificationToken,
		"instructions": "Add a TXT record with name '_url-verify' and value '" + verificationToken + "' to verify ownership.",
	})
}

// GetUserDomains returns all custom domains for the current user
func GetUserDomains(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var domains []db.CustomDomain
	if result := db.DB.Where("user_id = ?", userID).Find(&domains); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve domains"})
		return
	}
	
	c.JSON(http.StatusOK, domains)
}
