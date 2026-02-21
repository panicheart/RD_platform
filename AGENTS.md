# 微波工程部研发管理平台 (RDP) - 项目文档

> **文档用途**: 本文档供 AI 编程助手参考，用于理解项目架构、技术栈和开发规范。
> **最后更新**: 2026-02-21
> **版本**: V1.2（AI Agent开发版）

> **⚠️ 开发模式声明**: 本项目采用**AI Agent集群自主开发模式**，各功能模块由专门的AI Agent负责开发，PM-Agent（项目经理Agent）进行全局质量把控。本文档中所有"开发团队"、"负责人员"均指对应的AI Agent角色。

---

## 1. 项目概述

### 1.1 项目简介

**微波工程部研发管理平台** (R&D Platform, 简称 **RDP**) 是为微波系统研发部门设计的企业级综合管理平台。

| 属性 | 内容 |
|------|------|
| **系统名称** | 微波工程部研发管理平台 (RDP) |
| **目标用户** | 微波系统研发部门（约30-100人） |
| **部署环境** | 离线局域网，systemd裸机部署 |
| **开发语言** | 中文（界面）、英文（代码） |
| **文档位置** | `方案/` 目录为项目文档 |

### 1.2 项目定位

- **项目管理中心** - 统筹项目全生命周期（IPD方法论）
- **技术资产库** - 沉淀和复用技术成果
- **知识管理平台** - 系统化智力资产管理
- **协同工作空间** - 支持团队协作与沟通

### 1.3 目录结构

```
RD_platform/
├── 原始需求/                    # 原始需求文档（PDF）
│   └── 研发平台需求.pdf
└── 方案/                        # 项目文档
    ├── README.md                # 文档导航入口
    ├── 01_需求文档.md
    ├── 02_详细实施方案.md
    ├── 03_需求规格说明书.md
    ├── 技术架构分析报告.md
    └── pdf-html/                # PDF和HTML版本
        ├── 01-requirements-analysis.html
        ├── 01-需求梳理与补全文档.pdf
        ├── 02-implementation-plan.html
        ├── 02-详细实施方案.pdf
        ├── 03-requirements-specification.html
        └── 03-需求规格说明书.pdf
```

---

## 2. 技术栈

### 2.1 核心技术约束

| 层次 | 技术 | 版本 | 约束 |
|------|------|------|------|
| **前端框架** | React + TypeScript + Vite | React 18.x, TS 5.x, Vite 5.x | [MUST] 单体应用，非微前端 |
| **UI库** | Ant Design | 5.x | [MUST] 主UI库 |
| **后端API** | Go (Gin framework) | Go 1.22+, Gin 1.9+ | [MUST] 纯 Go 技术栈 |
| **数据库** | PostgreSQL | 16.x | [MUST] 系统包安装 |
| **缓存** | Redis | 7.x | [SHOULD] Phase 2 |
| **搜索** | PostgreSQL全文搜索(早期) / MeiliSearch(后期) | — | [MUST] P1用PG, P3用Meili |
| **Git服务** | Gitea | 1.22+ | [MUST] Phase 2 |
| **认证** | Casdoor | Latest | [MUST] Phase 1 |
| **IM** | 内置通知(早期) / Mattermost(后期) | — | [SHOULD] P1内置, P3外接 |
| **存储** | 本地文件系统 | — | [MUST] |
| **服务管理** | systemd | — | [MUST] 裸机部署 |

### 2.2 前端技术栈详情

| 技术 | 版本 | 用途 |
|------|------|------|
| React | 18.x+ | 前端框架 |
| TypeScript | 5.0+ | 类型系统 |
| Ant Design | 5.x | UI组件库 |
| Zustand | 4.x | 状态管理 |
| React Router | 6.x | 路由管理 |
| Axios | 1.6+ | HTTP客户端 |
| Vite | 5.0+ | 构建工具 |
| Tailwind CSS | 3.x | CSS框架 |
| ECharts | 5.4+ | 数据可视化 |
| AntV G6 | 4.8+ | 流程图/技术树 |

### 2.3 后端技术栈详情

| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.22+ | 编程语言 |
| Gin | 1.9+ | Web框架 |
| GORM | 2.0+ | ORM框架 |
| Casbin | 2.x+ | 权限管理 |
| Viper | 1.18+ | 配置管理 |
| Zap | 1.26+ | 日志框架 |
| Flowable | 7.0+ | 工作流引擎 |

---

## 3. 系统架构

### 3.1 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                        用户访问层                              │
│     门户界面  项目管理  项目开发  产品货架  技术货架              │
│     知识库    技术论坛  即时通讯                              │
├─────────────────────────────────────────────────────────────┤
│                        API网关层                              │
│              路由转发  认证鉴权  限流熔断                      │
├─────────────────────────────────────────────────────────────┤
│                        服务层                                  │
│     用户服务  项目服务  流程服务  文件服务  通知服务            │
│     搜索服务  报表服务  IM服务                                │
├─────────────────────────────────────────────────────────────┤
│                        外部系统集成层                          │
│     Casdoor(认证)  Gitea(Git)  Mattermost(IM)                │
├─────────────────────────────────────────────────────────────┤
│                        数据层                                  │
│     PostgreSQL    Redis    MeiliSearch    MinIO              │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 服务端口分配

| 服务 | 端口 | 说明 |
|------|------|------|
| Nginx HTTP | 80 | HTTP入口 |
| Nginx HTTPS | 443 | HTTPS入口 |
| RDP Portal | 3000 | 前端应用 |
| Gitea | 3002 | Git服务(Phase 2) |
| Casdoor | 8000 | 认证服务 |
| RDP API | 8080 | 主API服务 |
| PostgreSQL | 5432 | 主数据库 |
| Redis | 6379 | 缓存(Phase 2) |
| MeiliSearch | 7700 | 搜索引擎(Phase 3) |
| Mattermost | 8065 | 即时通讯(Phase 3) |

---

## 4. 功能模块

### 4.1 模块清单与实现阶段

| 模块 | 阶段 | 优先级 | 核心功能 |
|------|------|--------|----------|
| 门户界面 (Portal) | Phase 1 | MUST | 部门首页、个人工作台、快捷导航 |
| 用户管理 (User) | Phase 1 | MUST | 用户账户、组织架构、GitHub风格个人主页 |
| 项目管理 (PM) | Phase 1 | MUST | 项目创建向导、甘特图、项目文件管理 |
| 项目开发 (Dev) | Phase 2 | MUST | 流程执行、交付物管理、本地软件集成 |
| 产品货架 (Shelf) | Phase 2 | MUST | 产品分类浏览、选用、版本管理 |
| 技术货架 (Tech) | Phase 2 | MUST | 技术树展示、技术详情 |
| 知识库 (KB) | Phase 3 | MUST | 知识分类、全文搜索、Obsidian集成 |
| 技术论坛 (Forum) | Phase 3 | SHOULD | 帖子发布、回复、精华帖 |
| 即时通信 (IM) | Phase 3 | SHOULD | Mattermost集成(可选内置通知) |

### 4.2 核心业务流程

1. **项目创建流程**: 五步向导（信息录入→类别选择→流程绑定→团队分配→计划确认）
2. **项目开发流程**: 接收任务→查看活动定义→下载模板→本地开发→提交交付物→Git版本提交→交付物审核→完成活动→触发下一活动
3. **产品上架流程**: 项目完成→提交上架申请→技术评审→质量评审→审批上架→产品展示

---

## 5. 编码规范

### 5.1 通用规范

| 规范项 | 要求 |
|--------|------|
| **语言** | 代码注释和变量名使用英文；UI文案使用中文（i18n机制） |
| **API风格** | RESTful API；路径小写+连字符；版本前缀 `/api/v1/` |
| **错误处理** | 统一错误响应格式 `{"code": int, "message": string, "data": null}` |
| **认证方式** | JWT Bearer Token；Access Token有效期2小时；Refresh Token 7天 |
| **分页** | 统一分页参数 `?page=1&page_size=20`；响应含 total, page, page_size |
| **时间格式** | ISO 8601（UTC）：`2026-02-20T13:00:00Z` |
| **ID生成** | 雪花算法或ULID，不使用自增ID |

### 5.2 前端规范

- 状态管理：使用 Zustand 或 React Context，避免 Redux
- CSS方案：Tailwind CSS + Ant Design 主题定制
- 组件命名：PascalCase
- 变量/函数：camelCase
- 常量：UPPER_SNAKE_CASE

### 5.3 后端规范

- 包命名：小写，使用下划线分隔
- 接口命名：动词+名词，如 `GetUserList`
- 错误码：HTTP状态码 + 业务错误码（6xxx系列）
- 日志：使用 Zap 结构化日志

### 5.4 项目结构规范

```
rdp/                              # Monorepo根目录
├── apps/
│   ├── web/                      # 主前端应用(Shell)
│   ├── desktop/                  # 桌面辅助程序
│   └── modules/                  # 子模块前端
│       ├── user/
│       ├── project/
│       ├── dev/
│       ├── shelf/
│       ├── tech-shelf/
│       ├── knowledge/
│       └── forum/
├── services/
│   ├── api-gateway/              # API网关配置(Nginx)
│   ├── core/                     # 核心业务服务(Go)
│   ├── workflow/                 # 流程引擎服务
│   ├── file-service/             # 文件管理服务
│   ├── search-service/           # 搜索索引服务
│   └── notification/             # 通知服务
├── packages/
│   ├── shared-types/             # 共享TypeScript类型
│   ├── shared-utils/             # 共享工具函数
│   └── ui-components/            # 共享UI组件
├── database/
│   ├── migrations/               # 数据库迁移脚本
│   └── seeds/                    # 初始数据种子
├── deploy/
│   ├── systemd/                  # systemd服务配置
│   ├── nginx/                    # Nginx配置
│   └── scripts/                  # 部署脚本
└── docs/                         # 项目文档
```

---

## 6. API 设计规范

### 6.1 统一响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": 1704067200000,
  "requestId": "req_xxxxxxxx"
}
```

### 6.2 分页响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "pages": 5
  }
}
```

### 6.3 错误码定义

| 错误码 | 错误类型 | 说明 |
|--------|----------|------|
| 200 | SUCCESS | 成功 |
| 400 | BAD_REQUEST | 请求参数错误 |
| 401 | UNAUTHORIZED | 未认证 |
| 403 | FORBIDDEN | 无权限 |
| 404 | NOT_FOUND | 资源不存在 |
| 500 | INTERNAL_ERROR | 服务器内部错误 |

**业务错误码 (6xxx系列):**

| 错误码 | 说明 |
|--------|------|
| 6001 | 用户名已存在 |
| 6002 | 用户不存在 |
| 6003 | 密码错误 |
| 6101 | 项目编号已存在 |
| 6102 | 项目不存在 |
| 6103 | 项目状态不允许操作 |

---

## 7. 部署架构

### 7.1 部署方式

**重要约束**: 根据原始需求，**不使用 Docker 和 Kubernetes**，采用 **systemd裸机部署**。

部署目录结构:
```
/opt/rdp/
├── bin/                    # 可执行文件
├── config/                 # 配置文件
├── data/                   # 数据目录
│   ├── postgresql/
│   ├── redis/
│   ├── minio/
│   ├── gitea/
│   ├── mattermost/
│   └── meilisearch/
├── logs/                   # 日志目录
└── scripts/                # 运维脚本
```

### 7.2 核心服务清单

| 服务 | 类型 | 部署方式 |
|------|------|----------|
| PostgreSQL | 系统服务 | apt安装 + 数据目录迁移 |
| Redis | 系统服务 | apt安装 + 配置自定义 |
| RDP API | 自定义服务 | 二进制 + systemd |
| RDP Portal | Nginx托管 | 静态文件 |
| Casdoor | 自定义服务 | 二进制 + systemd |
| Gitea | 自定义服务 | 二进制 + systemd |
| Mattermost | 自定义服务 | 二进制 + systemd |
| MinIO | 自定义服务 | 二进制 + systemd |
| MeiliSearch | 自定义服务 | 二进制 + systemd |
| Nginx | 系统服务 | apt安装 |

### 7.3 运维命令

```bash
# 查看所有服务状态
systemctl status postgresql redis nginx rdp-api casdoor gitea mattermost meilisearch

# 查看日志
journalctl -u rdp-api -f                    # 实时查看API日志
journalctl -u casdoor --since "1 hour ago"  # 查看最近1小时日志

# 重启服务
sudo systemctl restart rdp-api

# 查看服务详情
sudo systemctl cat rdp-api

# 健康检查
curl http://localhost:8080/api/v1/health
curl http://localhost:8000/api/health

# 备份
/opt/rdp/scripts/backup.sh
```

---

## 8. AI Agent开发路线图（V1.2）

> **⚠️ 重要声明：本项目采用AI Agent集群自主开发模式**
> 
> 本项目区别于传统软件开发，采用**AI Agent自主开发**方式：
> - 各功能模块由专门的AI Agent负责（如PortalAgent、UserAgent等）
> - PM-Agent（项目经理Agent）负责全局质量把控和进度协调
> - 分期仅表示开发先后顺序，不设定具体时间限制
> - 只有通过PM-Agent质量审查的代码才能进入下一阶段

### 8.1 四期规划（AI Agent模式）

| 阶段 | 执行顺序 | 负责AI Agent | 交付内容 | 里程碑 |
|------|----------|-------------|----------|--------|
| **一期：基础骨架** | 第1批 | PortalAgent、UserAgent、ProjectAgent、InfraAgent | 门户框架、用户管理、RBAC权限、项目管理基础CRUD、内置通知、数据库设计、systemd部署 | 可登录使用的基础平台 |
| **二期：核心业务** | 第2批 | WorkflowAgent、DevAgent、ShelfAgent、DesktopAgent | 流程引擎(状态机)、项目开发模块、甘特图、产品/技术货架、Gitea集成、本地软件联动 | 研发流程上线 |
| **三期：知识智能** | 第3批 | KnowledgeAgent、SearchAgent、ForumAgent | 知识库模块、Obsidian/Zotero集成、全文搜索(MeiliSearch)、技术论坛、Mattermost集成(可选) | 知识管理上线 |
| **四期：优化完善** | 第4批 | AnalyticsAgent、各功能Agent | 数据分析仪表盘、流程优化、性能优化、用户体验提升 | 全功能稳定版本 |

### 8.2 AI Agent团队配置

| AI Agent角色 | 负责模块 | 技术栈 | 核心职责 |
|--------------|----------|--------|----------|
| **PortalAgent** | 门户界面 | React + TS + Ant Design | 前端应用开发 |
| **UserAgent** | 用户管理 | Go + Casdoor/Casbin | 认证授权服务 |
| **ProjectAgent** | 项目管理 | Go + Gitea API | 项目服务、Git集成 |
| **WorkflowAgent** | 流程引擎 | Go + 状态机 | 流程定义、活动流转 |
| **DevAgent** | 项目开发 | Go + rdp协议 | 开发模块、本地软件联动 |
| **ShelfAgent** | 产品/技术货架 | Go + PostgreSQL | 货架服务、选用管理 |
| **KnowledgeAgent** | 知识库 | Go + 文件同步 | 知识服务、Obsidian集成 |
| **SearchAgent** | 搜索服务 | Go + MeiliSearch | 全文搜索、索引管理 |
| **ForumAgent** | 技术论坛 | Go + React | 论坛服务、帖子管理 |
| **InfraAgent** | 基础设施 | Shell + systemd | 部署脚本、配置管理 |
| **DesktopAgent** | 桌面辅助程序 | Electron/Tauri | rdp协议、Git自动提交 |
| **AnalyticsAgent** | 数据分析 | Go + ECharts | 仪表盘、报表服务 |
| **PM-Agent** | **项目经理** | **全栈审查** | **代码审查、架构一致性、进度协调、质量把关** |

### 8.3 PM-Agent（项目经理Agent）核心职责

PM-Agent作为项目质量把控核心，负责：
1. **代码质量审查**：审查各Agent代码是否符合编码规范
2. **架构一致性检查**：确保接口设计一致、数据模型统一
3. **进度协调**：根据依赖关系决定开发顺序
4. **冲突解决**：决策不同Agent的实现方案选择
5. **集成测试**：组织阶段集成测试，确保模块协作正常
6. **质量门禁**：只有通过审查的代码才能进入下一阶段

---

## 9. 关键术语

| 术语 | 定义 |
|------|------|
| **IPD** | 集成产品开发（Integrated Product Development） |
| **DCP** | 决策检查点（Decision Check Point） |
| **L1-L4流程** | 流程分级，L1为最高层级，L4为具体活动层级 |
| **单机** | 独立完整的产品单元 |
| **模块** | 可复用的功能组件 |
| **货架** | 已成熟可复用的产品/技术展示平台 |
| **TRL** | 技术成熟度等级（Technology Readiness Level）1-9级 |
| **交付物** | 项目活动中产生的文档、图纸、代码等成果 |
| **活动** | 流程中的最小执行单元 |
| **AI Agent** | 人工智能代理，本项目中负责特定模块开发的智能程序 |
| **PM-Agent** | 项目经理Agent，负责全局质量把控和进度协调的AI Agent |
| **功能Agent** | 负责具体功能模块开发的AI Agent（如PortalAgent、UserAgent等） |
| **AI Agent开发模式** | 本项目采用的开发方式，各模块由专门AI Agent负责，PM-Agent统一协调 |

---

## 10. 参考资源

### 10.1 开源集成项目

| 项目 | 用途 | GitHub |
|------|------|--------|
| Casdoor | 统一认证/RBAC | https://github.com/casdoor/casdoor |
| Casbin | 权限管理 | https://github.com/casbin/casbin |
| Gitea | Git版本管理 | https://github.com/go-gitea/gitea |
| Mattermost | 即时通讯 | https://github.com/mattermost/mattermost |
| MeiliSearch | 全文搜索 | https://github.com/meilisearch/meilisearch |
| MinIO | 对象存储 | https://github.com/minio/minio |
| Wiki.js | 知识库参考 | https://github.com/requarks/wiki |

### 10.2 内部文档导航

| 文档 | 路径 | 说明 |
|------|------|------|
| 需求文档 | `方案/01_需求文档.md` | 需求分析、功能清单、验收标准 |
| 详细实施方案 | `方案/02_详细实施方案.md` | 架构设计、技术选型、实施计划 |
| 需求规格说明书 | `方案/03_需求规格说明书.md` | 功能规格、接口定义、质量要求 |
| 技术架构分析 | `方案/技术架构分析报告.md` | 技术分析、架构设计 |
| HTML/PDF版本 | `方案/pdf-html/` | 文档的HTML和PDF导出版本 |

---

## 11. 注意事项

1. **部署约束**: 严禁使用 Docker 或 Kubernetes，必须使用 systemd 裸机部署
2. **离线环境**: 系统部署在离线局域网，所有依赖需预下载
3. **Git集成**: 文件版本管理必须集成 Git，使用 Gitea 作为 Git 服务
4. **IPD规范**: 项目管理需符合 IPD 方法论，DCP 决策点不可跳过
5. **本地软件**: 支持调用本地专业软件（Altium Designer、Obsidian 等）

---

*微波研发部门研发管理平台 — 项目文档 V1.0*  
*© 2026 微波研发部门 | 内部文档*
