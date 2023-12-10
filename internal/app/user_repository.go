package app

type UserRepository interface {
	Insert(*User) (UserId, error)
	FindById(UserId) *User
	FindByEmail(email string) *User
}
