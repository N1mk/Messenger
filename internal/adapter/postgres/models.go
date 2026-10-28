package postgres

type RefreshToken struct {
	Token  string `gorm:"column:token"`
	UserID int    `gorm:"column:user_id"`
}
