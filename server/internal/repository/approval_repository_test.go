package repository_test

import (
	"context"
	"errors"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/repository"
	"regexp"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/suite"
)

type ApprovalRepositorySuite struct {
	suite.Suite
	mock pgxmock.PgxPoolIface
	repo *repository.ApprovalRepository
	ctx  context.Context
	now  time.Time
}

func (s *ApprovalRepositorySuite) SetupTest() {
	s.mock, _ = pgxmock.NewPool()
	s.repo = repository.NewApprovalRepository(s.mock)
	s.ctx = context.Background()
	s.now = time.Now()
}

func (s *ApprovalRepositorySuite) TearDownTest() {
	s.mock.Close()
}

func (s *ApprovalRepositorySuite) TestApprovalRepository_CreateTx() {
	notes := "dummy notes"

	tests := []struct {
		name     string
		mockFunc func(pgxmock.PgxPoolIface)
		param    *entity.Approval
		wantErr  error
	}{
		{
			name: "error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO approvals (expense_id, approver_id, status, notes, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
				)).
					WithArgs(uint64(1), uint64(1), pgxmock.AnyArg(), &notes, pgxmock.AnyArg()).
					WillReturnError(errors.New("something error"))
			},
			param: &entity.Approval{
				ExpenseID:  uint64(1),
				ApproverID: uint64(1),
				Status:     entity.ApprovalStatusApproved,
				Notes:      &notes,
			},
			wantErr: errors.New("something error"),
		},
		{
			name: "success",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO approvals (expense_id, approver_id, status, notes, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
				)).
					WithArgs(uint64(1), uint64(1), pgxmock.AnyArg(), &notes, pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(uint64(1)))
			},
			param: &entity.Approval{
				ExpenseID:  uint64(1),
				ApproverID: uint64(1),
				Status:     entity.ApprovalStatusApproved,
				Notes:      &notes,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			err := s.repo.CreateTx(s.ctx, s.mock, tt.param)
			s.Equal(tt.wantErr, err)
		})
	}
}

func TestApprovalRepositorySuite(t *testing.T) {
	suite.Run(t, new(ApprovalRepositorySuite))
}
