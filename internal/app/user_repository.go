package app

type UserRepository interface {
	Insert(u *User) (UserId, error)
	FindByEmail(email string) *User
}
