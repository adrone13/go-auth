package app

import (
	"auth/internal/app/common/repositories"
	"auth/internal/app/sessions"
	"auth/internal/app/tokens"
	"auth/internal/app/users"
	"auth/internal/config"
	"auth/internal/logger"
	"auth/internal/utils/arrays"
	"context"
	"errors"
	jwt "github.com/adrone13/gojwt"
)

type RefreshAuthUseCase struct {
	UserRepository     repositories.UserRepository
	SessionsRepository repositories.SessionsRepository
}

func (u *RefreshAuthUseCase) Execute(ctx context.Context, refreshToken string) (*Auth, error) {
	token, err := jwt.ParseCustomClaims(refreshToken, tokens.RefreshTokenClaims{})
	if err != nil {
		return nil, err
	}
	secret := config.Values.JwtSecret
	if !token.IsValid(secret) {
		return nil, errors.New("invalid_request")
	}

	user, err := u.UserRepository.FindById(ctx, users.UserId(token.Claims.Subject))
	if err != nil {
		logger.Error.Printf("User (id: %s) does not exist", token.Claims.Subject)

		return nil, err
	}
	session, err := u.SessionsRepository.FindById(ctx, sessions.SessionId(token.Claims.SessionId))
	if err != nil {
		logger.Error.Printf("Session (id: %s) does not exist", token.Claims.SessionId)

		return nil, err
	}
	if session.Expired() {
		logger.Error.Println("Logging out. Session has expired")

		// Logout user if session expired
		err = u.SessionsRepository.Delete(ctx, session.Id)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("access_denied")
	}

	if !arrays.Contains(session.RefreshTokens, refreshToken) {
		return nil, errors.New("invalid_request")
	}

	currentRefresh := arrays.Last(session.RefreshTokens)
	if refreshToken != currentRefresh {
		logger.Error.Printf("Logging out. Invalidated token provided: %s", refreshToken)

		// Logout user if invalidated refresh token received
		err = u.SessionsRepository.Delete(ctx, session.Id)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("access_denied")
	}

	access := tokens.CreateAccessToken(user)
	refresh := tokens.CreateRefreshToken(user, session)

	session.AddRefreshToken(refresh)
	err = u.SessionsRepository.Update(ctx, session)
	if err != nil {
		return nil, err
	}

	auth := &Auth{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "bearer",
		ExpiresIn:    config.Values.AccessTokenTtl,
	}

	return auth, nil
}
