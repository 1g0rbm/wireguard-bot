package main

import (
	"context"
	"wireguard-bot/internal/app"
)

func main() {
	ctx := context.Background()

	application := app.NewApp()

	application.Start(ctx)
}
