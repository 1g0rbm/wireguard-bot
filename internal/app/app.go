package app

import (
	"context"

	"github.com/go-telegram/bot"
)

// App is a container which provide interface to manging application.
// Run and stop methods are here.
type App struct {
	container *Container
	bot       *bot.Bot
}

// NewApp creates new instance of App.
func NewApp() *App {
	di := newContainer()
	b := di.Bot()

	return &App{
		container: di,
		bot:       b,
	}
}

func (a *App) Start(ctx context.Context) {
	a.initCommandHandlers()

	a.bot.Start(ctx)
}

func (a *App) initCommandHandlers() {
	a.bot.RegisterHandlerMatchFunc(a.container.StartHandler().Match, a.container.StartHandler().Handle)
	a.bot.RegisterHandlerMatchFunc(a.container.ConfigHandler().Match, a.container.ConfigHandler().Handle)
	a.bot.RegisterHandlerMatchFunc(a.container.QRCodeHandler().Match, a.container.QRCodeHandler().Handle)
}
