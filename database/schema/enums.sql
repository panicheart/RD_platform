-- Schema: Enumeration types for RDP platform
-- These enums are shared across multiple tables

-- User status enum
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'locked', 'pending');

-- Project status enum
CREATE TYPE project_status AS ENUM (
    'draft',
    'planning',
    'in_progress',
    'on_hold',
    'completed',
    'cancelled',
    'archived'
);

-- Project category enum (based on 7 categories from requirements)
CREATE TYPE project_category AS ENUM (
    'new_product',        -- 新产品开发
    'product_improvement', -- 产品改进
    'pre_research',       -- 预研项目
    'tech_platform',      -- 技术平台
    'component_development', -- 单机/模块开发
    'process_improvement',   -- 工艺改进
    'other'               -- 其他
);

-- Priority enum
CREATE TYPE priority_level AS ENUM ('low', 'medium', 'high', 'critical');

-- Data classification level (for security)
CREATE TYPE classification_level AS ENUM ('public', 'internal', 'confidential', 'secret');

-- Audit action type enum
CREATE TYPE audit_action AS ENUM ('create', 'read', 'update', 'delete', 'login', 'logout', 'export', 'import');

-- File status enum
CREATE TYPE file_status AS ENUM ('pending', 'active', 'archived', 'deleted');

-- Activity status enum (for workflow)
CREATE TYPE activity_status AS ENUM ('pending', 'in_progress', 'completed', 'blocked', 'cancelled');

-- Notification type enum
CREATE TYPE notification_type AS ENUM ('system', 'project', 'workflow', 'mention', 'deadline');

-- Organization type enum
CREATE TYPE org_type AS ENUM ('department', 'team', 'group', 'product_line');
