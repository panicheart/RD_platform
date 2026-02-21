# InfraAgent - 测试框架交付文档

## 交付信息

| 项目 | 内容 |
|------|------|
| **Agent** | InfraAgent |
| **任务** | 测试框架搭建 (P0-T1) |
| **状态** | ✅ 已完成 |
| **交付日期** | 2026-02-22 |

---

## 交付物清单

### 1. 前端测试框架

| 文件 | 说明 | 使用方式 |
|------|------|----------|
| `vitest.config.ts` | Vitest配置，含覆盖率阈值 | 运行 `npm run test` |
| `src/test/setup.ts` | 测试环境初始化 | 自动加载 |
| `src/test/README.md` | 测试指南文档 | 参考文档 |
| `src/utils/id.test.ts` | 工具函数测试示例 | 复制修改 |
| `src/types/index.test.ts` | 类型定义测试示例 | 复制修改 |

**更新到 package.json**:
```json
"scripts": {
  "test": "vitest",
  "test:ui": "vitest --ui",
  "test:coverage": "vitest --coverage"
}
```

**新增依赖**:
- vitest
- @testing-library/react
- @testing-library/jest-dom
- @testing-library/user-event
- @vitest/coverage-v8
- @vitest/ui
- jsdom
- msw (Mock Service Worker)

### 2. 后端测试框架

| 文件 | 说明 | 使用方式 |
|------|------|----------|
| `utils/id_test.go` | ULID工具测试 | 复制修改 |
| `handlers/response_test.go` | HTTP响应测试 | 复制修改 |
| `handlers/README_test.md` | Go测试指南 | 参考文档 |

**更新到 go.mod**:
```go
require (
    github.com/stretchr/testify v1.8.4
    // ... other deps
)
```

### 3. 测试脚本

**Makefile 新增命令**:
```makefile
test-frontend-watch:       # 监听模式运行测试
test-frontend-coverage:    # 生成覆盖率报告
test-frontend-ui:          # 图形界面运行测试
test-backend-coverage:     # 后端覆盖率
test-backend-coverage-html:# HTML格式覆盖率报告
test-coverage:             # 全量覆盖率
```

### 4. CI/CD配置

| 文件 | 说明 |
|------|------|
| `.github/workflows/ci.yml` | 持续集成：测试、构建、覆盖率 |
| `.github/workflows/security.yml` | 安全扫描：Trivy + CodeQL |

**CI流程**:
1. 前端测试 (type check + lint + test + coverage)
2. 后端测试 (lint + test with Postgres + coverage)
3. 构建检查 (确保能成功构建)
4. 安全扫描 (依赖漏洞 + 代码分析)

---

## 使用指南

### 前端Agent

1. **编写测试**:
   ```typescript
   import { describe, it, expect } from 'vitest';
   import { render, screen } from '@testing-library/react';
   
   describe('YourComponent', () => {
     it('should render', () => {
       render(<YourComponent />);
       expect(screen.getByText('Hello')).toBeInTheDocument();
     });
   });
   ```

2. **运行测试**:
   ```bash
   npm run test           # 单次运行
   npm run test --watch   # 监听模式
   ```

3. **覆盖率要求**: 最低60%

### 后端Agent

1. **编写测试**:
   ```go
   func TestYourHandler(t *testing.T) {
       router := setupTestRouter()
       router.GET("/test", YourHandler)
       
       w := httptest.NewRecorder()
       req, _ := http.NewRequest("GET", "/test", nil)
       router.ServeHTTP(w, req)
       
       assert.Equal(t, http.StatusOK, w.Code)
   }
   ```

2. **运行测试**:
   ```bash
   go test ./...
   go test -cover ./...
   ```

3. **覆盖率要求**: 最低60%

---

## 依赖关系

### 依赖本交付物的任务

| Agent | 任务 | 说明 |
|-------|------|------|
| PortalAgent | P1-T1~T4 | 使用前端测试框架编写组件测试 |
| UserAgent | P1-T5~T8 | 使用后端测试框架编写API测试 |
| ProjectAgent | P1-T9~T12 | 使用测试框架编写业务逻辑测试 |
| SecurityAgent | P1-T13~P1-T16 | 编写安全相关测试 |
| All Agents | - | CI/CD会自动运行所有测试 |

---

## 覆盖率阈值

| 指标 | 最低要求 | 目标 |
|------|----------|------|
| Lines | 60% | 80% |
| Functions | 60% | 80% |
| Branches | 50% | 70% |
| Statements | 60% | 80% |

---

## 检查清单

Reviewer Agent 请检查:

- [ ] `npm run test` 在 apps/web 目录可正常运行
- [ ] `go test ./...` 在 services/api 目录可正常运行
- [ ] 示例测试文件都能通过
- [ ] CI配置语法正确 (通过 `actionlint` 或 GitHub验证)
- [ ] Makefile 测试命令可用

---

*交付文档版本: 1.0*
*最后更新: 2026-02-22*
