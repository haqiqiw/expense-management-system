package config

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(ctx context.Context, env *Env) (*pgxpool.Pool, error) {
	params := url.Values{}

	if env.DBSSLMode != "" {
		params.Set("sslmode", env.DBSSLMode)
	} else {
		params.Set("sslmode", "disable")
	}

	if env.DBConnectTimeout > 0 {
		params.Set("connect_timeout", strconv.Itoa(env.DBConnectTimeout))
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		env.DBUsername,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
		params.Encode(),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	config.MaxConns = int32(env.DBMaxConn)
	config.MinConns = int32(env.DBMinConn)
	config.MaxConnLifetime = time.Second * time.Duration(env.DBMaxConnLifetime)

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create db connection pool: %w", err)
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return dbpool, nil
}
