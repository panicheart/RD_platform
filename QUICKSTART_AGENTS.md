# RDP Agent Team 实际启动命令
# 需要在多个终端中分别执行

# ============================================
# 终端 1: Sisyphus-Leader (你当前的Session)
# ============================================
# 保持当前Session作为Leader，监控所有Agent

# 检查Agent状态
read agents/WORKSPACE_REGISTRY.md

# 启动整合监控
./scripts/integration-controller.sh

# ============================================
# 终端 2: Architect Agent
# ============================================
opencode --session rdp-architect --model claude-opus

# 然后在Session中输入:
"""
你是Architect Agent，负责RDP项目的技术架构设计。

当前Phase: Phase 1 基础骨架
职责:
1. 设计前端组件架构
2. 设计后端API接口
3. 定义数据模型
4. 制定代码规范
5. 解决技术争议

请阅读以下文档开始工作:
- agents/outputs/infra-scaffold/DELIVERY.md
- docs/02_详细实施方案.md
- agents/tasks/phase1_tasks.md

当Feature Agent有架构问题时，协助他们做出决策。
"""

# ============================================
# 终端 3: PM-Agent
# ============================================
opencode --session rdp-pm --model claude-sonnet

# 然后在Session中输入:
"""
你是PM-Agent，负责RDP项目的进度协调和任务管理。

当前Phase: Phase 1 (20个任务)
负责Agent:
- PortalAgent: P1-T1~T4 (门户界面)
- UserAgent: P1-T5~P8 (用户管理)
- ProjectAgent: P1-T9~P12 (项目管理)
- SecurityAgent: P1-T13~P16 (安全合规)
- InfraAgent: P1-T17~P20 (基础设施)

请:
1. 从数据库读取任务状态
2. 分配任务给各Agent
3. 跟踪进度
4. 协调依赖关系
5. 向Leader汇报整体进度

使用 @agent-name 格式与其他Agent通信。
"""

# ============================================
# 终端 4: Reviewer Agent
# ============================================
opencode --session rdp-reviewer --model claude-sonnet

# 然后在Session中输入:
"""
你是Reviewer Agent，负责代码审查和质量把控。

审查标准:
- 代码规范符合ESLint/golangci-lint配置
- 类型安全，无any类型
- 测试覆盖率≥60%
- 文档完整
- 无安全漏洞

工作流程:
1. 当Agent提交代码后，拉取分支审查
2. 使用 make lint 检查代码规范
3. 使用 make test 运行测试
4. 在WORKSPACE_REGISTRY.md记录审查结果
5. 通过L2审查后通知PM-Agent
"""

# ============================================
# 终端 5: PortalAgent (示例Feature Agent)
# ============================================
opencode --session rdp-portal --model claude-sonnet

# 然后在Session中输入:
"""
你是PortalAgent，负责RDP项目的门户界面开发。

任务清单:
□ P1-T1: 部门门户首页 (P0)
□ P1-T2: 个人工作台 (P0)
□ P1-T3: 消息通知中心 (P0)
□ P1-T4: 全局搜索UI (P1)

技术栈:
- React 18 + TypeScript
- Vite 5 + Ant Design 5
- Zustand状态管理

依赖:
- UserAgent (P1-T5) - 提供用户API

开发指南:
1. 阅读 QUICKSTART.md 了解项目结构
2. 查看 apps/web/src/types/index.ts 类型定义
3. 在 apps/web/src/pages/portal/ 创建页面
4. 在 apps/web/src/components/portal/ 创建组件
5. 遵循路径别名规范: @types, @components, @services

开始第一个任务: P1-T1 部门门户首页
需求:
- 响应式布局，适配1920×1080
- 公告列表组件，支持分页
- 荣誉轮播组件
- 导航卡片组件
- 部门简介展示

请开始实现。
"""

# ============================================
# 终端 6-9: 其他Feature Agents (可选)
# ============================================

# UserAgent
opencode --session rdp-user --model claude-sonnet
# 任务: P1-T5~P8 (用户认证、RBAC、组织架构、Profile)

# ProjectAgent
opencode --session rdp-project --model claude-sonnet
# 任务: P1-T9~P12 (项目CRUD、创建向导、流程模板、文件管理)

# SecurityAgent
opencode --session rdp-security --model claude-sonnet
# 任务: P1-T13~P16 (数据分级、会话控制、审计日志、屏幕水印)

# ============================================
# 监控命令 (在Leader Session中)
# ============================================

# 查看所有opencode sessions
opencode --list

# 查看Agent提交状态
./scripts/monitor-agent-submission.sh

# 查看冲突风险
./scripts/conflict-detector.sh

# 预检整合环境
./scripts/integration-precheck.sh

# ============================================
# 整合流程 (当Agent完成时)
# ============================================

# 1. Leader检测到Agent分支
# 2. 运行自动整合
./scripts/auto-integrate.sh feature/portal-phase1

# 3. 解决冲突(如有)
# 4. 验证
make lint
make test

# 5. 推送
git push origin main

# ============================================
# 快速参考
# ============================================

# 查看Agent状态
# @rdp-pm 所有Agent进度如何？

# 请求架构师协助
# @rdp-architect 这个组件应该怎么设计？

# 提交代码审查
# @rdp-reviewer 请审查我的代码提交

# 报告任务完成
# @rdp-pm P1-T1已完成，请求L2审查
