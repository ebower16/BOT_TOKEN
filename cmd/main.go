package main

import (
	"log"
	"os"

	"botus/internal/bot"
	"botus/pkg/config"
)

func main() {
	// Загружаем конфигурацию
	if err := config.Load(); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	bot.RunBot(botToken)
}
