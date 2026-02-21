# RDP 项目 5-Agent 同步工作方案

基于 OpenCode + MCP 架构的精简 Agent 团队设计。

## 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    OpenCode Client                          │
│                    (主控界面)                                │
│              Session: Leader (Claude Opus)                  │
└─────────────────────────┬───────────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          ▼               ▼               ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  PM-Agent       │ │ Architect-Agent │ │ DevOps-Agent    │
│  (项目经理)      │ │  (架构师)        │ │ (运维部署)      │
│  Claude Sonnet  │ │ Claude Sonnet   │ │ Claude Sonnet  │
│  任务协调        │ │ 技术架构         │ │ 基础设施        │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
         └───────────────────┼───────────────────┘
                             │
          ┌──────────────────┴──────────────────┐
          ▼                                     ▼
┌─────────────────────┐              ┌─────────────────────┐
│  Backend-Agent      │              │  Frontend-Agent     │
│  (后端开发)          │              │  (前端开发)          │
│  Claude Sonnet      │              │  Claude Sonnet      │
│  Go + Gin + PG      │              │  React + TS + Vite  │
└─────────────────────┘              └─────────────────────┘
```

## 5个Agent职责

| Agent | 职责 | 核心任务 | 技术栈 |
|-------|------|----------|--------|
| **PM-Agent** | 项目经理 | 任务分配、进度跟踪、依赖协调、冲突仲裁 | - |
| **Architect-Agent** | 架构师 | 接口设计、数据模型、代码规范、技术选型 | Go/React |
| **Backend-Agent** | 后端开发 | API开发、数据库、业务逻辑、单元测试 | Go + Gin + GORM |
| **Frontend-Agent** | 前端开发 | UI组件、页面开发、状态管理、样式 | React + TS + Vite |
| **DevOps-Agent** | 运维部署 | 数据库、部署脚本、CI/CD、监控 | Shell + systemd |

## 启动命令

### 方法1: 使用启动脚本 (推荐)

```bash
# 启动所有5个Agent (在5个不同终端执行)
make agent-team-start

# 或分别启动
./agents/5-agent-team/start-pm.sh
./agents/5-agent-team/start-architect.sh
./agents/5-agent-team/start-backend.sh
./agents/5-agent-team/start-frontend.sh
./agents/5-agent-team/start-devops.sh
```

### 方法2: 手动启动

```bash
# 终端1: PM-Agent
opencode --session rdp-pm --model claude-sonnet

# 终端2: Architect-Agent
opencode --session rdp-architect --model claude-sonnet

# 终端3: Backend-Agent
opencode --session rdp-backend --model claude-sonnet

# 终端4: Frontend-Agent
opencode --session rdp-frontend --model claude-sonnet

# 终端5: DevOps-Agent
opencode --session rdp-devops --model claude-sonnet
```

## 任务分配策略

### Phase 1: 基础骨架

| 任务 | 负责Agent | 依赖 | 预计时间 |
|------|-----------|------|----------|
| 数据库设计 | Architect-Agent | - | 2h |
| 数据库初始化脚本 | DevOps-Agent | Architect-Agent | 1h |
| 用户管理API | Backend-Agent | Architect-Agent | 4h |
| 项目管理API | Backend-Agent | 用户管理API | 4h |
| 认证集成 | Backend-Agent | 用户管理API | 2h |
| 门户界面 | Frontend-Agent | - | 4h |
| 用户管理界面 | Frontend-Agent | 用户管理API | 3h |
| 项目管理界面 | Frontend-Agent | 项目管理API | 3h |
| 部署脚本 | DevOps-Agent | 所有开发完成 | 2h |

### 协调规则

1. **每日站会**: 各Agent汇报进度和阻塞问题
2. **API优先**: Backend-Agent优先定义接口，Frontend-Agent并行开发
3. **代码审查**: Architect-Agent审查所有代码
4. **冲突仲裁**: PM-Agent协调Agent间分歧

## 通信机制

Agent间通过 MCP Task Coordinator 通信：

```
Task Status Updates:
  Backend-Agent ──► 完成用户API ──► PM-Agent通知Frontend-Agent

Code Review:
  Frontend-Agent ──► 提交代码 ──► Architect-Agent审查

Blocker Escalation:
  Backend-Agent ──► 遇到技术问题 ──► Architect-Agent协助 ──► PM-Agent跟进
```
