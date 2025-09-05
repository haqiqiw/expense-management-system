package model

import "fmt"

type Event interface {
	GetID() string
}

type ExpenseApprovedEvent struct {
	ID             uint64 `json:"id"`
	UserID         uint64 `json:"user_id"`
	Amount         uint64 `json:"amount"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (u *ExpenseApprovedEvent) GetID() string {
	return fmt.Sprintf("expense-%d", u.ID)
}
