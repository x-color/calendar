package calendar

import (
	"context"

	"github.com/x-color/calendar/model/calendar"
	"github.com/x-color/calendar/service"
)

type Repogitory interface {
	// TxBegin()
	// TxCommit()
	// TxRollback()
	Calendar() CalendarRepogitory
	Plan() PlanRepogitory
	CalUser() UserRepogitory
}

type CalendarRepogitory interface {
	Create(ctx context.Context, cal calendar.Calendar) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, cal calendar.Calendar) error
	Find(ctx context.Context, id string) (calendar.Calendar, error)
}

type PlanRepogitory interface {
	Create(ctx context.Context, plan calendar.Plan) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, plan calendar.Plan) error
	Find(ctx context.Context, id string) (calendar.Plan, error)
}

type UserRepogitory interface {
	Create(ctx context.Context, user calendar.User) error
	Find(ctx context.Context, id string) (calendar.User, error)
}

type Service struct {
	repo Repogitory
	log  service.Logger
}

func NewService(repo Repogitory, log service.Logger) Service {
	return Service{
		repo: repo,
		log:  log,
	}
}
