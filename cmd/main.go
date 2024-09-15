package main

import (
	"log"
	"os"

	"botus/internal/app"
	"botus/internal/service"
	"botus/pkg/config"
)

func main() {
	if err := config.Load(".env"); err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN не установлен")
	}

	hashService := service.NewHashService()
	hashService.AddHash("secret", "5ebe2294ecd0e0f08eab7690d2a6ee69") // Добавляем соответствие "secret" -> MD5-хеш

	bot, err := app.NewBot(botToken, hashService)
	if err != nil {
		log.Fatalf("Не удалось создать бота: %v", err)
	}

	bot.Start()
}
