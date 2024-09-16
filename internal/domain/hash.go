package domain

import (
	"time"
)

type Hash struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Time        time.Time `json:"time" gorm:"autoCreateTime"`
	Description string    `json:"description" gorm:"size:255"`
}
