package events

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/consts"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/database/redis"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/types"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	redisHost     = os.Getenv("REDIS_HOST")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisDB, _    = strconv.Atoi(os.Getenv("REDIS_DB")) // TODO: Вероятно не самое лучшее решение, но что поделать
)

var redisClient = redis.GetClient(redisHost, redisPassword, redisDB)

var mutex sync.Mutex

func SendCurrentTrack(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var userId int
	if update.CallbackQuery != nil {
		userId = update.CallbackQuery.Message.From.ID
	} else {
		userId = update.Message.From.ID
	}

	mutex.Lock()
	track, err := redis.GetFromRedis(redisClient, userId)
	if err != nil {
		warnMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Для вас нет трека!")
		_, err = bot.Send(warnMessage)
		mutex.Unlock()
		return
	}
	mutex.Unlock()

	var artistsText string

	for _, artist := range track.Artists {
		artistsText += fmt.Sprintf("<a href='%s%s'>%s</a>, ", consts.YandexMusicBase, artist.Link, artist.Title)
	}

	artistsText = strings.TrimSuffix(artistsText, ", ")
	trackName := fmt.Sprintf("<a href='%s%s'>%s</a>", consts.YandexMusicBase, track.Link, track.Title)
	trackTime, _ := utils.FormatTime(track.UpdateTime)
	textMsg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\nСейчас слушает: %s\nИсполнитель: %s\nПоследний раз: %s", update.Message.From.FirstName, trackName, artistsText, trackTime))
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

func SetLastTrackHandler(w http.ResponseWriter, r *http.Request) {
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

	var userId = "580157064" // TODO
	mutex.Lock()
	err = redis.SaveToRedis(redisClient, userId, trackInfo) // TODO
	if err != nil {
		log.Fatal(err)
	}
	mutex.Unlock()

	fmt.Printf("TrackInfo: %+v\n", trackInfo)
}
