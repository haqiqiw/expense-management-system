package usecase_test

import (
	"context"
	"errors"
	"expense-management-system/internal/db"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/messaging"
	"expense-management-system/internal/mocks"
	"expense-management-system/internal/model"
	"expense-management-system/internal/usecase"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type AuMockFunc func(
	db pgxmock.PgxPoolIface,
	tx db.Transactioner,
	ar *mocks.ApprovalRepository,
	er *mocks.ExpenseRepository,
	p *mocks.Producer[*model.ExpenseApprovedEvent],
)

type ApprovalUsecaseSuite struct {
	suite.Suite
	log *zap.Logger
	ctx context.Context
}

func (s *ApprovalUsecaseSuite) SetupTest() {
	s.log = zap.NewNop()
	s.ctx = context.Background()
}

func (s *ApprovalUsecaseSuite) TestApprovalUsecase_Approve() {
	notes := "dummy notes"

	tests := []struct {
		name       string
		request    *model.ApprovalExpenseRequest
		mockFunc   AuMockFunc
		wantErrMsg string
	}{
		{
			name: "invalid role",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "employee",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
			},
			wantErrMsg: "Forbidden",
		},
		{
			name: "error on find expense with lock",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(nil, errors.New("something error"))
				db.ExpectRollback()
			},
			wantErrMsg: "failed to find expense by id (1) with lock = something error",
		},
		{
			name: "error on expense not found",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(nil, nil)
				db.ExpectRollback()
			},
			wantErrMsg: "Expense not found",
		},
		{
			name: "error on same user id",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(&entity.Expense{UserID: 1}, nil)
				db.ExpectRollback()
			},
			wantErrMsg: "Forbidden",
		},
		{
			name: "error on expense already processed",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(&entity.Expense{UserID: 2, Amount: 1500000, Status: entity.ExpenseStatusApproved}, nil)
				db.ExpectRollback()
			},
			wantErrMsg: "Expense already processed",
		},
		{
			name: "error on create approval",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(&entity.Expense{UserID: 2, Amount: 1500000, Status: entity.ExpenseStatusAwaitingApproval}, nil)
				ar.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("something error"))
				db.ExpectRollback()
			},
			wantErrMsg: "failed to to create approval for expense id (1) = something error",
		},
		{
			name: "error on update expense",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(&entity.Expense{UserID: 2, Amount: 1500000, Status: entity.ExpenseStatusAwaitingApproval}, nil)
				ar.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
				er.On("UpdateStatusByIDTx", mock.Anything, mock.Anything, mock.Anything, entity.ExpenseStatusApproved).
					Return(errors.New("something error"))
				db.ExpectRollback()
			},
			wantErrMsg: "failed to to update expense for id (1) = something error",
		},
		{
			name: "success but error send event",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(&entity.Expense{UserID: 2, Amount: 1500000, Status: entity.ExpenseStatusAwaitingApproval}, nil)
				ar.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
				er.On("UpdateStatusByIDTx", mock.Anything, mock.Anything, mock.Anything, entity.ExpenseStatusApproved).
					Return(nil)
				db.ExpectCommit()
				p.On("Send", mock.Anything).Return(errors.New("something error"))
			},
			wantErrMsg: "",
		},
		{
			name: "success",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(&entity.Expense{UserID: 2, Amount: 1500000, Status: entity.ExpenseStatusAwaitingApproval}, nil)
				ar.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
				er.On("UpdateStatusByIDTx", mock.Anything, mock.Anything, mock.Anything, entity.ExpenseStatusApproved).
					Return(nil)
				db.ExpectCommit()
				p.On("Send", mock.Anything).Return(nil)
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			dbMock, _ := pgxmock.NewPool()
			defer dbMock.Close()
			tx := db.NewTransactioner(dbMock)

			ar := mocks.NewApprovalRepository(s.T())
			er := mocks.NewExpenseRepository(s.T())

			p := mocks.NewProducer[*model.ExpenseApprovedEvent](s.T())
			eap := &messaging.ExpenseApprovedProducer{
				Producer: p,
			}

			usecase := usecase.NewApprovalUsecase(s.log, tx, ar, er, eap)
			tt.mockFunc(dbMock, tx, ar, er, p)

			err := usecase.Approve(s.ctx, tt.request)

			if tt.wantErrMsg != "" {
				s.Equal(tt.wantErrMsg, err.Error())
			} else {
				s.Nil(err)
			}
		})
	}
}

func (s *ApprovalUsecaseSuite) TestApprovalUsecase_Reject() {
	notes := "dummy notes"

	tests := []struct {
		name       string
		request    *model.ApprovalExpenseRequest
		mockFunc   AuMockFunc
		wantErrMsg string
	}{
		{
			name: "invalid role",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "employee",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
			},
			wantErrMsg: "Forbidden",
		},
		{
			name: "error on find expense with lock",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(nil, errors.New("something error"))
				db.ExpectRollback()
			},
			wantErrMsg: "failed to find expense by id (1) with lock = something error",
		},
		{
			name: "error on expense not found",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(nil, nil)
				db.ExpectRollback()
			},
			wantErrMsg: "Expense not found",
		},
		{
			name: "error on same user id",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(&entity.Expense{UserID: 1}, nil)
				db.ExpectRollback()
			},
			wantErrMsg: "Forbidden",
		},
		{
			name: "error on expense already processed",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(&entity.Expense{UserID: 2, Amount: 1500000, Status: entity.ExpenseStatusApproved}, nil)
				db.ExpectRollback()
			},
			wantErrMsg: "Expense already processed",
		},
		{
			name: "error on create approval",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(&entity.Expense{UserID: 2, Amount: 1500000, Status: entity.ExpenseStatusAwaitingApproval}, nil)
				ar.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("something error"))
				db.ExpectRollback()
			},
			wantErrMsg: "failed to to create approval for expense id (1) = something error",
		},
		{
			name: "error on update expense",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(&entity.Expense{UserID: 2, Amount: 1500000, Status: entity.ExpenseStatusAwaitingApproval}, nil)
				ar.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
				er.On("UpdateStatusByIDTx", mock.Anything, mock.Anything, mock.Anything, entity.ExpenseStatusRejected).
					Return(errors.New("something error"))
				db.ExpectRollback()
			},
			wantErrMsg: "failed to to update expense for id (1) = something error",
		},
		{
			name: "success",
			request: &model.ApprovalExpenseRequest{
				ID:       1,
				Notes:    &notes,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				tx db.Transactioner,
				ar *mocks.ApprovalRepository,
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				db.ExpectBegin()
				er.On("FindByIDWithLock", mock.Anything, mock.Anything, uint64(1)).
					Return(&entity.Expense{UserID: 2, Amount: 1500000, Status: entity.ExpenseStatusAwaitingApproval}, nil)
				ar.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
				er.On("UpdateStatusByIDTx", mock.Anything, mock.Anything, mock.Anything, entity.ExpenseStatusRejected).
					Return(nil)
				db.ExpectCommit()
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			dbMock, _ := pgxmock.NewPool()
			defer dbMock.Close()
			tx := db.NewTransactioner(dbMock)

			ar := mocks.NewApprovalRepository(s.T())
			er := mocks.NewExpenseRepository(s.T())

			p := mocks.NewProducer[*model.ExpenseApprovedEvent](s.T())
			eap := &messaging.ExpenseApprovedProducer{
				Producer: p,
			}

			usecase := usecase.NewApprovalUsecase(s.log, tx, ar, er, eap)
			tt.mockFunc(dbMock, tx, ar, er, p)

			err := usecase.Reject(s.ctx, tt.request)

			if tt.wantErrMsg != "" {
				s.Equal(tt.wantErrMsg, err.Error())
			} else {
				s.Nil(err)
			}
		})
	}
}

func TestApprovalUsecaseSuite(t *testing.T) {
	suite.Run(t, new(ApprovalUsecaseSuite))
}
