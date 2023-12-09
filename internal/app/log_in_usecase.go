package app

import (
	"fmt"
	"time"

	"auth/internal/config"
	"auth/internal/jwt"
)

/*
	TODO:
	- move JWT creation into separate library (WIP)
	- add DB implementation
	- implement registration
	- implement refresh token storage in sessions
	- implement JWT refreshing
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

// Controller constructor Constructor
type LogInUseCase struct {
	UserRepository UserRepository
}

func (c *LogInUseCase) Execute(cred Credentials) (*Auth, error) {
	user := c.UserRepository.FindByEmail(cred.Email)
	if user == nil {
		return nil, &UserNotFoundError{}
	}

	fmt.Printf("User: %+v\n", cred)

	if cred.Password != user.Password {
		return nil, &InvalidPasswordError{}
	}

	secret := config.GetString("JWT_SECRET")
	ttl := config.GetInt("JWT_TTL")

	claims := jwt.Claims{
		Iss:   "auth",
		Exp:   time.Now().Add(time.Second * time.Duration(ttl)).Unix(),
		Aud:   "todo",
		Sub:   string(user.Id),
		Name:  user.FullName,
		Roles: []string{"TODO"},
	}

	// claims := jwt.BuildClaims().
	// 	AddIss("auth").
	// 	AddExp(time.Now().Add(time.Second * time.Duration(ttl)).Unix()).
	// 	AddAud("todo").
	// 	AddSub(user.Id).
	// 	AddName(user.FullName).
	// 	AddRoles([]string{"TODO"})

	token := jwt.Sign(claims, secret)

	auth := new(Auth)
	auth.Token = token

	return auth, nil
}
