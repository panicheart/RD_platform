-- Workflow and Activity Tables Migration
-- Migration: 005_workflows.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Workflow states enum
CREATE TYPE workflow_state AS ENUM (
    'draft',
    'planning',
    'executing',
    'reviewing',
    'completed',
    'paused',
    'cancelled'
);

-- Activity status enum
CREATE TYPE activity_status AS ENUM (
    'pending',
    'ready',
    'running',
    'completed',
    'reviewing',
    'approved',
    'rejected',
    'skipped',
    'blocked'
);

-- Activity type enum
CREATE TYPE activity_type AS ENUM (
    'task',
    'milestone',
    'dcp',
    'review',
    'approval'
);

-- Review status enum
CREATE TYPE review_status AS ENUM (
    'pending',
    'submitted',
    'approved',
    'rejected',
    'revision'
);

-- Review type enum
CREATE TYPE review_type AS ENUM (
    'dcp',
    'code',
    'doc',
    'final'
);

-- Workflows table
CREATE TABLE workflows (
    id CHAR(26) PRIMARY KEY,
    project_id CHAR(26) NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    template_id CHAR(26) REFERENCES process_templates(id),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    state workflow_state NOT NULL DEFAULT 'draft',
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(26) REFERENCES users(id)
);

CREATE INDEX idx_workflows_project_id ON workflows(project_id);
CREATE INDEX idx_workflows_state ON workflows(state);
CREATE INDEX idx_workflows_created_at ON workflows(created_at);

-- Activities table
CREATE TABLE activities (
    id CHAR(26) PRIMARY KEY,
    workflow_id CHAR(26) NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    project_id CHAR(26) NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    parent_id CHAR(26) REFERENCES activities(id),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    type activity_type NOT NULL DEFAULT 'task',
    status activity_status NOT NULL DEFAULT 'pending',
    sequence INTEGER DEFAULT 0,
    priority INTEGER DEFAULT 0,
    assignee_id CHAR(26) REFERENCES users(id),
    planned_start TIMESTAMP WITH TIME ZONE,
    planned_end TIMESTAMP WITH TIME ZONE,
    actual_start TIMESTAMP WITH TIME ZONE,
    actual_end TIMESTAMP WITH TIME ZONE,
    progress INTEGER DEFAULT 0 CHECK (progress >= 0 AND progress <= 100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(26) REFERENCES users(id)
);

CREATE INDEX idx_activities_workflow_id ON activities(workflow_id);
CREATE INDEX idx_activities_project_id ON activities(project_id);
CREATE INDEX idx_activities_status ON activities(status);
CREATE INDEX idx_activities_assignee_id ON activities(assignee_id);
CREATE INDEX idx_activities_sequence ON activities(workflow_id, sequence);

-- Activity dependencies table
CREATE TABLE activity_dependencies (
    id CHAR(26) PRIMARY KEY,
    activity_id CHAR(26) NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    depends_on_id CHAR(26) NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    dependency_type VARCHAR(50) DEFAULT 'finish_to_start',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(activity_id, depends_on_id)
);

CREATE INDEX idx_activity_dependencies_activity_id ON activity_dependencies(activity_id);
CREATE INDEX idx_activity_dependencies_depends_on_id ON activity_dependencies(depends_on_id);

-- Deliverables table
CREATE TABLE deliverables (
    id CHAR(26) PRIMARY KEY,
    activity_id CHAR(26) NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    type VARCHAR(50),
    file_path VARCHAR(500),
    status VARCHAR(50) DEFAULT 'pending',
    submitted_at TIMESTAMP WITH TIME ZONE,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_deliverables_activity_id ON deliverables(activity_id);

-- Reviews table
CREATE TABLE reviews (
    id CHAR(26) PRIMARY KEY,
    activity_id CHAR(26) NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    project_id CHAR(26) NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    type review_type NOT NULL DEFAULT 'dcp',
    status review_status NOT NULL DEFAULT 'pending',
    reviewer_id CHAR(26) REFERENCES users(id),
    comments TEXT,
    score INTEGER CHECK (score >= 0 AND score <= 100),
    submitted_at TIMESTAMP WITH TIME ZONE,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(26) REFERENCES users(id)
);

CREATE INDEX idx_reviews_activity_id ON reviews(activity_id);
CREATE INDEX idx_reviews_project_id ON reviews(project_id);
CREATE INDEX idx_reviews_status ON reviews(status);
CREATE INDEX idx_reviews_reviewer_id ON reviews(reviewer_id);

-- Feedbacks table
CREATE TABLE feedbacks (
    id CHAR(26) PRIMARY KEY,
    review_id CHAR(26) NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    parent_id CHAR(26) REFERENCES feedbacks(id),
    content TEXT NOT NULL,
    author_id CHAR(26) NOT NULL REFERENCES users(id),
    mentions TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_feedbacks_review_id ON feedbacks(review_id);
CREATE INDEX idx_feedbacks_parent_id ON feedbacks(parent_id);
CREATE INDEX idx_feedbacks_author_id ON feedbacks(author_id);

-- Update timestamp trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_workflows_updated_at BEFORE UPDATE ON workflows
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_activities_updated_at BEFORE UPDATE ON activities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_deliverables_updated_at BEFORE UPDATE ON deliverables
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_feedbacks_updated_at BEFORE UPDATE ON feedbacks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE workflows IS 'Project workflow instances';
COMMENT ON TABLE activities IS 'Workflow activities and tasks';
COMMENT ON TABLE activity_dependencies IS 'Activity dependency relationships';
COMMENT ON TABLE deliverables IS 'Activity deliverables and outputs';
COMMENT ON TABLE reviews IS 'Activity review records including DCP';
COMMENT ON TABLE feedbacks IS 'Review feedback and comments';
