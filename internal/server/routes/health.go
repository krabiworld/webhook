package routes

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := fmt.Fprint(w, `{"status":"ok"}`); err != nil {
		log.Error().Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
	}
}
