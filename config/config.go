package config

import "os"

const (
	TWILIO_PHONE_NUMBER = "+00000000"
	APP_DOMAIN = "example.com"
)

func GetIsProduction() bool {
	return os.Getenv("STAGE") == "prod"
}

func GetPosthogKey() string {
	return os.Getenv("POSTHOG_KEY")
}

func GetResendKey() string {
	return os.Getenv("RESEND_KEY")
}

func GetURL() string {
	return APP_DOMAIN
}
