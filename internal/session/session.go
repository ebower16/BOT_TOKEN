package session

import "time"

// UserSession представляет сессию пользователя
type UserSession struct {
	UserID       int64     // Уникальный идентификатор пользователя
	Attempts     int       // Количество попыток входа
	Blocked      bool      // Флаг, указывающий, заблокирован ли пользователь
	BlockedUntil time.Time // Время, до которого пользователь заблокирован
	Username     string    // Имя пользователя
	Password     string    // Пароль (не хранится в сессии)
}

// NewUserSession создает новую сессию пользователя
func NewUserSession(userID int64, username string) *UserSession {
	return &UserSession{
		UserID:   userID,
		Username: username,
		Attempts: 0,
		Blocked:  false,
	}
}
