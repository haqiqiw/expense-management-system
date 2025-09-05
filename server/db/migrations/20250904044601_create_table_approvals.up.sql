CREATE TYPE approval_status AS ENUM (
    'approved',
    'rejected'
);

CREATE TABLE IF NOT EXISTS approvals (
    id BIGSERIAL PRIMARY KEY,
    expense_id BIGINT UNIQUE NOT NULL,
    approver_id BIGINT NOT NULL,
    status approval_status NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_approvals_expense_id
        FOREIGN KEY(expense_id)
        REFERENCES expenses(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_approvals_approver_id
        FOREIGN KEY(approver_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
);