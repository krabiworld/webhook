package server

import (
	"net/http"
	"time"
	"webhook/config"

	"github.com/rs/zerolog/log"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Debug().Str("addr", r.RemoteAddr).Str("method", r.Method).Str("path", r.URL.Path).Dur("ts", time.Since(start)).Send()
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
