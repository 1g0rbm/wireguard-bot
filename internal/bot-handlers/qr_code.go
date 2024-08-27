package bothandlers

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

const qrCodeCommand = "QR-код \uE1D8"

type QRCodeHandler struct {
	configService services.ConfigService
	logger        *slog.Logger
}

func NewQRCodeHandler(configService services.ConfigService, logger *slog.Logger) *QRCodeHandler {
	return &QRCodeHandler{
		configService: configService,
		logger:        logger,
	}
}

func (h *QRCodeHandler) Match(update *models.Update) bool {
	return update.Message.Text == qrCodeCommand
}

func (h *QRCodeHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	qrBytes, err := h.configService.GenerateQR(ctx, update.Message.Chat.ID)
	if err != nil {
		msgBytes, errRender := utils.Render("static/messages/something_went_wrong.tmp", nil)
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

	msgBytes, errRender := utils.Render("static/messages/sending_qr.tmp", nil)
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
		Filename: update.Message.Chat.Username + "_qr.png",
		Data:     bytes.NewReader(qrBytes),
	}
	_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: update.Message.Chat.ID,
		Photo:  document,
	})
	if err != nil {
		log.Fatalf("Sending message error.\nerr: %v \n", err)
	}
}
