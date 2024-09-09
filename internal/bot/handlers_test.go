package bot

import (
	"botus/internal/database"
	"botus/internal/session"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestHandleBlockedUser(t *testing.T) {
	session := &session.UserSession{
		BlockedUntil: time.Now().Add(time.Minute),
	}
	msg := &tgbotapi.MessageConfig{}

	handleBlockedUser(session, msg)

	if msg.Text == "" {
		t.Errorf("Expected message text, got empty string")
	}
}

func TestRegisterUser(t *testing.T) {
	db := InitDatabase()
	defer db.Close()

	msg := &tgbotapi.MessageConfig{}
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 123456},
			Text: "/register testuser testpassword",
		},
	}

	userSessions := make(map[int64]*session.UserSession)
	handleUserInput(update, msg, userSessions, update.Message.Chat.ID, db, adminChatID)

	if msg.Text == "" {
		t.Errorf("Expected message text for registration, got empty string")
	}
}

func TestLoginUser(t *testing.T) {
	db := InitDatabase()
	defer db.Close()

	database.RegisterUser("testuser", "testpassword", db, 123456)

	session := &session.UserSession{}
	msg := &tgbotapi.MessageConfig{}
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 123456},
			Text: "/login testuser testpassword",
		},
	}

	userSessions := make(map[int64]*session.UserSession)
	handleUserInput(update, msg, userSessions, update.Message.Chat.ID, db, adminChatID)

	if msg.Text == "" {
		t.Errorf("Expected message text for login, got empty string")
	}
}
