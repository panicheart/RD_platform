# 多Agent协作规范

## 1. 交付物文档格式

每个Agent完成任务后，必须在 `agents/outputs/{agent-name}/` 目录下创建交付文档:

```
agents/outputs/
├── portal-agent/
│   ├── DELIVERY.md          # 交付文档
│   ├── CHANGELOG.md         # 变更记录
│   └── CONFLICTS.md         # 冲突记录
├── user-agent/
│   ├── DELIVERY.md
│   └── ...
```

## 2. 交付文档模板

### DELIVERY.md 必须包含:
1. **交付信息**: Agent名称、任务ID、状态、日期
2. **交付物清单**: 文件列表、说明、依赖关系
3. **依赖关系**: 依赖谁、被谁依赖
4. **使用指南**: 其他Agent如何使用
5. **冲突处理**: 可能出现的冲突及预案

### CHANGELOG.md 记录:
- 每次修改的内容
- 修改原因
- 影响范围

### CONFLICTS.md 记录:
- 与其他Agent的冲突
- 解决方式
- 最终决策

## 3. 冲突处理流程

```
Agent A 发现冲突
    ↓
1. 暂停工作，记录冲突到 CONFLICTS.md
    ↓
2. 通知相关Agent (在CONFLICTS.md中@)
    ↓
3. 尝试协商解决
    ↓
协商成功? 
    ├─ 是 → 记录解决方案，继续工作
    └─ 否 → 上报 Architect Agent 裁决
                ↓
         Architect 决策
                ↓
         记录最终方案，强制执行
```

## 4. 文件所有权声明

在代码文件头部添加所有权注释:

**Go文件**:
```go
// File: handlers/user.go
// Owner: UserAgent
// Task: P1-T5
// LastModified: 2026-02-22
// Description: User authentication and CRUD handlers
```

**TypeScript文件**:
```typescript
/**
 * @file services/user.ts
 * @owner UserAgent
 * @task P1-T5
 * @lastModified 2026-02-22
 * @description User API service
 */
```

## 5. 共享资源修改规则

| 资源类型 | 修改规则 |
|----------|----------|
| `types/index.ts` | 可追加，不要修改已有类型 |
| `App.tsx` | 追加路由，不要修改现有路由 |
| `go.mod` | 追加依赖，版本需一致 |
| 数据库迁移 | 按序号递增，不要修改已有迁移 |

## 6. 每日同步机制

每个Agent每天结束时:
1. 更新 DELIVERY.md 状态
2. 提交代码到git
3. 检查其他Agent的交付文档
4. 更新依赖关系图

## 7. 紧急冲突处理

如果出现以下情况，立即上报人类监督者:
- 架构设计冲突
- 数据库表结构冲突
- API契约冲突
- 无法达成一致的规范冲突

---

*规范版本: 1.0*
