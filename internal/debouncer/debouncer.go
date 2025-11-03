package debouncer

import (
	"sync"
	"time"
)

type Debouncer struct {
	mu   sync.Mutex
	data map[string]bool
}

var debouncer = &Debouncer{data: make(map[string]bool)}

func Debounce(event, username, repository string, ttl time.Duration) bool {
	key := event + "-" + username + "-" + repository

	debouncer.mu.Lock()

	if _, ok := debouncer.data[key]; ok {
		debouncer.mu.Unlock()
		return false
	}

	debouncer.data[key] = true
	debouncer.mu.Unlock()

	time.AfterFunc(ttl, func() {
		debouncer.mu.Lock()
		delete(debouncer.data, key)
		debouncer.mu.Unlock()
	})

	return true
}
