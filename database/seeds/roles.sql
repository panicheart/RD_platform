-- Seed Data: Default Roles and Permissions
-- Description: Initialize default RBAC roles
-- Created: 2026-02-22

INSERT INTO roles (id, name, code, description, permissions, created_at, updated_at) VALUES
('01ROLE00000000000000000001', '系统管理员', 'admin', '系统管理员，拥有所有权限', 
 '["*"]'::jsonb, NOW(), NOW()),

('01ROLE00000000000000000002', '部门领导', 'dept_leader', '部门领导，可查看部门所有项目',
 '["project:*:read","project:*:approve","user:dept:read","report:dept:*"]'::jsonb, NOW(), NOW()),

('01ROLE00000000000000000003', '项目经理', 'project_manager', '项目经理，管理所属项目',
 '["project:own:*","team:own:*","report:project:read"]'::jsonb, NOW(), NOW()),

('01ROLE00000000000000000004', '研发工程师', 'engineer', '研发工程师，参与项目开发',
 '["project:assigned:read","activity:own:*","file:project:read"]'::jsonb, NOW(), NOW()),

('01ROLE00000000000000000005', '质量工程师', 'quality_engineer', '质量工程师，负责质量管理',
 '["quality:*","audit:read","report:quality:*"]'::jsonb, NOW(), NOW()),

('01ROLE00000000000000000006', '资料管理员', 'knowledge_manager', '资料管理员，管理知识库',
 '["knowledge:*","shelf:read","forum:moderate"]'::jsonb, NOW(), NOW());
