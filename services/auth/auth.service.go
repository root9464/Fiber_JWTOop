package auth

import (
	"root/services/consts"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(c *fiber.Ctx) error
    Login(c *fiber.Ctx) error
    Hello(c *fiber.Ctx) error
}

type DataAuthService struct {
	DB *gorm.DB
}

func (method *DataAuthService) Register(c *fiber.Ctx) error {
    user := new(User)
    if err := c.BodyParser(user); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(consts.BadReuest + err.Error())
    }
    if method.DB == nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorConnectToDb + method.DB.Error.Error())
    }
    result := method.DB.Create(&user)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorQueryDb + result.Error.Error())
    }
    
    // Create JWT tokens
    accessClaims := jwt.MapClaims{
        "user_id": user.ID,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    }
    
    refreshClaims := jwt.MapClaims{
        "user_id": user.ID,
        "exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
    }
    
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    
    access, _ := accessToken.SignedString([]byte("your_access_secret_key"))
    refresh, _ := refreshToken.SignedString([]byte("your_refresh_secret_key"))
    
    // Store tokens in the Tokens table
    token := &Token{
        UserID:         user.ID,
        JwtAccessToken: access,
        RefreshToken:   refresh,
        Expiry:         int(time.Now().Add(time.Hour * 24 * 7).Unix()),
    }
    method.DB.Create(&token)
    
    user.Token = *token // Update the user's token field
    
    return c.Status(fiber.StatusOK).JSON(user)
}

func (method *DataAuthService) Login(c *fiber.Ctx) error {
    // Извлечь электронную почту и пароль из запроса
    type LoginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    var request LoginRequest
    if err := c.BodyParser(&request); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(consts.BadReuest + err.Error())
    }

    email := request.Email
    password := request.Password
    
    // Найти пользователя по предоставленной электронной почте
    var user User
    if err := method.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or 1")
    }
    
    // Проверить пароль
    // if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
    //     return c.Status(fiber.StatusUnauthorized).SendString(user.Email + "\n" + email + "\n" + password + "\n" + user.Password)
    // }
    if user.Password != password {
        return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or 2")
    }

    // Создать и сохранить JWT-токены
    accessClaims := jwt.MapClaims{
        "user_id": user.ID,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    }
    refreshClaims := jwt.MapClaims{
        "user_id": user.ID,
        "exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
    }
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    access, _ := accessToken.SignedString([]byte("your_access_secret_key"))
    refresh, _ := refreshToken.SignedString([]byte("your_refresh_secret_key"))

    token := &Token{
        UserID: user.ID,
        JwtAccessToken: access,
        RefreshToken: refresh,
        Expiry: int(time.Now().Add(time.Hour * 24 * 7).Unix()),
    }
    if err := method.DB.Create(&token).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to create tokens")
    }

    user.Token = *token // Обновить поле токена пользователя
    return c.Status(fiber.StatusOK).JSON(user)
}

func (method *DataAuthService) Hello(c *fiber.Ctx) error {
    return c.SendString("Hello, auth!")
}
func MethodAuthService(db *gorm.DB) AuthService {
	return &DataAuthService{DB: db}
}