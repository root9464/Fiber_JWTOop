package users

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func AddUserRoutes(routes fiber.Router, gnu UserService, db *gorm.DB) {

    routes.Get("/users/hello", gnu.HelloUser)
    routes.Get("/users/allusers", gnu.GetAllUsers)
    routes.Get("/users/:id/details", gnu.GetUserByID)

    // routes.Use(UserExists(db))
    routes.Post("/users/adduser", gnu.AddUser)
    routes.Put("/users/:id", gnu.ChangeUserByID)
    routes.Delete("/users/:id/update", gnu.DeleteUserByID)
}