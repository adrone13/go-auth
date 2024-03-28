package tokens

import (
	"auth/internal/app/sessions"
	"auth/internal/app/users"
	"auth/internal/config"
	"auth/internal/logger"
	"fmt"
	"github.com/adrone13/gojwt"
	"time"
)

type RefreshTokenClaims struct {
	Expiration int64  `json:"exp"`
	Subject    string `json:"sub"`
	SessionId  string `json:"session_id"`
}

func CreateAccessToken(user *users.User) string {
	secret := config.Values.JwtSecret
	accessTtl := config.Values.AccessTokenTtl

	exp := time.Now().Add(time.Second * time.Duration(accessTtl))

	return jwt.Sign(jwt.Claims{
		Issuer:     "auth",
		Expiration: exp.Unix(),
		Audience:   "todo",
		Subject:    string(user.Id),
		Name:       user.FullName,
		Roles:      []string{"TODO"},
	}, secret)
}

func CreateRefreshToken(user *users.User, session *sessions.Session) string {
	secret := config.Values.JwtSecret
	refreshTtl := config.Values.RefreshTokenAbsoluteTtl

	now := time.Now()
	logger.Warn(fmt.Sprintf("Refresh token created: %v", now))

	return jwt.Sign(RefreshTokenClaims{
		Expiration: now.Add(time.Second * time.Duration(refreshTtl)).Unix(),
		Subject:    string(user.Id),
		SessionId:  string(session.Id),
	}, secret)
}
