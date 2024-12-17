package user

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"wireguard-bot/internal/repository/user"
	"wireguard-bot/internal/utils/dispatcher"
)

func (u *ServiceUser) LoginAdmin(ctx context.Context, username string) error {
	log.Println(username)
	usr, err := u.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("user_repository.login_admin")
	}
	if usr == nil || usr.Role != user.AdminRole {
		return fmt.Errorf("there is no user with username %s", username)
	}

	u.tgDispatChan <- dispatcher.TextMessage{
		Params: &bot.SendMessageParams{
			ChatID: usr.ID,
			Text:   "Для логина в админку жмакай 👇",
			ReplyMarkup: &models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{
						{Text: "Login 🔑", CallbackData: "callback.admin_login"},
					},
				},
			},
		},
	}

	return nil
}
