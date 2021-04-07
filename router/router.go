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
	auth.Post("/delete", middleware.AuthenticatedMiddleware, handlers.DeleteUser)
	auth.Post("/update", middleware.AuthenticatedMiddleware, handlers.ChangePassword)
	auth.Post("/user", middleware.AuthenticatedMiddleware, handlers.GetUserInfo)

	product := router.Group("/product", middleware.AuthenticatedMiddleware)
	product.Post("/create", handlers.CreateProduct)
	product.Post("/delete", handlers.DeleteProduct)
	product.Post("/products", handlers.GetProduct)
	product.Post("/id/:id", handlers.GetProductById)
	product.Post("/update", handlers.UpdateProduct)

	router.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

}
