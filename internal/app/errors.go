package app

import "fmt"

type InvalidPasswordError struct{}

func (e *InvalidPasswordError) Error() string {
	return "invalid password"
}

type UserNotFoundError struct {
	Criteria string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user not found by criteria: %s", e.Criteria)
}
