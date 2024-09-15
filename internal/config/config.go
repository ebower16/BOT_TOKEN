package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load(envFilePath string) error {
	if err := godotenv.Load(envFilePath); err != nil {
		log.Printf("Ошибка при загрузке .env файла: %v", err)
		return err
	}
	log.Println("Переменные окружения загружены из файла")

	if os.Getenv("BOT_TOKEN") == "" {
		log.Fatal("BOT_TOKEN не установлен")
	}

	if os.Getenv("DB_CONNECTION_STRING") == "" {
		log.Fatal("DB_CONNECTION_STRING не установлен")
	}

	return nil
}
