package router

import "github.com/gofiber/fiber/v2"

func InitalizeRoutes(router *fiber.App) {
	router.Get("/")
}