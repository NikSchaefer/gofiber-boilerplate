package router

import (
	"github.com/NikSchaefer/go-fiber/handlers"
	"github.com/NikSchaefer/go-fiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func InitalizeRoutes(router *fiber.App) {

	router.Use(middleware.Security)

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello, World!")
	})

	router.Use(middleware.JsonMiddleware)

	auth := router.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/logout", handlers.Logout)
	auth.Post("/create", handlers.CreateUser)
	auth.Post("/delete", handlers.DeleteUser)
	auth.Post("/update", handlers.ChangePasswordRoute)
	auth.Post("/user", handlers.GetUserInfo)

	product := router.Group("/product")
	
}
