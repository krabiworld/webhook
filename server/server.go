package server

import (
	"net/http"
	"strings"
	"time"
	"webhook/config"

	"github.com/rs/zerolog/log"
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Debug().Str("addr", r.RemoteAddr).Str("method", r.Method).Str("path", sanitizePath(r.URL.Path)).Dur("ts", time.Since(start)).Send()
	})
}

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/", webhook)

	err := http.ListenAndServe(config.Get().Address, loggingMiddleware(mux))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
