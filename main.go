package main

import (
	"webhook/client"
	"webhook/config"
	"webhook/logger"
	"webhook/server"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func main() {
	config.Init()

	logger.Init()

	log.Info().Msg("Initializing client")
	client.Init()

	log.Info().Str("addr", config.Get().Address).Msg("Starting server...")
	err := fasthttp.ListenAndServe(config.Get().Address, server.Handler)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
