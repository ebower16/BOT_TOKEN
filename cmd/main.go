package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

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

	for update := range updates {
		if update.Message == nil {
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message.Text == "/start" {
			msg.Text = "Welcome! Please choose an option:"
			msg.ReplyMarkup = getKeyboard()
		} else if update.Message.Text == "Hash secret" {
			msg.Text = processMessage("secret")
		} else if update.Message.Text == "md5" {
			msg.Text = processMessage(update.Message.Text)
		} else if update.Message.Text == "secret" {
			msg.Text = processMessage(update.Message.Text)
		} else if update.Message.Text == "Show Image" {
			msg.Text = "Here is an image for you:"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // Убираем клавиатуру
			imageMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "path/to/your/image.jpg")
			if _, err := bot.Send(imageMsg); err != nil {
				log.Printf("Failed to send image: %v", err)
			}
		} else {
			msg.Text = "Please choose a valid option from the buttons, send 'md5' to get its MD5 hash, or send 'secret' to get its MD5 hash and reverse."
		}

		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}

func getKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Hash secret"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("md5"),
			tgbotapi.NewKeyboardButton("secret"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Show Image"),
			tgbotapi.NewKeyboardButton("Help"),
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
	} else if text == "md5" {
		response = hashStr
	} else {
		response = "Неверный текст. Пожалуйста, выберите опцию из кнопок или отправьте 'secret' или 'md5'."
	}

	return response
}
