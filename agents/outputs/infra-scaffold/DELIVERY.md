# InfraAgent - 项目骨架交付文档

## 交付信息

| 项目 | 内容 |
|------|------|
| **Agent** | InfraAgent |
| **任务** | 项目骨架初始化 (P0-T0) |
| **状态** | ✅ 已完成 |
| **交付日期** | 2026-02-22 |
| **检查者** | 待 Reviewer Agent 审查 |

---

## 交付物清单

### 1. 前端脚手架 (`apps/web/`)

| 文件 | 说明 | 依赖 |
|------|------|------|
| `package.json` | 依赖配置 (React 18, AntD 5, Vite 5, Zustand) | 无 |
| `tsconfig.json` | TypeScript 配置，含路径别名 | 无 |
| `vite.config.ts` | Vite 配置，含代理规则 | 无 |
| `.eslintrc.cjs` | ESLint 规则 (TypeScript + React) | 无 |
| `.prettierrc` | 代码格式化配置 | 无 |
| `src/types/index.ts` | **核心类型定义** - 所有Agent必须使用 | 无 |
| `src/services/api.ts` | **API客户端模板** - 后端响应格式参考 | 无 |
| `src/stores/auth.ts` | 认证状态管理 | 无 |
| `src/App.tsx` | 路由占位 | 无 |

**路径别名规范** (所有前端Agent必须遵守):
```typescript
import { User } from '@types';           // 类型定义
import { apiClient } from '@services/api'; // API客户端
import { useAuthStore } from '@stores/auth'; // 状态管理
```

### 2. 后端脚手架 (`services/api/`)

| 文件 | 说明 | 依赖 |
|------|------|------|
| `go.mod` | Go模块定义 (Gin, GORM, ULID) | 无 |
| `.golangci.yml` | Lint规则 | 无 |
| `main.go` | 入口，含健康检查端点 | 无 |
| `handlers/response.go` | **统一响应格式** - 所有Handler必须使用 | 无 |
| `handlers/health.go` | 健康检查示例 | 无 |
| `utils/id.go` | **ULID生成器** - 所有模型必须使用 | 无 |

**API响应格式** (所有后端Agent必须遵守):
```go
{
  "code": 200,
  "message": "success", 
  "data": { ... }
}
```

**ID生成规范**:
```go
import "rdp-platform/rdp-api/utils"

id := utils.GenerateULID()  // 所有实体ID必须使用ULID
```

### 3. 数据库配置 (`database/`)

| 文件 | 说明 | 依赖 |
|------|------|------|
| `schema/enums.sql` | **枚举类型定义** - 所有表必须使用 | 无 |
| `seeds/process_templates.sql` | 7种项目类别模板 | 无 |
| `seeds/roles.sql` | 默认角色权限 | 无 |

**枚举类型** (所有Agent必须遵守):
- `user_status`: active, inactive, locked, pending
- `project_status`: draft, planning, in_progress, on_hold, completed, cancelled, archived
- `project_category`: new_product, product_improvement, pre_research, tech_platform, component_development, process_improvement, other
- `classification_level`: public, internal, confidential, secret
- `priority_level`: low, medium, high, critical

### 4. 部署配置 (`deploy/`)

| 文件 | 说明 | 依赖 |
|------|------|------|
| `systemd/rdp-api.service` | API服务配置 | 无 |
| `systemd/rdp-casdoor.service` | Casdoor服务配置 | 无 |
| `nginx/nginx.conf` | Nginx主配置 | 无 |
| `nginx/sites-available/rdp.conf` | RDP站点配置 | 无 |
| `scripts/install.sh` | 一键安装脚本 | 无 |
| `scripts/backup.sh` | 数据库备份脚本 | 无 |
| `scripts/health-check.sh` | 健康检查脚本 | 无 |

### 5. 通用配置

| 文件 | 说明 | 依赖 |
|------|------|------|
| `.editorconfig` | 编辑器统一配置 | 无 |
| `.gitignore` | Git忽略规则 | 无 |
| `Makefile` | 常用命令 | 无 |
| `.vscode/settings.json` | VS Code工作区配置 | 无 |

---

## 依赖关系

### 本交付物依赖
- 无 (这是项目第一层)

### 依赖本交付物的任务

| Agent | 任务 | 依赖项 |
|-------|------|--------|
| PortalAgent | P1-T1 ~ P1-T4 | 前端脚手架、类型定义 |
| UserAgent | P1-T5 ~ P1-T8 | 后端结构、枚举类型、响应格式 |
| ProjectAgent | P1-T9 ~ P1-T12 | 后端结构、枚举类型、响应格式 |
| SecurityAgent | P1-T13 ~ P1-T16 | 后端结构、中间件目录 |
| InfraAgent | P1-T17 ~ P1-T20 | 数据库目录、部署脚本模板 |

---

## 使用指南

### 前端Agent使用规范

1. **安装依赖**:
   ```bash
   cd apps/web && npm install
   ```

2. **类型定义**: 在 `src/types/index.ts` 中补充你的类型

3. **API服务**: 参考 `src/services/api.ts` 中的示例，在同级目录创建你的服务
   ```typescript
   export const yourModuleAPI = {
     getItems: () => apiClient.get('/your-module'),
   };
   ```

4. **页面组件**: 在 `src/pages/` 下创建目录，参考 `App.tsx` 中的路由配置

### 后端Agent使用规范

1. **初始化Go模块**:
   ```bash
   cd services/api && go mod tidy
   ```

2. **模型定义**: 在 `models/` 目录创建模型文件，使用ULID:
   ```go
   type YourModel struct {
       ID        string    `gorm:"primaryKey"`
       // ...
   }
   
   func (m *YourModel) BeforeCreate(tx *gorm.DB) error {
       if m.ID == "" {
           m.ID = utils.GenerateULID()
       }
       return nil
   }
   ```

3. **Handler实现**: 在 `handlers/` 目录创建，使用统一响应:
   ```go
   func YourHandler(c *gin.Context) {
       data := ...
       handlers.SuccessResponse(c, data)
   }
   ```

4. **路由注册**: 在 `main.go` 中添加路由

---

## 冲突处理预案

### 场景1: 文件路径冲突
**如果其他Agent创建了相同路径的文件**:
1. 检查是否为预期行为（如多个Agent都需要 types/index.ts）
2. 如果是扩展，合并内容
3. 如果是冲突，通知 Architect Agent 裁决

### 场景2: 依赖版本冲突
**如果 package.json/go.mod 依赖版本不一致**:
- 前端: 以本交付物的版本为基准 (React 18, AntD 5)
- 后端: 以本交付物的版本为基准 (Go 1.22, Gin 1.9)

### 场景3: 编码规范冲突
**如果其他Agent使用了不同的规范**:
- 必须遵守本交付物中的 ESLint/golangci-lint 配置
- 如有异议，提交 Reviewer Agent 审查

---

## 审查检查清单

Reviewer Agent 请检查:

- [ ] 前端配置完整，可执行 `npm install` 和 `npm run dev`
- [ ] 后端配置完整，可执行 `go mod tidy` 和 `go run main.go`
- [ ] 类型定义与需求规格一致
- [ ] 枚举类型覆盖所有需求场景
- [ ] 部署脚本语法正确
- [ ] Makefile 命令可用

---

## 后续扩展建议

1. **待 PortalAgent 完成后**: 添加前端路由守卫、布局组件
2. **待 UserAgent 完成后**: 添加 JWT 验证中间件
3. **待 InfraAgent 完成后**: 添加完整的数据库迁移脚本
4. **待 SecurityAgent 完成后**: 添加审计日志中间件

---

*交付文档版本: 1.0*
*最后更新: 2026-02-22*
