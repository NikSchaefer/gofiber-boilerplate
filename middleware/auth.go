package middleware

import (
	"github.com/NikSchaefer/go-fiber/handlers"
	"github.com/NikSchaefer/go-fiber/model"
	"github.com/gofiber/fiber/v2"
)

func AuthenticatedMiddleware(c *fiber.Ctx) error {
	json := new(model.Session)
	if err := c.BodyParser(json); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	user, status := handlers.GetUser(json.Sessionid)
	if status != 0 {
		return c.SendStatus(status)
	}
	c.Locals("user", user)
	return c.Next()
}
