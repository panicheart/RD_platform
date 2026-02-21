# Testing Guide

## 测试框架

- **单元测试**: Go testing + Testify
- **HTTP测试**: Gin test mode + httptest
- **覆盖率**: Go内置覆盖率工具
- **最低覆盖率**: 60%

## 快速开始

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./handlers

# 运行特定测试
go test -run TestHealthCheck ./handlers

# 运行测试并查看覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 测试规范

### 文件命名
- 测试文件: `*_test.go`
- 与源文件同目录

### 测试结构

```go
func TestFunctionName(t *testing.T) {
    t.Run("should do something", func(t *testing.T) {
        // Arrange
        input := "test"
        
        // Act
        result := FunctionName(input)
        
        // Assert
        assert.Equal(t, expected, result)
    })
}
```

### HTTP Handler 测试

```go
func TestYourHandler(t *testing.T) {
    // Setup
    router := setupTestRouter()
    router.GET("/endpoint", YourHandler)
    
    // Execute
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/endpoint", nil)
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## 常用断言

- `assert.Equal(t, expected, actual)`
- `assert.NoError(t, err)`
- `assert.Error(t, err)`
- `assert.Len(t, slice, length)`
- `assert.NotEmpty(t, value)`
- `require.NoError(t, err)` - 错误时立即停止

## Mock数据

使用 testify/mock 进行接口mock:

```go
import "github.com/stretchr/testify/mock"

type MockService struct {
    mock.Mock
}

func (m *MockService) GetUser(id string) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}
```

## 测试数据库

使用内存SQLite或测试容器:

```go
func setupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    db.AutoMigrate(&models.User{})
    return db
}
```

## CI集成

测试在CI中自动运行:

```bash
# Makefile
test-backend:
    cd services/api && go test -v ./...

test-coverage:
    cd services/api && go test -coverprofile=coverage.out ./...
```
