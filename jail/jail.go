package jail

import (
	"sync"
	"time"
)

type Jail struct {
	mu   sync.Mutex
	data map[string]bool
}

var jail = &Jail{data: make(map[string]bool)}

func Trap(event, username, repository string, timeout time.Duration) {
	key := event + "-" + username + "-" + repository

	jail.mu.Lock()

	if _, ok := jail.data[key]; ok {
		jail.mu.Unlock()
		return
	}

	jail.data[key] = true
	jail.mu.Unlock()

	time.AfterFunc(timeout, func() {
		jail.mu.Lock()
		delete(jail.data, key)
		jail.mu.Unlock()
	})
}
