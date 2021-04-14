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

	users := router.Group("/users")
	users.Post("/", handlers.CreateUser)
	users.Delete("/", middleware.AuthenticatedMiddleware, handlers.DeleteUser)
	users.Put("/", middleware.AuthenticatedMiddleware, handlers.ChangePassword)
	users.Get("/", middleware.AuthenticatedMiddleware, handlers.GetUserInfo)
	users.Post("/login", handlers.Login)
	users.Delete("/logout", handlers.Logout)

	products := router.Group("/products", middleware.AuthenticatedMiddleware)
	products.Post("/", handlers.CreateProduct)
	products.Get("/", handlers.GetProducts)
	products.Delete("/:id", handlers.DeleteProduct)
	products.Get("/:id", handlers.GetProductById)
	products.Put("/:id", handlers.UpdateProduct)

	router.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

}
