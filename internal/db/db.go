package db

import (
	"auth/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

var conn *pgxpool.Pool

func init() {
	var err error

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.Values.DbUser,
		config.Values.DbPassword,
		config.Values.DbHost,
		config.Values.DbPort,
		config.Values.DbName,
	)

	conn, err = pgxpool.New(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
}

func Ping(ctx context.Context) error {
	return conn.Ping(ctx)
}
