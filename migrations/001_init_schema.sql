-- +goose Up
-- Enable uuid extension (required for uuid_generate_v4)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- Users table
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255),
    image TEXT,
    is_super_admin BOOLEAN DEFAULT FALSE,
    username VARCHAR(255),
    github_id BIGINT,
    last_logged_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_github_id_unique ON users(github_id) WHERE github_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- ============================================
-- Posts table
-- ============================================
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    title VARCHAR(255) NOT NULL,
    created_by UUID NOT NULL,
    body TEXT,
    slug VARCHAR(255) NOT NULL,
    photo_url TEXT,
    published BOOLEAN DEFAULT TRUE,
    published_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS creator_and_slug_unique ON posts(created_by, slug) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_posts_deleted_at ON posts(deleted_at);
CREATE INDEX IF NOT EXISTS idx_posts_created_by ON posts(created_by);

ALTER TABLE posts
    ADD CONSTRAINT fk_posts_created_by
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE;

-- ============================================
-- Tags table
-- ============================================
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_name ON tags(name);

-- ============================================
-- Posts-to-Tags join table
-- ============================================
CREATE TABLE IF NOT EXISTS posts_to_tags (
    post_id UUID NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, tag_id)
);

ALTER TABLE posts_to_tags
    ADD CONSTRAINT fk_posts_to_tags_post_id
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;
ALTER TABLE posts_to_tags
    ADD CONSTRAINT fk_posts_to_tags_tag_id
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE;

-- ============================================
-- Profiles table
-- ============================================
CREATE TABLE IF NOT EXISTS profiles (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    bio TEXT,
    website TEXT,
    phone VARCHAR(50),
    location VARCHAR(255)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id);

ALTER TABLE profiles
    ADD CONSTRAINT fk_profiles_user_id
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- ============================================
-- Sessions table
-- ============================================
CREATE TABLE IF NOT EXISTS sessions (
    refresh_token TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ,
    user_agent TEXT,
    expires_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);

ALTER TABLE sessions
    ADD CONSTRAINT fk_sessions_user_id
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- ============================================
-- Post comments table
-- ============================================
CREATE TABLE IF NOT EXISTS post_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    text TEXT NOT NULL,
    post_id UUID NOT NULL,
    parent_comment_id UUID,
    created_by UUID NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_post_comments_post_id ON post_comments(post_id);
CREATE INDEX IF NOT EXISTS idx_post_comments_created_by ON post_comments(created_by);
CREATE INDEX IF NOT EXISTS idx_post_comments_parent_comment_id ON post_comments(parent_comment_id);
CREATE INDEX IF NOT EXISTS idx_post_comments_deleted_at ON post_comments(deleted_at);

ALTER TABLE post_comments
    ADD CONSTRAINT fk_post_comments_post_id
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;
ALTER TABLE post_comments
    ADD CONSTRAINT fk_post_comments_created_by
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE post_comments
    ADD CONSTRAINT fk_post_comments_parent
    FOREIGN KEY (parent_comment_id) REFERENCES post_comments(id) ON DELETE CASCADE;

-- ============================================
-- Files table
-- ============================================
CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255),
    path TEXT,
    size INTEGER,
    type VARCHAR(255),
    created_by UUID
);

CREATE INDEX IF NOT EXISTS idx_files_deleted_at ON files(deleted_at);
CREATE INDEX IF NOT EXISTS idx_files_created_by ON files(created_by);

ALTER TABLE files
    ADD CONSTRAINT fk_files_created_by
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS post_comments;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS posts_to_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";