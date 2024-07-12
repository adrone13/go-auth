package db

import (
	"auth/internal/app/sessions"
	"auth/internal/app/users"
	"auth/internal/logger"
	"auth/internal/utils/arrays"
	"context"
	"log"
)

type SessionRepository struct{}

func (r *SessionRepository) Insert(ctx context.Context, s *sessions.Session) error {
	if s.Id != "" {
		log.Panicf("failed to insert. session already has id: %s", s.Id)
	}

	sessionQuery := "INSERT INTO sessions (user_id, created_at, updated_at) VALUES ($1, $2, $3) RETURNING ID"
	err := conn.QueryRow(ctx, sessionQuery, s.UserId, s.CreatedAt, s.UpdatedAt).Scan(&s.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) Update(ctx context.Context, s *sessions.Session) error {
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
		logger.Error.Println(err)

		return err
	}

	refreshTokenQuery := "INSERT INTO refresh_tokens (session_id, token, created_at) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, refreshTokenQuery, s.Id, arrays.Last(s.RefreshTokens), s.UpdatedAt)
	if err != nil {
		logger.Error.Println(err)

		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) FindById(ctx context.Context, id sessions.SessionId) (*sessions.Session, error) {
	q := `select rt.token, s.id, s.user_id, s.created_at, s.updated_at from sessions s
	join refresh_tokens rt on s.id = rt.session_id
	where s.id = $1
	order by rt.created_at`

	rows, err := conn.Query(ctx, q, id)
	if err != nil {
		log.Fatalf("Failed to query: %s", err)
	}

	var sess *sessions.Session
	for rows.Next() {
		var token string

		if sess == nil {
			sess = new(sessions.Session)
			err = rows.Scan(&token, &sess.Id, &sess.UserId, &sess.CreatedAt, &sess.UpdatedAt)
			if err != nil {
				return nil, err
			}

			sess.RefreshTokens = []string{token}

			continue
		}

		err = rows.Scan(&token, &sess.Id, &sess.UserId, &sess.CreatedAt, &sess.UpdatedAt)
		if err != nil {
			return nil, err
		}

		sess.RefreshTokens = append(sess.RefreshTokens, token)
	}

	return sess, nil
}

func (r *SessionRepository) FindByUserId(ctx context.Context, userId users.UserId) (*sessions.Session, error) {
	q := `select rt.token, s.id, s.user_id, s.created_at, s.updated_at from sessions s
	join refresh_tokens rt on s.id = rt.session_id
	where s.user_id = $1
	order by rt.created_at`

	rows, err := conn.Query(ctx, q, userId)
	if err != nil {
		log.Fatalf("Failed to query: %s", err)
	}

	var sess *sessions.Session
	for rows.Next() {
		var token string

		if sess == nil {
			sess = new(sessions.Session)
			err = rows.Scan(&token, &sess.Id, &sess.UserId, &sess.CreatedAt, &sess.UpdatedAt)
			if err != nil {
				return nil, err
			}

			sess.RefreshTokens = []string{token}

			continue
		}

		err = rows.Scan(&token, &sess.Id, &sess.UserId, &sess.CreatedAt, &sess.UpdatedAt)
		if err != nil {
			return nil, err
		}

		sess.RefreshTokens = append(sess.RefreshTokens, token)
	}

	return sess, nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionId sessions.SessionId) error {
	_, err := conn.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, sessionId)
	if err != nil {
		return err
	}

	return nil
}
