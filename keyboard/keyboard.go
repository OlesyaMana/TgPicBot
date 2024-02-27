package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("start"),
			tgbotapi.NewKeyboardButton("stop"),
		),
	)
}
