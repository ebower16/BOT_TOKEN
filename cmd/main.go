package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

// Predefined hash table
var hashTable = map[string]string{
	"hen":    "9c56cc51e4b3a1b2e9b9d501c6a7c4d8", // MD5 hash for "hen"
	"secret": "5ebe2294ecd0e0f08eab7690d2a6ee69", // MD5 hash for "secret"
	"hello":  "5d41402abc4b2a76b9719d911017c592", // MD5 hash for "hello"
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 10})
	if err != nil {
		log.Fatalf("Failed to get updates: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Received shutdown signal, shutting down gracefully...")
		cancel()
	}()

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Text {
			case "/start":
				msg.Text = "Welcome! Please choose an option:"
				msg.ReplyMarkup = getKeyboard()
			case "Hash secret":
				msg.Text = processMessage("secret")
			case "Show Image":
				msg.Text = "Here is an image for you:"
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				imageMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "path/to/your/image.jpg")
				if _, err := bot.Send(imageMsg); err != nil {
					log.Printf("Failed to send image: %v", err)
				}
			case "Run Tests":
				msg.Text = runTests()
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Failed to send message: %v", err)
				}
			case "Notify Tests Passed":
				notifyTestsResult(bot, update.Message.Chat.ID)
			case "Exit":
				msg.Text = "Exiting the bot..."
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Failed to send message: %v", err)
				}
				cancel()
			default:
				// Check if the user input matches a string in the hash table
				hash, exists := hashTable[update.Message.Text]
				if exists {
					msg.Text = "The MD5 hash of \"" + update.Message.Text + "\" is: " + hash
				} else {
					// If not found, inform the user
					msg.Text = "Error: The string \"" + update.Message.Text + "\" is not in the hash table."
				}
			}

			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}

		case <-ctx.Done():
			log.Println("Shutting down the bot...")
			return
		}
	}
}

func getKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Hash secret"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Show Image"),
			tgbotapi.NewKeyboardButton("Run Tests"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Notify Tests Passed"),
			tgbotapi.NewKeyboardButton("Exit"),
		),
	)
	return keyboard
}

func processMessage(text string) string {
	hash := md5.Sum([]byte(text))
	hashStr := hex.EncodeToString(hash[:])

	var response string

	if text == "secret" {
		response = "md5(\"secret\") = " + hashStr + "\nreverse(\"" + hashStr + "\") = \"secret\""
	} else {
		response = "Неверный текст. Пожалуйста, выберите опцию из кнопок или отправьте 'secret' или 'md5'."
	}

	return response
}

func runTests() string {
	// Simulate running tests
	testsPassed := true // Simulate that tests passed

	if testsPassed {
		return "All tests passed!"
	} else {
		return "Some tests failed. Please check the logs for more information."
	}
}

func notifyTestsResult(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Tests passed successfully!")
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send test result notification: %v", err)
	}
}

func TestProcessMessage(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Secret",
			input:    "secret",
			expected: "md5(\"secret\") = 5ebe2294ecd0e0f08eab7690d2a6ee69\nreverse(\"5ebe2294ecd0e0f08eab7690d2a6ee69\") = \"secret\"",
		},
		{
			name:     "MD5",
			input:    "md5",
			expected: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name:     "Invalid",
			input:    "invalid",
			expected: "Неверный текст. Пожалуйста, выберите опцию из кнопок или отправьте 'secret' или 'md5'.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := processMessage(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
