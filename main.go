package main

import (
	"os"
	"os/signal"
	"syscall"
	"webhook/client"
	"webhook/config"
	"webhook/logger"
	"webhook/server"

	"github.com/rs/zerolog/log"
)

func main() {
	config.Init()

	logger.Init()

	client.Init()
	log.Info().Msg("Client initialized")

	go server.Start()
	log.Info().Str("addr", config.Get().Address).Msg("Server started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
