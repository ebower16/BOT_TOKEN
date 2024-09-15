package main

import (
	"database/sql"
	"log"
	"os"

	"botus/internal/app"
	"botus/internal/config"
	"botus/internal/service"
	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

const (
	maxRequestsPerHour = 100
)

func main() {
	if err := config.Load(".env"); err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN не установлен")
	}

	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	if dbConnStr == "" {
		log.Fatal("DB_CONNECTION_STRING не установлен")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	hashService := service.NewHashService(db)

	bot, err := app.NewBot(botToken, hashService, maxRequestsPerHour)
	if err != nil {
		log.Fatalf("Не удалось создать бота: %v", err)
	}

	bot.Start()
}
