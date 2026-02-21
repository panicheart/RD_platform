# Phase 1 整合计划

## 整合负责人
**InfraAgent** - 项目骨架 + 测试框架 + 协作规范

## 待整合的Agent

| Agent | 任务ID | 预计交付 | 冲突风险 |
|-------|--------|----------|----------|
| PortalAgent | P1-T1~T4 | 待定 | 低 |
| UserAgent | P1-T5~T8 | 待定 | 中 (types, handlers) |
| ProjectAgent | P1-T9~T12 | 待定 | 中 (types, handlers) |
| SecurityAgent | P1-T13~P1-T16 | 待定 | 低 |

## 当前状态

### InfraAgent 已完成 (已本地提交)
- ✅ 项目目录结构
- ✅ 前端脚手架 (React + TS + Vite)
- ✅ 后端脚手架 (Go + Gin)
- ✅ 数据库配置 (PostgreSQL enums)
- ✅ 部署配置 (systemd + Nginx)
- ✅ 测试框架 (Vitest + testify)
- ✅ CI/CD (GitHub Actions)
- ✅ 协作文档 (CHECKLIST, FAQ, COLLABORATION_GUIDE)

### 等待其他Agent
- ⏳ PortalAgent - 门户界面
- ⏳ UserAgent - 用户管理
- ⏳ ProjectAgent - 项目管理
- ⏳ SecurityAgent - 安全合规

## 潜在冲突点

### 1. 文件路径冲突
```
风险区域:
├── apps/web/src/types/index.ts          # 所有Agent都可能扩展
├── apps/web/src/services/api.ts         # 各Agent添加API
├── services/api/handlers/*.go           # 各Agent创建handler
├── services/api/models/*.go             # 各Agent创建model
└── database/migrations/*.sql            # 迁移文件序号
```

**解决策略**:
- `types/index.ts`: 可追加，不要修改已有类型
- `services/api.ts`: 各Agent创建独立的service文件
- `handlers/*.go`: 按模块命名，如 `user.go`, `project.go`
- `migrations/*.sql`: 使用序号递增，InfraAgent已预留 001-004

### 2. 依赖版本冲突
```
风险: package.json, go.mod 版本不一致
解决: 以InfraAgent设定的版本为基准
```

### 3. API路由冲突
```
风险: 多个Agent定义相同路由
解决: 按模块划分:
  - /api/v1/users/*     -> UserAgent
  - /api/v1/projects/*  -> ProjectAgent
  - /api/v1/auth/*      -> UserAgent (auth)
  - /api/v1/files/*     -> ProjectAgent (files)
  - /api/v1/security/*  -> SecurityAgent
```

## 整合步骤

### Step 1: 收集所有Agent工作
```bash
# 检查各Agent分支
 git branch -a
# 或检查远程分支
 git fetch origin
```

### Step 2: 创建整合分支
```bash
git checkout -b phase1-integration
```

### Step 3: 逐个合并Agent代码
```bash
# 假设各Agent在各自分支
git merge portal-agent-phase1
git merge user-agent-phase1
git merge project-agent-phase1
git merge security-agent-phase1
```

### Step 4: 解决冲突
冲突类型预判:
1. **package.json** - 合并依赖
2. **go.mod** - 合并依赖
3. **types/index.ts** - 保留所有类型定义
4. **App.tsx** - 合并路由配置
5. **main.go** - 合并路由注册

### Step 5: 验证整合
```bash
# 前端验证
cd apps/web
npm install
npm run typecheck
npm run lint
npm run test

# 后端验证
cd services/api
go mod tidy
go build ./...
go test ./...

# 全量验证
make lint
make test
```

### Step 6: 提交整合结果
```bash
git commit -m "integrate: merge all Phase 1 Agent work

- PortalAgent: portal pages and components
- UserAgent: user management and auth
- ProjectAgent: project CRUD and templates
- SecurityAgent: classification and audit
- InfraAgent: scaffold and testing framework

Co-authored-by: PortalAgent
Co-authored-by: UserAgent
Co-authored-by: ProjectAgent
Co-authored-by: SecurityAgent"
```

## 整合检查清单

### 代码检查
- [ ] 所有Agent代码已合并
- [ ] 无编译错误
- [ ] 无lint错误
- [ ] 所有测试通过
- [ ] 覆盖率达标 (60%+)

### 功能检查
- [ ] 前端可正常启动
- [ ] 后端可正常启动
- [ ] API健康检查通过
- [ ] 数据库迁移可执行
- [ ] 基本CRUD操作正常

### 文档检查
- [ ] 各Agent交付文档已更新
- [ ] API文档已更新
- [ ] README已更新
- [ ] WORKSPACE_REGISTRY已更新

## 整合时间线

| 阶段 | 预计时间 | 依赖 |
|------|----------|------|
| Agent提交 | 等待中 | Phase 1 Agent完成 |
| 冲突解决 | 2-4小时 | Agent提交后 |
| 验证测试 | 1-2小时 | 冲突解决后 |
| 最终提交 | 30分钟 | 验证通过后 |

## 问题上报

如果整合过程中遇到无法解决的问题:
1. 在 `CONFLICTS.md` 中记录
2. 联系相关Agent协商
3. 上报 PM-Agent 或 Architect Agent
4. 必要时上报人类监督者

---

*整合计划版本: 1.0*
*创建时间: 2026-02-22*
*最后更新: 2026-02-22*
