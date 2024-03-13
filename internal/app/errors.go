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
	return fmt.Sprintf("User not found. Criteria: %s", e.Criteria)
}
