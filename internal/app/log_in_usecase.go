package app

import (
	"auth/internal/app/common/repositories"
	"auth/internal/app/sessions"
	"auth/internal/app/tokens"
	"auth/internal/config"
	"auth/internal/crypto"
	"auth/internal/logger"
	"context"
	"fmt"
)

/*
	TODO:
	- implement logout and JWT invalidation
*/

type Credentials struct {
	Email    string
	Password string
}

type LogInUseCase struct {
	UserRepository     repositories.UserRepository
	SessionsRepository repositories.SessionsRepository
}

func (u *LogInUseCase) Execute(ctx context.Context, cred Credentials) (*Auth, error) {
	user, err := u.UserRepository.FindByEmail(ctx, cred.Email)
	if err != nil {
		return nil, err
	}

	if crypto.ComparePasswordHash(cred.Password, user.Password) {
		return nil, &InvalidPasswordError{}
	}

	accessTtl := config.Values.AccessTokenTtl

	session := sessions.NewSession(user.Id)
	err = u.SessionsRepository.Insert(ctx, session)
	if err != nil {
		logger.Error.Printf("Failed to insert Session. Error: %s\n", err)

		return nil, err
	}

	access := tokens.CreateAccessToken(user)
	refresh := tokens.CreateRefreshToken(user, session)

	session.AddRefreshToken(refresh)
	err = u.SessionsRepository.Update(ctx, session)
	if err != nil {
		logger.Error.Printf("Failed to update Session. Error: %s\n", err)

		return nil, err
	}

	auth := &Auth{
		AccessToken:  fmt.Sprintf("Bearer %s", access),
		TokenType:    "bearer",
		ExpiresIn:    accessTtl,
		RefreshToken: refresh,
	}

	return auth, nil
}
