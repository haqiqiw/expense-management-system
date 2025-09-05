package entity

import (
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID         uint64        `db:"id"`
	ExpenseID  uint64        `db:"expense_id"`
	ExternalID string        `db:"external_id"`
	PartnerID  *string       `db:"partner_id"`
	Amount     uint64        `db:"amount"`
	Status     PaymentStatus `db:"status"`
	Notes      *string       `db:"notes"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  time.Time     `db:"created_at"`
}
