package calendar

import (
	"context"
	"fmt"

	"github.com/x-color/calendar/model/calendar"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
	"github.com/x-color/calendar/service"
)

type Repogitory interface {
	// TxBegin()
	// TxCommit()
	// TxRollback()
	Calendar() CalendarRepogitory
}

type CalendarRepogitory interface {
	Create(ctx context.Context, cal calendar.Calendar) error
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

func (s *Service) MakeCalendar(ctx context.Context, name, color string) (calendar.Calendar, error) {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	userID := ctx.Value(cctx.UserIDKey).(string)
	cal, err := s.makeCalendar(ctx, userID, name, color)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Make calendar(%v)", cal.ID))
	}

	return cal, err
}

func (s *Service) makeCalendar(ctx context.Context, userID, name, color string) (calendar.Calendar, error) {
	if name == "" {
		return calendar.Calendar{}, cerror.NewInvalidContentError(
			nil,
			"name is empty",
		)
	}
	c, err := calendar.ConvertToColor(color)
	if err != nil {
		return calendar.Calendar{}, err
	}

	cal := calendar.NewCalendar(userID, name, c)

	err = s.repo.Calendar().Create(ctx, cal)
	if err != nil {
		return calendar.Calendar{}, err
	}
	return cal, nil
}
