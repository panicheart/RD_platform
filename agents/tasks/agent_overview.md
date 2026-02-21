# Agent 任务总览表与依赖关系

> **文档用途**: 本文档汇总所有Agent的任务概览，包含任务ID、输入、输出、检查者和依赖关系，供各Agent协同开发使用。

---

## 1. Agent 任务总览表

### Phase 1: 基础骨架 (20个任务)

| 任务ID | Agent | 任务名称 | 优先级 | 输入 | 输出 | 检查者 |
|--------|-------|----------|--------|------|------|--------|
| P1-T1 | PortalAgent | 部门门户首页 | P0 | 技术栈规范 | 门户首页组件 | Reviewer |
| P1-T2 | PortalAgent | 个人工作台 | P0 | 用户/项目数据 | 工作台页面 | Reviewer |
| P1-T3 | PortalAgent | 消息通知中心 | P0 | 通知API | 通知组件 | Reviewer |
| P1-T4 | PortalAgent | 全局搜索UI | P1 | 搜索API(后续) | 搜索组件 | Reviewer |
| P1-T5 | UserAgent | 用户认证与CRUD | P0 | 需求规格 | 用户API/模型 | Reviewer |
| P1-T6 | UserAgent | RBAC权限模型 | P0 | 权限矩阵 | 权限服务 | Reviewer |
| P1-T7 | UserAgent | 组织架构管理 | P1 | 用户数据 | 组织架构API | Reviewer |
| P1-T8 | UserAgent | 个人Profile页面 | P1 | 用户/项目数据 | Profile页面 | Reviewer |
| P1-T9 | ProjectAgent | 项目CRUD与看板 | P0 | 用户数据 | 项目API/看板 | Reviewer |
| P1-T10 | ProjectAgent | 五步创建向导 | P0 | 项目/用户数据 | 向导页面 | Reviewer |
| P1-T11 | ProjectAgent | 流程模板管理 | P0 | 需求规格 | 模板API | Reviewer |
| P1-T12 | ProjectAgent | 基础文件管理 | P0 | 项目数据 | 文件API | Reviewer |
| P1-T13 | SecurityAgent | 数据分级分类 | P0 | 需求规格 | 分级API | Reviewer |
| P1-T14 | SecurityAgent | 会话超时控制 | P0 | 需求规格 | 会话服务 | Reviewer |
| P1-T15 | SecurityAgent | 操作审计日志 | P0 | 需求规格 | 审计API | Reviewer |
| P1-T16 | SecurityAgent | 屏幕水印(基础) | P1 | 分级数据 | 水印组件 | Reviewer |
| P1-T17 | InfraAgent | 数据库初始化 | P0 | 各Agent模型 | SQL迁移 | Reviewer |
| P1-T18 | InfraAgent | systemd服务配置 | P0 | 服务定义 | .service文件 | Reviewer |
| P1-T19 | InfraAgent | Nginx反向代理 | P0 | 端口规划 | nginx.conf | Reviewer |
| P1-T20 | InfraAgent | 一键安装脚本 | P0 | 组件清单 | install.sh | Reviewer |

### Phase 2: 核心业务 (23个任务)

| 任务ID | Agent | 任务名称 | 优先级 | 输入 | 输出 | 检查者 |
|--------|-------|----------|--------|------|------|--------|
| P2-T1 | WorkflowAgent | 状态机引擎 | P0 | 模板数据 | 状态机服务 | Reviewer |
| P2-T2 | WorkflowAgent | 活动流转逻辑 | P0 | 状态机 | 活动API | Reviewer |
| P2-T3 | WorkflowAgent | DCP评审节点 | P0 | 活动数据 | 评审API | Reviewer |
| P2-T4 | ProjectAgent | Gitea API集成 | P0 | Gitea服务 | Git客户端 | Reviewer |
| P2-T5 | ProjectAgent | Git版本管理 | P0 | Gitea集成 | Git API | Reviewer |
| P2-T6 | ProjectAgent | 甘特图 | P0 | 活动数据 | 甘特图组件 | Reviewer |
| P2-T7 | DevAgent | 流程全景视图 | P0 | 流程数据 | 流程图组件 | Reviewer |
| P2-T8 | DevAgent | 本地软件联动 | P0 | 文件数据 | 协议服务 | Reviewer |
| P2-T9 | DevAgent | 活动执行面板 | P0 | 活动数据 | 执行页面 | Reviewer |
| P2-T10 | DevAgent | 评审反馈系统 | P0 | 用户数据 | 反馈API | Reviewer |
| P2-T11 | DevAgent | 变更管理 | P1 | 活动数据 | 变更API | Reviewer |
| P2-T12 | ShelfAgent | 产品浏览与筛选 | P0 | 项目数据 | 产品API | Reviewer |
| P2-T13 | ShelfAgent | 选用购物车 | P0 | 产品数据 | 购物车API | Reviewer |
| P2-T14 | ShelfAgent | 技术树可视化 | P0 | 技术数据 | 技术树组件 | Reviewer |
| P2-T15 | ShelfAgent | 版本管理与自动上架 | P1 | 产品/流程数据 | 版本服务 | Reviewer |
| P2-T16 | DesktopAgent | rdp协议注册 | P0 | 协议定义 | 桌面程序 | Reviewer |
| P2-T17 | DesktopAgent | 本地软件调用 | P0 | 映射配置 | 文件调用模块 | Reviewer |
| P2-T18 | DesktopAgent | Git自动提交 | P0 | 活动通知 | Git操作模块 | Reviewer |
| P2-T19 | DesktopAgent | 冲突检测与解决 | P1 | Git状态 | 冲突处理模块 | Reviewer |
| P2-T20 | QMAgent | 需求管理 | P1 | 项目数据 | 需求API | Reviewer |
| P2-T21 | QMAgent | 变更管理(ECR/ECO) | P1 | 需求数据 | 变更API | Reviewer |
| P2-T22 | QMAgent | 缺陷管理 | P1 | 项目数据 | 缺陷API | Reviewer |
| P2-T23 | QMAgent | 质量门禁 | P1 | DCP评审 | 门禁服务 | Reviewer |

### Phase 3: 知识智能 (12个任务)

| 任务ID | Agent | 任务名称 | 优先级 | 输入 | 输出 | 检查者 |
|--------|-------|----------|--------|------|------|--------|
| P3-T1 | KnowledgeAgent | 知识分类管理 | P0 | 需求规格 | 知识API | Reviewer |
| P3-T2 | KnowledgeAgent | Obsidian Vault同步 | P0 | 知识数据 | 同步服务 | Reviewer |
| P3-T3 | KnowledgeAgent | Zotero集成 | P0 | Zotero配置 | Zotero服务 | Reviewer |
| P3-T4 | KnowledgeAgent | Markdown渲染 | P0 | Markdown内容 | 渲染组件 | Reviewer |
| P3-T5 | KnowledgeAgent | 标签系统与审核 | P1 | 知识数据 | 标签服务 | Reviewer |
| P3-T6 | SearchAgent | MeiliSearch集成 | P0 | MeiliSearch | 搜索客户端 | Reviewer |
| P3-T7 | SearchAgent | 搜索API与高亮 | P0 | 搜索服务 | 搜索API | Reviewer |
| P3-T8 | SearchAgent | 跨模块索引 | P0 | 项目/知识/产品 | 索引器 | Reviewer |
| P3-T9 | ForumAgent | 论坛基础功能 | P1 | 需求规格 | 论坛API | Reviewer |
| P3-T10 | ForumAgent | 帖子发布与回复 | P1 | 用户数据 | 回复API | Reviewer |
| P3-T11 | ForumAgent | 搜索与标签 | P1 | 帖子数据 | 搜索服务 | Reviewer |
| P3-T12 | ForumAgent | 知识库关联 | P2 | 帖子/知识数据 | 关联服务 | Reviewer |

### Phase 4: 优化完善 (10个任务)

| 任务ID | Agent | 任务名称 | 优先级 | 输入 | 输出 | 检查者 |
|--------|-------|----------|--------|------|------|--------|
| P4-T1 | AnalyticsAgent | 数据仪表盘 | P1 | 项目/用户数据 | 分析API | Reviewer |
| P4-T2 | AnalyticsAgent | 项目统计报表 | P1 | 项目数据 | 统计服务 | Reviewer |
| P4-T3 | AnalyticsAgent | 人员绩效统计 | P1 | 用户/项目数据 | 绩效服务 | Reviewer |
| P4-T4 | AnalyticsAgent | 报表生成 | P1 | 统计数据 | 导出功能 | Reviewer |
| P4-T5 | MonitorAgent | 系统监控仪表盘 | P1 | 系统指标 | 监控API | Reviewer |
| P4-T6 | MonitorAgent | APM性能监控 | P1 | API性能 | APM服务 | Reviewer |
| P4-T7 | MonitorAgent | 日志集中管理 | P1 | 应用日志 | 日志服务 | Reviewer |
| P4-T8 | MonitorAgent | 告警机制 | P1 | 监控指标 | 告警服务 | Reviewer |
| P4-T9 | MonitorAgent | Prometheus集成 | P2 | Prometheus | 配置/面板 | Reviewer |
| P4-T10 | 全局 | 性能优化 | P1 | 性能数据 | 优化实施 | Reviewer |

---

## 2. Agent 依赖关系矩阵

### 2.1 Phase 1 依赖关系（完全并行）

```
PortalAgent ─────┐
                 ├── 无依赖，可同时启动
UserAgent   ────┤
                 │
ProjectAgent ───┤
                 │
SecurityAgent ──┤
                 │
InfraAgent  ────┘
```

### 2.2 Phase 2 依赖关系（分层并行）

#### Layer 1（并行）
```
WorkflowAgent ──┐
                ├── 依赖 Phase 1 完成
ProjectAgent ──┘   (Phase 1扩展)
```

#### Layer 2（依赖 Layer 1）
```
DevAgent ────────────────┐
                         ├── 依赖 WorkflowAgent + ProjectAgent
ShelfAgent ─────────────┤
                         │
DesktopAgent ───────────┤
                         │
QMAgent ────────────────┘
```

### 2.3 Phase 3 依赖关系（并行）

```
KnowledgeAgent ──┐
                 ├── 依赖 Phase 2 完成
SearchAgent  ────┤
                 │
ForumAgent  ─────┘
```

### 2.4 Phase 4 依赖关系（并行）

```
AnalyticsAgent ──┐
                 ├── 依赖 Phase 3 完成
MonitorAgent  ───┘
```

---

## 3. 数据流依赖关系

```
用户数据 (UserAgent)
    ↓
    ├─→ 项目管理 (ProjectAgent) ──→ 流程引擎 (WorkflowAgent) ──→ 项目开发 (DevAgent)
    │         ↓                           ↓                           ↓
    │      甘特图                    活动流转                   本地软件联动
    │         ↓                           ↓                           ↓
    │    Gitea集成                    DCP评审                   桌面程序
    │                                                   (DesktopAgent)
    │
    ├─→ 产品货架 (ShelfAgent)
    │         ↓
    │    技术货架
    │         ↓
    │    知识库 (KnowledgeAgent)
    │         ↓
    │    搜索服务 (SearchAgent)
    │
    ├─→ 论坛 (ForumAgent)
    │
    └─→ 安全合规 (SecurityAgent) ──→ 审计日志
                                            ↓
                                     监控 (MonitorAgent)
                                            ↓
                                     分析 (AnalyticsAgent)
```

---

## 4. 验收流程

### 4.1 单任务验收流程

```
Agent 完成开发
    ↓
    ↓ L1: 自审查 (代码规范、语法检查)
    ↓
    ↓ 提交交付物
    ↓
Reviewer Agent 审查 (L2)
    ↓
    ├─ 通过 → 交付完成
    └─ 需修改 → 打回修改 → 重新提交
```

### 4.2 Phase 验收流程

```
Phase 内所有任务完成
    ↓
    ↓ PM-Agent 集成测试 (L3)
    ↓
    ↓ 人类监督者验收 (L4)
    ↓
    ├─ 通过 → 进入下一Phase
    └─ 需修改 → 返回修改
```

---

## 5. 交付物格式规范

### 5.1 后端代码
| 类型 | 路径格式 | 示例 |
|------|----------|------|
| Handler | `services/api/handlers/{module}.go` | `services/api/handlers/user.go` |
| Service | `services/api/services/{module}.go` | `services/api/services/auth.go` |
| Model | `services/api/models/{module}.go` | `services/api/models/user.go` |
| Middleware | `services/api/middleware/{name}.go` | `services/api/middleware/auth.go` |
| Client | `services/api/clients/{service}.go` | `services/api/clients/gitea.go` |

### 5.2 前端代码
| 类型 | 路径格式 | 示例 |
|------|----------|------|
| Page | `apps/web/src/pages/{module}/{PageName}.tsx` | `apps/web/src/pages/portal/PortalPage.tsx` |
| Component | `apps/web/src/components/{category}/{ComponentName}.tsx` | `apps/web/src/components/portal/AnnouncementList.tsx` |
| Service | `apps/web/src/services/{module}.ts` | `apps/web/src/services/api.ts` |
| Store | `apps/web/src/stores/{module}.ts` | `apps/web/src/stores/user.ts` |
| Utils | `apps/web/src/utils/{name}.ts` | `apps/web/src/utils/watermark.ts` |

### 5.3 数据库
| 类型 | 路径格式 | 示例 |
|------|----------|------|
| Migration | `database/migrations/{序号}_{name}.sql` | `database/migrations/001_users.sql` |
| Seed | `database/seeds/{name}.sql` | `database/seeds/process_templates.sql` |
| Schema | `database/schema/{name}.sql` | `database/schema/enums.sql` |

### 5.4 部署配置
| 类型 | 路径格式 | 示例 |
|------|----------|------|
| systemd | `deploy/systemd/{service}.service` | `deploy/systemd/rdp-api.service` |
| Nginx | `deploy/nginx/{path}/{file}.conf` | `deploy/nginx/sites-available/rdp.conf` |
| Script | `deploy/scripts/{name}.sh` | `deploy/scripts/install.sh` |
| Config | `config/{name}.yaml` | `config/rdp-api.yaml` |

---

## 6. API 命名规范

### 6.1 路径规范
```
GET    /api/v1/users              # 列表
GET    /api/v1/users/:id           # 详情
POST   /api/v1/users               # 创建
PUT    /api/v1/users/:id           # 更新
DELETE /api/v1/users/:id           # 删除
```

### 6.2 响应格式
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### 6.3 错误码规范
```
200 - 成功
400 - 请求参数错误
401 - 未授权
403 - 禁止访问
404 - 资源不存在
500 - 服务器内部错误
```

---

## 7. 任务启动指南

### 7.1 Phase 1 启动命令

```bash
# 同时启动以下Agent：
task(agent="PortalAgent", ...)
task(agent="UserAgent", ...)
task(agent="ProjectAgent", ...)
task(agent="SecurityAgent", ...)
task(agent="InfraAgent", ...)
```

### 7.2 Phase 2 Layer 1 启动命令

```bash
# Phase 1 完成后，启动：
task(agent="WorkflowAgent", ...)
task(agent="ProjectAgent", ...)  # 扩展
```

### 7.3 Phase 2 Layer 2 启动命令

```bash
# Layer 1 完成后，启动：
task(agent="DevAgent", ...)
task(agent="ShelfAgent", ...)
task(agent="DesktopAgent", ...)
task(agent="QMAgent", ...)
```

### 7.4 Phase 3 启动命令

```bash
# Phase 2 完成后，启动：
task(agent="KnowledgeAgent", ...)
task(agent="SearchAgent", ...)
task(agent="ForumAgent", ...)
```

### 7.5 Phase 4 启动命令

```bash
# Phase 3 完成后，启动：
task(agent="AnalyticsAgent", ...)
task(agent="MonitorAgent", ...)
```

---

*文档结束 - 各Agent可根据此文档明确自己的任务范围和依赖关系*
