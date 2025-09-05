package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

//go:generate mockery --name=Transactioner --structname Transactioner --outpkg=mocks --output=./../mocks
type Transactioner interface {
	Do(ctx context.Context, fn func(Executor) error) error
}

type transactioner struct {
	db PgxIface
}

func NewTransactioner(db PgxIface) Transactioner {
	return &transactioner{
		db: db,
	}
}

func (h *transactioner) Do(ctx context.Context, fn func(Executor) error) (err error) {
	tx, err := h.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction error = %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("panic occured and rollback transaction error = %w", rbErr)
			}
			panic(p)
		}

		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("rollback transaction error = %w", rbErr)
			}
			return
		}

		if cmErr := tx.Commit(ctx); cmErr != nil {
			err = fmt.Errorf("commit transaction error = %w", cmErr)
		}
	}()

	err = fn(tx)

	return err
}
