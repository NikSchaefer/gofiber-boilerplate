package services

import (
	"context"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"github.com/NikSchaefer/go-fiber/ent"
	"github.com/NikSchaefer/go-fiber/ent/session"
	"github.com/NikSchaefer/go-fiber/internal/database"
)

// ValidateSession checks if a session is valid and returns the associated user
func ValidateSession(ctx context.Context, sessionID string) (*ent.User, error) {
	if sessionID == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "No session provided")
	}

	id, err := guuid.Parse(sessionID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid session format")
	}
	db := database.DB

	session, err := db.Session.Query().
		Where(session.ID(id)).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid session id")
	}

	user, err := session.QueryUser().
		Only(ctx)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid session user")
	}

	return user, nil
}
