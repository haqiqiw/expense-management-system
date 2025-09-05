package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"expense-management-system/internal/config"
	internalHttp "expense-management-system/internal/delivery/http"
	"expense-management-system/internal/mocks"
	"expense-management-system/internal/model"
	"expense-management-system/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type ApprovalControllerSuite struct {
	suite.Suite
	log      *zap.Logger
	validate *validator.Validate
}

func (s *ApprovalControllerSuite) SetupTest() {
	s.log = zap.NewNop()
	s.validate = validator.New()
}

func (s *ApprovalControllerSuite) TestApprovalController_Approve() {
	tests := []struct {
		name       string
		body       any
		mockFunc   func(a *mocks.ApprovalUsecase)
		wantStatus int
		wantRes    string
	}{
		{
			name: "custom error on approve",
			body: nil,
			mockFunc: func(a *mocks.ApprovalUsecase) {
				a.On("Approve", mock.Anything, mock.Anything).
					Return(model.ErrExpenseAlreadyProcessed)
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantRes:    `{"errors":[{"code":1006,"message":"Expense already processed"}],"meta":{"http_status":422}}`,
		},
		{
			name: "unexpected error on approve",
			body: nil,
			mockFunc: func(a *mocks.ApprovalUsecase) {
				a.On("Approve", mock.Anything, mock.Anything).
					Return(errors.New("something error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantRes:    `{"errors":[{"code":100,"message":"Internal server error"}],"meta":{"http_status":500}}`,
		},
		{
			name: "success",
			body: nil,
			mockFunc: func(a *mocks.ApprovalUsecase) {
				a.On("Approve", mock.Anything, mock.Anything).Return(nil)
			},
			wantStatus: http.StatusOK,
			wantRes:    `{"message":"Expense approved","meta":{"http_status":200}}`,
		},
		{
			name: "success with notes",
			body: map[string]interface{}{
				"notes": "dummy notes",
			},
			mockFunc: func(a *mocks.ApprovalUsecase) {
				a.On("Approve", mock.Anything, mock.Anything).Return(nil)
			},
			wantStatus: http.StatusOK,
			wantRes:    `{"message":"Expense approved","meta":{"http_status":200}}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			au := mocks.NewApprovalUsecase(s.T())
			tt.mockFunc(au)

			ac := internalHttp.NewApprovalController(s.log, au)

			app := config.NewGin(s.log)
			app.Use(test.NewAuthMiddleware(1, "manager"))
			app.PUT("/expenses/:id/approve", ac.Approve)

			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("PUT", "/expenses/1/approve", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, req)

			s.Equal(tt.wantStatus, rec.Code)
			s.Equal(tt.wantRes, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func (s *ApprovalControllerSuite) TestApprovalController_Reject() {
	tests := []struct {
		name       string
		body       any
		mockFunc   func(a *mocks.ApprovalUsecase)
		wantStatus int
		wantRes    string
	}{
		{
			name: "custom error on reject",
			body: nil,
			mockFunc: func(a *mocks.ApprovalUsecase) {
				a.On("Reject", mock.Anything, mock.Anything).
					Return(model.ErrExpenseAlreadyProcessed)
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantRes:    `{"errors":[{"code":1006,"message":"Expense already processed"}],"meta":{"http_status":422}}`,
		},
		{
			name: "unexpected error on reject",
			body: nil,
			mockFunc: func(a *mocks.ApprovalUsecase) {
				a.On("Reject", mock.Anything, mock.Anything).
					Return(errors.New("something error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantRes:    `{"errors":[{"code":100,"message":"Internal server error"}],"meta":{"http_status":500}}`,
		},
		{
			name: "success",
			body: nil,
			mockFunc: func(a *mocks.ApprovalUsecase) {
				a.On("Reject", mock.Anything, mock.Anything).Return(nil)
			},
			wantStatus: http.StatusOK,
			wantRes:    `{"message":"Expense rejected","meta":{"http_status":200}}`,
		},
		{
			name: "success with notes",
			body: map[string]interface{}{
				"notes": "dummy notes",
			},
			mockFunc: func(a *mocks.ApprovalUsecase) {
				a.On("Reject", mock.Anything, mock.Anything).Return(nil)
			},
			wantStatus: http.StatusOK,
			wantRes:    `{"message":"Expense rejected","meta":{"http_status":200}}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			au := mocks.NewApprovalUsecase(s.T())
			tt.mockFunc(au)

			ac := internalHttp.NewApprovalController(s.log, au)

			app := config.NewGin(s.log)
			app.Use(test.NewAuthMiddleware(1, "manager"))
			app.PUT("/expenses/:id/reject", ac.Reject)

			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("PUT", "/expenses/1/reject", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, req)

			s.Equal(tt.wantStatus, rec.Code)
			s.Equal(tt.wantRes, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestApprovalControllerSuite(t *testing.T) {
	suite.Run(t, new(ApprovalControllerSuite))
}
