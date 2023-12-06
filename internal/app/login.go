package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"auth/internal/user"
)

/* 
	TODO: 
	- move JWT creation into separate library
	- implement refresh token storage in sessions
	- implement JWT refreshing
	- implement logout and JWT invalidation (through cache)
*/

const (
	JWTDurationSeconds = 120
)

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

	JWT := newJWT(user)

	auth := new(Auth)
	auth.Token = JWT

	return auth, nil
}

func newJWT(u *user.User) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	claims := newClaims(u.Id, u.FullName, []string{"TODO"})
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		fmt.Println("failed to marshal claims")

		panic(err)
	}
	payload := base64.RawURLEncoding.EncodeToString([]byte(claimsJSON))

	signature := computeSignature(fmt.Sprintf("%s.%s", header, payload))
	encodedSignature := base64.RawURLEncoding.EncodeToString([]byte(signature))

	return fmt.Sprintf("%s.%s.%s", header, payload, encodedSignature)
}

func computeSignature(message string) []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Println("\"JWT_SECRET\" is not set up")

		panic("fai")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))

	return mac.Sum(nil)
}

func newClaims(sub string, name string, roles []string) *Claims {
	c := new(Claims)

	c.Iss = "auth"
	c.Exp = time.Now().Add(time.Second * JWTDurationSeconds).Unix()
	c.Sub = sub
	// c.Aud = "todo"
	c.Name = name
	c.Roles = roles

	return c
}

// https://www.iana.org/assignments/jwt/jwt.xhtml
type Claims struct {
	// Registered claims
	Iss string `json:"iss"` // Issuer (e.g. Auth service)
	Exp int64  `json:"exp"` // Expiration timestamp
	Aud string `json:"aud"` // Audience - service for which JWT is intended
	Sub string `json:"sub"` // Subject - user identity

	// Public claims
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}
