# Testing Guide

## 测试框架

- **单元测试**: Vitest + React Testing Library
- **覆盖率**: V8 覆盖率报告
- **最低覆盖率**: 60% (lines/functions/statements), 50% (branches)

## 快速开始

```bash
# 运行所有测试
npm run test

# 运行测试并查看覆盖率
npm run test:coverage

# 运行测试UI模式
npm run test:ui

# 运行特定测试文件
npm run test -- src/utils/example.test.ts
```

## 测试规范

### 文件命名
- 测试文件: `*.test.ts` 或 `*.spec.ts`
- 与源文件同目录或 `__tests__/` 子目录

### 示例测试

```typescript
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { YourComponent } from './YourComponent';

describe('YourComponent', () => {
  it('should render correctly', () => {
    render(<YourComponent />);
    expect(screen.getByText('Hello')).toBeInTheDocument();
  });
});
```

### API测试示例

```typescript
import { describe, it, expect, vi } from 'vitest';
import { userAPI } from '@services/api';

describe('userAPI', () => {
  it('should fetch users', async () => {
    const mockUsers = [{ id: '1', name: 'Test' }];
    global.fetch = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ code: 200, data: mockUsers }),
    });
    
    const result = await userAPI.getUsers();
    expect(result.data).toEqual(mockUsers);
  });
});
```

## 测试工具

### 可用API
- `render()` - 渲染组件
- `screen` - 查询DOM
- `fireEvent` / `userEvent` - 模拟用户交互
- `vi.fn()` - Mock函数
- `vi.mock()` - Mock模块

### 自定义Matcher
已配置 `@testing-library/jest-dom`，可用:
- `toBeInTheDocument()`
- `toHaveClass()`
- `toBeVisible()`
- `toBeDisabled()`
- ...

## 覆盖率报告

运行 `npm run test:coverage` 后查看:
- 终端输出
- `coverage/` 目录下的 HTML 报告
