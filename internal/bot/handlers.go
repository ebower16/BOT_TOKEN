package bot

import (
	"botus/internal/session"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// handleBlockedUser обрабатывает сообщения от заблокированных пользователей
func handleBlockedUser(session *session.UserSession, msg *tgbotapi.MessageConfig) {
	msg.Text = "🔒 Вы заблокированы. Пожалуйста, подождите, пока блокировка не будет снята."
}

// handleUserInput обрабатывает входящие сообщения от пользователей
func handleUserInput(update tgbotapi.Update, msg *tgbotapi.MessageConfig, userSessions map[int64]*session.UserSession, userID int64, db *sql.DB) {
	// Здесь добавьте логику обработки пользовательского ввода
	// Например, вы можете обновить сессию пользователя или выполнить какие-то действия
	msg.Text = "Ваше сообщение получено: " + update.Message.Text
}
