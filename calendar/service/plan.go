package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/x-color/calendar/calendar/model"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
	"github.com/x-color/slice/strs"
)

func (s *Service) Schedule(ctx context.Context, userID string, planPram model.Plan) (model.Plan, error) {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	planPram.UserID = userID

	plan, err := s.schedule(ctx, planPram)
	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Schedule plan(%v)", plan.ID))
	}

	return plan, err
}

func (s *Service) schedule(ctx context.Context, planPram model.Plan) (model.Plan, error) {
	if planPram.Name == "" || planPram.CalendarID == "" || len(planPram.Shares) == 0 {
		return model.Plan{}, cerror.NewInvalidContentError(
			nil,
			"some contents are empty",
		)
	}

	if !planPram.Period.IsAllDay && !planPram.Period.Begin.Before(planPram.Period.End) {
		return model.Plan{}, cerror.NewInvalidContentError(
			nil,
			"invalid period",
		)
	}

	cal, err := s.repo.Calendar().Find(ctx, planPram.CalendarID)
	if errors.Is(err, cerror.ErrNotFound) {
		return model.Plan{}, cerror.NewInvalidContentError(
			nil,
			"invalid calendar id",
		)
	} else if err != nil {
		return model.Plan{}, err
	}

	if !strs.Contains(cal.Shares, planPram.UserID) {
		return model.Plan{}, cerror.NewAuthorizationError(
			nil,
			fmt.Sprintf("user(%v) does not permit to access the calendar(%v)", planPram.UserID, planPram.CalendarID),
		)
	}

	for _, id := range planPram.Shares {
		cal, err := s.repo.Calendar().Find(ctx, id)
		if err != nil || !strs.Contains(cal.model().Shares, planPram.UserID) {
			return model.Plan{}, cerror.NewInvalidContentError(
				nil,
				"invalid calendar id in shares",
			)
		}
	}

	plan := model.NewPlan(
		planPram.CalendarID,
		planPram.UserID,
		planPram.Name,
		planPram.Memo,
		planPram.Color,
		planPram.Private,
		planPram.Shares,
		planPram.Period,
	)

	err = s.repo.Plan().Create(ctx, newPlanData(plan))
	if err != nil {
		return model.Plan{}, err
	}

	return plan, nil
}

func (s *Service) Unschedule(ctx context.Context, userID, calID, id string) error {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	err := s.unschedule(ctx, userID, calID, id)
	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Unschedule plan(%v)", id))
	}

	return err
}

func (s *Service) unschedule(ctx context.Context, userID, calID, id string) error {
	if calID == "" || id == "" {
		return cerror.NewInvalidContentError(
			nil,
			"some id are empty",
		)
	}

	plan, err := s.repo.Plan().Find(ctx, id)
	if err != nil {
		return err
	}

	// It changes not to share plan in the calendar if calID is not parent calendar for the plan.
	// If not, it deletes the plan in all calendars.
	if plan.UserID != userID || plan.CalendarID != calID {
		return s.unsharePlan(ctx, userID, calID, plan.model())
	}

	return s.repo.Plan().Delete(ctx, id)
}

func (s *Service) unsharePlan(ctx context.Context, userID, calID string, plan model.Plan) error {
	l, err := strs.RemoveE(plan.Shares, calID)
	if err != nil {
		return cerror.NewInvalidContentError(
			nil,
			"invalid calendar id",
		)
	}
	plan.Shares = l

	cal, err := s.repo.Calendar().Find(ctx, calID)
	if err != nil {
		return err
	}

	if cal.UserID != userID {
		return cerror.NewAuthorizationError(
			nil,
			fmt.Sprintf("user(%v) does not permit to access the plan(%v)", userID, plan.ID),
		)
	}

	return s.repo.Plan().Update(ctx, newPlanData(plan))
}
