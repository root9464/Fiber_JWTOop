package users

import "github.com/gofiber/fiber/v2"


func AddUserRoutes(routes fiber.Router, gnu UserService) {
    routes.Get("/users/hello", gnu.HelloUser)
    routes.Post("/users/adduser", gnu.AddUser)
    routes.Get("/users/allusers", gnu.GetAllUsers)
}