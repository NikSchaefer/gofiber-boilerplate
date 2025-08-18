package middleware

import "github.com/gofiber/fiber/v2"

func Json(c *fiber.Ctx) error {
	c.Accepts("application/json")
	return c.Next()
}
