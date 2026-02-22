# SearchAgent 启动指令 - Phase 3 搜索服务模块

## 任务概述

你是 SearchAgent，负责开发搜索服务相关模块。

## 负责任务

| 任务ID | 任务名称 | 优先级 |
|--------|----------|--------|
| P3-T6 | MeiliSearch集成 | P0 |
| P3-T7 | 搜索API与高亮 | P0 |
| P3-T8 | 跨模块索引 | P0 |

## 任务详情

### P3-T6: MeiliSearch集成
**输出**:
- `services/api/clients/meilisearch.go` - MeiliSearch客户端
- `services/api/services/search.go` - 搜索服务
- `services/api/handlers/search.go` - 搜索处理器
- `config/meilisearch.yaml` - MeiliSearch配置

**验收**: 连接正常、中文分词支持(jieba)

### P3-T7: 搜索API与高亮
**输出**:
- 搜索API Handler (已在P3-T6创建)
- `services/api/services/suggest.go` - 建议服务
- `apps/web/src/components/search/GlobalSearch.tsx` - 全局搜索组件

**验收**: 关键词高亮、搜索建议、自动补全、搜索响应<500ms

### P3-T8: 跨模块索引
**输出**:
- `services/api/services/indexer.go` - 索引服务
- `services/api/indexers/project_indexer.go` - 项目索引器
- `services/api/indexers/knowledge_indexer.go` - 知识索引器
- `services/api/indexers/product_indexer.go` - 产品索引器

**验收**: 增量更新、跨模块统一搜索

## 技术约束

1. **MeiliSearch**:
   - 使用 meilisearch-go 客户端
   - 支持中文分词 (jieba)
   - 索引配置: 可搜索属性、过滤属性、排序属性

2. **API风格**:
   - 路径: `/api/v1/search`, `/api/v1/search/suggest`
   - 方法: GET/POST
   - 响应: `{"code": int, "message": string, "data": object}`

3. **前端规范**:
   - 使用 Ant Design 组件
   - 复用现有布局
   - 支持搜索建议下拉

4. **数据库**:
   - 索引同步: 增量更新机制
   - 定时任务或Webhook触发

## 代码参考

- MeiliSearch官方: https://github.com/meilisearch/meilisearch-go
- 后端参考: `services/api/services/project.go`
- 前端参考: `apps/web/src/components/`

## 开发流程

1. **先读参考** - 阅读现有代码结构和MeiliSearch文档
2. **创建客户端** - meilisearch-go集成
3. **创建服务** - search service, indexer service
4. **创建处理器** - search handler
5. **创建前端** - GlobalSearch组件
6. **自验证** - 测试搜索功能

## 重要提醒

1. **开发完成后不要提交** - 等人类监督者确认后再提交
2. **遵循现有模式** - 复用 Phase 1/2 的代码结构
3. **代码注释** - 英文注释
4. **搜索性能** - 确保响应时间<500ms

---

## 参考文档

- Phase 3 任务卡片: `agents/tasks/phase3_tasks.md`
- MeiliSearch Go客户端: https://github.com/meilisearch/meilisearch-go
- 项目编码规范: `README.md`
