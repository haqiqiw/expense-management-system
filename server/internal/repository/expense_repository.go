package repository

import (
	"context"
	"database/sql"
	"errors"
	"expense-management-system/internal/db"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type ExpenseRepository struct {
	db db.PgxIface
}

func NewExpenseRepository(db db.PgxIface) *ExpenseRepository {
	return &ExpenseRepository{
		db: db,
	}
}

func (r *ExpenseRepository) Create(ctx context.Context, expense *entity.Expense) error {
	now := time.Now()
	query := `
		INSERT INTO expenses (user_id, amount, description, receipt_url, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := r.db.QueryRow(ctx, query,
		expense.UserID,
		expense.Amount,
		expense.Description,
		expense.ReceiptURL,
		expense.Status,
		now,
	).Scan(&expense.ID)

	expense.CreatedAt = now

	return err
}

func (r *ExpenseRepository) List(ctx context.Context, req *model.ListExpenseRequest) ([]entity.ExpenseWithUser, int, error) {
	var whereClauses []string
	whereArgs := make([]any, 0)
	argCount := 1

	baseSelectQuery := `
		SELECT
			e.id AS expense_id, e.user_id AS expense_user_id, e.amount AS expense_amount, e.description AS expense_description, 
			e.receipt_url AS expense_receipt_url, e.status AS expense_status, e.created_at AS expense_created_at, e.processed_at AS expense_processed_at,
			u.id AS user_id, u.email AS user_email, u.name AS user_name
		FROM expenses AS e
		JOIN users AS u ON e.user_id = u.id`

	baseCountQuery := `SELECT COUNT(*) FROM expenses AS e`

	switch req.View {
	case model.ExpenseViewPersonal:
		whereClauses = append(whereClauses, fmt.Sprintf("e.user_id = $%d", argCount))
		whereArgs = append(whereArgs, req.UserID)
		argCount++

		if req.Status != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("e.status = $%d", argCount))
			whereArgs = append(whereArgs, *req.Status)
			argCount++
		}

		if req.AutoApproved {
			whereClauses = append(whereClauses, fmt.Sprintf("e.amount < %d", entity.ApprovalThresholdAmount))
		}

	case model.ExpenseViewApprovalQueue:
		whereClauses = append(whereClauses, "e.status = 'awaiting_approval'")
		whereClauses = append(whereClauses, fmt.Sprintf("e.user_id != $%d", argCount))
		whereArgs = append(whereArgs, req.UserID)
		argCount++

	default:
		return nil, 0, errors.New("invalid view")
	}

	whereQuery := ""
	if len(whereClauses) > 0 {
		whereQuery = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var total int
	countQuery := baseCountQuery + whereQuery
	err := r.db.QueryRow(ctx, countQuery, whereArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, nil
	}

	selectQuery := baseSelectQuery + whereQuery
	selectQuery += fmt.Sprintf(" ORDER BY e.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	selectArgs := append(whereArgs, req.Limit, req.Offset)

	rows, err := r.db.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []entity.ExpenseWithUser
	for rows.Next() {
		var eu entity.ExpenseWithUser
		err := rows.Scan(
			&eu.Expense.ID, &eu.Expense.UserID, &eu.Expense.Amount, &eu.Expense.Description,
			&eu.Expense.ReceiptURL, &eu.Expense.Status, &eu.Expense.CreatedAt, &eu.Expense.ProcessedAt,
			&eu.User.ID, &eu.User.Email, &eu.User.Name,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, eu)
	}

	return results, total, nil
}

func (r *ExpenseRepository) FindDetailByID(ctx context.Context, id uint64) (*entity.ExpenseDetail, error) {
	query := `
		SELECT
			e.id AS expense_id, e.user_id AS expense_user_id, e.amount AS expense_amount, e.description AS expense_description, 
			e.receipt_url AS expense_receipt_url, e.status AS expense_status, e.created_at AS expense_created_at, e.processed_at AS expense_processed_at,
			ue.id AS user_id, ue.email AS user_email, ue.name AS user_name,
			a.id AS approval_id, a.approver_id, ua.email AS approver_email, ua.name AS approver_name,
			a.status AS approval_status, a.notes AS approval_notes, a.created_at AS approval_created_at
		FROM expenses AS e
		JOIN users AS ue ON e.user_id = ue.id
		LEFT JOIN approvals AS a ON e.id = a.expense_id
		LEFT JOIN users AS ua ON a.approver_id = ua.id
		WHERE e.id = $1`

	var detail entity.ExpenseDetail

	var approvalID, approvalApproverID sql.NullInt64
	var approvalStatus sql.NullString
	var approvalApproverEmail, approvalApproverName, approvalNotes sql.NullString
	var approvalCreatedAt sql.NullTime

	err := r.db.QueryRow(ctx, query, id).Scan(
		&detail.Expense.ID, &detail.Expense.UserID, &detail.Expense.Amount, &detail.Expense.Description,
		&detail.Expense.ReceiptURL, &detail.Expense.Status, &detail.Expense.CreatedAt, &detail.Expense.ProcessedAt,
		&detail.User.ID, &detail.User.Email, &detail.User.Name,
		&approvalID, &approvalApproverID, &approvalApproverEmail, &approvalApproverName, &approvalStatus, &approvalNotes, &approvalCreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if approvalID.Valid {
		detail.Approval = &entity.ApprovalDetail{
			ID:            uint64(approvalID.Int64),
			ApproverID:    uint64(approvalApproverID.Int64),
			ApproverEmail: approvalApproverEmail.String,
			ApproverName:  approvalApproverName.String,
			Status:        entity.ApprovalStatus(approvalStatus.String),
			Notes:         nullableStringPtr(approvalNotes),
			CreatedAt:     approvalCreatedAt.Time,
		}
	}

	return &detail, nil
}

func (r *ExpenseRepository) FindByID(ctx context.Context, id uint64) (*entity.Expense, error) {
	query := `SELECT id, user_id, amount, description, receipt_url, status, created_at, processed_at FROM expenses WHERE id = $1 LIMIT 1`

	var e entity.Expense
	err := r.db.QueryRow(ctx, query, id).Scan(&e.ID, &e.UserID, &e.Amount, &e.Description, &e.ReceiptURL, &e.Status, &e.CreatedAt, &e.ProcessedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &e, nil
}

func (r *ExpenseRepository) FindByIDWithLock(ctx context.Context, exec db.Executor, id uint64) (*entity.Expense, error) {
	query := `SELECT id, user_id, amount, description, receipt_url, status, created_at, processed_at FROM expenses WHERE id = $1 FOR UPDATE`

	var e entity.Expense
	err := exec.QueryRow(ctx, query, id).Scan(&e.ID, &e.UserID, &e.Amount, &e.Description, &e.ReceiptURL, &e.Status, &e.CreatedAt, &e.ProcessedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &e, nil
}

func (r *ExpenseRepository) UpdateStatusByIDTx(ctx context.Context, exec db.Executor, id uint64, status entity.ExpenseStatus) error {
	query := `UPDATE expenses SET status = $1 WHERE id = $2`

	_, err := exec.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *ExpenseRepository) UpdateStatusByID(ctx context.Context, id uint64, status entity.ExpenseStatus, processedAt time.Time) error {
	query := `UPDATE expenses SET status = $1, processed_at = $2 WHERE id = $3`

	_, err := r.db.Exec(ctx, query, status, processedAt, id)
	if err != nil {
		return err
	}

	return nil
}

func nullableStringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
