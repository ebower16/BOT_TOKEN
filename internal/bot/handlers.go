package bot

import (
	"botus/internal/database"
	"botus/internal/session"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// handleBlockedUser обрабатывает сообщения от заблокированных пользователей
func handleBlockedUser(session *session.UserSession, msg *tgbotapi.MessageConfig) {
	msg.Text = fmt.Sprintf("🔒 Вы заблокированы. Пожалуйста, подождите %d секунд, пока блокировка не будет снята.", session.BlockedUntil.Sub(time.Now()).Seconds())
}

// handleUserInput обрабатывает входящие сообщения от пользователей
func handleUserInput(update tgbotapi.Update, msg *tgbotapi.MessageConfig, userSessions map[int64]*session.UserSession, userID int64, db *sql.DB) {
	// Проверяем, существует ли сессия для пользователя
	if _, exists := userSessions[userID]; !exists {
		userSessions[userID] = session.NewUserSession(userID, "")
	}

	session := userSessions[userID]

	// Если пользователь заблокирован, обрабатываем блокировку
	if session.IsBlocked() {
		handleBlockedUser(session, msg)
		return
	}

	// Создаем клавиатуру с кнопками
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/register"),
			tgbotapi.NewKeyboardButton("/login"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/list_users"),
			tgbotapi.NewKeyboardButton("/calculate_parallel"),
		),
	)

	// Устанавливаем клавиатуру для сообщения
	msg.ReplyMarkup = keyboard

	// Обработка команд
	if strings.HasPrefix(update.Message.Text, "/register") {
		err := registerUser(update.Message.Text, db, update.Message.Chat.ID)
		if err != nil {
			msg.Text = "Ошибка регистрации: " + err.Error()
		} else {
			msg.Text = "Регистрация прошла успешно! Теперь вы можете войти, используя команду /login."
			notifyAdmin(fmt.Sprintf("Новый пользователь зарегистрировался: %s", strings.Fields(update.Message.Text)[1]))
		}
	} else if strings.HasPrefix(update.Message.Text, "/login") {
		err := loginUser(update.Message.Text, db, session)
		if err != nil {
			msg.Text = "Ошибка входа: " + err.Error()
		} else {
			msg.Text = "Вы успешно вошли в систему!"
		}
	} else if strings.HasPrefix(update.Message.Text, "/list_users") {
		users, err := database.GetAllUsers(db)
		if err != nil {
			msg.Text = "Ошибка получения списка пользователей: " + err.Error()
		} else {
			if len(users) == 0 {
				msg.Text = "Нет зарегистрированных пользователей."
			} else {
				msg.Text = "Зарегистрированные пользователи:\n" + strings.Join(users, "\n")
			}
		}
	} else if strings.HasPrefix(update.Message.Text, "/calculate_parallel") {
		n := 10000000 // Например, 10 миллионов
		start := time.Now()
		parallelResult := ParallelSumOfSquares(n)
		duration := time.Since(start)
		msg.Text = fmt.Sprintf("Параллельный результат: %d, время: %v", parallelResult, duration)
	} else {
		msg.Text = "Пожалуйста, выберите команду из меню."
	}
}

// registerUser регистрирует пользователя в базе данных
func registerUser(input string, db *sql.DB, chatID int64) error {
	// Разделяем входные данные на логин и пароль
	parts := strings.Fields(input)
	if len(parts) != 3 { // /register <логин> <пароль>
		return fmt.Errorf("неверный формат ввода")
	}
	username := parts[1]
	password := parts[2]

	// Регистрируем пользователя
	if !database.RegisterUser(username, password, db, chatID) {
		return fmt.Errorf("пользователь с таким именем уже существует")
	}
	return nil
}

// loginUser проверяет логин и пароль пользователя
func loginUser(input string, db *sql.DB, session *session.UserSession) error {
	// Разделяем входные данные на логин и пароль
	parts := strings.Fields(input)
	if len(parts) != 3 { // /login <логин> <пароль>
		return fmt.Errorf("неверный формат ввода")
	}
	username := parts[1]
	password := parts[2]

	// Проверяем логин и пароль
	if !database.CheckUserCredentials(username, password, db) {
		session.BlockUser(10 * time.Second)
		return fmt.Errorf("неверный логин или пароль")
	}

	session.Username = username
	session.ResetAttempts()
	return nil
}

// notifyAdmin отправляет уведомление администратору о регистрации нового пользователя
func notifyAdmin(message string) {
	botToken := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Printf("Failed to create bot for notification: %v", err)
		return
	}

	msg := tgbotapi.NewMessageToChannel(adminChatID, message)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send notification to admin: %v", err)
	}
}

// ParallelSumOfSquares вычисляет сумму квадратов чисел от 1 до n параллельно
func ParallelSumOfSquares(n int) int {
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunkSize := n / numCPU
	results := make([]int, numCPU)

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * chunkSize
			end := start + chunkSize
			if i == numCPU-1 {
				end = n // последний поток обрабатывает остаток
			}
			sum := 0
			for j := start + 1; j <= end; j++ {
				sum += j * j
			}
			results[i] = sum
		}(i)
	}

	wg.Wait()

	total := 0
	for _, result := range results {
		total += result
	}
	return total
}
