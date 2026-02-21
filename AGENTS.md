# 微波工程部研发管理平台 (RDP) - 项目文档

> **文档用途**: 本文档供 AI 编程助手参考，用于理解项目架构、技术栈和开发规范。
> **最后更新**: 2026-02-21

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
| **文档位置** | `方案3/` 目录为最新完整方案 |

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
├── 方案1/                       # 方案1：早期需求文档
│   ├── 需求文档.md
│   ├── 需求规格说明书.md
│   └── 实施方案.md
├── 方案2/                       # 方案2：中间版本文档
│   ├── 01-requirements-analysis.html
│   ├── 02-implementation-plan.html
│   └── 03-requirements-specification.html
└── 方案3/                       # 方案3：当前最终方案 ⭐
    ├── README.md                # 文档导航入口
    ├── index.html               # HTML文档主页
    ├── 01-implementation-plan.md    # 详细实施方案（1869行）
    ├── 02-requirements-specification.md  # 需求规格说明书（1074行）
    └── 03-deployment-correction.md   # 部署方案修正
```

**注意**: 方案3为当前采用方案，所有实现应以此为准。方案1和方案2为历史参考。

---

## 2. 技术栈

### 2.1 核心技术约束

| 层次 | 技术 | 版本 | 约束 |
|------|------|------|------|
| **前端框架** | React + TypeScript + Vite | React 18.x, TS 5.x, Vite 5.x | [MUST] 不可替换 |
| **UI库** | Ant Design | 5.x | [MUST] 主UI库 |
| **后端API** | Go (Gin framework) | Go 1.22+, Gin 1.9+ | [MUST] 核心 |
| **数据库** | PostgreSQL | 16.x | [MUST] |
| **缓存** | Redis | 7.x | [MUST] |
| **搜索** | MeiliSearch | 1.x | [MUST] Phase 3 |
| **Git服务** | Gitea | 1.22+ | [MUST] Phase 2 |
| **认证** | Casdoor | Latest | [MUST] Phase 1 |
| **IM** | Mattermost | Team Edition | [MUST] Phase 1 |
| **存储** | MinIO | Latest | [MUST] |
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
| RDP Portal | 3000 | 前端Shell |
| RDP Dev | 3001 | 项目开发模块 |
| Gitea | 3002 | Git服务 |
| Wiki.js | 3003 | 知识库（可选） |
| Casdoor | 8000 | 认证服务 |
| RDP API | 8080 | 主API服务 |
| Mattermost | 8065 | 即时通讯 |
| PostgreSQL | 5432 | 主数据库 |
| Redis | 6379 | 缓存/会话 |
| MinIO API | 9000 | 对象存储API |
| MinIO Console | 9001 | 对象存储控制台 |
| MeiliSearch | 7700 | 搜索引擎 |

---

## 4. 功能模块

### 4.1 模块清单与实现阶段

| 模块 | 阶段 | 优先级 | 核心功能 |
|------|------|--------|----------|
| 门户界面 (Portal) | Phase 1 | MUST | 部门首页、个人工作台、快捷导航 |
| 用户管理 (User) | Phase 1 | MUST | 用户账户、组织架构、GitHub风格个人主页 |
| 项目管理 (PM) | Phase 1/2 | MUST | 项目创建向导、甘特图、自动排程 |
| 项目开发 (Dev) | Phase 2 | MUST | 流程执行、交付物管理、本地软件集成 |
| 产品货架 (Shelf) | Phase 2 | MUST | 产品分类浏览、选用、版本管理 |
| 技术货架 (Tech) | Phase 2 | MUST | 技术树展示、技术详情 |
| 知识库 (KB) | Phase 3 | MUST | 知识分类、全文搜索、Obsidian集成 |
| 技术论坛 (Forum) | Phase 3 | SHOULD | 帖子发布、回复、精华帖 |
| 即时通信 (IM) | Phase 1 | MUST | Mattermost集成 |

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
systemctl status postgresql redis-server nginx rdp-api casdoor gitea mattermost minio meilisearch

# 查看日志
journalctl -u rdp-api -f                    # 实时查看API日志
journalctl -u casdoor --since "1 hour ago"  # 查看最近1小时日志

# 健康检查
/opt/rdp/scripts/health-check.sh

# 备份
/opt/rdp/scripts/backup.sh
```

---

## 8. 开发路线图

### 8.1 四期规划

| 阶段 | 周期 | 交付内容 | 里程碑 |
|------|------|----------|--------|
| **一期：基础骨架** | 3个月 | 门户框架、用户管理、RBAC权限、项目管理基础CRUD、Mattermost集成、数据库设计 | 可登录使用 |
| **二期：核心业务** | 4个月 | 流程引擎、项目开发模块、甘特图组件、产品/技术货架、Gitea集成、桌面辅助程序 | 研发流程上线 |
| **三期：知识智能** | 3个月 | 知识库模块、Obsidian集成、Zotero集成、全文搜索、技术论坛、标签智能关联 | 知识管理上线 |
| **四期：优化完善** | 2个月 | 数据分析仪表盘、MS Project导入导出、辅助程序完善、性能优化、安全加固 | 全功能GA |

### 8.2 团队分工建议

| 团队 | 人数 | 负责模块 | 技术栈 |
|------|------|----------|--------|
| 前端团队A | 2-3人 | 门户Shell + 用户管理 + 项目管理前端 | React + TypeScript + Ant Design |
| 前端团队B | 2-3人 | 项目开发 + 货架模块 + 知识库前端 | React + TypeScript + ECharts |
| 后端团队 | 3-4人 | API服务 + 流程引擎 + 数据库 | Go + PostgreSQL + Redis |
| 集成/DevOps | 1-2人 | 开源组件集成 + 部署 | Nginx + Shell + systemd |
| 桌面端 | 1人 | 辅助程序开发 | Electron/Tauri |

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
| 详细实施方案 | `方案3/01-implementation-plan.md` | 架构设计、技术选型、数据库设计 |
| 需求规格说明书 | `方案3/02-requirements-specification.md` | 功能规格、API规范、验收标准 |
| 部署方案修正 | `方案3/03-deployment-correction.md` | systemd裸机部署完整指南 |
| HTML文档主页 | `方案3/index.html` | 浏览器查看入口 |

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
