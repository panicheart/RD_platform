-- =====================================================
-- RDP Database Schema Initialization
-- Version: 1.0
-- Date: 2026-02-22
-- Description: Core database schema for Phase 1
-- =====================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =====================================================
-- ENUM Types
-- =====================================================

-- User roles
CREATE TYPE user_role AS ENUM ('admin', 'dept_leader', 'team_leader', 'designer', 'other');

-- Teams
CREATE TYPE team_type AS ENUM ('product_mgmt', 'product_dev', 'tech_dev', 'general_mgmt');

-- Product management specialty
CREATE TYPE pm_specialty AS ENUM ('model_mgmt', 'reliability');

-- Product development product line
CREATE TYPE pd_product_line AS ENUM ('line_a', 'line_b', 'line_c');

-- Technical development specialty
CREATE TYPE td_specialty AS ENUM ('antenna', 'rf', 'digital', 'power');

-- Designer title level
CREATE TYPE title_level AS ENUM ('designer_junior', 'assistant_eng', 'engineer', 'senior_eng', 'researcher');

-- Project category
CREATE TYPE project_category AS ENUM (
    'standalone',
    'module',
    'software',
    'tech_dev',
    'process_dev',
    'knowledge_dev',
    'product_launch'
);

-- Project status
CREATE TYPE project_status AS ENUM ('draft', 'planning', 'in_progress', 'review', 'completed', 'archived');

-- Activity status
CREATE TYPE activity_status AS ENUM ('pending', 'in_progress', 'review', 'completed', 'blocked');

-- Product maturity
CREATE TYPE product_maturity AS ENUM ('developing', 'prototype', 'engineering', 'qualified', 'mass_production');

-- Shelf type
CREATE TYPE shelf_type AS ENUM ('standalone', 'module', 'software', 'basic');

-- Basic subtype
CREATE TYPE basic_subtype AS ENUM ('component', 'fastener', 'material');

-- Knowledge category
CREATE TYPE knowledge_category AS ENUM (
    'theory',
    'standard',
    'regulation',
    'process',
    'case_positive',
    'case_negative',
    'simulation',
    'software_guide'
);

-- TRL level
CREATE TYPE trl_level AS ENUM ('TRL1', 'TRL2', 'TRL3', 'TRL4', 'TRL5', 'TRL6', 'TRL7', 'TRL8', 'TRL9');

-- Classification level
CREATE TYPE classification_level AS ENUM ('public', 'internal', 'secret', 'confidential');

-- =====================================================
-- Organizations Table
-- =====================================================

CREATE TABLE organizations (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(200) NOT NULL,
    code            VARCHAR(50) UNIQUE NOT NULL,
    parent_id       UUID REFERENCES organizations(id),
    level           INTEGER NOT NULL DEFAULT 1,
    description     TEXT,
    leader_id       UUID,
    sort_order      INTEGER DEFAULT 0,
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_organizations_parent ON organizations(parent_id);
CREATE INDEX idx_organizations_code ON organizations(code);
CREATE INDEX idx_organizations_level ON organizations(level);

-- =====================================================
-- Users Table
-- =====================================================

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username        VARCHAR(50) UNIQUE NOT NULL,
    display_name    VARCHAR(100) NOT NULL,
    email           VARCHAR(100),
    phone           VARCHAR(20),
    avatar_url      VARCHAR(500),
    role            user_role NOT NULL DEFAULT 'designer',
    team            team_type,
    specialty       VARCHAR(50),
    product_line    pd_product_line,
    title           title_level,
    organization_id UUID REFERENCES organizations(id),
    skills          JSONB DEFAULT '[]',
    honors          JSONB DEFAULT '[]',
    bio             TEXT,
    password_hash   VARCHAR(255),
    is_active       BOOLEAN DEFAULT true,
    casdoor_id      VARCHAR(100),
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_team ON users(team);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_organization ON users(organization_id);

-- =====================================================
-- Projects Table
-- =====================================================

CREATE TABLE projects (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code                VARCHAR(50) UNIQUE NOT NULL,
    name                VARCHAR(200) NOT NULL,
    description         TEXT,
    category            project_category NOT NULL,
    status              project_status NOT NULL DEFAULT 'draft',
    product_line        pd_product_line,
    team                team_type,
    
    -- Process binding
    process_template_id UUID,
    
    -- Dates
    start_date          DATE,
    end_date            DATE,
    actual_start_date   DATE,
    actual_end_date     DATE,
    
    -- Progress (0-100)
    progress            INTEGER DEFAULT 0,
    
    -- Git repository
    git_repo_id         VARCHAR(100),
    git_repo_url        VARCHAR(500),
    
    -- Classification
    classification_level classification_level DEFAULT 'internal',
    
    -- Team
    leader_id           UUID REFERENCES users(id),
    tech_leader_id     UUID REFERENCES users(id),
    product_leader_id  UUID REFERENCES users(id),
    
    -- Metadata
    created_by          UUID REFERENCES users(id),
    created_at         TIMESTAMPTZ DEFAULT NOW(),
    updated_at         TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_projects_code ON projects(code);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_category ON projects(category);
CREATE INDEX idx_projects_leader ON projects(leader_id);
CREATE INDEX idx_projects_created_at ON projects(created_at DESC);

-- =====================================================
-- Project Members Table
-- =====================================================

CREATE TABLE project_members (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(50) NOT NULL DEFAULT 'member',
    joined_at       TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(project_id, user_id)
);

CREATE INDEX idx_project_members_project ON project_members(project_id);
CREATE INDEX idx_project_members_user ON project_members(user_id);
CREATE INDEX idx_project_members_composite ON project_members(project_id, user_id);

-- =====================================================
-- Process Templates Table
-- =====================================================

CREATE TABLE process_templates (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(200) NOT NULL,
    code            VARCHAR(50) UNIQUE NOT NULL,
    category        project_category NOT NULL,
    description     TEXT,
    activities      JSONB NOT NULL DEFAULT '[]',
    is_default      BOOLEAN DEFAULT false,
    is_active       BOOLEAN DEFAULT true,
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_process_templates_category ON process_templates(category);
CREATE INDEX idx_process_templates_active ON process_templates(is_active);

-- =====================================================
-- Activities Table
-- =====================================================

CREATE TABLE activities (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    template_activity_id VARCHAR(50),
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    status          activity_status NOT NULL DEFAULT 'pending',
    sort_order      INTEGER DEFAULT 0,
    
    -- Dates
    start_date      DATE,
    end_date        DATE,
    actual_start_date DATE,
    actual_end_date DATE,
    
    -- Progress
    progress        INTEGER DEFAULT 0,
    
    -- Assignment
    assignee_id     UUID REFERENCES users(id),
    
    -- Dependencies
    depends_on      UUID[],
    
    -- Inputs/Outputs
    inputs          JSONB DEFAULT '[]',
    outputs         JSONB DEFAULT '[]',
    
    -- DCP Review
    require_review  BOOLEAN DEFAULT false,
    review_status   VARCHAR(20),
    
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_activities_project ON activities(project_id);
CREATE INDEX idx_activities_assignee ON activities(assignee_id);
CREATE INDEX idx_activities_status ON activities(status);
CREATE INDEX idx_activities_project_status ON activities(project_id, status);

-- =====================================================
-- Files Table
-- =====================================================

CREATE TABLE files (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id      UUID REFERENCES projects(id) ON DELETE CASCADE,
    parent_id       UUID REFERENCES files(id),
    name            VARCHAR(255) NOT NULL,
    path            VARCHAR(500) NOT NULL,
    size            BIGINT,
    mime_type       VARCHAR(100),
    is_directory    BOOLEAN DEFAULT false,
    
    -- Git info
    git_commit_hash  VARCHAR(40),
    git_commit_time TIMESTAMPTZ,
    
    -- Uploader
    uploaded_by     UUID REFERENCES users(id),
    
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_files_project ON files(project_id);
CREATE INDEX idx_files_parent ON files(parent_id);
CREATE INDEX idx_files_path ON files(path);
CREATE INDEX idx_files_uploaded_by ON files(uploaded_by);

-- =====================================================
-- Data Classification Table
-- =====================================================

CREATE TABLE data_classifications (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resource_type   VARCHAR(50) NOT NULL,
    resource_id     UUID NOT NULL,
    level           classification_level NOT NULL DEFAULT 'internal',
    owner_id        UUID REFERENCES users(id),
    approved_by     UUID REFERENCES users(id),
    approved_at     TIMESTAMPTZ,
    watermark       BOOLEAN DEFAULT false,
    export_limit    JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(resource_type, resource_id)
);

CREATE INDEX idx_data_classifications_resource ON data_classifications(resource_type, resource_id);
CREATE INDEX idx_data_classifications_level ON data_classifications(level);

-- =====================================================
-- Audit Logs Table
-- =====================================================

CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID REFERENCES users(id),
    action          VARCHAR(50) NOT NULL,
    resource_type   VARCHAR(50),
    resource_id     UUID,
    details         JSONB,
    ip_address      INET,
    user_agent      VARCHAR(500),
    status          VARCHAR(20),
    error_message   TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Create partitions for 2026
CREATE TABLE audit_logs_2026_01 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE audit_logs_2026_02 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE audit_logs_2026_03 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE audit_logs_2026_04 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
CREATE TABLE audit_logs_2026_05 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
CREATE TABLE audit_logs_2026_06 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE audit_logs_2026_07 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE audit_logs_2026_08 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');
CREATE TABLE audit_logs_2026_09 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');
CREATE TABLE audit_logs_2026_10 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');
CREATE TABLE audit_logs_2026_11 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE audit_logs_2026_12 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- =====================================================
-- Notifications Table
-- =====================================================

CREATE TABLE notifications (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type            VARCHAR(50) NOT NULL,
    title           VARCHAR(200) NOT NULL,
    content         TEXT,
    is_read         BOOLEAN DEFAULT false,
    related_id      VARCHAR(50),
    related_type    VARCHAR(50),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = false;

-- =====================================================
-- Announcements Table
-- =====================================================

CREATE TABLE announcements (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title           VARCHAR(200) NOT NULL,
    content         TEXT NOT NULL,
    author_id       UUID REFERENCES users(id),
    priority        VARCHAR(20) DEFAULT 'normal',
    is_pinned       BOOLEAN DEFAULT false,
    published_at    TIMESTAMPTZ,
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_announcements_published ON announcements(published_at DESC);

-- =====================================================
-- Honors Table
-- =====================================================

CREATE TABLE honors (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title           VARCHAR(200) NOT NULL,
    description     TEXT,
    award_year      INTEGER,
    award_month     INTEGER,
    recipient_id    UUID REFERENCES users(id),
    image_url       VARCHAR(500),
    sort_order      INTEGER DEFAULT 0,
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_honors_recipient ON honors(recipient_id);
CREATE INDEX idx_honors_active ON honors(is_active);

-- =====================================================
-- Default Admin User (password: admin123, MUST CHANGE ON FIRST LOGIN)
-- =====================================================

INSERT INTO organizations (id, name, code, level, sort_order) VALUES 
    (uuid_generate_v4(), '微波室', 'RD_DEPT', 1, 1),
    (uuid_generate_v4(), '产品管理组', 'PRODUCT_MGMT', 2, 1),
    (uuid_generate_v4(), '产品开发组', 'PRODUCT_DEV', 2, 2),
    (uuid_generate_v4(), '技术开发组', 'TECH_DEV', 2, 3),
    (uuid_generate_v4(), '综合管理组', 'GENERAL_MGMT', 2, 4);

-- Create default admin user
INSERT INTO users (
    id,
    username,
    display_name,
    email,
    role,
    team,
    password_hash,
    is_active,
    created_at
) VALUES (
    uuid_generate_v4(),
    'admin',
    '系统管理员',
    'admin@rdp.local',
    'admin',
    'general_mgmt',
    crypt('admin123', gen_salt('bf')),
    true,
    NOW()
);

-- =====================================================
-- Default Process Templates
-- =====================================================

INSERT INTO process_templates (id, name, code, category, description, activities, is_default, is_active) VALUES
    (
        uuid_generate_v4(),
        '单机产品开发流程',
        'PROCESS_STANDALONE',
        'standalone',
        '单机产品完整开发流程',
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
        true
    ),
    (
        uuid_generate_v4(),
        '模块开发流程',
        'PROCESS_MODULE',
        'module',
        '通用模块开发流程',
        '[
            {"id": "ACT001", "name": "需求分析", "duration": 7, "require_review": true},
            {"id": "ACT002", "name": "方案设计", "duration": 10, "require_review": true},
            {"id": "ACT003", "name": "详细设计", "duration": 15, "require_review": false},
            {"id": "ACT004", "name": "模块实现", "duration": 20, "require_review": false},
            {"id": "ACT005", "name": "测试验证", "duration": 10, "require_review": true},
            {"id": "ACT006", "name": "模块定型", "duration": 5, "require_review": true}
        ]'::jsonb,
        false,
        true
    ),
    (
        uuid_generate_v4(),
        '技术开发流程',
        'PROCESS_TECH',
        'tech_dev',
        '技术预研开发流程',
        '[
            {"id": "ACT001", "name": "技术调研", "duration": 15, "require_review": true},
            {"id": "ACT002", "name": "原理验证", "duration": 20, "require_review": true},
            {"id": "ACT003", "name": "方案设计", "duration": 15, "require_review": false},
            {"id": "ACT004", "name": "实验验证", "duration": 30, "require_review": true},
            {"id": "ACT005", "name": "技术总结", "duration": 10, "require_review": true}
        ]'::jsonb,
        false,
        true
    );

-- =====================================================
-- Function: Update updated_at timestamp
-- =====================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_organizations_updated_at BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_process_templates_updated_at BEFORE UPDATE ON process_templates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_activities_updated_at BEFORE UPDATE ON activities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_files_updated_at BEFORE UPDATE ON files
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_data_classifications_updated_at BEFORE UPDATE ON data_classifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- Grant Permissions
-- =====================================================

-- Grant table permissions
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO rdp_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO rdp_user;

-- Grant type permissions
GRANT USAGE ON TYPE user_role TO rdp_user;
GRANT USAGE ON TYPE team_type TO rdp_user;
GRANT USAGE ON TYPE project_category TO rdp_user;
GRANT USAGE ON TYPE project_status TO rdp_user;
GRANT USAGE ON TYPE activity_status TO rdp_user;
GRANT USAGE ON TYPE classification_level TO rdp_user;
GRANT USAGE ON TYPE title_level TO rdp_user;
GRANT USAGE ON TYPE pd_product_line TO rdp_user;
GRANT USAGE ON TYPE td_specialty TO rdp_user;
GRANT USAGE ON TYPE product_maturity TO rdp_user;
GRANT USAGE ON TYPE knowledge_category TO rdp_user;
GRANT USAGE ON TYPE trl_level TO rdp_user;

-- =====================================================
-- Comments
-- =====================================================

COMMENT ON TABLE users IS 'User table with role-based access control';
COMMENT ON TABLE projects IS 'Project management table';
COMMENT ON TABLE activities IS 'Project activity/task table';
COMMENT ON TABLE files IS 'Project file management with Git integration';
COMMENT ON TABLE audit_logs IS 'Security audit log with partitioning';
COMMENT ON TABLE notifications IS 'User notification system';
COMMENT ON TABLE process_templates IS 'Workflow template definitions';
