package config

import (
	"log"
	"os"
)

type Config struct {
	BotToken           string
	DBConnectionString string
}

func Load() (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Println("BOT_TOKEN not set; using default")
		botToken = "your_default_bot_token_here"
	}

	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	if dbConnStr == "" {
		log.Println("DB_CONNECTION_STRING not set; using default")
		dbConnStr = "postgres://user:password@db:5432/mydatabase?sslmode=disable"
	}

	return &Config{
		BotToken:           botToken,
		DBConnectionString: dbConnStr,
	}, nil
}
