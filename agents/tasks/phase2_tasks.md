# Agent 任务卡片 - Phase 2 核心业务

---

## P2-T1: WorkflowAgent - 状态机引擎

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T1 |
| **Agent** | WorkflowAgent |
| **模块** | 流程引擎 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 流程模板 | ProjectAgent (P1-T11) | 模板数据 |
| 项目数据 | ProjectAgent | 项目信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 状态机模型 | `services/api/models/workflow.go` | .go |
| 状态机服务 | `services/api/services/statmachine.go` | .go |
| 流程表迁移 | `database/migrations/005_workflows.sql` | .sql |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T1.1 | 定义状态枚举 | 状态定义完整 |
| P2-T1.2 | 实现状态转换规则 | 转换逻辑正确 |
| P2-T1.3 | 实现状态校验 | 非法转换被阻止 |

---

## P2-T2: WorkflowAgent - 活动流转逻辑

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T2 |
| **Agent** | WorkflowAgent |
| **模块** | 流程引擎 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 状态机 | P2-T1 | 状态管理 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 活动Handler | `services/api/handlers/activity.go` | .go |
| 活动服务 | `services/api/services/activity.go` | .go |
| 活动表迁移 | `database/migrations/005_activities.sql` | .sql |
| 前端活动组件 | `apps/web/src/components/workflow/ActivityCard.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T2.1 | 创建活动数据模型 | 模型定义正确 |
| P2-T2.2 | 实现活动启动 | 状态变更正确 |
| P2-T2.3 | 实现活动完成 | 触发下一活动 |
| P2-T2.4 | 实现活动审批 | 审批流程正常 |
| P2-T2.5 | 实现前端活动卡片 | 状态显示正确 |

---

## P2-T3: WorkflowAgent - DCP评审节点

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T3 |
| **Agent** | WorkflowAgent |
| **模块** | 流程引擎 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 活动数据 | P2-T2 | 活动信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 评审模型 | `services/api/models/review.go` | .go |
| 评审Handler | `services/api/handlers/review.go` | .go |
| 评审表迁移 | `database/migrations/005_reviews.sql` | .sql |
| 前端评审组件 | `apps/web/src/components/workflow/DcpReview.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T3.1 | 创建评审数据模型 | 模型定义正确 |
| P2-T3.2 | 实现评审提交 | 状态变更正确 |
| P2-T3.3 | 实现评审通过 | 触发下一活动 |
| P2-T3.4 | 实现评审驳回 | 返回修改意见 |
| P2-T3.5 | 实现前端评审UI | 操作响应正确 |

---

## P2-T4: ProjectAgent (扩展) - Gitea API集成

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T4 |
| **Agent** | ProjectAgent |
| **模块** | 项目管理 - Gitea集成 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| Gitea服务 | 部署环境 | Gitea服务 |
| 项目数据 | P1-T9 | 项目基础数据 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| Gitea客户端 | `services/api/clients/gitea.go` | .go |
| Git服务 | `services/api/services/git.go` | .go |
| Git仓库表迁移 | `database/migrations/006_git_repos.sql` | .sql |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent (API调用安全性)

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T4.1 | 实现Gitea客户端 | API封装完整 |
| P2-T4.2 | 实现仓库创建 | 自动创建仓库 |
| P2-T4.3 | 实现目录初始化 | 标准目录创建 |

---

## P2-T5: ProjectAgent (扩展) - Git版本管理

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T5 |
| **Agent** | ProjectAgent |
| **模块** | 项目管理 - Git版本 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| Gitea集成 | P2-T4 | Git操作 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| Git操作Handler | `services/api/handlers/git.go` | .go |
| Git历史组件 | `apps/web/src/components/projects/GitHistory.tsx` | .tsx |
| 版本对比组件 | `apps/web/src/components/projects/VersionDiff.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T5.1 | 实现文件上传并提交 | Git提交正确 |
| P2-T5.2 | 实现版本历史 | 提交列表正确 |
| P2-T5.3 | 实现版本对比 | Diff显示正确 |

---

## P2-T6: ProjectAgent (扩展) - 甘特图

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T6 |
| **Agent** | ProjectAgent |
| **模块** | 项目管理 - 甘特图 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 1 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 活动数据 | P2-T2 | 活动列表 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 甘特图组件 | `apps/web/src/components/projects/GanttChart.tsx` | .tsx |
| 甘特图数据转换 | `apps/web/src/utils/gantt.ts` | .ts |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T6.1 | 实现甘特图组件 | gantt-task-react集成 |
| P2-T6.2 | 实现时间粒度切换 | 日/周/月 |
| P2-T6.3 | 实现拖拽调整 | 任务时间可调整 |
| P2-T6.4 | 实现依赖连线 | 依赖关系显示 |

---

## P2-T7: DevAgent - 流程全景视图

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T7 |
| **Agent** | DevAgent |
| **模块** | 项目开发 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 流程数据 | WorkflowAgent | 活动列表 |
| 项目数据 | ProjectAgent | 项目信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 流程图组件 | `apps/web/src/components/development/ProcessFlow.tsx` | .tsx |
| 流程节点组件 | `apps/web/src/components/development/ProcessNode.tsx` | .tsx |
| 活动详情抽屉 | `apps/web/src/components/development/ActivityDrawer.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T7.1 | 实现流程图组件 | BPMN风格展示 |
| P2-T7.2 | 实现节点状态 | 高亮当前活动 |
| P2-T7.3 | 实现详情抽屉 | 面板展示正确 |
| P2-T7.4 | 实现缩放拖拽 | 交互流畅 |

---

## P2-T8: DevAgent - 本地软件联动(rdp协议)

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T8 |
| **Agent** | DevAgent |
| **模块** | 项目开发 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 文件数据 | ProjectAgent | 文件列表 |
| 桌面程序 | DesktopAgent (并行) | 协议处理 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 协议处理服务 | `services/api/services/rdp_protocol.go` | .go |
| 文件打开Handler | `services/api/handlers/file_open.go` | .go |
| 文件类型映射 | `config/file_types.yaml` | .yaml |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent (安全)

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T8.1 | 定义协议格式 | rdp://open?格式 |
| P2-T8.2 | 实现协议处理 | URL解析正确 |
| P2-T8.3 | 配置文件类型映射 | 映射表完整 |
| P2-T8.4 | 实现前端协议调用 | 调用正确 |

---

## P2-T9: DevAgent - 活动执行面板

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T9 |
| **Agent** | DevAgent |
| **模块** | 项目开发 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 活动数据 | WorkflowAgent | 活动信息 |
| 文件数据 | ProjectAgent | 交付物 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 活动执行页面 | `apps/web/src/pages/development/ActivityExecute.tsx` | .tsx |
| 交付物列表 | `apps/web/src/components/development/DeliverableList.tsx` | .tsx |
| 模板选择器 | `apps/web/src/components/development/TemplateSelector.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T9.1 | 实现活动执行页面 | 布局正确 |
| P2-T9.2 | 实现交付物列表 | 状态显示正确 |
| P2-T9.3 | 实现模板选择 | 模板加载正确 |
| P2-T9.4 | 实现活动完成 | 状态变更正确 |

---

## P2-T10: DevAgent - 评审反馈系统

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T10 |
| **Agent** | DevAgent |
| **模块** | 项目开发 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 用户数据 | UserAgent | 评审人员 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 反馈模型 | `services/api/models/feedback.go` | .go |
| 反馈Handler | `services/api/handlers/feedback.go` | .go |
| 反馈表迁移 | `database/migrations/008_feedbacks.sql` | .sql |
| 反馈组件 | `apps/web/src/components/development/FeedbackThread.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T10.1 | 创建反馈数据模型 | 模型正确 |
| P2-T10.2 | 实现反馈提交 | 提交成功 |
| P2-T10.3 | 实现@通知 | 通知发送 |
| P2-T10.4 | 实现反馈回复 | 楼中楼正常 |

---

## P2-T11: DevAgent - 变更管理

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T11 |
| **Agent** | DevAgent |
| **模块** | 项目开发 |
| **优先级** | P1 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 活动数据 | WorkflowAgent | 活动信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 变更模型 | `services/api/models/change.go` | .go |
| 变更Handler | `services/api/handlers/change.go` | .go |
| 变更表迁移 | `database/migrations/008_changes.sql` | .sql |
| 变更页面 | `apps/web/src/pages/development/ChangeRequest.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T11.1 | 创建变更模型 | ECR/ECO定义 |
| P2-T11.2 | 实现变更申请 | 申请流程正确 |
| P2-T11.3 | 实现影响分析 | 影响项列出 |
| P2-T11.4 | 实现审批闭环 | 审批流程完整 |

---

## P2-T12: ShelfAgent - 产品浏览与筛选

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T12 |
| **Agent** | ShelfAgent |
| **模块** | 产品货架 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 项目数据 | ProjectAgent | 来源项目 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 产品模型 | `services/api/models/product.go` | .go |
| 产品Handler | `services/api/handlers/product.go` | .go |
| 产品服务 | `services/api/services/product.go` | .go |
| 产品表迁移 | `database/migrations/009_products.sql` | .sql |
| 产品列表页 | `apps/web/src/pages/shelf/ProductList.tsx` | .tsx |
| 产品卡片 | `apps/web/src/components/shelf/ProductCard.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T12.1 | 创建产品数据模型 | 模型定义正确 |
| P2-T12.2 | 实现产品CRUD API | 接口正常 |
| P2-T12.3 | 实现产品列表 | 多维筛选正常 |
| P2-T12.4 | 实现产品卡片 | 成熟度标签显示 |
| P2-T12.5 | 实现产品详情 | 完整信息展示 |

---

## P2-T13: ShelfAgent - 选用购物车

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T13 |
| **Agent** | ShelfAgent |
| **模块** | 产品货架 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 产品数据 | P2-T12 | 产品信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 购物车服务 | `services/api/services/cart.go` | .go |
| 购物车Handler | `services/api/handlers/cart.go` | .go |
| 购物车表迁移 | `database/migrations/009_cart.sql` | .sql |
| 购物车组件 | `apps/web/src/components/shelf/ShoppingCart.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T13.1 | 创建购物车模型 | 模型定义正确 |
| P2-T13.2 | 实现加入购物车 | 添加成功 |
| P2-T13.3 | 实现一键导入 | 导入项目成功 |
| P2-T13.4 | 实现购物车UI | 列表显示正确 |

---

## P2-T14: ShelfAgent - 技术树可视化

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T14 |
| **Agent** | ShelfAgent |
| **模块** | 技术货架 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 技术数据 | 本模块 | 技术条目 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 技术模型 | `services/api/models/technology.go` | .go |
| 技术Handler | `services/api/handlers/technology.go` | .go |
| 技术表迁移 | `database/migrations/010_technologies.sql` | .sql |
| 技术树组件 | `apps/web/src/components/shelf/TechnologyTree.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T14.1 | 创建技术数据模型 | TRL等级定义 |
| P2-T14.2 | 实现技术CRUD API | 接口正常 |
| P2-T14.3 | 实现技术树展示 | 3层深度 |
| P2-T14.4 | 实现TRL标签 | 等级显示正确 |

---

## P2-T15: ShelfAgent - 版本管理与自动上架

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T15 |
| **Agent** | ShelfAgent |
| **模块** | 产品/技术货架 |
| **优先级** | P1 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 产品数据 | P2-T12 | 产品信息 |
| 项目流程 | WorkflowAgent | 完成状态 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 版本服务 | `services/api/services/version.go` | .go |
| 自动上架服务 | `services/api/services/auto_publish.go` | .go |
| 版本组件 | `apps/web/src/components/shelf/VersionTree.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T15.1 | 实现版本管理 | 版本链显示 |
| P2-T15.2 | 实现Fork创建 | 新版本创建 |
| P2-T15.3 | 实现自动上架 | 项目完成触发 |

---

## P2-T16: DesktopAgent - rdp协议注册

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T16 |
| **Agent** | DesktopAgent |
| **模块** | 桌面辅助程序 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 协议定义 | DevAgent (P2-T8) | 协议格式 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| Tauri应用结构 | `desktop/rdp-helper/` | 目录 |
| 协议注册模块 | `desktop/rdp-helper/src/protocol.rs` | .rs |
| Windows安装包 | `desktop/releases/rdp-helper.exe` | .exe |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T16.1 | 创建Tauri项目 | 项目结构正确 |
| P2-T16.2 | 实现协议注册 | 系统注册成功 |
| P2-T16.3 | 实现协议解析 | URL解析正确 |
| P2-T16.4 | 构建Windows安装包 | 可执行文件生成 |

---

## P2-T17: DesktopAgent - 本地软件调用

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T17 |
| **Agent** | DesktopAgent |
| **模块** | 桌面辅助程序 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 文件类型映射 | DevAgent | 映射配置 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 文件调用模块 | `desktop/rdp-helper/src/file_handler.rs` | .rs |
| 应用配置 | `desktop/rdp-helper/config/app_paths.yaml` | .yaml |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T17.1 | 实现文件定位 | 路径映射正确 |
| P2-T17.2 | 实现软件调用 | 应用启动正确 |
| P2-T17.3 | 配置应用路径 | 配置文件完整 |

---

## P2-T18: DesktopAgent - Git自动提交

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T18 |
| **Agent** | DesktopAgent |
| **模块** | 桌面辅助程序 |
| **优先级** | P0 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 活动完成通知 | DevAgent | WebSocket通知 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| Git操作模块 | `desktop/rdp-helper/src/git_ops.rs` | .rs |
| 提交服务 | `desktop/rdp-helper/src/commit_service.rs` | .rs |
| 托盘组件 | `desktop/rdp-helper/src/tray.rs` | .rs |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T18.1 | 实现Git操作 | commit/push正确 |
| P2-T18.2 | 实现自动提交 | 活动完成触发 |
| P2-T18.3 | 实现提交消息 | 格式正确 |
| P2-T18.4 | 实现状态托盘 | 状态显示正确 |

---

## P2-T19: DesktopAgent - 冲突检测与解决

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T19 |
| **Agent** | DesktopAgent |
| **模块** | 桌面辅助程序 |
| **优先级** | P1 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| Git操作 | P2-T18 | Git状态 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 冲突检测模块 | `desktop/rdp-helper/src/conflict.rs` | .rs |
| 冲突解决UI | `desktop/rdp-helper/src/ui/conflict.rs` | .rs |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T19.1 | 实现冲突检测 | 检测逻辑正确 |
| P2-T19.2 | 实现冲突提示 | 提示显示正确 |
| P2-T19.3 | 实现解决选项 | 选项操作正确 |

---

## P2-T20: QMAgent - 需求管理

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T20 |
| **Agent** | QMAgent |
| **模块** | 质量管理 |
| **优先级** | P1 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 项目数据 | ProjectAgent | 项目信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 需求模型 | `services/api/models/requirement.go` | .go |
| 需求Handler | `services/api/handlers/requirement.go` | .go |
| 需求表迁移 | `database/migrations/011_requirements.sql` | .sql |
| 需求页面 | `apps/web/src/pages/quality/Requirements.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T20.1 | 创建需求模型 | 模型定义正确 |
| P2-T20.2 | 实现需求CRUD | 接口正常 |
| P2-T20.3 | 实现需求分解 | 分解关系正确 |
| P2-T20.4 | 实现需求追溯 | 追溯关系显示 |

---

## P2-T21: QMAgent - 变更管理(ECR/ECO)

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T21 |
| **Agent** | QMAgent |
| **模块** | 质量管理 |
| **优先级** | P1 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 需求数据 | P2-T20 | 需求信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| ECR模型 | `services/api/models/ecr.go` | .go |
| ECO模型 | `services/api/models/eco.go` | .go |
| 变更表迁移 | `database/migrations/012_changes.sql` | .sql |
| 变更页面 | `apps/web/src/pages/quality/ChangeRequest.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent
- **L4**: 人类监督者

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T21.1 | 创建ECR模型 | ECR定义正确 |
| P2-T21.2 | 创建ECO模型 | ECO定义正确 |
| P2-T21.3 | 实现ECR申请 | 申请流程正确 |
| P2-T21.4 | 实现ECO执行 | 执行闭环正确 |

---

## P2-T22: QMAgent - 缺陷管理

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T22 |
| **Agent** | QMAgent |
| **模块** | 质量管理 |
| **优先级** | P1 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| 项目数据 | ProjectAgent | 项目信息 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 缺陷模型 | `services/api/models/defect.go` | .go |
| 缺陷Handler | `services/api/handlers/defect.go` | .go |
| 缺陷表迁移 | `database/migrations/013_defects.sql` | .sql |
| 缺陷页面 | `apps/web/src/pages/quality/Defects.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T22.1 | 创建缺陷模型 | 缺陷状态定义 |
| P2-T22.2 | 实现缺陷CRUD | 接口正常 |
| P2-T22.3 | 实现缺陷分配 | 分配正确 |
| P2-T22.4 | 实现缺陷关闭 | 闭环正确 |

---

## P2-T23: QMAgent - 质量门禁

### 基本信息
| 项目 | 内容 |
|------|------|
| **任务ID** | P2-T23 |
| **Agent** | QMAgent |
| **模块** | 质量管理 |
| **优先级** | P1 |
| **阶段** | Phase 2 - Layer 2 |

### 输入
| 依赖项 | 来源 | 说明 |
|--------|------|------|
| DCP评审 | WorkflowAgent | 评审状态 |

### 输出
| 交付物 | 路径 | 格式 |
|--------|------|------|
| 门禁服务 | `services/api/services/gate.go` | .go |
| 门禁配置 | `config/quality_gates.yaml` | .yaml |
| 门禁组件 | `apps/web/src/components/quality/GateIndicator.tsx` | .tsx |

### 检查者
- **L1**: Agent 自审查
- **L2**: Reviewer Agent

### 子任务清单
| 任务ID | 子任务 | 验收标准 |
|--------|--------|----------|
| P2-T23.1 | 定义门禁规则 | 规则配置正确 |
| P2-T23.2 | 实现门禁检查 | 检查逻辑正确 |
| P2-T23.3 | 实现阻止放行 | 未通过阻止 |
| P2-T23.4 | 实现门禁显示 | 状态指示正确 |

---

## Phase 2 集成测试清单

| 序号 | 测试项 | 验收条件 |
|------|--------|----------|
| IT-21 | 流程创建 | 项目绑定流程成功 |
| IT-22 | 活动流转 | 状态变更正确 |
| IT-23 | DCP评审 | 评审通过/驳回正常 |
| IT-24 | Git仓库创建 | 自动创建仓库成功 |
| IT-25 | 文件版本管理 | Git历史正确 |
| IT-26 | 甘特图交互 | 拖拽调整正常 |
| IT-27 | 流程全景视图 | BPMN图显示正确 |
| IT-28 | rdp协议调用 | 本地软件打开正确 |
| IT-29 | 活动完成提交 | Git自动提交成功 |
| IT-30 | 产品货架浏览 | 筛选正常 |
| IT-31 | 选用购物车 | 加入/导入正常 |
| IT-32 | 技术树展示 | 树形显示正确 |
| IT-33 | 桌面程序安装 | 协议注册成功 |
| IT-34 | 需求管理 | 需求CRUD正常 |
| IT-35 | 变更管理 | ECR/ECO流程正常 |
| IT-36 | 缺陷管理 | 缺陷生命周期正常 |
| IT-37 | 质量门禁 | 未通过阻止流程 |
