# 协作问题快速索引

## 🔴 常见问题

### Q1: 我应该从哪里开始？
**A**: 
1. 阅读 [QUICKSTART.md](../QUICKSTART.md)
2. 运行 `./scripts/validate-agent-startup.sh {AgentName} {TaskID}`
3. 查看 [CHECKLIST.md](../CHECKLIST.md)

### Q2: 如何知道其他Agent在做什么？
**A**: 
- 查看 [WORKSPACE_REGISTRY.md](../WORKSPACE_REGISTRY.md)
- 查看 `agents/outputs/*/DELIVERY.md`
- 运行 `git log` 查看提交记录

### Q3: 我发现我的任务和其他Agent有重叠怎么办？
**A**:
1. 先不要开始编码
2. 在 [CONFLICT_WARNING.md](../outputs/infra-scaffold/CONFLICT_WARNING.md) 中查看是否已声明
3. 在 [WORKSPACE_REGISTRY.md](../WORKSPACE_REGISTRY.md) 中 @ 相关Agent
4. 协商不成？上报 Architect Agent 或人类监督者

### Q4: 我能否修改别人创建的文件？
**A**:
- 配置文件 (package.json, go.mod): **不可以**，需协商
- 共享代码 (types/index.ts, handlers/response.go): **追加可以，修改不可以**
- 别人工作区的文件: **绝对不可以**

### Q5: 规范太严格了，我能改吗？
**A**:
- 如果是小问题（如函数命名），遵循现有规范
- 如果是大问题（如架构设计），在 `CONFLICTS.md` 中提出
- 重要变更需 Reviewer Agent 或人类监督者批准

### Q6: 我忘记了文档在哪，怎么办？
**A**: 看这个列表：
- [项目总览](../../README.md)
- [快速开始](../QUICKSTART.md)
- [协作规范](../outputs/COLLABORATION_GUIDE.md)
- [冲突预警](../outputs/infra-scaffold/CONFLICT_WARNING.md)
- [任务总览](agent_overview.md)

---

## 🟡 技术问题

### Q: 前端编译报错
**A**: 
```bash
cd apps/web
rm -rf node_modules package-lock.json
npm install
npm run dev
```

### Q: 后端编译报错
**A**:
```bash
cd services/api
go mod tidy
go run main.go
```

### Q: 数据库连接失败
**A**: 检查 `config/rdp-api.yaml` 中的数据库配置

---

## 🟢 流程问题

### Q: 我什么时候可以开始Phase 2？
**A**: 检查 [WORKSPACE_REGISTRY.md](../WORKSPACE_REGISTRY.md)，等待所有Phase 1 Agent标记为"已完成"

### Q: 我的代码什么时候会被审查？
**A**: 
1. 自审查 (L1): 你自己完成
2. 代码审查 (L2): 提交给 Reviewer Agent
3. 集成测试 (L3): PM-Agent 组织
4. 人类验收 (L4): 人类监督者

### Q: 发现Bug但不是我写的代码，怎么办？
**A**:
1. 在 [WORKSPACE_REGISTRY.md](../WORKSPACE_REGISTRY.md) 中记录
2. @ 相关Agent
3. 如果紧急，直接修复并记录变更

---

## 📞 上报流程

当问题无法解决时，按以下顺序上报：

```
1. 相关Agent之间协商
        ↓
2. 查阅文档 (FAQ, 协作规范)
        ↓
3. 在 CONFLICTS.md 中记录
        ↓
4. 请求 Architect Agent 裁决
        ↓
5. 上报人类监督者 (最终决策)
```

---

*最后更新: 2026-02-22*
