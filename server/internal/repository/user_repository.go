package repository

import (
	"context"
	"errors"
	"expense-management-system/internal/db"
	"expense-management-system/internal/entity"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db db.PgxIface
}

func NewUserRepository(db db.PgxIface) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	now := time.Now()
	query := `
		INSERT INTO users (email, name, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := r.db.QueryRow(ctx, query, user.Email, user.Name, user.PasswordHash, user.Role, now).Scan(&user.ID)
	if err != nil {
		return err
	}

	user.CreatedAt = now

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	query := `SELECT id, email, name, password_hash, role, created_at FROM users WHERE id = $1 LIMIT 1`

	var u entity.User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, email, name, password_hash, role, created_at FROM users WHERE email = $1 LIMIT 1`

	var u entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &u, nil
}

func (r *UserRepository) CountByEmail(ctx context.Context, email string) (int, error) {
	query := `SELECT COUNT(id) FROM users WHERE email = $1`

	var count int
	err := r.db.QueryRow(ctx, query, email).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
