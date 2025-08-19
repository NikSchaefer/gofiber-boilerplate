package router

import (
	auth_handlers "github.com/NikSchaefer/go-fiber/internal/handlers/auth"
	user_handlers "github.com/NikSchaefer/go-fiber/internal/handlers/users"
	"github.com/NikSchaefer/go-fiber/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	router.Use(middleware.Security)

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello, World!")
	})

	auth := router.Group("/auth")
	{
		// Login related
		auth.Post("/login/password", auth_handlers.LoginWithPassword)
		auth.Post("/login/otp/request", auth_handlers.RequestLoginWithOTP)
		auth.Post("/login/otp/verify", auth_handlers.VerifyLoginWithOTP)
		auth.Delete("/logout", middleware.Authenticated, auth_handlers.Logout)

		// OAuth routes
		auth.Post("/oauth/google", auth_handlers.GetGoogleAuthRedirect)
		auth.Post("/oauth/google/callback", auth_handlers.GetGoogleAuthCallback)
		auth.Post("/oauth/google/native/callback", auth_handlers.GetGoogleNativeAuthCallback)

		// Registration
		auth.Post("/signup", auth_handlers.SignUp)

		// Password management
		auth.Post("/password/change", middleware.Authenticated, auth_handlers.ChangePassword)
		auth.Post("/password/reset/request", auth_handlers.ResetPassword)
		auth.Post("/password/reset/verify", auth_handlers.VerifyResetPassword)
	}

	// User routes
	user := router.Group("/users", middleware.Authenticated)
	{
		user.Get("/me", user_handlers.GetCurrentUserInfo)
		user.Get("/profile", user_handlers.GetUserProfile)
		user.Patch("/", user_handlers.UpdateUser)
		user.Patch("/profile", user_handlers.UpdateProfile)
		user.Delete("/", user_handlers.DeleteUser)
	}
}
