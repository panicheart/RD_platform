# 研发管理平台文档

微波工程部研发管理平台 (RDP) 项目文档集。

## 文档列表

| 文档 | HTML | PDF |
|------|------|-----|
| 需求文档 | [01-requirements-analysis.html](pdf-html/01-requirements-analysis.html) | [01-需求梳理与补全文档.pdf](pdf-html/01-需求梳理与补全文档.pdf) |
| 详细实施方案 | [02-implementation-plan.html](pdf-html/02-implementation-plan.html) | [02-详细实施方案.pdf](pdf-html/02-详细实施方案.pdf) |
| 需求规格说明书 | [03-requirements-specification.html](pdf-html/03-requirements-specification.html) | [03-需求规格说明书.pdf](pdf-html/03-需求规格说明书.pdf) |
| 技术架构分析 | - | - |

## 技术栈

### 核心技术

| 层次 | 技术 | 说明 |
|------|------|------|
| **前端** | React 18 + TypeScript + Vite | 单体应用（非微前端） |
| **UI库** | Ant Design 5 | 企业级组件库 |
| **后端** | Go (Gin) | 纯 Go 技术栈，单二进制部署 |
| **数据库** | PostgreSQL 16 | 系统包安装，含全文搜索 |
| **部署** | systemd 裸机部署 | 离线局域网 |

### 分阶段引入组件

| 组件 | 引入阶段 | 说明 |
|------|---------|------|
| Casdoor | Phase 1 | 认证服务（必需） |
| Redis | Phase 2 | 缓存（可选，可用内存缓存替代） |
| Gitea | Phase 2 | Git 版本管理 |
| MeiliSearch | Phase 3 | 全文搜索引擎 |
| Mattermost | Phase 3 | 即时通讯（可选） |

### 不使用的技术

- ❌ Docker / Kubernetes
- ❌ 微前端架构
- ❌ Node.js 后端
- ❌ MinIO（文件量小，本地存储足够）

## 部署架构

```
┌─────────────────────────────────────────────────────────────┐
│                         Nginx                               │
│           (反向代理 + SSL + 静态资源托管)                      │
└─────────────┬───────────────────────────────┬───────────────┘
              │                               │
    ┌─────────▼──────────┐      ┌─────────────▼──────────────┐
    │   React 前端        │      │      Go API 服务          │
    │   (单体应用)       │      │      (systemd)            │
    └────────────────────┘      └─────────────┬──────────────┘
                                              │
    ┌─────────────────────────────────────────▼──────────────┐
    │              数据层 (PostgreSQL + 本地文件)              │
    └─────────────────────────────────────────────────────────┘

独立服务 (systemd):
├── Casdoor (认证)     :8000 (Phase 1)
├── Gitea (Git)        :3002 (Phase 2)
├── MeiliSearch        :7700 (Phase 3)
└── Mattermost         :8065 (Phase 3, 可选)
```

## 快速开始

```bash
# 1. 安装依赖（PostgreSQL）
sudo apt update
sudo apt install -y postgresql nginx

# 2. 下载离线安装包
cd /opt/rdp

# 3. 运行安装脚本
sudo ./install.sh

# 4. 启动服务
sudo systemctl start rdp-api
sudo systemctl start casdoor

# 5. 查看状态
sudo systemctl status rdp-api
```

## 服务管理

```bash
# 查看所有 RDP 服务
systemctl list-units --type=service --state=running | grep rdp

# 重启 API 服务
sudo systemctl restart rdp-api

# 查看日志
journalctl -u rdp-api -f

# 健康检查
curl http://localhost:8080/api/v1/health
```

## 目录结构

```
/opt/rdp/
├── bin/                    # 可执行文件
│   ├── rdp-api            # 主 API 服务（Go 二进制）
│   └── casdoor            # 认证服务（Go 二进制）
├── config/                 # 配置文件
│   ├── rdp-api.yaml
│   └── nginx.conf
├── data/                   # 数据目录
│   ├── postgresql/        # 数据库数据
│   ├── files/             # 上传文件
│   │   ├── projects/      # 项目文件
│   │   └── knowledge/     # 知识库文件
│   └── gitea/             # Git 仓库(Phase 2)
├── logs/                   # 日志目录
├── scripts/                # 运维脚本
│   ├── install.sh         # 安装脚本
│   ├── backup.sh          # 备份脚本
│   └── health-check.sh    # 健康检查
└── web/                   # 前端静态文件（React 构建产物）
```

## 项目结构

```
rdp/
├── apps/
│   └── web/               # 单体前端应用（React + Vite）
├── services/
│   └── api/               # 后端 API 服务（Go + Gin）
├── packages/
│   ├── shared-types/      # 共享 TypeScript 类型
│   └── shared-utils/      # 共享工具函数
├── database/
│   ├── migrations/        # 数据库迁移脚本
│   └── seeds/             # 初始数据
└── deploy/
    ├── systemd/           # systemd 服务配置
    ├── nginx/             # Nginx 配置
    └── scripts/           # 部署脚本
```

## AI Agent开发团队配置

> **⚠️ 开发模式**：本项目采用**AI Agent集群自主开发模式**，非人工开发。

### 功能Agent角色

| AI Agent | 负责模块 | 技术栈 | 阶段 |
|----------|----------|--------|------|
| **PortalAgent** | 门户界面（React前端） | React + TS + Ant Design | Phase 1/4 |
| **UserAgent** | 用户管理 + 权限 | Go + Casdoor/Casbin | Phase 1 |
| **ProjectAgent** | 项目管理 + Gitea集成 | Go + Gitea API | Phase 1/2 |
| **WorkflowAgent** | 流程引擎 + DCP评审 | Go + 状态机 | Phase 2 |
| **DevAgent** | 项目开发 + 本地软件联动 | Go + rdp协议 | Phase 2 |
| **ShelfAgent** | 产品/技术货架 | Go + PostgreSQL | Phase 2 |
| **KnowledgeAgent** | 知识库 + Obsidian集成 | Go + 文件同步 | Phase 3 |
| **SearchAgent** | 全文搜索 | Go + MeiliSearch | Phase 3 |
| **ForumAgent** | 技术论坛 | Go + React | Phase 3 |
| **InfraAgent** | 基础设施 + 部署 | Shell + systemd | Phase 1/3 |
| **DesktopAgent** | 桌面辅助程序 | Electron/Tauri | Phase 2 |
| **AnalyticsAgent** | 数据分析仪表盘 | Go + ECharts | Phase 4 |

### PM-Agent（项目经理Agent）

| 职责 | 说明 |
|------|------|
| **代码审查** | 审查各Agent代码质量、规范符合度 |
| **架构一致性** | 确保接口设计、数据模型统一 |
| **进度协调** | 根据依赖关系协调开发顺序 |
| **质量门禁** | 只有通过审查的代码才能进入下一阶段 |

### 开发模式特点

- **无固定时间**：分期仅表示先后顺序，不设定具体时间
- **质量优先**：PM-Agent根据质量把关结果决定进度
- **并行开发**：无依赖的模块可并行开发
- **迭代交付**：每个阶段交付可运行版本

---

*© 2026 微波研发部门 | 内部文档*  
*版本: V1.2 (AI Agent开发版)*
