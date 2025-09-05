package repository_test

import (
	"context"
	"encoding/json"
	"errors"
	"expense-management-system/internal/httpclient"
	"expense-management-system/internal/mocks"
	"expense-management-system/internal/model"
	"expense-management-system/internal/repository"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPaymentPartnerRepository_Execute(t *testing.T) {
	tests := []struct {
		name       string
		param      *model.PaymentPartnerRequest
		mockFunc   func(*mocks.APIClient)
		wantRes    *model.PaymentPartnerResponse
		wantErrMsg string
	}{
		{
			name: "error on post",
			param: &model.PaymentPartnerRequest{
				Amount:     15000,
				ExternalID: "EXP-000123ABC",
			},
			mockFunc: func(a *mocks.APIClient) {
				a.On("Post", mock.Anything, "/v1/payments", mock.Anything).
					Return(nil, errors.New("something error"))
			},
			wantRes:    nil,
			wantErrMsg: "something error",
		},
		{
			name: "internal server error",
			param: &model.PaymentPartnerRequest{
				Amount:     15000,
				ExternalID: "EXP-000123ABC",
			},
			mockFunc: func(a *mocks.APIClient) {
				a.On("Post", mock.Anything, "/v1/payments", mock.Anything).
					Return(&httpclient.APIResponse{StatusCode: http.StatusInternalServerError}, nil)
			},
			wantRes:    nil,
			wantErrMsg: "payment partner error with status code = 500",
		},
		{
			name: "error on unmarshall",
			param: &model.PaymentPartnerRequest{
				Amount:     15000,
				ExternalID: "EXP-000123ABC",
			},
			mockFunc: func(a *mocks.APIClient) {
				a.On("Post", mock.Anything, "/v1/payments", mock.Anything).
					Return(&httpclient.APIResponse{
						StatusCode: http.StatusOK,
						Body:       []byte(`invalid-json`),
					}, nil)
			},
			wantRes:    nil,
			wantErrMsg: "payment partner error parse = invalid character 'i' looking for beginning of value",
		},
		{
			name: "sucess but external id exists",
			param: &model.PaymentPartnerRequest{
				Amount:     15000,
				ExternalID: "EXP-000123ABC",
			},
			mockFunc: func(a *mocks.APIClient) {
				body, _ := json.Marshal(map[string]interface{}{
					"data": map[string]interface{}{
						"id":          "partner-123",
						"external_id": "EXP-000123ABC",
						"status":      "success",
					},
					"message": "external id already exists",
				})
				a.On("Post", mock.Anything, "/v1/payments", mock.Anything).
					Return(&httpclient.APIResponse{
						StatusCode: http.StatusBadRequest,
						Body:       body,
					}, nil)
			},
			wantRes: &model.PaymentPartnerResponse{
				PartnerID: "partner-123",
			},
			wantErrMsg: "",
		},
		{
			name: "sucess",
			param: &model.PaymentPartnerRequest{
				Amount:     15000,
				ExternalID: "EXP-000123ABC",
			},
			mockFunc: func(a *mocks.APIClient) {
				body, _ := json.Marshal(map[string]interface{}{
					"data": map[string]interface{}{
						"id":          "partner-123",
						"external_id": "EXP-000123ABC",
						"status":      "success",
					},
				})
				a.On("Post", mock.Anything, "/v1/payments", mock.Anything).
					Return(&httpclient.APIResponse{
						StatusCode: http.StatusOK,
						Body:       body,
					}, nil)
			},
			wantRes: &model.PaymentPartnerResponse{
				PartnerID: "partner-123",
			},
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := mocks.NewAPIClient(t)
			repo := repository.NewPaymentPartnerRepository(a)
			tt.mockFunc(a)

			res, err := repo.Execute(context.Background(), tt.param)

			if tt.wantErrMsg != "" {
				assert.Nil(t, res)
				assert.Equal(t, tt.wantErrMsg, err.Error())
			} else {
				assert.Equal(t, tt.wantRes, res)
				assert.Nil(t, err)
			}
		})
	}
}
