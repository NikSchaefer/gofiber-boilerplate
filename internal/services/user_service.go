package services

import (
	"context"
	"strings"

	"github.com/NikSchaefer/go-fiber/ent"
	"github.com/NikSchaefer/go-fiber/internal/database"
	"github.com/NikSchaefer/go-fiber/pkg/analytics"
	"github.com/NikSchaefer/go-fiber/pkg/utils"
	"github.com/NikSchaefer/go-fiber/pkg/validator"
)

type CreateUserStruct struct {
	Name                   string     `validate:"required,min=2,max=100"`
	Email                  string     `validate:"required,email"`
	Password               *string    `validate:"omitempty,min=8"`
	Phone                  *string    `validate:"omitempty,e164"`
}

func CreateUserWithTx(ctx context.Context, tx *ent.Tx, data CreateUserStruct) (*ent.User, error) {
	err := validator.Validate(data)
	if err != nil {
		return nil, err
	}

	var pw []byte
	if data.Password != nil {
		pw, err = utils.SaltAndVerifyPassword(*data.Password)
		if err != nil {
			return nil, err
		}
	}

	// Check if email exists
	err = validator.ValidateEmailUniqueness(ctx, tx.Client(), data.Email)
	if err != nil {
		return nil, err
	}

	// Check if phone exists (if provided)
	if data.Phone != nil {
		err = validator.ValidatePhoneUniqueness(ctx, tx.Client(), *data.Phone)
		if err != nil {
			return nil, err
		}
	}

	userEntity, err := tx.User.Create().
		SetEmail(strings.ToLower(data.Email)).
		SetNillablePhoneNumber(data.Phone).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	_, err = tx.Profile.Create().
		SetUser(userEntity).
		SetName(data.Name).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if data.Password != nil {
		_, err = tx.Account.Create().
			SetUser(userEntity).
			SetType("password").
			SetPasswordHash(pw).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	return userEntity, nil
}

func CreateUser(data CreateUserStruct) (*ent.User, error) {
	db := database.DB
	ctx := context.Background()

	// Create user and related entities in a single transaction
	tx, err := db.Tx(ctx)
	if err != nil {
		return nil, err
	}

	userEntity, err := CreateUserWithTx(ctx, tx, data)
	if err != nil {
		return nil, utils.RollbackTx(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	err = analytics.IdentifyUser(userEntity.ID.String(), map[string]interface{}{
		"name": data.Name,
	})
	if err != nil {
		return nil, err
	}

	return userEntity, nil
}

