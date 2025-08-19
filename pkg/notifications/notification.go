package notifications

import (
	"fmt"

	"github.com/NikSchaefer/go-fiber/pkg/notifications/templates"
)

// ChannelType represents the type of notification channel
type ChannelType string

const (
	EmailChannel ChannelType = "email"
	SMSChannel   ChannelType = "sms"
)

// NotificationRequest represents a complete notification request
// This is the main request that will be used to send notifications
// Only the fields that are 'to' is filled will be sent

type NotificationRequest struct {
	TemplateID   string
	Data         templates.TemplateData
	EmailAddress *string
	PhoneNumber  *string
}

// NotificationService handles sending notifications across different channels
type NotificationService struct {
	templates   map[string]templates.Template
	emailSender EmailSender
	smsSender   SMSSender
}

// Sender interface becomes channel-specific
type EmailSender interface {
	SendEmail(string, templates.EmailTemplateData) error
}

type SMSSender interface {
	SendSMS(string, templates.SMSTemplateData) error
}

// Service is the global notification service instance
var service *NotificationService

// InitService initializes the global notification service
func InitService() {
	service = &NotificationService{
		templates:   templates.Templates,
		emailSender: NewEmailSender(),
	}
}

// Send sends a notification to the specified channels based on provided contact information
func Send(req NotificationRequest) error {
	template, ok := service.templates[req.TemplateID]
	if !ok {
		return fmt.Errorf("template not found: %s", req.TemplateID)
	}

	// Add nil check for Data
	if req.Data == nil {
		return fmt.Errorf("template data cannot be nil")
	}

	// Validate the data
	if err := req.Data.Validate(); err != nil {
		return fmt.Errorf("invalid template data: %w", err)
	}

	if req.EmailAddress == nil && req.PhoneNumber == nil {
		return fmt.Errorf("no contact information provided")
	}

	// Send email if email address is provided
	if req.EmailAddress != nil {
		emailData, err := template.Email(req.Data)
		if err != nil {
			return fmt.Errorf("failed to generate email data: %w", err)
		}

		err = service.emailSender.SendEmail(*req.EmailAddress, emailData)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}

	// Send SMS if phone number is provided
	if req.PhoneNumber != nil && *req.PhoneNumber != "" {
		smsData, err := template.SMS(req.Data)
		if err != nil {
			return fmt.Errorf("failed to generate SMS data: %w", err)
		}

		err = service.smsSender.SendSMS(*req.PhoneNumber, smsData)
		if err != nil {
			return fmt.Errorf("failed to send SMS: %w", err)
		}
	}

	return nil
}
