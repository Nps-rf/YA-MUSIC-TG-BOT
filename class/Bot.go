package class

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/callbackHandlers"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/events"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/messageHandlers"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/types"
	"log"
)

type Bot struct {
	Config   types.BotConfig
	Instance *tgbotapi.BotAPI
}

func (bot *Bot) Init() {
	Instance, err := tgbotapi.NewBotAPI(bot.Config.Token)

	if err != nil {
		log.Panic(err)
	}

	Instance.Debug = bot.Config.Debug

	if Instance.Debug {
		log.Printf("Authorized on account %s", Instance.Self.UserName)
	}

	bot.Instance = Instance
}

func (bot *Bot) StartPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = bot.Config.Timeout

	updates, err := bot.Instance.GetUpdatesChan(u)
	if err != nil {
		println(err)
		return
	}

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			if messageHandlers.CommandsHandler(bot.Instance, update) {
				continue
			} else if update.CallbackQuery != nil {
				callbackHandlers.CallbackHandler(bot.Instance, update)
			} else {
				events.SendCurrentTrack(bot.Instance, update)
			}
		}
	}()
}
