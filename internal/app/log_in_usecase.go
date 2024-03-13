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

type Auth struct {
	Token string `json:"token"`
	// RefreshToken string
}

type LogInUseCase struct {
	UserRepository UserRepository
}

func (c *LogInUseCase) Execute(ctx context.Context, cred Credentials) (*Auth, error) {
	user, err := c.UserRepository.FindByEmail(ctx, cred.Email)
	if err != nil {
		return nil, err
	}

	if crypto.ComparePasswordHash(cred.Password, user.Password) {
		return nil, &InvalidPasswordError{}
	}

	secret := config.Values.JwtSecret
	ttl := config.Values.JwtTtl

	claims := jwt.Claims{
		Issuer:     "auth",
		Expiration: time.Now().Add(time.Second * time.Duration(ttl)).Unix(),
		Audience:   "todo",
		Subject:    string(user.Id),
		Name:       user.FullName,
		Roles:      []string{"TODO"},
	}

	token := jwt.Sign(claims, secret)

	auth := new(Auth)
	auth.Token = token

	return auth, nil
}
