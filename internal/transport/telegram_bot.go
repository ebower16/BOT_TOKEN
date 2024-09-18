package transport

import (
	"fmt"
	"sync"

	"botus/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api         *tgbotapi.BotAPI
	hashService *service.HashService
	mu          sync.Mutex
}

func NewBot(botToken string, hashService *service.HashService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	return &Bot{api: api, hashService: hashService}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !isEnglish(update.Message.Text) {
			continue
		}

		b.mu.Lock()

		if err := b.hashService.IncrementRequestCount(update.Message.From.ID); err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
			b.api.Send(msg)
			b.mu.Unlock()
			continue
		}

		if len(update.Message.Text) == 32 && isHexadecimal(update.Message.Text) {
			value, err := b.hashService.FindValueByHash(update.Message.Text)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
				b.api.Send(msg)
				b.mu.Unlock()
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Value for hash '%s': %s", update.Message.Text, value))
			b.api.Send(msg)
		} else {
			hash, err := b.hashService.AddHash(update.Message.Text)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
				b.api.Send(msg)
				b.mu.Unlock()
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("MD5 hash for '%s': %s", update.Message.Text, hash))
			b.api.Send(msg)
		}

		b.mu.Unlock()
	}
}

func isEnglish(input string) bool {
	for _, char := range input {
		if !(char >= 'a' && char <= 'z') && !(char >= 'A' && char <= 'Z') {
			return false
		}
	}
	return true
}

func isHexadecimal(input string) bool {
	for _, char := range input {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}
	return len(input) == 32
}
