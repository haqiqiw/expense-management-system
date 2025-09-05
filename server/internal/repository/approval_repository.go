package repository

import (
	"context"
	"expense-management-system/internal/db"
	"expense-management-system/internal/entity"
	"time"
)

type ApprovalRepository struct {
	db db.PgxIface
}

func NewApprovalRepository(db db.PgxIface) *ApprovalRepository {
	return &ApprovalRepository{
		db: db,
	}
}

func (r *ApprovalRepository) CreateTx(ctx context.Context, exec db.Executor, approval *entity.Approval) error {
	now := time.Now()
	query := `
		INSERT INTO approvals (expense_id, approver_id, status, notes, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := exec.QueryRow(ctx, query,
		approval.ExpenseID,
		approval.ApproverID,
		approval.Status,
		approval.Notes,
		now,
	).Scan(&approval.ID)

	approval.CreatedAt = now

	return err
}
