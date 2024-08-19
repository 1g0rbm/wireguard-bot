package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

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
	keyboard := &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{{Text: configCommand}},
			{{Text: qrCodeCommand}},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}

	exist, err := h.userService.IsUserExist(ctx, update.Message.Chat.ID)
	if err != nil {
		utils.SendMessage(
			func() ([]byte, error) {
				return utils.Render("static/something_went_wrong.tmp", nil)
			},
			func(msg []byte) error {
				_, err := b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   string(msg),
				})

				return fmt.Errorf("handler_start.handle %w", err)
			},
		)
		return
	}

	if exist {
		utils.SendMessage(
			func() ([]byte, error) {
				return utils.Render(
					"static/configuration_exist.tmp",
					map[string]string{"Username": update.Message.Chat.Username},
				)
			},
			func(msg []byte) error {
				_, err := b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:      update.Message.Chat.ID,
					Text:        string(msg),
					ReplyMarkup: keyboard,
				})

				return fmt.Errorf("handler_start.handle %w", err)
			},
		)
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
				return utils.Render("static/something_went_wrong.tmp", nil)
			},
			func(msg []byte) error {
				_, err := b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   string(msg),
				})

				return fmt.Errorf("handler_start.handle %w", err)
			},
		)
		return
	}

	utils.SendMessage(
		func() ([]byte, error) {
			return utils.Render(
				"static/configuration_generating.tmp",
				map[string]string{"Username": update.Message.Chat.Username},
			)
		},
		func(msg []byte) error {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        string(msg),
				ReplyMarkup: keyboard,
			})

			return fmt.Errorf("handler_start.handle %w", err)
		},
	)
}
