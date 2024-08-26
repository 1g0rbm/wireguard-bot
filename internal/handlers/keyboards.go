package handlers

import "github.com/go-telegram/bot/models"

var (
	configKeyboard = &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{{Text: configCommand}},
			{{Text: qrCodeCommand}},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}

	adminConfigKeyboard = &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{{Text: configCommand}},
			{{Text: qrCodeCommand}},
			{{Text: "GetDisabledUserList"}},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
)
