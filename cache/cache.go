package cache

import "time"

type Cache[T any] interface {
	Set(key string, value T, ttl time.Duration) error
	Get(key string) (T, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	Flush() error
}
