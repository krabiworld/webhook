package debouncer

import (
	"log/slog"
	"time"
	"webhook/internal/cache"
)

type Debouncer struct {
	cache.Cache[bool]
}

var debouncer *Debouncer

func Init() {
	debouncer = &Debouncer{cache.NewMemory[bool]()}

	slog.Info("Debouncer initialized")
}

func Debounce(event, username, repository string, ttl time.Duration) bool {
	key := event + "-" + username + "-" + repository

	if ok, _ := debouncer.Exists(key); ok {
		return false
	}

	err := debouncer.Set(key, true, ttl)
	if err != nil {
		slog.Error(err.Error(), "key", key)
		return false
	}

	slog.Debug("Event debounced", "key", key, "duration", ttl)

	return true
}
