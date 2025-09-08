package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"webhook/internal/client"
	"webhook/internal/config"
	"webhook/internal/debouncer"
	"webhook/internal/logger"
	"webhook/internal/server"
)

func main() {
	config.Init()

	logger.Init()

	client.Init()

	debouncer.Init()

	go server.Start()
	slog.Info("Server started", "addr", config.Get().Address)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
