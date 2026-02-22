# Phase 3 知识智能 - 5-Agent 开发团队

基于 OpenCode + MCP 架构的 Phase 3 开发团队设计。

## 团队架构

```
┌─────────────────────────────────────────────────────────────┐
│                    OpenCode Client                          │
│                    (主控界面)                                │
│              Session: Leader (Claude Opus)                 │
└─────────────────────────┬───────────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          ▼               ▼               ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  PM-Agent       │ │ KnowledgeAgent  │ │ SearchAgent     │
│  (项目经理)      │ │  (知识库)        │ │  (搜索服务)     │
│  Claude Sonnet  │ │ Claude Sonnet  │ │ Claude Sonnet  │
│  任务协调        │ │ 知识分类/Obsidian│ │ MeiliSearch    │
└────────┬────────┘ └─────────────────┘ └─────────────────┘
         │
         ▼
┌─────────────────┐ ┌─────────────────┐
│  ForumAgent     │ │ Reviewer-Agent  │
│  (技术论坛)      │ │ (代码审查)      │
│  Claude Sonnet  │ │ Claude Sonnet  │
│  板块/帖子/回复  │ │ 质量/测试       │
└─────────────────┘ └─────────────────┘
```

## Agent 职责

| Agent | 职责 | 核心任务 | 技术栈 |
|-------|------|----------|--------|
| **PM-Agent** | 项目经理 | 任务分配、进度跟踪、依赖协调、集成测试 | - |
| **KnowledgeAgent** | 知识库开发 | 知识分类、Obsidian同步、Zotero集成、Markdown渲染、标签系统 | Go + React |
| **SearchAgent** | 搜索服务开发 | MeiliSearch集成、搜索API、中文分词、跨模块索引 | Go + MeiliSearch |
| **ForumAgent** | 论坛开发 | 板块管理、帖子CRUD、回复、@通知、搜索 | Go + React |
| **Reviewer-Agent** | 质量审查 | 代码审查、单元测试、集成测试验证、问题修复 | - |

## Phase 3 任务分配

| 模块 | Agent | 任务ID | 优先级 |
|------|-------|--------|--------|
| 知识分类管理 | KnowledgeAgent | P3-T1 | P0 |
| Obsidian同步 | KnowledgeAgent | P3-T2 | P0 |
| Zotero集成 | KnowledgeAgent | P3-T3 | P0 |
| Markdown渲染 | KnowledgeAgent | P3-T4 | P0 |
| 标签系统 | KnowledgeAgent | P3-T5 | P1 |
| MeiliSearch集成 | SearchAgent | P3-T6 | P0 |
| 搜索API与高亮 | SearchAgent | P3-T7 | P0 |
| 跨模块索引 | SearchAgent | P3-T8 | P0 |
| 论坛基础功能 | ForumAgent | P3-T9 | P1 |
| 帖子发布与回复 | ForumAgent | P3-T10 | P1 |
| 搜索与标签 | ForumAgent | P3-T11 | P1 |
| 知识库关联 | ForumAgent | P3-T12 | P2 |

## 开发约束

### 技术约束
- 前端: React + TypeScript + Ant Design (复用现有框架)
- 后端: Go + Gin (复用现有服务)
- 数据库: PostgreSQL 16.x (复用现有实例)
- 搜索: MeiliSearch 1.x (新增)

### 编码规范 (强制)
| 规范项 | 要求 |
|--------|------|
| 代码注释 | 英文 |
| 变量命名 | 英文 |
| UI文案 | 中文 |
| API路径 | `/api/v1/{module}` 小写+连字符 |
| 错误响应 | `{"code": int, "message": string, "data": null}` |
| 时间格式 | ISO 8601 UTC |
| ID生成 | ULID 或雪花算法 |

### 目录约束
```
services/api/           # 后端Go代码
├── handlers/          # HTTP处理器
├── services/          # 业务逻辑
├── models/           # 数据模型
├── clients/           # 外部客户端
└── indexers/          # 索引器

apps/web/src/          # 前端React代码
├── pages/            # 页面组件
│   ├── knowledge/    # 知识库页面
│   ├── search/       # 搜索页面
│   └── forum/        # 论坛页面
└── components/       # 公共组件
```

## 开发流程

### 1. 并行开发阶段
- KnowledgeAgent、SearchAgent、ForumAgent 完全并行开发
- 各自独立完成模块的后端API和前端页面

### 2. 集成测试阶段
- Reviewer-Agent 验证所有模块
- PM-Agent 协调跨模块集成

### 3. 验收阶段
- Reviewer-Agent 执行集成测试清单
- 人类监督者 (用户) 验收确认

## 验证清单

| 序号 | 测试项 | 验收条件 |
|------|--------|----------|
| IT-38 | 知识分类管理 | 3级分类正常 |
| IT-39 | Obsidian同步 | 双向同步正常 |
| IT-40 | Zotero集成 | 文献读取正常 |
| IT-41 | Markdown渲染 | Wiki链接正确 |
| IT-42 | MeiliSearch | 搜索响应<500ms |
| IT-43 | 跨模块搜索 | 多范围搜索正常 |
| IT-44 | 论坛板块 | 增删改查正常 |
| IT-45 | 帖子发布 | Markdown正常 |
| IT-46 | 回复功能 | 楼中楼正常 |
| IT-47 | @通知 | 通知发送正确 |
| IT-48 | 搜索高亮 | 关键词高亮 |
| IT-49 | 搜索建议 | 自动补全正常 |

## 重要提醒

1. **开发完成后不要提交** - 等人类监督者确认后再提交
2. **先验证再报告** - 每个功能开发完成需自验证
3. **遵循现有模式** - 复用 Phase 1/2 的代码结构

---

*Phase 3 Team - 知识智能开发*
*启动日期: 2026-02-22*
