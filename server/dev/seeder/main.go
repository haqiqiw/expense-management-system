package main

import (
	"context"
	"expense-management-system/internal/config"
	"expense-management-system/internal/entity"
	"fmt"
	"log"
	"time"
)

// for seeding purposes, the password for each user is set to be
// the same as the prefix of their email (the part before '@')
// example: "john@mail.com" => password: "john"
func main() {
	ctx := context.Background()

	logger, err := config.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	env, err := config.NewEnv()
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to initialize env: %+v", err))
	}

	db, err := config.NewDatabase(ctx, env)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to initialize database: %+v", err))
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to start transaction: %+v", err))
	}

	logger.Info("starting database seeding ...")

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			logger.Fatal(fmt.Sprintf("database seeding failed: %+v", err))
		} else {
			tx.Commit(ctx)
			logger.Info("database seeding completed successfully!")
		}
	}()

	var (
		defaultReceiptURL = "https://placehold.co/500x700"

		users    []entity.User
		expenses []entity.Expense
	)

	// prepare data
	users = []entity.User{
		{ID: 1, Email: "john@mail.com", Name: "John", PasswordHash: "$2a$10$AVS.3jUpmFXCEgcN8J4Peerq4643.KOMYIhUk8pAgny06t/zxNqHm", Role: "manager", CreatedAt: time.Date(2025, 9, 1, 13, 2, 30, 000, time.UTC)},
		{ID: 2, Email: "wawan@mail.com", Name: "Wawan", PasswordHash: "$2a$10$KOIpI26eF/WTIE8s8GV56OCo6qK/GOHaIppLcNX8elxXnuUekOw82", Role: "manager", CreatedAt: time.Date(2025, 9, 1, 13, 3, 30, 000, time.UTC)},
		{ID: 3, Email: "budi@mail.com", Name: "Budi", PasswordHash: "$2a$10$ZnT7RBH14iLQXsQGN1Z99O7lASuaJYYiIZcyCpyj9oC8w6m7..Wxu", Role: "employee", CreatedAt: time.Date(2025, 9, 2, 13, 4, 30, 000, time.UTC)},
		{ID: 4, Email: "lala@mail.com", Name: "Lala", PasswordHash: "$2a$10$UFBubu4rYw7.ZvVd9rq75Otj12ppjaVOJO/VTBjyc0wkP.fhfBBsO", Role: "employee", CreatedAt: time.Date(2025, 9, 2, 13, 5, 30, 000, time.UTC)},
	}
	expenses = []entity.Expense{
		{ID: 1, UserID: 3, Amount: 150000, Description: "Snacks", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusCompleted, CreatedAt: time.Date(2025, 8, 2, 13, 2, 30, 000, time.UTC)},
		{ID: 2, UserID: 3, Amount: 15000, Description: "Transport KR", ReceiptURL: nil, Status: entity.ExpenseStatusApproved, CreatedAt: time.Date(2025, 8, 2, 13, 4, 30, 000, time.UTC)},
		{ID: 3, UserID: 3, Amount: 1750000, Description: "Transport train", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusCompleted, CreatedAt: time.Date(2025, 8, 2, 14, 5, 30, 000, time.UTC)},
		{ID: 4, UserID: 3, Amount: 1250000, Description: "Transport bus", ReceiptURL: nil, Status: entity.ExpenseStatusRejected, CreatedAt: time.Date(2025, 8, 2, 15, 1, 30, 000, time.UTC)},
		{ID: 5, UserID: 3, Amount: 7500000, Description: "Team dinner", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusApproved, CreatedAt: time.Date(2025, 8, 3, 16, 1, 30, 000, time.UTC)},
		{ID: 6, UserID: 3, Amount: 2500000, Description: "Transport plane", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusAwaitingApproval, CreatedAt: time.Date(2025, 8, 3, 17, 1, 30, 000, time.UTC)},

		{ID: 7, UserID: 4, Amount: 60000, Description: "Drinks", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusCompleted, CreatedAt: time.Date(2025, 8, 4, 1, 1, 30, 000, time.UTC)},
		{ID: 8, UserID: 4, Amount: 450000, Description: "Office supplies", ReceiptURL: nil, Status: entity.ExpenseStatusApproved, CreatedAt: time.Date(2025, 8, 4, 2, 2, 30, 000, time.UTC)},
		{ID: 9, UserID: 4, Amount: 1500000, Description: "Client dinner", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusAwaitingApproval, CreatedAt: time.Date(2025, 8, 4, 3, 1, 30, 000, time.UTC)},
		{ID: 10, UserID: 4, Amount: 2570000, Description: "Rent area", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusCompleted, CreatedAt: time.Date(2025, 8, 4, 4, 1, 30, 000, time.UTC)},

		{ID: 11, UserID: 1, Amount: 4100000, Description: "Meeting supplies", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusCompleted, CreatedAt: time.Date(2025, 8, 5, 5, 1, 30, 000, time.UTC)},
		{ID: 12, UserID: 1, Amount: 16500, Description: "Tranport ojol", ReceiptURL: nil, Status: entity.ExpenseStatusApproved, CreatedAt: time.Date(2025, 8, 5, 5, 2, 30, 000, time.UTC)},
		{ID: 13, UserID: 1, Amount: 2350000, Description: "Rent meeting room", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusAwaitingApproval, CreatedAt: time.Date(2025, 8, 5, 4, 1, 30, 000, time.UTC)},

		{ID: 14, UserID: 2, Amount: 4100000, Description: "Meeting supplies", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusCompleted, CreatedAt: time.Date(2025, 8, 5, 6, 1, 30, 000, time.UTC)},
		{ID: 15, UserID: 2, Amount: 16500, Description: "Tranport ojol", ReceiptURL: nil, Status: entity.ExpenseStatusApproved, CreatedAt: time.Date(2025, 8, 5, 7, 2, 30, 000, time.UTC)},
		{ID: 16, UserID: 2, Amount: 2350000, Description: "Rent meeting room", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusAwaitingApproval, CreatedAt: time.Date(2025, 8, 5, 8, 1, 30, 000, time.UTC)},

		{ID: 17, UserID: 3, Amount: 2350000, Description: "Transport plane", ReceiptURL: nil, Status: entity.ExpenseStatusRejected, CreatedAt: time.Date(2025, 8, 6, 1, 2, 30, 000, time.UTC)},
		{ID: 18, UserID: 3, Amount: 1200000, Description: "Rent meeting room", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusCompleted, CreatedAt: time.Date(2025, 8, 6, 2, 2, 30, 000, time.UTC)},
		{ID: 19, UserID: 3, Amount: 750000, Description: "Foods and drinks", ReceiptURL: nil, Status: entity.ExpenseStatusApproved, CreatedAt: time.Date(2025, 8, 6, 3, 2, 30, 000, time.UTC)},
		{ID: 20, UserID: 3, Amount: 1210000, Description: "Team dinner", ReceiptURL: &defaultReceiptURL, Status: entity.ExpenseStatusAwaitingApproval, CreatedAt: time.Date(2025, 8, 6, 4, 2, 30, 000, time.UTC)},
	}

	//  users table
	logger.Info("seeding users table ...")
	for _, u := range users {
		_, err = tx.Exec(ctx,
			`INSERT INTO users (id, email, name, password_hash, role, created_at) 
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			u.ID, u.Email, u.Name, u.PasswordHash, u.Role, u.CreatedAt,
		)
		if err != nil {
			return
		}
	}
	logger.Info("password for each user is the same as their email prefix (before @)")
	logger.Info("seeding users table completed")

	// expenses table
	logger.Info("seeding expenses table ...")
	for _, e := range expenses {
		var processedAt *time.Time
		if e.Status == entity.ExpenseStatusCompleted {
			t := e.CreatedAt.Add(10 * time.Second)
			processedAt = &t
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO expenses (id, user_id, amount, description, receipt_url, status, created_at, processed_at) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			e.ID, e.UserID, e.Amount, e.Description, e.ReceiptURL, e.Status, e.CreatedAt, processedAt,
		)
		if err != nil {
			return
		}
	}

	logger.Info("seeding approvals table ...")
	for _, e := range expenses {
		validStatus := e.Status == entity.ExpenseStatusCompleted ||
			e.Status == entity.ExpenseStatusRejected ||
			e.Status == entity.ExpenseStatusApproved
		if validStatus && e.RequiresApproval() {
			var (
				approverID uint64
				status     entity.ApprovalStatus
				notes      string
			)

			switch e.UserID {
			case 1:
				approverID = 2
			case 2:
				approverID = 1
			default:
				approverID = 1
			}

			switch e.Status {
			case entity.ExpenseStatusApproved, entity.ExpenseStatusCompleted:
				status = entity.ApprovalStatusApproved
				notes = "Approved from me"
			default:
				status = entity.ApprovalStatusRejected
				notes = "Please check again!"
			}

			_, err = tx.Exec(ctx,
				`INSERT INTO approvals (expense_id, approver_id, status, notes, created_at) 
				 VALUES ($1, $2, $3, $4, $5)`,
				e.ID, approverID, status, notes, e.CreatedAt.Add(5*time.Second),
			)
			if err != nil {
				return
			}
		}
	}
	logger.Info("seeding approvals table completed")

	// reset sequences
	logger.Info("reseting sequences ...")
	tables := []string{"users", "expenses", "approvals"}
	for _, t := range tables {
		query := fmt.Sprintf(`
			SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE(MAX(id), 1)) FROM %s
		`, t, t)

		_, err = tx.Exec(ctx, query)
		if err != nil {
			return
		}
	}
	log.Println("sequences reset to last inserted IDs")
}
