package app

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type SignUpUseCase struct {
	UserRepository UserRepository
}

type SignUpInput struct {
	FullName string
	Email    string
	Password string
}

func (u *SignUpUseCase) Execute(i SignUpInput) error {
	err := validateInput(i)
	if err != nil {
		return err
	}

	hashed, err := hashPassword(i.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := &User{
		FullName:   i.FullName,
		Email:      i.Email,
		Password:   hashed,
		IsVerified: false,
	}

	_, err = u.UserRepository.Insert(user)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
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
