-- Seed Data: Process Templates
-- Description: Predefined process templates for 7 project categories
-- Created: 2026-02-22

INSERT INTO process_templates (id, name, category, description, stages, is_active, created_at, updated_at) VALUES
('01H', '新产品开发流程', 'new_product', '完整的新产品开发流程，包含概念阶段、计划阶段、开发阶段、验证阶段、发布阶段', 
 '["概念阶段","计划阶段","开发阶段","验证阶段","发布阶段"]'::jsonb, true, NOW(), NOW()),

('02H', '产品改进流程', 'product_improvement', '针对现有产品的功能改进或性能优化', 
 '["需求确认","方案设计","开发实现","测试验证","发布上线"]'::jsonb, true, NOW(), NOW()),

('03H', '预研项目流程', 'pre_research', '技术研究与可行性验证', 
 '["技术调研","方案设计","原型开发","技术评审","成果归档"]'::jsonb, true, NOW(), NOW()),

('04H', '技术平台流程', 'tech_platform', '通用技术平台或基础组件开发', 
 '["需求分析","架构设计","核心开发","集成测试","平台发布"]'::jsonb, true, NOW(), NOW()),

('05H', '单机模块流程', 'component_development', '单机或功能模块开发', 
 '["需求分析","详细设计","模块开发","单元测试","集成验证"]'::jsonb, true, NOW(), NOW()),

('06H', '工艺改进流程', 'process_improvement', '生产工艺或流程改进', 
 '["问题定义","方案论证","试点实施","效果评估","推广固化"]'::jsonb, true, NOW(), NOW()),

('07H', '通用项目流程', 'other', '适用于其他类型项目的通用流程', 
 '["启动","计划","执行","监控","收尾"]'::jsonb, true, NOW(), NOW());
