# 企业级研发管理平台详细实施方案

> **版本**：v1.0  
> **日期**：2025年1月  
> **文档类型**：技术实施方案

---

## 目录

1. [系统整体架构](#1-系统整体架构)
2. [技术栈选型方案](#2-技术栈选型方案)
3. [数据库架构设计](#3-数据库架构设计)
4. [模块间接口与交互设计](#4-模块间接口与交互设计)
5. [核心业务流程图](#5-核心业务流程图)
6. [开源项目集成方案](#6-开源项目集成方案)
7. [部署架构设计](#7-部署架构设计)
8. [实施路线图](#8-实施路线图)

---

## 1. 系统整体架构

### 1.1 架构设计原则

| 原则 | 说明 |
|------|------|
| **模块化设计** | 各业务模块独立部署、独立扩展 |
| **微服务架构** | 核心服务采用微服务化，便于维护和升级 |
| **开源优先** | 优先采用成熟开源方案，降低开发成本 |
| **渐进式集成** | 支持分阶段实施，降低项目风险 |
| **统一认证** | 单点登录(SSO)，统一身份管理 |
| **数据隔离** | 支持多租户，数据权限精细控制 |

### 1.2 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        用户访问层 (Presentation)                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 门户界面 │  │ 项目管理 │  │ 产品货架 │  │ 知识库   │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 用户管理 │  │ 项目开发 │  │ 技术货架 │  │ 技术论坛 │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        API网关层 (Gateway)                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ 路由转发     │  │ 认证鉴权     │  │ 限流熔断     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        服务层 (Service)                           │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 用户服务 │  │ 项目服务 │  │ 流程服务 │  │ 文件服务 │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 通知服务 │  │ 搜索服务 │  │ 报表服务 │  │ IM服务   │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        数据层 (Data)                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │PostgreSQL│  │  Redis   │  │MeiliSearch│  │  MinIO  │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────────────────┘
```

### 1.3 分层架构说明

| 层级 | 职责 | 主要技术 |
|------|------|----------|
| **用户访问层** | 提供多终端访问入口 | Web、桌面、移动 |
| **接入层** | 流量分发、安全防护 | Nginx、CDN、WAF |
| **前端应用层** | 用户界面渲染和交互 | React + TypeScript + Vite + Ant Design |
| **微服务层** | 核心业务逻辑处理 | Go + Gin |
| **外部系统集成层** | 集成成熟开源系统 | Casdoor、Mattermost、Gitea等 |
| **数据存储层** | 数据持久化和检索 | PostgreSQL、Redis、MeiliSearch、MinIO |
| **基础设施层** | 服务管理和运维 | systemd、Prometheus |

### 1.4 模块依赖关系图

```
门户模块 ← 用户管理 ← 认证中心
    ↓
项目管理 ←→ 工作流引擎 ←→ 项目开发
    ↓           ↓              ↓
产品货架    技术货架 ←→ 文件存储
    ↓           ↓
    └────→ 技术论坛 ←── 即时通讯
```

---

## 2. 技术栈选型方案

### 2.1 前端技术栈

| 技术 | 版本 | 用途 | 选型理由 |
|------|------|------|----------|
| **React** | 18.x+ | 前端框架 | 组件化开发，生态丰富，企业级应用首选 |
| **TypeScript** | 5.0+ | 类型系统 | 增强代码可维护性，减少运行时错误 |
| **Ant Design** | 5.x | UI组件库 | 企业级组件丰富，设计规范统一 |
| **Zustand** | 4.x | 状态管理 | 轻量级，TypeScript友好，易于使用 |
| **React Router** | 6.x | 路由管理 | React官方路由，支持动态路由、嵌套路由 |
| **Axios** | 1.6+ | HTTP客户端 | 拦截器支持，请求取消，TypeScript支持 |
| **Vite** | 5.0+ | 构建工具 | 极速冷启动，HMR热更新，现代化构建 |
| **Tailwind CSS** | 3.x | CSS框架 | 原子化CSS，快速开发，自定义主题 |
| **ECharts** | 5.4+ | 数据可视化 | 丰富的图表类型，企业级可视化方案 |
| **AntV G6** | 4.8+ | 流程图/技术树 | 强大的图可视化引擎，适合流程展示 |
| **qiankun** | 2.x | 微前端框架 | 支持微前端架构，渐进式升级 |

### 2.2 后端技术栈

| 技术 | 版本 | 用途 | 选型理由 |
|------|------|------|----------|
| **Go** | 1.22+ | 编程语言 | 高性能，编译速度快，云原生友好 |
| **Gin** | 1.9+ | Web框架 | 高性能HTTP框架，路由灵活，中间件支持 |
| **GORM** | 2.0+ | ORM框架 | 功能完善，支持多种数据库，迁移方便 |
| **Casbin** | 2.x+ | 权限管理 | 支持RBAC/ABAC，策略灵活，性能好 |
| **Viper** | 1.18+ | 配置管理 | 支持多种格式，热加载，环境变量 |
| **Zap** | 1.26+ | 日志框架 | 高性能结构化日志，支持多种输出 |
| **Flowable** | 7.0+ | 工作流引擎 | BPMN 2.0完整支持，REST API |
| **Redis** | 7.0+ | 缓存/会话 | 高性能缓存，支持分布式锁 |
| **RabbitMQ** | 3.12+ | 消息队列 | 可靠消息传递，支持延迟队列 |
| **MeiliSearch** | 1.x+ | 全文检索 | 轻量级，中文分词优秀，毫秒级响应 |
| **MinIO** | 2024+ | 对象存储 | S3兼容，高性能，易于部署 |

### 2.3 数据库选型

| 数据库 | 用途 | 选型理由 |
|--------|------|----------|
| **PostgreSQL 16** | 主业务数据库 | 功能强大，JSON支持好，开源社区活跃 |
| **Redis 7.0** | 缓存/会话/锁 | 高性能，支持多种数据结构 |
| **MeiliSearch 1.x** | 全文检索 | 轻量级；中文分词优秀；离线部署简单（Rust编写） |
| **MongoDB 7.0** | 文档存储 | 知识库文档、日志等非结构化数据 |

### 2.4 基础设施技术栈

| 技术 | 用途 | 选型理由 |
|------|------|----------|
| **systemd** | 服务管理 | 无容器开销；进程级服务管理；符合无Docker约束 |
| **Nginx** | 反向代理/负载均衡 | 高性能，配置灵活 |
| **Prometheus** | 监控采集 | 云原生监控标准 |
| **Grafana** | 监控可视化 | 丰富的仪表盘模板 |
| **Loki + Grafana** | 日志收集分析 | 轻量级日志聚合，与Prometheus生态集成 |

### 2.5 技术选型对比分析

#### 前端框架对比

| 框架 | 学习成本 | 生态成熟度 | 性能 | 企业适用性 | 推荐度 |
|------|----------|------------|------|------------|--------|
| React 18 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ✅ 推荐 |
| Vue 3 | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | - |
| Angular 17 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | - |

#### 后端框架对比

| 框架 | 开发效率 | 性能 | 生态 | 微服务支持 | 推荐度 |
|------|----------|------|------|------------|--------|
| Go Gin | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ✅ 推荐 |
| Spring Boot | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | - |
| Node.js Nest | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | - |

---

## 3. 数据库架构设计

### 3.1 数据库整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                            应用层                                 │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌───────────┐ │
│  │  业务应用  │  │  外部系统  │  │  报表应用  │  │  管理后台 │ │
│  └──────┬─────┘  └──────┬─────┘  └──────┬─────┘  └─────┬─────┘ │
└─────────┼───────────────┼───────────────┼──────────────┼────────┘
          │               │               │              │
          └───────────────┴───────────────┴──────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                          数据库集群                               │
├─────────────────────────────────────────────────────────────────┤
│  PostgreSQL主从                                                      │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    │
│  │ PostgreSQL Master│───→│ PostgreSQL Slave │    │ PostgreSQL Slave │    │
│  │     业务写       │    │     业务读       │    │     报表读       │    │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘    │
├─────────────────────────────────────────────────────────────────┤
│  NoSQL集群                                                      │
│  ┌─────────┐  ┌─────────┐  ┌─────────────┐  ┌─────────┐        │
│  │  Redis  │  │ MongoDB │  │ MeiliSearch │  │  MinIO  │        │
│  │ 缓存/会话│  │ 文档存储 │  │  全文搜索   │  │ 对象存储 │        │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘            │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 核心实体关系图（ER图）

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

### 3.3 数据库表结构设计

#### 用户权限模块

```sql
-- 用户表
CREATE TABLE sys_user (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(100) NOT NULL COMMENT '密码',
    nickname VARCHAR(50) COMMENT '昵称',
    email VARCHAR(100) COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    avatar VARCHAR(200) COMMENT '头像URL',
    dept_id BIGINT COMMENT '部门ID',
    title_level TINYINT DEFAULT 1 COMMENT '职称等级 1-5',
    status TINYINT DEFAULT 1 COMMENT '状态 0-禁用 1-启用',
    login_ip VARCHAR(50) COMMENT '最后登录IP',
    login_date DATETIME COMMENT '最后登录时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_dept_id (dept_id),
    INDEX idx_status (status)
) COMMENT='系统用户表';

-- 部门表
CREATE TABLE sys_dept (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    parent_id BIGINT DEFAULT 0 COMMENT '父部门ID',
    ancestors VARCHAR(500) COMMENT '祖级列表',
    dept_name VARCHAR(50) NOT NULL COMMENT '部门名称',
    order_num INT DEFAULT 0 COMMENT '显示排序',
    leader VARCHAR(50) COMMENT '负责人',
    phone VARCHAR(20) COMMENT '联系电话',
    email VARCHAR(100) COMMENT '邮箱',
    status TINYINT DEFAULT 1 COMMENT '状态',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id)
) COMMENT='部门表';

-- 角色表
CREATE TABLE sys_role (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_name VARCHAR(50) NOT NULL COMMENT '角色名称',
    role_code VARCHAR(50) NOT NULL UNIQUE COMMENT '角色编码',
    role_sort INT DEFAULT 0 COMMENT '显示排序',
    data_scope TINYINT DEFAULT 1 COMMENT '数据范围 1-全部 2-本部门 3-本部门及以下 4-仅本人',
    status TINYINT DEFAULT 1 COMMENT '状态',
    remark VARCHAR(500) COMMENT '备注',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_role_code (role_code)
) COMMENT='角色表';

-- 用户角色关联表
CREATE TABLE sys_user_role (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_role (user_id, role_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id)
) COMMENT='用户角色关联表';

-- 权限表
CREATE TABLE sys_permission (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    parent_id BIGINT DEFAULT 0 COMMENT '父权限ID',
    perm_name VARCHAR(50) NOT NULL COMMENT '权限名称',
    perm_code VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
    perm_type TINYINT COMMENT '权限类型 1-菜单 2-按钮 3-接口',
    path VARCHAR(200) COMMENT '路由路径',
    component VARCHAR(200) COMMENT '组件路径',
    icon VARCHAR(100) COMMENT '图标',
    order_num INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态',
    INDEX idx_parent_id (parent_id)
) COMMENT='权限表';

-- 角色权限关联表
CREATE TABLE sys_role_permission (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    UNIQUE KEY uk_role_perm (role_id, permission_id)
) COMMENT='角色权限关联表';
```

#### 项目管理模块

```sql
-- 项目表
CREATE TABLE pm_project (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_code VARCHAR(50) NOT NULL UNIQUE COMMENT '项目编码',
    project_name VARCHAR(100) NOT NULL COMMENT '项目名称',
    project_type TINYINT NOT NULL COMMENT '项目类型 1-单机 2-模块 3-软件 4-技术开发 5-流程开发 6-知识库 7-产品上架',
    project_category VARCHAR(50) COMMENT '项目分类',
    dept_id BIGINT COMMENT '所属部门ID',
    manager_id BIGINT COMMENT '项目经理ID',
    status TINYINT DEFAULT 0 COMMENT '状态 0-草稿 1-待启动 2-执行中 3-评审中 4-已暂停 5-已完成 6-已归档',
    priority TINYINT DEFAULT 2 COMMENT '优先级 1-高 2-中 3-低',
    progress DECIMAL(5,2) DEFAULT 0 COMMENT '进度百分比',
    start_date DATE COMMENT '计划开始日期',
    end_date DATE COMMENT '计划结束日期',
    actual_start DATE COMMENT '实际开始日期',
    actual_end DATE COMMENT '实际结束日期',
    budget DECIMAL(15,2) COMMENT '预算',
    description TEXT COMMENT '项目描述',
    git_repo_url VARCHAR(200) COMMENT 'Git仓库地址',
    folder_path VARCHAR(500) COMMENT '项目文件夹路径',
    workflow_id BIGINT COMMENT '绑定的工作流ID',
    created_by BIGINT COMMENT '创建人',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_project_type (project_type),
    INDEX idx_status (status),
    INDEX idx_dept_id (dept_id),
    INDEX idx_manager_id (manager_id)
) COMMENT='项目表';

-- 项目成员表
CREATE TABLE pm_project_member (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '项目ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role_type TINYINT COMMENT '角色类型 1-项目经理 2-技术负责人 3-开发人员 4-测试人员',
    join_date DATE COMMENT '加入日期',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_project_user (project_id, user_id),
    INDEX idx_project_id (project_id),
    INDEX idx_user_id (user_id)
) COMMENT='项目成员表';

-- 项目阶段表
CREATE TABLE pm_project_phase (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '项目ID',
    phase_code VARCHAR(50) COMMENT '阶段编码',
    phase_name VARCHAR(100) NOT NULL COMMENT '阶段名称',
    phase_order INT DEFAULT 0 COMMENT '阶段顺序',
    parent_id BIGINT DEFAULT 0 COMMENT '父阶段ID',
    status TINYINT DEFAULT 0 COMMENT '状态 0-未开始 1-进行中 2-已完成',
    planned_start DATE COMMENT '计划开始日期',
    planned_end DATE COMMENT '计划结束日期',
    actual_start DATE COMMENT '实际开始日期',
    actual_end DATE COMMENT '实际结束日期',
    progress DECIMAL(5,2) DEFAULT 0 COMMENT '进度百分比',
    description TEXT COMMENT '阶段描述',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_project_id (project_id),
    INDEX idx_status (status)
) COMMENT='项目阶段表';

-- 项目活动表
CREATE TABLE pm_project_activity (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '项目ID',
    phase_id BIGINT COMMENT '所属阶段ID',
    activity_code VARCHAR(50) COMMENT '活动编码',
    activity_name VARCHAR(100) NOT NULL COMMENT '活动名称',
    activity_type TINYINT COMMENT '活动类型',
    assignee_id BIGINT COMMENT '执行人ID',
    reviewer_id BIGINT COMMENT '审核人ID',
    status TINYINT DEFAULT 0 COMMENT '状态 0-未开始 1-进行中 2-待审核 3-已完成 4-已跳过',
    priority TINYINT DEFAULT 2 COMMENT '优先级',
    planned_start DATE COMMENT '计划开始日期',
    planned_end DATE COMMENT '计划结束日期',
    actual_start DATE COMMENT '实际开始日期',
    actual_end DATE COMMENT '实际结束日期',
    estimated_hours INT COMMENT '预估工时',
    actual_hours INT COMMENT '实际工时',
    progress DECIMAL(5,2) DEFAULT 0 COMMENT '进度百分比',
    description TEXT COMMENT '活动描述',
    input_deliverables TEXT COMMENT '输入交付物(JSON)',
    output_deliverables TEXT COMMENT '输出交付物(JSON)',
    execution_guide TEXT COMMENT '执行指南',
    template_files TEXT COMMENT '模板文件(JSON)',
    workflow_instance_id VARCHAR(100) COMMENT '工作流实例ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_project_id (project_id),
    INDEX idx_phase_id (phase_id),
    INDEX idx_assignee_id (assignee_id),
    INDEX idx_status (status)
) COMMENT='项目活动表';

-- 项目交付物表
CREATE TABLE pm_deliverable (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '项目ID',
    activity_id BIGINT COMMENT '关联活动ID',
    deliverable_code VARCHAR(50) COMMENT '交付物编码',
    deliverable_name VARCHAR(100) NOT NULL COMMENT '交付物名称',
    deliverable_type TINYINT COMMENT '交付物类型 1-文档 2-图纸 3-代码 4-测试报告 5-其他',
    file_path VARCHAR(500) COMMENT '文件路径',
    file_size BIGINT COMMENT '文件大小',
    file_hash VARCHAR(64) COMMENT '文件哈希',
    version VARCHAR(20) COMMENT '版本号',
    status TINYINT DEFAULT 0 COMMENT '状态 0-草稿 1-已提交 2-审核中 3-已通过 4-已驳回',
    uploader_id BIGINT COMMENT '上传人ID',
    reviewer_id BIGINT COMMENT '审核人ID',
    upload_time DATETIME COMMENT '上传时间',
    review_time DATETIME COMMENT '审核时间',
    review_comment TEXT COMMENT '审核意见',
    git_commit_id VARCHAR(50) COMMENT 'Git提交ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_project_id (project_id),
    INDEX idx_activity_id (activity_id),
    INDEX idx_status (status)
) COMMENT='项目交付物表';

-- 项目里程碑表
CREATE TABLE pm_milestone (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    project_id BIGINT NOT NULL COMMENT '项目ID',
    milestone_name VARCHAR(100) NOT NULL COMMENT '里程碑名称',
    milestone_type TINYINT COMMENT '里程碑类型 1-DCP决策点 2-技术评审 3-阶段验收',
    planned_date DATE COMMENT '计划日期',
    actual_date DATE COMMENT '实际日期',
    status TINYINT DEFAULT 0 COMMENT '状态 0-未达成 1-已达成 2-已延期',
    description TEXT COMMENT '描述',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_project_id (project_id)
) COMMENT='项目里程碑表';
```

#### 产品货架模块

```sql
-- 产品分类表
CREATE TABLE ps_product_category (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    parent_id BIGINT DEFAULT 0 COMMENT '父分类ID',
    category_code VARCHAR(50) NOT NULL UNIQUE COMMENT '分类编码',
    category_name VARCHAR(100) NOT NULL COMMENT '分类名称',
    category_level TINYINT COMMENT '分类层级',
    category_type TINYINT COMMENT '分类类型 1-单机 2-模块 3-软件 4-基础',
    icon VARCHAR(100) COMMENT '图标',
    order_num INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id),
    INDEX idx_category_type (category_type)
) COMMENT='产品分类表';

-- 产品表
CREATE TABLE ps_product (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    product_code VARCHAR(50) NOT NULL UNIQUE COMMENT '产品编码',
    product_name VARCHAR(100) NOT NULL COMMENT '产品名称',
    category_id BIGINT NOT NULL COMMENT '分类ID',
    product_type TINYINT COMMENT '产品类型',
    maturity_level TINYINT DEFAULT 1 COMMENT '成熟度等级 TRL 1-9',
    developer_id BIGINT COMMENT '开发人员ID',
    version VARCHAR(20) DEFAULT '1.0.0' COMMENT '当前版本',
    specification TEXT COMMENT '规格参数(JSON)',
    performance_data TEXT COMMENT '性能数据',
    test_report_url VARCHAR(200) COMMENT '测试报告链接',
    image_urls TEXT COMMENT '产品图片(JSON)',
    document_urls TEXT COMMENT '文档链接(JSON)',
    application_history TEXT COMMENT '应用履历',
    authorization_type TINYINT DEFAULT 1 COMMENT '授权方式 1-公开 2-申请 3-授权',
    status TINYINT DEFAULT 1 COMMENT '状态 0-下架 1-上架 2-待审核',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    use_count INT DEFAULT 0 COMMENT '使用次数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category_id (category_id),
    INDEX idx_maturity_level (maturity_level),
    INDEX idx_status (status)
) COMMENT='产品表';

-- 产品版本表
CREATE TABLE ps_product_version (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    product_id BIGINT NOT NULL COMMENT '产品ID',
    version VARCHAR(20) NOT NULL COMMENT '版本号',
    version_type TINYINT COMMENT '版本类型 1-主要版本 2-次要版本 3-修订版本',
    change_log TEXT COMMENT '变更说明',
    file_urls TEXT COMMENT '文件链接(JSON)',
    is_current TINYINT DEFAULT 0 COMMENT '是否当前版本',
    created_by BIGINT COMMENT '创建人',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_product_id (product_id),
    INDEX idx_version (version)
) COMMENT='产品版本表';

-- 产品问题记录表
CREATE TABLE ps_product_issue (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    product_id BIGINT NOT NULL COMMENT '产品ID',
    issue_type TINYINT COMMENT '问题类型 1-设计问题 2-应用问题 3-批产问题',
    issue_title VARCHAR(200) NOT NULL COMMENT '问题标题',
    issue_description TEXT COMMENT '问题描述',
    reporter_id BIGINT COMMENT '报告人ID',
    status TINYINT DEFAULT 0 COMMENT '状态 0-待处理 1-处理中 2-已解决 3-已关闭',
    solution TEXT COMMENT '解决方案',
    resolved_by BIGINT COMMENT '解决人ID',
    resolved_at DATETIME COMMENT '解决时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_product_id (product_id),
    INDEX idx_status (status)
) COMMENT='产品问题记录表';

-- 技术分类表
CREATE TABLE ps_tech_category (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    parent_id BIGINT DEFAULT 0 COMMENT '父分类ID',
    category_code VARCHAR(50) NOT NULL UNIQUE COMMENT '分类编码',
    category_name VARCHAR(100) NOT NULL COMMENT '分类名称',
    category_type TINYINT COMMENT '分类类型 1-设计技术 2-试验技术 3-仿真技术',
    order_num INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id),
    INDEX idx_category_type (category_type)
) COMMENT='技术分类表';

-- 技术表
CREATE TABLE ps_technology (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tech_code VARCHAR(50) NOT NULL UNIQUE COMMENT '技术编码',
    tech_name VARCHAR(100) NOT NULL COMMENT '技术名称',
    category_id BIGINT NOT NULL COMMENT '分类ID',
    tech_type TINYINT COMMENT '技术类型',
    maturity_level TINYINT DEFAULT 1 COMMENT '成熟度等级',
    owner_id BIGINT COMMENT '技术负责人ID',
    principle TEXT COMMENT '技术原理',
    application_scenario TEXT COMMENT '应用场景',
    document_urls TEXT COMMENT '技术文档链接(JSON)',
    case_studies TEXT COMMENT '应用案例',
    related_products TEXT COMMENT '关联产品(JSON)',
    related_projects TEXT COMMENT '关联项目(JSON)',
    status TINYINT DEFAULT 1 COMMENT '状态',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category_id (category_id),
    INDEX idx_owner_id (owner_id)
) COMMENT='技术表';
```

#### 知识库模块

```sql
-- 知识分类表
CREATE TABLE kb_category (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    parent_id BIGINT DEFAULT 0 COMMENT '父分类ID',
    category_code VARCHAR(50) NOT NULL UNIQUE COMMENT '分类编码',
    category_name VARCHAR(100) NOT NULL COMMENT '分类名称',
    category_type TINYINT COMMENT '分类类型 1-理论知识 2-标准规范 3-制度文件 4-流程说明 5-案例 6-仿真模型 7-软件指南',
    icon VARCHAR(100) COMMENT '图标',
    order_num INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id),
    INDEX idx_category_type (category_type)
) COMMENT='知识分类表';

-- 知识条目表
CREATE TABLE kb_knowledge (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL COMMENT '标题',
    content LONGTEXT COMMENT '内容',
    summary TEXT COMMENT '摘要',
    category_id BIGINT NOT NULL COMMENT '分类ID',
    knowledge_type TINYINT COMMENT '知识类型',
    author_id BIGINT COMMENT '作者ID',
    status TINYINT DEFAULT 0 COMMENT '状态 0-草稿 1-待审核 2-已发布 3-已驳回 4-已归档',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    download_count INT DEFAULT 0 COMMENT '下载次数',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    version VARCHAR(20) DEFAULT '1.0' COMMENT '版本号',
    is_top TINYINT DEFAULT 0 COMMENT '是否置顶',
    is_essence TINYINT DEFAULT 0 COMMENT '是否精华',
    source_type TINYINT COMMENT '来源 1-原创 2-转载 3-翻译',
    source_url VARCHAR(500) COMMENT '来源链接',
    obsidian_path VARCHAR(500) COMMENT 'Obsidian文件路径',
    zotero_key VARCHAR(100) COMMENT 'Zotero引用键',
    tags TEXT COMMENT '标签(JSON)',
    attachments TEXT COMMENT '附件(JSON)',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category_id (category_id),
    INDEX idx_author_id (author_id),
    INDEX idx_status (status),
    INDEX idx_title (title),
    FULLTEXT INDEX ft_content (title, content, summary)
) COMMENT='知识条目表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 知识版本表
CREATE TABLE kb_knowledge_version (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    knowledge_id BIGINT NOT NULL COMMENT '知识ID',
    version VARCHAR(20) NOT NULL COMMENT '版本号',
    content LONGTEXT COMMENT '内容',
    change_log TEXT COMMENT '变更说明',
    created_by BIGINT COMMENT '创建人',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_knowledge_id (knowledge_id)
) COMMENT='知识版本表';

-- 知识评论表
CREATE TABLE kb_comment (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    knowledge_id BIGINT NOT NULL COMMENT '知识ID',
    parent_id BIGINT DEFAULT 0 COMMENT '父评论ID',
    content TEXT NOT NULL COMMENT '评论内容',
    author_id BIGINT COMMENT '评论人ID',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_knowledge_id (knowledge_id),
    INDEX idx_parent_id (parent_id)
) COMMENT='知识评论表';

-- 知识收藏表
CREATE TABLE kb_favorite (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    knowledge_id BIGINT NOT NULL COMMENT '知识ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_knowledge (user_id, knowledge_id)
) COMMENT='知识收藏表';
```

#### 论坛模块

```sql
-- 论坛版块表
CREATE TABLE forum_board (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    board_name VARCHAR(100) NOT NULL COMMENT '版块名称',
    board_code VARCHAR(50) NOT NULL UNIQUE COMMENT '版块编码',
    description TEXT COMMENT '版块描述',
    icon VARCHAR(100) COMMENT '图标',
    order_num INT DEFAULT 0 COMMENT '排序',
    topic_count INT DEFAULT 0 COMMENT '主题数',
    post_count INT DEFAULT 0 COMMENT '帖子数',
    moderators TEXT COMMENT '版主列表(JSON)',
    status TINYINT DEFAULT 1 COMMENT '状态',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_board_code (board_code)
) COMMENT='论坛版块表';

-- 论坛主题表
CREATE TABLE forum_topic (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    board_id BIGINT NOT NULL COMMENT '版块ID',
    author_id BIGINT NOT NULL COMMENT '作者ID',
    title VARCHAR(200) NOT NULL COMMENT '标题',
    content LONGTEXT COMMENT '内容',
    topic_type TINYINT DEFAULT 1 COMMENT '主题类型 1-普通 2-提问 3-分享 4-讨论 5-公告',
    is_top TINYINT DEFAULT 0 COMMENT '是否置顶',
    is_essence TINYINT DEFAULT 0 COMMENT '是否精华',
    is_locked TINYINT DEFAULT 0 COMMENT '是否锁定',
    view_count INT DEFAULT 0 COMMENT '浏览数',
    reply_count INT DEFAULT 0 COMMENT '回复数',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    last_reply_id BIGINT COMMENT '最后回复ID',
    last_reply_time DATETIME COMMENT '最后回复时间',
    status TINYINT DEFAULT 1 COMMENT '状态 0-删除 1-正常 2-待审核',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_board_id (board_id),
    INDEX idx_author_id (author_id),
    INDEX idx_topic_type (topic_type),
    INDEX idx_is_top (is_top),
    INDEX idx_status (status),
    FULLTEXT INDEX ft_title_content (title, content)
) COMMENT='论坛主题表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 论坛回帖表
CREATE TABLE forum_post (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    topic_id BIGINT NOT NULL COMMENT '主题ID',
    parent_id BIGINT DEFAULT 0 COMMENT '父帖子ID（回复回复）',
    author_id BIGINT NOT NULL COMMENT '作者ID',
    content LONGTEXT COMMENT '内容',
    is_best TINYINT DEFAULT 0 COMMENT '是否最佳回复',
    like_count INT DEFAULT 0 COMMENT '点赞数',
    status TINYINT DEFAULT 1 COMMENT '状态',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_topic_id (topic_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_author_id (author_id)
) COMMENT='论坛回帖表';
```

#### 消息通知模块

```sql
-- 消息表
CREATE TABLE msg_message (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    sender_id BIGINT COMMENT '发送人ID 0-系统',
    receiver_id BIGINT NOT NULL COMMENT '接收人ID',
    msg_type TINYINT COMMENT '消息类型 1-系统 2-任务 3-审批 4-提醒',
    msg_title VARCHAR(200) COMMENT '消息标题',
    msg_content TEXT COMMENT '消息内容',
    related_type VARCHAR(50) COMMENT '关联业务类型',
    related_id BIGINT COMMENT '关联业务ID',
    is_read TINYINT DEFAULT 0 COMMENT '是否已读',
    read_time DATETIME COMMENT '阅读时间',
    send_channel TINYINT DEFAULT 1 COMMENT '发送渠道 1-站内 2-邮件 3-即时通讯',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_receiver_id (receiver_id),
    INDEX idx_is_read (is_read),
    INDEX idx_msg_type (msg_type)
) COMMENT='消息表';
```

### 3.4 数据库分库分表策略

| 数据类型 | 分片策略 | 说明 |
|----------|----------|------|
| 用户数据 | 按用户ID哈希分片 | 用户表、用户角色关联表 |
| 项目数据 | 按项目ID哈希分片 | 项目表、项目成员表、项目任务表 |
| 文件数据 | 按时间范围分片 | 文件表按年月分表 |
| 日志数据 | 按时间范围分片 | 操作日志按天分区 |
| 消息数据 | 按接收者ID哈希分片 | 消息表 |

---

## 4. 模块间接口与交互设计

### 4.1 服务间调用关系图

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   门户模块   │────→│   用户服务   │←────│   认证中心   │
└──────┬──────┘     └──────┬──────┘     └─────────────┘
       │                   │
       ↓                   ↓
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   项目服务   │←───→│   工作流引擎 │←───→│   项目开发   │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       ↓                   ↓                   ↓
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   产品货架   │     │   技术货架   │←───→│   文件存储   │
└──────┬──────┘     └─────────────┘     └─────────────┘
       │
       ↓
┌─────────────┐
│   技术论坛   │←───→ 即时通讯服务
└─────────────┘
```

### 4.2 RESTful API 设计规范

#### 基础规范

| 项目 | 规范 |
|------|------|
| 协议 | HTTPS |
| 数据格式 | JSON |
| 字符编码 | UTF-8 |
| 日期格式 | ISO 8601 (yyyy-MM-ddTHH:mm:ss) |
| 分页参数 | page(页码), size(每页大小) |
| 认证方式 | JWT Token (Bearer) |

#### 统一响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": 1704067200000,
  "requestId": "req_xxxxxxxx"
}
```

#### 分页响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "size": 20,
    "pages": 5
  },
  "timestamp": 1704067200000
}
```

### 4.3 核心模块接口定义

#### 用户管理接口

```yaml
# 用户登录
POST /api/v1/auth/login
Request:
  username: string
  password: string
Response:
  accessToken: string
  refreshToken: string
  expiresIn: number
  user: UserInfo

# 获取当前用户信息
GET /api/v1/auth/userInfo
Response:
  id: number
  username: string
  realName: string
  email: string
  phone: string
  avatar: string
  orgId: number
  orgName: string
  titleLevel: number
  roles: Role[]
  permissions: string[]

# 用户列表查询
GET /api/v1/users?page=1&size=20&keyword=&deptId=
Response:
  list: User[]
  total: number

# 创建用户
POST /api/v1/users
Request:
  username: string
  password: string
  realName: string
  email: string
  phone: string
  deptId: number
  roleIds: number[]

# 更新用户
PUT /api/v1/users/{id}

# 删除用户
DELETE /api/v1/users/{id}
```

#### 项目管理接口

```yaml
# 项目列表查询
GET /api/v1/projects?page=1&size=20&status=&type=&managerId=
Response:
  list: Project[]
  total: number

# 创建项目
POST /api/v1/projects
Request:
  projectCode: string
  projectName: string
  projectType: number
  managerId: number
  startDate: string
  endDate: string
  description: string

# 获取项目详情
GET /api/v1/projects/{id}
Response:
  id: number
  projectCode: string
  projectName: string
  projectType: number
  status: number
  progress: number
  manager: UserInfo
  members: ProjectMember[]
  phases: ProjectPhase[]
  milestones: Milestone[]

# 更新项目
PUT /api/v1/projects/{id}

# 删除项目
DELETE /api/v1/projects/{id}

# 获取项目成员
GET /api/v1/projects/{id}/members

# 添加项目成员
POST /api/v1/projects/{id}/members
Request:
  userId: number
  roleType: number

# 获取项目任务
GET /api/v1/projects/{id}/tasks

# 创建项目任务
POST /api/v1/projects/{id}/tasks

# 更新任务进度
PUT /api/v1/projects/{id}/tasks/{taskId}/progress
Request:
  progress: number
  actualHours: number
```

#### 产品货架接口

```yaml
# 产品分类树
GET /api/v1/product-categories/tree
Response:
  id: number
  categoryCode: string
  categoryName: string
  children: Category[]

# 产品列表查询
GET /api/v1/products?page=1&size=20&categoryId=&keyword=
Response:
  list: Product[]
  total: number

# 获取产品详情
GET /api/v1/products/{id}
Response:
  id: number
  productCode: string
  productName: string
  category: Category
  maturityLevel: number
  developer: UserInfo
  specifications: object
  imageUrls: string[]
  documentUrls: string[]
  versionHistory: ProductVersion[]

# 产品选用
POST /api/v1/products/{id}/usage
Request:
  projectId: number
  usageDesc: string
```

#### 知识库接口

```yaml
# 知识分类树
GET /api/v1/kb-categories/tree

# 知识列表查询
GET /api/v1/knowledges?page=1&size=20&categoryId=&keyword=

# 获取知识详情
GET /api/v1/knowledges/{id}

# 创建知识
POST /api/v1/knowledges
Request:
  title: string
  content: string
  categoryId: number
  tags: string[]

# 更新知识
PUT /api/v1/knowledges/{id}

# 删除知识
DELETE /api/v1/knowledges/{id}

# 全文搜索
GET /api/v1/knowledges/search?keyword=&page=1&size=20

# 添加评论
POST /api/v1/knowledges/{id}/comments

# 收藏知识
POST /api/v1/knowledges/{id}/favorite
```

### 4.4 服务间通信方式

| 通信方式 | 使用场景 | 技术实现 |
|----------|----------|----------|
| **同步HTTP调用** | 实时性要求高的查询 | RESTful API |
| **异步消息队列** | 削峰填谷、解耦 | RabbitMQ |
| **事件总线** | 模块间事件通知 | 应用内事件 + 消息队列 |
| **WebSocket** | 实时推送、即时通讯 | STOMP over WebSocket |

### 4.5 消息队列设计

#### 队列划分

| 队列名称 | 用途 | 消费者 |
|----------|------|--------|
| `rdp.notification` | 通知消息 | 通知服务 |
| `rdp.email` | 邮件发送 | 邮件服务 |
| `rdp.search.index` | 搜索索引更新 | 搜索服务 |
| `rdp.file.convert` | 文件转换任务 | 文件服务 |
| `rdp.report.generate` | 报表生成 | 报表服务 |

#### 事件类型

```yaml
# 项目相关事件
ProjectCreated:
  projectId: number
  projectName: string
  creatorId: number
  
ProjectStatusChanged:
  projectId: number
  oldStatus: number
  newStatus: number
  operatorId: number

TaskAssigned:
  taskId: number
  taskName: string
  assigneeId: number
  
# 用户相关事件
UserRegistered:
  userId: number
  username: string
  
# 知识库相关事件
KnowledgePublished:
  knowledgeId: number
  title: string
  authorId: number
```

---

## 5. 核心业务流程图

### 5.1 项目创建流程

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ 开始     │───→│ 填写项目 │───→│ 选择项目 │───→│ 绑定流程 │───→│ 分配团队 │
│          │    │ 基本信息 │    │ 类别     │    │ 模板     │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
                                                                       │
                                                                       ↓
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ 结束     │←───│ 创建完成 │←───│ 确认计划 │←───│ 生成甘特 │←───│ 设置里程碑│
│          │    │          │    │          │    │ 图       │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
```

### 5.2 项目开发执行流程

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ 接收任务 │───→│ 查看活动 │───→│ 下载模板 │───→│ 本地开发 │───→│ 提交交付 │
│          │    │ 定义     │    │          │    │          │    │ 物       │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
                                                                       │
                                                                       ↓
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ 触发下一 │←───│ 完成活动 │←───│ 审核通过 │←───│ 交付物   │←───│ Git版本  │
│ 活动     │    │          │    │          │    │ 审核     │    │ 提交     │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
```

### 5.3 工作流审批流程

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ 提交审批 │───→│ 审批人   │───→│ 审批通过 │───→│ 流程继续 │
│          │    │ 审批     │    │          │    │          │
└──────────┘    └──────────┘    └─────┬────┘    └──────────┘
                                      │
                                      ↓ 驳回
                               ┌──────────┐
                               │ 退回修改 │
                               └──────────┘
```

### 5.4 产品上架流程

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ 项目完成 │───→│ 提交上架 │───→│ 技术评审 │───→│ 质量评审 │───→│ 审批上架 │
│          │    │ 申请     │    │          │    │          │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └─────┬────┘
                                                                        │
                                                                        ↓
                                                               ┌──────────┐
                                                               │ 产品上架 │
                                                               │ 展示     │
                                                               └──────────┘
```

### 5.5 知识贡献流程

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ 创建知识 │───→│ 编写内容 │───→│ 选择分类 │───→│ 提交审核 │───→│ 审核通过 │
│          │    │          │    │          │    │          │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └─────┬────┘
                                                                        │
                                                                        ↓
                                                               ┌──────────┐
                                                               │ 知识发布 │
                                                               │ 检索     │
                                                               └──────────┘
```

### 5.6 即时通讯集成流程

```
┌──────────┐         ┌──────────┐         ┌──────────┐
│ 项目事件 │────────→│ 消息队列 │────────→│ 通知服务 │
│ 触发     │         │          │         │          │
└──────────┘         └──────────┘         └────┬─────┘
                                               │
                          ┌────────────────────┼────────────────────┐
                          ↓                    ↓                    ↓
                   ┌──────────┐        ┌──────────┐        ┌──────────┐
                   │ 站内消息 │        │ 邮件通知 │        │ IM推送   │
                   └──────────┘        └──────────┘        └──────────┘
```

---

## 6. 开源项目集成方案

### 6.1 开源项目选型总览

| 功能模块 | 推荐开源项目 | 版本 | Stars | 集成方式 |
|----------|--------------|------|-------|----------|
| **用户认证/RBAC** | Casdoor + Casbin | Latest | 13K+19.8K | 独立服务 + API集成 |
| **即时通信** | Mattermost | Latest TE | 35.4K | 独立部署 + iframe/API集成 |
| **Git版本管理** | Gitea | 1.22+ | 53.8K | 独立服务 + API调用 |
| **知识库** | Wiki.js / BookStack | Latest | 27K/17K | 参考架构 + 自研 |
| **全文搜索** | MeiliSearch | 1.x | 52K | 独立服务 + SDK集成 |
| **甘特图** | gantt-task-react | Latest | 1K | 前端组件嵌入 |
| **流程引擎** | Flowable / 自研 | 7.0+ | 9K | 嵌入式引擎 |
| **组织架构图** | d3-org-chart | Latest | 2.3K | 前端组件嵌入 |
| **图表可视化** | Apache ECharts | 5.4+ | 64K | 前端组件嵌入 |
| **Markdown编辑** | ByteMD / Milkdown | Latest | 4K+ | 前端组件嵌入 |
| **论坛** | Flarum (参考) | Latest | 16K | 参考架构 + 自研 |
| **对象存储** | MinIO | Latest | 52K | 独立服务 |

### 6.2 开源项目详细说明

#### Casdoor - 统一认证

```yaml
项目信息:
  GitHub: https://github.com/casdoor/casdoor
  Stars: 13K
  技术栈: Go + React + PostgreSQL
  协议: Apache-2.0

主要功能:
  - 内置组织/团队管理
  - 多因素认证(MFA)
  - 130+社交登录支持
  - OAuth2/OIDC/SAML/LDAP支持
  - 多租户支持

集成方式:
  - 独立systemd服务部署
  - 通过OAuth2/OIDC与主平台对接
  - 组织/角色数据通过API同步

部署配置:
  端口: 8000
  数据库: PostgreSQL
  内存要求: 512MB
```

#### Mattermost - 即时通信

```yaml
项目信息:
  GitHub: https://github.com/mattermost/mattermost
  Stars: 35.4K
  技术栈: Go + React + PostgreSQL
  协议: MIT

主要功能:
  - 频道、私聊、群组
  - 文件共享
  - 700+集成
  - Bot框架
  - Webhook/API

集成方式:
  - systemd独立服务部署
  - Nginx反向代理集成
  - SSO与Casdoor对接
  - Webhook接收项目事件

部署配置:
  端口: 8065
  数据库: PostgreSQL
  内存要求: 2GB
```

#### Gitea - Git版本管理

```yaml
项目信息:
  GitHub: https://github.com/go-gitea/gitea
  Stars: 53.8K
  技术栈: Go
  协议: MIT

主要功能:
  - Git仓库管理
  - Issue跟踪
  - Pull Request
  - WebHook
  - 完善API

集成方式:
  - systemd独立服务部署
  - 每个项目自动创建Gitea仓库
  - 通过API实现文件上传/下载/版本管理
  - 自动commit

部署配置:
  端口: 3002
  数据库: PostgreSQL
  内存要求: 1GB
```

#### MeiliSearch - 全文搜索

```yaml
项目信息:
  GitHub: https://github.com/meilisearch/meilisearch
  Stars: 52K
  技术栈: Rust
  协议: MIT

主要功能:
  - 毫秒级搜索响应
  - 容错搜索
  - 中文分词
  - 同义词支持
  - 多租户

集成方式:
  - systemd独立服务部署
  - SDK集成
  - 异步索引更新

部署配置:
  端口: 7700
  数据库: 内置(LMDB)
  内存要求: 1GB
```

#### MinIO - 对象存储

```yaml
项目信息:
  GitHub: https://github.com/minio/minio
  Stars: 52K
  技术栈: Go
  协议: AGPL-3.0

主要功能:
  - S3兼容API
  - 分布式存储
  - 纠删码
  - 加密传输

集成方式:
  - systemd独立服务部署
  - S3 SDK调用

部署配置:
  API端口: 9000
  控制台端口: 9001
  存储: 本地磁盘
  内存要求: 4GB
```

---

## 7. 部署架构设计

### 7.1 局域网部署架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                     设计师工作站 (浏览器 + 辅助程序)              │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Nginx 反向代理 (1.25+)                       │
│           TLS终止 | 路由分发 | 静态资源 | 负载均衡                 │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   rdp-portal    │    │   rdp-api       │    │   rdp-dev       │
│   (前端Shell)   │    │   (Go API)      │    │   (项目开发)    │
│   :3000         │    │   :8080         │    │   :3001         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                       │                       │
          └───────────────────────┼───────────────────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Casdoor       │    │   Gitea         │    │   Mattermost    │
│   (认证/RBAC)   │    │   (Git服务)     │    │   (即时通讯)    │
│   :8000         │    │   :3002         │    │   :8065         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                       │                       │
          └───────────────────────┼───────────────────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │   Redis         │    │   MinIO         │
│   (主数据库)    │    │   (缓存/会话)   │    │   (对象存储)    │
│   :5432         │    │   :6379         │    │   :9000/9001    │
└─────────────────┘    └─────────────────┘    └─────────────────┘

辅助服务:
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   MeiliSearch   │    │   Wiki.js       │    │   备份服务      │
│   (搜索引擎)    │    │   (知识库)      │    │   (定时任务)    │
│   :7700         │    │   :3003         │    │   systemd timer │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 7.2 systemd服务配置示意

```ini
# /etc/systemd/system/rdp-api.service
[Unit]
Description=RDP Platform API Service
After=network.target postgresql.service redis.service
Wants=postgresql.service redis.service

[Service]
Type=simple
User=rdp-user
Group=rdp-user
WorkingDirectory=/opt/rdp
ExecStart=/opt/rdp/bin/rdp-api
Restart=always
RestartSec=5

# 资源限制
MemoryMax=2G
CPUQuota=200%

# 环境变量
Environment="DB_HOST=127.0.0.1"
Environment="DB_PORT=5432"
Environment="DB_USER=rdp"
Environment="DB_PASSWORD=your_password"
Environment="DB_NAME=rdp_db"
Environment="REDIS_HOST=127.0.0.1"
Environment="REDIS_PORT=6379"
Environment="REDIS_PASSWORD=your_redis_password"
Environment="PORT=8080"

# 日志
StandardOutput=journal
StandardError=journal
SyslogIdentifier=rdp-api

[Install]
WantedBy=multi-user.target
```

启用服务：
```bash
sudo systemctl enable rdp-api.service
sudo systemctl start rdp-api.service
```

### 7.3 Nginx统一配置

```nginx
user www-data;
worker_processes auto;
pid /run/nginx.pid;

events {
    worker_connections 4096;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # 日志格式
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /opt/rdp/logs/nginx/access.log main;
    error_log /opt/rdp/logs/nginx/error.log warn;

    # 性能优化
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    client_max_body_size 100M;

    # Gzip压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml;

    # 上游服务定义
    upstream rdp_api {
        server 127.0.0.1:8080;
    }

    upstream casdoor {
        server 127.0.0.1:8000;
    }

    upstream gitea {
        server 127.0.0.1:3002;
    }

    upstream mattermost {
        server 127.0.0.1:8065;
    }

    upstream minio {
        server 127.0.0.1:9000;
    }

    # 主站点配置
    server {
        listen 80;
        server_name _;
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name rdp.local;

        # SSL证书配置
        ssl_certificate /opt/rdp/config/nginx/ssl/rdp.crt;
        ssl_certificate_key /opt/rdp/config/nginx/ssl/rdp.key;
        ssl_protocols TLSv1.2 TLSv1.3;

        # 前端静态资源
        location / {
            root /opt/rdp/bin/rdp-portal/dist;
            try_files $uri $uri/ /index.html;
            expires 1d;
        }

        # API代理
        location /api/ {
            proxy_pass http://rdp_api/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Casdoor认证
        location /auth/ {
            proxy_pass http://casdoor/;
        }

        # Gitea Git服务
        location /git/ {
            proxy_pass http://gitea/;
        }

        # Mattermost即时通讯
        location /chat/ {
            proxy_pass http://mattermost/;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
        }

        # MinIO对象存储
        location /storage/ {
            proxy_pass http://minio/;
        }
    }
}
```

### 7.4 端口分配表

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

## 8. 实施路线图

### 8.1 四期实施路线图

| 阶段 | 周期 | 交付内容 | 里程碑 |
|------|------|----------|--------|
| **一期：基础骨架** | 3个月 (M1-M3) | 门户框架、用户管理(Casdoor)、RBAC权限、项目管理(基础CRUD)、Mattermost集成、**裸机部署方案**、数据库设计 | 可登录使用 |
| **二期：核心业务** | 4个月 (M4-M7) | 流程引擎、项目开发模块、甘特图组件、产品货架模块、技术货架模块、Gitea集成、桌面辅助程序(基础) | 研发流程上线 |
| **三期：知识智能** | 3个月 (M8-M10) | 知识库模块、Obsidian集成、Zotero集成、全文搜索(MeiliSearch)、技术论坛模块、标签智能关联 | 知识管理上线 |
| **四期：优化完善** | 2个月 (M11-M12) | 数据分析仪表盘、MS Project导入导出、辅助程序完善、性能优化、安全加固、用户培训 | 全功能GA |

### 8.2 团队分工建议

| 团队 | 人数 | 负责模块 | 技术栈要求 |
|------|------|----------|------------|
| 前端团队A | 2-3人 | 门户Shell + 用户管理 + 项目管理前端 | React + TypeScript + Ant Design |
| 前端团队B | 2-3人 | 项目开发 + 货架模块 + 知识库前端 | React + TypeScript + ECharts |
| 后端团队 | 3-4人 | API服务 + 流程引擎 + 数据库 | Go + PostgreSQL + Redis |
| 集成/DevOps | 1-2人 | Mattermost/Gitea/Casdoor集成 + 部署 | Nginx + Shell + **systemd** |
| 桌面端 | 1人 | 辅助程序开发 | Electron/Tauri + Rust/Node.js |

### 8.3 风险与应对

| 风险 | 等级 | 应对措施 |
|------|------|----------|
| 离线环境依赖管理复杂 | 🔴 高 | 预下载所有依赖包至本地仓库；制作离线安装包；自动化部署脚本 |
| 微前端技术复杂度高 | 🟡 中 | 先以Monorepo开发，成熟后再拆分微前端；保留降级方案 |
| 本地软件集成兼容性 | 🟡 中 | 桌面辅助程序做好版本适配；提供手动上传兜底方案 |
| 开源项目版本升级风险 | 🟢 低 | 锁定经过验证的版本号；建立内部Fork仓库 |
| 数据迁移和历史数据 | 🟡 中 | 设计通用导入接口；提供Excel/CSV批量导入工具 |
| 用户接受度和培训成本 | 🟡 中 | 分批培训；提供操作视频和帮助文档；设置"意见反馈"入口 |

---

*微波研发部门研发管理平台 — 详细实施方案 V1.0*  
*© 2026 微波研发部门 | 内部文档*

---

## 补充章节：详细开源项目集成方案

### A.1 开源项目参考方案

> **说明**：以下为备选开源项目参考，本项目最终采用 **React + TypeScript + Vite + Ant Design** 技术栈（详见第2章技术选型）。

#### A.1.1 RuoYi-Vue-Pro 参考方案（Vue技术栈）

| 属性 | 详情 |
|------|------|
| **GitHub** | https://github.com/YunaiV/ruoyi-vue-pro |
| **技术栈** | Spring Boot 3.x + Vue3 + Element Plus |
| **Stars** | 25k+ |
| **开源协议** | MIT License |
| **活跃度** | 高，持续更新 |
| **适用场景** | 快速启动Vue技术栈项目（本项目未采用）|

#### A.1.2 核心功能

- ✅ 完整的RBAC权限管理系统
- ✅ 支持SaaS多租户
- ✅ 内置工作流引擎（Flowable）
- ✅ 代码生成器（一键生成前后端代码）
- ✅ 支持多种登录方式（微信、钉钉、企业微信）
- ✅ 数据权限控制
- ✅ 支付、短信、商城等扩展模块

#### A.1.3 二次开发要点

| 模块 | 开发内容 | 工作量评估 |
|------|----------|------------|
| 门户首页 | 部门信息展示、新闻公告、快捷入口 | 2周 |
| 个人工作台 | 待办事项、消息通知、数据看板 | 3周 |
| 用户管理增强 | 职称体系、能力图谱、个人主页 | 2周 |
| 项目管理 | 项目创建向导、甘特图、Project集成 | 4周 |
| 项目开发 | 流程可视化、交付物管理、本地软件调用 | 4周 |
| 产品货架 | 分类浏览、版本管理、成熟度评估 | 3周 |
| 技术货架 | 技术树展示、技术分类、相似创建 | 2周 |

### A.2 OpenProject 集成方案

#### A.2.1 项目简介

| 属性 | 详情 |
|------|------|
| **官网** | https://www.openproject.org |
| **GitHub** | https://github.com/opf/openproject |
| **技术栈** | Ruby on Rails + Angular |
| **Stars** | 9k+ |
| **开源协议** | GPL-3.0 |

#### A.2.2 核心功能

- ✅ **甘特图**：完整的项目时间线视图
- ✅ **MS Project导入导出**：支持MPP文件
- ✅ **敏捷支持**：Scrum、Kanban看板
- ✅ **任务管理**：工作包、里程碑、依赖关系
- ✅ **时间跟踪**：工时记录和报告
- ✅ **成本管理**：预算和成本控制
- ✅ **文档管理**：Wiki、文件共享
- ✅ **会议管理**

#### A.2.3 集成方式

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   RDP平台   │←────→│   API网关   │←────→│ OpenProject │
└─────────────┘      └─────────────┘      └─────────────┘
      │                                          │
      │         REST API调用                      │
      │                                          │
      └────────── 项目同步 ──────────────────────┘
                   任务同步
                   进度同步
```

#### A.2.4 数据同步策略

| 数据类型 | 同步方向 | 同步频率 | 说明 |
|----------|----------|----------|------|
| 项目信息 | 双向同步 | 实时 | 项目创建/更新时同步 |
| 任务信息 | 双向同步 | 实时 | 任务状态变更时同步 |
| 甘特图数据 | OpenProject→RDP | 定时 | 每5分钟拉取一次 |
| 工时记录 | OpenProject→RDP | 定时 | 每小时拉取一次 |

### A.3 Flowable 集成方案

#### A.3.1 项目简介

| 属性 | 详情 |
|------|------|
| **GitHub** | https://github.com/flowable/flowable-engine |
| **技术栈** | Java |
| **Stars** | 8k+ |
| **开源协议** | Apache-2.0 |

#### A.3.2 核心功能

- ✅ 完整的BPMN 2.0支持
- ✅ DMN决策表
- ✅ CMMN案例管理
- ✅ 流程设计器（Web版）
- ✅ REST API
- ✅ 多租户支持
- ✅ 高性能异步执行

#### A.3.3 集成方式

**方案一：嵌入式集成（推荐）**
```java
// Flowable引擎嵌入Spring Boot
@Configuration
public class FlowableConfig {
    @Bean
    public ProcessEngine processEngine() {
        return ProcessEngineConfiguration
            .createStandaloneProcessEngineConfiguration()
            .setDatabaseType("postgres")
            .setDataSource(dataSource)
            .setAsyncExecutorActivate(true)
            .buildProcessEngine();
    }
}
```

**方案二：独立服务集成**
- Flowable作为独立服务部署
- 通过REST API调用
- 适合高并发场景

### A.4 Wiki.js 集成方案

#### A.4.1 项目简介

| 属性 | 详情 |
|------|------|
| **GitHub** | https://github.com/requarks/wiki |
| **官网** | https://js.wiki |
| **技术栈** | Node.js + Vue + PostgreSQL |
| **Stars** | 27k+ |
| **开源协议** | AGPL-3.0 |

#### A.4.2 核心功能

- ✅ Markdown/WYSIWYG编辑
- ✅ Git同步（双向）
- ✅ 多语言支持
- ✅ 全文搜索
- ✅ 权限管理
- ✅ 模块化存储
- ✅ 可视化编辑器

#### A.4.3 与Obsidian的集成

```
Obsidian Vault ←── Git同步 ──→ Wiki.js ←── API ──→ RDP平台
     │                                              │
     │         双向同步：                            │
     │         - Markdown文件                        │
     │         - Wiki链接                            │
     │         - 附件                                │
     │                                              │
     └────────── 标签关联 ──────────────────────────┘
```

### A.5 Mattermost 集成方案

#### A.5.1 部署配置

```yaml
# /etc/systemd/system/mattermost.service
[Unit]
Description=Mattermost Team Collaboration
After=network.target postgresql.service

[Service]
Type=notify
User=rdp-user
Group=rdp-user
WorkingDirectory=/opt/rdp/mattermost
ExecStart=/opt/rdp/mattermost/bin/mattermost
Restart=always
RestartSec=5
MemoryMax=2G

[Install]
WantedBy=multi-user.target
```

#### A.5.2 SSO集成配置

```json
// Mattermost SSO配置
{
  "Enable": true,
  "Secret": "casdoor-client-secret",
  "Id": "casdoor-client-id",
  "Scope": "openid profile email",
  "AuthEndpoint": "https://rdp.local/auth/login",
  "TokenEndpoint": "https://rdp.local/auth/token",
  "UserApiEndpoint": "https://rdp.local/auth/userinfo"
}
```

### A.6 Gitea 集成方案

#### A.6.1 部署配置

```yaml
# /etc/systemd/system/gitea.service
[Unit]
Description=Gitea Git Service
After=network.target postgresql.service

[Service]
Type=simple
User=rdp-user
Group=rdp-user
WorkingDirectory=/opt/rdp/data/gitea
ExecStart=/opt/rdp/bin/gitea web -c /opt/rdp/config/gitea/app.ini
Restart=always
RestartSec=5
MemoryMax=2G

[Install]
WantedBy=multi-user.target
```

#### A.6.2 API集成示例

```go
// 创建项目时自动创建Git仓库
func CreateProjectRepository(project *Project) error {
    giteaClient := gitea.NewClient("https://git.local", "token")
    
    _, _, err := giteaClient.CreateRepo(gitea.CreateRepoOption{
        Name:        project.ProjectCode,
        Description: project.ProjectName,
        Private:     true,
        AutoInit:    true,
    })
    
    return err
}
```

### A.7 统一认证（SSO）方案

#### A.7.1 Casdoor配置

```ini
# casdoor配置
appname = casdoor
httpport = 8000
runmode = prod

[database]
adapter = postgres
host = 127.0.0.1
port = 5432
user = rdp
password = your_password
database = casdoor

[oauth]
clientId = your-client-id
clientSecret = your-client-secret
```

#### A.7.2 应用系统集成

| 应用系统 | 集成协议 | 回调地址 |
|----------|----------|----------|
| RDP主平台 | OAuth2 | https://rdp.local/auth/callback |
| Mattermost | OAuth2 | https://rdp.local/chat/oauth/callback |
| Gitea | OAuth2 | https://rdp.local/git/user/oauth2/casdoor/callback |
| Wiki.js | OIDC | https://rdp.local/wiki/login/casdoor/callback |

