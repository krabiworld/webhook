package server

import "github.com/valyala/fasthttp"

func Handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/health":
		health(ctx)
	default:
		webhook(ctx)
	}
}
