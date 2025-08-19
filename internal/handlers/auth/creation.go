package auth_handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/NikSchaefer/go-fiber/internal/services"
	"github.com/NikSchaefer/go-fiber/pkg/notifications"
	"github.com/NikSchaefer/go-fiber/pkg/notifications/templates"
	"github.com/NikSchaefer/go-fiber/pkg/validator"
)

func SignUp(c *fiber.Ctx) error {
	type SignUpRequest struct {
		Name     string `json:"name" validate:"required,min=2,max=100"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	data := new(SignUpRequest)
	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON Sent")
	}

	err := validator.Validate(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	u, err := services.CreateUser(services.CreateUserStruct{
		Name:     data.Name,
		Email:    data.Email,
		Password: &data.Password,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, ""+err.Error())
	}

	otp, err := services.GenerateOTP(u)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	err = notifications.Send(notifications.NotificationRequest{
		TemplateID: "otp",
		Data: &templates.OTPTemplateData{
			OTP:  otp.Code,
			Name: data.Name,
		},
		EmailAddress: &u.Email,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	return c.JSON(u)
}

