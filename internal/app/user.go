package app

import "time"

type UserId string

type User struct {
	Id         UserId `json:"id"`
	FullName   string `json:"fullName"`
	Email      string `json:"-"`
	Password   string `json:"-"`
	IsVerified bool   `json:"isVerified"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
