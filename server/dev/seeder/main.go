package main

import (
	"context"
	"expense-management-system/internal/config"
	"fmt"
	"log"
)

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

	logger.Info("seeding users table ...")
	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, email, name, password_hash, role, created_at) VALUES
		(1, 'john@mail.com', 'John', '$2a$10$AVS.3jUpmFXCEgcN8J4Peerq4643.KOMYIhUk8pAgny06t/zxNqHm', 'manager', NOW()),
		(2, 'wawan@mail.com', 'Wawan', '$2a$10$KOIpI26eF/WTIE8s8GV56OCo6qK/GOHaIppLcNX8elxXnuUekOw82', 'manager', NOW()),
		(3, 'budi@mail.com', 'Budi', '$2a$10$ZnT7RBH14iLQXsQGN1Z99O7lASuaJYYiIZcyCpyj9oC8w6m7..Wxu', 'employee', NOW()),
		(4, 'lala@mail.com', 'Lala', '$2a$10$UFBubu4rYw7.ZvVd9rq75Otj12ppjaVOJO/VTBjyc0wkP.fhfBBsO', 'employee', NOW())
	`)
	if err != nil {
		return
	}
	logger.Info("password for each user is the same as their email prefix (before @)")
	logger.Info("seeding users table completed")

	logger.Info("seeding expenses table ...")
	_, err = tx.Exec(ctx, `
		INSERT INTO expenses (id, user_id, amount, description, receipt_url, status, created_at, processed_at) VALUES
		(1, 3, 150000, 'Snacks', 'https://example.com/receipt.jpg', 'completed', NOW(), NOW()),
		(2, 3, 15000, 'Transport KRL', NULL, 'approved', NOW(), NULL),
		(3, 3, 1750000, 'Transport train', 'https://example.com/receipt.jpg', 'completed', NOW(), NOW()), 
		(4, 3, 1250000, 'Transport bus', NULL, 'rejected', NOW(), NULL),
		(5, 3, 7500000, 'Team dinner', 'https://example.com/receipt.jpg', 'approved', NOW(), NULL),
		(6, 3, 2500000, 'Transport plane', 'https://example.com/receipt.jpg', 'awaiting_approval', NOW(), NULL),
		(7, 4, 60000, 'Drinks', 'https://example.com/receipt.jpg', 'completed', NOW(), NOW()),
		(8, 4, 450000, 'Office supplies', NULL, 'approved', NOW(), NULL),
		(9, 4, 1500000, 'Client dinner', 'https://example.com/receipt.jpg', 'awaiting_approval', NOW(), NULL),
		(10, 4, 2570000, 'Rent area', 'https://example.com/receipt.jpg', 'completed', NOW(),  NOW()),
		(11, 1, 4100000, 'Tranport plane', 'https://example.com/receipt.jpg', 'awaiting_approval', NOW(), NULL),
		(12, 2, 1035000, 'Tranport train', 'https://example.com/receipt.jpg', 'awaiting_approval', NOW(), NULL)
	`)
	if err != nil {
		return
	}
	logger.Info("seeding expenses table completed")

	logger.Info("seeding approvals table ...")
	_, err = tx.Exec(ctx, `
		INSERT INTO approvals (id, expense_id, approver_id, status, notes, created_at) VALUES
		(1, 3, 1, 'approved', 'Approved from me', NOW()),
		(2, 4, 1, 'rejected', 'Please attach the receipt', NOW()),
		(3, 5, 1, 'approved', 'Ok', NOW()),
		(4, 10, 2, 'approved', 'Lgtm', NOW())
	`)
	if err != nil {
		return
	}
	logger.Info("seeding approvals table completed")
}
