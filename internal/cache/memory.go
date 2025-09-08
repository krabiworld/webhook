package cache

import (
	"log/slog"
	"sync"
	"time"
)

type Memory[T any] struct {
	mu   sync.Mutex
	data map[string]T
}

func NewMemory[T any]() *Memory[T] {
	return &Memory[T]{data: make(map[string]T)}
}

func (m *Memory[T]) Set(key string, value T, ttl time.Duration) error {
	m.mu.Lock()
	m.data[key] = value
	m.mu.Unlock()

	if ttl > 0 {
		time.AfterFunc(ttl, func() {
			if err := m.Delete(key); err != nil {
				slog.Error("Failed to delete key", "key", key, "err", err.Error())
				return
			}

			slog.Debug("Key deleted", "key", key, "duration", ttl)
		})
	}

	return nil
}

func (m *Memory[T]) Delete(key string) error {
	m.mu.Lock()
	delete(m.data, key)
	m.mu.Unlock()
	return nil
}

func (m *Memory[T]) Exists(key string) (bool, error) {
	m.mu.Lock()
	_, ok := m.data[key]
	m.mu.Unlock()
	return ok, nil
}
