package messageHandlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func CommandsHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {

	if update.Message.Text == "/start" {
		start(bot, update)
		return true
	} else {
		return false
	}
}
