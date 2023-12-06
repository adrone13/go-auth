package database

import "auth/internal/user"

type UserRepositoryImpl struct{}

/* 
	TODO:
	For now it's just a mock
	Need to implement integration with Postgres
*/
func (r *UserRepositoryImpl) FindByEmail(email string) *user.User {
	if email != "alex@gmail.com" {
		return nil
	}

	return &user.User{
		Id:       "uuid",
		FullName: "Alex The Mad",
		Email:    "alex@gmail.com",
		Password: "for now not hashed pass",
		Secret:   "",
	}
}
