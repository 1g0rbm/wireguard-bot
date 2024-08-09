package handlers

import (
	"bytes"
	"context"
	"log"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"wireguard-api/internal/services"
	"wireguard-api/internal/utils"
)

const configCommand = "Конфиг </>"

type VpnConfig struct {
	UserPrivateKey      string
	UserAddress         string
	DNS                 string
	ServerPublicKey     string
	ServerEndpoint      string
	AllowedIPs          string
	PersistentKeepalive int
}

type ConfigHandler struct {
	configService services.ConfigService
	logger        *slog.Logger
}

func NewConfigHandler(configService services.ConfigService, logger *slog.Logger) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
		logger:        logger,
	}
}

func (h *ConfigHandler) Match(update *models.Update) bool {
	return update.Message.Text == configCommand
}

func (h *ConfigHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	cfgBytes, err := h.configService.GenerateConf(ctx, update.Message.Chat.ID)
	if err != nil {
		msgBytes, errRender := utils.Render("static/something_went_wrong.tmp", nil)
		if errRender != nil {
			h.logger.ErrorContext(ctx, "Render message error.", "error", errRender)
		}
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   string(msgBytes),
		})
		if err != nil {
			log.Fatalf("Sending message error.\nerr: %v\n", err)
		}
	}

	msgBytes, errRender := utils.Render("static/sending_conf.tmp", nil)
	if errRender != nil {
		h.logger.ErrorContext(ctx, "Render message error.", "error", errRender)
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   string(msgBytes),
	})
	if err != nil {
		log.Fatalf("Sending message error.\nerr: %v \n", err)
	}

	document := &models.InputFileUpload{
		Filename: update.Message.Chat.Username + ".conf",
		Data:     bytes.NewReader(cfgBytes),
	}
	_, err = b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID:   update.Message.Chat.ID,
		Document: document,
	})
	if err != nil {
		log.Fatalf("Sending message error.\nerr: %v \n", err)
	}
}
