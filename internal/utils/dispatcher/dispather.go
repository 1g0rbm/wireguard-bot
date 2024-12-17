package dispatcher

import (
	"context"
	"log"
	"log/slog"

	"github.com/go-telegram/bot"
)

type MessageCh chan Sendable

type Dispatcher struct {
	ch     MessageCh
	logger *slog.Logger
}

func NewDispatcher(logger *slog.Logger) (*Dispatcher, chan<- Sendable) {
	ch := make(MessageCh)
	return &Dispatcher{ch: ch, logger: logger}, ch
}

func (d *Dispatcher) Run(ctx context.Context, b *bot.Bot) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				d.logger.InfoContext(ctx, "Dispathcer goroutine stopped.")
				return
			case msg, ok := <-d.ch:
				if !ok {
					d.logger.WarnContext(ctx, "Dispathcer channel closed")
					return
				}

				if err := msg.Send(ctx, b); err != nil {
					log.Fatalf("Error sending message: %v", err)
				}
			}
		}
	}()
}
