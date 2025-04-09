package model

import (
	"auth-service/app/internal/passhash"
	"time"
)

type User struct {
	ID        uint      `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserRegisterRequest struct {
	Username string `form:"username" validate:"required,min=5,max=15"`
	Email    string `form:"email"    validate:"required,email"`
	Password string `form:"password" validate:"required,min=8,max=64"`
}

func (u *UserRegisterRequest) ToUser() (*User, error) {
	hashPassword, err := passhash.Hash(u.Password)
	if err != nil {
		return nil, err
	}
	return &User{Username: u.Username, Email: u.Email, Password: hashPassword}, nil
}

type UserLoginRequest struct {
	Login    string `form:"login" validate:"login"`
	Password string `form:"password" validate:"required,min=8,max=64"`
}

type UserTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}
