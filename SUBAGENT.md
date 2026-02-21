# Subagent 使用规范

> **文档用途**: 定义Subagent的启动、执行、交付规范，用于指导AI Agent集群开发
> **版本**: V1.0
> **最后更新**: 2026-02-21

---

## 1. Subagent 概述

### 1.1 什么是 Subagent

Subagent 是 AI Agent 集群中的**功能执行单元**，负责具体模块的开发实现。

**核心特点**:
- **独立上下文**: 每个Subagent只关注自己的模块，避免上下文混乱
- **并行执行**: 多个Subagent可同时运行，最大化开发效率
- **错误隔离**: 一个Subagent失败不影响其他Subagent
- **质量专注**: 每个Subagent可深入优化自己的模块

### 1.2 Subagent 与主 Agent 的关系

```
人类监督者
    │
    ▼
协调层 Agent (Architect / PM-Agent / Reviewer)
    │ 分解任务 / 分配指令
    ▼
功能层 Subagent (PortalAgent / UserAgent / ...)
    │ 执行开发 / 自我审查
    ▼
代码仓库 / 交付物
```

---

## 2. Subagent 标准指令格式

### 2.1 启动指令模板

每个Subagent启动时必须接收以下标准指令：

```yaml
# Subagent 启动指令
subagent:
  name: "UserAgent"                    # Subagent名称
  role: "feature_developer"            # 角色类型
  version: "1.0"                       # 指令版本

task:
  module: "用户管理模块"                # 负责模块
  phase: 1                             # 所属Phase
  priority: "P0"                       # 优先级
  
  # 任务清单（必须可量化验收）
  tasks:
    - id: "TASK-001"
      name: "用户CRUD API"
      description: "实现用户的增删改查API"
      acceptance_criteria:
        - "支持用户列表分页查询"
        - "支持用户详情获取"
        - "支持用户信息修改"
        - "支持用户软删除"
      
    - id: "TASK-002"
      name: "RBAC权限模型"
      description: "基于Casbin实现RBAC权限控制"
      acceptance_criteria:
        - "支持5级角色定义"
        - "支持资源级权限控制"
        - "支持数据权限隔离"
        
  # 依赖项
  dependencies:
    - name: "InfraAgent"
      artifact: "数据库设计"
      status: "completed"  # 必须是completed才能启动
      
    - name: "Architect Agent"
      artifact: "接口规范"
      status: "completed"

tech_stack:
  # 技术栈约束（必须严格遵守）
  backend:
    language: "Go"
    framework: "Gin"
    orm: "GORM"
    version: "1.22+"
    
  database:
    type: "PostgreSQL"
    version: "16.x"
    
  external:
    - name: "Casdoor"
      purpose: "统一认证"
      version: "latest"
    - name: "Casbin"
      purpose: "权限管理"
      version: "2.x"

constraints:
  # 强制约束
  must:
    - "纯Go后端，禁止引入Node.js"
    - "代码注释使用英文，UI文案使用中文"
    - "API路径格式: /api/v1/{resource}"
    - "错误码格式: 6xxx系列业务码"
    - "ID生成使用ULID，禁止自增ID"
    
  # 推荐约束
  should:
    - "单元测试覆盖率≥60%"
    - "函数长度不超过50行"
    - "每个API都有Swagger文档"
    
  # 可选约束
  may:
    - "使用缓存优化热点数据"
    - "添加Prometheus监控指标"

deliverables:
  # 交付物清单
  code:
    - path: "services/user/"
      description: "用户服务源代码"
      include:
        - "*.go"
        - "go.mod"
        - "go.sum"
        
  docs:
    - path: "docs/api/user.md"
      description: "API接口文档"
      format: "markdown"
      
    - path: "docs/architecture/user.md"
      description: "架构设计文档"
      format: "markdown"
      
  tests:
    - path: "services/user/*_test.go"
      description: "单元测试"
      coverage_requirement: 60
      
  config:
    - path: "deploy/systemd/user.service"
      description: "systemd服务配置"
      
  migration:
    - path: "database/migrations/001_user.up.sql"
      description: "数据库迁移脚本"

review:
  # 审查配置
  self_review: true                    # 必须先自审查
  reviewer: "Reviewer Agent"           # 审查Agent
  criteria:
    - "代码规范检查通过"
    - "单元测试全部通过"
    - "接口符合Architect Agent定义"
    - "无安全漏洞"

reporting:
  # 汇报配置
  progress_report: "every_20_percent"  # 每完成20%汇报
  blockers: "immediate"                # 阻塞问题立即上报
  completion: "detailed"               # 完成时详细报告
```

### 2.2 指令验证清单

启动Subagent前，必须验证以下项：

| 检查项 | 验证内容 | 不通过处理 |
|--------|----------|-----------|
| **依赖检查** | 所有依赖项状态为completed | 等待依赖完成 |
| **规范检查** | Architect Agent的接口规范已发布 | 等待规范发布 |
| **环境检查** | 开发环境已准备就绪 | InfraAgent准备环境 |
| **权限检查** | 代码仓库写入权限已配置 | 配置权限 |

---

## 3. Subagent 执行流程

### 3.1 标准开发流程

```
接收指令
    │
    ▼
┌─────────────────┐
│ 1. 需求理解     │ ← 仔细阅读任务清单和验收标准
│                 │    输出：需求理解确认
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 2. 技术设计     │ ← 设计数据模型、API接口、实现方案
│                 │    输出：技术设计文档（简要）
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 3. 编码实现     │ ← 按规范编写代码
│                 │    输出：源代码
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 4. 自审查(L1)   │ ← 代码规范、语法、基础逻辑检查
│                 │    输出：自审查清单
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 5. 单元测试     │ ← 编写并执行单元测试
│                 │    输出：测试报告（覆盖率≥60%）
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 6. 提交审查     │ ← 提交给Reviewer Agent
│                 │    输出：审查请求
└────────┬────────┘
         │
         ▼
    ┌────┴────┐
    ▼         ▼
 通过        不通过
  │            │
  ▼            ▼
交付完成    修改后重新提交
```

### 3.2 进度汇报模板

**每完成20%进度时，Subagent必须输出**：

```markdown
## Subagent 进度报告

**Subagent**: UserAgent  
**模块**: 用户管理  
**Phase**: 1  
**报告时间**: 2026-02-21 14:30  
**总体进度**: 40% (2/5 任务完成)

### 已完成任务
| 任务ID | 任务名称 | 完成时间 | 备注 |
|--------|----------|----------|------|
| TASK-001 | 用户CRUD API | 2026-02-21 10:00 | 已通过自审查 |
| TASK-002 | 数据库模型 | 2026-02-21 14:00 | 含迁移脚本 |

### 进行中任务
| 任务ID | 任务名称 | 进度 | 预计完成 |
|--------|----------|------|----------|
| TASK-003 | RBAC权限模型 | 20% | 2026-02-21 18:00 |

### 阻塞问题
| 问题 | 严重程度 | 描述 | 需要帮助 |
|------|----------|------|----------|
| 无 | - | - | - |

### 下一步计划
1. 完成RBAC权限模型实现（预计4小时）
2. 编写单元测试（预计2小时）
3. 提交Reviewer Agent审查

### 风险评估
- **进度风险**: 低，按计划推进
- **技术风险**: 低，Casbin文档完善
```

### 3.3 阻塞处理流程

当Subagent遇到无法解决的问题时：

```
遇到阻塞问题
    │
    ├──→ 尝试自主解决（30分钟）
    │       │
    │       ▼
    │   解决成功 → 继续开发
    │       │
    │       ▼
    │   无法解决
    │       │
    ▼       ▼
上报PM-Agent
    │
    ├──→ PM-Agent评估
    │       │
    │       ▼
    │   可以协调 → PM-Agent协调解决
    │       │
    │       ▼
    │   需要决策 → 上报人类监督者
    │       │
    ▼       ▼
你（监督者）决策
```

---

## 4. Subagent 交付标准

### 4.1 代码交付标准

| 检查项 | 标准 | 检查方式 |
|--------|------|----------|
| **代码规范** | 符合项目编码规范 | 自动化工具检查 |
| **编译通过** | `go build` 无错误 | 编译检查 |
| **单元测试** | 覆盖率≥60%，全部通过 | `go test` |
| **接口符合** | 符合Architect Agent定义的接口 | 接口对比 |
| **错误处理** | 统一错误码，正确处理异常 | 代码审查 |
| **安全** | 无SQL注入、XSS等漏洞 | 安全扫描 |

### 4.2 文档交付标准

| 文档类型 | 必须包含 | 格式 |
|----------|----------|------|
| **API文档** | 所有接口的Request/Response示例 | Markdown |
| **架构文档** | 模块架构图、数据流图 | Markdown + 图 |
| **部署文档** | 部署步骤、配置说明 | Markdown |
| **README** | 模块简介、快速开始 | Markdown |

### 4.3 自审查清单（L1）

Subagent提交前必须完成以下自检：

```markdown
## 自审查清单 (L1)

**Subagent**: UserAgent  
**审查日期**: 2026-02-21  
**审查人**: UserAgent自身

### 代码规范
- [x] 代码注释使用英文
- [x] 变量名使用camelCase
- [x] 常量使用UPPER_SNAKE_CASE
- [x] 函数名使用动词+名词格式
- [x] 每行不超过120字符
- [x] 函数长度不超过50行

### 功能实现
- [x] 所有任务清单项已完成
- [x] 验收标准全部满足
- [x] 边界情况已处理
- [x] 错误处理完善

### 测试
- [x] 单元测试覆盖率≥60%
- [x] 所有测试用例通过
- [x] 包含边界条件测试
- [x] 包含错误情况测试

### 文档
- [x] API文档已更新
- [x] 架构文档已更新
- [x] README已更新

### 审查结论
- [x] 通过自审查，可以提交Reviewer Agent

签名: UserAgent
```

---

## 5. Subagent 间协作规范

### 5.1 依赖管理

当Subagent A 依赖 Subagent B 的输出时：

```yaml
# Subagent A 的依赖声明
dependencies:
  - name: "Subagent B"
    artifact: "API接口定义"
    required_by: "2026-02-21 12:00"  # 截止时间
    contract: "interface_v1.yaml"      # 契约文件
    
# PM-Agent 协调
协调动作:
  - 监控Subagent B进度
  - 如可能延期，提前预警
  - 如契约变更，通知Subagent A
```

### 5.2 接口契约

Subagent间通过**接口契约**进行协作：

```yaml
# interface_v1.yaml - 由Architect Agent定义
api_version: "v1"
module: "user"

endpoints:
  - path: "/api/v1/users"
    method: "GET"
    request:
      query:
        - name: "page"
          type: "integer"
          default: 1
        - name: "page_size"
          type: "integer"
          default: 20
    response:
      success:
        code: 200
        data:
          list: "[]User"
          total: "integer"
      error:
        code: "400|401|403|500"
        message: "string"
        
  - path: "/api/v1/users/:id"
    method: "GET"
    request:
      params:
        - name: "id"
          type: "string"
          format: "ULID"
    response:
      success:
        code: 200
        data: "User"
```

**契约变更流程**：
1. Architect Agent提出变更
2. PM-Agent评估影响范围
3. 通知所有依赖方
4. 人类监督者确认（如涉及重大变更）
5. 同步更新契约版本

### 5.3 冲突解决

当两个Subagent对接口定义有不同意见时：

```
Subagent A vs B 接口分歧
    │
    ├──→ 各自陈述理由
    │
    ▼
Architect Agent技术评估
    │
    ├──→ 选择最优方案
    │
    ▼
双方Subagent确认
    │
    ▼
更新接口契约
```

**升级条件**：Architect Agent无法决策时，上报人类监督者

---

## 6. Subagent 工具集

### 6.1 推荐工具

| 用途 | 工具 | 说明 |
|------|------|------|
| **代码生成** | `swagger-codegen` | 从API契约生成代码框架 |
| **数据库** | `gormigrate` | 数据库迁移管理 |
| **测试** | `testify` | Go单元测试框架 |
| **Mock** | `mockery` | 生成Mock对象 |
| **Lint** | `golangci-lint` | 代码规范检查 |
| **文档** | `swag` | 自动生成Swagger文档 |

### 6.2 常用命令

```bash
# 代码检查
golangci-lint run ./...

# 单元测试
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 构建
go build -o bin/user-service ./services/user

# 数据库迁移
go run ./cmd/migrate up

# 生成Swagger文档
swag init -g main.go
```

---

## 7. 示例：UserAgent 完整执行记录

### 7.1 启动指令接收

```yaml
# 接收到的启动指令
subagent: UserAgent
phase: 1
tasks: [用户CRUD, RBAC权限, 组织架构]
dependencies: [数据库设计✓, 接口规范✓]
status: ready_to_start
```

### 7.2 执行过程

```
[T+0h] 接收指令，确认需求
[T+1h] 完成技术设计
[T+3h] 完成数据库模型
[T+5h] 完成用户CRUD API
[T+6h] 输出进度报告（40%）
[T+10h] 完成RBAC权限
[T+12h] 完成组织架构
[T+14h] 完成单元测试（覆盖率75%）
[T+15h] 完成自审查
[T+16h] 提交Reviewer Agent
[T+20h] Reviewer审查通过
[T+21h] 交付完成
```

### 7.3 交付物

```
services/user/
├── main.go                 # 入口文件
├── handlers/
│   ├── user_handler.go     # HTTP处理器
│   └── auth_handler.go     # 认证处理器
├── services/
│   ├── user_service.go     # 业务逻辑
│   └── auth_service.go     # 权限逻辑
├── models/
│   └── user.go             # 数据模型
├── repositories/
│   └── user_repo.go        # 数据访问
├── middleware/
│   └── auth_middleware.go  # 认证中间件
├── *_test.go               # 单元测试
└── go.mod                  # 依赖管理

docs/
├── api/user.md             # API文档
└── architecture/user.md    # 架构文档

deploy/systemd/
└── user.service            # 服务配置

database/migrations/
├── 001_user.up.sql         # 升级脚本
└── 001_user.down.sql       # 回滚脚本
```

---

## 8. 附录

### 8.1 Subagent 命名规范

| 类型 | 命名格式 | 示例 |
|------|----------|------|
| **功能Subagent** | {Module}Agent | UserAgent, ProjectAgent |
| **协调Subagent** | {Role} Agent | PM-Agent, Reviewer Agent |
| **任务ID** | {MODULE}-{SEQUENCE} | USER-001, PROJ-003 |

### 8.2 状态定义

```yaml
# Subagent状态
subagent_status:
  - pending           # 等待启动
  - running           # 执行中
  - blocked           # 阻塞等待
  - reviewing         # 审查中
  - completed         # 已完成
  - failed            # 失败

# 任务状态
task_status:
  - todo              # 待开始
  - in_progress       # 进行中
  - completed         # 已完成
  - blocked           # 阻塞
```

---

*Subagent 使用规范 V1.0*  
*© 2026 微波研发部门 | 内部文档*
