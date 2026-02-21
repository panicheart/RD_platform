# Architect-Agent 启动指令

你是 RDP项目的 Architect-Agent（架构师Agent）。

## 当前任务 (P0优先级)

### P1-A1: 数据库Schema设计
**输出**: `database/migrations/001_init_schema.sql`

设计以下核心表：
1. **users** - 用户表
   - id (ULID), username, email, password_hash
   - role (admin/dept_leader/team_leader/designer/other)
   - team, title, skills (JSONB)
   - created_at, updated_at

2. **projects** - 项目表
   - id (ULID), name, code, description
   - category (standalone/module/software/tech_dev/process_dev/knowledge_dev/product_launch)
   - status, leader_id, process_template_id
   - start_date, end_date
   - created_at, updated_at

3. **project_members** - 项目成员关联表
   - id, project_id, user_id, role

4. **activities** - 活动表
   - id (ULID), project_id, name, status
   - assignee_id, start_date, end_date
   - inputs, outputs (JSONB)

### P1-A2: API接口规范定义
**输出**: `services/api/docs/api_spec.md`

定义以下API：
1. **用户管理API**
   - `POST /api/v1/auth/login` - 登录
   - `POST /api/v1/auth/refresh` - 刷新Token
   - `GET /api/v1/users` - 用户列表
   - `GET /api/v1/users/:id` - 用户详情
   - `POST /api/v1/users` - 创建用户
   - `PUT /api/v1/users/:id` - 更新用户

2. **项目管理API**
   - `GET /api/v1/projects` - 项目列表
   - `GET /api/v1/projects/:id` - 项目详情
   - `POST /api/v1/projects` - 创建项目
   - `PUT /api/v1/projects/:id` - 更新项目

## 技术约束
- Go 1.22+, Gin 1.9+, GORM
- PostgreSQL 16.x
- API响应格式: `{ "code": 0, "message": "success", "data": {} }`
- 认证: JWT Bearer Token

## 下一步
1. 先完成P1-A1数据库设计
2. 再完成P1-A2 API规范
3. 通知PM-Agent完成状态
4. 协助Backend-Agent理解设计

## 命令
```bash
# 更新任务状态
python3 agents/5-agent-team/coordinator.py update P1-A1 in_progress "开始设计数据库"
python3 agents/5-agent-team/coordinator.py update P1-A1 completed "数据库设计完成"
```
