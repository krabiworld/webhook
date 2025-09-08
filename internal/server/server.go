package server

import (
	"log/slog"
	"net/http"
	"os"
	"webhook/internal/config"
	"webhook/internal/server/routes"
)

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", routes.Health)
	mux.HandleFunc("POST /{id}/{token}", routes.Webhook)

	err := http.ListenAndServe(config.Get().Address, mux)
	if err != nil {
		slog.Error("Failed to start server", "err", err.Error())
		os.Exit(1)
	}
}
