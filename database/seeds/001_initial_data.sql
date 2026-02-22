-- =====================================================
-- RDP Initial Seed Data
-- Version: 1.0
-- Date: 2026-02-21
-- Description: Initial data for RDP platform
-- =====================================================

-- =====================================================
-- 1. Default Organization Structure
-- =====================================================

-- Insert root department
INSERT INTO organizations (id, name, code, level, sort_order, description, is_active)
VALUES 
    ('01ORG000000000000000000001', '微波室', 'RD_DEPT', 1, 1, '微波室研发管理部门', true)
ON CONFLICT (code) DO NOTHING;

-- Insert sub-departments
INSERT INTO organizations (id, name, code, parent_id, level, sort_order, description, is_active)
VALUES 
    ('01ORG000000000000000000002', '产品管理组', 'PRODUCT_MGMT', '01ORG000000000000000000001', 2, 1, '负责产品规划与管理', true),
    ('01ORG000000000000000000003', '产品开发组', 'PRODUCT_DEV', '01ORG000000000000000001', 2, 2, '负责产品设计与开发', true),
    ('01ORG000000000000000000004', '技术开发组', 'TECH_DEV', '01ORG000000000000000001', 2, 3, '负责技术预研与平台开发', true),
    ('01ORG000000000000000000005', '综合管理组', 'GENERAL_MGMT', '01ORG000000000000000001', 2, 4, '负责部门综合管理', true)
ON CONFLICT (code) DO NOTHING;

-- =====================================================
-- 2. Default Admin User
-- Password: admin123 (bcrypt hashed)
-- MUST CHANGE ON FIRST LOGIN
-- =====================================================

INSERT INTO users (
    id, username, display_name, email, phone, 
    role, team, title, organization_id,
    password_hash, is_active, created_at, updated_at
) VALUES (
    '01USER00000000000000000001',
    'admin',
    '系统管理员',
    'admin@rdp.local',
    '13800000000',
    'admin',
    'general_mgmt',
    'senior_eng',
    '01ORG000000000000000000001',
    crypt('admin123', gen_salt('bf')),
    true,
    NOW(),
    NOW()
) ON CONFLICT (username) DO UPDATE SET
    password_hash = crypt('admin123', gen_salt('bf')),
    is_active = true,
    updated_at = NOW();

-- =====================================================
-- 3. Sample Users for Testing
-- =====================================================

INSERT INTO users (
    id, username, display_name, email, 
    role, team, title, organization_id,
    password_hash, is_active, created_at
) VALUES 
    ('01USER00000000000000000002', 'zhangsan', '张三', 'zhangsan@rdp.local', 
     'dept_leader', 'product_mgmt', 'researcher', '01ORG000000000000000000002',
     crypt('test123', gen_salt('bf')), true, NOW()),
    
    ('01USER00000000000000000003', 'lisi', '李四', 'lisi@rdp.local',
     'team_leader', 'product_dev', 'senior_eng', '01ORG000000000000000000003',
     crypt('test123', gen_salt('bf')), true, NOW()),
    
    ('01USER00000000000000000004', 'wangwu', '王五', 'wangwu@rdp.local',
     'designer', 'tech_dev', 'engineer', '01ORG000000000000000000004',
     crypt('test123', gen_salt('bf')), true, NOW())
ON CONFLICT (username) DO NOTHING;

-- =====================================================
-- 4. Process Templates (7 Project Categories)
-- =====================================================

INSERT INTO process_templates (
    id, name, code, category, description, 
    activities, is_default, is_active, created_at, updated_at
) VALUES
    -- 单机产品开发流程
    (
        '01PROC00000000000000000001',
        '单机产品开发流程',
        'PROCESS_STANDALONE',
        'standalone',
        '单机产品完整开发流程，包含需求分析到产品定型',
        '[
            {"id": "ACT001", "name": "需求分析", "duration": 10, "require_review": true},
            {"id": "ACT002", "name": "方案设计", "duration": 15, "require_review": true},
            {"id": "ACT003", "name": "详细设计", "duration": 20, "require_review": false},
            {"id": "ACT004", "name": "硬件实现", "duration": 30, "require_review": false},
            {"id": "ACT005", "name": "软件实现", "duration": 30, "require_review": false},
            {"id": "ACT006", "name": "系统集成", "duration": 15, "require_review": true},
            {"id": "ACT007", "name": "测试验证", "duration": 20, "require_review": true},
            {"id": "ACT008", "name": "产品定型", "duration": 10, "require_review": true}
        ]'::jsonb,
        true,
        true,
        NOW(),
        NOW()
    ),
    -- 模块开发流程
    (
        '01PROC00000000000000000002',
        '模块开发流程',
        'PROCESS_MODULE',
        'module',
        '通用模块开发流程，适用于可复用模块',
        '[
            {"id": "ACT001", "name": "需求分析", "duration": 7, "require_review": true},
            {"id": "ACT002", "name": "方案设计", "duration": 10, "require_review": true},
            {"id": "ACT003", "name": "详细设计", "duration": 15, "require_review": false},
            {"id": "ACT004", "name": "模块实现", "duration": 20, "require_review": false},
            {"id": "ACT005", "name": "测试验证", "duration": 10, "require_review": true},
            {"id": "ACT006", "name": "模块定型", "duration": 5, "require_review": true}
        ]'::jsonb,
        false,
        true,
        NOW(),
        NOW()
    ),
    -- 软件开发流程
    (
        '01PROC00000000000000000003',
        '软件开发流程',
        'PROCESS_SOFTWARE',
        'software',
        '软件项目开发流程，包含敏捷元素',
        '[
            {"id": "ACT001", "name": "需求分析", "duration": 7, "require_review": true},
            {"id": "ACT002", "name": "架构设计", "duration": 10, "require_review": true},
            {"id": "ACT003", "name": "详细设计", "duration": 10, "require_review": false},
            {"id": "ACT004", "name": "编码实现", "duration": 25, "require_review": false},
            {"id": "ACT005", "name": "单元测试", "duration": 7, "require_review": false},
            {"id": "ACT006", "name": "集成测试", "duration": 10, "require_review": true},
            {"id": "ACT007", "name": "系统测试", "duration": 10, "require_review": true},
            {"id": "ACT008", "name": "发布上线", "duration": 3, "require_review": true}
        ]'::jsonb,
        false,
        true,
        NOW(),
        NOW()
    ),
    -- 技术开发流程
    (
        '01PROC00000000000000000004',
        '技术开发流程',
        'PROCESS_TECH_DEV',
        'tech_dev',
        '技术预研开发流程，用于新技术研究',
        '[
            {"id": "ACT001", "name": "技术调研", "duration": 15, "require_review": true},
            {"id": "ACT002", "name": "原理验证", "duration": 20, "require_review": true},
            {"id": "ACT003", "name": "方案设计", "duration": 15, "require_review": false},
            {"id": "ACT004", "name": "实验验证", "duration": 30, "require_review": true},
            {"id": "ACT005", "name": "技术总结", "duration": 10, "require_review": true}
        ]'::jsonb,
        false,
        true,
        NOW(),
        NOW()
    ),
    -- 工艺开发流程
    (
        '01PROC00000000000000000005',
        '工艺开发流程',
        'PROCESS_PROCESS_DEV',
        'process_dev',
        '工艺改进与开发流程',
        '[
            {"id": "ACT001", "name": "问题定义", "duration": 5, "require_review": true},
            {"id": "ACT002", "name": "方案论证", "duration": 10, "require_review": true},
            {"id": "ACT003", "name": "试点实施", "duration": 20, "require_review": false},
            {"id": "ACT004", "name": "效果评估", "duration": 10, "require_review": true},
            {"id": "ACT005", "name": "推广固化", "duration": 5, "require_review": true}
        ]'::jsonb,
        false,
        true,
        NOW(),
        NOW()
    ),
    -- 知识开发流程
    (
        '01PROC00000000000000000006',
        '知识开发流程',
        'PROCESS_KNOWLEDGE_DEV',
        'knowledge_dev',
        '知识沉淀与文档开发流程',
        '[
            {"id": "ACT001", "name": "知识梳理", "duration": 5, "require_review": false},
            {"id": "ACT002", "name": "文档编写", "duration": 10, "require_review": false},
            {"id": "ACT003", "name": "评审修订", "duration": 5, "require_review": true},
            {"id": "ACT004", "name": "发布归档", "duration": 3, "require_review": true}
        ]'::jsonb,
        false,
        true,
        NOW(),
        NOW()
    ),
    -- 产品化开发流程
    (
        '01PROC00000000000000000007',
        '产品化开发流程',
        'PROCESS_PRODUCT_LAUNCH',
        'product_launch',
        '研究成果产品化开发流程',
        '[
            {"id": "ACT001", "name": "市场分析", "duration": 10, "require_review": true},
            {"id": "ACT002", "name": "产品定义", "duration": 10, "require_review": true},
            {"id": "ACT003", "name": "原型开发", "duration": 20, "require_review": false},
            {"id": "ACT004", "name": "用户验证", "duration": 15, "require_review": true},
            {"id": "ACT005", "name": "批量生产", "duration": 15, "require_review": true}
        ]'::jsonb,
        false,
        true,
        NOW(),
        NOW()
    )
ON CONFLICT (code) DO UPDATE SET
    activities = EXCLUDED.activities,
    updated_at = NOW();

-- =====================================================
-- 5. Sample Project (for demonstration)
-- =====================================================

INSERT INTO projects (
    id, code, name, description, category, status,
    leader_id, created_by, start_date, end_date,
    classification_level, progress, created_at, updated_at
) VALUES (
    '01PROJ00000000000000000001',
    'SAMPLE-2026-001',
    '示例项目 - 射频模块开发',
    '这是一个示例项目，用于演示系统功能',
    'module',
    'planning',
    '01USER00000000000000000002',
    '01USER00000000000000000001',
    CURRENT_DATE,
    CURRENT_DATE + INTERVAL '90 days',
    'internal',
    0,
    NOW(),
    NOW()
) ON CONFLICT (code) DO NOTHING;

-- Add project member
INSERT INTO project_members (project_id, user_id, role, joined_at)
VALUES 
    ('01PROJ00000000000000000001', '01USER00000000000000000002', 'leader', NOW()),
    ('01PROJ00000000000000000001', '01USER00000000000000000003', 'member', NOW()),
    ('01PROJ00000000000000000001', '01USER00000000000000000004', 'member', NOW())
ON CONFLICT DO NOTHING;

-- =====================================================
-- 6. System Announcements
-- =====================================================

INSERT INTO announcements (id, title, content, author_id, priority, is_pinned, published_at)
VALUES (
    uuid_generate_v4(),
    '🎉 RDP系统正式上线',
    '欢迎使用微波室研发管理平台！本系统支持项目管理、流程执行、文档协作等功能。如有问题请联系系统管理员。',
    '01USER00000000000000000001',
    'high',
    true,
    NOW()
) ON CONFLICT DO NOTHING;

-- =====================================================
-- 7. Sample Notifications for Admin
-- =====================================================

INSERT INTO notifications (id, user_id, type, title, content, is_read, created_at)
VALUES 
    (
        uuid_generate_v4(),
        '01USER00000000000000000001',
        'system',
        '欢迎使用 RDP 系统',
        '系统初始化完成，您可以开始使用了。',
        false,
        NOW()
    ),
    (
        uuid_generate_v4(),
        '01USER00000000000000000001',
        'project',
        '示例项目已创建',
        '系统已自动创建一个示例项目供您参考。',
        false,
        NOW() - INTERVAL '1 hour'
    )
ON CONFLICT DO NOTHING;

-- =====================================================
-- 8. Sample Honors
-- =====================================================

INSERT INTO honors (id, title, description, award_year, award_month, recipient_id, is_active, sort_order, created_at)
VALUES (
    uuid_generate_v4(),
    '年度优秀研发团队',
    '在2025年度研发工作中表现突出，获得优秀团队称号',
    2025,
    12,
    '01USER00000000000000000002',
    true,
    1,
    NOW()
) ON CONFLICT DO NOTHING;

-- =====================================================
-- Seed Data Complete
-- =====================================================

DO $$
BEGIN
    RAISE NOTICE '种子数据导入完成！';
    RAISE NOTICE '默认管理员: admin / admin123';
    RAISE NOTICE '测试用户: zhangsan/lisi/wangwu / test123';
END $$;
