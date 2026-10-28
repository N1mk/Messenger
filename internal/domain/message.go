package domain

import (
	"time"
)

type Message struct {
	ID        int       `gorm:"primaryKey" json:"id" validate:"required,min=0"`
	ChatID    int       `json:"chat_id" validate:"required,min=0"`
	UserID    int       `json:"user_id" validate:"required,min=0"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
