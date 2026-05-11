-- +goose Up
-- Add bookmark_count column to posts table
ALTER TABLE posts ADD COLUMN IF NOT EXISTS bookmark_count BIGINT DEFAULT 0;

-- ============================================
-- Bookmark folders table
-- ============================================
CREATE TABLE IF NOT EXISTS bookmark_folders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_bookmark_folders_user_name ON bookmark_folders(user_id, name);
CREATE INDEX IF NOT EXISTS idx_bookmark_folders_user_id ON bookmark_folders(user_id);

ALTER TABLE bookmark_folders
    ADD CONSTRAINT fk_bookmark_folders_user_id
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- ============================================
-- Post bookmarks table
-- ============================================
CREATE TABLE IF NOT EXISTS post_bookmarks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,
    folder_id UUID,
    name VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_post_bookmarks_unique_user_post
    ON post_bookmarks(post_id, user_id);
CREATE INDEX IF NOT EXISTS idx_post_bookmarks_user_id ON post_bookmarks(user_id);
CREATE INDEX IF NOT EXISTS idx_post_bookmarks_folder_id ON post_bookmarks(folder_id);

ALTER TABLE post_bookmarks
    ADD CONSTRAINT fk_post_bookmarks_post_id
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;
ALTER TABLE post_bookmarks
    ADD CONSTRAINT fk_post_bookmarks_user_id
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE post_bookmarks
    ADD CONSTRAINT fk_post_bookmarks_folder_id
    FOREIGN KEY (folder_id) REFERENCES bookmark_folders(id) ON DELETE SET NULL;

-- ============================================
-- Trigger: auto-update posts.bookmark_count
-- ============================================
CREATE OR REPLACE FUNCTION update_post_bookmark_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE posts SET bookmark_count = bookmark_count + 1 WHERE id = NEW.post_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE posts SET bookmark_count = bookmark_count - 1 WHERE id = OLD.post_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_post_bookmark_count_insert
    AFTER INSERT ON post_bookmarks
    FOR EACH ROW
    EXECUTE FUNCTION update_post_bookmark_count();

CREATE TRIGGER trigger_update_post_bookmark_count_delete
    AFTER DELETE ON post_bookmarks
    FOR EACH ROW
    EXECUTE FUNCTION update_post_bookmark_count();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_update_post_bookmark_count_delete ON post_bookmarks;
DROP TRIGGER IF EXISTS trigger_update_post_bookmark_count_insert ON post_bookmarks;
DROP FUNCTION IF EXISTS update_post_bookmark_count();

ALTER TABLE post_bookmarks DROP CONSTRAINT IF EXISTS fk_post_bookmarks_folder_id;
ALTER TABLE post_bookmarks DROP CONSTRAINT IF EXISTS fk_post_bookmarks_user_id;
ALTER TABLE post_bookmarks DROP CONSTRAINT IF EXISTS fk_post_bookmarks_post_id;

ALTER TABLE bookmark_folders DROP CONSTRAINT IF EXISTS fk_bookmark_folders_user_id;

DROP TABLE IF EXISTS post_bookmarks;
DROP TABLE IF EXISTS bookmark_folders;

ALTER TABLE posts DROP COLUMN IF EXISTS bookmark_count;