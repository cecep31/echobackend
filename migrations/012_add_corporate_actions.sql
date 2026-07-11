-- +goose Up
-- ============================================
-- Corporate actions table (dividend & RUPS calendar, persisted from IDX)
-- ============================================
CREATE TABLE IF NOT EXISTS corporate_actions (
    id BIGSERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    name TEXT,
    type TEXT NOT NULL,
    event_date DATE NOT NULL,
    pay_date DATE,
    amount NUMERIC(18,4),
    currency TEXT,
    note TEXT,
    market TEXT NOT NULL DEFAULT 'IDX',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_corporate_actions_symbol_type_date UNIQUE (symbol, type, event_date),
    CONSTRAINT chk_corporate_actions_type CHECK (type IN ('dividend', 'rups'))
);

CREATE INDEX IF NOT EXISTS idx_corporate_actions_event_date ON corporate_actions(event_date);

-- +goose Down
DROP TABLE IF EXISTS corporate_actions;
