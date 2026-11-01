package domain

type User struct {
	ID           int    `json:"id" gorm:"primaryKey" validate:"required,min=0"`
	Username     string `json:"username" validate:"required,max=64" gorm:"notnull;uniqueIndex"`
	Email        string `json:"email" validate:"required,email" gorm:"notnull;uniqueIndex"`
	PasswordHash []byte `json:"password" gorm:"notnull"`
	Chats        []Chat `json:"chats" gorm:"many2many:user_chats;"`
}

// Переработать
type UserWithoutSlices struct {
	ID       int    `json:"id" gorm:"primaryKey" validate:"required,min=0"`
	Username string `json:"username" validate:"required,max=64" gorm:"notnull;uniqueIndex"`
	Email    string `json:"email" validate:"required,email" gorm:"notnull;uniqueIndex"`
	//PasswordHash []byte `json:"password" gorm:"notnull"`
}
