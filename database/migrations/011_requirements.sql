-- Quality Management Tables Migration
-- Migration: 011_requirements.sql

-- Requirement type enum
CREATE TYPE requirement_type AS ENUM ('functional', 'non_functional', 'interface', 'safety');

-- Requirement status enum
CREATE TYPE requirement_status AS ENUM ('draft', 'reviewed', 'approved', 'rejected', 'implemented', 'verified');

-- Requirements table
CREATE TABLE requirements (
    id CHAR(26) PRIMARY KEY,
    project_id CHAR(26) NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    parent_id CHAR(26) REFERENCES requirements(id),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    type requirement_type NOT NULL,
    priority INTEGER DEFAULT 2 CHECK (priority >= 1 AND priority <= 4),
    status requirement_status NOT NULL DEFAULT 'draft',
    rationale TEXT,
    source VARCHAR(200),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(26) REFERENCES users(id)
);

CREATE INDEX idx_requirements_project_id ON requirements(project_id);
CREATE INDEX idx_requirements_status ON requirements(status);
CREATE INDEX idx_requirements_parent_id ON requirements(parent_id);

-- Change request type enum
CREATE TYPE change_request_type AS ENUM ('ecr', 'eco');

-- Change request status enum
CREATE TYPE change_request_status AS ENUM ('draft', 'submitted', 'evaluated', 'approved', 'rejected', 'implemented', 'closed');

-- Change requests table
CREATE TABLE change_requests (
    id CHAR(26) PRIMARY KEY,
    project_id CHAR(26) NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    type change_request_type NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    reason TEXT,
    impact_analysis TEXT,
    affected_items JSONB,
    status change_request_status NOT NULL DEFAULT 'draft',
    requester_id CHAR(26) REFERENCES users(id),
    approver_id CHAR(26) REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,
    implemented_at TIMESTAMP WITH TIME ZONE,
    closed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_change_requests_project_id ON change_requests(project_id);
CREATE INDEX idx_change_requests_status ON change_requests(status);
CREATE INDEX idx_change_requests_type ON change_requests(type);

-- Defect severity enum
CREATE TYPE defect_severity AS ENUM ('critical', 'high', 'medium', 'low');

-- Defect status enum
CREATE TYPE defect_status AS ENUM ('new', 'assigned', 'in_progress', 'resolved', 'closed', 'reopened');

-- Defects table
CREATE TABLE defects (
    id CHAR(26) PRIMARY KEY,
    project_id CHAR(26) NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    severity defect_severity NOT NULL,
    status defect_status NOT NULL DEFAULT 'new',
    reporter_id CHAR(26) REFERENCES users(id),
    assignee_id CHAR(26) REFERENCES users(id),
    resolution TEXT,
    reported_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE,
    closed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_defects_project_id ON defects(project_id);
CREATE INDEX idx_defects_status ON defects(status);
CREATE INDEX idx_defects_severity ON defects(severity);
CREATE INDEX idx_defects_assignee_id ON defects(assignee_id);

CREATE TRIGGER update_requirements_updated_at BEFORE UPDATE ON requirements
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_change_requests_updated_at BEFORE UPDATE ON change_requests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_defects_updated_at BEFORE UPDATE ON defects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
