package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"webhook/internal/client"
	"webhook/internal/config"
	"webhook/internal/server"
)

func main() {
	config.Init()

	client.Init()

	go server.Start()
	fmt.Println("Server started, addr", config.Get().Address)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
