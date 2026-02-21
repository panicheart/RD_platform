# PM-Agent 启动指令模板

## 启动 Phase 1 Agent 的标准流程

### Step 1: 确认前提条件

```bash
# 1.1 验证项目骨架
./scripts/validate-agent-startup.sh PM-Agent P0

# 1.2 确保所有文档已提交
git add agents/outputs/
git commit -m "chore: add collaboration docs and scaffold"

# 1.3 更新主README
git add README.md QUICKSTART.md
git commit -m "docs: update README with quickstart guide"
```

### Step 2: 广播启动通知

向所有 Phase 1 Agent 发送启动指令：

---

## 📢 启动指令: Phase 1 Agent

**致: PortalAgent, UserAgent, ProjectAgent, SecurityAgent**

项目骨架已搭建完成！请按以下步骤启动：

### 必须执行（强制）

1. **阅读文档**（5分钟）:
   - [QUICKSTART.md](../../QUICKSTART.md)
   - [agents/CHECKLIST.md](../../CHECKLIST.md)
   - [agents/outputs/COLLABORATION_GUIDE.md](../../outputs/COLLABORATION_GUIDE.md)
   - [agents/outputs/infra-scaffold/CONFLICT_WARNING.md](../../outputs/infra-scaffold/CONFLICT_WARNING.md)

2. **运行启动验证**:
   ```bash
   ./scripts/validate-agent-startup.sh {YourAgentName} {YourTaskID}
   ```

3. **创建工作分支**:
   ```bash
   git checkout -b feature/{agent-name}-phase1
   ```

4. **创建交付文档框架**:
   ```bash
   mkdir -p agents/outputs/{agent-name}
   touch agents/outputs/{agent-name}/DELIVERY.md
   touch agents/outputs/{agent-name}/CHANGELOG.md
   ```

### 然后开始开发

根据你的任务卡片 (agents/tasks/phase1_tasks.md) 开始实现。

### 每日必须

- 更新 `CHANGELOG.md`
- 提交代码 `git commit`
- 检查其他Agent的交付文档

---

## Step 3: 并行启动Agent

```bash
# 并行启动所有Phase 1 Agent（无依赖关系）

task(agent="PortalAgent", ...)
task(agent="UserAgent", ...)
task(agent="ProjectAgent", ...)
task(agent="SecurityAgent", ...)
```

### 启动参数模板

```yaml
agent_role: "PortalAgent"
module: "门户界面"
tasks: ["P1-T1", "P1-T2", "P1-T3", "P1-T4"]
prerequisites:
  - "Read QUICKSTART.md"
  - "Read COLLABORATION_GUIDE.md"
  - "Read CONFLICT_WARNING.md"
  - "Run validate-agent-startup.sh"
deliverables:
  - "agents/outputs/portal-agent/DELIVERY.md"
  - "agents/outputs/portal-agent/CHANGELOG.md"
  - "apps/web/src/pages/portal/*"
  - "apps/web/src/components/portal/*"
reviewer: "Reviewer Agent"
```

---

## Step 4: 监控进度

### 每日检查清单

- [ ] 检查所有Agent是否提交了代码
- [ ] 检查是否有冲突报告
- [ ] 更新项目进度看板
- [ ] 向人类监督者汇报进度

### 冲突处理流程

```
Agent A 报告冲突
    ↓
PM-Agent 记录冲突
    ↓
尝试协调 (24小时内)
    ↓
协调成功?
    ├─ 是 → 记录解决方案
    └─ 否 → 上报 Architect Agent
```

---

## Step 5: Phase 1 验收

### 验收检查清单

- [ ] 所有Agent交付文档完整
- [ ] 所有代码通过L2审查
- [ ] 集成测试通过
- [ ] 无未解决的冲突
- [ ] 人类监督者验收签字

---

*PM-Agent 执行此流程启动 Phase 1*
