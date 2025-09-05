package repository_test

import (
	"context"
	"errors"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/model"
	"expense-management-system/internal/repository"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/suite"
)

type ExpenseRepositorySuite struct {
	suite.Suite
	mock pgxmock.PgxPoolIface
	repo *repository.ExpenseRepository
	ctx  context.Context
	now  time.Time
}

func (s *ExpenseRepositorySuite) SetupTest() {
	s.mock, _ = pgxmock.NewPool()
	s.repo = repository.NewExpenseRepository(s.mock)
	s.ctx = context.Background()
	s.now = time.Now()
}

func (s *ExpenseRepositorySuite) TearDownTest() {
	s.mock.Close()
}

func (s *ExpenseRepositorySuite) TestExpenseRepository_Create() {
	description := "dummy description"
	receiptUrl := "https://example.com/receipt.jpg"

	tests := []struct {
		name     string
		mockFunc func(pgxmock.PgxPoolIface)
		param    *entity.Expense
		wantErr  error
	}{
		{
			name: "error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO expenses (user_id, amount, description, receipt_url, status, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
				)).
					WithArgs(uint64(1), uint64(15000), description, &receiptUrl, pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("something error"))
			},
			param: &entity.Expense{
				UserID:      uint64(1),
				Amount:      uint64(15000),
				Description: description,
				ReceiptURL:  &receiptUrl,
				Status:      entity.ExpenseStatusApproved,
			},
			wantErr: errors.New("something error"),
		},
		{
			name: "success",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO expenses (user_id, amount, description, receipt_url, status, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
				)).
					WithArgs(uint64(1), uint64(15000), description, &receiptUrl, pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(uint64(1)))
			},
			param: &entity.Expense{
				UserID:      uint64(1),
				Amount:      uint64(15000),
				Description: description,
				ReceiptURL:  &receiptUrl,
				Status:      entity.ExpenseStatusApproved,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			err := s.repo.Create(s.ctx, tt.param)
			s.Equal(tt.wantErr, err)
		})
	}
}

func (s *ExpenseRepositorySuite) TestExpenseRepository_List() {
	now := time.Date(2025, 8, 13, 10, 0, 0, 0, time.UTC)
	description := "dummy description"
	status := "approved"

	tests := []struct {
		name      string
		mockFunc  func(pgxmock.PgxPoolIface)
		param     *model.ListExpenseRequest
		wantRes   []entity.ExpenseWithUser
		wantTotal int
		wantErr   error
	}{
		{
			name: "not found",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				countQuery := `SELECT COUNT(*) FROM expenses AS e WHERE e.user_id = $1`
				m.ExpectQuery(regexp.QuoteMeta(countQuery)).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))
			},
			param: &model.ListExpenseRequest{
				UserID:   uint64(1),
				UserRole: "manager",
				View:     model.ExpenseViewPersonal,
				Limit:    10,
				Offset:   0,
			},
			wantRes:   nil,
			wantTotal: 0,
			wantErr:   nil,
		},
		{
			name: "count error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				countQuery := `SELECT COUNT(*) FROM expenses AS e WHERE e.user_id = $1`
				m.ExpectQuery(regexp.QuoteMeta(countQuery)).
					WithArgs(uint64(1)).
					WillReturnError(errors.New("something error"))
			},
			param: &model.ListExpenseRequest{
				UserID:   uint64(1),
				UserRole: "manager",
				View:     model.ExpenseViewPersonal,
				Limit:    10,
				Offset:   0,
			},
			wantRes:   nil,
			wantTotal: 0,
			wantErr:   errors.New("something error"),
		},
		{
			name: "select error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				countQuery := `SELECT COUNT(*) FROM expenses AS e WHERE e.user_id = $1`
				selectQuery := `
		SELECT
			e.id AS expense_id, e.user_id AS expense_user_id, e.amount AS expense_amount, e.description AS expense_description,
			e.receipt_url AS expense_receipt_url, e.status AS expense_status, e.created_at AS expense_created_at, e.processed_at AS expense_processed_at,
			u.id AS user_id, u.email AS user_email, u.name AS user_name
		FROM expenses AS e
		JOIN users AS u ON e.user_id = u.id WHERE e.user_id = $1 ORDER BY e.created_at DESC LIMIT $2 OFFSET $3`

				m.ExpectQuery(regexp.QuoteMeta(countQuery)).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))
				m.ExpectQuery(regexp.QuoteMeta(selectQuery)).
					WithArgs(uint64(1), 10, 0).
					WillReturnError(errors.New("something error"))
			},
			param: &model.ListExpenseRequest{
				UserID:   uint64(1),
				UserRole: "manager",
				View:     model.ExpenseViewPersonal,
				Limit:    10,
				Offset:   0,
			},
			wantRes:   nil,
			wantTotal: 0,
			wantErr:   errors.New("something error"),
		},
		{
			name: "success personal with params user_id",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				countQuery := `SELECT COUNT(*) FROM expenses AS e WHERE e.user_id = $1`
				selectQuery := `
		SELECT
			e.id AS expense_id, e.user_id AS expense_user_id, e.amount AS expense_amount, e.description AS expense_description,
			e.receipt_url AS expense_receipt_url, e.status AS expense_status, e.created_at AS expense_created_at, e.processed_at AS expense_processed_at,
			u.id AS user_id, u.email AS user_email, u.name AS user_name
		FROM expenses AS e
		JOIN users AS u ON e.user_id = u.id WHERE e.user_id = $1 ORDER BY e.created_at DESC LIMIT $2 OFFSET $3`

				m.ExpectQuery(regexp.QuoteMeta(countQuery)).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))
				rows := pgxmock.NewRows([]string{
					"expense_id", "expense_user_id", "expense_amount", "expense_description",
					"expense_receipt_url", "expense_status", "expense_created_at", "expense_processed_at",
					"user_id", "user_email", "user_name",
				}).AddRow(
					uint64(1), uint64(1), uint64(15000), description,
					nil, entity.ExpenseStatusApproved, now, nil,
					uint64(1), "john@mail.com", "John Doe",
				)
				m.ExpectQuery(regexp.QuoteMeta(selectQuery)).
					WithArgs(uint64(1), 10, 0).
					WillReturnRows(rows)
			},
			param: &model.ListExpenseRequest{
				UserID:   uint64(1),
				UserRole: "manager",
				View:     model.ExpenseViewPersonal,
				Limit:    10,
				Offset:   0,
			},
			wantRes: []entity.ExpenseWithUser{
				{
					Expense: entity.Expense{
						ID:          uint64(1),
						UserID:      uint64(1),
						Amount:      uint64(15000),
						Description: description,
						ReceiptURL:  nil,
						Status:      entity.ExpenseStatusApproved,
						CreatedAt:   now,
						ProcessedAt: nil,
					},
					User: entity.UserSimple{
						ID:    1,
						Email: "john@mail.com",
						Name:  "John Doe",
					},
				},
			},
			wantTotal: 1,
			wantErr:   nil,
		},
		{
			name: "success personal with params user_id and status",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				countQuery := `SELECT COUNT(*) FROM expenses AS e WHERE e.user_id = $1 AND e.status = $2`
				selectQuery := `
		SELECT
			e.id AS expense_id, e.user_id AS expense_user_id, e.amount AS expense_amount, e.description AS expense_description,
			e.receipt_url AS expense_receipt_url, e.status AS expense_status, e.created_at AS expense_created_at, e.processed_at AS expense_processed_at,
			u.id AS user_id, u.email AS user_email, u.name AS user_name
		FROM expenses AS e
		JOIN users AS u ON e.user_id = u.id WHERE e.user_id = $1 AND e.status = $2 ORDER BY e.created_at DESC LIMIT $3 OFFSET $4`

				m.ExpectQuery(regexp.QuoteMeta(countQuery)).
					WithArgs(uint64(1), "approved").
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))
				rows := pgxmock.NewRows([]string{
					"expense_id", "expense_user_id", "expense_amount", "expense_description",
					"expense_receipt_url", "expense_status", "expense_created_at", "expense_processed_at",
					"user_id", "user_email", "user_name",
				}).AddRow(
					uint64(1), uint64(1), uint64(15000), description,
					nil, entity.ExpenseStatusApproved, now, nil,
					uint64(1), "john@mail.com", "John Doe",
				)
				m.ExpectQuery(regexp.QuoteMeta(selectQuery)).
					WithArgs(uint64(1), "approved", 10, 0).
					WillReturnRows(rows)
			},
			param: &model.ListExpenseRequest{
				UserID:   uint64(1),
				UserRole: "manager",
				View:     model.ExpenseViewPersonal,
				Status:   &status,
				Limit:    10,
				Offset:   0,
			},
			wantRes: []entity.ExpenseWithUser{
				{
					Expense: entity.Expense{
						ID:          uint64(1),
						UserID:      uint64(1),
						Amount:      uint64(15000),
						Description: description,
						ReceiptURL:  nil,
						Status:      entity.ExpenseStatusApproved,
						CreatedAt:   now,
						ProcessedAt: nil,
					},
					User: entity.UserSimple{
						ID:    1,
						Email: "john@mail.com",
						Name:  "John Doe",
					},
				},
			},
			wantTotal: 1,
			wantErr:   nil,
		},
		{
			name: "success personal with params user_id and status and auto_approved",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				countQuery := `SELECT COUNT(*) FROM expenses AS e WHERE e.user_id = $1 AND e.status = $2 AND e.amount < 1000000`
				selectQuery := `
		SELECT
			e.id AS expense_id, e.user_id AS expense_user_id, e.amount AS expense_amount, e.description AS expense_description,
			e.receipt_url AS expense_receipt_url, e.status AS expense_status, e.created_at AS expense_created_at, e.processed_at AS expense_processed_at,
			u.id AS user_id, u.email AS user_email, u.name AS user_name
		FROM expenses AS e
		JOIN users AS u ON e.user_id = u.id WHERE e.user_id = $1 AND e.status = $2 AND e.amount < 1000000 ORDER BY e.created_at DESC LIMIT $3 OFFSET $4`

				m.ExpectQuery(regexp.QuoteMeta(countQuery)).
					WithArgs(uint64(1), "approved").
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))
				rows := pgxmock.NewRows([]string{
					"expense_id", "expense_user_id", "expense_amount", "expense_description",
					"expense_receipt_url", "expense_status", "expense_created_at", "expense_processed_at",
					"user_id", "user_email", "user_name",
				}).AddRow(
					uint64(1), uint64(1), uint64(15000), description,
					nil, entity.ExpenseStatusApproved, now, nil,
					uint64(1), "john@mail.com", "John Doe",
				)
				m.ExpectQuery(regexp.QuoteMeta(selectQuery)).
					WithArgs(uint64(1), "approved", 10, 0).
					WillReturnRows(rows)
			},
			param: &model.ListExpenseRequest{
				UserID:       uint64(1),
				UserRole:     "manager",
				View:         model.ExpenseViewPersonal,
				Status:       &status,
				AutoApproved: true,
				Limit:        10,
				Offset:       0,
			},
			wantRes: []entity.ExpenseWithUser{
				{
					Expense: entity.Expense{
						ID:          uint64(1),
						UserID:      uint64(1),
						Amount:      uint64(15000),
						Description: description,
						ReceiptURL:  nil,
						Status:      entity.ExpenseStatusApproved,
						CreatedAt:   now,
						ProcessedAt: nil,
					},
					User: entity.UserSimple{
						ID:    1,
						Email: "john@mail.com",
						Name:  "John Doe",
					},
				},
			},
			wantTotal: 1,
			wantErr:   nil,
		},
		{
			name: "success approval_queue",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				countQuery := `SELECT COUNT(*) FROM expenses AS e WHERE e.status = 'awaiting_approval' AND e.user_id != $1`
				selectQuery := `
		SELECT
			e.id AS expense_id, e.user_id AS expense_user_id, e.amount AS expense_amount, e.description AS expense_description,
			e.receipt_url AS expense_receipt_url, e.status AS expense_status, e.created_at AS expense_created_at, e.processed_at AS expense_processed_at,
			u.id AS user_id, u.email AS user_email, u.name AS user_name
		FROM expenses AS e
		JOIN users AS u ON e.user_id = u.id WHERE e.status = 'awaiting_approval' AND e.user_id != $1 ORDER BY e.created_at DESC LIMIT $2 OFFSET $3`

				m.ExpectQuery(regexp.QuoteMeta(countQuery)).
					WithArgs(uint64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))
				rows := pgxmock.NewRows([]string{
					"expense_id", "expense_user_id", "expense_amount", "expense_description",
					"expense_receipt_url", "expense_status", "expense_created_at", "expense_processed_at",
					"user_id", "user_email", "user_name",
				}).AddRow(
					uint64(1), uint64(1), uint64(15000), description,
					nil, entity.ExpenseStatusApproved, now, nil,
					uint64(1), "john@mail.com", "John Doe",
				)
				m.ExpectQuery(regexp.QuoteMeta(selectQuery)).
					WithArgs(uint64(1), 10, 0).
					WillReturnRows(rows)
			},
			param: &model.ListExpenseRequest{
				UserID:   uint64(1),
				UserRole: "manager",
				View:     model.ExpenseViewApprovalQueue,
				Limit:    10,
				Offset:   0,
			},
			wantRes: []entity.ExpenseWithUser{
				{
					Expense: entity.Expense{
						ID:          uint64(1),
						UserID:      uint64(1),
						Amount:      uint64(15000),
						Description: description,
						ReceiptURL:  nil,
						Status:      entity.ExpenseStatusApproved,
						CreatedAt:   now,
						ProcessedAt: nil,
					},
					User: entity.UserSimple{
						ID:    1,
						Email: "john@mail.com",
						Name:  "John Doe",
					},
				},
			},
			wantTotal: 1,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			res, total, err := s.repo.List(s.ctx, tt.param)

			s.Equal(tt.wantRes, res)
			s.Equal(tt.wantTotal, total)
			s.Equal(tt.wantErr, err)
		})
	}
}

func (s *ExpenseRepositorySuite) TestExpenseRepository_FindDetailByID() {
	now := time.Date(2025, 8, 13, 10, 0, 0, 0, time.UTC)
	description := "dummy description"

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

	tests := []struct {
		name     string
		mockFunc func(pgxmock.PgxPoolIface)
		paramID  uint64
		wantRes  *entity.ExpenseDetail
		wantErr  error
	}{
		{
			name: "error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(uint64(1)).
					WillReturnError(errors.New("something error"))
			},
			paramID: uint64(1),
			wantRes: nil,
			wantErr: errors.New("something error"),
		},
		{
			name: "not found",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(uint64(1)).
					WillReturnError(pgx.ErrNoRows)
			},
			paramID: uint64(1),
			wantRes: nil,
			wantErr: nil,
		},
		{
			name: "success with approval",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"expense_id", "expense_user_id", "expense_amount", "expense_description",
					"expense_receipt_url", "expense_status", "expense_created_at", "expense_processed_at",
					"user_id", "user_email", "user_name",
					"approval_id", "approver_id", "approver_email", "approver_name",
					"approval_status", "approval_notes", "approval_created_at",
				}).AddRow(
					uint64(1), uint64(1), uint64(15000), description,
					nil, entity.ExpenseStatusApproved, now, nil,
					uint64(1), "john@mail.com", "John Doe",
					uint64(1), uint64(2), "budi@mail.com", "Budi",
					entity.ApprovalStatusApproved, nil, now,
				)
				m.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(uint64(1)).
					WillReturnRows(rows)
			},
			paramID: uint64(1),
			wantRes: &entity.ExpenseDetail{
				Expense: entity.Expense{
					ID:          uint64(1),
					UserID:      uint64(1),
					Amount:      uint64(15000),
					Description: description,
					ReceiptURL:  nil,
					Status:      entity.ExpenseStatusApproved,
					CreatedAt:   now,
					ProcessedAt: nil,
				},
				User: entity.UserSimple{
					ID:    uint64(1),
					Email: "john@mail.com",
					Name:  "John Doe",
				},
				Approval: &entity.ApprovalDetail{
					ID:            uint64(1),
					ApproverID:    uint64(2),
					ApproverEmail: "budi@mail.com",
					ApproverName:  "Budi",
					Status:        entity.ApprovalStatusApproved,
					Notes:         nil,
					CreatedAt:     now,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			res, err := s.repo.FindDetailByID(s.ctx, tt.paramID)

			s.Equal(tt.wantRes, res)
			s.Equal(tt.wantErr, err)
		})
	}
}

func (s *ExpenseRepositorySuite) TestExpenseRepository_FindByID() {
	description := "dummy description"
	receiptUrl := "https://example.com/receipt.jpg"

	tests := []struct {
		name     string
		mockFunc func(pgxmock.PgxPoolIface)
		paramID  uint64
		wantRes  *entity.Expense
		wantErr  error
	}{
		{
			name: "error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, user_id, amount, description, receipt_url, status, created_at, processed_at FROM expenses WHERE id = $1 LIMIT 1`,
				)).
					WithArgs(uint64(1)).
					WillReturnError(errors.New("something error"))
			},
			paramID: uint64(1),
			wantRes: nil,
			wantErr: errors.New("something error"),
		},
		{
			name: "not found",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, user_id, amount, description, receipt_url, status, created_at, processed_at FROM expenses WHERE id = $1 LIMIT 1`,
				)).
					WithArgs(uint64(1)).
					WillReturnError(pgx.ErrNoRows)
			},
			paramID: uint64(1),
			wantRes: nil,
			wantErr: nil,
		},
		{
			name: "success",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "user_id", "amount", "description", "receipt_url", "status", "created_at", "processed_at"}).
					AddRow(uint64(1), uint64(1), uint64(15000), description, &receiptUrl, entity.ExpenseStatusApproved, s.now, nil)
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, user_id, amount, description, receipt_url, status, created_at, processed_at FROM expenses WHERE id = $1 LIMIT 1`,
				)).
					WithArgs(uint64(1)).
					WillReturnRows(rows)
			},
			paramID: uint64(1),
			wantRes: &entity.Expense{
				ID:          uint64(1),
				UserID:      uint64(1),
				Amount:      uint64(15000),
				Description: description,
				ReceiptURL:  &receiptUrl,
				Status:      entity.ExpenseStatusApproved,
				CreatedAt:   s.now,
				ProcessedAt: nil,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			res, err := s.repo.FindByID(s.ctx, tt.paramID)

			s.Equal(tt.wantRes, res)
			s.Equal(tt.wantErr, err)
		})
	}
}

func (s *ExpenseRepositorySuite) TestExpenseRepository_FindByIDWithLock() {
	description := "dummy description"
	receiptUrl := "https://example.com/receipt.jpg"

	tests := []struct {
		name     string
		mockFunc func(pgxmock.PgxPoolIface)
		paramID  uint64
		wantRes  *entity.Expense
		wantErr  error
	}{
		{
			name: "error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, user_id, amount, description, receipt_url, status, created_at, processed_at FROM expenses WHERE id = $1 FOR UPDATE`,
				)).
					WithArgs(uint64(1)).
					WillReturnError(errors.New("something error"))
			},
			paramID: uint64(1),
			wantRes: nil,
			wantErr: errors.New("something error"),
		},
		{
			name: "not found",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, user_id, amount, description, receipt_url, status, created_at, processed_at FROM expenses WHERE id = $1 FOR UPDATE`,
				)).
					WithArgs(uint64(1)).
					WillReturnError(pgx.ErrNoRows)
			},
			paramID: uint64(1),
			wantRes: nil,
			wantErr: nil,
		},
		{
			name: "success",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "user_id", "amount", "description", "receipt_url", "status", "created_at", "processed_at"}).
					AddRow(uint64(1), uint64(1), uint64(15000), description, &receiptUrl, entity.ExpenseStatusApproved, s.now, nil)
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, user_id, amount, description, receipt_url, status, created_at, processed_at FROM expenses WHERE id = $1 FOR UPDATE`,
				)).
					WithArgs(uint64(1)).
					WillReturnRows(rows)
			},
			paramID: uint64(1),
			wantRes: &entity.Expense{
				ID:          uint64(1),
				UserID:      uint64(1),
				Amount:      uint64(15000),
				Description: description,
				ReceiptURL:  &receiptUrl,
				Status:      entity.ExpenseStatusApproved,
				CreatedAt:   s.now,
				ProcessedAt: nil,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			res, err := s.repo.FindByIDWithLock(s.ctx, s.mock, tt.paramID)

			s.Equal(tt.wantRes, res)
			s.Equal(tt.wantErr, err)
		})
	}
}

func (s *ExpenseRepositorySuite) TestExpenseRepository_UpdateStatusByIDTx() {
	tests := []struct {
		name        string
		mockFunc    func(pgxmock.PgxPoolIface)
		paramID     uint64
		paramStatus entity.ExpenseStatus
		wantErr     error
	}{
		{
			name: "error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(regexp.QuoteMeta(`UPDATE expenses SET status = $1 WHERE id = $2`)).
					WithArgs(pgxmock.AnyArg(), uint64(1)).
					WillReturnError(errors.New("something error"))
			},
			paramID:     uint64(1),
			paramStatus: entity.ExpenseStatusApproved,
			wantErr:     errors.New("something error"),
		},
		{
			name: "success",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(regexp.QuoteMeta(`UPDATE expenses SET status = $1 WHERE id = $2`)).
					WithArgs(pgxmock.AnyArg(), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			paramID:     uint64(1),
			paramStatus: entity.ExpenseStatusApproved,
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			err := s.repo.UpdateStatusByIDTx(s.ctx, s.mock, tt.paramID, tt.paramStatus)
			s.Equal(tt.wantErr, err)
		})
	}
}

func (s *ExpenseRepositorySuite) TestExpenseRepository_UpdateStatusByID() {
	tests := []struct {
		name        string
		mockFunc    func(pgxmock.PgxPoolIface)
		paramID     uint64
		paramStatus entity.ExpenseStatus
		wantErr     error
	}{
		{
			name: "error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(regexp.QuoteMeta(`UPDATE expenses SET status = $1, processed_at = $2 WHERE id = $3`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), uint64(1)).
					WillReturnError(errors.New("something error"))
			},
			paramID:     uint64(1),
			paramStatus: entity.ExpenseStatusApproved,
			wantErr:     errors.New("something error"),
		},
		{
			name: "success",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(regexp.QuoteMeta(`UPDATE expenses SET status = $1, processed_at = $2 WHERE id = $3`)).
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), uint64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			paramID:     uint64(1),
			paramStatus: entity.ExpenseStatusApproved,
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			err := s.repo.UpdateStatusByID(s.ctx, tt.paramID, tt.paramStatus, time.Now())
			s.Equal(tt.wantErr, err)
		})
	}
}

func TestExpenseRepositorySuite(t *testing.T) {
	suite.Run(t, new(ExpenseRepositorySuite))
}
