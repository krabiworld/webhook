package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
	"webhook/internal/codec"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Redis[T any] struct {
	client *redis.Client
	ctx    context.Context
	codec  codec.Codec[T]
}

func NewRedis[T any](url string, codec codec.Codec[T]) (*Redis[T], error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Redis[T]{
		client: client,
		ctx:    ctx,
		codec:  codec,
	}, nil
}

func (r *Redis[T]) Set(key string, value T, ttl time.Duration) error {
	s, err := r.codec.Encode(value)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, key, s, ttl).Err()
}

func (r *Redis[T]) Get(key string) (T, error) {
	var zero T

	s, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return zero, nil
		}
		return zero, err
	}

	return r.codec.Decode(s)
}

func (r *Redis[T]) Delete(key string) error {
	result, err := r.client.Del(r.ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	if result == 0 {
		log.Debug().Str("key", key).Msg("Key not found for deletion")
	} else {
		log.Debug().Str("key", key).Msg("Key deleted from Redis")
	}

	return nil
}

func (r *Redis[T]) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}

	return result > 0, nil
}

func (r *Redis[T]) Flush() error {
	err := r.client.FlushDB(r.ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush Redis: %w", err)
	}

	log.Info().Msg("Redis flushed")
	return nil
}
