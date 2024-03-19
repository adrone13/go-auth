package app

import (
	"auth/internal/crypto"
	"context"
	"time"

	"auth/internal/config"
	"github.com/adrone13/gojwt"
)

/*
	TODO:
	- move JWT creation into separate library +
	- add DB implementation +
	- implement registration use case +
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
	UserRepository UserRepository
}

func (u *LogInUseCase) Execute(ctx context.Context, cred Credentials) (*Auth, error) {
	user, err := u.UserRepository.FindByEmail(ctx, cred.Email)
	if err != nil {
		return nil, err
	}

	if crypto.ComparePasswordHash(cred.Password, user.Password) {
		return nil, &InvalidPasswordError{}
	}

	secret := config.Values.JwtSecret
	accessTtl := config.Values.JwtTtl
	refreshTtl := config.Values.RefreshTokenAbsoluteTtl

	access := jwt.Sign(jwt.Claims{
		Issuer:     "auth",
		Expiration: time.Now().Add(time.Second * time.Duration(accessTtl)).Unix(),
		Audience:   "todo",
		Subject:    string(user.Id),
		Name:       user.FullName,
		Roles:      []string{"TODO"},
	}, secret)

	refresh := jwt.Sign(jwt.Claims{
		Subject:    string(user.Id),
		Expiration: time.Now().Add(time.Second * time.Duration(refreshTtl)).Unix(),
	}, secret)

	auth := &Auth{
		AccessToken:  access,
		TokenType:    "bearer",
		ExpiresIn:    accessTtl,
		RefreshToken: refresh,
	}

	return auth, nil
}
