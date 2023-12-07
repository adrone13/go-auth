package app

import (
	"fmt"
	"time"

	"auth/internal/config"
	"auth/internal/jwt"
	"auth/internal/user"
)

/*
	TODO:
	- add DB implementation
	- move JWT creation into separate library
	- implement registration
	- implement refresh token storage in sessions
	- implement JWT refreshing
	- implement logout and JWT invalidation (through cache)
*/

type UserRepository interface {
	FindByEmail(email string) *user.User
}

type Credentials struct {
	Email    string
	Password string
}

type Auth struct {
	Token string `json:"token"`
	// RefreshToken string
}

// Controller constructor Constructor
type LoginController struct {
	UserRepository UserRepository
}

func (c *LoginController) Execute(cred Credentials) (*Auth, error) {
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
	claims := jwt.BuildClaims().
		AddIss("auth").
		AddExp(time.Now().Add(time.Second * time.Duration(ttl)).Unix()).
		AddAud("todo").
		AddSub(user.Id).
		AddName(user.FullName).
		AddRoles([]string{"TODO"})

	token := jwt.Sign(*claims, secret)

	auth := new(Auth)
	auth.Token = token

	return auth, nil
}
