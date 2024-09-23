package app

import (
	"botus/internal/domain"
	"testing"
)

type MockHashService struct{}

func (m *MockHashService) FindHash(value string) (domain.Hash, error) {
	return domain.Hash{Value: value}, nil
}

func TestBotStart(t *testing.T) {
	botToken := "test_token"
	hashService := &MockHashService{}

	bot, err := NewBot(botToken, hashService)
	if err != nil {
		t.Fatalf("Failed to create bot: %v", err)
	}

	if bot.api.Token != botToken {
		t.Fatalf("Expected token %s, got %s", botToken, bot.api.Token)
	}
}
