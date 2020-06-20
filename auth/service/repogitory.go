package service

import (
	"context"
	"time"

	"github.com/x-color/calendar/auth/model"
)

type Repogitory interface {
	User() UserRepogitory
	Session() SessionRepogitory
}

type UserRepogitory interface {
	FindByName(ctx context.Context, name string) (UserData, error)
	Create(ctx context.Context, user UserData) error
}

type SessionRepogitory interface {
	Find(ctx context.Context, id string) (SessionData, error)
	Create(ctx context.Context, session SessionData) error
	Delete(ctx context.Context, id string) error
}

type UserData struct {
	ID       string
	Name     string
	Password string
}

func newUserData(user model.User) UserData {
	return UserData{
		ID:       user.ID,
		Name:     user.Name,
		Password: user.Password,
	}
}

func (u *UserData) model() model.User {
	return model.User{
		ID:       u.ID,
		Name:     u.Name,
		Password: u.Password,
	}
}

type SessionData struct {
	ID      string
	UserID  string
	Expires int64
}

func newSessionData(s model.Session) SessionData {
	return SessionData{
		ID:      s.ID,
		UserID:  s.UserID,
		Expires: s.Expires.Unix(),
	}
}

func (s *SessionData) model() model.Session {
	return model.Session{
		ID:      s.ID,
		UserID:  s.UserID,
		Expires: time.Unix(s.Expires, 0),
	}
}
