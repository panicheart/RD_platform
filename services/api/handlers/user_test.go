package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"rdp-platform/rdp-api/middleware"
	"rdp-platform/rdp-api/models"
)

// MockUserService 模拟用户服务
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx interface{}, req models.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx interface{}, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByUsername(ctx interface{}, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserList(ctx interface{}, query models.UserListQuery) (*models.PaginatedResponse, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaginatedResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx interface{}, id string, req models.UpdateUserRequest, isAdmin bool) (*models.User, error) {
	args := m.Called(ctx, id, req, isAdmin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx interface{}, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) ValidateCredentials(ctx interface{}, username, password string) (*models.User, error) {
	args := m.Called(ctx, username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GenerateTokenPair(user *models.User) (*models.LoginResponse, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockUserService) RefreshAccessToken(ctx interface{}, refreshToken string) (*models.RefreshTokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshTokenResponse), args.Error(1)
}

func (m *MockUserService) Logout(ctx interface{}, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockUserService) ValidateToken(ctx interface{}, tokenString string) (*models.JWTClaims, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JWTClaims), args.Error(1)
}

func (m *MockUserService) UpdatePassword(ctx interface{}, userID, oldPassword, newPassword string) error {
	args := m.Called(ctx, userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserService) ResetPassword(ctx interface{}, userID, newPassword string) error {
	args := m.Called(ctx, userID, newPassword)
	return args.Error(0)
}

func (m *MockUserService) GetUserStats(ctx interface{}, userID string) (*models.UserStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserStats), args.Error(1)
}

func (m *MockUserService) GetUserContributions(ctx interface{}, userID string, year int) ([]models.UserContribution, error) {
	args := m.Called(ctx, userID, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserContribution), args.Error(1)
}

type UserHandlerTestSuite struct {
	suite.Suite
	mockService *MockUserService
	handler     *UserHandler
	router      *gin.Engine
}

func (s *UserHandlerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	s.mockService = new(MockUserService)
	s.handler = NewUserHandler(s.mockService)
	s.router = gin.New()
}

func (s *UserHandlerTestSuite) TearDownTest() {
	s.mockService.AssertExpectations(s.T())
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

// TestLogin_Success 测试成功登录
func (s *UserHandlerTestSuite) TestLogin_Success() {
	reqBody := map[string]string{
		"username": "testuser",
		"password": "TestPass123",
	}
	body, _ := json.Marshal(reqBody)

	mockUser := &models.User{
		ID:          "test-id",
		Username:    "testuser",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	mockResponse := &models.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    7200,
		TokenType:    "Bearer",
		User:         mockUser.ToResponse(),
	}

	s.mockService.On("ValidateCredentials", mock.Anything, "testuser", "TestPass123").Return(mockUser, nil)
	s.mockService.On("GenerateTokenPair", mockUser).Return(mockResponse, nil)

	s.router.POST("/login", s.handler.Login)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(s.T(), 200, response.Code)
}

// TestLogin_InvalidCredentials 测试无效凭据
func (s *UserHandlerTestSuite) TestLogin_InvalidCredentials() {
	reqBody := map[string]string{
		"username": "testuser",
		"password": "WrongPass",
	}
	body, _ := json.Marshal(reqBody)

	s.mockService.On("ValidateCredentials", mock.Anything, "testuser", "WrongPass").
		Return(nil, errors.New("invalid credentials"))

	s.router.POST("/login", s.handler.Login)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
}

// TestLogin_MissingFields 测试缺少字段
func (s *UserHandlerTestSuite) TestLogin_MissingFields() {
	reqBody := map[string]string{
		"username": "testuser",
		// missing password
	}
	body, _ := json.Marshal(reqBody)

	s.router.POST("/login", s.handler.Login)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// TestGetUser_Success 测试获取用户详情
func (s *UserHandlerTestSuite) TestGetUser_Success() {
	mockUser := &models.User{
		ID:          "test-id",
		Username:    "testuser",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	s.mockService.On("GetUserByID", mock.Anything, "test-id").Return(mockUser, nil)

	s.router.GET("/users/:id", s.handler.GetUser)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/users/test-id", nil)

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(s.T(), 200, response.Code)
}

// TestGetUser_NotFound 测试用户不存在
func (s *UserHandlerTestSuite) TestGetUser_NotFound() {
	s.mockService.On("GetUserByID", mock.Anything, "non-existent").
		Return(nil, errors.New("user not found"))

	s.router.GET("/users/:id", s.handler.GetUser)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/users/non-existent", nil)

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusNotFound, w.Code)
}

// TestCreateUser_Success 测试创建用户
func (s *UserHandlerTestSuite) TestCreateUser_Success() {
	reqBody := models.CreateUserRequest{
		Username:    "newuser",
		Password:    "TestPass123",
		DisplayName: "New User",
		Email:       "new@example.com",
		Role:        models.RoleDesigner,
	}
	body, _ := json.Marshal(reqBody)

	mockUser := &models.User{
		ID:          "new-id",
		Username:    "newuser",
		DisplayName: "New User",
		Email:       "new@example.com",
		Role:        models.RoleDesigner,
	}

	s.mockService.On("CreateUser", mock.Anything, reqBody).Return(mockUser, nil)

	s.router.POST("/users", s.handler.CreateUser)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(s.T(), 201, response.Code)
}

// TestGetUserList_Success 测试获取用户列表
func (s *UserHandlerTestSuite) TestGetUserList_Success() {
	mockResponse := &models.PaginatedResponse{
		List: []map[string]interface{}{
			{"id": "1", "username": "user1"},
			{"id": "2", "username": "user2"},
		},
		Total:    2,
		Page:     1,
		PageSize: 20,
		Pages:    1,
	}

	s.mockService.On("GetUserList", mock.Anything, mock.AnythingOfType("models.UserListQuery")).
		Return(mockResponse, nil)

	s.router.GET("/users", func(c *gin.Context) {
		c.Set("currentUser", &models.JWTClaims{
			UserID: "test-id",
			Role:   models.RoleAdmin,
		})
		s.handler.GetUserList(c)
	})

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/users?page=1&page_size=10", nil)

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(s.T(), 200, response.Code)
}

// TestUpdateUser_Success 测试更新用户
func (s *UserHandlerTestSuite) TestUpdateUser_Success() {
	reqBody := models.UpdateUserRequest{
		DisplayName: "Updated Name",
	}
	body, _ := json.Marshal(reqBody)

	mockUser := &models.User{
		ID:          "test-id",
		Username:    "testuser",
		DisplayName: "Updated Name",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	s.mockService.On("UpdateUser", mock.Anything, "test-id", reqBody, true).Return(mockUser, nil)

	s.router.PUT("/users/:id", func(c *gin.Context) {
		c.Set("currentUser", &models.JWTClaims{
			UserID: "test-id",
			Role:   models.RoleAdmin,
		})
		s.handler.UpdateUser(c)
	})

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/users/test-id", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestDeleteUser_Success 测试删除用户
func (s *UserHandlerTestSuite) TestDeleteUser_Success() {
	s.mockService.On("DeleteUser", mock.Anything, "test-id").Return(nil)

	s.router.DELETE("/users/:id", s.handler.DeleteUser)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/users/test-id", nil)

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestRefreshToken_Success 测试刷新令牌
func (s *UserHandlerTestSuite) TestRefreshToken_Success() {
	reqBody := map[string]string{
		"refresh_token": "valid-refresh-token",
	}
	body, _ := json.Marshal(reqBody)

	mockResponse := &models.RefreshTokenResponse{
		AccessToken: "new-access-token",
		ExpiresIn:   7200,
		TokenType:   "Bearer",
	}

	s.mockService.On("RefreshAccessToken", mock.Anything, "valid-refresh-token").
		Return(mockResponse, nil)

	s.router.POST("/refresh", s.handler.RefreshToken)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(s.T(), 200, response.Code)
}

// TestChangePassword_Success 测试修改密码
func (s *UserHandlerTestSuite) TestChangePassword_Success() {
	reqBody := map[string]string{
		"old_password": "OldPass123",
		"new_password": "NewPass456",
	}
	body, _ := json.Marshal(reqBody)

	s.mockService.On("UpdatePassword", mock.Anything, "test-id", "OldPass123", "NewPass456").
		Return(nil)

	s.router.PUT("/users/:id/password", func(c *gin.Context) {
		c.Set("currentUser", &models.JWTClaims{
			UserID: "test-id",
			Role:   models.RoleDesigner,
		})
		s.handler.ChangePassword(c)
	})

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/users/test-id/password", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestGetCurrentUser_Success 测试获取当前用户
func (s *UserHandlerTestSuite) TestGetCurrentUser_Success() {
	mockUser := &models.User{
		ID:          "test-id",
		Username:    "testuser",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	s.mockService.On("GetUserByID", mock.Anything, "test-id").Return(mockUser, nil)

	s.router.GET("/users/me", func(c *gin.Context) {
		c.Set("currentUser", &models.JWTClaims{
			UserID: "test-id",
			Role:   models.RoleDesigner,
		})
		s.handler.GetCurrentUser(c)
	})

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/users/me", nil)

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestGetUserProfile_Success 测试获取用户Profile
func (s *UserHandlerTestSuite) TestGetUserProfile_Success() {
	mockUser := &models.User{
		ID:          "test-id",
		Username:    "testuser",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	mockStats := &models.UserStats{
		TotalProjects:      10,
		CompletedProjects:  8,
		TotalContributions: 100,
	}

	s.mockService.On("GetUserByID", mock.Anything, "test-id").Return(mockUser, nil)
	s.mockService.On("GetUserStats", mock.Anything, "test-id").Return(mockStats, nil)
	s.mockService.On("GetUserContributions", mock.Anything, "test-id", 0).
		Return([]models.UserContribution{}, nil)

	s.router.GET("/users/:id/profile", s.handler.GetUserProfile)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/users/test-id/profile", nil)

	s.router.ServeHTTP(w, httpReq)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

// TestSuccessResponse 测试成功响应
func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SuccessResponse(c, map[string]string{"key": "value"})

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "success", response.Message)
}

// TestErrorResponse 测试错误响应
func TestErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ErrorResponse(c, http.StatusBadRequest, 40001, "Bad request")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 40001, response.Code)
	assert.Equal(t, "Bad request", response.Message)
}

// TestBadRequestResponse 测试400响应
func TestBadRequestResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	BadRequestResponse(c, "Invalid input")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUnauthorizedResponse 测试401响应
func TestUnauthorizedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	UnauthorizedResponse(c, "Unauthorized")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestForbiddenResponse 测试403响应
func TestForbiddenResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ForbiddenResponse(c, "Forbidden")

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestNotFoundResponse 测试404响应
func TestNotFoundResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	NotFoundResponse(c, "Not found")

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestInternalServerErrorResponse 测试500响应
func TestInternalServerErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	InternalServerErrorResponse(c, "Server error")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestUserToResponse 测试用户响应转换
func TestUserToResponse(t *testing.T) {
	user := &models.User{
		ID:          "test-id",
		Username:    "testuser",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
		Skills:      []string{"Go", "Python"},
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	response := user.ToResponse()
	assert.Equal(t, "test-id", response["id"])
	assert.Equal(t, "testuser", response["username"])
	assert.Equal(t, models.RoleDesigner, response["role"])
}

// TestUserIsAdmin 测试管理员检查
func TestUserIsAdmin(t *testing.T) {
	admin := &models.User{Role: models.RoleAdmin}
	assert.True(t, admin.IsAdmin())

	normal := &models.User{Role: models.RoleDesigner}
	assert.False(t, normal.IsAdmin())
}

// TestUserCheckPassword 测试密码检查
func TestUserCheckPassword(t *testing.T) {
	user := &models.User{}
	err := user.SetPassword("TestPass123")
	assert.NoError(t, err)

	assert.True(t, user.CheckPassword("TestPass123"))
	assert.False(t, user.CheckPassword("WrongPass"))
}

// TestMiddlewareGetCurrentUser 测试中间件获取当前用户
func TestMiddlewareGetCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 测试没有用户的情况
	_, ok := middleware.GetCurrentUser(c)
	assert.False(t, ok)

	// 测试有用户的情况
	claims := &models.JWTClaims{
		UserID:   "test-id",
		Username: "testuser",
		Role:     models.RoleAdmin,
	}
	c.Set("currentUser", claims)

	user, ok := middleware.GetCurrentUser(c)
	assert.True(t, ok)
	assert.Equal(t, "test-id", user.UserID)
}

// TestMiddlewareIsAdmin 测试中间件管理员检查
func TestMiddlewareIsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 没有用户
	assert.False(t, middleware.IsAdmin(c))

	// 管理员用户
	c.Set("currentUser", &models.JWTClaims{Role: models.RoleAdmin})
	assert.True(t, middleware.IsAdmin(c))

	// 普通用户
	c.Set("currentUser", &models.JWTClaims{Role: models.RoleDesigner})
	assert.False(t, middleware.IsAdmin(c))
}

// TestMiddlewareGetUserID 测试中间件获取用户ID
func TestMiddlewareGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 没有用户ID
	_, ok := middleware.GetUserID(c)
	assert.False(t, ok)

	// 有用户ID
	c.Set("userID", "test-id")
	id, ok := middleware.GetUserID(c)
	assert.True(t, ok)
	assert.Equal(t, "test-id", id)
}
