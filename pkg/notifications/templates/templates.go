package templates

import "github.com/matcornic/hermes"

type EmailTemplateData struct {
	Subject string
	Body    string
	Intros  []string
	Actions []hermes.Action
	Outros  []string
	Name    string
}

type SMSTemplateData struct {
	Message string
}

// Template defines content for each supported channel
type Template struct {
	ID    string
	Email func(data TemplateData) (EmailTemplateData, error)
	SMS   func(data TemplateData) (SMSTemplateData, error)
}

// TemplateData is the interface that all template data must implement
type TemplateData interface {
	Validate() error
}

var Templates = map[string]Template{
	"otp":                          OTPTemplate,
	"reset_password":               ResetPasswordTemplate,
}
