CREATE TYPE expense_status AS ENUM (
    'awaiting_approval',
    'approved',
    'rejected',
    'completed'
);

CREATE TABLE IF NOT EXISTS expenses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    description TEXT NOT NULL,
    receipt_url TEXT,
    status expense_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,

    CONSTRAINT fk_expenses_user_id
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CONSTRAINT amount_check CHECK (amount > 0)
);

CREATE INDEX idx_expenses_user_id_status ON expenses (user_id, status);

CREATE INDEX idx_expenses_status_user_id ON expenses (status, user_id);