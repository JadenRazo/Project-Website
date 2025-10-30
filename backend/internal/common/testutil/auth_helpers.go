package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/common/security"
	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
)

type AuthTestSuite struct {
	DB          *gorm.DB
	AuthService *auth.Service
	AdminAuth   *auth.AdminAuth
	SecurityMgr *security.Manager
	Router      *gin.Engine
	TestUser    *entity.User
	TestAdmin   *entity.User
	JWTSecret   string
}

func SetupAuthTestSuite(t *testing.T) *AuthTestSuite {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Silent,
	})
	require.NoError(t, err)

	err = db.AutoMigrate(&entity.User{})
	require.NoError(t, err)

	jwtSecret := "test-secret-key-for-testing-only"

	authService := auth.NewService(db, jwtSecret)
	adminAuth := auth.NewAdminAuth(authService)

	securityConfig := &security.Config{
		Environment:       "test",
		RateLimit:         1000,
		EnableCORS:        true,
		AllowedOrigins:    []string{"http://localhost:3000"},
		EnableIPWhitelist: false,
		AuditLogLevel:     "info",
	}
	securityMgr := security.NewManager(securityConfig)

	router := gin.New()
	router.Use(gin.Recovery())

	testUser := &entity.User{
		Email:      "test@example.com",
		Password:   "TestPassword123!",
		Role:       entity.UserRole,
		IsActive:   true,
		IsVerified: true,
	}

	testAdmin := &entity.User{
		Email:      "admin@jadenrazo.dev",
		Password:   "AdminPassword123!",
		Role:       entity.AdminRole,
		IsActive:   true,
		IsVerified: true,
	}

	return &AuthTestSuite{
		DB:          db,
		AuthService: authService,
		AdminAuth:   adminAuth,
		SecurityMgr: securityMgr,
		Router:      router,
		TestUser:    testUser,
		TestAdmin:   testAdmin,
		JWTSecret:   jwtSecret,
	}
}

func (suite *AuthTestSuite) CreateTestUser(t *testing.T, user *entity.User) *entity.User {
	hashedPassword, err := auth.HashPassword(user.Password)
	require.NoError(t, err)

	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = suite.DB.Create(user).Error
	require.NoError(t, err)

	return user
}

func (suite *AuthTestSuite) GenerateValidJWT(t *testing.T, userID string, role entity.Role) string {
	token, err := suite.AuthService.GenerateTokens(userID, string(role))
	require.NoError(t, err)
	return token.AccessToken
}

func (suite *AuthTestSuite) GenerateExpiredJWT(t *testing.T, userID string, role entity.Role) string {
	jwtService := auth.NewJWTService(suite.JWTSecret)

	claims := &auth.Claims{
		UserID: userID,
		Role:   string(role),
		StandardClaims: auth.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
			IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
			Issuer:    "devpanel-test",
			Audience:  "devpanel-users",
		},
	}

	token, err := jwtService.GenerateToken(claims)
	require.NoError(t, err)
	return token
}

func (suite *AuthTestSuite) MakeAuthenticatedRequest(t *testing.T, method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, url, reqBody)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	return w
}

func (suite *AuthTestSuite) AssertJSONError(t *testing.T, resp *httptest.ResponseRecorder, expectedCode int, expectedMessage string) {
	require.Equal(t, expectedCode, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Contains(t, response, "error")
	require.Equal(t, expectedMessage, response["error"])
}

func (suite *AuthTestSuite) AssertSuccessResponse(t *testing.T, resp *httptest.ResponseRecorder, expectedCode int) {
	require.Equal(t, expectedCode, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err)

	require.NotContains(t, response, "error")
}

func (suite *AuthTestSuite) Cleanup(t *testing.T) {
	db, err := suite.DB.DB()
	require.NoError(t, err)
	err = db.Close()
	require.NoError(t, err)
}

type MockUserRepository struct {
	users map[string]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (m *MockUserRepository) Create(user *entity.User) error {
	if _, exists := m.users[user.Email]; exists {
		return fmt.Errorf("user already exists")
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByEmail(email string) (*entity.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) GetByID(id string) (*entity.User, error) {
	for _, user := range m.users {
		if user.ID.String() == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) Update(user *entity.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(id string) error {
	for email, user := range m.users {
		if user.ID.String() == id {
			delete(m.users, email)
			return nil
		}
	}
	return fmt.Errorf("user not found")
}

type SecurityTestHelper struct {
	requestCounts map[string]int
}

func NewSecurityTestHelper() *SecurityTestHelper {
	return &SecurityTestHelper{
		requestCounts: make(map[string]int),
	}
}

func (h *SecurityTestHelper) SimulateRateLimitRequests(suite *AuthTestSuite, endpoint string, count int) []*httptest.ResponseRecorder {
	responses := make([]*httptest.ResponseRecorder, count)

	for i := 0; i < count; i++ {
		resp := suite.MakeAuthenticatedRequest(nil, "GET", endpoint, nil, "")
		responses[i] = resp
		h.requestCounts[endpoint]++
	}

	return responses
}

func (h *SecurityTestHelper) GetRequestCount(endpoint string) int {
	return h.requestCounts[endpoint]
}
