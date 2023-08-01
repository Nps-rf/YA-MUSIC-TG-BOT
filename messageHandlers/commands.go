package messageHandlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	startButton := tgbotapi.NewInlineKeyboardButtonData("Nikolai Pikalov", "580157064")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			startButton,
		))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, выберите человека:")

	msg.ReplyMarkup = keyboard

	bot.Send(msg)

}
