package app

import "context"

type GetUserUseCase struct {
	UserRepository UserRepository
}

func (u *GetUserUseCase) Execute(ctx context.Context, id UserId) (*User, error) {
	return u.UserRepository.FindById(ctx, id)
}
