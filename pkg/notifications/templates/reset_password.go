package templates

import (
	"fmt"

	"github.com/matcornic/hermes"
	"github.com/NikSchaefer/go-fiber/config"
)

// ResetPasswordTemplateData is a specific template data type for password reset notifications
type ResetPasswordTemplateData struct {
	ResetCode string
	Email     string
	Name      string
}

// Validate implements TemplateData interface for ResetPasswordTemplateData
func (d *ResetPasswordTemplateData) Validate() error {
	if d.ResetCode == "" {
		return fmt.Errorf("reset code cannot be empty")
	}
	return nil
}

var ResetPasswordTemplate = Template{
	ID: "reset_password",
	Email: func(data TemplateData) (EmailTemplateData, error) {
		resetData, ok := data.(*ResetPasswordTemplateData)
		if !ok {
			return EmailTemplateData{}, fmt.Errorf("invalid template data type")
		}

		return EmailTemplateData{
			Subject: "Reset Your Password",
			Name:    resetData.Name,
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password:",
					Button: hermes.Button{
						Text:      "Reset Password",
						Link:      fmt.Sprintf("%s/reset-password/verify?code=%s&email=%s", config.GetURL(), resetData.ResetCode, resetData.Email),
					},
				},
			},
			Outros: []string{
				"If you did not request a password reset, please ignore this email.",
			},
		}, nil
	},
	SMS: func(data TemplateData) (SMSTemplateData, error) {
		return SMSTemplateData{}, fmt.Errorf("sms notifications not supported for password reset")
	},
}
