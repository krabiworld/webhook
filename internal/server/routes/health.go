package routes

import (
	"fmt"
	"log/slog"
	"net/http"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := fmt.Fprint(w, `{"status":"ok"}`); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
