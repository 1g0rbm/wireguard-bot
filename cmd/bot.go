package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wireguard-bot/internal/app"
)

func main() {
	ctx := context.Background()

	application := app.NewApp()

	application.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Application stopped.")
}
