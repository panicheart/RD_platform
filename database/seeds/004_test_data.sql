-- =====================================================
-- RDP Test Seed Data
-- Version: 1.1
-- Date: 2026-02-23
-- Description: Test data for development and testing
-- =====================================================

-- =====================================================
-- 1. Test Organization
-- =====================================================

INSERT INTO organizations (id, name, code, level, sort_order, description, is_active)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440000', '微波室', 'RD_DEPT', 1, 1, '微波室研发管理部门', true)
ON CONFLICT (code) DO NOTHING;

-- =====================================================
-- 2. Test Users
-- Password for all: test123
-- =====================================================

-- Admin user
INSERT INTO users (
    id, username, display_name, email, phone, 
    role, team, title, organization_id,
    password_hash, is_active, created_at, updated_at
) VALUES (
    '550e8400-e29b-41d4-a716-446655440001',
    'admin',
    '系统管理员',
    'admin@rdp.local',
    '13800000001',
    'admin',
    'general_mgmt',
    'senior_eng',
    '550e8400-e29b-41d4-a716-446655440000',
    '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEv.F3zUHJe3vJhL2a', -- bcrypt hash for 'test123'
    true,
    NOW(),
    NOW()
) ON CONFLICT (username) DO UPDATE SET
    password_hash = '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEv.F3zUHJe3vJhL2a',
    is_active = true,
    updated_at = NOW();

-- Department leader
INSERT INTO users (
    id, username, display_name, email, phone,
    role, team, title, organization_id,
    password_hash, is_active, created_at
) VALUES (
    '550e8400-e29b-41d4-a716-446655440002',
    'dept_leader',
    '部门领导',
    'leader@rdp.local',
    '13800000002',
    'dept_leader',
    'product_mgmt',
    'researcher',
    '550e8400-e29b-41d4-a716-446655440000',
    '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEv.F3zUHJe3vJhL2a',
    true,
    NOW()
) ON CONFLICT (username) DO NOTHING;

-- Team leader
INSERT INTO users (
    id, username, display_name, email, phone,
    role, team, title, organization_id,
    password_hash, is_active, created_at
) VALUES (
    '550e8400-e29b-41d4-a716-446655440003',
    'team_leader',
    '团队组长',
    'teamlead@rdp.local',
    '13800000003',
    'team_leader',
    'product_dev',
    'senior_eng',
    '550e8400-e29b-41d4-a716-446655440000',
    '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEv.F3zUHJe3vJhL2a',
    true,
    NOW()
) ON CONFLICT (username) DO NOTHING;

-- Designer user
INSERT INTO users (
    id, username, display_name, email, phone,
    role, team, title, organization_id,
    password_hash, is_active, created_at
) VALUES (
    '550e8400-e29b-41d4-a716-446655440004',
    'designer',
    '设计师',
    'designer@rdp.local',
    '13800000004',
    'designer',
    'tech_dev',
    'engineer',
    '550e8400-e29b-41d4-a716-446655440000',
    '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEv.F3zUHJe3vJhL2a',
    true,
    NOW()
) ON CONFLICT (username) DO NOTHING;

-- =====================================================
-- 3. Test Projects
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
        '550e8400-e29b-41d4-a716-446655440002',
        'internal',
        '550e8400-e29b-41d4-a716-446655440001',
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
        '550e8400-e29b-41d4-a716-446655440003',
        'internal',
        '550e8400-e29b-41d4-a716-446655440001',
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
        '550e8400-e29b-41d4-a716-446655440002',
        'internal',
        '550e8400-e29b-41d4-a716-446655440001',
        NOW(),
        NOW()
    )
ON CONFLICT (code) DO NOTHING;

-- =====================================================
-- 4. Test Forum Boards
-- =====================================================

INSERT INTO forum_boards (id, name, description, category, icon, sort_order, is_active, created_at)
VALUES 
    ('550e8400-e29b-41d4-a716-446655442001', '技术讨论', '技术相关问题讨论和交流', 'tech', 'code', 1, true, NOW()),
    ('550e8400-e29b-41d4-a716-446655442002', '求助问答', '遇到问题可以在这里提问', 'help', 'question-circle', 2, true, NOW()),
    ('550e8400-e29b-41d4-a716-446655442003', '综合讨论', '综合话题讨论区', 'general', 'message', 3, true, NOW()),
    ('550e8400-e29b-41d4-a716-446655442004', '公告通知', '官方公告和重要通知', 'announcement', 'notification', 0, true, NOW())
ON CONFLICT DO NOTHING;

-- =====================================================
-- 5. Test Forum Posts
-- =====================================================

INSERT INTO forum_posts (
    id, board_id, title, content, author_id, author_name,
    view_count, reply_count, is_pinned, is_locked, tags, created_at, updated_at
) VALUES 
    (
        '550e8400-e29b-41d4-a716-446655443001',
        '550e8400-e29b-41d4-a716-446655442001',
        '欢迎使用技术论坛',
        '这是技术论坛的第一篇帖子，欢迎大家积极交流技术问题！\n\n支持Markdown格式：\n- **粗体文字**\n- *斜体文字*\n- 代码块\n\n```go\nfmt.Println(\"Hello, RDP!\")\n```',
        '550e8400-e29b-41d4-a716-446655440001',
        '系统管理员',
        128,
        5,
        true,
        false,
        '[\"公告\",\"欢迎\"]',
        NOW() - INTERVAL '7 days',
        NOW()
    ),
    (
        '550e8400-e29b-41d4-a716-446655443002',
        '550e8400-e29b-41d4-a716-446655442001',
        '如何优化Go代码性能？',
        '最近在做性能优化，想请教一下大家有什么好的建议？\n\n目前遇到的问题是：\n1. 内存占用过高\n2. GC频率太高\n\n谢谢！',
        '550e8400-e29b-41d4-a716-446655440004',
        '设计师',
        56,
        3,
        false,
        false,
        '[\"go\",\"性能优化\"]',
        NOW() - INTERVAL '3 days',
        NOW()
    ),
    (
        '550e8400-e29b-41d4-a716-446655443003',
        '550e8400-e29b-41d4-a716-446655442002',
        '新人求助：怎么创建新项目？',
        '刚入职不久，想了解一下怎么在系统里创建新项目？\n\n有没有详细的操作指南？',
        '550e8400-e29b-41d4-a716-446655440004',
        '设计师',
        32,
        2,
        false,
        false,
        '[\"新人\",\"求助\"]',
        NOW() - INTERVAL '1 day',
        NOW()
    )
ON CONFLICT DO NOTHING;

-- =====================================================
-- 6. Test Forum Replies
-- =====================================================

INSERT INTO forum_replies (id, post_id, content, author_id, author_name, created_at, updated_at)
VALUES 
    (
        '550e8400-e29b-41d4-a716-446655444001',
        '550e8400-e29b-41d4-a716-446655443001',
        '感谢分享！期待更多技术讨论。',
        '550e8400-e29b-41d4-a716-446655440002',
        '部门领导',
        NOW() - INTERVAL '6 days',
        NOW()
    ),
    (
        '550e8400-e29b-41d4-a716-446655444002',
        '550e8400-e29b-41d4-a716-446655443001',
        '建议增加一些实际案例分享。',
        '550e8400-e29b-41d4-a716-446655440003',
        '团队组长',
        NOW() - INTERVAL '5 days',
        NOW()
    ),
    (
        '550e8400-e29b-41d4-a716-446655444003',
        '550e8400-e29b-41d4-a716-446655443002',
        '可以参考官方的性能优化指南。',
        '550e8400-e29b-41d4-a716-446655440003',
        '团队组长',
        NOW() - INTERVAL '2 days',
        NOW()
    )
ON CONFLICT DO NOTHING;

-- =====================================================
-- 7. Test Knowledge Categories
-- =====================================================

INSERT INTO knowledge_categories (id, name, description, parent_id, sort_order, is_active, created_at)
VALUES 
    ('550e8400-e29b-41d4-a716-446655445001', '技术文档', '技术相关文档资料', NULL, 1, true, NOW()),
    ('550e8400-e29b-41d4-a716-446655445002', '设计规范', '产品设计规范文档', NULL, 2, true, NOW()),
    ('550e8400-e29b-41d4-a716-446655445003', '流程制度', '部门流程和管理制度', NULL, 3, true, NOW()),
    ('550e8400-e29b-41d4-a716-446655445004', 'Go语言', 'Go语言编程指南', '550e8400-e29b-41d4-a716-446655445001', 1, true, NOW()),
    ('550e8400-e29b-41d4-a716-446655445005', '前端开发', '前端开发技术文档', '550e8400-e29b-41d4-a716-446655445001', 2, true, NOW())
ON CONFLICT DO NOTHING;

-- =====================================================
-- 8. Test Knowledge Items
-- =====================================================

INSERT INTO knowledge (
    id, title, content, category_id, author_id, author_name,
    status, view_count, version, source, created_at, updated_at
) VALUES 
    (
        '550e8400-e29b-41d4-a716-446655446001',
        'Go语言编码规范',
        '# Go语言编码规范\n\n## 1. 代码格式\n\n- 使用gofmt格式化代码\n- 每行不超过120个字符\n- 使用4个空格缩进\n\n## 2. 命名规范\n\n- 包名：小写，简洁\n- 函数名：驼峰命名\n- 常量：全大写，下划线分隔',
        '550e8400-e29b-41d4-a716-446655445004',
        '550e8400-e29b-41d4-a716-446655440003',
        '团队组长',
        'published',
        256,
        1,
        'manual',
        NOW() - INTERVAL '30 days',
        NOW()
    ),
    (
        '550e8400-e29b-41d4-a716-446655446002',
        'React开发最佳实践',
        '# React开发最佳实践\n\n## 组件设计\n\n1. 单一职责原则\n2. Props向下传递\n3. 状态提升\n\n## 性能优化\n\n- 使用React.memo\n- 避免不必要的重渲染',
        '550e8400-e29b-41d4-a716-446655445005',
        '550e8400-e29b-41d4-a716-446655440004',
        '设计师',
        'published',
        128,
        2,
        'manual',
        NOW() - INTERVAL '15 days',
        NOW()
    ),
    (
        '550e8400-e29b-41d4-a716-446655446003',
        '项目开发流程',
        '# 项目开发流程\n\n## 阶段一：需求分析\n\n1. 收集需求\n2. 可行性分析\n3. 编写需求文档\n\n## 阶段二：方案设计\n\n1. 概要设计\n2. 详细设计\n3. 设计评审',
        '550e8400-e29b-41d4-a716-446655445003',
        '550e8400-e29b-41d4-a716-446655440002',
        '部门领导',
        'published',
        512,
        1,
        'manual',
        NOW() - INTERVAL '60 days',
        NOW()
    )
ON CONFLICT DO NOTHING;

-- =====================================================
-- 9. Test Notifications
-- =====================================================

INSERT INTO notifications (id, user_id, type, title, content, is_read, related_id, related_type, created_at)
VALUES 
    (
        '550e8400-e29b-41d4-a716-446655447001',
        '550e8400-e29b-41d4-a716-446655440004',
        'system',
        '欢迎使用RDP系统',
        '欢迎使用微波室研发管理平台！如有问题请联系管理员。',
        false,
        NULL,
        NULL,
        NOW() - INTERVAL '1 day'
    ),
    (
        '550e8400-e29b-41d4-a716-446655447002',
        '550e8400-e29b-41d4-a716-446655440004',
        'forum',
        '你的帖子收到了回复',
        '你发表的帖子"如何优化Go代码性能？"收到了新回复。',
        false,
        '550e8400-e29b-41d4-a716-446655443002',
        'forum_post',
        NOW() - INTERVAL '2 hours'
    ),
    (
        '550e8400-e29b-41d4-a716-446655447003',
        '550e8400-e29b-41d4-a716-446655440004',
        'project',
        '你被分配到新项目',
        '你被分配参与项目"测试项目一"。',
        true,
        '550e8400-e29b-41d4-a716-446655441001',
        'project',
        NOW() - INTERVAL '3 days'
    )
ON CONFLICT DO NOTHING;

-- =====================================================
-- 10. Verify data insertion
-- =====================================================

SELECT 'Users created' as info, COUNT(*) as count FROM users
UNION ALL
SELECT 'Projects created', COUNT(*) FROM projects
UNION ALL
SELECT 'Forum boards created', COUNT(*) FROM forum_boards
UNION ALL
SELECT 'Forum posts created', COUNT(*) FROM forum_posts
UNION ALL
SELECT 'Forum replies created', COUNT(*) FROM forum_replies
UNION ALL
SELECT 'Knowledge categories created', COUNT(*) FROM knowledge_categories
UNION ALL
SELECT 'Knowledge items created', COUNT(*) FROM knowledge
UNION ALL
SELECT 'Notifications created', COUNT(*) FROM notifications;