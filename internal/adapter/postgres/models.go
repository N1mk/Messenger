package postgres

type RefreshToken struct {
	Token  string `gorm:"column:token"`
	UserID uint   `gorm:"column:user_id"`
}
