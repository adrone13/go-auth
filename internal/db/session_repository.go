package db

import (
	"auth/internal/app"
	"auth/internal/logger"
	"auth/internal/utils/arrays"
	"context"
	"log"
)

type SessionRepository struct{}

func (r *SessionRepository) Insert(ctx context.Context, s *app.Session) error {
	if s.Id != "" {
		log.Panicf("failed to insert. session already has id: %s", s.Id)
	}
	if len(s.RefreshTokens) != 1 {
		log.Panicf("invalid number of refresh tokens. expected 1, received %d\n", len(s.RefreshTokens))
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sessionQuery := "INSERT INTO sessions (user_id, created_at, updated_at) VALUES ($1, $2, $3) RETURNING ID"
	err = tx.QueryRow(ctx, sessionQuery, s.UserId, s.CreatedAt, s.UpdatedAt).Scan(&s.Id)
	if err != nil {
		return err
	}

	refreshTokenQuery := "INSERT INTO refresh_tokens (session_id, token, created_at) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, refreshTokenQuery, s.Id, s.RefreshTokens[0], s.CreatedAt)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) Update(ctx context.Context, s *app.Session) error {
	if s.Id == "" {
		log.Panicln("failed to update. session missing id")
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sessionQuery := "UPDATE sessions SET updated_at = $1 WHERE id = $2"
	_, err = tx.Exec(ctx, sessionQuery, s.UpdatedAt, s.Id)
	if err != nil {
		logger.Error(err)

		return err
	}

	refreshTokenQuery := "INSERT INTO refresh_tokens (session_id, token, created_at) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, refreshTokenQuery, s.Id, arrays.Last(s.RefreshTokens), s.UpdatedAt)
	if err != nil {
		logger.Error(err)

		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) FindByUserId(ctx context.Context, userId app.UserId) (*app.Session, error) {
	q := `select rt.token, s.id, s.user_id, s.created_at, s.updated_at from sessions s
	join refresh_tokens rt on s.id = rt.session_id
	where s.user_id = $1
	order by rt.created_at`

	rows, err := conn.Query(ctx, q, userId)
	if err != nil {
		log.Fatalf("Failed to query: %s", err)
	}

	var session *app.Session
	for rows.Next() {
		var token string

		if session == nil {
			session = new(app.Session)
			err = rows.Scan(&token, &session.Id, &session.UserId, &session.CreatedAt, &session.UpdatedAt)
			if err != nil {
				return nil, err
			}

			session.RefreshTokens = []string{token}

			continue
		}

		err = rows.Scan(&token, &session.Id, &session.UserId, &session.CreatedAt, &session.UpdatedAt)
		if err != nil {
			return nil, err
		}

		session.RefreshTokens = append(session.RefreshTokens, token)
	}

	return session, nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionId app.SessionId) error {
	_, err := conn.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, sessionId)
	if err != nil {
		return err
	}

	return nil
}
