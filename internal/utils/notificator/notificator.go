package notificator

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"

	"wireguard-bot/internal/repository"
	"wireguard-bot/internal/repository/user"
	"wireguard-bot/internal/utils/dispatcher"
)

const AdminEvent = "admin"

type Notificator struct {
	dispatChan chan<- dispatcher.Sendable
	userRepo   repository.UserRepository
}

func NewNotificator(ch chan<- dispatcher.Sendable, r repository.UserRepository) *Notificator {
	return &Notificator{
		dispatChan: ch,
		userRepo:   r,
	}
}

func (n *Notificator) Notify(ctx context.Context, event string) error {
	if event != AdminEvent {
		return fmt.Errorf("invalid notify event name %s", event)
	}

	admins, err := n.userRepo.FindByRole(ctx, user.AdminRole)
	if err != nil {
		return fmt.Errorf("notificator.notify.find_user_by_role %w", err)
	}

	for _, admin := range admins {
		n.dispatChan <- dispatcher.TextMessage{
			Params: &bot.SendMessageParams{
				ChatID: admin.ID,
				Text:   "Одмэн!\nНовый халявщик встал в очередь за VPN.\nПрими меры!",
			},
		}
	}

	return nil
}
