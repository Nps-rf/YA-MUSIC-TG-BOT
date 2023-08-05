package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/class"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/events"
	"github.com/nps-rf/YA-MUSIC-TG-BOT/types"
	"log"
	"os"
)

var _ = godotenv.Load()

var (
	env = os.Getenv("ENV")
)

func main() {

	if "" == env {
		env = "development"
	}

	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var Debug bool

	if env == "development" {
		Debug = true
	} else {
		Debug = false
	}

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

	e.Start(":8080")
}
