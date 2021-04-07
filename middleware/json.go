package middleware

import "github.com/gofiber/fiber/v2"

func JsonMiddleware(c *fiber.Ctx) error {
	c.Accepts("application/json")
	return c.Next()
}
