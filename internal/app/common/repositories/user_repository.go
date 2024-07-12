package repositories

import (
	"auth/internal/app/users"
	"context"
)

type UserRepository interface {
	Insert(ctx context.Context, u *users.User) error
	FindById(ctx context.Context, id users.UserId) (*users.User, error)
	FindByEmail(ctx context.Context, email string) (*users.User, error)
}
