package mocks

import (
	"auth/internal/app"
	"auth/internal/app/users"
	"context"
	"fmt"
	"log"
	"time"
)

type UserRepositoryMock struct{}

func (r *UserRepositoryMock) Insert(_ context.Context, _ *users.User) error {
	return nil
}

func (r *UserRepositoryMock) FindById(_ context.Context, id users.UserId) (*users.User, error) {
	if id != "user_uuid" {
		return nil, &app.UserNotFoundError{Criteria: fmt.Sprintf("Id = %s", id)}
	}

	createdAt, err := time.Parse(time.RFC3339, "2024-04-08T12:00:00.000Z")
	updatedAt, err := time.Parse(time.RFC3339, "2024-04-08T12:05:00.000Z")
	if err != nil {
		log.Fatalln(err)
	}

	return &users.User{
		Id:         id,
		FullName:   "Full Name",
		Email:      "fullname@mail.com",
		Password:   "hash_pass",
		IsVerified: true,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func (r *UserRepositoryMock) FindByEmail(_ context.Context, _ string) (*users.User, error) {
	return nil, nil
}
