package notifications

import (
	"strings"

	"github.com/NikSchaefer/go-fiber/config"
	"github.com/NikSchaefer/go-fiber/pkg/notifications/templates"
	"github.com/matcornic/hermes"
	resend "github.com/resend/resend-go/v2"
)

const (
	sender          = "info@example.com"
	displayName     = "YourAppName"
	appDownloadLink = "https://your-app-domain.com"
	signature       = "Best regards"
	logo            = "https://your-app-domain.com/logo.svg"
)

type emailSender struct {
	hermesConfig hermes.Hermes
	resendClient *resend.Client
}

func NewEmailSender() *emailSender {
	return &emailSender{
		hermesConfig: hermes.Hermes{
			Product: hermes.Product{
				Name:      displayName,
				Link:      appDownloadLink,
				Logo:      logo,
				Copyright: "Copyright © 2025 YourAppName LLC. All rights reserved.",
			},
		},
		resendClient: resend.NewClient(config.GetResendKey()),
	}
}

func (s *emailSender) SendEmail(emailAddress string, data templates.EmailTemplateData) error {
	email := hermes.Email{
		Body: hermes.Body{
			Name:      strings.Split(data.Name, " ")[0],
			Intros:    append(data.Intros, data.Body),
			Actions:   data.Actions,
			Outros:    data.Outros,
			Signature: signature,
		},
	}

	html, err := s.hermesConfig.GenerateHTML(email)
	if err != nil {
		return err
	}

	text, err := s.hermesConfig.GeneratePlainText(email)
	if err != nil {
		return err
	}

	resendRequest := &resend.SendEmailRequest{
		From:    sender,
		To:      []string{emailAddress},
		Subject: data.Subject,
		Html:    html,
		Text:    text,
	}

	_, err = s.resendClient.Emails.Send(resendRequest)

	if err != nil {
		return err
	}

	return nil
}
