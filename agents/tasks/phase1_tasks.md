# 微波室研发管理平台 (RDP)
# Agent 任务卡片 - Phase 1 基础骨架

---

## P1-T1: PortalAgent - 部门门户首页

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T1 |
| **Agent** | PortalAgent |
| **模块** | 门户界面 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 技术栈规范 | `AGENTS.md` | React + TS + Ant Design |
| UI设计参考 | `01_需求文档.md` | 参考 department_homepage |
| 后端API | UserAgent (待开发) | 获取公告数据 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 部门门户首页组件 | `apps/web/src/pages/portal/PortalPage.tsx` | .tsx |
| 公告列表组件 | `apps/web/src/components/portal/AnnouncementList.tsx` | .tsx |
| 荣誉展示组件 | `apps/web/src/components/portal/HonorCarousel.tsx` | .tsx |
| 导航卡片组件 | `apps/web/src/components/portal/NavCards.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查（代码规范、语法检查）
- **L2**: Reviewer Agent（UI一致性、组件规范）
- **L4**: 人类监督者（功能验收）

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T1.1 | 创建门户布局组件 | 响应式布局，适配1920×1080 |
| P1-T1.2 | 实现公告列表 | 支持分页，数据接口mock |
| P1-T1.3 | 实现荣誉轮播 | 轮播动画流畅 |
| P1-T1.4 | 实现导航卡片 | 点击跳转正确 |
| P1-T1.5 | 集成部门简介 | 富文本渲染正常 |

---

## P1-T2: PortalAgent - 个人工作台

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T2 |
| **Agent** | PortalAgent |
| **模块** | 门户界面 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 用户认证信息 | UserAgent | 用户登录态 |
| 项目数据 | ProjectAgent | 我的项目列表 |
| 通知数据 | 本模块 | 消息通知 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 工作台页面组件 | `apps/web/src/pages/workbench/WorkbenchPage.tsx` | .tsx |
| 待办事项组件 | `apps/web/src/components/workbench/TodoList.tsx` | .tsx |
| 我的项目组件 | `apps/web/src/components/workbench/MyProjects.tsx` | .tsx |
| 快捷操作面板 | `apps/web/src/components/workbench/QuickActions.tsx` | .tsx |
| 统计概览卡片 | `apps/web/src/components/workbench/StatsOverview.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T2.1 | 创建工作台布局 | 布局结构正确 |
| P1-T2.2 | 实现待办事项列表 | 按紧急程度排序 |
| P1-T2.3 | 实现我的项目卡片 | 显示进度百分比 |
| P1-T2.4 | 实现快捷操作 | 操作响应正确 |
| P1-T2.5 | 实现统计概览 | 数据展示正确 |

---

## P1-T3: PortalAgent - 消息通知中心

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T3 |
| **Agent** | PortalAgent |
| **模块** | 门户界面 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 通知API | UserAgent | 通知列表接口 |
| WebSocket | 本模块 | 实时推送 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 通知中心组件 | `apps/web/src/components/notification/NotificationCenter.tsx` | .tsx |
| 通知列表项 | `apps/web/src/components/notification/NotificationItem.tsx` | .tsx |
| 通知图标Badge | `apps/web/src/components/notification/NotificationBadge.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T3.1 | 实现通知中心抽屉 | 抽屉展开/收起动画正常 |
| P1-T3.2 | 实现通知列表 | 支持未读/已读筛选 |
| P1-T3.3 | 实现通知Badge | 数字显示正确 |
| P1-T3.4 | 实现通知已读 | 点击标记已读 |

---

## P1-T4: PortalAgent - 全局搜索UI

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T4 |
| **Agent** | PortalAgent |
| **模块** | 门户界面 |
| **优先级** | P1 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 搜索API | SearchAgent (Phase 3) | 搜索接口（可先mock） |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 全局搜索栏组件 | `apps/web/src/components/search/GlobalSearch.tsx` | .tsx |
| 搜索结果下拉 | `apps/web/src/components/search/SearchDropdown.tsx` | .tsx |
| 搜索结果页面 | `apps/web/src/pages/search/SearchResultPage.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T4.1 | 实现搜索栏 | 支持回车搜索 |
| P1-T4.2 | 实现防抖搜索 | 300ms防抖 |
| P1-T4.3 | 实现结果下拉 | 显示搜索建议 |
| P1-T4.4 | 实现结果页面 | 分页展示结果 |

---

## P1-T5: UserAgent - 用户认证与CRUD

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T5 |
| **Agent** | UserAgent |
| **模块** | 用户管理 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 需求规格 | `03_需求规格说明书.md` | SRS-USER-001 |
| Casdoor配置 | `02_详细实施方案.md` | Casdoor集成方案 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 用户模型 | `services/api/models/user.go` | .go |
| 用户Handler | `services/api/handlers/user.go` | .go |
| 认证服务 | `services/api/services/auth.go` | .go |
| 用户表迁移 | `database/migrations/001_users.sql` | .sql |
| API文档 | `docs/api/user-api.md` | .md |

### 检查者
- **L1**: Agent 自审查（代码规范、语法检查）
- **L2**: Reviewer Agent（权限控制、安全漏洞）
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T5.1 | 创建用户数据模型 | 符合数据模型规范 |
| P1-T5.2 | 实现用户CRUD API | 接口符合RESTful规范 |
| P1-T5.3 | 集成Casdoor认证 | 登录/登出正常 |
| P1-T5.4 | 实现JWT Token管理 | Token刷新正常 |
| P1-T5.5 | 创建用户表迁移 | SQL执行无错误 |
| P1-T5.6 | 编写API文档 | 文档完整准确 |

---

## P1-T6: UserAgent - RBAC权限模型

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T6 |
| **Agent** | UserAgent |
| **模块** | 用户管理 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 权限矩阵 | `03_需求规格说明书.md` | SRS-USER-002 |
| Casbin配置 | `02_详细实施方案.md` | 权限模型定义 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 权限服务 | `services/api/services/permission.go` | .go |
| Casbin模型 | `config/casbin_model.conf` | .conf |
| 权限中间件 | `services/api/middleware/auth.go` | .go |
| 权限表迁移 | `database/migrations/001_permissions.sql` | .sql |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent（权限控制测试）
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T6.1 | 配置Casbin权限模型 | 模型定义正确 |
| P1-T6.2 | 实现权限服务 | 权限检查正常 |
| P1-T6.3 | 实现权限中间件 | 拦截未授权请求 |
| P1-T6.4 | 实现数据权限隔离 | 团队/产品线隔离生效 |

---

## P1-T7: UserAgent - 组织架构管理

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T7 |
| **Agent** | UserAgent |
| **模块** | 用户管理 |
| **优先级** | P1 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 用户数据 | P1-T5 | 用户基础数据 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 组织架构模型 | `services/api/models/organization.go` | .go |
| 组织架构Handler | `services/api/handlers/organization.go` | .go |
| 组织架构表迁移 | `database/migrations/001_organizations.sql` | .sql |
| 前端组织树组件 | `apps/web/src/components/org/OrgTree.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T7.1 | 创建组织架构模型 | 模型定义正确 |
| P1-T7.2 | 实现组织架构API | CRUD正常 |
| P1-T7.3 | 实现前端组织树 | 树形展示正确 |

---

## P1-T8: UserAgent - 个人Profile页面

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T8 |
| **Agent** | UserAgent |
| **模块** | 用户管理 |
| **优先级** | P1 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 用户数据 | P1-T5 | 用户基础数据 |
| 项目数据 | ProjectAgent | 项目履历 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| Profile页面 | `apps/web/src/pages/profile/ProfilePage.tsx` | .tsx |
| 能力雷达图 | `apps/web/src/components/profile/AbilityRadar.tsx` | .tsx |
| 贡献热力图 | `apps/web/src/components/profile/ContributionHeatmap.tsx` | .tsx |
| 项目履历 | `apps/web/src/components/profile/ProjectHistory.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T8.1 | 实现Profile页面布局 | 布局符合设计 |
| P1-T8.2 | 实现能力雷达图 | ECharts渲染正确 |
| P1-T8.3 | 实现贡献热力图 | 数据展示正确 |
| P1-T8.4 | 实现项目履历 | 列表展示正确 |

---

## P1-T9: ProjectAgent - 项目CRUD与看板

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T9 |
| **Agent** | ProjectAgent |
| **模块** | 项目管理 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 需求规格 | `03_需求规格说明书.md` | SRS-PM-001/002 |
| 用户数据 | UserAgent | 项目成员 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 项目模型 | `services/api/models/project.go` | .go |
| 项目Handler | `services/api/handlers/project.go` | .go |
| 项目服务 | `services/api/services/project.go` | .go |
| 项目表迁移 | `database/migrations/002_projects.sql` | .sql |
| 项目看板页面 | `apps/web/src/pages/projects/ProjectBoard.tsx` | .tsx |
| 项目列表组件 | `apps/web/src/components/projects/ProjectList.tsx` | .tsx |
| 项目卡片组件 | `apps/web/src/components/projects/ProjectCard.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T9.1 | 创建项目数据模型 | 模型定义正确 |
| P1-T9.2 | 实现项目CRUD API | 接口符合规范 |
| P1-T9.3 | 实现项目看板 | 多维筛选正常 |
| P1-T9.4 | 实现项目列表 | 分页展示正常 |
| P1-T9.5 | 实现项目卡片 | 状态/进度显示正确 |

---

## P1-T10: ProjectAgent - 五步项目创建向导

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T10 |
| **Agent** | ProjectAgent |
| **模块** | 项目管理 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 项目数据 | P1-T9 | 项目基础数据 |
| 用户数据 | UserAgent | 成员选择 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 创建向导页面 | `apps/web/src/pages/projects/CreateWizard.tsx` | .tsx |
| 基本信息步骤 | `apps/web/src/components/wizard/StepBasicInfo.tsx` | .tsx |
| 类别选择步骤 | `apps/web/src/components/wizard/StepCategory.tsx` | .tsx |
| 流程绑定步骤 | `apps/web/src/components/wizard/StepProcess.tsx` | .tsx |
| 团队分配步骤 | `apps/web/src/components/wizard/StepTeam.tsx` | .tsx |
| 计划确认步骤 | `apps/web/src/components/wizard/StepConfirm.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T10.1 | 实现向导布局 | 步骤切换正常 |
| P1-T10.2 | 实现基本信息步骤 | 表单验证正确 |
| P1-T10.3 | 实现类别选择步骤 | 模板推荐正确 |
| P1-T10.4 | 实现流程绑定步骤 | 活动预览正常 |
| P1-T10.5 | 实现团队分配步骤 | 成员选择正常 |
| P1-T10.6 | 实现计划确认步骤 | 甘特图预览正常 |

---

## P1-T11: ProjectAgent - 流程模板管理

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T11 |
| **Agent** | ProjectAgent |
| **模块** | 项目管理 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 需求规格 | `03_需求规格说明书.md` | 流程模板定义 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 流程模板模型 | `services/api/models/process_template.go` | .go |
| 流程模板Handler | `services/api/handlers/process_template.go` | .go |
| 流程模板表迁移 | `database/migrations/002_process_templates.sql` | .sql |
| 种子数据 | `database/seeds/process_templates.sql` | .sql |
| 前端模板管理页面 | `apps/web/src/pages/projects/ProcessTemplates.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T11.1 | 创建流程模板模型 | 模型定义正确 |
| P1-T11.2 | 实现模板CRUD API | 接口正常 |
| P1-T11.3 | 预置种子数据 | 7种项目类别模板 |
| P1-T11.4 | 实现模板管理页面 | 增删改查正常 |

---

## P1-T12: ProjectAgent - 基础文件管理

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T12 |
| **Agent** | ProjectAgent |
| **模块** | 项目管理 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 项目数据 | P1-T9 | 项目信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 文件模型 | `services/api/models/file.go` | .go |
| 文件Handler | `services/api/handlers/file.go` | .go |
| 文件服务 | `services/api/services/file.go` | .go |
| 文件表迁移 | `database/migrations/002_files.sql` | .sql |
| 文件管理页面 | `apps/web/src/components/projects/FileManager.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T12.1 | 创建文件数据模型 | 模型定义正确 |
| P1-T12.2 | 实现文件上传API | 支持大文件分片 |
| P1-T12.3 | 实现文件下载API | 下载正常 |
| P1-T12.4 | 实现文件浏览 | 目录树展示正常 |

---

## P1-T13: SecurityAgent - 数据分级分类

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T13 |
| **Agent** | SecurityAgent |
| **模块** | 安全合规 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 需求规格 | `03_需求规格说明书.md` | SRS-SECURITY-001 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 分级模型 | `services/api/models/classification.go` | .go |
| 分级Handler | `services/api/handlers/classification.go` | .go |
| 分级服务 | `services/api/services/classification.go` | .go |
| 分级表迁移 | `database/migrations/003_classification.sql` | .sql |
| 前端分级组件 | `apps/web/src/components/security/ClassificationTag.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T13.1 | 创建分级数据模型 | 4级密定义正确 |
| P1-T13.2 | 实现分级设置API | CRUD正常 |
| P1-T13.3 | 实现分级权限检查 | 未授权无法访问 |
| P1-T13.4 | 实现前端分级标签 | 标签显示正确 |

---

## P1-T14: SecurityAgent - 会话超时控制

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T14 |
| **Agent** | SecurityAgent |
| **模块** | 安全合规 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 需求规格 | `03_需求规格说明书.md` | SRS-SECURITY-002 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 会话服务 | `services/api/services/session.go` | .go |
| 会话中间件 | `services/api/middleware/session.go` | .go |
| 前端心跳组件 | `apps/web/src/utils/heartbeat.ts` | .ts |
| 前端超时提示 | `apps/web/src/components/security/SessionWarning.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T14.1 | 实现会话服务 | 会话管理正常 |
| P1-T14.2 | 实现超时中间件 | 30分钟自动登出 |
| P1-T14.3 | 实现前端心跳 | 定时发送心跳 |
| P1-T14.4 | 实现超时警告 | 25分钟弹出提示 |

---

## P1-T15: SecurityAgent - 操作审计日志

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T15 |
| **Agent** | SecurityAgent |
| **模块** | 安全合规 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 需求规格 | `03_需求规格说明书.md` | SRS-SECURITY-003 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 审计模型 | `services/api/models/audit.go` | .go |
| 审计Handler | `services/api/handlers/audit.go` | .go |
| 审计服务 | `services/api/services/audit.go` | .go |
| 审计表迁移 | `database/migrations/004_audit_logs.sql` | .sql |
| 审计查询页面 | `apps/web/src/pages/admin/AuditLogs.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T15.1 | 创建审计数据模型 | 模型定义正确 |
| P1-T15.2 | 实现审计日志记录 | 关键操作记录完整 |
| P1-T15.3 | 实现审计查询API | 支持多条件筛选 |
| P1-T15.4 | 实现审计页面 | 分页展示正常 |

---

## P1-T16: SecurityAgent - 屏幕水印(基础)

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T16 |
| **Agent** | SecurityAgent |
| **模块** | 安全合规 |
| **优先级** | P1 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 分级数据 | P1-T13 | 密级信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 水印组件 | `apps/web/src/components/security/Watermark.tsx` | .tsx |
| 水印工具 | `apps/web/src/utils/watermark.ts` | .ts |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T16.1 | 实现CSS水印 | 覆盖整个页面 |
| P1-T16.2 | 实现动态水印内容 | 包含用户名时间 |
| P1-T16.3 | 集成密级判断 | 秘密/机密级显示 |

---

## P1-T17: InfraAgent - 数据库初始化

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T17 |
| **Agent** | InfraAgent |
| **模块** | 基础设施 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 数据模型 | 各Agent模型 | 汇总所有表结构 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 初始化SQL | `database/init.sql` | .sql |
| 枚举类型SQL | `database/schema/enums.sql` | .sql |
| 用户表迁移 | `database/migrations/001_users.sql` | .sql |
| 组织表迁移 | `database/migrations/001_organizations.sql` | .sql |
| 项目表迁移 | `database/migrations/002_projects.sql` | .sql |
| 文件表迁移 | `database/migrations/002_files.sql` | .sql |
| 分级表迁移 | `database/migrations/003_classification.sql` | .sql |
| 审计表迁移 | `database/migrations/004_audit_logs.sql` | .sql |

### 检查者
- **L1**: Agent 自审查（SQL语法检查）
- **L2**: Reviewer Agent（表结构评审）

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T17.1 | 创建枚举类型 | 类型定义完整 |
| P1-T17.2 | 创建初始化SQL | 包含所有表 |
| P1-T17.3 | 创建索引SQL | 索引设计合理 |
| P1-T17.4 | 创建种子数据 | 初始用户/组织 |

---

## P1-T18: InfraAgent - systemd服务配置

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T18 |
| **Agent** | InfraAgent |
| **模块** | 基础设施 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 服务配置 | `02_详细实施方案.md` | 服务定义 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| API服务配置 | `deploy/systemd/rdp-api.service` | .service |
| Casdoor服务配置 | `deploy/systemd/rdp-casdoor.service` | .service |
| 服务管理脚本 | `deploy/scripts/servicectl.sh` | .sh |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T18.1 | 创建API服务配置 | 配置正确可执行 |
| P1-T18.2 | 创建Casdoor配置 | 配置正确可执行 |
| P1-T18.3 | 创建服务管理脚本 | 启停脚本正常 |

---

## P1-T19: InfraAgent - Nginx反向代理

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T19 |
| **Agent** | InfraAgent |
| **模块** | 基础设施 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 端口规划 | `02_详细实施方案.md` | 端口定义 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| Nginx主配置 | `deploy/nginx/nginx.conf` | .conf |
| RDP站点配置 | `deploy/nginx/sites-available/rdp.conf` | .conf |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T19.1 | 创建主配置文件 | 配置正确 |
| P1-T19.2 | 创建站点配置 | 路由规则正确 |

---

## P1-T20: InfraAgent - 一键安装脚本

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P1-T20 |
| **Agent** | InfraAgent |
| **模块** | 基础设施 |
| **优先级** | P0 |
| **阶段** | Phase 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 依赖组件 | `02_详细实施方案.md` | 组件清单 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 安装脚本 | `deploy/scripts/install.sh` | .sh |
| 备份脚本 | `deploy/scripts/backup.sh` | .sh |
| 健康检查 | `deploy/scripts/health-check.sh` | .sh |
| 应用配置 | `config/rdp-api.yaml` | .yaml |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent (脚本可执行性)
- **L4**: 人类监督者 (最终验收)

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P1-T20.1 | 创建安装脚本 | Ubuntu/CentOS支持 |
| P1-T20.2 | 创建备份脚本 | 每日自动备份 |
| P1-T20.3 | 创建健康检查 | 服务状态检测 |
| P1-T20.4 | 创建应用配置 | 配置项完整 |

---

## Phase 1 集成测试清单

| 序号 | 测试项 | 验收条件 |
|------|--------|----------|
| IT-01 | 用户登录 | 登录成功，获取Token |
| IT-02 | 用户权限 | 不同角色权限正确 |
| IT-03 | 项目创建 | 五步向导可完成 |
| IT-04 | 项目列表 | 筛选/分页正常 |
| IT-05 | 文件上传 | 大文件上传成功 |
| IT-06 | 密级控制 | 未授权无法访问 |
| IT-07 | 审计日志 | 操作记录完整 |
| IT-08 | 会话超时 | 30分钟自动登出 |
| IT-09 | 门户首页 | 公告/荣誉正常显示 |
| IT-10 | 工作台 | 待办/项目正常显示 |
| IT-11 | 服务部署 | 一键安装可执行 |
| IT-12 | Nginx代理 | 访问正常 |
