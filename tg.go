package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

type ArtistInfo struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type TrackInfo struct {
	Artists []ArtistInfo `json:"artists"`
	Image   string       `json:"cover"`
	Title   string       `json:"title"`
	Link    string       `json:"link"`
}

var lastTrackInfo TrackInfo
var mutex sync.Mutex

func setLastTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var trackInfo TrackInfo
	err := decoder.Decode(&trackInfo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	lastTrackInfo = trackInfo
	mutex.Unlock()

	fmt.Printf("TrackInfo: %+v\n", trackInfo)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")
	ownerName := os.Getenv("OWNER_NAME")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			mutex.Lock()
			track := lastTrackInfo
			mutex.Unlock()

			var artistsText string
			for _, artist := range track.Artists {
				artistsText += fmt.Sprintf("<a href='https://music.yandex.ru%s'>%s</a>, ", artist.Link, artist.Title)
			}
			artistsText = strings.TrimSuffix(artistsText, ", ")
			trackName := fmt.Sprintf("<a href='https://music.yandex.ru%s'>%s</a>", track.Link, track.Title)
			textMsg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\nСейчас слушает: %s\nИсполнитель: %s", ownerName, trackName, artistsText))
			textMsg.ParseMode = "HTML"
			bot.Send(textMsg)

			// Отправка изображения
			photoMsg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, strings.ReplaceAll(track.Image, "%%", "300x300"))
			bot.Send(photoMsg)
		}
	}()

	http.HandleFunc("/set-last-track", setLastTrackHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
