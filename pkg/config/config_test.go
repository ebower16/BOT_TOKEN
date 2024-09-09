package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Создаем временный .env файл для тестирования
	err := os.WriteFile(".env", []byte("BOT_TOKEN=test_token\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env") // Удаляем файл после теста

	if err := Load(); err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if os.Getenv("BOT_TOKEN") != "test_token" {
		t.Errorf("Expected BOT_TOKEN to be 'test_token', got '%s'", os.Getenv("BOT_TOKEN"))
	}
}
