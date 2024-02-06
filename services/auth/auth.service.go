package auth

import (
	"log"
	"root/services/consts"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(c *fiber.Ctx) error
    Login(c *fiber.Ctx) error
    Hello(c *fiber.Ctx) error
    AccessTokenUpdate(c *fiber.Ctx) error
}

type DataAuthService struct {
	DB *gorm.DB
}

func (method *DataAuthService) Register(c *fiber.Ctx) error {
    user := new(User)
    if err := c.BodyParser(user); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(consts.BadReuest + err.Error())
    }

    password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorGenHash + err.Error())
    }else{
        user.Password = string(password)
    }

    if method.DB == nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorConnectToDb + method.DB.Error.Error())
    }
    result := method.DB.Create(&user)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorQueryDb + result.Error.Error())
    }
    
    access, err := method.CreateAccessToken(user.ID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorCreateAccessToken + err.Error())
    }

    refresh, err := method.CreateRefreshToken(user.ID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorCreateRefreshToken + err.Error())
    }
    
    access = "Bearer" + access


    token := &Token{
        UserID:         user.ID,
        JwtAccessToken: access,
        JwtRefreshToken:   refresh,
        Expiry:         int(time.Now().Add(time.Hour * 24 * 7).Unix()),
    }


    method.DB.Create(&token)
    
    user.Token = *token 
    
    return c.Status(fiber.StatusOK).JSON(user)
}

func (method *DataAuthService) Login(c *fiber.Ctx) error {
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

    var user User
    if err := method.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return c.Status(fiber.StatusUnauthorized).SendString("данного пользователя нету в бд" + err.Error())
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return c.Status(fiber.StatusUnauthorized).SendString(user.Email + "\n" + email + "\n" + password + "\n" + user.Password)
    }

    var token Token
    if err := method.DB.Where("user_id = ?", user.ID).First(&token).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorQueryDb + err.Error())
    }else {
        user.Token = token
    }


    tokenData, err := jwt.Parse(token.JwtAccessToken, func(token *jwt.Token) (interface{}, error) {
        return []byte("your_access_secret_key"), nil
    }); if err != nil {
        if _, ok := tokenData.Claims.(jwt.MapClaims); !ok || !tokenData.Valid && time.Now().Unix() > int64(tokenData.Claims.(jwt.MapClaims)["exp"].(float64)){
            newAccessToken, newRefreshToken, err := method.CreateTokens(user.ID)
            if err != nil {
                return c.Status(fiber.StatusInternalServerError).SendString("fff" + err.Error())
            }
    
            token.JwtAccessToken = newAccessToken
            token.JwtRefreshToken = newRefreshToken
        }else {
            log.Println("токены не устарели")
        }
    }

    return c.Status(fiber.StatusOK).JSON(user)
}




func (method *DataAuthService) Hello(c *fiber.Ctx) error {
    return c.SendString("Hello, auth!")
}
func MethodAuthService(db *gorm.DB) AuthService {
	return &DataAuthService{DB: db}
}