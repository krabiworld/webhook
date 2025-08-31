package debouncer

import (
	"time"
	"webhook/cache"
	"webhook/codec"
	"webhook/config"

	"github.com/rs/zerolog/log"
)

type Debouncer struct {
	cache.Cache[bool]
}

var debouncer *Debouncer

func Init() {
	backend := config.Get().StorageBackend

	switch backend {
	case "memory":
		debouncer = &Debouncer{cache.NewMemory[bool]()}
	case "redis":
		c, err := cache.NewRedis(config.Get().RedisUrl, codec.BoolCodec{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create redis client")
		}

		debouncer = &Debouncer{c}
	default:
		log.Fatal().Msg("Unsupported storage backend")
	}

	log.Info().Str("backend", backend).Msg("Debouncer initialized")
}

func Debounce(event, username, repository string, ttl time.Duration) bool {
	key := event + "-" + username + "-" + repository

	if ok, _ := debouncer.Exists(key); ok {
		return false
	}

	err := debouncer.Set(key, true, ttl)
	if err != nil {
		log.Error().Err(err).Str("key", key).Send()
		return false
	}

	log.Debug().Str("key", key).Dur("duration", ttl).Msg("Event debounced")

	return true
}
