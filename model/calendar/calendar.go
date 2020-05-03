package calendar

import (
	"github.com/google/uuid"
)

type Calendar struct {
	ID     string
	UserID string
	Name   string
	Color  Color
	Plans  []Plan
	Shares []string
}

func NewCalendar(userID, name string, color Color) Calendar {
	return Calendar{
		ID:     uuid.New().String(),
		UserID: userID,
		Name:   name,
		Color:  color,
		Plans:  []Plan{},
		Shares: []string{userID},
	}
}
