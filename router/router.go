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

	auth := router.Group("/users")
	auth.Post("/", handlers.CreateUser)
	auth.Delete("/", middleware.AuthenticatedMiddleware, handlers.DeleteUser)
	auth.Put("/", middleware.AuthenticatedMiddleware, handlers.ChangePassword)
	auth.Get("/", middleware.AuthenticatedMiddleware, handlers.GetUserInfo)
	auth.Post("/login", handlers.Login)
	auth.Delete("/logout", handlers.Logout)

	product := router.Group("/products", middleware.AuthenticatedMiddleware)
	product.Post("/", handlers.CreateProduct)
	product.Get("/", handlers.GetProducts)
	product.Delete("/:id", handlers.DeleteProduct)
	product.Get("/:id", handlers.GetProductById)
	product.Put("/:id", handlers.UpdateProduct)

	router.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

}
