package user

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"

	"wireguard-bot/internal/repository/user"
)

func (u *ServiceUser) Enable(ctx context.Context, userID int64) error {
	userModel, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user_service.enable %w", err)
	}

	userModel.State = user.EnabledState
	if err := u.userRepo.UpdateUser(ctx, userModel); err != nil {
		return fmt.Errorf("user_service.enable %w", err)
	}

	u.outTxtMsgChan <- &bot.SendMessageParams{
		ChatID: userModel.ID,
		Text:   "Активировано. Можешь юзать впн!",
	}

	return nil
}
