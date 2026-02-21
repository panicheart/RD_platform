-- Git Repository Integration Migration
-- Migration: 006_git_repos.sql

CREATE TABLE git_repos (
    id CHAR(26) PRIMARY KEY,
    project_id CHAR(26) NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    repo_name VARCHAR(200) NOT NULL,
    repo_url VARCHAR(500) NOT NULL,
    repo_full_name VARCHAR(500) NOT NULL,
    default_branch VARCHAR(100) DEFAULT 'main',
    gitea_id BIGINT,
    is_private BOOLEAN DEFAULT true,
    initialized BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_git_repos_project_id ON git_repos(project_id);
CREATE INDEX idx_git_repos_repo_name ON git_repos(repo_name);

CREATE TRIGGER update_git_repos_updated_at BEFORE UPDATE ON git_repos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE git_repos IS 'Git repository links for projects';
