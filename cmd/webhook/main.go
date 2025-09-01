package main

import (
	"os"
	"os/signal"
	"syscall"
	"webhook/internal/client"
	"webhook/internal/config"
	"webhook/internal/debouncer"
	"webhook/internal/logger"
	"webhook/internal/server"

	"github.com/rs/zerolog/log"
)

func main() {
	config.Init()

	logger.Init()

	client.Init()

	debouncer.Init()

	go server.Start()
	log.Info().Str("addr", config.Get().Address).Msg("Server started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
