# InfraAgent - 冲突预警表

> **重要**: 其他Agent在开始工作前，请检查此表以避免冲突

---

## 已占用的文件/目录

以下文件/目录由 InfraAgent 创建，**其他Agent不应修改或覆盖**:

### 配置文件 (全局)
| 文件 | 说明 | 修改需协商 |
|------|------|------------|
| `.editorconfig` | 编辑器配置 | 是 |
| `.gitignore` | Git忽略规则 | 是 |
| `Makefile` | 构建命令 | 是 |

### 前端配置 (apps/web/)
| 文件 | 说明 | 修改需协商 |
|------|------|------------|
| `package.json` | 依赖版本 | **必须** |
| `tsconfig.json` | TS配置 | 是 |
| `vite.config.ts` | Vite配置 | 是 |
| `.eslintrc.cjs` | ESLint规则 | 是 |
| `.prettierrc` | 格式化规则 | 是 |

### 前端共享代码 (apps/web/src/)
| 文件 | 说明 | 使用方式 |
|------|------|----------|
| `types/index.ts` | 类型定义 | **追加**，不要修改 |
| `services/api.ts` | API客户端 | 参考示例扩展 |
| `stores/auth.ts` | 认证状态 | 使用，不要修改 |
| `App.tsx` | 根组件 | **追加路由**，不要修改现有 |
| `main.tsx` | 入口 | 不要修改 |

### 后端配置 (services/api/)
| 文件 | 说明 | 修改需协商 |
|------|------|------------|
| `go.mod` | Go模块 | **必须** |
| `.golangci.yml` | Lint规则 | 是 |
| `main.go` | 入口 | 追加路由 |

### 后端共享代码 (services/api/)
| 文件 | 说明 | 使用方式 |
|------|------|----------|
| `handlers/response.go` | 响应格式 | 使用函数，不要修改 |
| `handlers/health.go` | 健康检查 | 不要修改 |
| `utils/id.go` | ID生成 | 使用函数，不要修改 |

### 数据库 (database/)
| 文件 | 说明 | 使用方式 |
|------|------|----------|
| `schema/enums.sql` | 枚举类型 | **使用**，不要修改 |
| `migrations/000_*.sql` | 初始化 | 不要修改 |
| `seeds/*.sql` | 种子数据 | 可追加新文件 |

### 部署配置 (deploy/)
| 文件 | 说明 | 修改需协商 |
|------|------|------------|
| `systemd/*.service` | 服务配置 | 是 |
| `nginx/*.conf` | Nginx配置 | 是 |
| `scripts/install.sh` | 安装脚本 | 是 |
| `scripts/backup.sh` | 备份脚本 | 是 |
| `scripts/health-check.sh` | 健康检查 | 是 |

### 文档 (agents/outputs/)
| 文件 | 说明 | 权限 |
|------|------|------|
| `infra-scaffold/DELIVERY.md` | 交付文档 | 只读 |
| `COLLABORATION_GUIDE.md` | 协作规范 | 可追加 |

---

## 安全的扩展区域

以下区域是**安全的**，各Agent可以自由创建文件:

### PortalAgent (P1-T1 ~ P1-T4)
```
apps/web/src/pages/portal/
apps/web/src/components/portal/
apps/web/src/components/workbench/
apps/web/src/components/notification/
apps/web/src/components/search/
```

### UserAgent (P1-T5 ~ P1-T8)
```
services/api/models/user.go
services/api/models/organization.go
services/api/handlers/user.go
services/api/handlers/organization.go
services/api/services/auth.go
services/api/services/permission.go
services/api/middleware/auth.go
database/migrations/001_users.sql
database/migrations/001_organizations.sql
database/migrations/001_permissions.sql
apps/web/src/pages/profile/
apps/web/src/components/org/
```

### ProjectAgent (P1-T9 ~ P1-T12)
```
services/api/models/project.go
services/api/models/process_template.go
services/api/models/file.go
services/api/handlers/project.go
services/api/handlers/process_template.go
services/api/handlers/file.go
services/api/services/project.go
database/migrations/002_projects.sql
database/migrations/002_process_templates.sql
database/migrations/002_files.sql
apps/web/src/pages/projects/
apps/web/src/components/projects/
apps/web/src/components/wizard/
```

### SecurityAgent (P1-T13 ~ P1-T16)
```
services/api/models/classification.go
services/api/models/audit.go
services/api/handlers/classification.go
services/api/handlers/audit.go
services/api/services/classification.go
services/api/services/session.go
services/api/middleware/session.go
database/migrations/003_classification.sql
database/migrations/004_audit_logs.sql
apps/web/src/components/security/
apps/web/src/pages/admin/
apps/web/src/utils/watermark.ts
```

### InfraAgent (P1-T17 ~ P1-T20)
```
database/init.sql
database/migrations/003_*.sql
database/migrations/004_*.sql
config/*.yaml
deploy/scripts/servicectl.sh
```

---

## 需要协商的变更

如果必须修改以下配置，请通过此流程:

1. **创建 Issue 记录**:
   ```markdown
   ## 变更请求
   - Agent: [YourAgent]
   - 目标文件: [file path]
   - 变更原因: [reason]
   - 影响范围: [impact]
   ```

2. **通知相关Agent**:
   - 在 CONFLICTS.md 中记录
   - 等待 24 小时反馈

3. **决策**:
   - 无反对 → 实施变更
   - 有反对 → 提交 Architect Agent 裁决

---

## 紧急联系方式

如果出现冲突:
1. 检查本表确认是否为已占用资源
2. 查阅 `agents/outputs/infra-scaffold/DELIVERY.md`
3. 查阅 `agents/outputs/COLLABORATION_GUIDE.md`
4. 仍无法解决 → 上报人类监督者

---

*最后更新: 2026-02-22*
