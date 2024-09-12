package bothandlers

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type DefaultHandler struct {
	adminLoginHandler *AdminLoginCallbackHandler
}

func NewDefaultHandler(adminLoginHandler *AdminLoginCallbackHandler) *DefaultHandler {
	return &DefaultHandler{
		adminLoginHandler: adminLoginHandler,
	}
}

func (h *DefaultHandler) Handle(ctx context.Context, bot *bot.Bot, update *models.Update) {
	if h.adminLoginHandler.Match(update) {
		h.adminLoginHandler.Handle(ctx, bot, update)
	} else {
		log.Printf(
			"CahtId: %d\n Username:%s\n Text:%s\n",
			update.Message.Chat.ID, update.Message.Chat.Username, update.Message.Text,
		)
	}
}
