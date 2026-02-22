# 本地开发部署指南

## 快速启动

### 方式1: Docker 一键启动 (推荐)

```bash
# 1. 进入部署目录
cd deploy/docker

# 2. 启动所有服务
docker-compose -f docker-compose.dev.yml up -d

# 3. 等待服务启动 (约30秒)

# 4. 访问应用
# 前端: http://localhost:3000
# 后端 API: http://localhost:8080
# API 文档: http://localhost:8080/swagger/index.html
```

### 方式2: 本地开发模式

#### 前提条件
- Node.js 18+
- Go 1.22+
- PostgreSQL 16

#### 启动步骤

**1. 启动数据库**
```bash
# macOS
brew install postgresql@16
brew services start postgresql@16

# 创建数据库
createdb rdp_db
```

**2. 启动后端**
```bash
cd services/api
cp .env.example .env
# 编辑 .env 配置数据库连接

go mod download
go run main.go
```

**3. 启动前端**
```bash
cd apps/web
npm install
npm run dev

# 访问 http://localhost:5173
```

---

## 默认账号

| 账号 | 密码 | 角色 |
|------|------|------|
| admin | Admin@123 | 管理员 |
| user | User@123 | 普通用户 |

---

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端 | 3000 | React + Vite |
| 后端 API | 8080 | Go + Gin |
| 数据库 | 5432 | PostgreSQL |

---

## 常用命令

```bash
# 查看日志
docker-compose -f docker-compose.dev.yml logs -f

# 停止服务
docker-compose -f docker-compose.dev.yml down

# 重置数据库
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d

# 进入数据库容器
docker exec -it rdp-postgres psql -U rdp -d rdp_db
```

---

## 开发工作流

### 修改后端代码
- 文件位于 `services/api/`
- 修改后自动热重载 (air)

### 修改前端代码
- 文件位于 `apps/web/`
- 修改后自动热重载 (Vite HMR)

---

## 故障排除

### 端口被占用
```bash
# 检查端口占用
lsof -i :3000
lsof -i :8080
lsof -i :5432

# 杀掉进程
kill -9 <PID>
```

### 数据库连接失败
```bash
# 检查 PostgreSQL 状态
docker ps | grep rdp-postgres

# 查看数据库日志
docker logs rdp-postgres
```

### 前端构建失败
```bash
cd apps/web
rm -rf node_modules package-lock.json
npm install
```

---

## 生产部署

生产环境请使用：
```bash
docker-compose -f docker-compose.prod.yml up -d
```

详见: `deploy/docker/README.md`
