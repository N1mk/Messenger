package domain

type Chat struct {
	ID       int       `json:"id" gorm:"primaryKey" validate:"required,min=0"`
	OwnerID  int       `json:"ownerID" validate:"required,min=0"`
	Name     string    `json:"name"`
	Users    []User    `json:"users" gorm:"many2many:user_chats;constraint:OnDelete:CASCADE;"`
	Messages []Message `json:"messages" gorm:"foreignKey:ChatID"`
}

type ChatWithoutSlices struct {
	ID   int    `json:"id" gorm:"primaryKey" validate:"required,min=0"`
	Name string `json:"name"`
	//Owner    User      `json:"owner" gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
