package bot_handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"wireguard-api/internal/repository/user"
	"wireguard-api/internal/services"
	"wireguard-api/internal/utils"
)

const startCommand = "/start"

type StartHandler struct {
	userService services.UserService
	logger      *slog.Logger
}

func NewStartHandler(userService services.UserService, logger *slog.Logger) *StartHandler {
	return &StartHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *StartHandler) Match(update *models.Update) bool {
	return update.Message.Text == startCommand
}

func (h *StartHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	utils.SendMessage(
		func() ([]byte, error) {
			return utils.Render(
				"static/messages/greetings.tmp",
				map[string]string{"Username": update.Message.Chat.Username},
			)
		},
		func(msg []byte) error {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      string(msg),
				ParseMode: models.ParseModeMarkdown,
			})
			if err != nil {
				return fmt.Errorf("handler_start.handle %w", err)
			}
			return nil
		},
	)

	utils.SendMessage(
		func() ([]byte, error) { return utils.Render("static/messages/about.tmp", nil) },
		func(msg []byte) error {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      string(msg),
				ParseMode: models.ParseModeMarkdown,
			})
			if err != nil {
				return fmt.Errorf("handler_start.handle %w", err)
			}
			return nil
		},
	)

	userModel, err := h.userService.FindUser(ctx, update.Message.Chat.ID)
	if err != nil {
		utils.SendMessage(
			func() ([]byte, error) {
				return utils.Render("static/messages/something_went_wrong.tmp", nil)
			},
			func(msg []byte) error {
				_, err := b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: string(msg)})
				if err != nil {
					return fmt.Errorf("handler_start.handle %w", err)
				}
				return nil
			},
		)
		return
	}

	if userModel != nil && userModel.StateIs(user.EnabledState) {
		handleEnabledUser(ctx, b, update, userModel)
		return
	}

	err = h.userService.Create(
		ctx,
		0,
		update.Message.Chat.ID,
		update.Message.Chat.Username,
		update.Message.Chat.FirstName,
		update.Message.Chat.LastName,
	)
	if err != nil {
		utils.SendMessage(
			func() ([]byte, error) {
				return utils.Render("static/messages/something_went_wrong.tmp", nil)
			},
			func(msg []byte) error {
				_, err := b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: string(msg)})
				if err != nil {
					return fmt.Errorf("handler_start.handle %w", err)
				}
				return nil
			},
		)
		return
	}

	utils.SendMessage(
		func() ([]byte, error) {
			return utils.Render(
				"static/messages/user_created.tmp",
				map[string]string{"Username": update.Message.Chat.Username},
			)
		},
		func(msg []byte) error {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      string(msg),
				ParseMode: models.ParseModeMarkdown,
			})
			if err != nil {
				return fmt.Errorf("handler_start.handle %w", err)
			}
			return nil
		},
	)
}

func handleEnabledUser(ctx context.Context, b *bot.Bot, update *models.Update, userModel *user.Model) {
	var keyboard *models.ReplyKeyboardMarkup
	if userModel.RoleIs(user.AdminRole) {
		keyboard = &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{{Text: configCommand}},
				{{Text: qrCodeCommand}},
				{{Text: "GetDisabledUserList"}},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		}
	} else {
		keyboard = &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{{Text: configCommand}},
				{{Text: qrCodeCommand}},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		}
	}

	utils.SendMessage(
		func() ([]byte, error) {
			return utils.Render(
				"static/messages/user_already_enabled.tmp",
				map[string]string{"Username": update.Message.Chat.Username},
			)
		},
		func(msg []byte) error {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        string(msg),
				ParseMode:   models.ParseModeMarkdown,
				ReplyMarkup: keyboard,
			})
			if err != nil {
				return fmt.Errorf("handler_start.handle %w", err)
			}
			return nil
		},
	)
}
