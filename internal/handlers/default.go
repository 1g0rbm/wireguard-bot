package handlers

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type DefaultHandler struct {
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (h *DefaultHandler) Match(_ *models.Update) bool {
	return true
}

func (h *DefaultHandler) Handle(_ context.Context, _ *bot.Bot, update *models.Update) {
	fmt.Println(update.Message.Chat.ID, update.Message.Chat.Username, update.Message.Text)
}
