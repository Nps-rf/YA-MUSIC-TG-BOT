package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/callbackHandlers"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/events"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/messageHandlers"
	"log"
	"net/http"
	"os"
)

var _ = godotenv.Load()

var (
	botToken = os.Getenv("BOT_TOKEN")
)

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	f := func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			if err != nil {
				println(err)
				continue
			}

			if messageHandlers.CommandsHandler(bot, update) {
				continue
			} else if update.CallbackQuery != nil {
				callbackHandlers.CallbackHandler(bot, update)

			} else {
				events.SendCurrentTrack(bot, update)
			}
		}
	}
	go f()

	http.HandleFunc("/set-last-track", events.SetLastTrackHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
