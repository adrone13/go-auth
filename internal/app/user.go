package app

type UserId string

type User struct {
	Id         UserId
	FullName   string
	Email      string
	Password   string
	IsVerified bool
}
