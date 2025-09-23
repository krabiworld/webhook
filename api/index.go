package handler

import (
	"net/http"
	"webhook/internal/config"
	"webhook/internal/server/routes"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	config.Init()

	routes.Webhook(w, r)
}
