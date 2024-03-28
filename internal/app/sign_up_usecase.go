package app

import (
	"auth/internal/app/common/repositories"
	"auth/internal/app/users"
	"auth/internal/crypto"
	"context"
	"errors"
	"fmt"
)

type SignUpUseCase struct {
	UserRepo repositories.UserRepository
}

type SignUpInput struct {
	FullName string
	Email    string
	Password string
}

func (u *SignUpUseCase) Execute(ctx context.Context, i SignUpInput) error {
	err := validateInput(i)
	if err != nil {
		return err
	}

	hashed, err := crypto.HashPassword(i.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := users.New(i.FullName, i.Email, hashed)

	err = u.UserRepo.Insert(ctx, user)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create user. error: %s", err))
	}

	return nil
}

func validateInput(i SignUpInput) error {
	if i.FullName == "" {
		return errors.New("full name required")
	}
	if i.Email == "" {
		return errors.New("email required")
	}
	if i.Password == "" {
		return errors.New("password required")
	}
	// Use rainbow table
	if len(i.Password) < 8 {
		return errors.New("password not safe")
	}

	return nil
}
