package app

import "context"

type UserRepository interface {
	Insert(ctx context.Context, u *User) error
	FindById(ctx context.Context, u UserId) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}
