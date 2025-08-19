package auth_handlers

import (
	"github.com/NikSchaefer/go-fiber/ent"
	"github.com/NikSchaefer/go-fiber/ent/account"
	"github.com/NikSchaefer/go-fiber/ent/schema"
	"github.com/NikSchaefer/go-fiber/ent/user"
	"github.com/NikSchaefer/go-fiber/internal/database"
	"github.com/NikSchaefer/go-fiber/pkg/analytics"
	"github.com/NikSchaefer/go-fiber/pkg/utils"
	"github.com/NikSchaefer/go-fiber/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func LoginWithPassword(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email" validate:"omitempty,email"`
		Password string `json:"password" validate:"required,min=8"`
		Phone    string `json:"phone" validate:"omitempty,e164"`
	}
	data := new(LoginRequest)
	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON sent")
	}

	err := validator.Validate(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if data.Email == "" && data.Phone == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email or Phone is required")
	}

	db := database.DB
	u, err := db.User.Query().Where(
		user.Or(
			user.EmailEQ(data.Email),
			user.PhoneNumberEQ(data.Phone),
		),
	).
		WithAccounts(func(q *ent.AccountQuery) {
			q.Where(account.TypeEQ(account.TypePassword))
		}).
		Only(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// Check if user has a password account
	if len(u.Edges.Accounts) == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "No password found for this user")
	}

	if !utils.ComparePasswords(u.Edges.Accounts[0].PasswordHash, []byte(data.Password)) {
		return fiber.NewError(fiber.StatusUnauthorized, "Password is incorrect")
	}

	// create session
	s, err := db.Session.Create().
		SetUser(u).
		Save(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    s.ID.String(),
		Expires:  schema.GetTokenExpiration(),
		HTTPOnly: true,
	})

	analytics.TrackEventWithUser("password_login", map[string]interface{}{
		"user": u.ID,
	}, u)

	return c.JSON(fiber.Map{
		"session": s.ID,
	})
}

func Logout(c *fiber.Ctx) error {
	sessionID := c.Cookies("session")
	if sessionID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "No session provided")
	}

	id, err := uuid.Parse(sessionID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid session format")
	}
	db := database.DB

	err = db.Session.DeleteOneID(id).Exec(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
