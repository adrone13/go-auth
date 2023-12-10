package server

import (
	"auth/internal/config"
	"auth/internal/jwt"
	"errors"
	"net/http"
	"strings"
)

/*
TODO:
Think about putting authentication into middleware

https://go.dev/blog/context#TOC_3.1.
https://stackoverflow.com/questions/51224251/authentication-using-request-vs-context
*/
func authenticate(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	h := r.Header.Get("Authorization")
	parts := strings.Split(h, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("unauthorized")
	}

	token, err := jwt.Parse(parts[1])
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	secret := config.GetString("JWT_SECRET")
	if !token.IsValid(secret) {
		return nil, errors.New("unauthorized")
	}

	return token, nil
}
