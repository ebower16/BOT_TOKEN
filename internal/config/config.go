package config

import (
	"log"
	"os"
)

// Config holds the configuration values for the application.
type Config struct {
	BotToken string // Telegram Bot Token
}

// Load loads environment variables into Config struct.
func Load() (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Println("BOT_TOKEN not set; using default")
		botToken = "your_default_bot_token_here" // Set your default bot token here if needed
	}

	return &Config{
		BotToken: botToken,
	}, nil
}
