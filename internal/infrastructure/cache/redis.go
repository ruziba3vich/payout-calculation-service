package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ruziba3vich/payout-calculation-service/internal/config"
	"github.com/ruziba3vich/payout-calculation-service/internal/domain/errs"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, cfg config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, errs.Wrap(err, "redis: ping")
	}

	return &Redis{client: client}, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Get(ctx context.Context, key string, dest any) (bool, error) {
	b, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, errs.Wrap(err, "redis: get")
	}
	if err := json.Unmarshal(b, dest); err != nil {
		return false, errs.Wrap(err, "redis: decode")
	}
	return true, nil
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return errs.Wrap(err, "redis: encode")
	}
	return errs.Wrap(r.client.Set(ctx, key, b, ttl).Err(), "redis: set")
}

func (r *Redis) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return errs.Wrap(r.client.Del(ctx, keys...).Err(), "redis: del")
}

func (r *Redis) Incr(ctx context.Context, key string) (int64, error) {
	n, err := r.client.Incr(ctx, key).Result()
	return n, errs.Wrap(err, "redis: incr")
}

func (r *Redis) GetInt(ctx context.Context, key string) (int64, error) {
	n, err := r.client.Get(ctx, key).Int64()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return n, errs.Wrap(err, "redis: get int")
}
