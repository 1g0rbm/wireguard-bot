package app

import (
	"log"

	"github.com/go-telegram/bot"

	"wireguard-api/internal/closer"
	"wireguard-api/internal/config"
	"wireguard-api/internal/config/env"
	"wireguard-api/internal/db"
	"wireguard-api/internal/db/pg"
	"wireguard-api/internal/db/tx"
	"wireguard-api/internal/handlers/h_start"
	"wireguard-api/internal/repository"
	"wireguard-api/internal/repository/user"
)

type Container struct {
	closer *closer.Closer

	botCfg config.BotConfig
	pgCfg  config.PGConfig

	bot       *bot.Bot
	db        db.Client
	txManager db.TxManager

	startHandler *h_start.Handler

	userRepo repository.UserRepository
}

func newContainer() *Container {
	return &Container{}
}

func (c *Container) TxManager() db.TxManager {
	if c.txManager == nil {
		c.txManager = tx.NewTxManager(c.DB().DB())
	}

	return c.txManager
}

func (c *Container) getCloser() *closer.Closer {
	if c.closer == nil {
		c.closer = closer.NewCloser()
	}

	return c.closer
}

func (c *Container) DB() db.Client {
	if c.db == nil {
		dbc, err := pg.NewClient(c.PgCfg().DSN())
		if err != nil {
			log.Fatalf("DB initialization err: %v\n", err)
		}

		c.getCloser().Add(func() error {
			return dbc.Close()
		})

		c.db = dbc
	}

	return c.db
}

func (c *Container) UserRepo() repository.UserRepository {
	if c.userRepo == nil {
		c.userRepo = user.NewRepository(c.DB())
	}

	return c.userRepo
}

func (c *Container) BotCfg() config.BotConfig {
	if c.botCfg != nil {
		return c.botCfg
	}

	c.botCfg = env.NewBotConfig()

	return c.botCfg
}

func (c *Container) PgCfg() config.PGConfig {
	if c.pgCfg != nil {
		return c.pgCfg
	}

	c.pgCfg = env.NewPgConfig()

	return c.pgCfg
}

func (c *Container) Bot() *bot.Bot {
	if c.bot != nil {
		return c.bot
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(nil),
	}

	b, err := bot.New(c.BotCfg().Token(), opts...)
	if err != nil {
		log.Fatalf("Initialize bot error: %s", err)
	}

	c.bot = b

	return c.bot
}

func (c *Container) StartHandler() *h_start.Handler {
	if c.startHandler == nil {
		c.startHandler = h_start.NewHandler(c.UserRepo())
	}

	return c.startHandler
}

func (c *Container) EnableStartHandler() {
	c.Bot().RegisterHandlerMatchFunc(c.StartHandler().Match, c.StartHandler().Handle)
}
