package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	internalHttp "expense-management-system/internal/delivery/http"
	"expense-management-system/internal/mocks"
	"expense-management-system/internal/model"
	"expense-management-system/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type ExpenseControllerSuite struct {
	suite.Suite
	log      *zap.Logger
	validate *validator.Validate
}

func (s *ExpenseControllerSuite) SetupTest() {
	s.log = zap.NewNop()
	s.validate = validator.New()
}

func (s *ExpenseControllerSuite) TestExpenseController_Create() {
	tests := []struct {
		name       string
		body       any
		mockFunc   func(a *mocks.ExpenseUsecase)
		wantStatus int
		wantRes    string
	}{
		{
			name:       "empty body",
			body:       nil,
			mockFunc:   func(a *mocks.ExpenseUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantRes: `{"errors":[{"code":2000,"message":"AmountIDR failed on the 'required' rule"},` +
				`{"code":2001,"message":"Description failed on the 'required' rule"}],"meta":{"http_status":400}}`,
		},
		{
			name: "error on validate body",
			body: map[string]interface{}{
				"amount_idr":  0,
				"description": "",
				"receipt_url": "",
			},
			mockFunc:   func(a *mocks.ExpenseUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantRes: `{"errors":[{"code":2000,"message":"AmountIDR failed on the 'required' rule"},` +
				`{"code":2001,"message":"Description failed on the 'required' rule"}],"meta":{"http_status":400}}`,
		},
		{
			name: "error on create",
			body: map[string]interface{}{
				"amount_idr":  10000,
				"description": "Supplies",
				"receipt_url": "https://example.com/receipt.jpg",
			},
			mockFunc: func(a *mocks.ExpenseUsecase) {
				a.On("Create", mock.Anything, mock.Anything).
					Return(nil, errors.New("something error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantRes:    `{"errors":[{"code":100,"message":"Internal server error"}],"meta":{"http_status":500}}`,
		},
		{
			name: "success",
			body: map[string]interface{}{
				"amount_idr":  10000,
				"description": "Supplies",
				"receipt_url": "https://example.com/receipt.jpg",
			},
			mockFunc: func(a *mocks.ExpenseUsecase) {
				now := time.Date(2025, 10, 27, 13, 7, 31, 000, time.UTC)
				receipt := "https://example.com/receipt.jpg"

				a.On("Create", mock.Anything, mock.Anything).Return(&model.ExpenseCreateResponse{
					ID:               1,
					AmountIDR:        10000,
					Description:      "Supplies",
					ReceiptURL:       &receipt,
					Status:           "approved",
					RequiresApproval: false,
					AutoApproved:     true,
					CreatedAt:        now.Format(time.RFC3339),
				}, nil)
			},
			wantStatus: http.StatusCreated,
			wantRes: `{"data":{"id":1,"amount_idr":10000,"description":"Supplies","receipt_url":"https://example.com/receipt.jpg",` +
				`"status":"approved","requires_approval":false,"auto_approved":true,"created_at":"2025-10-27T13:07:31Z"},"meta":{"http_status":201}}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			eu := mocks.NewExpenseUsecase(s.T())
			tt.mockFunc(eu)

			ec := internalHttp.NewExpenseController(s.log, s.validate, eu)

			app := test.NewApi(s.log)
			app.Use(test.NewAuthMiddleware(1, "manager"))
			app.POST("/api/expenses", ec.Create)

			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/expenses", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, req)

			s.Equal(tt.wantStatus, rec.Code)
			s.Equal(tt.wantRes, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func (s *ExpenseControllerSuite) TestExpenseController_List() {
	tests := []struct {
		name       string
		mockFunc   func(a *mocks.ExpenseUsecase)
		wantStatus int
		wantRes    string
	}{
		{
			name: "error on list",
			mockFunc: func(a *mocks.ExpenseUsecase) {
				a.On("List", mock.Anything, mock.Anything).
					Return([]model.ExpenseDetailResponse{}, 0, errors.New("something error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantRes:    `{"errors":[{"code":100,"message":"Internal server error"}],"meta":{"http_status":500}}`,
		},
		{
			name: "success",
			mockFunc: func(a *mocks.ExpenseUsecase) {
				now := time.Date(2025, 10, 27, 13, 7, 31, 000, time.UTC)
				description := "dummy description"
				receipt := "https://example.com/receipt.jpg"

				a.On("List", mock.Anything, mock.Anything).
					Return([]model.ExpenseWithUserResponse{
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
					}, 1, nil)
			},
			wantStatus: http.StatusOK,
			wantRes: `{"data":[{"id":1,"amount_idr":10000,"description":"dummy description","receipt_url":"https://example.com/receipt.jpg",` +
				`"status":"approved","requires_approval":false,"auto_approved":true,"created_at":"2025-10-27T13:07:31Z",` +
				`"user":{"id":1,"email":"john@mail.com","name":"John Doe"}}],"meta":{"limit":10,"offset":0,"total":1,"http_status":200}}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			eu := mocks.NewExpenseUsecase(s.T())
			tt.mockFunc(eu)

			ec := internalHttp.NewExpenseController(s.log, s.validate, eu)

			app := test.NewApi(s.log)
			app.Use(test.NewAuthMiddleware(1, "manager"))
			app.GET("/api/expenses", ec.List)

			req := httptest.NewRequest("GET", "/api/expenses", nil)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, req)

			s.Equal(tt.wantStatus, rec.Code)
			s.Equal(tt.wantRes, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func (s *ExpenseControllerSuite) TestExpenseController_Get() {
	tests := []struct {
		name       string
		mockFunc   func(a *mocks.ExpenseUsecase)
		wantStatus int
		wantRes    string
	}{
		{
			name: "error on get",
			mockFunc: func(a *mocks.ExpenseUsecase) {
				a.On("FindByID", mock.Anything, mock.Anything).
					Return(nil, errors.New("something error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantRes:    `{"errors":[{"code":100,"message":"Internal server error"}],"meta":{"http_status":500}}`,
		},
		{
			name: "success",
			mockFunc: func(a *mocks.ExpenseUsecase) {
				now := time.Date(2025, 10, 27, 13, 7, 31, 000, time.UTC)
				description := "dummy description"
				receipt := "https://example.com/receipt.jpg"
				notes := "dummy notes"

				a.On("FindByID", mock.Anything, mock.Anything).Return(&model.ExpenseDetailResponse{
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
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantRes: `{"data":{"id":1,"amount_idr":10000,"description":"dummy description","receipt_url":"https://example.com/receipt.jpg",` +
				`"status":"approved","requires_approval":false,"auto_approved":true,"created_at":"2025-10-27T13:07:31Z","processed_at":null,` +
				`"user":{"id":1,"email":"john@mail.com","name":"John Doe"},"approval":{"id":1,"approver_id":1,"approver_email":"john@mail.com",` +
				`"approver_name":"John Doe","status":"approved","notes":"dummy notes","created_at":"2025-10-27T13:07:31Z"}},"meta":{"http_status":200}}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			eu := mocks.NewExpenseUsecase(s.T())
			tt.mockFunc(eu)

			ec := internalHttp.NewExpenseController(s.log, s.validate, eu)

			app := test.NewApi(s.log)
			app.Use(test.NewAuthMiddleware(1, "manager"))
			app.GET("/api/expenses/:id", ec.Get)

			req := httptest.NewRequest("GET", "/api/expenses/1", nil)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, req)

			s.Equal(tt.wantStatus, rec.Code)
			s.Equal(tt.wantRes, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestExpenseControllerSuite(t *testing.T) {
	suite.Run(t, new(ExpenseControllerSuite))
}
