package app

import (
	"botus/internal/service"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api         *tgbotapi.BotAPI
	hashService *service.HashService
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
		if update.Message == nil {
			continue
		}

		input := update.Message.Text

		if isEnglish(input) {
			hash := b.hashService.GenerateMD5(input)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("MD5-хеш для '%s': %s", input, hash))
			b.api.Send(msg)
		} else if isHexadecimal(input) && len(input) == 32 {
			value, err := b.hashService.FindValueByHash(input)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				b.api.Send(msg)
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Значение для хеша '%s': %s", input, value))
			b.api.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, вводите только английские буквы или MD5-хеши.")
			b.api.Send(msg)
		}
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
	return true
}
