# ForumAgent 启动指令 - Phase 3 技术论坛模块

## 任务概述

你是 ForumAgent，负责开发技术论坛相关模块。

## 负责任务

| 任务ID | 任务名称 | 优先级 |
|--------|----------|--------|
| P3-T9 | 论坛基础功能 | P1 |
| P3-T10 | 帖子发布与回复 | P1 |
| P3-T11 | 搜索与标签 | P1 |
| P3-T12 | 知识库关联 | P2 |

## 任务详情

### P3-T9: 论坛基础功能
**输出**:
- `services/api/models/forum.go` - 论坛数据模型
- `services/api/services/forum.go` - 论坛服务
- `services/api/handlers/forum.go` - HTTP处理器
- `database/migrations/015_forum.sql` - 数据库迁移
- `apps/web/src/pages/forum/BoardList.tsx` - 板块列表页
- `apps/web/src/pages/forum/PostList.tsx` - 帖子列表页

**验收**: 板块CRUD正常、帖子CRUD正常

### P3-T10: 帖子发布与回复
**输出**:
- `services/api/models/reply.go` - 回复数据模型
- `services/api/services/reply.go` - 回复服务
- `services/api/handlers/reply.go` - 回复处理器
- `database/migrations/015_replies.sql` - 数据库迁移
- `apps/web/src/components/forum/PostEditor.tsx` - 发帖编辑器
- `apps/web/src/components/forum/ReplyThread.tsx` - 回复组件

**验收**: Markdown支持、楼中楼回复、@通知、最佳答案标记

### P3-T11: 搜索与标签
**输出**:
- `services/api/services/forum_search.go` - 论坛搜索服务
- `services/api/services/forum_tag.go` - 论坛标签服务
- `apps/web/src/pages/forum/ForumSearch.tsx` - 搜索页面

**验收**: 论坛全文搜索、标签分类、标签筛选

### P3-T12: 知识库关联
**输出**:
- `services/api/services/knowledge_link.go` - 知识关联服务
- `apps/web/src/components/forum/ArchiveToKB.tsx` - 归档组件

**验收**: 帖子归档到知识库、知识推荐

## 技术约束

1. **API风格**:
   - 板块: `/api/v1/forum/boards`
   - 帖子: `/api/v1/forum/posts`
   - 回复: `/api/v1/forum/replies`
   - 响应: `{"code": int, "message": string, "data": object}`

2. **前端规范**:
   - 使用 Ant Design 组件
   - 复用现有布局
   - Markdown编辑器 (使用现有或引入库)

3. **数据库**:
   - 使用 PostgreSQL
   - ID生成: ULID
   - 时间: ISO 8601 UTC

4. **通知**:
   - 使用现有通知服务 (`services/api/services/notification.go`)
   - @提及触发通知

## 代码参考

- 后端参考: `services/api/handlers/user.go`, `services/api/services/project.go`
- 前端参考: `apps/web/src/pages/users/UsersPage.tsx`, `apps/web/src/pages/projects/ProjectsPage.tsx`
- 数据库迁移: `database/migrations/`

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
