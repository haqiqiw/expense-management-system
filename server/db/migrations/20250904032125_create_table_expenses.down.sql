DROP INDEX IF EXISTS idx_expenses_user_id_status;

DROP INDEX IF EXISTS idx_expenses_status_user_id;

DROP TABLE IF EXISTS expenses;

DROP TYPE IF EXISTS expense_status;