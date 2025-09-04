package model_test

import (
	"expense-management-system/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomError_Error(t *testing.T) {
	tests := []struct {
		name      string
		customErr *model.CustomError
		wantMsg   string
	}{
		{
			name:      "single error",
			customErr: model.ErrEmailAlreadyExist,
			wantMsg:   "email already exist",
		},
		{
			name: "multiple errors",
			customErr: func() *model.CustomError {
				err := model.ErrInvalidPassword
				err.Append(model.ErrorItem{
					Code:    9999,
					Message: "another error",
				})
				return err
			}(),
			wantMsg: "invalid password",
		},
		{
			name: "empty errors",
			customErr: &model.CustomError{
				HTTPStatus: http.StatusBadRequest,
				Errors:     []model.ErrorItem{},
			},
			wantMsg: "empty error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.customErr.Error()

			assert.Equal(t, tt.wantMsg, msg)
		})
	}
}
