package users

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserExists(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id") 
		var user User
		result := db.Where("id = ?", id).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return c.Status(404).SendString("нету пользователя")
			}
			return c.Status(403).SendString("нету доступа")
		}
		c.Locals("user", user)
		return c.Next()
	}
}