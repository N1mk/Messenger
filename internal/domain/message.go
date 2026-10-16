package domain

import (
	"time"
)

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ChatID    uint      `json:"chat_id"`
	Text      string    `json:"text"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
