package app

import (
	"auth/internal/app/common/repositories"
	"auth/internal/app/users"
	"context"
)

type GetUserUseCase struct {
	UserRepository repositories.UserRepository
}

func (u *GetUserUseCase) Execute(ctx context.Context, id users.UserId) (*users.User, error) {
	return u.UserRepository.FindById(ctx, id)
}
