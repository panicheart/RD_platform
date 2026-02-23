-- Migration: 020_fix_forum_user_refs
-- Description: Fix forum tables to properly reference UUID user IDs
-- Date: 2026-02-23

-- Alter forum_posts table author_id to support UUID
ALTER TABLE forum_posts 
ALTER COLUMN author_id TYPE VARCHAR(36);

-- Alter forum_replies table author_id to support UUID
ALTER TABLE forum_replies
ALTER COLUMN author_id TYPE VARCHAR(36);

COMMENT ON COLUMN forum_posts.author_id IS 'References users.id (UUID format)';
COMMENT ON COLUMN forum_replies.author_id IS 'References users.id (UUID format)';
