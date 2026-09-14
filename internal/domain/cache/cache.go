package cache

import (
	"context"
	"time"
)

// Store is a small key-value cache. Misses return found=false, never an error.
// Callers treat any error as a miss so the cache can never break a request.
type Store interface {
	Get(ctx context.Context, key string, dest any) (found bool, err error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	// Incr bumps a version counter used to invalidate a whole key family.
	Incr(ctx context.Context, key string) (int64, error)
	GetInt(ctx context.Context, key string) (int64, error)
}
