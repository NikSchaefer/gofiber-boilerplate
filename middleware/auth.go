package middleware

import (
	"encoding/json"

	"github.com/NikSchaefer/go-fiber/handlers"
	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
)

func Authenticated(c *fiber.Ctx) error {
	type AuthRequest struct {
		Sessionid string `json:"sessionid"`
	}
	data := new(AuthRequest)
	body := c.Body()
	if err := json.Unmarshal(body, &data); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON Sent",
		})
	}
	id, err := guuid.Parse(data.Sessionid)
	if err != nil {
		return c.Status(400).SendString("Invalid Session Format")
	}
	user, status := handlers.GetUser(id)
	if status != 0 {
		return c.SendStatus(status)
	}
	c.Locals("user", user)
	return c.Next()
}
