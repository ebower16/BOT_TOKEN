package config

import (
	"log"
	"os"
)


type Config struct {
	BotToken string 
}

func Load() (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Println("BOT_TOKEN not set; using default")
		botToken = "your_default_bot_token_here" 
	}

	return &Config{
		BotToken: botToken,
	}, nil
}
