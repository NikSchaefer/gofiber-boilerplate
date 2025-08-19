package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ConfigGoogle to set config of oauth
func ConfigGoogle() *oauth2.Config {
	clientID := GetGoogleClientID()
	clientSecret := GetGoogleClientSecret()
	redirectURL := GetGoogleRedirectURL()

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		panic("GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, or GOOGLE_REDIRECT_URL is not set")
	}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}
