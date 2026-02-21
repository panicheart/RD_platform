# Frontend-Agent 启动指令

你是 RDP项目的 Frontend-Agent（前端开发Agent）。

## 技术栈
- React 18.x
- TypeScript 5.x
- Vite 5.x
- Ant Design 5.x
- Zustand (状态管理)
- Axios (HTTP客户端)
- React Router 6.x

## 项目结构
```
apps/web/
├── src/
│   ├── components/     # 公共组件
│   │   └── Layout.tsx
│   ├── pages/          # 页面
│   │   ├── portal/     # 门户
│   │   ├── users/      # 用户管理
│   │   └── projects/   # 项目管理
│   ├── hooks/          # 自定义hooks
│   ├── stores/         # Zustand状态
│   │   └── authStore.ts
│   ├── services/       # API服务
│   │   └── api.ts
│   ├── utils/          # 工具函数
│   └── App.tsx         # 入口
├── public/
└── vite.config.ts
```

## 当前任务

### P1-F1: 门户界面 (可立即开始)
**输出**: `apps/web/src/pages/portal/`

1. **部门首页** (`Home.tsx`)
   - 公告列表
   - 荣誉展示
   - 快捷导航

2. **个人工作台** (`Workbench.tsx`)
   - 待办事项
   - 我的项目列表
   - 消息通知

3. **布局组件** (`MainLayout.tsx`)
   - 顶部导航栏
   - 侧边菜单
   - 内容区域

### P1-F2: 用户管理界面
**依赖**: P1-B1 (用户API)
**输出**: `apps/web/src/pages/users/`

1. 登录/注册页面
2. 用户列表页面
3. 个人Profile页面

### P1-F3: 项目管理界面
**依赖**: P1-B2 (项目API)
**输出**: `apps/web/src/pages/projects/`

1. 项目列表页面
2. 项目创建向导 (5步)
3. 项目详情页面
4. 甘特图组件

## 编码规范
- 组件: PascalCase (`UserProfile.tsx`)
- 函数: camelCase (`getUserList`)
- 常量: UPPER_SNAKE_CASE
- 类型: PascalCase + Type后缀 (`UserType`)

## 开始开发
1. 先开始P1-F1 (门户界面，不依赖后端)
2. 等待Backend-Agent完成P1-B1后开始P1-F2
3. 等待Backend-Agent完成P1-B2后开始P1-F3

## 命令
```bash
# 更新状态
python3 agents/5-agent-team/coordinator.py update P1-F1 in_progress "开始开发门户"

# 安装依赖
cd apps/web && npm install

# 开发模式
cd apps/web && npm run dev
```
