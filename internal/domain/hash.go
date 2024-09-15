package domain

import (
	"time"
)

// Hash представляет собой доменную модель для хранения хеша и его значения.
type Hash struct {
	ID          int       // Уникальный идентификатор
	Time        time.Time // Время создания
	Description string    // Описание или значение, связанное с хешем
}
