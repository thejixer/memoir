package models

import "time"

type UserRepository interface {
	Create(name, email, password string, isEmailVerified bool) (*User, error)
	FindById(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	VerifyEmail(email string) error
}

type User struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	Password        string    `json:"password"`
	CreatedAt       time.Time `json:"createdAt"`
}

type UserDto struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type LL_UserDto struct {
	Total  int       `json:"total"`
	Result []UserDto `json:"result"`
}
