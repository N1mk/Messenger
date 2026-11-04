package repository

type RefreshToken struct {
	TokenHash []byte `gorm:"column:token_hash"`
	UserID    int    `gorm:"column:user_id"`
}
