package user

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"wireguard-bot/internal/repository/user"
	"wireguard-bot/internal/utils"
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

	msg, err := utils.Render(
		"static/messages/user_enabled.tmp",
		map[string]string{"Username": userModel.Username},
	)
	if err != nil {
		return fmt.Errorf("user_service.enable.message_render %w", err)
	}

	u.outTxtMsgChan <- &bot.SendMessageParams{
		ChatID:    userModel.ID,
		Text:      string(msg),
		ParseMode: models.ParseModeMarkdown,
	}

	return nil
}
