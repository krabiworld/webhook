package debouncer

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Debouncer struct {
	mu   sync.Mutex
	data map[string]bool
}

var debouncer = &Debouncer{data: make(map[string]bool)}

func Debounce(event, username, repository string, timeout time.Duration) bool {
	key := event + "-" + username + "-" + repository

	debouncer.mu.Lock()

	if _, ok := debouncer.data[key]; ok {
		debouncer.mu.Unlock()
		return false
	}

	debouncer.data[key] = true
	debouncer.mu.Unlock()

	log.Debug().Str("key", key).Dur("duration", timeout).Msg("Event debounced")

	time.AfterFunc(timeout, func() {
		debouncer.mu.Lock()
		delete(debouncer.data, key)
		debouncer.mu.Unlock()

		log.Debug().Str("key", key).Dur("duration", timeout).Msg("Event released from debouncer")
	})

	return true
}
