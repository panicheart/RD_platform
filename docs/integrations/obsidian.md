# Obsidian 集成指南

## 概述

Obsidian 集成模块提供双向同步功能，允许用户在 RDP 平台和 Obsidian 个人知识库之间无缝同步 Markdown 文档。

## 功能特性

- **WebDAV 服务**: 提供标准 WebDAV 协议支持，Obsidian 可通过 WebDAV 插件连接
- **双向同步**: 支持平台 ↔ Obsidian Vault 的双向文档同步
- **YAML Frontmatter**: 自动解析和生成 YAML 元数据（标题、标签、作者等）
- **Wiki 链接转换**: 支持 Obsidian 的 `[[内部链接]]` 格式
- **标签同步**: 自动同步 #标签 和 frontmatter 中的标签
- **冲突检测**: 智能检测文件冲突并提供解决策略

## 快速开始

### 1. 配置 Vault

```bash
# 通过 API 创建 Vault 映射
POST /api/v1/obsidian/vaults
{
  "vault_path": "/obsidian/my-vault",
  "local_path": "/home/user/Obsidian Vault",
  "category_id": "01H...",
  "auto_sync": true
}
```

### 2. 在 Obsidian 中配置 WebDAV

1. 安装 **Remotely Save** 插件或 **WebDAV Sync** 插件
2. 配置 WebDAV 地址: `https://your-rdp-server/api/v1/obsidian/vaults/{vault_id}/`
3. 输入认证 Token
4. 启用自动同步

### 3. 执行同步

```bash
# 触发手动同步
POST /api/v1/obsidian/vaults/{vault_id}/sync
```

## API 接口

### Vault 管理

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/obsidian/vaults` | 列出所有 Vault |
| POST | `/api/v1/obsidian/vaults` | 创建 Vault 映射 |
| GET | `/api/v1/obsidian/vaults/:id` | 获取 Vault 详情 |
| PUT | `/api/v1/obsidian/vaults/:id` | 更新 Vault 配置 |
| DELETE | `/api/v1/obsidian/vaults/:id` | 删除 Vault 映射 |

### 同步操作

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/obsidian/vaults/:id/sync` | 执行双向同步 |
| POST | `/api/v1/obsidian/vaults/:id/import` | 导入 Markdown 文件 |
| POST | `/api/v1/obsidian/vaults/:id/export/:knowledgeId` | 导出到 Vault |

### WebDAV 协议

| 方法 | 描述 |
|------|------|
| OPTIONS | 获取支持的方法 |
| PROPFIND | 列出目录内容 |
| GET | 获取文件内容 |
| PUT | 上传文件 |
| DELETE | 删除文件 |
| MKCOL | 创建目录 |
| MOVE | 移动/重命名文件 |

## 同步机制

### 导入流程 (Vault → Platform)

1. 扫描 Vault 目录中的 `.md` 文件
2. 解析 YAML frontmatter:
   ```yaml
   ---
   title: 文档标题
   author: 作者ID
   tags: 标签1, 标签2
   created: 2026-01-01T00:00:00Z
   updated: 2026-01-02T00:00:00Z
   ---
   ```
3. 提取正文 Markdown 内容
4. 提取 `#标签` 格式的标签
5. 创建或更新知识库条目

### 导出流程 (Platform → Vault)

1. 获取知识库条目
2. 生成 YAML frontmatter:
   ```yaml
   ---
   title: 知识标题
   id: 01HABC...
   author: 用户ID
   tags: 标签1, 标签2
   created: 2026-01-01T00:00:00Z
   updated: 2026-01-02T00:00:00Z
   ---
   ```
3. 写入 Markdown 内容
4. 保存到 Vault 目录

### 冲突解决

当同一文件在两端都有更新时，支持以下冲突解决策略:

| 策略 | 说明 |
|------|------|
| `platform_wins` | 平台版本优先 |
| `vault_wins` | Vault 版本优先 |
| `newer_wins` | 较新的版本优先 (默认) |

## 文件结构映射

```
Obsidian Vault/
├── 01-理论基础/
│   ├── 微波原理.md
│   └── 天线设计.md
├── 02-项目文档/
│   └── X波段组件设计.md
└── 附件/
    └── diagram.png
```

映射到平台分类:
- `01-理论基础` → 对应知识分类
- `02-项目文档` → 对应知识分类

## 安全考虑

- **路径安全检查**: 防止目录遍历攻击，所有操作必须在 Vault 目录内
- **文件类型限制**: 仅同步 `.md` 文件
- **权限控制**: 通过认证中间件验证用户权限

## 故障排查

### 同步失败

1. 检查 Vault 路径是否正确
2. 确认文件系统权限
3. 查看日志中的错误信息

### 文件未同步

1. 确认文件是 `.md` 格式
2. 检查是否在排除列表中
3. 验证文件编码为 UTF-8

### 中文乱码

确保所有 Markdown 文件使用 UTF-8 编码保存。

## 最佳实践

1. **定期同步**: 建议开启 `auto_sync` 自动同步
2. **分类管理**: 在 Vault 中使用文件夹组织文档
3. **标签规范**: 使用统一的标签命名规范
4. **版本控制**: 重要文档在导出前创建版本

## 技术细节

### 数据模型

```go
type ObsidianMapping struct {
    ID          string    // ULID
    VaultPath   string    // WebDAV 路径
    LocalPath   string    // 本地文件系统路径
    CategoryID  string    // 关联的知识分类
    AutoSync    bool      // 是否自动同步
    LastSyncAt  *time.Time
}
```

### WebDAV 实现

基于 RFC 4918 标准实现，支持:
- 属性查询 (PROPFIND)
- 文件读写 (GET/PUT)
- 目录操作 (MKCOL)
- 文件移动 (MOVE)

## 参考资料

- [Obsidian Help - Sync](https://help.obsidian.md/Obsidian+Sync)
- [WebDAV RFC 4918](https://tools.ietf.org/html/rfc4918)
- [Remotely Save Plugin](https://github.com/remotely-save/remotely-save)

---

*RDP Platform - Obsidian Integration v1.0*
