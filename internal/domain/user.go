package domain

type User struct {
	ID       int    `json:"id" gorm:"primaryKey" validate:"required,min=0"`
	Username string `json:"username" validate:"required,max=20" gorm:"notnull;uniqueIndex"`
	Email    string `json:"email" validate:"required,email" gorm:"notnull;uniqueIndex"`
	Password string `json:"password" validate:"required,min=6,max=20" gorm:"notnull"`
	Chats    []Chat `json:"chats" gorm:"many2many:user_chats;"`
}

type UserWithoutSlices struct {
	ID       int    `json:"id" gorm:"primaryKey" validate:"required,min=0"`
	Username string `json:"username" validate:"required,max=20" gorm:"notnull"`
	Email    string `json:"email" validate:"required,email" gorm:"notnull"`
	Password string `json:"password" validate:"required,min:6,max:20" gorm:"notnull"`
}
