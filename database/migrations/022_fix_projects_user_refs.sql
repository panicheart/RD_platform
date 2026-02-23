-- Migration: 022_fix_projects_user_refs
-- Description: Change projects user references to support UUID format
-- Date: 2026-02-23

-- Drop tables first
DROP TABLE IF EXISTS project_members CASCADE;
DROP TABLE IF EXISTS projects CASCADE;

-- Recreate projects table with proper user ID references (VARCHAR(36))
CREATE TABLE projects (
    id VARCHAR(26) PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    product_line VARCHAR(50),
    team VARCHAR(50),
    process_template_id VARCHAR(26),
    start_date DATE,
    end_date DATE,
    actual_start_date DATE,
    actual_end_date DATE,
    progress INTEGER DEFAULT 0,
    git_repo_id VARCHAR(100),
    git_repo_url VARCHAR(500),
    classification_level VARCHAR(50) DEFAULT 'internal',
    leader_id VARCHAR(36),
    tech_leader_id VARCHAR(36),
    product_leader_id VARCHAR(36),
    created_by VARCHAR(36),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Recreate project_members table
CREATE TABLE project_members (
    id VARCHAR(26) PRIMARY KEY,
    project_id VARCHAR(26) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    role VARCHAR(50) DEFAULT 'member',
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, user_id)
);

-- Create indexes
CREATE INDEX idx_projects_code ON projects(code);
CREATE INDEX idx_projects_category ON projects(category);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_project_members_project ON project_members(project_id);
CREATE INDEX idx_project_members_user ON project_members(user_id);

-- Add trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert test data
INSERT INTO projects (
    id, code, name, description, category, status,
    start_date, end_date, progress, leader_id,
    classification_level, created_by, created_at, updated_at
) VALUES 
    (
        '01HPVJ8ZM00000000000000501',
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
        '01HPVJ8ZM00000000000000502',
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
        '01HPVJ8ZM00000000000000503',
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
    );

-- Verify data
SELECT 'Projects recreated' as info, COUNT(*) as count FROM projects;
