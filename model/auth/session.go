package auth

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID      string
	UserID  string
	Expires time.Time
}

func NewSession(userID string, expired time.Time) Session {
	return Session{
		ID:      uuid.New().String(),
		UserID:  userID,
		Expires: expired,
	}
}
