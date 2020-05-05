package model

type User struct {
	ID string
}

func NewUser(id string) User {
	return User{
		ID: id,
	}
}
