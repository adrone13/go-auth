package app

import (
	"auth/internal/app"
	"auth/internal/app/users"
	"auth/internal/config"
	"auth/tests/mocks"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestUserNotFound(t *testing.T) {
	config.Values = mocks.ConfigMock

	useCase := &app.GetUserUseCase{UserRepository: &mocks.UserRepositoryMock{}}

	_, err := useCase.Execute(context.Background(), "non_existent_id")
	if err == nil {
		t.Errorf("Expected to return error")
	}

	expected := &app.UserNotFoundError{Criteria: fmt.Sprintf("Id = %s", "non_existent_id")}
	if !errors.As(err, &expected) {
		t.Errorf("Expected error to be of type app.UserNotFoundError")
	}

	if err.Error() != expected.Error() {
		t.Errorf("Expected error message to be user not found by criteria: Id = non_existent_id")
	}
}

func TestUserFound(t *testing.T) {
	config.Values = mocks.ConfigMock

	useCase := &app.GetUserUseCase{UserRepository: &mocks.UserRepositoryMock{}}

	receivedUser, err := useCase.Execute(context.Background(), "user_uuid")
	if err != nil {
		t.Errorf("Expected to return user, go error instead")
	}

	createdAt, err := time.Parse(time.RFC3339, "2024-04-08T12:00:00.000Z")
	updatedAt, err := time.Parse(time.RFC3339, "2024-04-08T12:05:00.000Z")

	expectedUser := &users.User{
		Id:         "user_uuid",
		FullName:   "Full Name",
		Email:      "fullname@mail.com",
		Password:   "hash_pass",
		IsVerified: true,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	if *receivedUser != *expectedUser {
		t.Errorf("Received invalid user")
	}
}
