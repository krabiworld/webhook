package jail

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Jail struct {
	mu   sync.Mutex
	data map[string]bool
}

var jail = &Jail{data: make(map[string]bool)}

func Trap(event, username, repository string, timeout time.Duration) bool {
	key := event + "-" + username + "-" + repository

	jail.mu.Lock()

	if _, ok := jail.data[key]; ok {
		jail.mu.Unlock()
		return false
	}

	jail.data[key] = true
	jail.mu.Unlock()

	log.Debug().Str("key", key).Dur("duration", timeout).Msg("Caught in a trap")

	time.AfterFunc(timeout, func() {
		jail.mu.Lock()
		delete(jail.data, key)
		jail.mu.Unlock()

		log.Debug().Str("key", key).Dur("duration", timeout).Msg("Rescued from the trap")
	})

	return true
}
