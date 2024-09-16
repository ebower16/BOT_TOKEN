package config

import (
	"fmt"
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
		log.Println("BOT_TOKEN not set")
		return nil, fmt.Errorf("BOT_TOKEN not set")
	}

	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	if dbConnStr == "" {
		log.Println("DB_CONNECTION_STRING not set")
		return nil, fmt.Errorf("DB_CONNECTION_STRING not set")
	}

	return &Config{
		BotToken:           botToken,
		DBConnectionString: dbConnStr,
	}, nil
}
