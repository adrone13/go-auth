package server

import (
	"auth/internal/config"
	"errors"
	"github.com/adrone13/gojwt"
	"net/http"
	"strings"
	"time"
)

// Think about putting authentication into middleware:
// https://go.dev/blog/context#TOC_3.1.
// https://stackoverflow.com/questions/51224251/authentication-using-request-vs-context
func authenticate(_ http.ResponseWriter, r *http.Request) (*jwt.Token[jwt.Claims], error) {
	h := r.Header.Get("Authorization")
	parts := strings.Split(h, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("invalid token format")
	}

	token, err := jwt.Parse(parts[1])
	if err != nil {
		return nil, err
	}
	secret := config.Values.JwtSecret
	if !token.IsValid(secret) {
		return nil, errors.New("unauthorized")
	}
	now := time.Now().Unix()
	if token.Claims.Expiration < now {
		return nil, errors.New("token expired")
	}

	return token, nil
}
