package users

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserService interface {
    AddUser(context *fiber.Ctx) error
    GetAllUsers(context *fiber.Ctx) error
    HelloUser(context *fiber.Ctx) error
}

type DataUserService struct {
    DB *gorm.DB
}

func (method *DataUserService) AddUser(c *fiber.Ctx) error {
    user := new(User)
    if err := c.BodyParser(user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }
    if method.DB == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection is nil"})
    }
    result := method.DB.Create(&user)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(user)
}

func (method *DataUserService) HelloUser(c *fiber.Ctx) error {
    return c.SendString("Hello, users!")
}

func (method *DataUserService) GetAllUsers(c *fiber.Ctx) error {
    var users []User
    if method.DB == nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection is nil"})
    }
    result := method.DB.Find(&users)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
    }
    return c.Status(fiber.StatusOK).JSON(users)
}

func MethodUserService(db *gorm.DB) UserService {
    return &DataUserService{DB: db}
}
