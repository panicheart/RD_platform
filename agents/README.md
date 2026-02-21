# 微波研发部门研发管理平台 (RDP)
# Agent 开发任务文档导航

> **文档用途**: 本目录包含所有Agent的任务卡片，明确每个Agent的输入、输出、检查者和具体子任务，供AI Agent集群协同开发使用。

---

## 文档导航

### 核心文档

| 文档 | 说明 |
|------|------|
| **[任务总览表](./tasks/agent_overview.md)** | 所有Agent任务汇总、依赖关系、数据流、验收流程 |

### 分期任务卡片

| Phase | 文档 | 任务数量 |
|-------|------|----------|
| Phase 1: 基础骨架 | [tasks/phase1_tasks.md](./tasks/phase1_tasks.md) | 20个任务 |
| Phase 2: 核心业务 | [tasks/phase2_tasks.md](./tasks/phase2_tasks.md) | 23个任务 |
| Phase 3: 知识智能 | [tasks/phase3_tasks.md](./tasks/phase3_tasks.md) | 12个任务 |
| Phase 4: 优化完善 | [tasks/phase4_tasks.md](./tasks/phase4_tasks.md) | 10个任务 |

---

## 任务统计

| Phase | Agent数量 | 任务数量 | 并行度 |
|-------|-----------|----------|--------|
| Phase 1 | 5 | 20 | 完全并行 |
| Phase 2 | 5 | 23 | 分层并行 |
| Phase 3 | 3 | 12 | 完全并行 |
| Phase 4 | 2 | 10 | 完全并行 |
| **总计** | **15** | **65** | - |

---

## 快速开始

### 1. 确认开发范围

选择要启动的Phase：
- **Phase 1**: 基础骨架（5个Agent并行）
- **Phase 2**: 核心业务（5个Agent分层）
- **Phase 3**: 知识智能（3个Agent并行）
- **Phase 4**: 优化完善（2个Agent并行）

### 2. 阅读任务卡片

启动Agent前，需要Agent阅读对应的任务卡片：
- 输入：明确需要依赖哪些模块的数据
- 输出：明确需要交付哪些文件
- 检查者：明确由谁进行审查
- 子任务：明确具体实现内容

### 3. 执行开发

按照任务卡片中的子任务清单逐项实现，每个子任务完成后进行自审查（L1）。

### 4. 提交审查

开发完成后，将交付物提交给Reviewer Agent进行L2审查。

---

## 项目目录结构

```
RD_platform/
├── apps/
│   └── web/                    # React前端
│       ├── src/
│       │   ├── pages/          # 页面组件
│       │   ├── components/     # 公共组件
│       │   ├── services/       # API服务
│       │   ├── stores/         # 状态管理
│       │   ├── types/          # 类型定义
│       │   └── utils/          # 工具函数
│       └── dist/               # 构建产物
├── services/
│   └── api/                    # Go后端
│       ├── handlers/           # HTTP处理器
│       ├── services/           # 业务逻辑
│       ├── models/             # 数据模型
│       ├── middleware/         # 中间件
│       ├── clients/            # 外部服务客户端
│       └── utils/              # 工具函数
├── database/
│   ├── migrations/             # 数据库迁移
│   ├── seeds/                 # 种子数据
│   └── schema/                # 架构定义
├── deploy/
│   ├── systemd/               # systemd服务配置
│   ├── nginx/                 # Nginx配置
│   └── scripts/               # 部署脚本
├── config/                     # 配置文件
├── agents/
│   ├── instructions/          # Agent指令（待创建）
│   ├── outputs/               # Agent交付物（待创建）
│   └── tasks/                 # 任务卡片
│       ├── agent_overview.md  # 任务总览表
│       ├── phase1_tasks.md   # Phase 1任务
│       ├── phase2_tasks.md   # Phase 2任务
│       ├── phase3_tasks.md   # Phase 3任务
│       └── phase4_tasks.md   # Phase 4任务
└── 方案/                       # 项目文档
```

---

## 编码规范（强制）

| 规范项 | 要求 |
|--------|------|
| 代码注释 | 英文 |
| 变量命名 | 英文 |
| UI文案 | 中文 |
| API路径 | `/api/v1/{module}` 小写+连字符 |
| 错误响应 | `{"code": int, "message": string, "data": null}` |
| 时间格式 | ISO 8601 UTC |
| ID生成 | ULID或雪花算法 |

---

## 联系方式

如有疑问，请查阅：
- 需求文档: `方案/01_需求文档.md`
- 实施方案: `方案/02_详细实施方案.md`
- 需求规格: `方案/03_需求规格说明书.md`
- 开发指南: `AGENTS.md`

---

*文档版本: V1.0*
*最后更新: 2026-02-22*
