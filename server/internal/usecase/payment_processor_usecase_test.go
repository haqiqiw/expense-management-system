package usecase_test

import (
	"context"
	"errors"
	"expense-management-system/internal/db"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/mocks"
	"expense-management-system/internal/model"
	"expense-management-system/internal/usecase"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type PpMockFunc func(
	db pgxmock.PgxPoolIface,
	rc *mocks.RedisClient,
	tx db.Transactioner,
	er *mocks.ExpenseRepository,
	ppr *mocks.PaymentPartnerRepository,
)

type PaymentProcessorUsecaseSuite struct {
	suite.Suite
	log *zap.Logger
	ctx context.Context
}

func (s *PaymentProcessorUsecaseSuite) SetupTest() {
	s.log = zap.NewNop()
	s.ctx = context.Background()
}

func (s *PaymentProcessorUsecaseSuite) TestPaymentProcessorUsecase_Execute() {
	tests := []struct {
		name       string
		request    *model.PaymentProcessorRequest
		mockFunc   PpMockFunc
		wantErrMsg string
	}{
		{
			name: "error on find expense",
			request: &model.PaymentProcessorRequest{
				ID:             1,
				UserID:         1,
				Amount:         17000,
				IdempotencyKey: "EXP-000ABC123",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				rc *mocks.RedisClient,
				tx db.Transactioner,
				er *mocks.ExpenseRepository,
				ppr *mocks.PaymentPartnerRepository,
			) {
				er.On("FindByID", mock.Anything, uint64(1)).
					Return(nil, errors.New("something error"))
			},
			wantErrMsg: "failed to find expense by id (1) = something error",
		},
		{
			name: "success with expense not found",
			request: &model.PaymentProcessorRequest{
				ID:             1,
				UserID:         1,
				Amount:         17000,
				IdempotencyKey: "EXP-000ABC123",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				rc *mocks.RedisClient,
				tx db.Transactioner,
				er *mocks.ExpenseRepository,
				ppr *mocks.PaymentPartnerRepository,
			) {
				er.On("FindByID", mock.Anything, uint64(1)).
					Return(nil, nil)
			},
			wantErrMsg: "",
		},
		{
			name: "success with invalid expense status",
			request: &model.PaymentProcessorRequest{
				ID:             1,
				UserID:         1,
				Amount:         17000,
				IdempotencyKey: "EXP-000ABC123",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				rc *mocks.RedisClient,
				tx db.Transactioner,
				er *mocks.ExpenseRepository,
				ppr *mocks.PaymentPartnerRepository,
			) {
				er.On("FindByID", mock.Anything, uint64(1)).
					Return(&entity.Expense{ID: 1, Status: entity.ExpenseStatusCompleted}, nil)
			},
			wantErrMsg: "",
		},
		{
			name: "error on acquire lock",
			request: &model.PaymentProcessorRequest{
				ID:             1,
				UserID:         1,
				Amount:         17000,
				IdempotencyKey: "EXP-000ABC123",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				rc *mocks.RedisClient,
				tx db.Transactioner,
				er *mocks.ExpenseRepository,
				ppr *mocks.PaymentPartnerRepository,
			) {
				er.On("FindByID", mock.Anything, uint64(1)).
					Return(&entity.Expense{ID: 1, Status: entity.ExpenseStatusApproved}, nil)

				boolCmd := redis.NewBoolCmd(context.Background())
				boolCmd.SetErr(errors.New("something error"))
				rc.On("SetNX", mock.Anything, "expense-payment:lock:1", "lock", mock.Anything).
					Return(boolCmd)
			},
			wantErrMsg: "failed to set lock for expense id (1) = something error",
		},
		{
			name: "success with lock already acquired",
			request: &model.PaymentProcessorRequest{
				ID:             1,
				UserID:         1,
				Amount:         17000,
				IdempotencyKey: "EXP-000ABC123",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				rc *mocks.RedisClient,
				tx db.Transactioner,
				er *mocks.ExpenseRepository,
				ppr *mocks.PaymentPartnerRepository,
			) {
				er.On("FindByID", mock.Anything, uint64(1)).
					Return(&entity.Expense{ID: 1, Status: entity.ExpenseStatusApproved}, nil)

				boolCmd := redis.NewBoolCmd(context.Background())
				boolCmd.SetVal(false)
				rc.On("SetNX", mock.Anything, "expense-payment:lock:1", "lock", mock.Anything).
					Return(boolCmd)
			},
			wantErrMsg: "",
		},
		{
			name: "error on update expense",
			request: &model.PaymentProcessorRequest{
				ID:             1,
				UserID:         1,
				Amount:         17000,
				IdempotencyKey: "EXP-000ABC123",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				rc *mocks.RedisClient,
				tx db.Transactioner,
				er *mocks.ExpenseRepository,
				ppr *mocks.PaymentPartnerRepository,
			) {
				er.On("FindByID", mock.Anything, uint64(1)).
					Return(&entity.Expense{ID: 1, Status: entity.ExpenseStatusApproved}, nil)

				boolCmd := redis.NewBoolCmd(context.Background())
				boolCmd.SetVal(true)
				rc.On("SetNX", mock.Anything, "expense-payment:lock:1", "lock", mock.Anything).
					Return(boolCmd)

				db.ExpectBegin()
				er.On("CompleteByIDTx", mock.Anything, mock.Anything, uint64(1), mock.Anything).
					Return(errors.New("something error"))
				db.ExpectRollback()

				intCmd := redis.NewIntCmd(context.Background())
				rc.On("Del", mock.Anything, "expense-payment:lock:1").Return(intCmd)
			},
			wantErrMsg: "failed to update expense for id (1) = something error",
		},
		{
			name: "error on payment partner",
			request: &model.PaymentProcessorRequest{
				ID:             1,
				UserID:         1,
				Amount:         17000,
				IdempotencyKey: "EXP-000ABC123",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				rc *mocks.RedisClient,
				tx db.Transactioner,
				er *mocks.ExpenseRepository,
				ppr *mocks.PaymentPartnerRepository,
			) {
				er.On("FindByID", mock.Anything, uint64(1)).
					Return(&entity.Expense{ID: 1, Status: entity.ExpenseStatusApproved}, nil)

				boolCmd := redis.NewBoolCmd(context.Background())
				boolCmd.SetVal(true)
				rc.On("SetNX", mock.Anything, "expense-payment:lock:1", "lock", mock.Anything).
					Return(boolCmd)

				db.ExpectBegin()
				er.On("CompleteByIDTx", mock.Anything, mock.Anything, uint64(1), mock.Anything).
					Return(nil)
				ppr.On("Execute", mock.Anything, mock.Anything).
					Return(nil, errors.New("something error"))
				db.ExpectRollback()

				intCmd := redis.NewIntCmd(context.Background())
				rc.On("Del", mock.Anything, "expense-payment:lock:1").Return(intCmd)
			},
			wantErrMsg: "failed to call partner for expense id (1) = something error",
		},
		{
			name: "success",
			request: &model.PaymentProcessorRequest{
				ID:             1,
				UserID:         1,
				Amount:         17000,
				IdempotencyKey: "EXP-000ABC123",
			},
			mockFunc: func(
				db pgxmock.PgxPoolIface,
				rc *mocks.RedisClient,
				tx db.Transactioner,
				er *mocks.ExpenseRepository,
				ppr *mocks.PaymentPartnerRepository,
			) {
				er.On("FindByID", mock.Anything, uint64(1)).
					Return(&entity.Expense{ID: 1, Status: entity.ExpenseStatusApproved}, nil)

				boolCmd := redis.NewBoolCmd(context.Background())
				boolCmd.SetVal(true)
				rc.On("SetNX", mock.Anything, "expense-payment:lock:1", "lock", mock.Anything).
					Return(boolCmd)

				db.ExpectBegin()
				er.On("CompleteByIDTx", mock.Anything, mock.Anything, uint64(1), mock.Anything).
					Return(nil)
				ppr.On("Execute", mock.Anything, mock.Anything).
					Return(&model.PaymentPartnerResponse{
						PartnerID: "sample-id",
					}, nil)
				db.ExpectCommit()

				intCmd := redis.NewIntCmd(context.Background())
				rc.On("Del", mock.Anything, "expense-payment:lock:1").Return(intCmd)
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			dbMock, _ := pgxmock.NewPool()
			defer dbMock.Close()
			tx := db.NewTransactioner(dbMock)

			rc := mocks.NewRedisClient(s.T())
			er := mocks.NewExpenseRepository(s.T())
			ppr := mocks.NewPaymentPartnerRepository(s.T())

			usecase := usecase.NewPaymentProcessorUsecase(s.log, rc, tx, er, ppr, 1)
			tt.mockFunc(dbMock, rc, tx, er, ppr)

			err := usecase.Execute(s.ctx, tt.request)

			if tt.wantErrMsg != "" {
				s.Equal(tt.wantErrMsg, err.Error())
			} else {
				s.Nil(err)
			}
		})
	}
}

func TestPaymentProcessorUsecaseSuite(t *testing.T) {
	suite.Run(t, new(PaymentProcessorUsecaseSuite))
}
