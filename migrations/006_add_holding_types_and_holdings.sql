-- +goose Up
-- ============================================
-- Holding types table (global catalog)
-- ============================================
CREATE TABLE IF NOT EXISTS holding_types (
    id SMALLSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    notes TEXT
);

-- Seed default holding types
INSERT INTO holding_types (code, name, notes) VALUES
    ('STOCK', 'Stock', NULL),
    ('CRYPTO', 'Cryptocurrency', NULL),
    ('BOND', 'Bond', NULL),
    ('ETF', 'ETF', NULL),
    ('MUTUAL_FUND', 'Mutual Fund', NULL),
    ('REAL_ESTATE', 'Real Estate', NULL),
    ('COMMODITY', 'Commodity', NULL),
    ('OTHER', 'Other', NULL)
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- Holdings table
-- ============================================
CREATE TABLE IF NOT EXISTS holdings (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    symbol TEXT,
    platform TEXT NOT NULL,
    holding_type_id SMALLINT NOT NULL,
    currency CHAR(3) NOT NULL,
    invested_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
    current_value NUMERIC(18,2) NOT NULL DEFAULT 0,
    gain_amount NUMERIC(18,2) GENERATED ALWAYS AS (current_value - invested_amount) STORED,
    gain_percent NUMERIC(18,2) GENERATED ALWAYS AS (
        CASE
            WHEN invested_amount = 0 THEN 0
            ELSE ((current_value - invested_amount) / invested_amount) * 100
        END
    ) STORED,
    units NUMERIC(24,3),
    avg_buy_price NUMERIC(18,8),
    current_price NUMERIC(18,8),
    last_updated TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    month INTEGER NOT NULL DEFAULT 1,
    year INTEGER NOT NULL DEFAULT EXTRACT(YEAR FROM CURRENT_DATE)::INT,

    CONSTRAINT holdings_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT holdings_holding_type_id_fkey FOREIGN KEY (holding_type_id) REFERENCES holding_types(id) ON DELETE RESTRICT,
    CONSTRAINT chk_holdings_positive_amounts CHECK (invested_amount >= 0 AND current_value >= 0),
    CONSTRAINT chk_holdings_valid_month CHECK (month >= 1 AND month <= 12),
    CONSTRAINT chk_holdings_valid_year CHECK (year >= 2000)
);

CREATE INDEX IF NOT EXISTS idx_holdings_user ON holdings(user_id);
CREATE INDEX IF NOT EXISTS idx_holdings_holding_type_id ON holdings(holding_type_id);
CREATE INDEX IF NOT EXISTS idx_holdings_month_year ON holdings(year, month);
CREATE INDEX IF NOT EXISTS idx_holdings_user_month_year ON holdings(user_id, year, month);
CREATE INDEX IF NOT EXISTS idx_holdings_user_type ON holdings(user_id, holding_type_id);

-- +goose Down
DROP TABLE IF EXISTS holdings;
DROP TABLE IF EXISTS holding_types;