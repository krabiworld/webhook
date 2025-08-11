package server

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := fmt.Fprint(w, `{"status":"ok"}`); err != nil {
		log.Error().Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
	}
}
