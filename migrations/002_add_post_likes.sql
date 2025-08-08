-- Add like_count column to posts table
ALTER TABLE posts ADD COLUMN like_count BIGINT DEFAULT 0;

-- Create post_likes table
CREATE TABLE IF NOT EXISTS post_likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for post_likes
CREATE INDEX IF NOT EXISTS idx_post_likes_post_id ON post_likes(post_id);
CREATE INDEX IF NOT EXISTS idx_post_likes_user_id ON post_likes(user_id);
CREATE INDEX IF NOT EXISTS idx_post_likes_deleted_at ON post_likes(deleted_at);
CREATE INDEX IF NOT EXISTS idx_post_likes_created_at ON post_likes(created_at);

-- Create unique index to prevent duplicate likes from same user on same post
CREATE UNIQUE INDEX IF NOT EXISTS idx_post_likes_unique_user_post 
ON post_likes(post_id, user_id) 
WHERE deleted_at IS NULL;

-- Add foreign key constraints
ALTER TABLE post_likes 
ADD CONSTRAINT fk_post_likes_post_id 
FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;

ALTER TABLE post_likes 
ADD CONSTRAINT fk_post_likes_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Create function to update like_count when likes are added/removed
CREATE OR REPLACE FUNCTION update_post_like_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Increment like count
        UPDATE posts 
        SET like_count = like_count + 1 
        WHERE id = NEW.post_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        -- Decrement like count
        UPDATE posts 
        SET like_count = like_count - 1 
        WHERE id = OLD.post_id;
        RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
        -- Handle soft delete (when deleted_at changes)
        IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
            -- Soft delete: decrement count
            UPDATE posts 
            SET like_count = like_count - 1 
            WHERE id = NEW.post_id;
        ELSIF OLD.deleted_at IS NOT NULL AND NEW.deleted_at IS NULL THEN
            -- Restore: increment count
            UPDATE posts 
            SET like_count = like_count + 1 
            WHERE id = NEW.post_id;
        END IF;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers to automatically update like_count
CREATE TRIGGER trigger_update_post_like_count_insert
    AFTER INSERT ON post_likes
    FOR EACH ROW
    EXECUTE FUNCTION update_post_like_count();

CREATE TRIGGER trigger_update_post_like_count_delete
    AFTER DELETE ON post_likes
    FOR EACH ROW
    EXECUTE FUNCTION update_post_like_count();

CREATE TRIGGER trigger_update_post_like_count_update
    AFTER UPDATE ON post_likes
    FOR EACH ROW
    EXECUTE FUNCTION update_post_like_count();