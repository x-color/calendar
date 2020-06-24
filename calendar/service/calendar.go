package service

import (
	"context"
	"fmt"

	"github.com/x-color/calendar/calendar/model"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
	"github.com/x-color/slice/strs"
)

func (s *Service) GetCalendars(ctx context.Context, userID string) ([]model.Calendar, error) {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	cals, err := s.getCalendars(ctx, userID)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Get calendars for user(%v)", userID))
	}

	return cals, err
}

func (s *Service) getCalendars(ctx context.Context, userID string) ([]model.Calendar, error) {
	cl, err := s.repo.Calendar().FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	cals := make([]model.Calendar, len(cl))
	for i, cal := range cl {
		cals[i] = cal.model()
		pl, err := s.repo.Plan().FindByCalendarID(ctx, cal.ID)
		if err != nil {
			return nil, err
		}
		plans := make([]model.Plan, len(pl))
		for j, p := range pl {
			plan := p.model()
			if plan.Private && plan.UserID != userID {
				plan, _ = maskPlan(plan, cal.ID)
			}
			plans[j] = plan
		}
		cals[i].Plans = plans
	}

	return cals, nil
}

func (s *Service) MakeCalendar(ctx context.Context, userID, name, color string) (model.Calendar, error) {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	cal, err := s.makeCalendar(ctx, userID, name, color)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Make calendar(%v)", cal.ID))
	}

	return cal, err
}

func (s *Service) makeCalendar(ctx context.Context, userID, name, color string) (model.Calendar, error) {
	if name == "" {
		return model.Calendar{}, cerror.NewInvalidContentError(
			nil,
			"name is empty",
		)
	}
	c, err := model.ConvertToColor(color)
	if err != nil {
		return model.Calendar{}, err
	}

	cal := model.NewCalendar(userID, name, c)

	err = s.repo.Calendar().Create(ctx, newCalendarData(cal))
	if err != nil {
		return model.Calendar{}, err
	}
	return cal, nil
}

func (s *Service) RemoveCalendar(ctx context.Context, userID, id string) error {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	err := s.removeCalendar(ctx, userID, id)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Remove calendar(%v)", id))
	}

	return err
}

func (s *Service) removeCalendar(ctx context.Context, userID, id string) error {
	if id == "" {
		return cerror.NewInvalidContentError(
			nil,
			"id is empty",
		)
	}

	cal, err := s.repo.Calendar().Find(ctx, id)
	if err != nil {
		return err
	}

	// User is not owner of the model.
	if userID != cal.UserID {
		return s.unshareCalendar(ctx, userID, cal.model())
	}

	return s.repo.Calendar().Delete(ctx, id)

}

func (s *Service) unshareCalendar(ctx context.Context, userID string, cal model.Calendar) error {
	l, err := strs.RemoveE(cal.Shares, userID)
	if err != nil {
		return cerror.NewAuthorizationError(
			nil,
			fmt.Sprintf("user(%v) does not permit to delete calendar(%v)", userID, cal.ID),
		)
	}
	cal.Shares = l

	return s.repo.Calendar().Update(ctx, newCalendarData(cal))
}

func (s *Service) ChangeCalendar(ctx context.Context, userID string, calPram model.Calendar) error {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	err := s.changeCalendar(ctx, userID, calPram)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Change calendar(%v)", calPram.ID))
	}

	return err
}

func (s *Service) changeCalendar(ctx context.Context, userID string, calPram model.Calendar) error {
	if calPram.ID == "" {
		return cerror.NewInvalidContentError(
			nil,
			"id is empty",
		)
	}

	if calPram.Name == "" {
		return cerror.NewInvalidContentError(
			nil,
			"name is empty",
		)
	}

	if !strs.Contains(calPram.Shares, userID) {
		return cerror.NewInvalidContentError(
			nil,
			"owner is not in shares",
		)
	}

	for _, uid := range calPram.Shares {
		if _, err := s.repo.User().Find(ctx, uid); err != nil {
			return cerror.NewInvalidContentError(
				nil,
				"invalid user in shares",
			)
		}
	}

	c, err := s.repo.Calendar().Find(ctx, calPram.ID)
	if err != nil {
		return err
	}

	if userID != c.UserID {
		return cerror.NewAuthorizationError(
			nil,
			fmt.Sprintf("user(%v) does not permit to change calendar(%v)", userID, calPram.ID),
		)
	}

	calPram.UserID = c.UserID

	return s.repo.Calendar().Update(ctx, newCalendarData(calPram))
}

func maskPlan(plan model.Plan, calID string) (model.Plan, error) {
	if !strs.Contains(plan.Shares, calID) {
		return model.Plan{}, cerror.NewInvalidContentError(
			nil,
			fmt.Sprintf("publishing destination calendar(%v) is not shared", calID),
		)
	}
	plan.Name = ""
	plan.Memo = ""
	plan.Shares = []string{calID}
	return plan, nil
}
