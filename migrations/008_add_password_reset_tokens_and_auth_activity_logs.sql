-- +goose Up
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token ON password_reset_tokens(token);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);

ALTER TABLE password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_user_id_users_id_fk
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS auth_activity_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    activity_type VARCHAR(50) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'success',
    error_message TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_auth_activity_logs_user_id ON auth_activity_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_activity_logs_activity_type ON auth_activity_logs(activity_type);
CREATE INDEX IF NOT EXISTS idx_auth_activity_logs_created_at ON auth_activity_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_auth_activity_logs_user_activity ON auth_activity_logs(user_id, activity_type);
CREATE INDEX IF NOT EXISTS idx_auth_activity_logs_user_status_created ON auth_activity_logs(user_id, status, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS auth_activity_logs;
DROP TABLE IF EXISTS password_reset_tokens;