package oauth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/common/cache"
	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
)

type OAuthHandlers struct {
	manager    *OAuthManager
	cache      cache.Cache
	db         *gorm.DB
	jwtManager *auth.JWTManager
	config     *auth.OAuth2Config
	encryptKey []byte
}

func NewOAuthHandlers(
	manager *OAuthManager,
	cache cache.Cache,
	db *gorm.DB,
	jwtManager *auth.JWTManager,
	config *auth.OAuth2Config,
) *OAuthHandlers {
	return &OAuthHandlers{
		manager:    manager,
		cache:      cache,
		db:         db,
		jwtManager: jwtManager,
		config:     config,
		encryptKey: []byte(config.EncryptionKey),
	}
}

// InitiateOAuth starts the OAuth flow by redirecting to the provider
func (h *OAuthHandlers) InitiateOAuth(c *gin.Context) {
	providerName := c.Param("provider")

	provider, err := h.manager.GetProvider(providerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth provider"})
		return
	}

	sessionID := generateSessionID()
	redirectURL := c.Query("redirect")
	if redirectURL == "" {
		redirectURL = "/devpanel"
	}

	if !isValidRedirectURL(redirectURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid redirect URL"})
		return
	}

	oauthState, err := GenerateOAuthState(
		c.Request.Context(),
		h.cache,
		providerName,
		sessionID,
		redirectURL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize OAuth flow"})
		return
	}

	authURLResp, err := provider.GetAuthURL(oauthState.State, oauthState.Nonce)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authorization URL"})
		return
	}

	if authURLResp.Verifier != "" {
		if err := oauthState.StoreVerifier(c.Request.Context(), h.cache, authURLResp.Verifier); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store PKCE verifier"})
			return
		}
	}

	c.Redirect(http.StatusTemporaryRedirect, authURLResp.URL)
}

// HandleCallback processes the OAuth provider callback
func (h *OAuthHandlers) HandleCallback(c *gin.Context) {
	providerName := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		errorDesc := c.Query("error_description")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             errorParam,
			"error_description": errorDesc,
		})
		return
	}

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state parameter"})
		return
	}

	oauthState, err := ValidateOAuthState(c.Request.Context(), h.cache, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired state parameter"})
		return
	}

	if oauthState.Provider != providerName {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider mismatch"})
		return
	}

	provider, err := h.manager.GetProvider(providerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	tokenResp, err := provider.ExchangeCode(c.Request.Context(), code, oauthState.Verifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange authorization code"})
		return
	}

	userInfo, err := provider.GetUserInfo(c.Request.Context(), tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
		return
	}

	if !h.config.IsEmailAllowed(userInfo.Email) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Email address not authorized for admin access"})
		return
	}

	if !userInfo.EmailVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Email address not verified with OAuth provider"})
		return
	}

	user, err := h.findOrCreateUser(c.Request.Context(), userInfo, providerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update user"})
		return
	}

	if err := h.storeOAuthTokens(c.Request.Context(), user.ID.String(), tokenResp, providerName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OAuth tokens"})
		return
	}

	tokenPair, err := h.jwtManager.GenerateTokenPair(
		user.ID.String(),
		user.Username,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication tokens"})
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	redirectURL := fmt.Sprintf(
		"%s%s#access_token=%s&token_type=Bearer&expires_in=%d",
		frontendURL,
		oauthState.RedirectURL,
		tokenPair.AccessToken,
		tokenPair.ExpiresIn,
	)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// findOrCreateUser finds an existing user or creates a new one
func (h *OAuthHandlers) findOrCreateUser(ctx context.Context, userInfo *UserInfo, provider string) (*entity.User, error) {
	var user entity.User

	err := h.db.Where("email = ?", userInfo.Email).First(&user).Error
	if err == nil {
		user.FullName = userInfo.Name
		user.AvatarURL = userInfo.Picture
		user.LastLogin = timePtr(time.Now())
		user.TwoFactorEnabled = true
		user.TwoFactorProvider = &provider
		user.TwoFactorProviderID = &userInfo.ProviderUserID

		if err := h.db.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
		return &user, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("database error: %w", err)
	}

	username := generateUsernameFromEmail(userInfo.Email)
	user = entity.User{
		ID:                  uuid.New(),
		Username:            username,
		Email:               userInfo.Email,
		HashedPassword:      generateRandomHash(),
		FullName:            userInfo.Name,
		AvatarURL:           userInfo.Picture,
		IsActive:            true,
		IsVerified:          true,
		Role:                entity.RoleAdmin,
		LastLogin:           timePtr(time.Now()),
		TwoFactorEnabled:    true,
		TwoFactorProvider:   &provider,
		TwoFactorProviderID: &userInfo.ProviderUserID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := h.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// storeOAuthTokens encrypts and stores OAuth tokens in the database
func (h *OAuthHandlers) storeOAuthTokens(ctx context.Context, userID string, tokenResp *TokenResponse, provider string) error {
	encryptedAccess, err := EncryptToken(tokenResp.AccessToken, h.encryptKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}

	var encryptedRefresh string
	if tokenResp.RefreshToken != "" {
		encryptedRefresh, err = EncryptToken(tokenResp.RefreshToken, h.encryptKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt refresh token: %w", err)
		}
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	tokenExpiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	oauthToken := &OAuthToken{
		ID:                   uuid.New(),
		UserID:               userUUID,
		Provider:             provider,
		EncryptedAccessToken: encryptedAccess,
		EncryptedRefreshToken: encryptedRefresh,
		TokenExpiry:          &tokenExpiry,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	result := h.db.Exec(`
		INSERT INTO oauth_tokens (id, user_id, provider, encrypted_access_token, encrypted_refresh_token, token_expiry, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (user_id, provider)
		DO UPDATE SET
			encrypted_access_token = EXCLUDED.encrypted_access_token,
			encrypted_refresh_token = EXCLUDED.encrypted_refresh_token,
			token_expiry = EXCLUDED.token_expiry,
			updated_at = EXCLUDED.updated_at
	`, oauthToken.ID, oauthToken.UserID, oauthToken.Provider, oauthToken.EncryptedAccessToken,
		oauthToken.EncryptedRefreshToken, oauthToken.TokenExpiry, oauthToken.CreatedAt, oauthToken.UpdatedAt)

	if result.Error != nil {
		return fmt.Errorf("failed to store OAuth tokens: %w", result.Error)
	}

	return nil
}

// setAuthCookies sets httpOnly cookies for authentication
func (h *OAuthHandlers) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	secure := c.Request.Header.Get("X-Forwarded-Proto") == "https" || c.Request.TLS != nil

	c.SetCookie(
		"access_token",
		accessToken,
		900,
		"/",
		"",
		secure,
		true,
	)
	c.SetSameSite(http.SameSiteStrictMode)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		604800,
		"/api/v1/auth/refresh",
		"",
		secure,
		true,
	)
	c.SetSameSite(http.SameSiteStrictMode)
}

// OAuthToken represents the database model for stored OAuth tokens
type OAuthToken struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primary_key"`
	UserID                uuid.UUID  `gorm:"type:uuid;not null"`
	Provider              string     `gorm:"type:varchar(50);not null"`
	EncryptedAccessToken  string     `gorm:"type:text;not null"`
	EncryptedRefreshToken string     `gorm:"type:text"`
	TokenExpiry           *time.Time `gorm:"type:timestamp"`
	CreatedAt             time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt             time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (OAuthToken) TableName() string {
	return "oauth_tokens"
}

// Helper functions

func generateSessionID() string {
	return uuid.New().String()
}

func generateUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		username := strings.ReplaceAll(parts[0], ".", "_")
		username = strings.ReplaceAll(username, "+", "_")
		return username
	}
	return fmt.Sprintf("user_%s", uuid.New().String()[:8])
}

func generateRandomHash() string {
	return fmt.Sprintf("oauth_user_%s", uuid.New().String())
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func isValidRedirectURL(redirectURL string) bool {
	if strings.HasPrefix(redirectURL, "/") && !strings.HasPrefix(redirectURL, "//") {
		return true
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL != "" {
		if strings.HasPrefix(redirectURL, frontendURL+"/") || redirectURL == frontendURL {
			return true
		}
	}

	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"https://jadenrazo.dev",
		"https://www.jadenrazo.dev",
	}

	for _, origin := range allowedOrigins {
		if strings.HasPrefix(redirectURL, origin+"/") || redirectURL == origin {
			return true
		}
	}

	return false
}
