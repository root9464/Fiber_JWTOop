package root

import (
	"log"
	"root/services/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

// Функция для установки маршрутов
func AllRo(app *fiber.App, db *gorm.DB) {
    routes := app.Group("/api")
    routes.Use(SetDBState(db))
    userService := users.MethodUserService(db)
    users.AddUserRoutes(routes, userService)
}

// Функция для запуска сервера
func Root(db *gorm.DB) {
    app := fiber.New()

    // Добавление обработчиков запросов и маршрутов
    AllRo(app, db)

    // Запуск сервера
    app.Use(cors.New())
    app.Use(func(c *fiber.Ctx) error {
        return c.SendStatus(404) // => 404 "Not Found"
    })

    log.Fatal(app.Listen(":3000"))
}