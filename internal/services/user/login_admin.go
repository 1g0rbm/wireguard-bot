package user

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (u *ServiceUser) LoginAdmin(ctx context.Context, username string) error {
	log.Println(username)
	user, err := u.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("user_repository.login_admin")
	}
	if user == nil {
		return fmt.Errorf("there is no user with username %s", username)
	}

	u.outTxtMsgChan <- &bot.SendMessageParams{
		ChatID: user.ID,
		Text:   "Для логина в админку жмакай 👇",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Login 🔑", CallbackData: "callback.admin_login"},
				},
			},
		},
	}

	return nil
}
