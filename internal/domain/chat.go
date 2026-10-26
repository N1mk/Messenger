package domain

type Chat struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	OwnerID  uint      `json:"ownerID"`
	Name     string    `json:"name"`
	Users    []User    `json:"users" gorm:"many2many:user_chats;constraint:OnDelete:CASCADE;"`
	Messages []Message `json:"messages" gorm:"foreignKey:ChatID"`
}

type ChatWithoutSlices struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	//Owner    User      `json:"owner" gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
