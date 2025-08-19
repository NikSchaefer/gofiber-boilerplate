package config

import "os"

// Database Configuration
func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

// Server Configuration
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8000"
	}
	return port
}

func GetIsProduction() bool {
	return os.Getenv("STAGE") == "prod"
}

func GetAllowedOrigins() string {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		return "http://localhost:3000,http://localhost:3001"
	}
	return origins
}

// External Services
func GetPosthogKey() string {
	return os.Getenv("POSTHOG_KEY")
}

func GetResendKey() string {
	return os.Getenv("RESEND_KEY")
}

// Twilio Configuration
func GetTwilioAccountSID() string {
	return os.Getenv("TWILIO_ACCOUNT_SID")
}

func GetTwilioAuthToken() string {
	return os.Getenv("TWILIO_AUTH_TOKEN")
}

func GetTwilioPhoneNumber() string {
	phoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")
	if phoneNumber == "" {
		return "+00000000"
	}
	return phoneNumber
}

// OAuth Configuration
func GetGoogleClientID() string {
	return os.Getenv("GOOGLE_CLIENT_ID")
}

func GetGoogleClientSecret() string {
	return os.Getenv("GOOGLE_CLIENT_SECRET")
}

func GetGoogleRedirectURL() string {
	return os.Getenv("GOOGLE_REDIRECT_URL")
}

// Application Configuration
func GetAppDomain() string {
	domain := os.Getenv("APP_DOMAIN")
	if domain == "" {
		return "localhost:8000"
	}
	return domain
}

func GetURL() string {
	return GetAppDomain()
}
