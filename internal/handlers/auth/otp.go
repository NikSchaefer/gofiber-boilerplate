package auth_handlers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/NikSchaefer/go-fiber/ent/otp"
	"github.com/NikSchaefer/go-fiber/ent/predicate"
	"github.com/NikSchaefer/go-fiber/ent/schema"
	"github.com/NikSchaefer/go-fiber/ent/user"
	"github.com/NikSchaefer/go-fiber/internal/database"
	"github.com/NikSchaefer/go-fiber/internal/services"
	"github.com/NikSchaefer/go-fiber/pkg/analytics"
	"github.com/NikSchaefer/go-fiber/pkg/notifications"
	"github.com/NikSchaefer/go-fiber/pkg/notifications/templates"
	"github.com/NikSchaefer/go-fiber/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func RequestLoginWithOTP(c *fiber.Ctx) error {
	type RequestOTPRequest struct {
		Email string `json:"email" validate:"omitempty,email"`
		Phone string `json:"phone" validate:"omitempty,e164"`
	}

	data := new(RequestOTPRequest)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data sent")
	}

	err := validator.Validate(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if data.Email == "" && data.Phone == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email or Phone is required")
	}

	if data.Email != "" && data.Phone != "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email and Phone cannot be used together")
	}

	db := database.DB
	var conditions []predicate.User

	if data.Email != "" {
		data.Email = strings.ToLower(data.Email)
		conditions = append(conditions, user.EmailEQ(data.Email))
	}
	if data.Phone != "" {
		conditions = append(conditions, user.PhoneNumberEQ(data.Phone))
	}

	u, err := db.User.Query().
		Where(user.Or(conditions...)).
		WithProfile().
		Only(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	otp, err := services.GenerateOTP(u)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	var emailToSend *string
	var phoneToSend *string

	if data.Email != "" {
		emailToSend = &u.Email
	} else {
		phoneToSend = &u.PhoneNumber
	}

	err = notifications.Send(notifications.NotificationRequest{
		TemplateID: "otp",
		Data: &templates.OTPTemplateData{
			OTP:  otp.Code,
			Name: u.Edges.Profile.Name,
		},
		EmailAddress: emailToSend,
		PhoneNumber:  phoneToSend,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func VerifyLoginWithOTP(c *fiber.Ctx) error {
	type VerifyOTPRequest struct {
		Email string `json:"email" validate:"omitempty,email"`
		Phone string `json:"phone" validate:"omitempty,e164"`
		Code  string `json:"code" validate:"required"`
	}
	data := new(VerifyOTPRequest)
	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON Sent")
	}

	err := validator.Validate(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if data.Email == "" && data.Phone == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email or Phone is required")
	}

	if data.Email != "" && data.Phone != "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email and Phone cannot be used together")
	}

	db := database.DB

	var conditions []predicate.User

	if data.Email != "" {
		data.Email = strings.ToLower(data.Email)
		conditions = append(conditions, user.EmailEQ(data.Email))
	}
	if data.Phone != "" {
		conditions = append(conditions, user.PhoneNumberEQ(data.Phone))
	}

	u, err := db.User.Query().Where(user.Or(conditions...)).Only(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	o, err := db.OTP.Query().Where(
		otp.And(
			otp.CodeEQ(data.Code),
			otp.HasUserWith(user.IDEQ(u.ID)),
			otp.ExpiresAtGTE(time.Now()),
			otp.TypeEQ(otp.TypeLogin),
		),
	).Only(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Code not found")
	}

	if o.Used {
		return fiber.NewError(fiber.StatusUnauthorized, "Code already used")
	}

	update := u.Update()

	if data.Email != "" {
		update.SetEmailVerified(true)
	}

	if data.Phone != "" {
		update.SetPhoneNumberVerified(true)
	}

	_, err = update.Save(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	_, err = o.Update().SetUsed(true).Save(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
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

	analytics.TrackEventWithUser("otp_login", map[string]interface{}{
		"user": u.ID,
	}, u)

	return c.JSON(u)
}
