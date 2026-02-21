package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"rdp-platform/rdp-api/models"
)

type UserServiceTestSuite struct {
	suite.Suite
	db          *gorm.DB
	userService *UserService
	ctx         context.Context
}

func (s *UserServiceTestSuite) SetupSuite() {
	var err error
	// 使用内存SQLite数据库进行测试
	s.db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		s.T().Fatal(err)
	}

	// 自动迁移
	err = s.db.AutoMigrate(&models.User{}, &models.TokenBlacklist{})
	if err != nil {
		s.T().Fatal(err)
	}

	// 创建服务
	authConfig := models.AuthConfig{
		JWTSecret:       "test-secret",
		AccessTokenTTL:  2 * time.Hour,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "rdp-api-test",
		Audience:        "rdp-users-test",
	}

	s.userService = NewUserService(s.db, authConfig)
	s.ctx = context.Background()
}

func (s *UserServiceTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	sqlDB.Close()
}

func (s *UserServiceTestSuite) SetupTest() {
	// 清空表
	s.db.Exec("DELETE FROM token_blacklists")
	s.db.Exec("DELETE FROM users")
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

// TestCreateUser 测试创建用户
func (s *UserServiceTestSuite) TestCreateUser() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
		Team:        models.TeamProductDev,
	}

	user, err := s.userService.CreateUser(s.ctx, req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), req.Username, user.Username)
	assert.Equal(s.T(), req.DisplayName, user.DisplayName)
	assert.Equal(s.T(), req.Email, user.Email)
	assert.Equal(s.T(), req.Role, user.Role)
	assert.NotEmpty(s.T(), user.ID)
	assert.NotEmpty(s.T(), user.PasswordHash)
	assert.True(s.T(), user.IsActive)

	// 验证密码已加密
	assert.NotEqual(s.T(), req.Password, user.PasswordHash)
}

// TestCreateUser_DuplicateUsername 测试重复用户名
func (s *UserServiceTestSuite) TestCreateUser_DuplicateUsername() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User",
		Email:       "test1@example.com",
		Role:        models.RoleDesigner,
	}

	_, err := s.userService.CreateUser(s.ctx, req)
	assert.NoError(s.T(), err)

	// 创建相同用户名的用户
	req2 := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User 2",
		Email:       "test2@example.com",
		Role:        models.RoleDesigner,
	}

	_, err = s.userService.CreateUser(s.ctx, req2)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrUsernameExists, err)
}

// TestCreateUser_DuplicateEmail 测试重复邮箱
func (s *UserServiceTestSuite) TestCreateUser_DuplicateEmail() {
	req := models.CreateUserRequest{
		Username:    "testuser1",
		Password:    "TestPass123",
		DisplayName: "Test User 1",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	_, err := s.userService.CreateUser(s.ctx, req)
	assert.NoError(s.T(), err)

	// 创建相同邮箱的用户
	req2 := models.CreateUserRequest{
		Username:    "testuser2",
		Password:    "TestPass123",
		DisplayName: "Test User 2",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	_, err = s.userService.CreateUser(s.ctx, req2)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrEmailExists, err)
}

// TestCreateUser_InvalidRole 测试无效角色
func (s *UserServiceTestSuite) TestCreateUser_InvalidRole() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        "invalid_role",
	}

	_, err := s.userService.CreateUser(s.ctx, req)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrInvalidRole, err)
}

// TestGetUserByID 测试根据ID获取用户
func (s *UserServiceTestSuite) TestGetUserByID() {
	// 先创建用户
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	created, _ := s.userService.CreateUser(s.ctx, req)

	// 获取用户
	user, err := s.userService.GetUserByID(s.ctx, created.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), created.ID, user.ID)
	assert.Equal(s.T(), created.Username, user.Username)
}

// TestGetUserByID_NotFound 测试用户不存在
func (s *UserServiceTestSuite) TestGetUserByID_NotFound() {
	_, err := s.userService.GetUserByID(s.ctx, "non-existent-id")
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrUserNotFound, err)
}

// TestGetUserByUsername 测试根据用户名获取用户
func (s *UserServiceTestSuite) TestGetUserByUsername() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	created, _ := s.userService.CreateUser(s.ctx, req)

	user, err := s.userService.GetUserByUsername(s.ctx, "testuser")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), created.ID, user.ID)
}

// TestGetUserList 测试获取用户列表
func (s *UserServiceTestSuite) TestGetUserList() {
	// 创建多个用户
	for i := 0; i < 5; i++ {
		req := models.CreateUserRequest{
			Username:    "testuser" + string(rune('0'+i)),
			Password:    "TestPass123",
			DisplayName: "Test User " + string(rune('0'+i)),
			Email:       "test" + string(rune('0'+i)) + "@example.com",
			Role:        models.RoleDesigner,
		}
		_, err := s.userService.CreateUser(s.ctx, req)
		assert.NoError(s.T(), err)
	}

	query := models.DefaultUserListQuery()
	query.Page = 1
	query.PageSize = 10

	resp, err := s.userService.GetUserList(s.ctx, query)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), int64(5), resp.Total)
	assert.Equal(s.T(), 1, resp.Page)
}

// TestGetUserList_WithKeyword 测试关键词搜索
func (s *UserServiceTestSuite) TestGetUserList_WithKeyword() {
	req1 := models.CreateUserRequest{
		Username:    "alice",
		Password:    "TestPass123",
		DisplayName: "Alice Smith",
		Email:       "alice@example.com",
		Role:        models.RoleDesigner,
	}
	_, _ = s.userService.CreateUser(s.ctx, req1)

	req2 := models.CreateUserRequest{
		Username:    "bob",
		Password:    "TestPass123",
		DisplayName: "Bob Johnson",
		Email:       "bob@example.com",
		Role:        models.RoleDesigner,
	}
	_, _ = s.userService.CreateUser(s.ctx, req2)

	query := models.DefaultUserListQuery()
	query.Keyword = "alice"

	resp, err := s.userService.GetUserList(s.ctx, query)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), resp.Total)
}

// TestUpdateUser 测试更新用户
func (s *UserServiceTestSuite) TestUpdateUser() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	created, _ := s.userService.CreateUser(s.ctx, req)

	updateReq := models.UpdateUserRequest{
		DisplayName: "Updated Name",
		Bio:         "New bio",
	}

	updated, err := s.userService.UpdateUser(s.ctx, created.ID, updateReq, false)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Name", updated.DisplayName)
	assert.Equal(s.T(), "New bio", updated.Bio)
}

// TestUpdateUser_AsAdmin 测试管理员更新用户
func (s *UserServiceTestSuite) TestUpdateUser_AsAdmin() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	created, _ := s.userService.CreateUser(s.ctx, req)

	updateReq := models.UpdateUserRequest{
		Role:     models.RoleTeamLeader,
		IsActive: boolPtr(false),
	}

	updated, err := s.userService.UpdateUser(s.ctx, created.ID, updateReq, true)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), models.RoleTeamLeader, updated.Role)
	assert.False(s.T(), updated.IsActive)
}

// TestDeleteUser 测试删除用户
func (s *UserServiceTestSuite) TestDeleteUser() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	created, _ := s.userService.CreateUser(s.ctx, req)

	err := s.userService.DeleteUser(s.ctx, created.ID)
	assert.NoError(s.T(), err)

	// 确认用户已删除
	_, err = s.userService.GetUserByID(s.ctx, created.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrUserNotFound, err)
}

// TestDeleteUser_Admin 测试不能删除管理员
func (s *UserServiceTestSuite) TestDeleteUser_Admin() {
	req := models.CreateUserRequest{
		Username:    "adminuser",
		Password:    "TestPass123",
		DisplayName: "Admin User",
		Email:       "admin@example.com",
		Role:        models.RoleAdmin,
	}

	created, _ := s.userService.CreateUser(s.ctx, req)

	err := s.userService.DeleteUser(s.ctx, created.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrCannotDeleteAdmin, err)
}

// TestValidateCredentials 测试验证凭据
func (s *UserServiceTestSuite) TestValidateCredentials() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	_, _ = s.userService.CreateUser(s.ctx, req)

	// 正确凭据
	user, err := s.userService.ValidateCredentials(s.ctx, "testuser", "TestPass123")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)

	// 错误密码
	_, err = s.userService.ValidateCredentials(s.ctx, "testuser", "WrongPass")
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrInvalidCredentials, err)

	// 不存在的用户
	_, err = s.userService.ValidateCredentials(s.ctx, "nonexistent", "TestPass123")
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrInvalidCredentials, err)
}

// TestValidateCredentials_DisabledUser 测试禁用用户无法登录
func (s *UserServiceTestSuite) TestValidateCredentials_DisabledUser() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	created, _ := s.userService.CreateUser(s.ctx, req)

	// 禁用用户
	updateReq := models.UpdateUserRequest{
		IsActive: boolPtr(false),
	}
	_, _ = s.userService.UpdateUser(s.ctx, created.ID, updateReq, true)

	// 尝试登录
	_, err := s.userService.ValidateCredentials(s.ctx, "testuser", "TestPass123")
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrUserDisabled, err)
}

// TestGenerateTokenPair 测试生成令牌对
func (s *UserServiceTestSuite) TestGenerateTokenPair() {
	user := &models.User{
		ID:          "test-id",
		Username:    "testuser",
		DisplayName: "Test User",
		Role:        models.RoleDesigner,
	}

	resp, err := s.userService.GenerateTokenPair(user)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), resp.AccessToken)
	assert.NotEmpty(s.T(), resp.RefreshToken)
	assert.Equal(s.T(), "Bearer", resp.TokenType)
	assert.Equal(s.T(), 7200, resp.ExpiresIn)
	assert.NotNil(s.T(), resp.User)
}

// TestRefreshAccessToken 测试刷新访问令牌
func (s *UserServiceTestSuite) TestRefreshAccessToken() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "TestPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	created, _ := s.userService.CreateUser(s.ctx, req)

	// 生成令牌对
	resp, _ := s.userService.GenerateTokenPair(created)

	// 刷新令牌
	refreshResp, err := s.userService.RefreshAccessToken(s.ctx, resp.RefreshToken)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), refreshResp.AccessToken)
	assert.Equal(s.T(), "Bearer", refreshResp.TokenType)
}

// TestRefreshAccessToken_InvalidToken 测试无效令牌刷新
func (s *UserServiceTestSuite) TestRefreshAccessToken_InvalidToken() {
	_, err := s.userService.RefreshAccessToken(s.ctx, "invalid-token")
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrInvalidToken, err)
}

// TestLogout 测试登出
func (s *UserServiceTestSuite) TestLogout() {
	user := &models.User{
		ID:          "test-id",
		Username:    "testuser",
		DisplayName: "Test User",
		Role:        models.RoleDesigner,
	}

	resp, _ := s.userService.GenerateTokenPair(user)

	// 登出
	err := s.userService.Logout(s.ctx, resp.AccessToken)
	assert.NoError(s.T(), err)

	// 验证令牌已被加入黑名单
	_, err = s.userService.ValidateToken(s.ctx, resp.AccessToken)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrTokenBlacklisted, err)
}

// TestUpdatePassword 测试更新密码
func (s *UserServiceTestSuite) TestUpdatePassword() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "OldPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	created, _ := s.userService.CreateUser(s.ctx, req)

	// 更新密码
	err := s.userService.UpdatePassword(s.ctx, created.ID, "OldPass123", "NewPass456")
	assert.NoError(s.T(), err)

	// 使用新密码登录
	_, err = s.userService.ValidateCredentials(s.ctx, "testuser", "NewPass456")
	assert.NoError(s.T(), err)

	// 使用旧密码登录应该失败
	_, err = s.userService.ValidateCredentials(s.ctx, "testuser", "OldPass123")
	assert.Error(s.T(), err)
}

// TestResetPassword 测试重置密码
func (s *UserServiceTestSuite) TestResetPassword() {
	req := models.CreateUserRequest{
		Username:    "testuser",
		Password:    "OldPass123",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Role:        models.RoleDesigner,
	}

	created, _ := s.userService.CreateUser(s.ctx, req)

	// 重置密码
	err := s.userService.ResetPassword(s.ctx, created.ID, "NewPass456")
	assert.NoError(s.T(), err)

	// 使用新密码登录
	_, err = s.userService.ValidateCredentials(s.ctx, "testuser", "NewPass456")
	assert.NoError(s.T(), err)
}

// TestIsAdmin 测试管理员检查
func (s *UserServiceTestSuite) TestIsAdmin() {
	adminUser := &models.User{Role: models.RoleAdmin}
	assert.True(s.T(), adminUser.IsAdmin())

	normalUser := &models.User{Role: models.RoleDesigner}
	assert.False(s.T(), normalUser.IsAdmin())
}

// TestHasRole 测试角色权限检查
func (s *UserServiceTestSuite) TestHasRole() {
	adminUser := &models.User{Role: models.RoleAdmin}
	assert.True(s.T(), adminUser.HasRole(models.RoleDesigner))
	assert.True(s.T(), adminUser.HasRole(models.RoleAdmin))

	teamLeader := &models.User{Role: models.RoleTeamLeader}
	assert.True(s.T(), teamLeader.HasRole(models.RoleDesigner))
	assert.True(s.T(), teamLeader.HasRole(models.RoleTeamLeader))
	assert.False(s.T(), teamLeader.HasRole(models.RoleDeptLeader))

	designer := &models.User{Role: models.RoleDesigner}
	assert.True(s.T(), designer.HasRole(models.RoleDesigner))
	assert.False(s.T(), designer.HasRole(models.RoleTeamLeader))
}

// Helper function
func boolPtr(b bool) *bool {
	return &b
}
