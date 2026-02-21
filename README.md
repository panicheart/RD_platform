# 微波工程部研发管理平台 (RDP)

> 微波系统研发部门的企业级综合管理平台

[![License](https://img.shields.io/badge/license-Internal-blue.svg)]()
[![Version](https://img.shields.io/badge/version-V1.0-green.svg)]()
[![Platform](https://img.shields.io/badge/platform-Offline_LAN-orange.svg)]()

---

## 项目简介

**微波工程部研发管理平台** (R&D Platform, 简称 **RDP**) 是为微波系统研发部门设计的企业级综合管理平台。

采用 **AI Agent 集群自主开发模式**，所有代码由 AI Agent 编写，人类仅做监督和决策。

### 核心目标

| 目标 | 说明 |
|------|------|
| **G1-流程数字化** | 将L1-L4级研发流程固化到系统中，实现流程驱动的项目执行 |
| **G2-知识资产化** | 建立组织级知识库，实现知识的采集、分类、检索和智能推荐 |
| **G3-产品货架化** | 将成熟产品和技术上架管理，支持一键选用和集成 |
| **G4-协同高效化** | 提供即时通信、论坛、项目协同等工具，打通团队沟通壁垒 |
| **G5-管理可视化** | 通过仪表盘、甘特图、统计报表等实现管理决策的数据支撑 |

---

## 技术栈

### 核心技术

| 层次 | 技术 | 版本 |
|------|------|------|
| **前端框架** | React + TypeScript + Vite | React 18.x, TS 5.x, Vite 5.x |
| **UI库** | Ant Design | 5.x |
| **后端框架** | Go (Gin) | Go 1.22+, Gin 1.9+ |
| **数据库** | PostgreSQL | 16.x |
| **认证服务** | Casdoor | Latest |
| **部署方式** | systemd | 裸机部署 |

### 分阶段引入组件

| 组件 | 阶段 | 说明 |
|------|------|------|
| Casdoor | Phase 1 | 认证服务（必需） |
| Redis | Phase 2 | 缓存（可选） |
| Gitea | Phase 2 | Git 版本管理 |
| MeiliSearch | Phase 3 | 全文搜索引擎 |
| Mattermost | Phase 3 | 即时通讯（可选） |

---

## 项目结构

```
RD_platform/
├── apps/                      # 前端应用
│   └── web/                   # React 单体应用
├── services/                  # 后端服务
│   └── api/                   # Go API 服务
├── database/                  # 数据库
│   ├── migrations/             # 迁移脚本
│   └── seeds/                 # 种子数据
├── deploy/                    # 部署配置
│   ├── systemd/               # systemd 服务
│   ├── nginx/                 # Nginx 配置
│   └── scripts/               # 运维脚本
├── config/                    # 配置文件
├── packages/                  # 共享包
├── agents/                    # Agent 任务文档
│   └── tasks/                 # 任务卡片
├── docs/                     # 项目文档
│   ├── 原始需求/              # 原始需求 PDF
│   ├── 01_需求文档.md
│   ├── 02_详细实施方案.md
│   ├── 03_需求规格说明书.md
│   ├── 技术架构分析报告.md
│   └── 任务拆解文档.md
└── README.md                 # 本文件
```

---

## 开发阶段

### Phase 1: 基础骨架
- [x] 门户界面 (PortalAgent)
- [x] 用户管理 (UserAgent)
- [x] 项目管理 (ProjectAgent)
- [x] 安全合规 (SecurityAgent)
- [x] 基础设施 (InfraAgent)

### Phase 2: 核心业务
- [ ] 流程引擎 (WorkflowAgent)
- [ ] 项目开发 (DevAgent)
- [ ] 产品/技术货架 (ShelfAgent)
- [ ] 桌面辅助程序 (DesktopAgent)
- [ ] 质量管理 (QMAgent)

### Phase 3: 知识智能
- [ ] 知识库 (KnowledgeAgent)
- [ ] 搜索服务 (SearchAgent)
- [ ] 技术论坛 (ForumAgent)

### Phase 4: 优化完善
- [ ] 数据分析 (AnalyticsAgent)
- [ ] 运维监控 (MonitorAgent)

---

## 快速开始

### 环境要求

- **操作系统**: Ubuntu Server 22.04 LTS / CentOS 8+ / 麒麟OS
- **硬件**: 8核CPU / 32GB内存 / 1TB SSD / 千兆网络
- **数据库**: PostgreSQL 16.x

### 安装步骤

```bash
# 1. 克隆项目
git clone <repository-url> /opt/rdp

# 2. 运行安装脚本
cd /opt/rdp/deploy/scripts
sudo ./install.sh

# 3. 启动服务
sudo systemctl start rdp-api
sudo systemctl start rdp-casdoor

# 4. 访问系统
# 前端: http://<server-ip>
# API:  http://<server-ip>:8080
# Casdoor: http://<server-ip>:8000
```

---

## 功能模块

| 模块 | 说明 | 状态 |
|------|------|------|
| 门户首页 | 部门简介、新闻公告、荣誉展示 | Phase 1 |
| 个人工作台 | 待办事项、我的项目、消息通知 | Phase 1 |
| 用户管理 | 认证授权、RBAC权限、组织架构 | Phase 1 |
| 项目管理 | 项目创建、甘特图、文件管理 | Phase 1/2 |
| 流程引擎 | 状态机、活动流转、DCP评审 | Phase 2 |
| 产品货架 | 产品浏览、选用购物车、版本管理 | Phase 2 |
| 技术货架 | 技术树、TRL等级、fork创建 | Phase 2 |
| 知识库 | 分类管理、Obsidian同步、Zotero集成 | Phase 3 |
| 全文搜索 | 中文分词、跨模块搜索、结果高亮 | Phase 3 |
| 技术论坛 | 板块管理、发帖回帖、@通知 | Phase 3 |
| 数据分析 | 仪表盘、统计报表、导出功能 | Phase 4 |
| 运维监控 | 系统监控、APM、日志、告警 | Phase 4 |

---

## 文档导航

| 文档 | 说明 |
|------|------|
| [需求文档](docs/01_需求文档.md) | 需求分析、功能清单、验收标准 |
| [详细实施方案](docs/02_详细实施方案.md) | 架构设计、技术选型、实施计划 |
| [需求规格说明书](docs/03_需求规格说明书.md) | 功能规格、接口定义、质量要求 |
| [技术架构分析报告](docs/技术架构分析报告.md) | 技术分析、架构设计 |
| [任务拆解文档](docs/任务拆解文档.md) | Agent 任务拆分、依赖关系 |
| [Agent 任务卡片](agents/tasks/) | 各 Agent 详细任务输入输出 |

---

## 编码规范

| 规范项 | 要求 |
|--------|------|
| **代码注释** | 英文 |
| **变量命名** | 英文 |
| **UI文案** | 中文 |
| **API路径** | `/api/v1/{module}` 小写+连字符 |
| **错误响应** | `{"code": int, "message": string, "data": null}` |
| **时间格式** | ISO 8601 UTC |
| **ID生成** | ULID 或雪花算法 |

---

## 相关资源

- [Casdoor](https://github.com/casdoor/casdoor) - 身份认证
- [Casbin](https://github.com/casbin/casbin) - 权限管理
- [Gitea](https://github.com/go-gitea/gitea) - Git 服务
- [MeiliSearch](https://github.com/meilisearch/meilisearch) - 全文搜索
- [Ant Design](https://ant.design/) - UI 组件库

---

## 团队

| 角色 | 说明 |
|------|------|
| **Architect Agent** | 架构设计、技术选型 |
| **PM-Agent** | 项目协调、进度跟踪 |
| **Reviewer Agent** | 代码审查、质量把关 |
| **Feature Agents** | 各功能模块开发 |

---

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| V1.0 | 2026-02-22 | 项目初始化，任务拆解完成 |

---

*微波研发部门研发管理平台 - AI Agent 开发模式*
*© 2026 微波研发部门 | 内部文档*
