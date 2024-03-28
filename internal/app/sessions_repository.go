package app

import (
	"auth/internal/app/sessions"
	"auth/internal/app/users"
	"context"
)

type SessionsRepository interface {
	Insert(ctx context.Context, s *sessions.Session) error
	Update(ctx context.Context, s *sessions.Session) error
	FindById(ctx context.Context, id sessions.SessionId) (*sessions.Session, error)
	FindByUserId(ctx context.Context, userId users.UserId) (*sessions.Session, error)
	Delete(ctx context.Context, sessionId sessions.SessionId) error
}
