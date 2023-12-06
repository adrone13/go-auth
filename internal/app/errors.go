package app

type InvalidPasswordError struct{}

func (e *InvalidPasswordError) Error() string {
	return "invalid password"
}

type UserNotFoundError struct{}

func (e *UserNotFoundError) Error() string {
	return "User not found"
}
