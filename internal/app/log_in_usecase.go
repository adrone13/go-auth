package app

import (
	"auth/internal/app/sessions"
	"auth/internal/app/tokens"
	"auth/internal/config"
	"auth/internal/crypto"
	"context"
	"fmt"
	"log"
)

/*
	TODO:
	- replace custom JWT library with third party to allow adding custom claims
	- write session id (or refresh token id) to refresh token JWT to enable multiple sessions
	for user (for example from different devices)
	- implement session creation in log_in_usecase
	- implement refresh token storage in sessions
	- implement JWT refreshing use case
	- implement logout and JWT invalidation (through cache)
	- implement HTTP wrappers to cut boilerplate code
*/

type Credentials struct {
	Email    string
	Password string
}

type LogInUseCase struct {
	UserRepository     UserRepository
	SessionsRepository SessionsRepository
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
		log.Printf("Failed to insert Session. Error: %s\n", err)

		return nil, err
	}

	access := tokens.CreateAccessToken(user)
	refresh := tokens.CreateRefreshToken(user, session)

	session.AddRefreshToken(refresh)
	err = u.SessionsRepository.Update(ctx, session)
	if err != nil {
		log.Printf("Failed to update Session. Error: %s\n", err)

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
