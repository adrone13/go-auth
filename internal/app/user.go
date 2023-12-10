package app

type UserId string

type User struct {
	Id         UserId `json:"id"`
	FullName   string `json:"fullName"`
	Email      string `json:"-"`
	Password   string `json:"-"`
	IsVerified bool   `json:"isVerified"`
}
