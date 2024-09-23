package app

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/streadway/amqp"
)


type Bot struct {
	api        *tgbotapi.BotAPI
	mu         sync.Mutex
	rabbitConn *amqp.Connection
}


func NewBot(botToken string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	rabbitConn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return &Bot{api: api, rabbitConn: rabbitConn}, nil
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

		if err := b.processMessage(update.Message); err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
			b.api.Send(msg)
		}

		b.mu.Unlock()
	}
}


func (b *Bot) processMessage(message *tgbotapi.Message) error {
	if len(message.Text) == 32 && isHexadecimal(message.Text) {
		
		return nil // Placeholder for actual logic.
	} else {
		return b.sendToRabbitMQ(message.Text)
	}
}


func (b *Bot) sendToRabbitMQ(text string) error {
	ch, err := b.rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"md5_queue", 
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	body := []byte(text)
	err = ch.Publish(
		"",    
		q.Name, 
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return b.logStatistics(text) // Log statistics after sending to RabbitMQ.
}


func (b *Bot) logStatistics(request string) error {
	data := []byte(fmt.Sprintf(`{"request": "%s"}`, request))
	resp, err := http.Post("http://localhost:8080/stats", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Error logging statistics: %v", err)
	}
	defer resp.Body.Close()
	return nil
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
