package db

import (
	"auth/internal/app"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

type UserRepositoryImpl struct{}

func (r *UserRepositoryImpl) Insert(ctx context.Context, u *app.User) error {
	if u.Id != "" {
		log.Panicln("user is already inserted")
	}

	q := "INSERT INTO users (full_name, email, password, is_verified) VALUES($1, $2, $3, $4) RETURNING id;"
	err := conn.QueryRow(ctx, q, u.FullName, u.Email, u.Password, u.IsVerified).Scan(&u.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) FindById(ctx context.Context, id app.UserId) (*app.User, error) {
	rows, err := conn.Query(ctx, "SELECT * FROM users WHERE id = $1 LIMIT 1", id)
	if err != nil {
		log.Fatalf("Failed to query: %s", err)
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[app.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &app.UserNotFoundError{Criteria: fmt.Sprintf("Id = %s", id)}
		}

		return nil, err
	}

	return u, err
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*app.User, error) {
	rows, err := conn.Query(ctx, "SELECT * FROM users WHERE email = $1 LIMIT 1", email)
	if err != nil {
		log.Fatalf("Failed to query: %s", err)
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[app.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &app.UserNotFoundError{Criteria: fmt.Sprintf("email = %s", email)}
		}

		return nil, err
	}

	return u, err
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context) []*app.User {
	rows, err := conn.Query(ctx, "SELECT * FROM users")
	if err != nil {
		log.Fatalf("Failed to execute query. Error: %s", err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[app.User])
	if err != nil {
		log.Fatalf("Failed to collect rows. Error: %s", err)
	}

	return users
}
