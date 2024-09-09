package session

import "time"

// UserSession представляет сессию пользователя
type UserSession struct {
	UserID       int64     // Уникальный идентификатор пользователя
	Attempts     int       // Количество попыток входа
	Blocked      bool      // Флаг, указывающий, заблокирован ли пользователь
	BlockedUntil time.Time // Время, до которого пользователь заблокирован
	Username     string    // Имя пользователя
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

// IsBlocked проверяет, заблокирован ли в данный момент пользователь
func (s *UserSession) IsBlocked() bool {
	return s.Blocked && time.Now().Before(s.BlockedUntil)
}

// BlockUser помечает пользователя как заблокированного до указанного времени
func (s *UserSession) BlockUser(duration time.Duration) {
	s.Blocked = true
	s.BlockedUntil = time.Now().Add(duration)
}

// UnblockUser разблокирует пользователя
func (s *UserSession) UnblockUser() {
	s.Blocked = false
	s.BlockedUntil = time.Time{} // Сбросить до нулевого времени
}

// IncrementAttempts увеличивает количество попыток входа
func (s *UserSession) IncrementAttempts() {
	s.Attempts++
}

// ResetAttempts сбрасывает количество попыток входа
func (s *UserSession) ResetAttempts() {
	s.Attempts = 0
}
