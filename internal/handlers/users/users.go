package users_handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/NikSchaefer/go-fiber/ent"
	"github.com/NikSchaefer/go-fiber/internal/database"
	"github.com/NikSchaefer/go-fiber/pkg/validator"
)

type UpdateUserRequest struct {
	PhoneNumber *string `json:"phoneNumber" validate:"omitempty,e164"`
}

func UpdateUser(c *fiber.Ctx) error {
	data := new(UpdateUserRequest)
	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format")
	}

	err := validator.Validate(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	u := c.Locals("user").(*ent.User)

	update := u.Update()
	if data.PhoneNumber == nil {
		update.ClearPhoneNumber()
	} else {
		if *data.PhoneNumber == u.PhoneNumber {
			return c.JSON(u)
		}

		err := validator.ValidatePhoneUniqueness(c.Context(), database.DB, *data.PhoneNumber)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		update.SetPhoneNumber(*data.PhoneNumber)
	}

	_, err = update.Save(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	return c.JSON(u)
}

type UpdateProfileRequest struct {
	Name     string     `json:"name" validate:"required,min=2,max=100"`
	Birthday *time.Time `json:"birthday,omitempty"`
}

func UpdateProfile(c *fiber.Ctx) error {
	data := new(UpdateProfileRequest)
	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format")
	}

	if err := validator.Validate(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	u := c.Locals("user").(*ent.User)
	ctx := c.Context()

	pro, err := u.QueryProfile().Only(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	_, err = pro.Update().
		SetName(data.Name).
		SetNillableBirthday(data.Birthday).
		Save(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	return c.JSON(pro)
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DB
	u := c.Locals("user").(*ent.User)

	err := db.User.DeleteOneID(u.ID).Exec(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func GetCurrentUserInfo(c *fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)
	return c.JSON(user)
}

func GetUserProfile(c *fiber.Ctx) error {
	u := c.Locals("user").(*ent.User)
	pro, err := u.QueryProfile().Only(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}
	return c.JSON(pro)
}
