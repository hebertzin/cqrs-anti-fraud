CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS accounts (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID        NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'blocked', 'flagged')),
    risk_level  DECIMAL(5, 4) NOT NULL DEFAULT 0,
    blocked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts (user_id);
CREATE INDEX IF NOT EXISTS idx_accounts_status   ON accounts (status);

CREATE TABLE IF NOT EXISTS transactions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id  UUID          NOT NULL REFERENCES accounts (id),
    amount      DECIMAL(20,2) NOT NULL,
    currency    CHAR(3)       NOT NULL,
    merchant_id VARCHAR(255)  NOT NULL,
    location    VARCHAR(10)   NOT NULL,
    status      VARCHAR(20)   NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'approved', 'declined', 'flagged')),
    risk_score  DECIMAL(5, 4) NOT NULL DEFAULT 0,
    metadata    JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_account_id  ON transactions (account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status      ON transactions (status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at  ON transactions (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_risk_score  ON transactions (risk_score DESC);
