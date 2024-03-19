package app

import "context"

type SessionsRepository interface {
	Insert(ctx context.Context, s *Session) error
	Update(ctx context.Context, s *Session) error
	FindByUserId(ctx context.Context, userId UserId) (*Session, error)
	Delete(ctx context.Context, sessionId SessionId) error
}
