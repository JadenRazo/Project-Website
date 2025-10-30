package mock

import (
	"context"
	"fmt"
	"sync"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
)

type Claims struct {
	UserID         string `json:"user_id"`
	Role           string `json:"role"`
	StandardClaims StandardClaims
}

type StandardClaims struct {
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	Issuer    string `json:"iss"`
	Audience  string `json:"aud"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type SessionInfo struct {
	UserID    string    `json:"user_id"`
	LoginTime time.Time `json:"login_time"`
	LastSeen  time.Time `json:"last_seen"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
}

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func HashPassword(password string) (string, error) {
	if len(password) > 72 {
		return "", fmt.Errorf("password too long")
	}

	if len(password) == 0 {
		return "", fmt.Errorf("password cannot be empty")
	}

	for _, b := range []byte(password) {
		if b == 0 {
			return "", fmt.Errorf("password contains null byte")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPassword(password, hash string) bool {
	if hash == "" || password == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type MockUserRepository struct {
	users     map[string]*entity.User
	usersById map[uuid.UUID]*entity.User
	mu        sync.RWMutex
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:     make(map[string]*entity.User),
		usersById: make(map[uuid.UUID]*entity.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.Email]; exists {
		return fmt.Errorf("user already exists")
	}

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	m.users[user.Email] = user
	m.usersById[user.ID] = user
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.usersById[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.usersById[user.ID]; !exists {
		return fmt.Errorf("user not found")
	}

	user.UpdatedAt = time.Now()
	m.users[user.Email] = user
	m.usersById[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, exists := m.usersById[id]
	if !exists {
		return fmt.Errorf("user not found")
	}

	delete(m.users, user.Email)
	delete(m.usersById, id)
	return nil
}

type MockJWTService struct {
	secret            string
	validTokens       map[string]*Claims
	blacklistedTokens map[string]bool
	mu                sync.RWMutex
}

func NewMockJWTService(secret string) *MockJWTService {
	return &MockJWTService{
		secret:            secret,
		validTokens:       make(map[string]*Claims),
		blacklistedTokens: make(map[string]bool),
	}
}

func (m *MockJWTService) GenerateToken(claims *Claims) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	tokenID := fmt.Sprintf("mock-token-%d", time.Now().UnixNano())
	m.validTokens[tokenID] = claims
	return tokenID, nil
}

func (m *MockJWTService) ValidateToken(tokenString string) (*Claims, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.blacklistedTokens[tokenString] {
		return nil, fmt.Errorf("token is blacklisted")
	}

	claims, exists := m.validTokens[tokenString]
	if !exists {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.StandardClaims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

func (m *MockJWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	newClaims := &Claims{
		UserID: claims.UserID,
		Role:   claims.Role,
		StandardClaims: StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    claims.StandardClaims.Issuer,
			Audience:  claims.StandardClaims.Audience,
		},
	}

	return m.GenerateToken(newClaims)
}

func (m *MockJWTService) BlacklistToken(tokenString string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.blacklistedTokens[tokenString] = true
	return nil
}

type MockAuthService struct {
	userRepo   UserRepository
	jwtService *MockJWTService
	sessions   map[string]*SessionInfo
	mu         sync.RWMutex
}

func NewMockAuthService(userRepo UserRepository, jwtService *MockJWTService) *MockAuthService {
	return &MockAuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
		sessions:   make(map[string]*SessionInfo),
	}
}

func (m *MockAuthService) Register(ctx context.Context, email, password, username string) (*entity.User, error) {
	if !IsValidPassword(password) {
		return nil, fmt.Errorf("invalid password")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:          email,
		Username:       username,
		HashedPassword: hashedPassword,
		IsActive:       true,
		IsVerified:     false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = m.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := m.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	if !CheckPassword(password, user.HashedPassword) {
		return nil, fmt.Errorf("invalid credentials")
	}

	claims := &Claims{
		UserID: user.ID.String(),
		Role:   "user",
		StandardClaims: StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "devpanel-test",
			Audience:  "devpanel-users",
		},
	}

	accessToken, err := m.jwtService.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	refreshClaims := &Claims{
		UserID: user.ID.String(),
		Role:   "refresh",
		StandardClaims: StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "devpanel-test",
			Audience:  "devpanel-users",
		},
	}

	refreshToken, err := m.jwtService.GenerateToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	sessionInfo := &SessionInfo{
		UserID:    user.ID.String(),
		LoginTime: time.Now(),
		LastSeen:  time.Now(),
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	}

	m.mu.Lock()
	m.sessions[accessToken] = sessionInfo
	m.mu.Unlock()

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    24 * 60 * 60,
	}, nil
}

func (m *MockAuthService) ValidateToken(tokenString string) (*Claims, error) {
	return m.jwtService.ValidateToken(tokenString)
}

func (m *MockAuthService) RefreshToken(refreshToken string) (*TokenPair, error) {
	newAccessToken, err := m.jwtService.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    24 * 60 * 60,
	}, nil
}

func (m *MockAuthService) Logout(tokenString string) error {
	return m.jwtService.BlacklistToken(tokenString)
}

func (m *MockAuthService) GetSessionInfo(tokenString string) (*SessionInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[tokenString]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	return session, nil
}

func (m *MockAuthService) GenerateTokens(userID, role string) (*TokenPair, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "devpanel-test",
			Audience:  "devpanel-users",
		},
	}

	accessToken, err := m.jwtService.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	refreshClaims := &Claims{
		UserID: userID,
		Role:   "refresh",
		StandardClaims: StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "devpanel-test",
			Audience:  "devpanel-users",
		},
	}

	refreshToken, err := m.jwtService.GenerateToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    24 * 60 * 60,
	}, nil
}

type MockRateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
	mu       sync.RWMutex
}

func NewMockRateLimiter(limit int, window time.Duration) *MockRateLimiter {
	return &MockRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (m *MockRateLimiter) Allow(identifier string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-m.window)

	requests := m.requests[identifier]

	filtered := make([]time.Time, 0)
	for _, req := range requests {
		if req.After(cutoff) {
			filtered = append(filtered, req)
		}
	}

	if len(filtered) >= m.limit {
		return false
	}

	filtered = append(filtered, now)
	m.requests[identifier] = filtered

	return true
}

func (m *MockRateLimiter) GetRequestCount(identifier string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.requests[identifier])
}

func (m *MockRateLimiter) Reset(identifier string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.requests, identifier)
}
