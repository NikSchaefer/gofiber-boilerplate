package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/NikSchaefer/go-fiber/internal/services"
)

// Authenticated middleware verifies that a user has a valid session and is verified
func Authenticated(c *fiber.Ctx) error {
	// check if user is already in locals (Authenticated via other methods)
	if c.Locals("user") != nil {
		return c.Next()
	}

	user, err := services.ValidateSession(c.Context(), c.Cookies("session"))
	if err != nil {
		return err
	}

	c.Locals("auth_type", "session")
	c.Locals("user", user)
	return c.Next()
}
