package main

import (
	"fmt"
	"gohook/internal/client"
	"gohook/internal/config"
	"gohook/internal/server"
	"os"
	"os/signal"
	"syscall"
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
