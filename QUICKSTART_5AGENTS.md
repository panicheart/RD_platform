# RDP 项目 5-Agent Team 快速启动指南

基于 OpenCode + MCP 的精简 Agent 团队方案，5个 Agent 同步开展 Phase 1 开发。

## 🎯 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                    OpenCode Client (你)                      │
│                    Session: Leader                           │
└─────────────────────────┬───────────────────────────────────┘
                          │
    ┌─────────────────────┼─────────────────────┐
    ▼                     ▼                     ▼
┌──────────┐      ┌──────────────┐      ┌──────────┐
│ PM-Agent │◄────►│ Architect-A  │◄────►│DevOps-A  │
│(项目协调)│      │  (架构设计)   │      │(运维部署)│
└────┬─────┘      └──────┬───────┘      └────┬─────┘
     │                   │                   │
     └───────────────────┼───────────────────┘
                         │
        ┌────────────────┴────────────────┐
        ▼                                  ▼
┌───────────────────┐           ┌───────────────────┐
│  Backend-Agent    │◄─────────►│  Frontend-Agent   │
│  (Go后端开发)      │   API对接  │  (React前端开发)   │
│  • 用户管理API     │           │  • 门户界面        │
│  • 项目管理API     │           │  • 用户管理界面    │
│  • 认证授权       │           │  • 项目管理界面    │
└───────────────────┘           └───────────────────┘
```

## 🚀 快速启动 (5分钟)

### 步骤1: 初始化任务

```bash
# 在当前终端初始化Phase 1任务
make agent-team-init
```

输出示例：
```
🚀 初始化5-Agent Team任务...
✅ 数据库初始化完成: agents/data/5agent_tasks.db
✅ 任务已添加: P1-A1 -> Architect-Agent
✅ 任务已添加: P1-B1 -> Backend-Agent
...
✅ Phase 1 任务初始化完成，共 9 个任务
```

### 步骤2: 查看任务看板

```bash
make agent-team-status
```

### 步骤3: 启动5个Agent

打开 **5个终端窗口**，分别执行：

```bash
# 终端1: PM-Agent (项目经理)
make agent-pm

# 终端2: Architect-Agent (架构师)
make agent-architect

# 终端3: Backend-Agent (后端开发)
make agent-backend

# 终端4: Frontend-Agent (前端开发)
make agent-frontend

# 终端5: DevOps-Agent (运维部署)
make agent-devops
```

## 📋 各Agent启动指令

启动后，将以下内容粘贴到对应Agent会话中：

### PM-Agent (终端1)

```
你是 RDP项目的 PM-Agent（项目经理Agent）。

## 当前Phase: Phase 1 - 基础骨架

## 立即执行
1. 运行任务看板: python3 agents/5-agent-team/coordinator.py status
2. 查看所有任务: python3 agents/5-agent-team/coordinator.py list

## 你的职责
1. 任务分配: 确保各Agent收到正确任务
2. 进度跟踪: 每30分钟询问一次进度
3. 依赖协调: 当Backend-Agent完成API后，立即通知Frontend-Agent
4. 冲突仲裁: 当Agent间有分歧时做出决策

## 当前优先级
P0任务:
- P1-A1: Architect-Agent - 数据库Schema设计
- P1-A2: Architect-Agent - API接口规范
- P1-B1: Backend-Agent - 用户管理API
- P1-F1: Frontend-Agent - 门户界面

## 输出文件
- agents/outputs/pm/task_assignments.md
- agents/outputs/pm/progress_reports.md
```

### Architect-Agent (终端2)

```
你是 RDP项目的 Architect-Agent（架构师Agent）。

## 当前任务 (P0)
1. P1-A1: 数据库Schema设计
   - 设计users表(用户管理)
   - 设计projects表(项目管理)
   - 设计activities表(活动跟踪)
   - 输出: database/migrations/001_init_schema.sql

2. P1-A2: API接口规范
   - 定义用户管理API (/api/v1/users)
   - 定义项目管理API (/api/v1/projects)
   - 输出: services/api/docs/api_spec.md

## 技术约束
- Go 1.22+, Gin 1.9+, GORM
- React 18.x, TypeScript 5.x, Ant Design 5.x
- PostgreSQL 16.x

## 下一步
完成后通知PM-Agent，并协助Backend-Agent理解设计
```

### Backend-Agent (终端3)

```
你是 RDP项目的 Backend-Agent（后端开发Agent）。

## 当前任务 (P0)
1. P1-B1: 用户管理API
   - 用户CRUD (GET/POST/PUT/DELETE /api/v1/users)
   - JWT认证 (/api/v1/auth/login)
   - RBAC权限控制
   - 输出: services/api/handlers/user.go

2. P1-B2: 项目管理API
   - 项目CRUD (/api/v1/projects)
   - 项目成员管理
   - 输出: services/api/handlers/project.go

## 项目结构
services/api/
├── handlers/     # HTTP处理器
├── services/     # 业务逻辑
├── models/       # GORM模型
└── middleware/   # 中间件

## 依赖
等待Architect-Agent完成P1-A1和P1-A2后开始开发

## 完成后
通知PM-Agent和Frontend-Agent
```

### Frontend-Agent (终端4)

```
你是 RDP项目的 Frontend-Agent（前端开发Agent）。

## 当前任务 (P0)
1. P1-F1: 门户界面
   - 部门首页 (公告、荣誉展示)
   - 个人工作台 (待办、项目列表)
   - 消息通知中心
   - 输出: apps/web/src/pages/portal/

2. P1-F2: 用户管理界面
   - 登录/注册页面
   - 用户列表页面
   - 个人Profile页面
   - 输出: apps/web/src/pages/users/

3. P1-F3: 项目管理界面
   - 项目列表页面
   - 项目创建向导 (5步)
   - 项目详情页面
   - 输出: apps/web/src/pages/projects/

## 技术栈
- React 18.x, TypeScript 5.x, Vite 5.x
- Ant Design 5.x, Zustand

## 依赖
- P1-F1 可立即开始
- P1-F2 依赖Backend-Agent的P1-B1
- P1-F3 依赖Backend-Agent的P1-B2
```

### DevOps-Agent (终端5)

```
你是 RDP项目的 DevOps-Agent（运维部署Agent）。

## 当前任务
1. P1-D1: 数据库初始化脚本 (P0)
   - database/init.sql
   - deploy/scripts/init-db.sh
   - 依赖Architect-Agent的P1-A1

2. P1-D2: systemd服务配置 (P1)
   - deploy/systemd/rdp-api.service
   - deploy/systemd/rdp-casdoor.service

3. P1-D3: 部署脚本 (P1)
   - deploy/scripts/install.sh
   - deploy/scripts/backup.sh
   - 依赖所有开发完成

## 技术栈
- PostgreSQL 16.x
- systemd
- Nginx 1.25+
- Shell脚本
```

## 📊 任务协调

### 查看任务状态

在任意终端执行：

```bash
# 查看所有任务
python3 agents/5-agent-team/coordinator.py list

# 查看指定Agent任务
python3 agents/5-agent-team/coordinator.py list Backend-Agent

# 更新任务状态
python3 agents/5-agent-team/coordinator.py update P1-B1 completed "API开发完成，测试通过"
```

### Agent间通信

Agent通过共享数据库通信：

1. **Backend-Agent完成API后**：
   - 更新任务状态: `coordinator.py update P1-B1 completed`
   - PM-Agent会收到通知，转发给Frontend-Agent

2. **遇到技术问题**：
   - Backend-Agent向Architect-Agent提问
   - Architect-Agent协助解决

3. **代码审查**：
   - 各Agent完成任务后，代码提交给Architect-Agent审查

## ⏱️ 开发时序

```
时间轴 ──────────────────────────────────────────────────────►

Architect-Agent
  ├─ P1-A1 [数据库设计]      ████
  └─ P1-A2 [API规范]              ████

Backend-Agent
  ├─ P1-B1 [用户API]                  ████████
  └─ P1-B2 [项目API]                          ████████

Frontend-Agent
  ├─ P1-F1 [门户]    ████████
  ├─ P1-F2 [用户界面]                     ██████
  └─ P1-F3 [项目界面]                             ██████

DevOps-Agent
  ├─ P1-D1 [DB脚本]      ██
  ├─ P1-D2 [systemd]         ████
  └─ P1-D3 [部署脚本]                                     ████
```

## 🎯 Phase 1 交付物

完成以下9个任务：

| 任务 | Agent | 交付物 |
|------|-------|--------|
| P1-A1 | Architect | database/migrations/001_init_schema.sql |
| P1-A2 | Architect | services/api/docs/api_spec.md |
| P1-B1 | Backend | services/api/handlers/user.go + tests |
| P1-B2 | Backend | services/api/handlers/project.go + tests |
| P1-F1 | Frontend | apps/web/src/pages/portal/ |
| P1-F2 | Frontend | apps/web/src/pages/users/ |
| P1-F3 | Frontend | apps/web/src/pages/projects/ |
| P1-D1 | DevOps | database/init.sql |
| P1-D3 | DevOps | deploy/scripts/install.sh |

## 🆘 故障排除

### 问题1: Agent无法启动

```bash
# 检查OpenCode是否安装
which opencode

# 检查工作目录
cd /Users/tancong/Code/RD_platform
pwd
```

### 问题2: 任务状态不更新

```bash
# 检查数据库
ls -la agents/data/5agent_tasks.db

# 重新初始化
rm agents/data/5agent_tasks.db
make agent-team-init
```

### 问题3: Agent间不同步

```bash
# 手动查看消息
python3 -c "
import sqlite3
conn = sqlite3.connect('agents/data/5agent_tasks.db')
cursor = conn.cursor()
cursor.execute('SELECT * FROM messages ORDER BY created_at DESC LIMIT 10')
for row in cursor.fetchall():
    print(row)
"
```

## 📚 相关文档

- [5-Agent详细说明](agents/5-agent-team/README.md)
- [MCP配置](.opencode/mcp.json)
- [需求文档](docs/01_需求文档.md)
- [实施方案](docs/02_详细实施方案.md)

## 🎉 成功标准

Phase 1 成功标志：
- ✅ 数据库Schema设计完成
- ✅ API接口规范定义完成
- ✅ 用户管理API可调用
- ✅ 门户界面可访问
- ✅ 一键部署脚本可运行

祝开发顺利！🚀
