package config

import (
	"github.com/joho/godotenv"
)

// Load загружает конфигурацию из .env файла
func Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
