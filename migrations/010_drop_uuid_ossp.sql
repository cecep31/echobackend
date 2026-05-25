-- +goose Up
-- uuid-ossp is no longer needed after switching defaults to native uuidv7() (PostgreSQL 18+).
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'user_tag_follows'
    ) THEN
        ALTER TABLE user_tag_follows ALTER COLUMN id SET DEFAULT uuidv7();
    END IF;
END $$;
-- +goose StatementEnd

DROP EXTENSION IF EXISTS "uuid-ossp";

-- +goose Down
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'user_tag_follows'
    ) THEN
        ALTER TABLE user_tag_follows ALTER COLUMN id SET DEFAULT uuid_generate_v4();
    END IF;
END $$;
-- +goose StatementEnd
