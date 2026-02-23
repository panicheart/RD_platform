-- Migration: 016_analytics
-- Description: Create analytics and reporting tables

-- Analytics dashboards table
CREATE TABLE analytics_dashboards (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    layout TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_by VARCHAR(26),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Project statistics table
CREATE TABLE project_stats (
    id VARCHAR(26) PRIMARY KEY,
    date DATE NOT NULL,
    total_projects INTEGER DEFAULT 0,
    active_projects INTEGER DEFAULT 0,
    completed_projects INTEGER DEFAULT 0,
    delayed_projects INTEGER DEFAULT 0,
    avg_progress DECIMAL(5,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User statistics table
CREATE TABLE user_stats (
    id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL,
    date DATE NOT NULL,
    tasks_completed INTEGER DEFAULT 0,
    tasks_created INTEGER DEFAULT 0,
    projects_joined INTEGER DEFAULT 0,
    reviews_done INTEGER DEFAULT 0,
    contribution DECIMAL(10,2) DEFAULT 0,
    work_hours DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Report templates table
CREATE TABLE report_templates (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) CHECK (type IN ('project', 'user', 'system')),
    format VARCHAR(20) CHECK (format IN ('pdf', 'excel')),
    template TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_by VARCHAR(26),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_project_stats_date ON project_stats(date);
CREATE INDEX idx_user_stats_user ON user_stats(user_id);
CREATE INDEX idx_user_stats_date ON user_stats(date);
CREATE INDEX idx_user_stats_user_date ON user_stats(user_id, date);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_analytics_dashboards_updated_at BEFORE UPDATE ON analytics_dashboards
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_report_templates_updated_at BEFORE UPDATE ON report_templates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
