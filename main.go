package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/class"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/events"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/types"
	"os"
)

func main() {
	var Debug bool

	//if os.Getenv('ENV') == "development" {
	//	Debug = true
	//} else {
	//	Debug = false
	//}

	cfg := types.BotConfig{
		Token:   os.Getenv("BOT_TOKEN"),
		Debug:   Debug,
		Timeout: 60,
	}

	bot := class.Bot{
		Config: cfg,
	}

	bot.Init()

	go bot.StartPolling()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/set-last-track", events.SetLastTrack)

	err := e.Start(":" + os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}
}
