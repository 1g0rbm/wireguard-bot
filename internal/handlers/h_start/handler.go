package h_start

import (
	"context"
	"fmt"
	"wireguard-api/internal/repository/user"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"wireguard-api/internal/repository"
)

const command = "/start"

type Handler struct {
	userRepo repository.UserRepository
}

func NewHandler(userRepo repository.UserRepository) *Handler {
	return &Handler{
		userRepo: userRepo,
	}
}

func (h *Handler) Match(update *models.Update) bool {
	return update.Message.Text == command
}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	userModel := &user.Model{
		Id:         update.Message.Chat.ID,
		Username:   update.Message.Chat.Username,
		FirstName:  update.Message.Chat.FirstName,
		LastName:   update.Message.Chat.LastName,
		Role:       1,
		PrivateKey: "test_private_key",
		PublicKey:  "test_public_key",
	}

	var text string

	err := h.userRepo.CreateUser(ctx, userModel)
	fmt.Println(err)
	if err != nil {
		text = "error!"
		fmt.Println(err)
	} else {
		text = "hi"
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})

	fmt.Printf("sending message error: %v\n", err)
}
