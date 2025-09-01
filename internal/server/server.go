package server

import (
	"net/http"
	"strings"
	"time"
	"webhook/internal/config"

	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/tomasen/realip"
)

var routes = map[string]string{
	"/health": "/health",
	"/":       "/:id/:token",
}

func sanitizePath(path string) string {
	if v, ok := routes[path]; ok {
		return v
	}

	for prefix, pattern := range routes {
		if prefix == "/" {
			continue
		}
		if strings.HasPrefix(path, prefix) {
			return pattern
		}
	}

	if pattern, ok := routes["/"]; ok {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) == 2 {
			return pattern
		}
	}

	return path
}

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/", webhook)

	c := alice.New()

	c = c.Append(hlog.NewHandler(log.Logger))
	c = c.Append(hlog.MethodHandler("method"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("ip", realip.FromRequest(r)).
			Dur("duration", duration).
			Int("size", size).
			Int("status", status).
			Str("url", sanitizePath(r.URL.Path)).
			Msg("Request")
	}))

	h := c.Then(mux)

	err := http.ListenAndServe(config.Get().Address, h)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
