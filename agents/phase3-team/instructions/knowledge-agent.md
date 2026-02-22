# KnowledgeAgent 启动指令 - Phase 3 知识库模块

## 任务概述

你是 KnowledgeAgent，负责开发知识库相关模块。

## 负责任务

| 任务ID | 任务名称 | 优先级 |
|--------|----------|--------|
| P3-T1 | 知识分类管理 | P0 |
| P3-T2 | Obsidian Vault同步 | P0 |
| P3-T3 | Zotero集成 | P0 |
| P3-T4 | Markdown渲染 | P0 |
| P3-T5 | 标签系统与审核 | P1 |

## 任务详情

### P3-T1: 知识分类管理
**输出**:
- `services/api/models/knowledge.go` - 知识数据模型
- `services/api/services/knowledge.go` - 知识服务
- `services/api/handlers/knowledge.go` - HTTP处理器
- `database/migrations/014_knowledge.sql` - 数据库迁移
- `apps/web/src/pages/knowledge/KnowledgeList.tsx` - 知识列表页
- `apps/web/src/components/knowledge/CategoryTree.tsx` - 分类树组件

**验收**: 3级分类支持、CRUD正常

### P3-T2: Obsidian Vault同步
**输出**:
- `services/api/services/obsidian.go` - Obsidian服务
- `services/api/handlers/obsidian_sync.go` - 文件同步处理器
- `config/obsidian.yaml` - 路径映射配置

**验收**: 双向同步正常

### P3-T3: Zotero集成
**输出**:
- `services/api/services/zotero.go` - Zotero服务
- `services/api/handlers/zotero.go` - Zotero处理器
- `apps/web/src/components/knowledge/ZoteroLibrary.tsx` - Zotero组件

**验收**: 文献读取正常、PDF预览正常

### P3-T4: Markdown渲染
**输出**:
- `services/api/services/markdown.go` - Markdown服务
- `apps/web/src/components/knowledge/MarkdownRender.tsx` - 渲染组件
- `apps/web/src/utils/wiki_link.ts` - Wiki链接解析

**验收**: Wiki链接跳转、Callout渲染、公式渲染正常

### P3-T5: 标签系统与审核
**输出**:
- `services/api/services/tag.go` - 标签服务
- `services/api/services/review.go` - 审核服务
- `apps/web/src/components/knowledge/TagInput.tsx` - 标签组件

**验收**: 多标签支持、审核流程正常

## 技术约束

1. **API风格**:
   - 路径: `/api/v1/knowledge`, `/api/v1/knowledge/categories`, `/api/v1/obsidian`, `/api/v1/zotero`
   - 方法: GET/POST/PUT/DELETE
   - 响应: `{"code": int, "message": string, "data": object}`

2. **前端规范**:
   - 使用 Ant Design 组件
   - 复用现有布局 (`MainLayout.tsx`)
   - 使用现有认证 (`useAuth.tsx`)

3. **数据库**:
   - 使用 PostgreSQL
   - ID生成: ULID
   - 时间: ISO 8601 UTC

## 代码参考

- 后端参考: `services/api/handlers/user.go`, `services/api/services/project.go`
- 前端参考: `apps/web/src/pages/users/UsersPage.tsx`, `apps/web/src/pages/projects/ProjectsPage.tsx`
- 数据库迁移: `database/migrations/` 目录

## 开发流程

1. **先读参考** - 阅读现有代码结构
2. **创建后端** - models → services → handlers → routes
3. **创建前端** - pages → components
4. **自验证** - 测试API和页面功能

## 重要提醒

1. **开发完成后不要提交** - 等人类监督者确认后再提交
2. **遵循现有模式** - 复用 Phase 1/2 的代码结构
3. **代码注释** - 英文注释

---

## 参考文档

- Phase 3 任务卡片: `agents/tasks/phase3_tasks.md`
- 项目编码规范: `README.md`
