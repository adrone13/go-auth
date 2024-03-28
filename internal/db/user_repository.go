package db

import (
	"auth/internal/app"
	"auth/internal/app/users"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

type UserRepository struct{}

func (r *UserRepository) Insert(ctx context.Context, u *users.User) error {
	if u.Id != "" {
		log.Panicln("user is already inserted")
	}

	q :=
		`INSERT INTO users (full_name, email, password, is_verified, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6) RETURNING id;
`
	err := conn.QueryRow(
		ctx,
		q,
		u.FullName,
		u.Email,
		u.Password,
		u.IsVerified,
		u.CreatedAt,
		u.UpdatedAt,
	).Scan(&u.Id)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindById(ctx context.Context, id users.UserId) (*users.User, error) {
	rows, err := conn.Query(ctx, "SELECT * FROM users WHERE id = $1 LIMIT 1", id)
	if err != nil {
		log.Fatalf("Failed to query: %s", err)
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[users.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &app.UserNotFoundError{Criteria: fmt.Sprintf("Id = %s", id)}
		}

		return nil, err
	}

	return u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	rows, err := conn.Query(ctx, "SELECT * FROM users WHERE email = $1 LIMIT 1", email)
	if err != nil {
		log.Fatalf("Failed to query: %s", err)
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[users.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &app.UserNotFoundError{Criteria: fmt.Sprintf("email = %s", email)}
		}

		return nil, err
	}

	return u, err
}

func (r *UserRepository) FindAll(ctx context.Context) []*users.User {
	rows, err := conn.Query(ctx, "SELECT * FROM users")
	if err != nil {
		log.Fatalf("Failed to execute query. Error: %s", err)
	}

	collectedUsers, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[users.User])
	if err != nil {
		log.Fatalf("Failed to collect rows. Error: %s", err)
	}

	return collectedUsers
}
