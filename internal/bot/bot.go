package bot

import (
	"botus/internal/database"
	"botus/internal/session"
	"context"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const maxAttempts = 3 // Максимальное количество попыток входа

// RunBot инициализирует бота и начинает прослушивание обновлений
func RunBot(botToken string) {
	db := database.InitDatabase()
	defer db.Close()

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := getUpdatesChan(bot)
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := setupSignalHandler(cancel)

	userSessions := make(map[int64]*session.UserSession)

	runBot(updates, ctx, signalChan, db, bot, userSessions)
}

// getUpdatesChan получает обновления от бота
func getUpdatesChan(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 10})
	if err != nil {
		log.Fatalf("Failed to get updates: %v", err)
	}
	return updates
}

// setupSignalHandler настраивает обработчик сигналов для корректного завершения
func setupSignalHandler(cancel context.CancelFunc) chan os.Signal {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Received shutdown signal, shutting down gracefully...")
		cancel()
	}()

	return signalChan
}

// runBot обрабатывает входящие обновления
func runBot(updates tgbotapi.UpdatesChannel, ctx context.Context, signalChan chan os.Signal, db *sql.DB, bot *tgbotapi.BotAPI, userSessions map[int64]*session.UserSession) {
	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			userID := update.Message.Chat.ID

			// Проверяем, существует ли сессия для пользователя
			if _, exists := userSessions[userID]; !exists {
				userSessions[userID] = session.NewUserSession(userID, "")
			}

			// Используем функции для обработки
			if session, exists := userSessions[userID]; exists && session.Blocked {
				handleBlockedUser(session, &msg) // Передаем указатель на msg
			} else {
				handleUserInput(update, &msg, userSessions, userID, db) // Передаем указатель на msg
			}

			if msg.Text != "" {
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Failed to send message: %v", err)
				}
			}

		case <-ctx.Done():
			log.Println("Shutting down the bot...")
			return
		}
	}
}
