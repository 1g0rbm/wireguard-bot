package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"wireguard-bot/internal/server-handlers/middleaware"
	"wireguard-bot/internal/utils/msgs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-telegram/bot"

	bothandlers "wireguard-bot/internal/bot-handlers"
	"wireguard-bot/internal/closer"
	"wireguard-bot/internal/config"
	"wireguard-bot/internal/config/env"
	"wireguard-bot/internal/db"
	"wireguard-bot/internal/db/pg"
	"wireguard-bot/internal/db/tx"
	"wireguard-bot/internal/repository"
	"wireguard-bot/internal/repository/server"
	"wireguard-bot/internal/repository/session"
	"wireguard-bot/internal/repository/user"
	"wireguard-bot/internal/repository/users2servers"
	serverhandlers "wireguard-bot/internal/server-handlers"
	"wireguard-bot/internal/services"
	configService "wireguard-bot/internal/services/config"
	sessionService "wireguard-bot/internal/services/session"
	userService "wireguard-bot/internal/services/user"
	"wireguard-bot/internal/utils/dhcp"
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

	server    chi.Router
	bot       *bot.Bot
	db        db.Client
	txManager db.TxManager

	startHandler   *bothandlers.StartHandler
	defaultHandler *bothandlers.DefaultHandler
	configHandler  *bothandlers.ConfigHandler
	qrHandler      *bothandlers.QRCodeHandler

	rootHandler  *serverhandlers.RootHandler
	loginHandler *serverhandlers.LoginHandler

	authMiddleware *middleaware.Auth

	userRepo          repository.UserRepository
	sessionRepo       repository.SessionRepository
	serverRepo        repository.ServerRepository
	users2serversRepo repository.Users2Servers

	configService  services.ConfigService
	userService    services.UserService
	sessionService services.SessionService

	txtMsgChan chan *bot.SendMessageParams

	dhcp *dhcp.DHCP
	msgs *msgs.Sender
}

func newContainer() *Container {
	return &Container{}
}

func (c *Container) Server() chi.Router {
	if c.server == nil {
		c.server = chi.NewRouter()

		c.server.Use(middleware.Logger)
		c.server.Use(middleware.Recoverer)
	}

	return c.server
}

func (c *Container) RootHandler() *serverhandlers.RootHandler {
	if c.rootHandler == nil {
		c.rootHandler = serverhandlers.NewRootHandler()
	}

	return c.rootHandler
}

func (c *Container) LoginHandler() *serverhandlers.LoginHandler {
	if c.loginHandler == nil {
		c.loginHandler = serverhandlers.NewLoginHandler(c.UserService(), c.Logger())
	}

	return c.loginHandler
}

func (c *Container) AuthMiddleware() *middleaware.Auth {
	if c.authMiddleware == nil {
		c.authMiddleware = middleaware.NewAuth(c.SessionService())
	}

	return c.authMiddleware
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

func (c *Container) SessionRepo() repository.SessionRepository {
	if c.sessionRepo == nil {
		c.sessionRepo = session.NewRepository(c.DB())
	}

	return c.sessionRepo
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
		c.userService = userService.NewServiceUser(
			c.UserRepo(),
			c.Users2ServersRepo(),
			c.TxManager(), c.DHCP(),
			c.TxtMsgChan(),
		)
	}

	return c.userService
}

func (c *Container) SessionService() services.SessionService {
	if c.sessionService == nil {
		c.sessionService = sessionService.NewServiceSession(c.SessionRepo(), c.TxManager())
	}

	return c.sessionService
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

func (c *Container) StartHandler() *bothandlers.StartHandler {
	if c.startHandler == nil {
		c.startHandler = bothandlers.NewStartHandler(c.UserService(), c.Logger())
	}

	return c.startHandler
}

func (c *Container) ConfigHandler() *bothandlers.ConfigHandler {
	if c.configHandler == nil {
		c.configHandler = bothandlers.NewConfigHandler(c.ConfigService(), c.Logger())
	}

	return c.configHandler
}

func (c *Container) QRCodeHandler() *bothandlers.QRCodeHandler {
	if c.qrHandler == nil {
		c.qrHandler = bothandlers.NewQRCodeHandler(c.ConfigService(), c.Logger())
	}

	return c.qrHandler
}

func (c *Container) DefaultHandler() *bothandlers.DefaultHandler {
	if c.defaultHandler == nil {
		c.defaultHandler = bothandlers.NewDefaultHandler()
	}

	return c.defaultHandler
}

func (c *Container) MsgS() *msgs.Sender {
	if c.msgs == nil {
		c.msgs = msgs.NewSender(c.Bot(), c.TxtMsgChan())
	}

	return c.msgs
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

func (c *Container) TxtMsgChan() chan *bot.SendMessageParams {
	if c.txtMsgChan == nil {
		c.txtMsgChan = make(chan *bot.SendMessageParams)
		c.getCloser().Add(func() error {
			close(c.txtMsgChan)
			return nil
		})
	}

	return c.txtMsgChan
}
