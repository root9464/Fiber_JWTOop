package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddAuthrRoutes(routes fiber.Router, gnu AuthService, db *gorm.DB) {
	routes.Post("/auth/register", gnu.Register)
	routes.Post("/auth/login", gnu.Login)
	routes.Get("/auth/hello", gnu.Hello)
	routes.Post("/auth/refresh", gnu.AccessTokenUpdate)
}