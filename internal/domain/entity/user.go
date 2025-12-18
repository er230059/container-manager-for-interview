package entity

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrEmptyPassword = errors.New("password cannot be empty")

type User struct {
	ID       int64
	Username string
	Password string
}

func NewUser(id int64, username, plainPassword string) (*User, error) {
	if plainPassword == "" {
		return nil, ErrEmptyPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       id,
		Username: username,
		Password: string(hashedPassword),
	}, nil
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
