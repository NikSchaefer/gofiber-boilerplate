package middleware

import (
	"github.com/NikSchaefer/go-fiber/handlers"
	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
)

func Authenticated(c *fiber.Ctx) error {
	type AuthRequest struct {
		Sessionid guuid.UUID `json:"sessionid"`
	}
	json := new(AuthRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(422).JSON(fiber.Map{
			"err":  err,
			"body": c.Body(),
		})
	}
	user, status := handlers.GetUser(json.Sessionid)
	if status != 0 {
		return c.SendStatus(status)
	}
	c.Locals("user", user)
	return c.Next()
}
