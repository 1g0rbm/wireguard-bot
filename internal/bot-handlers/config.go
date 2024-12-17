package bothandlers

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"wireguard-bot/internal/services"
	"wireguard-bot/internal/utils"
	"wireguard-bot/internal/utils/dispatcher"
)

const ConfigCommand = "Конфиг </>"

type ConfigHandler struct {
	dispatChan    chan<- dispatcher.Sendable
	configService services.ConfigService
}

func NewConfigHandler(d chan<- dispatcher.Sendable, s services.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		dispatChan:    d,
		configService: s,
	}
}

func (h *ConfigHandler) Match(update *models.Update) bool {
	if update.Message == nil {
		return false
	}

	return update.Message.Text == ConfigCommand
}

func (h *ConfigHandler) Handle(ctx context.Context, update *models.Update) error {
	cfgBytes, err := h.configService.GenerateConf(ctx, update.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("handler_config.handle.generate_config %w", err)
	}

	msg, err := utils.Render("static/messages/sending_conf.tmp", nil)
	if err != nil {
		return fmt.Errorf("handler_config.handle.sending_conf_render %w", err)
	}
	h.dispatChan <- dispatcher.TextMessage{
		Params: &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   string(msg),
		},
	}

	confName := update.Message.Chat.Username
	if confName == "" {
		confName = fmt.Sprintf("%d", update.Message.Chat.ID)
	}

	h.dispatChan <- dispatcher.DocumentMessage{
		Params: &bot.SendDocumentParams{
			ChatID: update.Message.Chat.ID,
			Document: &models.InputFileUpload{
				Filename: confName + ".conf",
				Data:     bytes.NewReader(cfgBytes),
			},
		},
	}

	return nil
}
