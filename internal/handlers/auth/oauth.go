package auth_handlers

import (
	"encoding/json"
	"time"

	"github.com/NikSchaefer/go-fiber/config"
	"github.com/NikSchaefer/go-fiber/ent/schema"
	"github.com/NikSchaefer/go-fiber/internal/database"
	"github.com/NikSchaefer/go-fiber/internal/services"
	"github.com/NikSchaefer/go-fiber/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

// GetGoogleAuthRedirect initiates the OAuth flow by redirecting the user to Google's consent page
func GetGoogleAuthRedirect(c *fiber.Ctx) error {
	// Get the OAuth configuration for Google
	conf := config.ConfigGoogle()

	// Generate a random state string to prevent CSRF attacks
	state := oauth2.GenerateVerifier()

	// Generate the authorization URL with the state parameter
	// AccessTypeOffline allows getting a refresh token
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Store the state in a cookie for verification when Google calls back
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state_google",
		Value:    state,
		Expires:  time.Now().Add(1 * time.Hour),
		HTTPOnly: true, // Prevents JavaScript access to the cookie
	})

	// Send user to Google's consent page
	return c.JSON(url)
}

type GoogleAuthCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

// GetGoogleAuthCallback handles the response from Google after user consents
func GetGoogleAuthCallback(c *fiber.Ctx) error {
	data := new(GoogleAuthCallbackRequest)
	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	err := validator.Validate(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Get the OAuth configuration
	conf := config.ConfigGoogle()

	// Exchange the authorization code for an access token
	token, err := conf.Exchange(c.Context(), data.Code)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Get the state from the cookie
	savedState := c.Cookies("oauth_state_google")
	if savedState == "" {
		return fiber.NewError(fiber.StatusBadRequest, "state not found")
	}
	if savedState != data.State {
		return fiber.NewError(fiber.StatusBadRequest, "state does not match")
	}

	// Use the access token to fetch the user's information from Google
	resp, err := conf.Client(c.Context(), token).Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	// Define a struct to parse the user information from Google
	var userInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	// Parse the JSON response from Google into the userInfo struct
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user info")
	}

	// Check if the email is verified
	if !userInfo.VerifiedEmail {
		return fiber.NewError(fiber.StatusBadRequest, "Google email not verified")
	}

	u, err := services.HandleOauthLogin(services.HandleOauthLoginStruct{
		Email:      userInfo.Email,
		Name:       userInfo.Name,
		Type:       "google",
		ProviderID: userInfo.ID,
		AvatarURL:  &userInfo.Picture,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	db := database.DB
	// create session
	s, err := db.Session.Create().
		SetUser(u).
		Save(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    s.ID.String(),
		Expires:  schema.GetTokenExpiration(),
		HTTPOnly: true,
	})

	return c.JSON(u)
}

type GoogleNativeAuthCallbackRequest struct {
	AccessToken string `json:"accessToken" validate:"required"`
}

func GetGoogleNativeAuthCallback(c *fiber.Ctx) error {
	data := new(GoogleNativeAuthCallbackRequest)
	if err := c.BodyParser(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	err := validator.Validate(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	conf := config.ConfigGoogle()

	// Create a token struct with the provided access token
	token := &oauth2.Token{
		AccessToken: data.AccessToken,
	}

	// Use the token to fetch the user's information from Google
	resp, err := conf.Client(c.Context(), token).Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	// Define a struct to parse the user information from Google
	var userInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	// Parse the JSON response from Google into the userInfo struct
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse user info")
	}

	if userInfo.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Google email not returned: "+userInfo.Email)
	}

	// // Check if the email is verified
	// if !userInfo.VerifiedEmail {
	// 	return fiber.NewError(fiber.StatusBadRequest, "Google email not verified")
	// }

	u, err := services.HandleOauthLogin(services.HandleOauthLoginStruct{
		Email:      userInfo.Email,
		Name:       userInfo.Name,
		Type:       "google",
		ProviderID: userInfo.ID,
		AvatarURL:  &userInfo.Picture,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	db := database.DB
	// create session
	s, err := db.Session.Create().
		SetUser(u).
		Save(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    s.ID.String(),
		Expires:  schema.GetTokenExpiration(),
		HTTPOnly: true,
	})

	return c.JSON(u)
}
