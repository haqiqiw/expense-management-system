package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=RedisClient --structname RedisClient --outpkg=mocks --output=./../mocks
type RedisClient interface {
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}
