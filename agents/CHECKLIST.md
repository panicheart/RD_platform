# Agent 预检清单

## 强制阅读清单 ⭐

**每个Agent启动前必须阅读，并在下方签字确认：**

| 文档 | 必须阅读 | 已确认 |
|------|----------|--------|
| [QUICKSTART.md](../QUICKSTART.md) | ✅ | □ |
| [COLLABORATION_GUIDE.md](../outputs/COLLABORATION_GUIDE.md) | ✅ | □ |
| [CONFLICT_WARNING.md](../outputs/infra-scaffold/CONFLICT_WARNING.md) | ✅ | □ |
| [DELIVERY.md](../outputs/infra-scaffold/DELIVERY.md) | ✅ | □ |

**签字确认方式**：在 [Agent启动验证表](#agent启动验证表) 中勾选

---

## Agent 启动验证表

### PortalAgent (P1-T1 ~ P1-T4)

**前置条件：**
- [ ] 已阅读协作规范
- [ ] 已检查冲突预警表
- [ ] 已确认工作目录无冲突

**启动命令：**
```bash
./scripts/validate-agent-startup.sh PortalAgent P1-T1
cd apps/web && npm install && npm run dev
```

**Agent签字：** `_________________` 日期：`_______`

---

### UserAgent (P1-T5 ~ P1-T8)

**前置条件：**
- [ ] 已阅读协作规范
- [ ] 已检查冲突预警表
- [ ] 已确认工作目录无冲突

**启动命令：**
```bash
./scripts/validate-agent-startup.sh UserAgent P1-T5
cd services/api && go mod tidy && go run main.go
```

**Agent签字：** `_________________` 日期：`_______`

---

### ProjectAgent (P1-T9 ~ P1-T12)

**前置条件：**
- [ ] 已阅读协作规范
- [ ] 已检查冲突预警表
- [ ] 已确认工作目录无冲突

**启动命令：**
```bash
./scripts/validate-agent-startup.sh ProjectAgent P1-T9
cd services/api && go mod tidy && go run main.go
```

**Agent签字：** `_________________` 日期：`_______`

---

### SecurityAgent (P1-T13 ~ P1-T16)

**前置条件：**
- [ ] 已阅读协作规范
- [ ] 已检查冲突预警表
- [ ] 已确认工作目录无冲突

**启动命令：**
```bash
./scripts/validate-agent-startup.sh SecurityAgent P1-T13
cd services/api && go mod tidy && go run main.go
```

**Agent签字：** `_________________` 日期：`_______`

---

## 强制检查点

### 检查点 1: 代码提交前
每个Agent在提交代码前必须确认：

- [ ] 已创建/更新 `agents/outputs/{agent-name}/DELIVERY.md`
- [ ] 已检查是否与其他Agent有文件冲突
- [ ] 代码通过 lint 检查 (`make lint`)
- [ ] 已阅读冲突预警表的最新版本

### 检查点 2: 每日同步
每个Agent每天结束工作前：

- [ ] 已提交代码到git
- [ ] 已更新 CHANGELOG.md
- [ ] 已检查其他Agent的交付文档更新
- [ ] 已更新依赖关系图（如有变更）

### 检查点 3: Phase 完成前
Phase内所有Agent在Phase完成前：

- [ ] 所有交付文档已更新
- [ ] 所有冲突已解决
- [ ] 已运行集成测试
- [ ] 已提交给 Reviewer Agent

---

## 不遵守的后果

如果Agent未阅读文档就开始工作：

1. **发现冲突** → 必须重构代码（时间损失）
2. **违反规范** → Reviewer Agent 打回修改
3. **覆盖他人代码** → 承担合并责任
4. **严重违规** → 上报人类监督者裁决

---

## 快速 FAQ

**Q: 我可以用脚本自动检查吗？**  
A: 可以！运行 `./scripts/validate-agent-startup.sh {AgentName} {TaskID}`

**Q: 如果我发现文档有误怎么办？**  
A: 在 `agents/outputs/infra-scaffold/CONFLICTS.md` 中记录，并通知 InfraAgent

**Q: 我可以跳过检查吗？**  
A: 技术上可以，但如果因此产生冲突，你需要承担修复责任

**Q: 如何知道其他Agent更新了什么？**  
A: 每天查看 `git log` 和 `agents/outputs/*/DELIVERY.md`

---

*此清单强制执行 - 请在开始工作前确认所有检查项*
