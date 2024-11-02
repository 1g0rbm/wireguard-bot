package msgs

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
)

type Sender struct {
	bot    *bot.Bot
	logger *slog.Logger

	txtMsgChan <-chan *bot.SendMessageParams
}

func NewSender(
	bot *bot.Bot,
	logger *slog.Logger,
	txtMsgChan <-chan *bot.SendMessageParams,
) *Sender {
	return &Sender{
		bot:    bot,
		logger: logger,

		txtMsgChan: txtMsgChan,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("msgs.run %w", ctx.Err())
		case msg := <-s.txtMsgChan:
			_, err := s.bot.SendMessage(ctx, msg)
			if err != nil {
				s.logger.ErrorContext(ctx,
					"Impossible to send message",
					"msg", msg,
					"text", msg.Text,
					"chat id", msg.ChatID,
					"parse mode", msg.ParseMode,
				)
				return fmt.Errorf("msgs.run %w", err)
			}
		}
	}
}
