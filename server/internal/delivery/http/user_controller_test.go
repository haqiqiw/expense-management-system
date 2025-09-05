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
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type UserControllerSuite struct {
	suite.Suite
	log      *zap.Logger
	validate *validator.Validate
}

func (s *UserControllerSuite) SetupTest() {
	s.log = zap.NewNop()
	s.validate = validator.New()
}

func (s *UserControllerSuite) TestUserController_Create() {
	tests := []struct {
		name       string
		body       any
		mockFunc   func(a *mocks.UserUsecase)
		wantStatus int
		wantRes    string
	}{
		{
			name:       "empty body",
			body:       nil,
			mockFunc:   func(a *mocks.UserUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantRes: `{"errors":[{"code":2000,"message":"Email failed on the 'required' rule"},` +
				`{"code":2001,"message":"Name failed on the 'required' rule"},` +
				`{"code":2002,"message":"Password failed on the 'required' rule"},` +
				`{"code":2003,"message":"Role failed on the 'required' rule"}],"meta":{"http_status":400}}`,
		},
		{
			name: "error on validate body",
			body: map[string]interface{}{
				"email":    "",
				"name":     "",
				"password": "",
				"role":     "",
			},
			mockFunc:   func(a *mocks.UserUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantRes: `{"errors":[{"code":2000,"message":"Email failed on the 'required' rule"},` +
				`{"code":2001,"message":"Name failed on the 'required' rule"},` +
				`{"code":2002,"message":"Password failed on the 'required' rule"},` +
				`{"code":2003,"message":"Role failed on the 'required' rule"}],"meta":{"http_status":400}}`,
		},
		{
			name: "error on create",
			body: map[string]interface{}{
				"email":    "john@mail.com",
				"name":     "John Doe",
				"password": "password",
				"role":     "manager",
			},
			mockFunc: func(a *mocks.UserUsecase) {
				a.On("Create", mock.Anything, mock.Anything).
					Return(nil, errors.New("something error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantRes:    `{"errors":[{"code":100,"message":"Internal server error"}],"meta":{"http_status":500}}`,
		},
		{
			name: "success",
			body: map[string]interface{}{
				"email":    "john@mail.com",
				"name":     "John Doe",
				"password": "password",
				"role":     "manager",
			},
			mockFunc: func(a *mocks.UserUsecase) {
				now := time.Date(2025, 10, 27, 13, 7, 31, 000, time.UTC)
				a.On("Create", mock.Anything, mock.Anything).Return(&model.UserResponse{
					ID:        1,
					Email:     "john@mail.com",
					Name:      "John Doe",
					Role:      "manager",
					CreatedAt: now.Format(time.RFC3339),
				}, nil)
			},
			wantStatus: http.StatusCreated,
			wantRes: `{"data":{"id":1,"email":"john@mail.com","name":"John Doe","role":"manager",` +
				`"created_at":"2025-10-27T13:07:31Z"},"meta":{"http_status":201}}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			uu := mocks.NewUserUsecase(s.T())
			tt.mockFunc(uu)

			uc := internalHttp.NewUserController(s.log, s.validate, uu)

			app := config.NewGin(s.log)
			app.POST("/api/users", uc.Register)

			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/users", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, req)

			s.Equal(tt.wantStatus, rec.Code)
			s.Equal(tt.wantRes, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func (s *UserControllerSuite) TestUserController_Get() {
	tests := []struct {
		name       string
		mockFunc   func(a *mocks.UserUsecase)
		wantStatus int
		wantRes    string
	}{
		{
			name: "error on get",
			mockFunc: func(a *mocks.UserUsecase) {
				a.On("FindByID", mock.Anything, mock.Anything).
					Return(nil, errors.New("something error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantRes:    `{"errors":[{"code":100,"message":"Internal server error"}],"meta":{"http_status":500}}`,
		},
		{
			name: "success",
			mockFunc: func(a *mocks.UserUsecase) {
				now := time.Date(2025, 10, 27, 13, 7, 31, 000, time.UTC)
				a.On("FindByID", mock.Anything, mock.Anything).Return(&model.UserResponse{
					ID:        1,
					Email:     "john@mail.com",
					Name:      "John Doe",
					Role:      "manager",
					CreatedAt: now.Format(time.RFC3339),
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantRes: `{"data":{"id":1,"email":"john@mail.com","name":"John Doe","role":"manager",` +
				`"created_at":"2025-10-27T13:07:31Z"},"meta":{"http_status":200}}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			uu := mocks.NewUserUsecase(s.T())
			tt.mockFunc(uu)

			uc := internalHttp.NewUserController(s.log, s.validate, uu)

			app := config.NewGin(s.log)
			app.Use(test.NewAuthMiddleware(1, "manager"))
			app.GET("/api/users/me", uc.Me)

			req := httptest.NewRequest("GET", "/api/users/me", nil)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, req)

			s.Equal(tt.wantStatus, rec.Code)
			s.Equal(tt.wantRes, strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerSuite))
}
