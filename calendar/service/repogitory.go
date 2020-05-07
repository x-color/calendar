package service

import (
	"context"
	"time"

	"github.com/x-color/calendar/calendar/model"
)

type Repogitory interface {
	// TxBegin()
	// TxCommit()
	// TxRollback()
	Calendar() CalendarRepogitory
	Plan() PlanRepogitory
	User() UserRepogitory
}

type CalendarRepogitory interface {
	Create(ctx context.Context, cal CalendarData) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, cal CalendarData) error
	Find(ctx context.Context, id string) (CalendarData, error)
}

type PlanRepogitory interface {
	Create(ctx context.Context, plan PlanData) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, plan PlanData) error
	Find(ctx context.Context, id string) (PlanData, error)
}

type UserRepogitory interface {
	Create(ctx context.Context, user UserData) error
	Find(ctx context.Context, id string) (UserData, error)
}

type UserData struct {
	ID string
}

func newUserData(user model.User) UserData {
	return UserData{
		ID: user.ID,
	}
}

func (u *UserData) model() model.User {
	return model.User{
		ID: u.ID,
	}
}

type CalendarData struct {
	ID     string
	UserID string
	Name   string
	Color  string
	Shares []string
}

func newCalendarData(cal model.Calendar) CalendarData {
	return CalendarData{
		ID:     cal.ID,
		UserID: cal.UserID,
		Name:   cal.Name,
		Color:  string(cal.Color),
		Shares: cal.Shares,
	}
}

func (c *CalendarData) model() model.Calendar {
	return model.Calendar{
		ID:     c.ID,
		UserID: c.UserID,
		Name:   c.Name,
		Color:  model.Color(c.Color),
		Plans:  []model.Plan{},
		Shares: c.Shares,
	}
}

type PlanData struct {
	ID         string
	CalendarID string
	UserID     string
	Name       string
	Memo       string
	Color      string
	Private    bool
	Shares     []string
	IsAllDay   bool
	Begin      int64
	End        int64
}

func newPlanData(plan model.Plan) PlanData {
	return PlanData{
		ID:         plan.ID,
		CalendarID: plan.CalendarID,
		UserID:     plan.UserID,
		Name:       plan.Name,
		Memo:       plan.Memo,
		Color:      string(plan.Color),
		Private:    plan.Private,
		Shares:     plan.Shares,
		IsAllDay:   plan.Period.IsAllDay,
		Begin:      plan.Period.Begin.Unix(),
		End:        plan.Period.End.Unix(),
	}
}

func (p *PlanData) model() model.Plan {
	return model.Plan{
		ID:         p.ID,
		CalendarID: p.CalendarID,
		UserID:     p.UserID,
		Name:       p.Name,
		Memo:       p.Memo,
		Color:      model.Color(p.Color),
		Private:    p.Private,
		Shares:     p.Shares,
		Period: model.Period{
			IsAllDay: p.IsAllDay,
			Begin:    time.Unix(p.Begin, 0),
			End:      time.Unix(p.End, 0),
		},
	}
}
