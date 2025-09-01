package cache

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
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
				log.Error().Err(err).Msg("Failed to delete key")
				return
			}

			log.Debug().Str("key", key).Dur("duration", ttl).Msg("Key deleted")
		})
	}

	return nil
}

func (m *Memory[T]) Get(key string) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data[key], nil
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

func (m *Memory[T]) Flush() error {
	m.mu.Lock()
	m.data = make(map[string]T)
	m.mu.Unlock()
	return nil
}
