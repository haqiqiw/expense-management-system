package model

type PaymentPartnerRequest struct {
	Amount     uint64 `json:"amount"`
	ExternalID string `json:"external_id"`
}

type PaymentPartnerResponse struct {
	PartnerID string `json:"partner_id"`
}

type PaymentProcessorRequest struct {
	ID             uint64 `json:"id"`
	UserID         uint64 `json:"user_id"`
	Amount         uint64 `json:"amount"`
	IdempotencyKey string `json:"idempotency_key"`
}
