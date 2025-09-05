package entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ExpenseStatus string

const (
	MinExpenseAmount        = 10_000
	MaxExpenseAmount        = 50_000_000
	ApprovalThresholdAmount = 1_000_000

	ExpenseStatusAwaitingApproval ExpenseStatus = "awaiting_approval"
	ExpenseStatusApproved         ExpenseStatus = "approved"
	ExpenseStatusRejected         ExpenseStatus = "rejected"
	ExpenseStatusCompleted        ExpenseStatus = "completed"

	keyPrefix = "EXP-"
)

type Expense struct {
	ID          uint64        `db:"id"`
	UserID      uint64        `db:"user_id"`
	Amount      uint64        `db:"amount"`
	Description string        `db:"description"`
	ReceiptURL  *string       `db:"receipt_url"`
	Status      ExpenseStatus `db:"status"`
	CreatedAt   time.Time     `db:"created_at"`
	ProcessedAt *time.Time    `db:"processed_at"`
}

func (e *Expense) RequiresApproval() bool {
	if e != nil {
		return e.Amount >= ApprovalThresholdAmount
	}

	return false
}

func (e *Expense) AutoApproved() bool {
	if e != nil {
		return e.Amount < ApprovalThresholdAmount
	}

	return false
}

func (e *Expense) GetKey() string {
	if e != nil {
		return fmt.Sprintf("%s%09s", keyPrefix, strings.ToUpper(strconv.FormatUint(e.ID, 36)))
	}

	return ""
}

func ParseExpenseStatus(str string) (ExpenseStatus, error) {
	switch str {
	case "awaiting_approval":
		return ExpenseStatusAwaitingApproval, nil
	case "approved":
		return ExpenseStatusApproved, nil
	case "rejected":
		return ExpenseStatusRejected, nil
	case "completed":
		return ExpenseStatusCompleted, nil
	default:
		return "", fmt.Errorf("invalid status: %s", str)
	}
}

type ExpenseWithUser struct {
	Expense
	User UserSimple
}

type ExpenseDetail struct {
	Expense
	User     UserSimple
	Approval *ApprovalDetail
}
