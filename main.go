package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/database/redis"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/messageHandlers"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/types"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	YandexMusicBase string = "https://music.yandex.ru"
)

var _ = godotenv.Load()
var (
	botToken      = os.Getenv("BOT_TOKEN")
	ownerName     = os.Getenv("OWNER_NAME")
	redisHost     = os.Getenv("REDIS_HOST")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisDB, _    = strconv.Atoi(os.Getenv("REDIS_DB")) // TODO: Вероятно не самое лучшее решение, но что поделать
)

var redisClient = redis.GetClient(redisHost, redisPassword, redisDB)

var mutex sync.Mutex

func FormatTime(t string) (string, error) {
	trackTime, err := time.Parse("2006-01-02T15:04:05", t)
	if err != nil {
		return "", err
	}

	now := time.Now()
	if now.Year() == trackTime.Year() && now.YearDay() == trackTime.YearDay() {
		return "Сегодня в " + trackTime.Format("15:04"), nil
	}

	if now.Year() == trackTime.Year() && now.YearDay()-1 == trackTime.YearDay() {
		return "Вчера в " + trackTime.Format("15:04"), nil
	}

	return trackTime.Format("02.01.2006 в 15:04"), nil
}

func setLastTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var trackInfo types.TrackInfo
	err := decoder.Decode(&trackInfo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	trackInfo.UpdateTime = time.Now().Format("2006-01-02T15:04:05")

	var userId = trackInfo.User.Id
	mutex.Lock()
	err = redis.SaveToRedis(redisClient, userId, trackInfo) // TODO
	if err != nil {
		log.Fatal(err)
	}
	mutex.Unlock()

	fmt.Printf("TrackInfo: %+v\n", trackInfo)
}

func sendCurrentTrack(bot *tgbotapi.BotAPI, update tgbotapi.Update, ownerName string) {
	var userId = update.Message.From.ID

	mutex.Lock()
	track, err := redis.GetFromRedis(redisClient, userId) // TODO
	if err != nil {
		log.Fatal(err)
	}
	mutex.Unlock()

	var artistsText string

	for _, artist := range track.Artists {
		artistsText += fmt.Sprintf("<a href='%s%s'>%s</a>, ", YandexMusicBase, artist.Link, artist.Title)
	}

	artistsText = strings.TrimSuffix(artistsText, ", ")
	trackName := fmt.Sprintf("<a href='%s%s'>%s</a>", YandexMusicBase, track.Link, track.Title)
	trackTime, _ := FormatTime(track.UpdateTime)
	textMsg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\nСейчас слушает: %s\nИсполнитель: %s\nПоследний раз: %s", ownerName, trackName, artistsText, trackTime))
	textMsg.ParseMode = "HTML"
	_, err = bot.Send(textMsg)
	if err != nil {
		log.Panic(err)
		return
	}

	// Отправка изображения
	photoMsg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, strings.ReplaceAll(track.Image, "%%", "300x300"))
	_, err = bot.Send(photoMsg)
	if err != nil {
		log.Panic(err)
		return
	}
}

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

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			if messageHandlers.CommandsHandler(bot, update) {
				continue
			} else {
				sendCurrentTrack(bot, update, ownerName)
			}
		}
	}()

	http.HandleFunc("/set-last-track", setLastTrackHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
