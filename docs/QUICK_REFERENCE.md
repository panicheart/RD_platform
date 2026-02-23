# 开发环境快速参考

> **日期**: 2026-02-23  
> **版本**: 基于 Phase 1/2 稳定版

---

## 🚀 快速启动

### 1. 启动数据库

```bash
# 确保 Docker/Colima 已运行
colima status || colima start

# 启动 PostgreSQL
cd deploy/docker
docker-compose -f docker-compose.dev.yml up -d postgres
```

### 2. 启动后端

```bash
# 设置环境变量
export RDP_DB_USER=rdp
export RDP_DB_PASSWORD=rdp123
export RDP_DB_NAME=rdp_db
export RDP_JWT_SECRET=your-secret-key

# 运行后端
cd services/api
/tmp/rdp-api

# 或使用 go run
go run main.go
```

### 3. 启动前端

```bash
cd apps/web
npm run dev
```

---

## 🔧 常用命令

### 数据库

```bash
# 查看容器状态
docker ps --filter "name=rdp"

# 进入数据库
docker exec -it rdp-postgres psql -U rdp -d rdp_db

# 查看日志
docker logs rdp-postgres

# 停止数据库
docker-compose -f deploy/docker/docker-compose.dev.yml down
```

### 后端

```bash
# 构建后端
cd services/api
go build -o /tmp/rdp-api main.go

# 测试 API
curl http://localhost:8080/api/v1/health

# 查看后端进程
ps aux | grep rdp-api

# 停止后端
pkill rdp-api
```

### 前端

```bash
# 安装依赖
cd apps/web
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 运行测试
npm run test
```

---

## 📊 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端 | 5173 | Vite 开发服务器 |
| 后端 API | 8080 | Go Gin 服务 |
| 数据库 | 5432 | PostgreSQL 16 |

---

## 🔗 API 端点

### 健康检查
```bash
GET http://localhost:8080/api/v1/health
```

### 用户管理
```bash
# 列出用户
GET http://localhost:8080/api/v1/users

# 获取当前用户
GET http://localhost:8080/api/v1/users/me
```

### 项目管理
```bash
# 列出项目
GET http://localhost:8080/api/v1/projects

# 创建项目
POST http://localhost:8080/api/v1/projects

# 获取项目详情
GET http://localhost:8080/api/v1/projects/:id
```

### 产品货架
```bash
# 列出产品
GET http://localhost:8080/api/v1/products

# 获取 TRL 等级
GET http://localhost:8080/api/v1/trl-levels
```

---

## 🐛 故障排除

### 问题：后端无法连接数据库

**症状**: `Failed to connect to database: failed SASL auth`

**解决**:
```bash
# 检查环境变量
export RDP_DB_USER=rdp
export RDP_DB_PASSWORD=rdp123
export RDP_DB_NAME=rdp_db

# 检查数据库容器
docker ps | grep rdp-postgres

# 重置数据库
docker-compose -f deploy/docker/docker-compose.dev.yml down -v
docker-compose -f deploy/docker/docker-compose.dev.yml up -d postgres
```

### 问题：前端无法连接后端

**症状**: API 请求失败

**解决**:
- 检查后端是否运行在 8080 端口
- 检查 Vite 代理配置 `apps/web/vite.config.ts`

### 问题：Colima 无法启动

**症状**: `colima start` 失败

**解决**:
```bash
# 删除并重新创建 Colima 实例
colima delete
colima start --cpu 4 --memory 8 --disk 50
```

---

## 📁 项目结构

```
RD_platform/
├── apps/
│   └── web/                    # React + Vite 前端
├── services/
│   └── api/                    # Go + Gin 后端
│       ├── main.go             # 主程序入口
│       ├── handlers/           # HTTP 处理器
│       ├── services/           # 业务逻辑
│       ├── models/             # 数据模型
│       ├── routes/             # 路由配置
│       ├── middleware/         # 中间件
│       └── clients/            # 外部客户端
├── database/
│   ├── migrations/             # 迁移脚本
│   └── seeds/                  # 种子数据
├── deploy/
│   └── docker/                 # Docker 配置
├── docs/                       # 项目文档
└── config/                     # 配置文件
```

---

## 📝 环境变量

### 后端必需

```bash
# 数据库
RDP_DB_HOST=localhost
RDP_DB_PORT=5432
RDP_DB_USER=rdp
RDP_DB_PASSWORD=rdp123
RDP_DB_NAME=rdp_db
RDP_DB_SSLMODE=disable

# JWT
RDP_JWT_SECRET=your-secret-key-change-in-production
RDP_ACCESS_TOKEN_TTL=2h
RDP_REFRESH_TOKEN_TTL=168h

# 服务器
RDP_API_PORT=8080
RDP_ENV=development
```

### 前端必需

```bash
# API 地址
VITE_API_URL=http://localhost:8080
```

---

## 🔒 备份文件

Phase 3/4 未完成的功能代码已备份到 `/tmp/rdp-phase3-backup/`：

```bash
/tmp/rdp-phase3-backup/
├── analytics.go      # 数据分析服务
├── audit.go          # 审计中间件
├── forum.go          # 论坛服务
├── knowledge.go      # 知识库服务
├── markdown.go       # Markdown 处理
├── monitor.go        # 监控处理器
├── obsidian.go       # Obsidian 集成
├── search.go         # 搜索服务
└── zotero.go         # Zotero 集成
```

**注意**: 这些文件包含编译错误，需要修复后才能使用。

---

## 🎯 开发建议

1. **代码规范**: 所有代码注释使用英文，UI 文案使用中文
2. **测试**: 在提交前运行 `go test ./...` 和 `npm run test`
3. **文档**: 修改 API 时更新相关文档
4. **提交**: 遵循项目的 Git 提交规范

---

## 📚 相关文档

- [SETUP_COMPLETION.md](./SETUP_COMPLETION.md) - 安装与修复详情
- [README.md](../README.md) - 项目概述
- [AGENTS.md](../AGENTS.md) - Agent 开发指南

---

*本文档由 Sisyphus AI Agent 自动生成*
