package templates

import (
	"fmt"

	"github.com/matcornic/hermes"
)

// OTPTemplateData is a specific template data type for OTP notifications
type OTPTemplateData struct {
	OTP  string
	Name string
}

// Validate implements TemplateData interface for OTPTemplateData
func (d *OTPTemplateData) Validate() error {
	if d.OTP == "" {
		return fmt.Errorf("OTP cannot be empty")
	}
	return nil
}

var OTPTemplate = Template{
	ID: "otp",
	Email: func(data TemplateData) (EmailTemplateData, error) {
		otpData, ok := data.(*OTPTemplateData)
		if !ok {
			return EmailTemplateData{}, fmt.Errorf("invalid template data type")
		}

		return EmailTemplateData{
			Name:    otpData.Name,
			Subject: "Your verification code",
			Actions: []hermes.Action{
				{
					Instructions: "Please use the following code to verify your account:",
					InviteCode:   otpData.OTP,
				},
			},
			Outros: []string{
				"If you did not request this code, please ignore this email.",
			},
		}, nil
	},
	SMS: func(data TemplateData) (SMSTemplateData, error) {
		otpData, ok := data.(*OTPTemplateData)
		if !ok {
			return SMSTemplateData{}, fmt.Errorf("invalid template data type")
		}

		return SMSTemplateData{
			Message: fmt.Sprintf("Your verification code is: %s. Don't share this code with anyone.", otpData.OTP),
		}, nil
	},
}
