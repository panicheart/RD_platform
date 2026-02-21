# 研发管理平台文档

微波工程部研发管理平台 (RDP) 项目文档集。

## 文档列表

| 序号 | 文档 | 大小 | 说明 |
|------|------|------|------|
| 01 | [implementation-plan.md](01-implementation-plan.md) | 75KB | 详细实施方案（架构设计、技术选型、数据库、模块方案、开源集成） |
| 02 | [requirements-specification.md](02-requirements-specification.md) | 29KB | 需求规格说明书（功能规格、API规范、验收标准） |
| 03 | [deployment-correction.md](03-deployment-correction.md) | 26KB | 部署方案修正（systemd裸机部署完整指南） |
| 00 | **[index.html](index.html)** | **HTML** | **📱 文档中心主页（推荐入口）** |
| 01 | [01-implementation-plan.html](01-implementation-plan.html) | HTML | 详细实施方案（浏览器查看格式） |
| 02 | [02-requirements-specification.html](02-requirements-specification.html) | HTML | 需求规格说明书（浏览器查看格式） |
| — | [specification.html](specification.html) | HTML | 技术规格说明书（浏览器查看格式） |

## 技术栈

- **前端**: React + Vite + TypeScript + Tailwind CSS
- **后端**: Go / Gin
- **数据库**: PostgreSQL + Redis
- **搜索**: MeiliSearch
- **存储**: MinIO
- **部署**: **裸机部署（systemd服务）** — 详见 [03-deployment-correction.md](03-deployment-correction.md)

## 开源集成

- [Casdoor](https://github.com/casdoor/casdoor) — 统一认证
- [Casbin](https://github.com/casbin/casbin) — RBAC 权限
- [Gitea](https://github.com/go-gitea/gitea) — Git 版本管理
- [Mattermost](https://github.com/mattermost/mattermost) — 即时通讯
- [Wiki.js](https://github.com/requarks/wiki) — 知识库参考

## 文档内容概览

### 01-implementation-plan.md (1869行)
- 系统整体架构（分层架构、模块依赖）
- 技术栈选型方案（前后端、数据库、基础设施）
- 数据库架构设计（ER图、表结构、分库分表策略）
- 模块间接口与交互设计（RESTful API、消息队列）
- 核心业务流程图（项目创建、开发执行、工作流审批等）
- 开源项目集成方案（Casdoor、Mattermost、Gitea等）
- **补充章节**：详细开源集成方案（RuoYi、OpenProject、Flowable、Wiki.js、SSO）
- 部署架构设计（systemd服务、Nginx配置）
- 实施路线图（四期规划、团队分工、风险应对）

### 02-requirements-specification.md (1074行)
- 系统约束与技术规范
- 数据模型规范（实体定义、字段说明）
- 功能规格（门户、用户、项目、开发、货架、知识库、论坛、IM）
- 非功能规格（性能、安全、可用性、兼容性）
- API设计规范（错误码、分页、认证）
- 验收测试矩阵

### 03-deployment-correction.md (850行)
- 部署方案变更说明（Docker→systemd）
- 裸机部署架构图
- 服务部署清单（PostgreSQL、Redis、各应用服务）
- Nginx统一配置
- 运维脚本（安装、备份、健康检查）
- 端口分配表

## 🌐 浏览器查看（推荐）

所有文档已提供现代化的 HTML 版本，可直接在浏览器中查看：

1. **打开 `index.html`** — 文档中心主页，包含：
   - 文档导航卡片
   - 功能模块概览
   - 技术栈展示
   - 开源项目链接

2. **各文档页面特点**：
   - 左侧固定导航栏，快速跳转到各章节
   - 卡片式布局，视觉层次清晰
   - 响应式设计，支持移动端
   - 平滑滚动和交互动效

## 重要说明

⚠️ **部署方式**：所有文档已统一更新为 **systemd裸机部署** 方案

- **Markdown文件**：包含完整详细内容，供开发参考
- **HTML文件**：供浏览器查看格式用（基于 Tailwind CSS + Font Awesome）

---

*微波研发部门研发管理平台 — 文档集 V1.0*  
*© 2026 微波研发部门 | 内部文档*
