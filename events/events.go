package events

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/labstack/echo/v4"
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
		log.Panic(err)
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
	}
}

func SetLastTrack(c echo.Context) error {
	var requestData types.RequestData

	if err := json.NewDecoder(c.Request().Body).Decode(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}
	defer c.Request().Body.Close()

	requestData.TrackInfo.UpdateTime = time.Now().Format("2006-01-02T15:04:05")

	mutex.Lock()
	err := redis.SaveToRedis(redisClient, requestData.UserInfo.Id, requestData.TrackInfo) // TODO
	mutex.Unlock()

	if err != nil {
		_ = c.JSON(http.StatusInternalServerError, err.Error())
		log.Panic(err)
	}

	fmt.Printf("TrackInfo: %+v\n", requestData)
	return c.JSON(http.StatusOK, requestData)
}
