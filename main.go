package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/class"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/events"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/types"
	"os"
)

var _ = godotenv.Load()

var (
	botToken = os.Getenv("BOT_TOKEN")
)

func main() {
	cfg := types.BotConfig{
		Token:   botToken,
		Debug:   false,
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

	e.Start(":8080")
}
