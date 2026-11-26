package domain

type ChatInvite struct {
	ChatID int `json:"chat_id" gorm:"column:chat_id" validate:"required"`
	UserID int `json:"user_id" gorm:"column:user_id" validate:"required"`
}
