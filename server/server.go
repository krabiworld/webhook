package server

import (
	"net/http"
	"webhook/config"

	"github.com/rs/zerolog/log"
)

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/", webhook)

	err := http.ListenAndServe(config.Get().Address, mux)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
