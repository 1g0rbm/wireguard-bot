package main

import (
	"context"
	"wireguard-api/internal/app"
)

func main() {
	ctx := context.Background()

	application := app.NewApp()

	application.Start(ctx)
}
