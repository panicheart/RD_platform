# 用户管理API 实现文件清单

## 任务 P1-B1: 用户管理API 完成总结

### 创建/更新的文件

#### 1. 数据模型
- **文件**: `services/api/models/user.go`
- **内容**: User结构体定义、表名方法、密码加密方法、JWT claims结构体
- **关键类型**: UserRole, TeamType, TitleLevel, User, UserStats, UserProfile

#### 2. 认证相关模型
- **文件**: `services/api/models/auth.go`
- **内容**: JWTClaims, LoginRequest/Response, CreateUserRequest, UpdateUserRequest等

#### 3. 项目模型（依赖）
- **文件**: `services/api/models/project.go`
- **内容**: Project模型定义（UserService依赖）

#### 4. 服务层
- **文件**: `services/api/services/user_service.go`
- **方法**:
  - CreateUser - 创建用户
  - GetUserByID - 根据ID获取用户
  - GetUserList - 获取用户列表（支持分页、筛选）
  - UpdateUser - 更新用户
  - DeleteUser - 删除用户
  - ValidateCredentials - 验证登录凭据
  - GenerateTokenPair - 生成JWT令牌对
  - RefreshAccessToken - 刷新访问令牌
  - Logout - 用户登出

#### 5. HTTP处理器
- **文件**: `services/api/handlers/user.go`
- **方法**:
  - Register - 注册处理器
  - Login - 登录处理器
  - RefreshToken - 刷新令牌
  - Logout - 登出处理器
  - GetUsers - 获取用户列表
  - GetUser - 获取用户详情
  - CreateUser - 创建用户（Admin）
  - UpdateUser - 更新用户
  - DeleteUser - 删除用户（Admin）
  - GetUserProfile - 获取用户Profile
  - GetCurrentUser - 获取当前用户
  - ChangePassword - 修改密码

#### 6. 认证中间件
- **文件**: `services/api/middleware/auth.go`
- **功能**:
  - JWTAuth - JWT验证中间件
  - RequireRole - 角色权限检查中间件
  - RequireAdmin - 管理员权限检查
  - RequireAdminOrSelf - 管理员或本人检查
  - OptionalAuth - 可选认证
  - SecurityHeaders - 安全头部中间件
  - CORS - 跨域中间件

#### 7. 路由配置
- **文件**: `services/api/routes/routes.go`
- **内容**: 注册用户相关路由和认证路由

#### 8. 配置模块
- **文件**: `services/api/config/config.go`
- **内容**: 应用配置、数据库配置、认证配置

#### 9. 主入口
- **文件**: `services/api/main.go`
- **内容**: 初始化Gin引擎、数据库连接、注册路由、启动服务

#### 10. 健康检查处理器
- **文件**: `services/api/handlers/health.go`
- **内容**: HealthCheck, HealthDetailed, ReadinessCheck, LivenessCheck

#### 11. 单元测试 - 服务层
- **文件**: `services/api/services/user_service_test.go`
- **测试覆盖**:
  - TestCreateUser - 创建用户
  - TestCreateUser_DuplicateUsername - 重复用户名
  - TestCreateUser_DuplicateEmail - 重复邮箱
  - TestCreateUser_InvalidRole - 无效角色
  - TestGetUserByID - 获取用户
  - TestGetUserByID_NotFound - 用户不存在
  - TestGetUserList - 用户列表
  - TestGetUserList_WithKeyword - 关键词搜索
  - TestUpdateUser - 更新用户
  - TestDeleteUser - 删除用户
  - TestDeleteUser_Admin - 不能删除管理员
  - TestValidateCredentials - 验证凭据
  - TestGenerateTokenPair - 生成令牌
  - TestRefreshAccessToken - 刷新令牌
  - TestLogout - 登出
  - TestUpdatePassword - 更新密码
  - TestResetPassword - 重置密码

#### 12. 单元测试 - 处理器
- **文件**: `services/api/handlers/user_test.go`
- **测试覆盖**:
  - TestLogin_Success
  - TestLogin_InvalidCredentials
  - TestGetUser_Success
  - TestCreateUser_Success
  - TestGetUserList_Success
  - TestUpdateUser_Success
  - TestDeleteUser_Success
  - TestRefreshToken_Success
  - TestChangePassword_Success
  - 响应方法测试

### API端点列表

| 方法 | 路径 | 描述 | 权限 |
|------|------|------|------|
| POST | /api/v1/auth/login | 登录 | 公开 |
| POST | /api/v1/auth/refresh | 刷新Token | 公开 |
| POST | /api/v1/auth/logout | 登出 | 认证用户 |
| GET | /api/v1/users | 用户列表 | 认证用户 |
| GET | /api/v1/users/me | 当前用户 | 认证用户 |
| POST | /api/v1/users | 创建用户 | Admin |
| GET | /api/v1/users/:id | 用户详情 | 认证用户 |
| GET | /api/v1/users/:id/profile | 用户Profile | 认证用户 |
| PUT | /api/v1/users/:id | 更新用户 | Admin/本人 |
| DELETE | /api/v1/users/:id | 删除用户 | Admin |
| PUT | /api/v1/users/:id/password | 修改密码 | Admin/本人 |

### 技术栈

- Go 1.22+
- Gin 1.9+
- GORM (PostgreSQL)
- golang-jwt/jwt/v5
- golang.org/x/crypto/bcrypt
- oklog/ulid/v2 (ULID生成)

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| RDP_ENV | development | 运行环境 |
| RDP_API_PORT | 8080 | API服务端口 |
| RDP_DB_HOST | localhost | 数据库主机 |
| RDP_DB_PORT | 5432 | 数据库端口 |
| RDP_DB_USER | rdp_user | 数据库用户 |
| RDP_DB_PASSWORD | rdp_password | 数据库密码 |
| RDP_DB_NAME | rdp_db | 数据库名 |
| RDP_JWT_SECRET | change-this-secret-in-production | JWT密钥 |
| RDP_ACCESS_TOKEN_TTL | 2h | Access Token有效期 |
| RDP_REFRESH_TOKEN_TTL | 168h | Refresh Token有效期 |

### 默认管理员

系统启动时会自动创建默认管理员：
- 用户名: admin
- 密码: Admin@123
- 角色: admin

