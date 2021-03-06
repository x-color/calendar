package model

import "github.com/google/uuid"

type User struct {
	ID       string
	Name     string
	Password string
}

func NewUser(name string, password string) User {
	return User{
		ID:       uuid.New().String(),
		Name:     name,
		Password: password,
	}
}
