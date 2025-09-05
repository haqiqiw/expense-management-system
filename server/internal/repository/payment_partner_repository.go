package repository

import (
	"context"
	"encoding/json"
	"expense-management-system/internal/httpclient"
	"expense-management-system/internal/model"
	"fmt"
	"net/http"
)

const (
	paymentURL              = "/v1/payments"
	externalIDExistsMessage = "external id already exists"
)

type PaymentPartnerRepository struct {
	client httpclient.APIClient
}

func NewPaymentPartnerRepository(client httpclient.APIClient) *PaymentPartnerRepository {
	return &PaymentPartnerRepository{
		client: client,
	}
}

func (r *PaymentPartnerRepository) Execute(ctx context.Context, req *model.PaymentPartnerRequest) (*model.PaymentPartnerResponse, error) {
	reqBody := map[string]interface{}{
		"amount":      req.Amount,
		"external_id": req.ExternalID,
	}

	resp, err := r.client.Post(ctx, paymentURL, reqBody)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("payment partner error with status code = %d", resp.StatusCode)
	}

	var parsedRes struct {
		Data struct {
			ID         string `json:"id"`
			ExternalID string `json:"external_id"`
			Status     string `json:"status"`
		} `json:"data"`
		Message string `json:"message,omitempty"`
	}

	err = json.Unmarshal(resp.Body, &parsedRes)
	if err != nil {
		return nil, fmt.Errorf("payment partner error parse = %w", err)
	}

	if resp.StatusCode == http.StatusOK ||
		(resp.StatusCode == http.StatusBadRequest && parsedRes.Message == externalIDExistsMessage) {
		return &model.PaymentPartnerResponse{
			PartnerID: parsedRes.Data.ID,
		}, nil
	}

	return nil, fmt.Errorf("payment partner error with status code = %d", resp.StatusCode)
}
