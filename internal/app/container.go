package app

import (
	"log"
	"log/slog"
	"os"

	"github.com/go-telegram/bot"

	"wireguard-api/internal/closer"
	"wireguard-api/internal/config"
	"wireguard-api/internal/config/env"
	"wireguard-api/internal/db"
	"wireguard-api/internal/db/pg"
	"wireguard-api/internal/db/tx"
	"wireguard-api/internal/handlers"
	"wireguard-api/internal/repository"
	"wireguard-api/internal/repository/server"
	"wireguard-api/internal/repository/user"
	"wireguard-api/internal/services"
	configService "wireguard-api/internal/services/config"
)

type Container struct {
	closer *closer.Closer
	logger *slog.Logger

	botCfg    config.BotConfig
	pgCfg     config.PGConfig
	loggerCfg config.LoggerConfig

	bot       *bot.Bot
	db        db.Client
	txManager db.TxManager

	startHandler   *handlers.StartHandler
	defaultHandler *handlers.DefaultHandler
	configHandler  *handlers.ConfigHandler
	qrHandler      *handlers.QRCodeHandler

	userRepo   repository.UserRepository
	serverRepo repository.ServerRepository

	configService services.ConfigService
}

func newContainer() *Container {
	return &Container{}
}

func (c *Container) Logger() *slog.Logger {
	if c.logger != nil {
		return c.logger
	}

	file, err := os.OpenFile(
		c.LogCfg().LogFilepath(),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	c.getCloser().Add(func() error {
		return file.Close()
	})

	c.logger = slog.New(slog.NewJSONHandler(file, nil))

	return c.logger
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

func (c *Container) ServerRepo() repository.ServerRepository {
	if c.serverRepo == nil {
		c.serverRepo = server.NewRepository(c.DB())
	}

	return c.serverRepo
}

func (c *Container) ConfigService() services.ConfigService {
	if c.configService == nil {
		c.configService = configService.NewConfigService(c.UserRepo(), c.ServerRepo())
	}

	return c.configService
}

func (c *Container) LogCfg() config.LoggerConfig {
	if c.loggerCfg == nil {
		c.loggerCfg = env.NewLoggerConfig()
	}

	return c.loggerCfg
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
		bot.WithDefaultHandler(c.DefaultHandler().Handle),
	}

	b, err := bot.New(c.BotCfg().Token(), opts...)
	if err != nil {
		log.Fatalf("Initialize bot error: %s", err)
	}

	c.bot = b

	return c.bot
}

func (c *Container) StartHandler() *handlers.StartHandler {
	if c.startHandler == nil {
		c.startHandler = handlers.NewStartHandler(c.UserRepo())
	}

	return c.startHandler
}

func (c *Container) ConfigHandler() *handlers.ConfigHandler {
	if c.configHandler == nil {
		c.configHandler = handlers.NewConfigHandler(c.ConfigService(), c.Logger())
	}

	return c.configHandler
}

func (c *Container) QRCodeHandler() *handlers.QRCodeHandler {
	if c.qrHandler == nil {
		c.qrHandler = handlers.NewQRCodeHandler(c.ConfigService(), c.Logger())
	}

	return c.qrHandler
}

func (c *Container) DefaultHandler() *handlers.DefaultHandler {
	if c.defaultHandler == nil {
		c.defaultHandler = handlers.NewDefaultHandler()
	}

	return c.defaultHandler
}
