-- +goose Up
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    data TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, read);
CREATE INDEX IF NOT EXISTS idx_notifications_user_read_created_at ON notifications(user_id, read, created_at DESC);

ALTER TABLE notifications
    ADD CONSTRAINT notifications_user_id_users_id_fk
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_user_id_users_id_fk;
DROP TABLE IF EXISTS notifications;
