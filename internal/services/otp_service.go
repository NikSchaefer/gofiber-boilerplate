package services

import (
	"context"

	"github.com/NikSchaefer/go-fiber/ent"
	"github.com/NikSchaefer/go-fiber/internal/database"
)

func GenerateOTP(user *ent.User) (*ent.OTP, error) {
	db := database.DB
	ctx := context.Background()

	return db.OTP.Create().
		SetUser(user).
		Save(ctx)
}
