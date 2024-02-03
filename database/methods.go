package database

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func (method *Data) Migration(data interface{}) error {
	db := method.DB
	err := db.AutoMigrate(data); if err != nil {
		log.Fatal("ошибка: не удалось мигрировать. \n", err)
	}
	return nil
}

func (method *Data) CreateTestData( c *fiber.Ctx) error {
	dto := method.DB
	if err := c.BodyParser(dto); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	method.DB.Create(&dto)
	
	return c.Status(200).JSON(dto)
}
func MethodDB() Database {
	return &Data{}
}

