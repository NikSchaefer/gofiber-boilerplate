package notifications

import (
	"errors"

	"github.com/NikSchaefer/go-fiber/config"
	"github.com/NikSchaefer/go-fiber/pkg/notifications/templates"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type smsSender struct {
	client          *twilio.RestClient
	fromPhoneNumber string
}

func NewSMSSender() SMSSender {
	return &smsSender{
		client:          twilio.NewRestClient(),
		fromPhoneNumber: config.TWILIO_PHONE_NUMBER,
	}
}

func (s *smsSender) SendSMS(phoneNumber string, data templates.SMSTemplateData) error {
	if phoneNumber == "" {
		return errors.New("phone number is required")
	}

	params := &api.CreateMessageParams{
		To:   &phoneNumber,
		Body: &data.Message,
		From: &s.fromPhoneNumber,
	}

	_, err := s.client.Api.CreateMessage(params)

	return err
}
