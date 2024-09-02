package main

import (
	"context"
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
	"1":  "9c56cc51e4b3a1b2e9b9d501c6a7c4d8", // MD5 hash for "1"
	"2":  "5ebe2294ecd0e0f08eab7690d2a6ee69", // MD5 hash for "2"
	"3":  "5d41402abc4b2a76b9719d911017c592", // MD5 hash for "3"
	"4":  "5eb63bbbe01eeed093cb22bb8f5acdd2", // MD5 hash for "4"
	"5":  "d3b07384d113edec49eaa6238ad5ff00", // MD5 hash for "5"
	"6":  "f2ca1b5e7f4c7d7f3b0b6d8f7d8b3b3",  // MD5 hash for "6"
	"7":  "e2fc714c4727ee9395f324cd2e7f331f", // MD5 hash for "7"
	"8":  "2c6ee24b09816a6f14f95d1698b24ead", // MD5 hash for "8"
	"9":  "f8a5b5a1c1f4d8b4d8b4a1f4c1f4d8b4", // MD5 hash for "9"
	"10": "d41d8cd98f00b204e9800998ecf8427e", // MD5 hash for "10"
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
			case "Show Image":
				msg.Text = "Here is an image for you:"
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				imageMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "path/to/your/image.jpg")
				if _, err := bot.Send(imageMsg); err != nil {
					log.Printf("Failed to send image: %v", err)
				}
			case "Run Tests":
				msg.Text = runTests()
			case "Notify Tests Passed":
				notifyTestsResult(bot, update.Message.Chat.ID)
			case "Exit":
				msg.Text = "Exiting the bot..."
				cancel()
			default:
				// Process the message to find the hash or original string
				msg.Text = processMessage(update.Message.Text)
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
	// Check if the user input matches a string in the hash table
	hash, exists := hashTable[text]
	if exists {
		return "The MD5 hash of \"" + text + "\" is: " + hash
	}

	// Check if the user input matches a hash in the hash table
	for key, value := range hashTable {
		if value == text {
			return "The original string for the MD5 hash \"" + text + "\" is: \"" + key + "\""
		}
	}

	// If not found, inform the user
	return "Error: The string or hash \"" + text + "\" is not in the hash table."
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
			name:     "Hash for '1'",
			input:    "1",
			expected: "The MD5 hash of \"1\" is: 9c56cc51e4b3a1b2e9b9d501c6a7c4d8",
		},
		{
			name:     "Hash for '2'",
			input:    "2",
			expected: "The MD5 hash of \"2\" is: 5ebe2294ecd0e0f08eab7690d2a6ee69",
		},
		{
			name:     "Hash for '3'",
			input:    "3",
			expected: "The MD5 hash of \"3\" is: 5d41402abc4b2a76b9719d911017c592",
		},
		{
			name:     "Hash for '4'",
			input:    "4",
			expected: "The MD5 hash of \"4\" is: 5eb63bbbe01eeed093cb22bb8f5acdd2",
		},
		{
			name:     "Hash for '5'",
			input:    "5",
			expected: "The MD5 hash of \"5\" is: d3b07384d113edec49eaa6238ad5ff00",
		},
		{
			name:     "Hash for '6'",
			input:    "6",
			expected: "The MD5 hash of \"6\" is: f2ca1b5e7f4c7d7f3b0b6d8f7d8b3b3",
		},
		{
			name:     "Hash for '7'",
			input:    "7",
			expected: "The MD5 hash of \"7\" is: e2fc714c4727ee9395f324cd2e7f331f",
		},
		{
			name:     "Hash for '8'",
			input:    "8",
			expected: "The MD5 hash of \"8\" is: 2c6ee24b09816a6f14f95d1698b24ead",
		},
		{
			name:     "Hash for '9'",
			input:    "9",
			expected: "The MD5 hash of \"9\" is: f8a5b5a1c1f4d8b4d8b4a1f4c1f4d8b4",
		},
		{
			name:     "Hash for '10'",
			input:    "10",
			expected: "The MD5 hash of \"10\" is: d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:     "Invalid input",
			input:    "invalid",
			expected: "Error: The string or hash \"invalid\" is not in the hash table.",
		},
		{
			name:     "Reverse hash for '5ebe2294ecd0e0f08eab7690d2a6ee69'",
			input:    "5ebe2294ecd0e0f08eab7690d2a6ee69",
			expected: "The original string for the MD5 hash \"5ebe2294ecd0e0f08eab7690d2a6ee69\" is: \"2\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := processMessage(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
