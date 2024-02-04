package users

import (
	"root/services/consts"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserService interface {
    AddUser(context *fiber.Ctx) error
    GetAllUsers(context *fiber.Ctx) error
    GetUserByID(context *fiber.Ctx) error
    HelloUser(context *fiber.Ctx) error
    ChangeUserByID(context *fiber.Ctx) error
    DeleteUserByID(context *fiber.Ctx) error
}

type DataUserService struct {
    DB *gorm.DB
}

func (method *DataUserService) AddUser(c *fiber.Ctx) error {
    user := new(User)
    if err := c.BodyParser(user); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(consts.BadReuest + err.Error())
    }
    if method.DB == nil {
		return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorConnectToDb +  method.DB.Error.Error())
    }
    result := method.DB.Create(&user)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorQueryDb + result.Error.Error())
    }

    return c.Status(fiber.StatusOK).JSON(user)
}

func (method *DataUserService) HelloUser(c *fiber.Ctx) error {
    return c.SendString("Hello, users!")
}

func (method *DataUserService) GetAllUsers(c *fiber.Ctx) error {
    var users []User
    if method.DB == nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorNil +  method.DB.Error.Error())
    }
    result := method.DB.Find(&users)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(consts.ErrorQueryDb + result.Error.Error())
    }
    return c.Status(fiber.StatusOK).JSON(users)
}

func (method *DataUserService) GetUserByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var user User
	result := method.DB.Where("id = ?", id).First(&user)

	if id == "" {
		return c.Status(400).SendString(consts.ErrorNil)
	}
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return c.Status(404).SendString(consts.ErrorNil)
        }
        return c.Status(500).SendString(consts.ErrorQueryDb + result.Error.Error())
    }
	return c.Status(200).JSON(user)
}

func (method *DataUserService) ChangeUserByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var user User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(400).SendString(consts.BadReuest + err.Error())
    }
    result := method.DB.Model(&user).Where("id = ?", id).Updates(&user)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return c.Status(404).SendString(consts.ErrorNil)
        }
        return c.Status(500).SendString(consts.ErrorQueryDb + result.Error.Error())
    }
    return c.Status(200).JSON(user)
}

func (method *DataUserService) DeleteUserByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var user User
    result := method.DB.Where("id = ?", id).Delete(&user)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return c.Status(404).SendString(consts.ErrorNil)
        }
        return c.Status(500).SendString(consts.ErrorQueryDb + result.Error.Error())
    }
    return c.Status(200).SendString("Пользователь успешно удален")
}


func MethodUserService(db *gorm.DB) UserService {
    return &DataUserService{DB: db}
}
