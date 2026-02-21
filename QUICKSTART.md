# RDP Platform - 项目骨架快速启动

## 已完成的准备工作

✅ **项目目录结构** - 前后端、数据库、部署配置框架  
✅ **前端脚手架** - React + TS + Vite + Ant Design  
✅ **后端脚手架** - Go + Gin + GORM  
✅ **编码规范** - ESLint + golangci-lint  
✅ **数据库框架** - PostgreSQL 迁移、枚举类型、种子数据  
✅ **部署配置** - systemd + Nginx + 安装脚本  

---

## 各Agent快速开始

### PortalAgent - 前端页面开发

```bash
# 1. 进入前端目录
cd apps/web

# 2. 安装依赖
npm install

# 3. 启动开发服务器
npm run dev
# 访问 http://localhost:3000
```

**你的任务区域**:
- `src/pages/portal/` - 门户首页
- `src/pages/workbench/` - 个人工作台
- `src/components/portal/` - 门户组件
- `src/components/workbench/` - 工作台组件
- `src/components/notification/` - 通知组件
- `src/components/search/` - 搜索组件

**类型定义参考**: `src/types/index.ts`  
**API调用参考**: `src/services/api.ts`

---

### UserAgent - 用户管理后端

```bash
# 1. 进入后端目录
cd services/api

# 2. 初始化Go模块
go mod tidy

# 3. 运行服务
go run main.go
# API 运行在 http://localhost:8080
```

**你的任务区域**:
- `models/user.go` - 用户模型
- `models/organization.go` - 组织架构模型
- `handlers/user.go` - 用户API
- `handlers/organization.go` - 组织API
- `services/auth.go` - 认证服务
- `services/permission.go` - 权限服务
- `middleware/auth.go` - 认证中间件
- `database/migrations/001_*.sql` - 用户相关表

**规范参考**:
- ID生成: `utils/id.go` (使用 `GenerateULID()`)
- 响应格式: `handlers/response.go` (使用 `SuccessResponse()`)
- 枚举类型: `database/schema/enums.sql`

---

### ProjectAgent - 项目管理

```bash
# 同上，使用已搭建好的后端框架
```

**你的任务区域**:
- `models/project.go` - 项目模型
- `models/process_template.go` - 流程模板模型
- `models/file.go` - 文件模型
- `handlers/project.go` - 项目API
- `handlers/process_template.go` - 模板API
- `handlers/file.go` - 文件API
- `services/project.go` - 项目服务
- `database/migrations/002_*.sql` - 项目相关表

---

### SecurityAgent - 安全合规

```bash
# 同上
```

**你的任务区域**:
- `models/classification.go` - 分级模型
- `models/audit.go` - 审计日志模型
- `handlers/classification.go` - 分级API
- `handlers/audit.go` - 审计API
- `services/classification.go` - 分级服务
- `services/session.go` - 会话服务
- `middleware/session.go` - 会话中间件
- `database/migrations/003_*.sql` - 分级表
- `database/migrations/004_*.sql` - 审计表

---

### InfraAgent - 基础设施完善

**你的任务区域**:
- `database/init.sql` - 数据库初始化脚本
- `database/migrations/` - 补充迁移脚本
- `config/` - 配置文件完善
- `deploy/scripts/` - 部署脚本完善

---

## 协作要点

### 1. 共享资源 (谨慎修改)

**类型定义** (`apps/web/src/types/index.ts`):
- ✅ 可以追加新类型
- ❌ 不要修改已有类型

**API客户端** (`apps/web/src/services/api.ts`):
- ✅ 参考示例创建你的服务文件
- ❌ 不要修改 `apiClient` 基础类

**数据库枚举** (`database/schema/enums.sql`):
- ✅ 使用已有枚举
- ❌ 不要修改枚举定义

### 2. 规范遵守

**前端**:
- 代码注释使用英文
- UI文案使用中文
- 遵循 ESLint 规则

**后端**:
- 代码注释使用英文
- 使用 ULID 生成 ID
- 使用统一响应格式

### 3. 冲突避免

**检查冲突预警表**:
```
agents/outputs/infra-scaffold/CONFLICT_WARNING.md
```

**如果发现冲突**:
1. 查阅 `COLLABORATION_GUIDE.md`
2. 在 `CONFLICTS.md` 中记录
3. 协商解决或上报裁决

---

## 常用命令

```bash
# 前端开发
cd apps/web
npm install        # 安装依赖
npm run dev        # 开发模式
npm run build      # 构建
npm run lint       # 代码检查
npm run format     # 代码格式化

# 后端开发
cd services/api
go mod tidy        # 下载依赖
go run main.go     # 运行
go build -o bin/api main.go  # 构建
golangci-lint run  # 代码检查

# 项目构建
make install       # 安装所有依赖
make dev-frontend  # 启动前端
make dev-backend   # 启动后端
make build         # 构建所有
make lint          # 检查所有代码

# 部署
cd deploy/scripts
sudo ./install.sh          # 安装
sudo ./health-check.sh     # 健康检查
sudo ./backup.sh           # 备份
```

---

## 文档索引

| 文档 | 位置 | 说明 |
|------|------|------|
| 需求文档 | `docs/01_需求文档.md` | 功能需求 |
| 实施方案 | `docs/02_详细实施方案.md` | 技术架构 |
| 需求规格 | `docs/03_需求规格说明书.md` | 详细规格 |
| 任务总览 | `agents/tasks/agent_overview.md` | 所有Agent任务 |
| 协作规范 | `agents/outputs/COLLABORATION_GUIDE.md` | 协作规则 |
| 冲突预警 | `agents/outputs/infra-scaffold/CONFLICT_WARNING.md` | 资源占用情况 |
| 本交付文档 | `agents/outputs/infra-scaffold/DELIVERY.md` | 详细交付说明 |

---

## 问题反馈

遇到问题？按优先级寻求帮助:
1. 查阅本文档和交付文档
2. 查阅协作规范和冲突预警表
3. 联系相关Agent协商
4. 上报 Architect Agent 或人类监督者

---

**祝开发顺利！** 🚀

*最后更新: 2026-02-22*
