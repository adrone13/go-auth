package app

import (
	"auth/internal/crypto"
	"context"
	"errors"
)

type SignUpUseCase struct {
	UserRepo UserRepository
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

	user := &User{
		FullName:   i.FullName,
		Email:      i.Email,
		Password:   hashed,
		IsVerified: false,
	}

	err = u.UserRepo.Insert(ctx, user)
	if err != nil {
		return errors.New("failed to create user")
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
