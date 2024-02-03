package users

type User struct {
	Id       int    `json:"id" gorm:"uniqueIndex;not null"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password" gorm:"not null"`
	Token    string `json:"token"`
}
