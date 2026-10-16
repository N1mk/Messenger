package domain

type User struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Username   string `json:"username" validate:"required,max=20" gorm:"notnull"`
	Email      string `json:"email" validate:"required,email" gorm:"notnull;uniqueIndex"`
	Password   string `json:"password" validate:"required,min=6,max=20" gorm:"notnull"`
	Chats      []Chat `json:"chats" gorm:"many2many:user_chats;"`
	OwnedChats []Chat `gorm:"foreignKey:OwnerID"`
}

type UserWithoutSlices struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" validate:"required,max=20" gorm:"notnull"`
	Email    string `json:"email" validate:"required,email" gorm:"notnull"`
	Password string `json:"password" validate:"required,min:6,max:20" gorm:"notnull"`
}
