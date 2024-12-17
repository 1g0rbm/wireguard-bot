package bothandlers

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"wireguard-bot/internal/utils"
	"wireguard-bot/internal/utils/dispatcher"
)

type Handler interface {
	Match(update *models.Update) bool
	Handle(ctx context.Context, update *models.Update) error
}

type DefaultHandler struct {
	handlers   []Handler
	dispatChan chan<- dispatcher.Sendable
	logger     *slog.Logger
}

func NewDefaultHandler(ch chan<- dispatcher.Sendable, h []Handler, l *slog.Logger) *DefaultHandler {
	return &DefaultHandler{
		handlers:   h,
		dispatChan: ch,
		logger:     l,
	}
}

func (dh *DefaultHandler) Handle(ctx context.Context, _ *bot.Bot, update *models.Update) {
	matched := false
	for _, h := range dh.handlers {
		if h.Match(update) {
			matched = true
			if err := h.Handle(ctx, update); err != nil {
				dh.logger.ErrorContext(ctx, "Handling message error.", "err", err)
				msg, err := utils.Render("static/messages/something_went_wrong.tmp", nil)
				if err != nil {
					msg = []byte("a")
				}

				dh.dispatChan <- dispatcher.TextMessage{
					Params: &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   string(msg),
					},
				}
			}
		}
	}

	if !matched {
		dh.logger.WarnContext(ctx, "Wrong command.", "msg", update)
	}
}
