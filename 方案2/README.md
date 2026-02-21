# 研发管理平台文档

微波工程部研发管理平台 (RDP) 项目文档集。

## 文档列表

| 文档 | HTML | PDF |
|------|------|-----|
| 需求梳理与补全 | [01-requirements-analysis.html](01-requirements-analysis.html) | [01-需求梳理与补全文档.pdf](01-需求梳理与补全文档.pdf) |
| 详细实施方案 | [02-implementation-plan.html](02-implementation-plan.html) | [02-详细实施方案.pdf](02-详细实施方案.pdf) |
| 需求规格说明书 | [03-requirements-specification.html](03-requirements-specification.html) | [03-需求规格说明书.pdf](03-需求规格说明书.pdf) |

## 技术栈

- **前端**: React + Vite + TypeScript + Tailwind CSS
- **后端**: Go / Gin
- **数据库**: PostgreSQL + Redis
- **搜索**: MeiliSearch
- **存储**: MinIO
- **部署**: Docker Compose（离线局域网）

## 开源集成

- [Casdoor](https://github.com/casdoor/casdoor) — 统一认证
- [Casbin](https://github.com/casbin/casbin) — RBAC 权限
- [Gitea](https://github.com/go-gitea/gitea) — Git 版本管理
- [Mattermost](https://github.com/mattermost/mattermost) — 即时通讯
- [Wiki.js](https://github.com/requarks/wiki) — 知识库参考
