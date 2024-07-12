package server

import "context"

type Database interface {
	Ping(ctx context.Context) error
}
