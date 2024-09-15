package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load(envFilePath string) error {
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		log.Printf("Файл %s не найден, загружаем переменные окружения", envFilePath)
	} else {
		if err := godotenv.Load(envFilePath); err != nil {
			log.Printf("Ошибка при загрузке .env файла: %v", err)
			return err
		}
		log.Println("Переменные окружения загружены из файла")
	}

	if os.Getenv("BOT_TOKEN") == "" {
		log.Fatal("BOT_TOKEN не установлен в окружении")
	}

	return nil
}
