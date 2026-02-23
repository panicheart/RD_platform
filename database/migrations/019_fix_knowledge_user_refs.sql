-- Migration: 019_fix_knowledge_user_refs
-- Description: Fix knowledge table to properly reference UUID user IDs
-- Date: 2026-02-23

-- Alter knowledge table author_id to support UUID
ALTER TABLE knowledge 
ALTER COLUMN author_id TYPE VARCHAR(36);

-- Alter knowledge_reviews table reviewer_id to support UUID  
ALTER TABLE knowledge_reviews
ALTER COLUMN reviewer_id TYPE VARCHAR(36);

-- Add foreign key constraints if they don't exist
-- Note: We keep these as loose references since users table uses UUID type
-- and knowledge tables use VARCHAR. The application layer should enforce integrity.

COMMENT ON COLUMN knowledge.author_id IS 'References users.id (UUID format)';
COMMENT ON COLUMN knowledge_reviews.reviewer_id IS 'References users.id (UUID format)';
