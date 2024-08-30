package msgs

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

type Sender struct {
	bot *bot.Bot

	txtMsgChan <-chan *bot.SendMessageParams
}

func NewSender(bot *bot.Bot, txtMsgChan <-chan *bot.SendMessageParams) *Sender {
	return &Sender{
		bot:        bot,
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
				return fmt.Errorf("msgs.run %w", err)
			}
		}
	}
}
