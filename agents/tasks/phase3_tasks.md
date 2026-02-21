# Agent 任务卡片 - Phase 3 知识智能

---

## P3-T1: KnowledgeAgent - 知识分类管理

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T1 |
| **Agent** | KnowledgeAgent |
| **模块** | 知识库 |
| **优先级** | P0 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 需求规格 | `03_需求规格说明书.md` | SRS-KB-001 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 知识模型 | `services/api/models/knowledge.go` | .go |
| 知识Handler | `services/api/handlers/knowledge.go` | .go |
| 知识服务 | `services/api/services/knowledge.go` | .go |
| 知识表迁移 | `database/migrations/014_knowledge.sql` | .sql |
| 知识列表页 | `apps/web/src/pages/knowledge/KnowledgeList.tsx` | .tsx |
| 知识树组件 | `apps/web/src/components/knowledge/CategoryTree.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T1.1 | 创建知识数据模型 | 模型定义正确 |
| P3-T1.2 | 实现知识CRUD API | 接口正常 |
| P3-T1.3 | 实现分类树管理 | 3级深度支持 |
| P3-T1.4 | 实现前端分类树 | 展开/折叠正常 |
| P3-T1.5 | 实现知识列表 | 筛选/分页正常 |

---

## P3-T2: KnowledgeAgent - Obsidian Vault同步

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T2 |
| **Agent** | KnowledgeAgent |
| **模块** | 知识库 |
| **优先级** | P0 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 知识数据 | P3-T1 | 知识条目 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| Obsidian服务 | `services/api/services/obsidian.go` | .go |
| 文件同步Handler | `services/api/handlers/obsidian_sync.go` | .go |
| 路径映射配置 | `config/obsidian.yaml` | .yaml |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent (文件系统安全)

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T2.1 | 实现Vault路径映射 | 路径对应正确 |
| P3-T2.2 | 实现文件读取 | Markdown读取正常 |
| P3-T2.3 | 实现文件写入 | 保存到Vault |
| P3-T2.4 | 配置同步规则 | 规则正确 |

---

## P3-T3: KnowledgeAgent - Zotero集成

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T3 |
| **Agent** | KnowledgeAgent |
| **模块** | 知识库 |
| **优先级** | P0 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| Zotero配置 | 部署环境 | Zotero数据 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| Zotero服务 | `services/api/services/zotero.go` | .go |
| Zotero Handler | `services/api/handlers/zotero.go` | .go |
| Zotero组件 | `apps/web/src/components/knowledge/ZoteroLibrary.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T3.1 | 实现Zotero连接 | 数据库读取正常 |
| P3-T3.2 | 实现分类读取 | 分类列表正确 |
| P3-T3.3 | 实现文献读取 | 条目列表正确 |
| P3-T3.4 | 实现PDF预览 | 预览正常 |

---

## P3-T4: KnowledgeAgent - Markdown渲染

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T4 |
| **Agent** | KnowledgeAgent |
| **模块** | 知识库 |
| **优先级** | P0 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| Markdown内容 | P3-T2 | Obsidian文件 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| Markdown服务 | `services/api/services/markdown.go` | .go |
| 渲染组件 | `apps/web/src/components/knowledge/MarkdownRender.tsx` | .tsx |
| Wiki链接解析 | `apps/web/src/utils/wiki_link.ts` | .ts |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T4.1 | 实现基础Markdown渲染 | 语法正确 |
| P3-T4.2 | 实现Wiki链接解析 | [[]]链接跳转 |
| P3-T4.3 | 实现Callout渲染 | > [!NOTE]正确 |
| P3-T4.4 | 实现公式渲染 | LaTeX正确 |

---

## P3-T5: KnowledgeAgent - 标签系统与知识审核

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T5 |
| **Agent** | KnowledgeAgent |
| **模块** | 知识库 |
| **优先级** | P1 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 知识数据 | P3-T1 | 知识条目 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 标签服务 | `services/api/services/tag.go` | .go |
| 审核服务 | `services/api/services/review.go` | .go |
| 标签组件 | `apps/web/src/components/knowledge/TagInput.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T5.1 | 实现标签管理 | 多标签支持 |
| P2 | 实现标签自动关联 | 关联推荐 |
| P3-T5.3 | 实现审核流程 | 审核状态正确3-T5. |
| P3-T5.4 | 实现版本管理 | 历史版本保存 |

---

## P3-T6: SearchAgent - MeiliSearch集成

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T6 |
| **Agent** | SearchAgent |
| **模块** | 搜索服务 |
| **优先级** | P0 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| MeiliSearch | 部署环境 | 搜索引擎 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 搜索客户端 | `services/api/clients/meilisearch.go` | .go |
| 搜索服务 | `services/api/services/search.go` | .go |
| 搜索Handler | `services/api/handlers/search.go` | .go |
| MeiliSearch配置 | `config/meilisearch.yaml` | .yaml |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T6.1 | 实现MeiliSearch客户端 | 连接正常 |
| P3-T6.2 | 实现索引创建 | 索引配置正确 |
| P3-T6.3 | 实现中文分词 | jieba集成 |
| P3-T6.4 | 实现索引更新 | 增量更新 |

---

## P3-T7: SearchAgent - 搜索API与高亮

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T7 |
| **Agent** | SearchAgent |
| **模块** | 搜索服务 |
| **优先级** | P0 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 搜索服务 | P3-T6 | 搜索引擎 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 搜索API Handler | `services/api/handlers/search.go` | .go |
| 建议服务 | `services/api/services/suggest.go` | .go |
| 搜索组件升级 | `apps/web/src/components/search/GlobalSearch.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T7.1 | 实现搜索API | 多范围搜索 |
| P3-T7.2 | 实现关键词高亮 | 命中词高亮 |
| P3-T7.3 | 实现搜索建议 | 自动补全 |
| P3-T7.4 | 实现搜索结果页 | 分页/筛选 |

---

## P3-T8: SearchAgent - 跨模块索引

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T8 |
| **Agent** | SearchAgent |
| **模块** | 搜索服务 |
| **优先级** | P0 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 项目数据 | ProjectAgent | 项目索引 |
| 知识数据 | KnowledgeAgent | 知识索引 |
| 产品数据 | ShelfAgent | 产品索引 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 索引服务 | `services/api/services/indexer.go` | .go |
| 项目索引器 | `services/api/indexers/project_indexer.go` | .go |
| 知识索引器 | `services/api/indexers/knowledge_indexer.go` | .go |
| 产品索引器 | `services/api/indexers/product_indexer.go` | .go |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T8.1 | 实现项目索引器 | 索引更新正确 |
| P3-T8.2 | 实现知识索引器 | 索引更新正确 |
| P3-T8.3 | 实现产品索引器 | 索引更新正确 |
| P3-T8.4 | 实现统一搜索入口 | 跨模块搜索 |

---

## P3-T9: ForumAgent - 论坛基础功能

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T9 |
| **Agent** | ForumAgent |
| **模块** | 技术论坛 |
| **优先级** | P1 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 需求规格 | `03_需求规格说明书.md` | SRS-FORUM-001 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 论坛模型 | `services/api/models/forum.go` | .go |
| 论坛Handler | `services/api/handlers/forum.go` | .go |
| 论坛服务 | `services/api/services/forum.go` | .go |
| 论坛表迁移 | `database/migrations/015_forum.sql` | .sql |
| 板块列表页 | `apps/web/src/pages/forum/BoardList.tsx` | .tsx |
| 帖子列表页 | `apps/web/src/pages/forum/PostList.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T9.1 | 创建论坛数据模型 | 模型定义正确 |
| P3-T9.2 | 实现板块CRUD | 接口正常 |
| P3-T9.3 | 实现帖子CRUD | 接口正常 |
| P3-T9.4 | 实现板块列表 | 显示正确 |
| P3-T9.5 | 实现帖子列表 | 分页正确 |

---

## P3-T10: ForumAgent - 帖子发布与回复

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T10 |
| **Agent** | ForumAgent |
| **模块** | 技术论坛 |
| **优先级** | P1 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 用户数据 | UserAgent | 发帖用户 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 回复模型 | `services/api/models/reply.go` | .go |
| 回复Handler | `services/api/handlers/reply.go` | .go |
| 回复表迁移 | `database/migrations/015_replies.sql` | .sql |
| 发帖编辑器 | `apps/web/src/components/forum/PostEditor.tsx` | .tsx |
| 回复组件 | `apps/web/src/components/forum/ReplyThread.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T10.1 | 创建回复数据模型 | 模型定义正确 |
| P3-T10.2 | 实现回复API | 接口正常 |
| P3-T10.3 | 实现富文本编辑器 | Markdown支持 |
| P3-T10.4 | 实现@通知 | 通知发送 |
| P3-T10.5 | 实现最佳答案 | 标记功能 |

---

## P3-T11: ForumAgent - 搜索与标签

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T11 |
| **Agent** | ForumAgent |
| **模块** | 技术论坛 |
| **优先级** | P1 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 帖子数据 | P3-T9 | 帖子信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 论坛搜索 | `services/api/services/forum_search.go` | .go |
| 标签服务 | `services/api/services/forum_tag.go` | .go |
| 搜索页面 | `apps/web/src/pages/forum/ForumSearch.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T11.1 | 实现论坛全文搜索 | 搜索正常 |
| P3-T11.2 | 实现标签管理 | 标签分类 |
| P3-T11.3 | 实现标签筛选 | 筛选正确 |

---

## P3-T12: ForumAgent - 知识库关联

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P3-T12 |
| **Agent** | ForumAgent |
| **模块** | 技术论坛 |
| **优先级** | P2 |
| **阶段** | Phase 3 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 帖子数据 | P3-T9 | 帖子信息 |
| 知识数据 | KnowledgeAgent | 知识条目 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 知识关联服务 | `services/api/services/knowledge_link.go` | .go |
| 归档组件 | `apps/web/src/components/forum/ArchiveToKB.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P3-T12.1 | 实现帖子归档 | 归档到知识库 |
| P3-T12.2 | 实现关联推荐 | 知识推荐 |

---

## Phase 3 集成测试清单

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
