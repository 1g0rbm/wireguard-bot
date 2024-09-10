package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-telegram/bot"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 15 * time.Second
)

// App is a container which provide interface to manging application.
// Run and stop methods are here.
type App struct {
	container *Container
	bot       *bot.Bot
	server    chi.Router
}

// NewApp creates new instance of App.
func NewApp() *App {
	di := newContainer()

	return &App{
		container: di,
		bot:       di.Bot(),
		server:    di.Server(),
	}
}

func (a *App) Start(ctx context.Context) {
	a.initBotCommandHandlers()
	a.initServerHandlers()

	go a.bot.Start(ctx)

	go func() {
		if err := a.container.MsgS().Run(ctx); err != nil {
			log.Fatalf("Sending message error: %s", err)
		}
	}()

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      a.server,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Starting server error: %s", err)
		}
	}()
}

func (a *App) initServerHandlers() {
	a.server.Group(func(r chi.Router) {
		r.Use(a.container.authMiddleware.HandleFunc)
		a.container.RootHandler().Register(r)
	})

	a.container.LoginHandler().Register(a.server)
}

func (a *App) initBotCommandHandlers() {
	a.bot.RegisterHandlerMatchFunc(a.container.StartHandler().Match, a.container.StartHandler().Handle)
	a.bot.RegisterHandlerMatchFunc(a.container.ConfigHandler().Match, a.container.ConfigHandler().Handle)
	a.bot.RegisterHandlerMatchFunc(a.container.QRCodeHandler().Match, a.container.QRCodeHandler().Handle)
}
