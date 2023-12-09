package db

import (
	"auth/internal/app"
	"errors"
	"fmt"
)

type UserRepositoryImpl struct{}

/*
	TODO:
	For now it's just a mock
	Need to implement integration with Postgres
*/

func (r *UserRepositoryImpl) Insert(u *app.User) (app.UserId, error) {
	if u.Id != "" {
		return "", errors.New("user already exists")
	}

	fmt.Printf("Saved user: %+v\n", u)

	return "uuid", nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) *app.User {
	if email != "alex@gmail.com" {
		return nil
	}

	return &app.User{
		Id:         "uuid",
		FullName:   "Alex The Mad",
		Email:      "alex@gmail.com",
		Password:   "for now not hashed pass",
		IsVerified: false,
	}
}
