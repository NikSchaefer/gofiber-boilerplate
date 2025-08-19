package services

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/NikSchaefer/go-fiber/ent"
	"github.com/NikSchaefer/go-fiber/ent/account"
	"github.com/NikSchaefer/go-fiber/ent/user"
	"github.com/NikSchaefer/go-fiber/internal/database"
	"github.com/NikSchaefer/go-fiber/pkg/analytics"
	"github.com/NikSchaefer/go-fiber/pkg/utils"
	"github.com/NikSchaefer/go-fiber/pkg/validator"
)

type CreateOAuthUserStruct struct {
	Email      string    `validate:"required,email"`
	Name       string    `validate:"required,min=2,max=100"`
	Type       string    `validate:"required,oneof=google apple"`
	ProviderID string    `validate:"required"`
	AvatarURL  *string   `validate:"omitempty,url"`
}

func CreateOAuthUser(data CreateOAuthUserStruct) (*ent.User, error) {
	err := validator.Validate(data)
	if err != nil {
		return nil, err
	}

	db := database.DB
	ctx := context.Background()

	err = validator.ValidateEmailUniqueness(ctx, db, data.Email)
	if err != nil {
		return nil, err
	}

	tx, err := db.Tx(ctx)
	if err != nil {
		return nil, err
	}

	userEntity, err := tx.User.Create().
		SetEmail(strings.ToLower(data.Email)).
		SetEmailVerified(true).
		Save(ctx)
	if err != nil {
		return nil, utils.RollbackTx(tx, err)
	}

	_, err = tx.Profile.Create().
		SetUser(userEntity).
		SetName(data.Name).
		Save(ctx)
	if err != nil {
		return nil, utils.RollbackTx(tx, err)
	}

	_, err = tx.Account.Create().
		SetUser(userEntity).
		SetType(account.Type(data.Type)).
		SetProviderID(data.ProviderID).
		Save(ctx)
	if err != nil {
		return nil, utils.RollbackTx(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return userEntity, nil
}

type AddOauthAccountToUserStruct struct {
	UserID     uuid.UUID `validate:"required"`
	Type       string    `validate:"required,oneof=google apple"`
	ProviderID string    `validate:"required"`
	AvatarURL  *string   `validate:"omitempty,url"`
}

func AddOauthAccountToUser(data AddOauthAccountToUserStruct) (*ent.User, error) {
	err := validator.Validate(data)
	if err != nil {
		return nil, err
	}

	db := database.DB
	ctx := context.Background()

	user, err := db.User.Query().Where(user.ID(data.UserID)).WithProfile().Only(ctx)
	if err != nil {
		return nil, err
	}

	_, err = db.Account.Create().
		SetUser(user).
		SetType(account.Type(data.Type)).
		SetProviderID(data.ProviderID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	analytics.TrackEventWithUser("oauth_account_added", map[string]interface{}{
		"type": data.Type,
	}, user)

	return user, nil
}

type HandleOauthLoginStruct struct {
	Email      string  `validate:"omitempty,email"`
	Name       string  `validate:"omitempty,min=2,max=100"`
	Type       string  `validate:"required,oneof=google apple"`
	ProviderID string  `validate:"required"`
	AvatarURL  *string `validate:"omitempty,url"`
}

func HandleOauthLogin(data HandleOauthLoginStruct) (*ent.User, error) {
	err := validator.Validate(data)
	if err != nil {
		return nil, err
	}

	db := database.DB
	ctx := context.Background()

	// First try to find by provider ID
	providerAccount, err := db.Account.Query().
		Where(account.And(
			account.TypeEQ(account.Type(data.Type)),
			account.ProviderIDEQ(data.ProviderID),
		)).
		WithUser().
		Only(ctx)
	if err == nil {
		// Found existing OAuth account
		return providerAccount.QueryUser().Only(ctx)
	} else if !ent.IsNotFound(err) {
		return nil, err
	}

	// Then try to find by email
	u, err := db.User.Query().Where(user.Email(data.Email)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// No existing account found, create new one
			return CreateOAuthUser(CreateOAuthUserStruct{
				Email:      data.Email,
				Name:       data.Name,
				Type:       data.Type,
				ProviderID: data.ProviderID,
				AvatarURL:  data.AvatarURL,
			})
		}
		return nil, err
	}

	// Found existing user by email - link the accounts
	return AddOauthAccountToUser(AddOauthAccountToUserStruct{
		UserID:     u.ID,
		Type:       data.Type,
		ProviderID: data.ProviderID,
		AvatarURL:  data.AvatarURL,
	})
}
