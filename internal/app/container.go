package app

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/go-telegram/bot"

	"wireguard-api/internal/bot_handlers"
	"wireguard-api/internal/closer"
	"wireguard-api/internal/config"
	"wireguard-api/internal/config/env"
	"wireguard-api/internal/db"
	"wireguard-api/internal/db/pg"
	"wireguard-api/internal/db/tx"
	"wireguard-api/internal/repository"
	"wireguard-api/internal/repository/server"
	"wireguard-api/internal/repository/user"
	"wireguard-api/internal/repository/users2servers"
	"wireguard-api/internal/services"
	configService "wireguard-api/internal/services/config"
	userService "wireguard-api/internal/services/user"
	"wireguard-api/internal/utils/dhcp"
)

const (
	logFilePerms = 0600
	mask         = "10.0.0.0/24"
	gateway      = "10.0.0.1"
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

	startHandler   *bot_handlers.StartHandler
	defaultHandler *bot_handlers.DefaultHandler
	configHandler  *bot_handlers.ConfigHandler
	qrHandler      *bot_handlers.QRCodeHandler

	userRepo          repository.UserRepository
	serverRepo        repository.ServerRepository
	users2serversRepo repository.Users2Servers

	configService services.ConfigService
	userService   services.UserService

	dhcp *dhcp.DHCP
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
		logFilePerms,
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

func (c *Container) Users2ServersRepo() repository.Users2Servers {
	if c.users2serversRepo == nil {
		c.users2serversRepo = users2servers.NewRepository(c.DB())
	}

	return c.users2serversRepo
}

func (c *Container) ConfigService() services.ConfigService {
	if c.configService == nil {
		c.configService = configService.NewConfigService(c.Users2ServersRepo())
	}

	return c.configService
}

func (c *Container) UserService() services.UserService {
	if c.userService == nil {
		c.userService = userService.NewServiceUser(c.UserRepo(), c.Users2ServersRepo(), c.TxManager(), c.DHCP())
	}

	return c.userService
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

func (c *Container) StartHandler() *bot_handlers.StartHandler {
	if c.startHandler == nil {
		c.startHandler = bot_handlers.NewStartHandler(c.UserService(), c.Logger())
	}

	return c.startHandler
}

func (c *Container) ConfigHandler() *bot_handlers.ConfigHandler {
	if c.configHandler == nil {
		c.configHandler = bot_handlers.NewConfigHandler(c.ConfigService(), c.Logger())
	}

	return c.configHandler
}

func (c *Container) QRCodeHandler() *bot_handlers.QRCodeHandler {
	if c.qrHandler == nil {
		c.qrHandler = bot_handlers.NewQRCodeHandler(c.ConfigService(), c.Logger())
	}

	return c.qrHandler
}

func (c *Container) DefaultHandler() *bot_handlers.DefaultHandler {
	if c.defaultHandler == nil {
		c.defaultHandler = bot_handlers.NewDefaultHandler()
	}

	return c.defaultHandler
}

func (c *Container) DHCP() *dhcp.DHCP {
	if c.dhcp == nil {
		ctx := context.Background()

		allocatedIPs, err := c.Users2ServersRepo().GetAllAllocatedIPsByServerAlias(ctx, "astana_1")
		if err != nil {
			log.Fatalf("DHCP initialize error: %s", err)
		}

		m := make(map[string]bool)
		for _, allocatedIP := range allocatedIPs {
			m[allocatedIP] = true
		}

		d, err := dhcp.NewDHCP(mask, gateway, m)
		if err != nil {
			log.Fatalf("DHCP initialize error: %s", err)
		}

		c.dhcp = d
	}

	return c.dhcp
}
