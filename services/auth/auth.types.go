package auth

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

type User struct {
	ID       int    `gorm:"primaryKey; autoIncrement" json:"id" `
	Email    string `gorm:"unique_index" json:"email"`
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Role     Role   `json:"role"`
	Token    Token  `gorm:"foreignKey:UserID"`
}

type Token struct {
	ID              int    `gorm:"primaryKey; autoIncrement" json:"id"`
	UserID          int    `gorm:"not null" json:"user_id"`
	JwtAccessToken  string `gorm:"not null" json:"jwt_access_token"`
	JwtRefreshToken string `gorm:"not null" json:"jwt_refresh_token"`
	Expiry          int    `json:"expiry"`
}