package callbackHandlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/events"
	"regexp"
)

func CallbackHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	println("CallbackHandler ", update.CallbackQuery.Data)

	const matchPattern = "^/([0-9]+)$" // 9 digits pattern

	if match, _ := regexp.MatchString(matchPattern, update.CallbackQuery.Data); match {
		events.SendCurrentTrack(bot, update)
	} else { // Other callbacks probably
		return false
	}

	return true
}
