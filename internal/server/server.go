package server

import (
	"fmt"
	"gohook/internal/config"
	"gohook/internal/server/routes"
	"os"

	"github.com/valyala/fasthttp"
)

func Start() {
	handler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/health":
			routes.Health(ctx)
		default:
			routes.Webhook(ctx)
		}
	}

	err := fasthttp.ListenAndServe(config.Get().Address, handler)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
