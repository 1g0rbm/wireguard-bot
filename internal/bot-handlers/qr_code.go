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

const qrCodeCommand = "QR-код \uE1D8"

type QRCodeHandler struct {
	dispatChan    chan<- dispatcher.Sendable
	configService services.ConfigService
}

func NewQRCodeHandler(d chan<- dispatcher.Sendable, s services.ConfigService) *QRCodeHandler {
	return &QRCodeHandler{
		dispatChan:    d,
		configService: s,
	}
}

func (h *QRCodeHandler) Match(update *models.Update) bool {
	if update.Message == nil {
		return false
	}

	return update.Message.Text == qrCodeCommand
}

func (h *QRCodeHandler) Handle(ctx context.Context, update *models.Update) error {
	qrBytes, err := h.configService.GenerateQR(ctx, update.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("handler_qr.handle.generate_qr %w", err)
	}

	msg, err := utils.Render("static/messages/sending_qr.tmp", nil)
	if err != nil {
		return fmt.Errorf("handler_qr.handle.sending_qr_render %w", err)
	}
	h.dispatChan <- dispatcher.TextMessage{
		Params: &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   string(msg),
		},
	}

	h.dispatChan <- dispatcher.PhotoMessage{
		Params: &bot.SendPhotoParams{
			ChatID: update.Message.Chat.ID,
			Photo: &models.InputFileUpload{
				Filename: update.Message.Chat.Username + "_qr.png",
				Data:     bytes.NewReader(qrBytes),
			},
		},
	}

	return nil
}
