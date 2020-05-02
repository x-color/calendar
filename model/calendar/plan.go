package calendar

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

func NewPeriod(begin, end time.Time) Period {
	return Period{
		IsAllDay: false,
		Begin:    begin,
		End:      end,
	}
}

type Plan struct {
	ID         string
	CalendarID string
	UserID     string
	Name       string
	Memo       string
	Color      Color
	Private    bool
	Period     Period
}

func NewPlan(calendarID, userID, name, memo string, color Color, pricate bool, period Period) Plan {
	return Plan{
		ID:         uuid.New().String(),
		CalendarID: calendarID,
		UserID:     userID,
		Name:       name,
		Memo:       memo,
		Color:      color,
		Private:    true,
		Period:     period,
	}
}
