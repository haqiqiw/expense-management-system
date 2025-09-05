package usecase_test

import (
	"context"
	"errors"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/messaging"
	"expense-management-system/internal/mocks"
	"expense-management-system/internal/model"
	"expense-management-system/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type ExpenseUsecaseSuite struct {
	suite.Suite
	log *zap.Logger
	ctx context.Context
}

func (s *ExpenseUsecaseSuite) SetupTest() {
	s.log = zap.NewNop()
	s.ctx = context.Background()
}

func (s *ExpenseUsecaseSuite) TestExpenseUsecase_Create() {
	receiptUrl := "https://example.com/receipt.jpg"

	tests := []struct {
		name     string
		request  *model.CreateExpenseRequest
		mockFunc func(
			er *mocks.ExpenseRepository,
			p *mocks.Producer[*model.ExpenseApprovedEvent],
		)
		wantErrMsg string
	}{
		{
			name: "error on min amount",
			request: &model.CreateExpenseRequest{
				UserID:      1,
				AmountIDR:   500,
				Description: "dummy description",
				ReceiptURL:  &receiptUrl,
			},
			mockFunc: func(
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
			},
			wantErrMsg: "Amount can't be less than Rp 10.000",
		},
		{
			name: "error on max amount",
			request: &model.CreateExpenseRequest{
				UserID:      1,
				AmountIDR:   250000000,
				Description: "dummy description",
				ReceiptURL:  &receiptUrl,
			},
			mockFunc: func(
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
			},
			wantErrMsg: "Amount can't be greater than Rp 50.000.000",
		},
		{
			name: "error on create",
			request: &model.CreateExpenseRequest{
				UserID:      1,
				AmountIDR:   15500,
				Description: "dummy description",
				ReceiptURL:  &receiptUrl,
			},
			mockFunc: func(
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {

				er.On("Create", mock.Anything, mock.Anything).
					Return(errors.New("something error"))
			},
			wantErrMsg: "failed to create expense = something error",
		},
		{
			name: "success but error send event",
			request: &model.CreateExpenseRequest{
				UserID:      1,
				AmountIDR:   15500,
				Description: "dummy description",
				ReceiptURL:  &receiptUrl,
			},
			mockFunc: func(
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				er.On("Create", mock.Anything, mock.Anything).Return(nil)
				p.On("Send", mock.Anything).Return(errors.New("something error"))
			},
			wantErrMsg: "",
		},
		{
			name: "success",
			request: &model.CreateExpenseRequest{
				UserID:      1,
				AmountIDR:   15500,
				Description: "dummy description",
				ReceiptURL:  &receiptUrl,
			},
			mockFunc: func(
				er *mocks.ExpenseRepository,
				p *mocks.Producer[*model.ExpenseApprovedEvent],
			) {
				er.On("Create", mock.Anything, mock.Anything).Return(nil)
				p.On("Send", mock.Anything).Return(nil)
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			er := mocks.NewExpenseRepository(s.T())
			p := mocks.NewProducer[*model.ExpenseApprovedEvent](s.T())
			eap := &messaging.ExpenseApprovedProducer{
				Producer: p,
			}

			usecase := usecase.NewExpenseUsecase(s.log, er, eap)
			tt.mockFunc(er, p)

			_, err := usecase.Create(s.ctx, tt.request)

			if tt.wantErrMsg != "" {
				s.Equal(tt.wantErrMsg, err.Error())
			} else {
				s.Nil(err)
			}
		})
	}
}

func (s *ExpenseUsecaseSuite) TestExpenseUsecase_List() {
	now := time.Date(2025, 8, 13, 10, 0, 0, 0, time.UTC)
	description := "dummy description"
	receipt := "https://example.com/receipt.jpg"

	tests := []struct {
		name       string
		request    *model.ListExpenseRequest
		mockFunc   func(er *mocks.ExpenseRepository)
		wantRes    []model.ExpenseWithUserResponse
		wantTotal  int
		wantErrMsg string
	}{
		{
			name: "error on list",
			request: &model.ListExpenseRequest{
				Offset: 0,
				Limit:  10,
			},
			mockFunc: func(er *mocks.ExpenseRepository) {
				er.On("List", mock.Anything, mock.Anything).
					Return(nil, 0, errors.New("something error"))
			},
			wantRes:    []model.ExpenseWithUserResponse{},
			wantTotal:  0,
			wantErrMsg: "failed to get expenses = something error",
		},
		{
			name: "success but empty",
			request: &model.ListExpenseRequest{
				Offset: 0,
				Limit:  10,
			},
			mockFunc: func(er *mocks.ExpenseRepository) {
				er.On("List", mock.Anything, mock.Anything).
					Return([]entity.ExpenseWithUser{}, 0, nil)
			},
			wantRes:    []model.ExpenseWithUserResponse{},
			wantTotal:  0,
			wantErrMsg: "",
		},
		{
			name: "success",
			request: &model.ListExpenseRequest{
				Offset: 0,
				Limit:  10,
			},
			mockFunc: func(er *mocks.ExpenseRepository) {
				er.On("List", mock.Anything, mock.Anything).Return([]entity.ExpenseWithUser{
					{
						Expense: entity.Expense{
							ID:          1,
							UserID:      1,
							Amount:      10000,
							Description: description,
							ReceiptURL:  &receipt,
							Status:      entity.ExpenseStatusApproved,
							CreatedAt:   now,
						},
						User: entity.UserSimple{
							ID:    1,
							Email: "john@mail.com",
							Name:  "John Doe",
						},
					},
				}, 1, nil)
			},
			wantRes: []model.ExpenseWithUserResponse{
				{
					ID:               1,
					AmountIDR:        10000,
					Description:      description,
					ReceiptURL:       &receipt,
					Status:           "approved",
					RequiresApproval: false,
					AutoApproved:     true,
					CreatedAt:        now.Format(time.RFC3339),
					User: model.UserSimpleResponse{
						ID:    1,
						Email: "john@mail.com",
						Name:  "John Doe",
					},
				},
			},
			wantTotal:  1,
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			er := mocks.NewExpenseRepository(s.T())
			p := mocks.NewProducer[*model.ExpenseApprovedEvent](s.T())
			eap := &messaging.ExpenseApprovedProducer{
				Producer: p,
			}

			usecase := usecase.NewExpenseUsecase(s.log, er, eap)
			tt.mockFunc(er)

			res, total, err := usecase.List(s.ctx, tt.request)

			if tt.wantErrMsg != "" {
				s.Empty(res)
				s.Equal(tt.wantErrMsg, err.Error())
			} else {
				if len(res) > 0 {
					s.Equal(tt.wantRes[0], res[0])
				}
				s.Len(res, len(tt.wantRes))
				s.Nil(err)
			}
			s.Equal(tt.wantTotal, total)
		})
	}
}

func (s *ExpenseUsecaseSuite) TestExpenseUsecase_FindByID() {
	now := time.Date(2025, 8, 13, 10, 0, 0, 0, time.UTC)
	description := "dummy description"
	notes := "dummy notes"
	receipt := "https://example.com/receipt.jpg"

	tests := []struct {
		name       string
		request    *model.GetExpenseRequest
		mockFunc   func(er *mocks.ExpenseRepository)
		wantRes    *model.ExpenseDetailResponse
		wantErrMsg string
	}{
		{
			name: "error on list",
			request: &model.GetExpenseRequest{
				ID:       1,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(er *mocks.ExpenseRepository) {
				er.On("FindDetailByID", mock.Anything, uint64(1)).
					Return(nil, errors.New("something error"))
			},
			wantRes:    &model.ExpenseDetailResponse{},
			wantErrMsg: "failed to find expense by id (1) = something error",
		},
		{
			name: "error on expense not found",
			request: &model.GetExpenseRequest{
				ID:       1,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(er *mocks.ExpenseRepository) {
				er.On("FindDetailByID", mock.Anything, uint64(1)).
					Return(nil, nil)
			},
			wantRes:    nil,
			wantErrMsg: "Expense not found",
		},
		{
			name: "error on invalid access",
			request: &model.GetExpenseRequest{
				ID:       1,
				UserID:   1,
				UserRole: "employee",
			},
			mockFunc: func(er *mocks.ExpenseRepository) {
				er.On("FindDetailByID", mock.Anything, uint64(1)).
					Return(&entity.ExpenseDetail{
						Expense: entity.Expense{ID: 1, UserID: 2},
					}, nil)
			},
			wantRes:    nil,
			wantErrMsg: "Forbidden",
		},
		{
			name: "success",
			request: &model.GetExpenseRequest{
				ID:       1,
				UserID:   1,
				UserRole: "manager",
			},
			mockFunc: func(er *mocks.ExpenseRepository) {
				er.On("FindDetailByID", mock.Anything, uint64(1)).
					Return(&entity.ExpenseDetail{
						Expense: entity.Expense{
							ID:          1,
							UserID:      2,
							Amount:      10000,
							Description: description,
							ReceiptURL:  &receipt,
							Status:      entity.ExpenseStatusApproved,
							CreatedAt:   now,
						},
						User: entity.UserSimple{
							ID:    1,
							Email: "john@mail.com",
							Name:  "John Doe",
						},
						Approval: &entity.ApprovalDetail{
							ID:            1,
							ApproverID:    1,
							ApproverEmail: "john@mail.com",
							ApproverName:  "John Doe",
							Status:        entity.ApprovalStatusApproved,
							Notes:         &notes,
							CreatedAt:     now,
						},
					}, nil)
			},
			wantRes: &model.ExpenseDetailResponse{
				ID:               1,
				AmountIDR:        10000,
				Description:      description,
				ReceiptURL:       &receipt,
				Status:           "approved",
				RequiresApproval: false,
				AutoApproved:     true,
				CreatedAt:        now.Format(time.RFC3339),
				User: model.UserSimpleResponse{
					ID:    1,
					Email: "john@mail.com",
					Name:  "John Doe",
				},
				Approval: &model.ApprovalDetailResponse{
					ID:            1,
					ApproverID:    1,
					ApproverEmail: "john@mail.com",
					ApproverName:  "John Doe",
					Status:        "approved",
					Notes:         &notes,
					CreatedAt:     now.Format(time.RFC3339),
				},
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			er := mocks.NewExpenseRepository(s.T())
			p := mocks.NewProducer[*model.ExpenseApprovedEvent](s.T())
			eap := &messaging.ExpenseApprovedProducer{
				Producer: p,
			}

			usecase := usecase.NewExpenseUsecase(s.log, er, eap)
			tt.mockFunc(er)

			res, err := usecase.FindByID(s.ctx, tt.request)

			if tt.wantErrMsg != "" {
				s.Empty(res)
				s.Equal(tt.wantErrMsg, err.Error())
			} else {
				s.Equal(*tt.wantRes, *res)
				s.Nil(err)
			}
		})
	}
}

func TestExpenseUsecaseSuite(t *testing.T) {
	suite.Run(t, new(ExpenseUsecaseSuite))
}
