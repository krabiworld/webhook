package server

import (
	"net/http"
	"strings"
	"time"
	"webhook/internal/config"
	"webhook/internal/middlewares"
	"webhook/internal/server/routes"
	"webhook/internal/server/routes/api"

	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/tomasen/realip"
)

var r = map[string]string{
	"/health": "/health",
	"/":       "/:id/:token",
}

func sanitizePath(path string) string {
	if v, ok := r[path]; ok {
		return v
	}

	for prefix, pattern := range r {
		if prefix == "/" {
			continue
		}
		if strings.HasPrefix(path, prefix) {
			return pattern
		}
	}

	if pattern, ok := r["/"]; ok {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) == 2 {
			return pattern
		}
	}

	return path
}

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", routes.Health)
	mux.HandleFunc("POST /{id}/{token}", routes.Webhook)
	//mux.HandleFunc("POST /webhook/{id}", routes.WebhookNew)

	// Webhooks
	mux.HandleFunc("POST /api/webhooks", api.CreateWebhook)
	//mux.HandleFunc("GET /api/webhooks", api.GetWebhooks)
	//mux.HandleFunc("GET /api/webhooks/{id}", api.GetWebhook)
	//mux.HandleFunc("PUT /api/webhooks/{id}", api.PutWebhook)
	//mux.HandleFunc("DELETE /api/webhooks/{id}", api.DeleteWebhook)

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

	c = c.Append(middlewares.Authentication)

	h := c.Then(mux)

	err := http.ListenAndServe(config.Get().Address, h)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
