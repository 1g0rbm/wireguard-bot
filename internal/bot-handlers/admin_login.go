package bothandlers

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"wireguard-bot/internal/services"
	"wireguard-bot/internal/utils"
	"wireguard-bot/internal/utils/dispatcher"
)

const AdminLoginCallback = "callback.admin_login"

type AdminLoginCallbackHandler struct {
	dispatChan     chan<- dispatcher.Sendable
	sessionService services.SessionService
}

func NewAdminLoginCallbackHandler(d chan<- dispatcher.Sendable, s services.SessionService) *AdminLoginCallbackHandler {
	return &AdminLoginCallbackHandler{
		dispatChan:     d,
		sessionService: s,
	}
}

func (h *AdminLoginCallbackHandler) Match(update *models.Update) bool {
	if update.CallbackQuery == nil {
		return false
	}

	return update.CallbackQuery.Data == AdminLoginCallback
}

func (h *AdminLoginCallbackHandler) Handle(ctx context.Context, update *models.Update) error {
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	if err := h.sessionService.CreateOrUpdate(ctx, chatID); err != nil {
		return fmt.Errorf("handler_admin_login.handle.create_session %w", err)
	}

	msg, err := utils.Render("static/messages/login_admin.tmp", nil)
	if err != nil {
		return fmt.Errorf("handler_admin_login.handle.login_admin_render %w", err)
	}

	h.dispatChan <- dispatcher.EditMessage{
		Params: &bot.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text:      string(msg),
		},
	}

	return nil
}
