package handlers

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"wireguard-api/internal/repository"
	"wireguard-api/internal/repository/user"
	"wireguard-api/internal/utils"
)

const startCommand = "/start"

type StartHandler struct {
	userRepo repository.UserRepository
}

func NewStartHandler(userRepo repository.UserRepository) *StartHandler {
	return &StartHandler{
		userRepo: userRepo,
	}
}

func (h *StartHandler) Match(update *models.Update) bool {
	return update.Message.Text == startCommand
}

func (h *StartHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {

	keyboard := &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{{Text: configCommand}},
			{{Text: "QR-код \uE1D8"}},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}

	userModel, err := h.userRepo.GetUserById(ctx, update.Message.Chat.ID)
	if err != nil {
		fmt.Println(err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Что-то пошло не так. \n Попробуй позже.",
		})
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if userModel != nil {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        fmt.Sprintf("Привет, %s!\n У тебя уже есть конфигурация.", userModel.Username),
			ReplyMarkup: keyboard,
		})
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	privateKey, publicKey, err := utils.GenerateKeyPair()
	if err != nil {
		fmt.Println(err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ну удалось сгенерировать конфигурацию. \n Попробуй позже.",
		})
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	userModel = &user.Model{
		Id:         update.Message.Chat.ID,
		Username:   update.Message.Chat.Username,
		FirstName:  update.Message.Chat.FirstName,
		LastName:   update.Message.Chat.LastName,
		Role:       1,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	if err := h.userRepo.CreateUser(ctx, userModel); err != nil {
		fmt.Println(err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ну удалось сгенерировать конфигурацию. \n Попробуйте позже.",
		})
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        fmt.Sprintf("Привет, %s!\n Твоя конфигурация успешно сгенерирована!", userModel.Username),
		ReplyMarkup: keyboard,
	})
	if err != nil {
		fmt.Println(err)
	}
}
