package calendar

import (
	"context"
	"errors"
	"fmt"

	"github.com/x-color/calendar/model/calendar"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
)

func (s *Service) Schedule(ctx context.Context, planPram calendar.Plan) (calendar.Plan, error) {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	userID := ctx.Value(cctx.UserIDKey).(string)
	planPram.UserID = userID

	plan, err := s.schedule(ctx, planPram)
	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Schedule plan(%v)", plan.ID))
	}

	return plan, err
}

func (s *Service) schedule(ctx context.Context, planPram calendar.Plan) (calendar.Plan, error) {
	if planPram.Name == "" || planPram.CalendarID == "" || len(planPram.Shares) == 0 {
		return calendar.Plan{}, cerror.NewInvalidContentError(
			nil,
			"some contents are empty",
		)
	}

	if !planPram.Period.IsAllDay && !planPram.Period.Begin.Before(planPram.Period.End) {
		return calendar.Plan{}, cerror.NewInvalidContentError(
			nil,
			"invalid period",
		)
	}

	cal, err := s.repo.Calendar().Find(ctx, planPram.CalendarID)
	if errors.Is(err, cerror.ErrNotFound) {
		return calendar.Plan{}, cerror.NewInvalidContentError(
			nil,
			"invalid calendar id",
		)
	} else if err != nil {
		return calendar.Plan{}, err
	}

	in := false
	for _, u := range cal.Shares {
		if u == planPram.UserID {
			in = true
			break
		}
	}
	if !in {
		return calendar.Plan{}, cerror.NewAuthorizationError(
			nil,
			fmt.Sprintf("user(%v) does not permit to access the calendar(%v)", planPram.UserID, planPram.CalendarID),
		)
	}

	// TODO: Checking user id in planPram.Shares

	plan := calendar.NewPlan(
		planPram.CalendarID,
		planPram.UserID,
		planPram.Name,
		planPram.Memo,
		planPram.Color,
		planPram.Private,
		planPram.Shares,
		planPram.Period,
	)

	err = s.repo.Plan().Create(ctx, plan)
	if err != nil {
		return calendar.Plan{}, err
	}

	return plan, nil
}
