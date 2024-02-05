package auth

import (
	"log"
	"root/services/consts"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(c *fiber.Ctx) error
    Login(c *fiber.Ctx) error
    
    Hello(c *fiber.Ctx) error
    
    AccessTokenUpdate(c *fiber.Ctx) error

    CreateAccessToken(userID int) (string, error)
    CreateRefreshToken(userID int) (string, error)
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

    // Создание токена обновления
    refresh, err := method.CreateRefreshToken(user.ID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorCreateRefreshToken + err.Error())
    }
    
    access = "Bearer " + access


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


    //! сдесь ошибка при обновлении токена
    
    log.Print(user)

    expiry := user.Token.Expiry
    log.Print(expiry)


    newExpiry := int(time.Now().Unix())
    log.Print(newExpiry)


    newToken := Token{
        UserID:         user.ID,
        JwtAccessToken: user.Token.JwtAccessToken,
        JwtRefreshToken: user.Token.JwtRefreshToken,
        Expiry:         newExpiry,
    }



    if expiry < newExpiry || user.Token == (Token{
        UserID:         0,
        JwtAccessToken: "",
        JwtRefreshToken: "",
        Expiry:         user.Token.Expiry,
    }) {
        access, err := method.CreateAccessToken(user.ID)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorCreateAccessToken + err.Error())
        }else{
            access = "Bearer " + access
            newToken.JwtAccessToken = access
        }

        refresh, err := method.CreateRefreshToken(user.ID)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorCreateRefreshToken + err.Error())
        }else{
            newToken.JwtRefreshToken = refresh
        }
      
        log.Print(newToken)
        newToken.Expiry = newExpiry
        method.DB.Model(&Token{}).Where("user_id = ?", user.ID).Update("jwt_access_token", newToken)

        newToken.Expiry = int(time.Now().Add(time.Hour * 24 * 7).Unix())
        user.Token = newToken
    }

    return c.Status(fiber.StatusOK).JSON(user)
}




func (method *DataAuthService) Hello(c *fiber.Ctx) error {
    return c.SendString("Hello, auth!")
}
func MethodAuthService(db *gorm.DB) AuthService {
	return &DataAuthService{DB: db}
}