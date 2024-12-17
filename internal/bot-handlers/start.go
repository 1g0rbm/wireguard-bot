package bothandlers

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"wireguard-bot/internal/services"
	"wireguard-bot/internal/utils"
	"wireguard-bot/internal/utils/dispatcher"
)

const startCommand = "/start"

type StartHandler struct {
	dispatChan  chan<- dispatcher.Sendable
	userService services.UserService
}

func NewStartHandler(d chan<- dispatcher.Sendable, s services.UserService) *StartHandler {
	return &StartHandler{
		dispatChan:  d,
		userService: s,
	}
}

func (h *StartHandler) Match(update *models.Update) bool {
	if update.Message == nil {
		return false
	}

	return update.Message.Text == startCommand
}

func (h *StartHandler) Handle(ctx context.Context, update *models.Update) error {
	if err := h.handleGreetings(update); err != nil {
		return fmt.Errorf("handler_start.handle.greetings %w", err)
	}

	userModel, err := h.userService.GetOrCreate(
		ctx,
		0,
		update.Message.Chat.ID,
		update.Message.Chat.Username,
		update.Message.Chat.FirstName,
		update.Message.Chat.LastName,
	)
	if err != nil {
		return fmt.Errorf("handler_start.handle.get_or_create %w", err)
	}

	if userModel.Enabled() {
		err = h.handleEnabledUser(update)
	} else {
		err = h.handleNewUser(update)
	}

	if err != nil {
		return fmt.Errorf("handler_start.handle %w", err)
	}

	return nil
}

func (h *StartHandler) handleNewUser(update *models.Update) error {
	msg, err := utils.Render(
		"static/messages/user_created.tmp",
		map[string]string{"Username": update.Message.Chat.Username},
	)
	if err != nil {
		return fmt.Errorf("handler_start.create_user.render_user_created %w", err)
	}
	h.dispatChan <- dispatcher.TextMessage{
		Params: &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      string(msg),
			ParseMode: models.ParseModeMarkdown,
		},
	}

	return nil
}

func (h *StartHandler) handleEnabledUser(update *models.Update) error {
	keyboard := &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{{Text: configCommand}},
			{{Text: qrCodeCommand}},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}

	msg, err := utils.Render(
		"static/messages/user_already_enabled.tmp",
		map[string]string{"Username": update.Message.Chat.Username},
	)
	if err != nil {
		return fmt.Errorf("handler_start.handle_enabled_user.render %w", err)
	}
	h.dispatChan <- dispatcher.TextMessage{
		Params: &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        string(msg),
			ParseMode:   models.ParseModeMarkdown,
			ReplyMarkup: keyboard,
		},
	}

	return nil
}

func (h *StartHandler) handleGreetings(update *models.Update) error {
	msg, err := utils.Render(
		"static/messages/greetings.tmp",
		map[string]string{"Username": update.Message.Chat.Username},
	)
	if err != nil {
		return fmt.Errorf("handler_start.greetings.render_greetings %w", err)
	}
	h.dispatChan <- dispatcher.TextMessage{
		Params: &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   string(msg),
		},
	}

	msg, err = utils.Render("static/messages/about.tmp", nil)
	if err != nil {
		return fmt.Errorf("handler_start.greetings.render_about_message %w", err)
	}
	h.dispatChan <- dispatcher.TextMessage{
		Params: &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      string(msg),
			ParseMode: models.ParseModeMarkdown,
		},
	}

	return nil
}
