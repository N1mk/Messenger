package domain

type RefreshToken struct {
	TokenHash []byte `json:"token_hash" gorm:"column:token_hash" validate:"required"`
	UserID    int    `json:"user_id" gorm:"column:user_id" validate:"required"`
}
