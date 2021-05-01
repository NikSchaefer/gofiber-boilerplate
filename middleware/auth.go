package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/NikSchaefer/go-fiber/handlers"
	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
)

// func Authenticated(c *fiber.Ctx) error {
// 	type AuthRequest struct {
// 		Sessionid guuid.UUID `json:"sessionid"`
// 	}
// 	json := new(AuthRequest)
// 	if err := c.BodyParser(json); err != nil {
// 		return c.Status(422).JSON(err)
// 	}
// 	user, status := handlers.GetUser(json.Sessionid)
// 	if status != 0 {
// 		return c.SendStatus(status)
// 	}
// 	c.Locals("user", user)
// 	return c.Next()
// }

func Authenticated(c *fiber.Ctx) error {
	type AuthRequest struct {
		Sessionid string `json:"sessionid"`
	}
	data := new(AuthRequest)
	body := c.Body()
	if err := json.Unmarshal(body, &data); err != nil {
		var msg string
		switch t := err.(type) {
		case *json.SyntaxError:
			jsn := string(body[0:t.Offset])
			jsn += "<--(Invalid Character)"
			msg = fmt.Sprintf("Invalid character at offset %v\n %s", t.Offset, jsn)
		case *json.UnmarshalTypeError:
			jsn := string(body[0:t.Offset])
			jsn += "<--(Invalid Type)"
			msg = fmt.Sprintf("Invalid value at offset %v\n %s", t.Offset, jsn)
		default:
			msg = err.Error()
		}
		return c.Status(200).JSON(fiber.Map{
			"err":  err.Error(),
			"msg":  msg,
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
