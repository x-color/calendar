package service

import (
	"context"

	"github.com/x-color/calendar/calendar/model"
	"github.com/x-color/calendar/logging"
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
	Create(ctx context.Context, cal model.Calendar) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, cal model.Calendar) error
	Find(ctx context.Context, id string) (model.Calendar, error)
}

type PlanRepogitory interface {
	Create(ctx context.Context, plan model.Plan) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, plan model.Plan) error
	Find(ctx context.Context, id string) (model.Plan, error)
}

type UserRepogitory interface {
	Create(ctx context.Context, user model.User) error
	Find(ctx context.Context, id string) (model.User, error)
}

type Service struct {
	repo Repogitory
	log  logging.Logger
}

func NewService(repo Repogitory, log logging.Logger) Service {
	return Service{
		repo: repo,
		log:  log,
	}
}
