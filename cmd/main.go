package main

import (
	"github.com/yourusername/botus/internal/app"
	"github.com/yourusername/botus/internal/config"
	"log"
)

func main() {
	configData, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	bot, err := app.NewBot(configData.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Start()
}
