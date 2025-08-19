package auth_handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/NikSchaefer/go-fiber/ent"
	"github.com/NikSchaefer/go-fiber/ent/account"
	"github.com/NikSchaefer/go-fiber/ent/otp"
	"github.com/NikSchaefer/go-fiber/ent/user"
	"github.com/NikSchaefer/go-fiber/internal/database"
	"github.com/NikSchaefer/go-fiber/pkg/notifications"
	"github.com/NikSchaefer/go-fiber/pkg/notifications/templates"
	"github.com/NikSchaefer/go-fiber/pkg/utils"
	"github.com/NikSchaefer/go-fiber/pkg/validator"
)

func ChangePassword(c *fiber.Ctx) error {
	type ChangePasswordRequest struct {
		Password    string `json:"password" validate:"required,min=8"`
		NewPassword string `json:"newPassword" validate:"required,min=8"`
	}

	data := new(ChangePasswordRequest)
	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON sent")
	}

	err := validator.Validate(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	db := database.DB
	u := c.Locals("user").(*ent.User)

	acc, err := db.Account.Query().
		Where(account.HasUserWith(user.IDEQ(u.ID))).
		Where(account.TypeEQ(account.TypePassword)).
		Only(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	if !utils.ComparePasswords(acc.PasswordHash, []byte(data.Password)) {
		return fiber.NewError(fiber.StatusUnauthorized, "Current password is incorrect")
	}

	newHash, err := utils.HashAndSalt([]byte(data.NewPassword))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	_, err = acc.Update().
		SetPasswordHash(newHash).
		Save(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}

func ResetPassword(c *fiber.Ctx) error {
	type ResetPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}
	data := new(ResetPasswordRequest)
	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format")
	}

	err := validator.Validate(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	db := database.DB
	u, err := db.User.Query().
		Where(user.EmailEQ(data.Email)).
		WithProfile().
		Only(c.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	o, err := db.OTP.Create().
		SetType(otp.TypePasswordReset).
		SetUser(u).
		Save(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	err = notifications.Send(notifications.NotificationRequest{
		TemplateID: "reset_password",
		Data: &templates.ResetPasswordTemplateData{
			ResetCode: o.Code,
			Name:      u.Edges.Profile.Name,
			Email:     u.Email,
		},
		EmailAddress: &u.Email,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to send reset email: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Reset email sent successfully",
	})
}

func VerifyResetPassword(c *fiber.Ctx) error {
	type VerifyResetPasswordRequest struct {
		Email       string `json:"email" validate:"required,email"`
		Code        string `json:"code" validate:"required"`
		NewPassword string `json:"newPassword" validate:"required,min=8"`
	}
	data := new(VerifyResetPasswordRequest)
	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format")
	}

	err := validator.Validate(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	db := database.DB

	u, err := db.User.Query().
		Where(user.EmailEQ(data.Email)).
		Only(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "User not found")
	}

	o, err := db.OTP.Query().
		Where(otp.And(
			otp.CodeEQ(data.Code),
			otp.TypeEQ(otp.TypePasswordReset),
			otp.ExpiresAtGTE(time.Now()),
			otp.HasUserWith(user.IDEQ(u.ID)),
		)).
		Only(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid code")
	}

	if o.Used {
		return fiber.NewError(fiber.StatusBadRequest, "Code already used")
	}

	acc, err := db.Account.Query().
		Where(account.HasUserWith(user.IDEQ(u.ID))).
		Where(account.TypeEQ(account.TypePassword)).
		Only(c.Context())
	if err != nil {
		if !ent.IsNotFound(err) {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
		}
		// Create new password account if it doesn't exist
		newHash, err := utils.HashAndSalt([]byte(data.NewPassword))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
		}
		_, err = db.Account.Create().
			SetType(account.TypePassword).
			SetUser(u).
			SetPasswordHash(newHash).
			Save(c.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
		}
	} else {
		// Update existing account
		newHash, err := utils.HashAndSalt([]byte(data.NewPassword))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
		}
		_, err = acc.Update().SetPasswordHash(newHash).Save(c.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
		}
	}

	_, err = o.Update().SetUsed(true).Save(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Reset code verified successfully",
	})
}
