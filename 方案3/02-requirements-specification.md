# 微波研发部门研发管理平台需求规格说明书 (SRS)

> **文档编号**: RDP-SRS-2026-001  
> **版本**: V1.0  
> **编制日期**: 2026年2月20日  
> **密级**: 内部公开  
> **目标读者**: AI Agent 实现集群

⚠️ **Agent 实现指引**：本文档是面向 AI Agent 集群的需求规格说明书。每个功能规格包含明确的输入/输出/验收标准/数据模型/API规范，Agent应严格按照规格实现，不得遗漏。对于标注为 **[MUST]** 的条目为强制实现，**[SHOULD]** 为推荐实现，**[MAY]** 为可选实现。

---

## 目录

1. [引言](#1-引言)
2. [系统约束与技术规范](#2-系统约束与技术规范)
3. [数据模型规范](#3-数据模型规范)
4. [功能规格：门户界面 (SRS-PORTAL)](#4-功能规格门户界面-srs-portal)
5. [功能规格：用户管理 (SRS-USER)](#5-功能规格用户管理-srs-user)
6. [功能规格：项目管理 (SRS-PM)](#6-功能规格项目管理-srs-pm)
7. [功能规格：项目开发 (SRS-DEV)](#7-功能规格项目开发-srs-dev)
8. [功能规格：产品货架 (SRS-SHELF)](#8-功能规格产品货架-srs-shelf)
9. [功能规格：技术货架 (SRS-TECH)](#9-功能规格技术货架-srs-tech)
10. [功能规格：知识库 (SRS-KB)](#10-功能规格知识库-srs-kb)
11. [功能规格：技术论坛 (SRS-FORUM)](#11-功能规格技术论坛-srs-forum)
12. [功能规格：即时通信 (SRS-IM)](#12-功能规格即时通信-srs-im)
13. [非功能规格](#13-非功能规格)
14. [API 设计规范](#14-api-设计规范)
15. [验收测试矩阵](#15-验收测试矩阵)

---

## 1. 引言

### 1.1 系统标识

| 项目 | 内容 |
|------|------|
| 系统名称 | 微波研发部门研发管理平台（R&D Platform, 简称RDP） |
| 系统代号 | RDP |
| 版本目标 | V1.0（覆盖一至四期全部功能） |
| 部署环境 | 离线局域网，**systemd裸机部署** ⚠️ 已修正（原Docker Compose） |
| 用户规模 | 30-100人 |

### 1.2 Agent 实现约定

**实现顺序约定**：
1. Phase 1 标记的功能必须在第一批次实现（门户+用户+项目管理基础+IM集成）
2. Phase 2 标记的功能在第二批次实现（流程引擎+项目开发+货架）
3. Phase 3 标记的功能在第三批次实现（知识库+论坛+搜索）
4. Phase 4 标记的功能在第四批次实现（优化+高级特性）

每个功能规格包含：功能ID | 描述 | 输入 | 输出 | 前置条件 | 数据模型 | API端点 | 验收标准 | 阶段标记

### 1.3 术语定义

| 术语 | 定义 |
|------|------|
| IPD | 集成产品开发（Integrated Product Development） |
| DCP | 决策检查点（Decision Check Point） |
| L1-L4流程 | 流程分级，L1为最高层级，L4为具体活动层级 |
| 单机 | 独立完整的产品单元 |
| 模块 | 可复用的功能组件 |
| 货架 | 已成熟可复用的产品/技术展示平台 |
| FPGA | 现场可编程门阵列 |
| ADC | 模数转换器 |
| SDR | 软件定义无线电 |
| RFSoC | 射频片上系统 |
| TRL | 技术成熟度等级（Technology Readiness Level）1-9级 |

---

## 2. 系统约束与技术规范

### 2.1 技术栈约束（强制）

| 层次 | 技术 | 版本约束 | 备注 |
|------|------|----------|------|
| 前端框架 | React + TypeScript + Vite | React 18.x, TS 5.x, Vite 5.x | [MUST] 不可替换 |
| UI库 | Ant Design 5 | 5.x | [MUST] 主UI库 |
| 后端API | Go (Gin framework) | Go 1.22+, Gin 1.9+ | [MUST] 核心 |
| 数据库 | PostgreSQL | 16.x | [MUST] |
| 缓存 | Redis | 7.x | [MUST] |
| 搜索 | MeiliSearch | 1.x | [MUST] Phase 3 |
| Git服务 | Gitea (独立部署,API集成) | 1.22+ | [MUST] Phase 2 |
| 认证 | Casdoor (独立部署) | Latest | [MUST] Phase 1 |
| IM | Mattermost (独立部署) | Latest Team Edition | [MUST] Phase 1 |
| 对象存储 | MinIO | Latest | [MUST] |
| **服务管理** | **systemd + 裸机二进制部署** | — | **[MUST]** ⚠️ 已修正（原Docker） |

### 2.2 编码规范约束

| 规范项 | 要求 |
|--------|------|
| 语言 | [MUST] 代码注释和变量名使用英文；UI文案使用中文（i18n机制） |
| API风格 | [MUST] RESTful API；路径小写+连字符；版本前缀 /api/v1/ |
| 错误处理 | [MUST] 统一错误响应格式 {"code": int, "message": string, "data": null} |
| 认证方式 | [MUST] JWT Bearer Token；Access Token有效期2小时；Refresh Token 7天 |
| 分页 | [MUST] 统一分页参数 ?page=1&page_size=20 ；响应含 total, page, page_size |
| 时间格式 | [MUST] ISO 8601（UTC）： 2026-02-20T13:00:00Z |
| ID生成 | [MUST] 雪花算法或ULID，不使用自增ID |
| 前端状态 | [SHOULD] Zustand 或 React Context；避免Redux |
| CSS方案 | [SHOULD] Tailwind CSS + Ant Design 主题定制 |
| 测试覆盖 | [SHOULD] 后端核心逻辑单元测试覆盖率 ≥ 60% |

### 2.3 项目结构约束

```
rdp/                              # Monorepo根目录
├── apps/
│   ├── web/                      # 主前端应用(Shell)
│   ├── desktop/                  # 桌面辅助程序(Electron/Tauri)
│   └── modules/
│       ├── user/                 # 用户管理子模块前端
│       ├── project/              # 项目管理子模块前端
│       ├── dev/                  # 项目开发子模块前端
│       ├── shelf/                # 产品货架子模块前端
│       ├── tech-shelf/           # 技术货架子模块前端
│       ├── knowledge/            # 知识库子模块前端
│       └── forum/                # 论坛子模块前端
├── services/
│   ├── api-gateway/              # API网关配置(Nginx)
│   ├── core/                     # 核心业务服务(Go)
│   ├── workflow/                 # 流程引擎服务(Go)
│   ├── file-service/             # 文件管理服务(Go)
│   ├── search-service/           # 搜索索引服务(Go)
│   └── notification/             # 通知服务(Go)
├── packages/
│   ├── shared-types/             # 共享TypeScript类型定义
│   ├── shared-utils/             # 共享工具函数
│   └── ui-components/            # 共享UI组件
├── database/
│   ├── migrations/               # 数据库迁移脚本
│   └── seeds/                    # 初始数据种子
├── deploy/
│   ├── systemd/                  # systemd服务配置 ⚠️ 已修正（原docker）
│   ├── nginx/                    # Nginx配置
│   └── scripts/                  # 部署脚本
└── docs/                         # 项目文档
```

---

## 3. 数据模型规范

### 3.1 实体关系图

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│    User     │◄─────►│    Role     │◄─────►│  Permission │
└──────┬──────┘       └─────────────┘       └─────────────┘
       │
       │ 1:N
       ▼
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│Organization │◄─────►│   Project   │◄─────►│   Task      │
└─────────────┘       └──────┬──────┘       └─────────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
              ▼              ▼              ▼
       ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
       │  Deliverable│ │   Product   │ │  Technology │
       └─────────────┘ └─────────────┘ └─────────────┘
                             │
                             ▼
                      ┌─────────────┐
                      │    File     │
                      └─────────────┘
```

### 3.2 核心实体定义

#### 用户实体 (User)

| 字段名 | 数据类型 | 长度 | 是否为空 | 默认值 | 说明 |
|--------|----------|------|----------|--------|------|
| id | BIGINT | - | 否 | 自增 | 主键 |
| username | VARCHAR | 64 | 否 | - | 用户名(唯一) |
| password | VARCHAR | 256 | 否 | - | 加密密码 |
| real_name | VARCHAR | 64 | 否 | - | 真实姓名 |
| email | VARCHAR | 128 | 是 | NULL | 邮箱 |
| phone | VARCHAR | 32 | 是 | NULL | 手机号 |
| avatar | VARCHAR | 512 | 是 | NULL | 头像URL |
| org_id | BIGINT | - | 否 | - | 所属组织ID |
| title_level | TINYINT | - | 否 | 1 | 职称等级(1-5) |
| status | TINYINT | - | 否 | 1 | 状态(0-禁用,1-启用,2-离职) |
| last_login_time | TIMESTAMP | - | 是 | NULL | 最后登录时间 |
| last_login_ip | VARCHAR | 64 | 是 | NULL | 最后登录IP |
| online_duration | BIGINT | - | 否 | 0 | 在线时长(分钟) |
| created_at | TIMESTAMP | - | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | - | 否 | CURRENT_TIMESTAMP | 更新时间 |

#### 项目实体 (Project)

| 字段名 | 数据类型 | 长度 | 是否为空 | 默认值 | 说明 |
|--------|----------|------|----------|--------|------|
| id | BIGINT | - | 否 | 自增 | 主键 |
| project_code | VARCHAR | 64 | 否 | - | 项目编号(唯一) |
| project_name | VARCHAR | 256 | 否 | - | 项目名称 |
| project_type | TINYINT | - | 否 | - | 类型(1-单机,2-模块,3-软件,4-技术,5-流程,6-知识库,7-产品上架) |
| category_id | BIGINT | - | 否 | - | 类别ID |
| status | TINYINT | - | 否 | 1 | 状态(1-草稿,2-待启动,3-执行中,4-评审中,5-已暂停,6-已完成,7-已归档) |
| priority | TINYINT | - | 否 | 2 | 优先级(1-高,2-中,3-低) |
| manager_id | BIGINT | - | 否 | - | 项目经理ID |
| start_date | DATE | - | 是 | NULL | 计划开始日期 |
| end_date | DATE | - | 是 | NULL | 计划结束日期 |
| actual_start_date | DATE | - | 是 | NULL | 实际开始日期 |
| actual_end_date | DATE | - | 是 | NULL | 实际结束日期 |
| progress | DECIMAL | 5,2 | 否 | 0.00 | 进度百分比 |
| budget | DECIMAL | 15,2 | 是 | NULL | 预算金额 |
| description | TEXT | - | 是 | NULL | 项目描述 |
| folder_path | VARCHAR | 512 | 是 | NULL | 项目文件夹路径 |
| workflow_id | BIGINT | - | 是 | NULL | 绑定流程ID |
| parent_id | BIGINT | - | 是 | NULL | 父项目ID |
| created_at | TIMESTAMP | - | 否 | CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | - | 否 | CURRENT_TIMESTAMP | 更新时间 |

---

## 4. 功能规格：门户界面 (SRS-PORTAL)

### SRS-PORTAL-001 [MUST] [Phase 1] 部门首页展示

**功能描述**  
展示部门介绍、新闻动态、通知公告、荣誉展示等信息。

**输入**  
无

**输出**  
渲染部门首页，包含：
- 部门Logo和简介
- 轮播新闻（支持自动播放和手动切换）
- 通知公告列表
- 荣誉墙展示
- 快速入口导航

**前置条件**  
无

**数据模型**  
- `portal_news` 新闻表
- `portal_notice` 通知表
- `portal_honor` 荣誉表

**API端点**  
```
GET /api/v1/portal/news?limit=5
GET /api/v1/portal/notices?limit=10
GET /api/v1/portal/honors
```

**验收标准**  
- [ ] 页面加载时间 < 3秒
- [ ] 支持响应式布局（适配桌面、平板）
- [ ] 新闻轮播支持自动播放（间隔5秒）和手动切换
- [ ] 通知公告按时间倒序排列
- [ ] 荣誉墙展示部门近年获得的奖项

---

### SRS-PORTAL-002 [MUST] [Phase 1] 个人工作台

**功能描述**  
类似飞书工作台，集成各模块快捷入口，展示待办事项、消息通知、数据看板。

**输入**  
用户登录态

**输出**  
个人工作台页面，包含：
- 快捷入口卡片（可拖拽排序）
- 待办事项列表
- 消息通知中心
- 个人数据看板（项目进度、任务完成率等）
- 日历视图（项目里程碑、会议安排）

**前置条件**  
用户已登录

**数据模型**  
- `user_workbench_config` 工作台配置表
- `sys_notification` 通知表

**API端点**  
```
GET /api/v1/workbench/config
PUT /api/v1/workbench/config
GET /api/v1/workbench/todos
GET /api/v1/workbench/notifications
GET /api/v1/workbench/dashboard
```

**验收标准**  
- [ ] 工作台布局支持拖拽调整
- [ ] 快捷入口支持自定义添加/删除
- [ ] 待办事项实时更新（WebSocket推送）
- [ ] 消息通知支持一键已读
- [ ] 数据看板图表可交互

---

## 5. 功能规格：用户管理 (SRS-USER)

### SRS-USER-001 [MUST] [Phase 1] 用户登录/登出

**功能描述**  
用户通过用户名密码登录系统，集成Casdoor OAuth2认证。

**输入**  
```json
{
  "username": "zhangsan",
  "password": "encrypted_password",
  "captcha": "1234",
  "captchaKey": "cap_xxx"
}
```

**输出**  
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
    "expiresIn": 7200,
    "tokenType": "Bearer",
    "user": {
      "id": 1,
      "username": "zhangsan",
      "realName": "张三",
      "avatar": "https://...",
      "orgName": "产品开发-A产品线",
      "roles": ["A产品线设计师", "工程师"]
    }
  }
}
```

**前置条件**  
用户已注册

**API端点**  
```
POST /api/v1/auth/login
POST /api/v1/auth/logout
POST /api/v1/auth/refresh
```

**验收标准**  
- [ ] 支持用户名密码登录
- [ ] 支持验证码（连续失败3次后启用）
- [ ] Token有效期2小时，Refresh Token有效期7天
- [ ] 登出后Token失效

---

### SRS-USER-002 [MUST] [Phase 1] GitHub风格个人主页

**功能描述**  
每个用户拥有个人主页，展示基本信息、荣誉、项目、能力图谱等。

**输入**  
用户ID

**输出**  
个人主页包含：
- 基本信息区（头像、姓名、职称、所属团队）
- 荣誉展示区（奖项、专利、论文时间轴）
- 项目经历区（已完成/进行中项目列表）
- 能力图谱区（六维能力雷达图）
- 贡献统计区（GitHub风格热力图）
- 活动统计区（近期平台使用时长）

**前置条件**  
用户已登录

**数据模型**  
- `sys_user_honor` 用户荣誉表
- `sys_user_skill` 用户技能表

**API端点**  
```
GET /api/v1/users/{id}/profile
GET /api/v1/users/{id}/honors
GET /api/v1/users/{id}/projects
GET /api/v1/users/{id}/contributions
```

**验收标准**  
- [ ] 页面设计参考GitHub个人页
- [ ] 荣誉时间轴支持按年份筛选
- [ ] 能力雷达图使用ECharts实现
- [ ] 贡献热力图展示近12个月活跃度
- [ ] 页面加载 < 2秒

---

### SRS-USER-003 [MUST] [Phase 1] 组织架构树

**功能描述**  
可视化展示部门组织架构，支持展开折叠、搜索定位。

**输入**  
无

**输出**  
组织架构树形图，包含：
- 部门层级（部门→团队→小组）
- 人员卡片（头像、姓名、职称）
- 展开/折叠节点
- 搜索高亮

**前置条件**  
用户已登录

**API端点**  
```
GET /api/v1/orgs/tree
GET /api/v1/orgs/{id}/members
```

**验收标准**  
- [ ] 支持拖拽缩放
- [ ] 节点支持展开/折叠
- [ ] 支持按姓名搜索定位
- [ ] 点击人员卡片跳转个人主页

---

## 6. 功能规格：项目管理 (SRS-PM)

### SRS-PM-001 [MUST] [Phase 1] 项目创建向导（五步流程）

**功能描述**  
通过五步向导引导用户完成项目创建：信息录入→类别选择→流程绑定→团队分配→计划确认。

**输入**  
```json
{
  "projectName": "X波段TR组件",
  "projectCode": "PRJ-2026-001",
  "projectType": 1,
  "managerId": 10,
  "startDate": "2026-03-01",
  "endDate": "2026-12-31",
  "description": "项目描述...",
  "teamMembers": [
    {"userId": 11, "roleType": 2},
    {"userId": 12, "roleType": 3}
  ]
}
```

**输出**  
```json
{
  "code": 200,
  "data": {
    "id": 100,
    "projectCode": "PRJ-2026-001",
    "status": 2,
    "folderPath": "/data/projects/PRJ-2026-001",
    "gitRepoUrl": "http://git.local/rdp/PRJ-2026-001"
  }
}
```

**前置条件**  
用户具有项目创建权限

**API端点**  
```
POST /api/v1/projects
GET /api/v1/project-types
GET /api/v1/workflow-templates
```

**验收标准**  
- [ ] 支持五步向导流程
- [ ] 项目编码自动生成（支持自定义）
- [ ] 自动创建标准化文件夹结构
- [ ] 自动创建Gitea Git仓库
- [ ] 自动绑定对应流程模板
- [ ] 支持中途保存草稿

---

### SRS-PM-002 [MUST] [Phase 2] 甘特图展示与编辑

**功能描述**  
支持项目WBS的甘特图展示，可拖拽调整工期、设置依赖关系。

**输入**  
项目ID

**输出**  
甘特图视图，包含：
- 任务列表（左侧）
- 时间轴（顶部）
- 任务条（可拖拽调整起止时间）
- 依赖关系线

**前置条件**  
项目已创建

**API端点**  
```
GET /api/v1/projects/{id}/gantt
PUT /api/v1/projects/{id}/gantt/tasks/{taskId}
POST /api/v1/projects/{id}/gantt/dependencies
```

**验收标准**  
- [ ] 支持MS Project文件导入导出
- [ ] 任务条支持拖拽调整工期
- [ ] 支持设置任务依赖关系（FS/SS/FF/SF）
- [ ] 支持关键路径计算
- [ ] 支持基线对比

---

### SRS-PM-003 [MUST] [Phase 1] 项目总览仪表盘

**功能描述**  
展示部门所有项目的总览信息，支持多维度筛选。

**输入**  
筛选条件（产品线、时间、负责人、状态）

**输出**  
项目列表+统计图表：
- 项目卡片网格/列表视图
- 状态分布饼图
- 进度趋势折线图
- 负责人工作量柱状图

**前置条件**  
用户已登录

**API端点**  
```
GET /api/v1/projects?page=1&size=20&status=&type=&managerId=
GET /api/v1/projects/statistics
```

**验收标准**  
- [ ] 支持按产品线/时间/负责人/状态筛选
- [ ] 支持卡片/列表视图切换
- [ ] 支持导出Excel
- [ ] 支持项目收藏

---

## 7. 功能规格：项目开发 (SRS-DEV)

### SRS-DEV-001 [MUST] [Phase 2] 流程全景视图

**功能描述**  
图形化展示项目完整开发流程，当前活动高亮标识。

**输入**  
项目ID

**输出**  
BPMN流程图，包含：
- 流程节点（活动、网关、事件）
- 当前活动高亮
- 已完成节点标记
- 点击节点查看详情

**前置条件**  
项目已绑定流程模板

**API端点**  
```
GET /api/v1/projects/{id}/workflow
GET /api/v1/projects/{id}/activities/{activityId}
```

**验收标准**  
- [ ] 支持BPMN 2.0规范图形
- [ ] 当前活动高亮显示
- [ ] 支持缩放/拖拽
- [ ] 点击节点显示活动定义、输入输出、执行指南

---

### SRS-DEV-002 [MUST] [Phase 2] 活动执行面板

**功能描述**  
任务认领、进度填报、交付物上传、完成确认。

**输入**  
```json
{
  "activityId": 1001,
  "progress": 100,
  "actualHours": 8,
  "deliverables": [
    {"name": "设计文档.docx", "fileId": 2001}
  ],
  "comment": "已完成"
}
```

**输出**  
活动状态更新，触发下一活动

**前置条件**  
用户被分配为该活动执行人

**API端点**  
```
POST /api/v1/activities/{id}/claim
POST /api/v1/activities/{id}/complete
POST /api/v1/activities/{id}/progress
```

**验收标准**  
- [ ] 支持任务认领/转交
- [ ] 支持进度填报（0-100%）
- [ ] 支持交付物上传（多文件）
- [ ] 支持关联Git Commit
- [ ] 完成时自动触发下一活动

---

### SRS-DEV-003 [MUST] [Phase 2] 本地软件集成

**功能描述**  
通过自定义协议调用本地软件（Altium Designer、Obsidian等）打开项目文件。

**输入**  
点击模板文件，触发协议调用：
```
rdp://open?file=/data/projects/PRJ-001/template.schdoc&type=altium
obsidian://open?vault=project-vault&file=note.md
```

**输出**  
本地软件打开对应文件

**前置条件**  
用户已安装桌面辅助程序

**验收标准**  
- [ ] 支持rdp://自定义协议
- [ ] 支持obsidian://协议
- [ ] 支持zotero://协议
- [ ] 支持文件保存后自动Git Commit
- [ ] 支持文件同步到服务器

---

## 8. 功能规格：产品货架 (SRS-SHELF)

### SRS-SHELF-001 [MUST] [Phase 2] 产品分类浏览

**功能描述**  
按技术领域、应用场景、成熟度等级多维组织产品展示。

**输入**  
分类ID + 筛选条件

**输出**  
产品卡片列表，包含：
- 产品图片
- 产品名称、型号
- 成熟度等级（TRL）
- 开发人员
- 使用次数

**前置条件**  
无

**API端点**  
```
GET /api/v1/product-categories/tree
GET /api/v1/products?page=1&size=20&categoryId=&trlLevel=
```

**验收标准**  
- [ ] 支持分类树导航
- [ ] 支持TRL等级筛选
- [ ] 支持关键词搜索
- [ ] 支持列表/网格视图切换

---

### SRS-SHELF-002 [MUST] [Phase 2] 产品详情与选用

**功能描述**  
展示产品详细信息，支持一键申请集成。

**输入**  
产品ID

**输出**  
产品详情页，包含：
- 产品规格参数
- 性能数据图表
- 测试报告
- 应用履历
- 问题记录
- 选用按钮

**前置条件**  
用户已登录

**API端点**  
```
GET /api/v1/products/{id}
POST /api/v1/products/{id}/usage-request
```

**验收标准**  
- [ ] 规格参数以表格形式展示
- [ ] 支持下载数据手册
- [ ] 支持查看历史版本
- [ ] 选用申请需审批

---

## 9. 功能规格：技术货架 (SRS-TECH)

### SRS-TECH-001 [MUST] [Phase 2] 技术树展示

**功能描述**  
分类、分领域展示相关技术，技术树可视化。

**输入**  
技术分类ID

**输出**  
技术树形图，包含：
- 技术分类节点
- 技术详情卡片
- 展开/折叠节点

**前置条件**  
无

**API端点**  
```
GET /api/v1/tech-categories/tree
GET /api/v1/technologies?page=1&size=20&categoryId=
```

**验收标准**  
- [ ] 使用ECharts Tree或D3.js实现
- [ ] 支持节点展开/折叠
- [ ] 支持缩放/拖拽
- [ ] 点击节点查看技术详情

---

## 10. 功能规格：知识库 (SRS-KB)

### SRS-KB-001 [MUST] [Phase 3] 知识分类管理

**功能描述**  
支持理论知识、标准规范、制度文件、流程说明、案例等分类管理。

**输入**  
```json
{
  "title": "微波电路设计规范",
  "categoryId": 2,
  "content": "# 设计规范\\n...",
  "tags": ["微波", "设计规范"],
  "attachments": [{"name": "附件.pdf", "fileId": 3001}]
}
```

**输出**  
创建知识条目

**前置条件**  
用户具有知识编辑权限

**API端点**  
```
GET /api/v1/kb-categories/tree
GET /api/v1/knowledges?page=1&size=20&categoryId=
POST /api/v1/knowledges
```

**验收标准**  
- [ ] 支持7种知识分类
- [ ] 支持Markdown编辑
- [ ] 支持附件上传
- [ ] 支持版本管理

---

### SRS-KB-002 [MUST] [Phase 3] 全文搜索

**功能描述**  
基于MeiliSearch的全文搜索，支持关键词高亮。

**输入**  
搜索关键词

**输出**  
搜索结果列表，包含：
- 标题（关键词高亮）
- 摘要（关键词高亮）
- 分类标签
- 相关性评分

**前置条件**  
无

**API端点**  
```
GET /api/v1/knowledges/search?keyword=&page=1&size=20
```

**验收标准**  
- [ ] 搜索响应时间 < 200ms
- [ ] 支持中文分词
- [ ] 关键词高亮显示
- [ ] 支持按分类筛选

---

### SRS-KB-003 [MUST] [Phase 3] Obsidian集成

**功能描述**  
支持Obsidian Vault同步，Wiki链接解析。

**输入**  
Obsidian Vault路径

**输出**  
知识库与Obsidian双向同步

**前置条件**  
用户已配置Obsidian Vault

**API端点**  
```
POST /api/v1/kb/sync/obsidian
GET /api/v1/kb/obsidian/status
```

**验收标准**  
- [ ] 支持Obsidian Markdown语法
- [ ] 支持Wiki链接（[[笔记名]]）
- [ ] 支持Callout语法
- [ ] 支持双向同步

---

## 11. 功能规格：技术论坛 (SRS-FORUM)

### SRS-FORUM-001 [SHOULD] [Phase 3] 帖子发布与回复

**功能描述**  
按技术领域划分版块，支持发帖、回帖、搜索。

**输入**  
```json
{
  "boardId": 1,
  "title": "如何优化射频电路的噪声系数？",
  "content": "问题描述...",
  "topicType": 2
}
```

**输出**  
创建主题帖

**前置条件**  
用户已登录

**API端点**  
```
GET /api/v1/forum/boards
GET /api/v1/forum/topics?page=1&size=20&boardId=
POST /api/v1/forum/topics
POST /api/v1/forum/topics/{id}/replies
```

**验收标准**  
- [ ] 支持富文本编辑
- [ ] 支持代码块（语法高亮）
- [ ] 支持@用户
- [ ] 支持附件上传
- [ ] 支持帖子置顶/精华标记

---

## 12. 功能规格：即时通信 (SRS-IM)

### SRS-IM-001 [MUST] [Phase 1] Mattermost集成

**功能描述**  
集成Mattermost即时通讯，支持频道、私聊、文件共享。

**输入**  
无

**输出**  
嵌入Mattermost界面或独立窗口

**前置条件**  
用户已登录，已配置Mattermost

**API端点**  
```
GET /api/v1/im/mattermost/url
```

**验收标准**  
- [ ] 支持SSO登录
- [ ] 支持iframe嵌入
- [ ] 项目事件自动推送至频道
- [ ] 支持Bot消息推送

---

## 13. 非功能规格

### 13.1 性能要求

| 指标 | 要求 | 优先级 |
|------|------|--------|
| 页面首屏加载 | < 3秒 | [MUST] |
| API响应时间(P95) | < 200ms | [MUST] |
| 并发用户数 | 支持200人同时在线 | [MUST] |
| 文件上传速度 | > 10MB/s (局域网) | [SHOULD] |
| 搜索响应时间 | < 200ms | [MUST] |
| 报表生成时间 | < 30秒 | [SHOULD] |

### 13.2 安全要求

| 要求 | 说明 | 优先级 |
|------|------|--------|
| 传输加密 | HTTPS/TLS 1.2+ | [MUST] |
| 密码策略 | 8位以上，含大小写+数字+特殊字符 | [MUST] |
| 会话管理 | 超时自动退出(30分钟) | [MUST] |
| 操作审计 | 记录所有关键操作 | [MUST] |
| 数据加密 | 敏感数据AES-256加密存储 | [MUST] |
| 访问控制 | 基于RBAC的细粒度权限 | [MUST] |

### 13.3 可用性要求

| 要求 | 说明 | 优先级 |
|------|------|--------|
| 系统可用性 | ≥ 99.5% | [MUST] |
| 数据备份 | 每日全量备份 + 实时增量 | [MUST] |
| 故障恢复 | < 30分钟 | [SHOULD] |
| 灾难恢复 | RTO < 4小时 | [SHOULD] |

### 13.4 兼容性要求

| 要求 | 说明 | 优先级 |
|------|------|--------|
| 浏览器兼容 | Chrome, Edge, Firefox最新2个版本 | [MUST] |
| 操作系统 | Windows 10+, macOS, Linux | [MUST] |
| 移动端 | 响应式设计，支持平板访问 | [SHOULD] |

---

## 14. API 设计规范

### 14.1 错误码定义

| 错误码 | 错误类型 | 说明 |
|--------|----------|------|
| 200 | SUCCESS | 成功 |
| 400 | BAD_REQUEST | 请求参数错误 |
| 401 | UNAUTHORIZED | 未认证 |
| 403 | FORBIDDEN | 无权限 |
| 404 | NOT_FOUND | 资源不存在 |
| 409 | CONFLICT | 资源冲突 |
| 422 | UNPROCESSABLE | 业务校验失败 |
| 429 | TOO_MANY_REQUESTS | 请求过于频繁 |
| 500 | INTERNAL_ERROR | 服务器内部错误 |
| 503 | SERVICE_UNAVAILABLE | 服务不可用 |

**业务错误码 (6xxx系列):**

| 错误码 | 说明 |
|--------|------|
| 6001 | 用户名已存在 |
| 6002 | 用户不存在 |
| 6003 | 密码错误 |
| 6004 | 用户已禁用 |
| 6101 | 项目编号已存在 |
| 6102 | 项目不存在 |
| 6103 | 项目状态不允许操作 |
| 6201 | 任务不存在 |
| 6202 | 任务状态不允许操作 |
| 6301 | 交付物不存在 |

### 14.2 分页规范

```
请求参数:
  page: 当前页码（从1开始）
  page_size: 每页条数（默认20，最大100）
  sort: 排序字段（可选，如：-created_at表示按创建时间倒序）

响应格式:
{
  "code": 200,
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "pages": 5
  }
}
```

### 14.3 认证规范

```
请求头:
  Authorization: Bearer <access_token>

Token刷新:
  POST /api/v1/auth/refresh
  Body: { "refreshToken": "..." }
```

---

## 15. 验收测试矩阵

| 功能模块 | 测试项 | 测试方法 | 通过标准 |
|----------|--------|----------|----------|
| 用户管理 | 用户CRUD | 自动化测试 | 100%通过 |
| 用户管理 | 登录认证 | 自动化测试 | 100%通过 |
| 用户管理 | 权限控制 | 自动化测试 | 100%通过 |
| 项目管理 | 项目创建 | 自动化测试 | 100%通过 |
| 项目管理 | 甘特图操作 | 手工测试 | 核心场景通过 |
| 项目开发 | 流程执行 | 自动化测试 | 100%通过 |
| 项目开发 | 本地软件集成 | 手工测试 | 核心场景通过 |
| 产品货架 | 产品CRUD | 自动化测试 | 100%通过 |
| 知识库 | 知识CRUD | 自动化测试 | 100%通过 |
| 知识库 | 全文搜索 | 自动化测试 | 100%通过 |
| 论坛 | 帖子CRUD | 自动化测试 | 100%通过 |
| IM | Mattermost集成 | 手工测试 | 核心场景通过 |
| 性能 | 负载测试 | 性能测试 | 达到性能指标 |
| 安全 | 渗透测试 | 安全测试 | 无高危漏洞 |

---

## 附录：部署方式修正说明

> ⚠️ **重要变更**（2026-02-21）

根据原始需求约束（"**不使用Docker和K8s虚拟化技术**"），本文档中的部署方式已修正：

| 项目 | 原文档 | 修正后 |
|------|--------|--------|
| 部署环境 | 离线局域网，Docker Compose单节点 | 离线局域网，**systemd裸机部署** |
| 服务管理 | Docker + Docker Compose 24.x [MUST] | **systemd + 裸机二进制部署** [MUST] |

详细部署方案请参考：[03-deployment-correction.md](03-deployment-correction.md)

---

*微波研发部门研发管理平台 — 需求规格说明书 V1.0*  
*© 2026 微波研发部门 | 内部文档*
