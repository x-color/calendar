package model

import (
	"time"

	"github.com/google/uuid"
)

var AllDay = Period{IsAllDay: true}

type Period struct {
	IsAllDay bool
	Begin    time.Time
	End      time.Time
}

type Plan struct {
	ID         string
	CalendarID string
	UserID     string
	Name       string
	Memo       string
	Color      Color
	Private    bool
	Shares     []string
	Period     Period
}

func NewPlan(calendarID, userID, name, memo string, color Color, private bool, shares []string, period Period) Plan {
	return Plan{
		ID:         uuid.New().String(),
		CalendarID: calendarID,
		UserID:     userID,
		Name:       name,
		Memo:       memo,
		Color:      color,
		Private:    private,
		Shares:     shares,
		Period:     period,
	}
}
