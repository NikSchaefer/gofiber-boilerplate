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
		return c.Status(422).JSON(err)
	}
	user, status := handlers.GetUser(json.Sessionid)
	if status != 0 {
		return c.SendStatus(status)
	}
	c.Locals("user", user)
	return c.Next()
}

/*
func Authenticated(c *fiber.Ctx) error {
	type AuthRequest struct {
		Sessionid guuid.UUID `json:"sessionid"`
	}
	js := new(AuthRequest)
	b := c.Body()
	if err := json.Unmarshal(c.Body(), &js); err != nil {
		var msg string
		switch t := err.(type) {
		case *json.SyntaxError:
			jsn := string(b[0:t.Offset])
			jsn += "<--(Invalid Character)"
			msg = fmt.Sprintf("Invalid character at offset %v\n %s", t.Offset, jsn)
		case *json.UnmarshalTypeError:
			jsn := string(b[0:t.Offset])
			jsn += "<--(Invalid Type)"
			msg = fmt.Sprintf("Invalid value at offset %v\n %s", t.Offset, jsn)
		default:
			msg = err.Error()
		}
		return c.Status(200).JSON(fiber.Map{
			"err":  err.Error(),
			"data": js,
			"msg":  msg,
		})

	}

*/
