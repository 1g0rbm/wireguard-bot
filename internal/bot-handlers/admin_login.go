package bothandlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"wireguard-bot/internal/services"
	"wireguard-bot/internal/utils"
)

const AdminLoginCallback = "callback.admin_login"

type AdminLoginCallbackHandler struct {
	sessionService services.SessionService
	logger         *slog.Logger
}

func NewAdminLoginCallbackHandler(
	sessionService services.SessionService, logger *slog.Logger,
) *AdminLoginCallbackHandler {
	return &AdminLoginCallbackHandler{
		sessionService: sessionService,
		logger:         logger,
	}
}

func (h *AdminLoginCallbackHandler) Match(update *models.Update) bool {
	if update.CallbackQuery == nil {
		return false
	}

	return update.CallbackQuery.Data == AdminLoginCallback
}

func (h *AdminLoginCallbackHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	if err := h.sessionService.CreateOrUpdate(ctx, chatID); err != nil {
		utils.SendMessage(
			func() ([]byte, error) {
				return utils.Render("static/messages/something_went_wrong.tmp", nil)
			},
			func(msg []byte) error {
				_, err := b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: string(msg)})
				if err != nil {
					return fmt.Errorf("admin_login.handle %w", err)
				}
				return nil
			},
		)
		return
	}

	utils.SendMessage(
		func() ([]byte, error) {
			return utils.Render("static/messages/login_admin.tmp", nil)
		},
		func(msg []byte) error {
			_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    chatID,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text:      string(msg),
			})
			if err != nil {
				return fmt.Errorf("admin_login.handle %w", err)
			}
			return nil
		},
	)
}
