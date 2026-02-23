-- =====================================================
-- RDP Test Seed Data - Version 3.0 (Short IDs)
-- Date: 2026-02-23
-- Description: Test data with short ULID format IDs
-- =====================================================

-- =====================================================
-- 1. Categories (for knowledge base)
-- =====================================================

INSERT INTO categories (id, name, description, parent_id, level, sort_order, created_at)
VALUES 
    ('01HPVJ8ZM00000000000000001', '技术文档', '技术相关文档资料', NULL, 1, 1, NOW()),
    ('01HPVJ8ZM00000000000000002', '设计规范', '产品设计规范文档', NULL, 1, 2, NOW()),
    ('01HPVJ8ZM00000000000000003', '流程制度', '部门流程和管理制度', NULL, 1, 3, NOW()),
    ('01HPVJ8ZM00000000000000004', 'Go语言', 'Go语言编程指南', '01HPVJ8ZM00000000000000001', 2, 1, NOW()),
    ('01HPVJ8ZM00000000000000005', '前端开发', '前端开发技术文档', '01HPVJ8ZM00000000000000001', 2, 2, NOW())
ON CONFLICT DO NOTHING;

-- =====================================================
-- 2. Knowledge Items
-- =====================================================

INSERT INTO knowledge (
    id, title, content, category_id, author_id,
    status, view_count, version, source, created_at, updated_at
) VALUES 
    (
        '01HPVJ8ZM00000000000000101',
        'Go语言编码规范',
        '# Go语言编码规范

## 1. 代码格式

- 使用gofmt格式化代码
- 每行不超过120个字符
- 使用4个空格缩进

## 2. 命名规范

- 包名：小写，简洁
- 函数名：驼峰命名
- 常量：全大写，下划线分隔',
        '01HPVJ8ZM00000000000000004',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        'published',
        256,
        1,
        'manual',
        NOW() - INTERVAL '30 days',
        NOW()
    ),
    (
        '01HPVJ8ZM00000000000000102',
        'React开发最佳实践',
        '# React开发最佳实践

## 组件设计

1. 单一职责原则
2. Props向下传递
3. 状态提升

## 性能优化

- 使用React.memo
- 避免不必要的重渲染',
        '01HPVJ8ZM00000000000000005',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        'published',
        128,
        2,
        'manual',
        NOW() - INTERVAL '15 days',
        NOW()
    ),
    (
        '01HPVJ8ZM00000000000000103',
        '项目开发流程',
        '# 项目开发流程

## 阶段一：需求分析

1. 收集需求
2. 可行性分析
3. 编写需求文档

## 阶段二：方案设计

1. 概要设计
2. 详细设计
3. 设计评审',
        '01HPVJ8ZM00000000000000003',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        'published',
        512,
        1,
        'manual',
        NOW() - INTERVAL '60 days',
        NOW()
    )
ON CONFLICT DO NOTHING;

-- =====================================================
-- 3. Forum Boards
-- =====================================================

INSERT INTO forum_boards (id, name, description, category, icon, sort_order, is_active, created_at)
VALUES 
    ('01HPVJ8ZM00000000000000201', '技术讨论', '技术相关问题讨论和交流', 'tech', 'code', 1, true, NOW()),
    ('01HPVJ8ZM00000000000000202', '求助问答', '遇到问题可以在这里提问', 'help', 'question-circle', 2, true, NOW()),
    ('01HPVJ8ZM00000000000000203', '综合讨论', '综合话题讨论区', 'general', 'message', 3, true, NOW()),
    ('01HPVJ8ZM00000000000000204', '公告通知', '官方公告和重要通知', 'announcement', 'notification', 0, true, NOW())
ON CONFLICT DO NOTHING;

-- =====================================================
-- 4. Forum Posts
-- =====================================================

INSERT INTO forum_posts (
    id, board_id, title, content, author_id, author_name,
    view_count, reply_count, is_pinned, is_locked, tags, created_at, updated_at
) VALUES 
    (
        '01HPVJ8ZM00000000000000301',
        '01HPVJ8ZM00000000000000204',
        '欢迎使用技术论坛',
        '这是技术论坛的第一篇帖子，欢迎大家积极交流技术问题！

支持Markdown格式：
- **粗体文字**
- *斜体文字*
- 代码块

```go
fmt.Println("Hello, RDP!")
```',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        '系统管理员',
        128,
        5,
        true,
        false,
        '["公告","欢迎"]',
        NOW() - INTERVAL '7 days',
        NOW()
    ),
    (
        '01HPVJ8ZM00000000000000302',
        '01HPVJ8ZM00000000000000201',
        '如何优化Go代码性能？',
        '最近在做性能优化，想请教一下大家有什么好的建议？

目前遇到的问题是：
1. 内存占用过高
2. GC频率太高

谢谢！',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        '管理员',
        56,
        3,
        false,
        false,
        '["go","性能优化"]',
        NOW() - INTERVAL '3 days',
        NOW()
    ),
    (
        '01HPVJ8ZM00000000000000303',
        '01HPVJ8ZM00000000000000202',
        '新人求助：怎么创建新项目？',
        '刚入职不久，想了解一下怎么在系统里创建新项目？

有没有详细的操作指南？',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        '管理员',
        32,
        2,
        false,
        false,
        '["新人","求助"]',
        NOW() - INTERVAL '1 day',
        NOW()
    )
ON CONFLICT DO NOTHING;

-- =====================================================
-- 5. Forum Replies
-- =====================================================

INSERT INTO forum_replies (id, post_id, content, author_id, author_name, created_at, updated_at)
VALUES 
    (
        '01HPVJ8ZM00000000000000401',
        '01HPVJ8ZM00000000000000301',
        '感谢分享！期待更多技术讨论。',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        '管理员',
        NOW() - INTERVAL '6 days',
        NOW()
    ),
    (
        '01HPVJ8ZM00000000000000402',
        '01HPVJ8ZM00000000000000301',
        '建议增加一些实际案例分享。',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        '管理员',
        NOW() - INTERVAL '5 days',
        NOW()
    ),
    (
        '01HPVJ8ZM00000000000000403',
        '01HPVJ8ZM00000000000000302',
        '可以参考官方的性能优化指南。',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        '管理员',
        NOW() - INTERVAL '2 days',
        NOW()
    )
ON CONFLICT DO NOTHING;

-- =====================================================
-- 6. Test Projects
-- =====================================================

INSERT INTO projects (
    id, code, name, description, category, status,
    start_date, end_date, progress, leader_id,
    classification_level, created_by, created_at, updated_at
) VALUES 
    (
        '550e8400-e29b-41d4-a716-446655441001',
        'PROJ-2026-001',
        '测试项目一',
        '用于功能测试的示例项目',
        'standalone',
        'in_progress',
        '2026-01-01',
        '2026-12-31',
        45,
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        'internal',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        NOW(),
        NOW()
    ),
    (
        '550e8400-e29b-41d4-a716-446655441002',
        'PROJ-2026-002',
        '测试项目二',
        '第二个测试项目',
        'module',
        'planning',
        '2026-03-01',
        '2026-09-30',
        10,
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        'internal',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        NOW(),
        NOW()
    ),
    (
        '550e8400-e29b-41d4-a716-446655441003',
        'PROJ-2026-003',
        '已完成项目',
        '已完成的示例项目',
        'software',
        'completed',
        '2025-06-01',
        '2025-12-31',
        100,
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        'internal',
        'e7aafda0-2de9-44b8-8e18-1f8a86025449',
        NOW(),
        NOW()
    )
ON CONFLICT (code) DO NOTHING;

-- =====================================================
-- 7. Verify data insertion
-- =====================================================

SELECT 'Categories created' as info, COUNT(*) as count FROM categories
UNION ALL
SELECT 'Knowledge items created', COUNT(*) FROM knowledge
UNION ALL
SELECT 'Forum boards created', COUNT(*) FROM forum_boards
UNION ALL
SELECT 'Forum posts created', COUNT(*) FROM forum_posts
UNION ALL
SELECT 'Forum replies created', COUNT(*) FROM forum_replies
UNION ALL
SELECT 'Projects created', COUNT(*) FROM projects;
