package service

import (
	"context"
	"fmt"

	"github.com/x-color/calendar/calendar/model"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
)

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
	for i, uid := range cal.Shares {
		if userID == uid {
			if i == len(cal.Shares)-1 {
				cal.Shares = cal.Shares[:i]
			} else {
				cal.Shares = append(cal.Shares[:i], cal.Shares[i+1:]...)
			}
			err := s.repo.Calendar().Update(ctx, newCalendarData(cal))
			if err != nil {
				return err
			}
			return nil
		}
	}

	return cerror.NewAuthorizationError(
		nil,
		fmt.Sprintf("user(%v) does not permit to delete calendar(%v)", userID, cal.ID),
	)
}

func (s *Service) ChangeCalendar(ctx context.Context, userID string, cal model.Calendar) error {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	err := s.changeCalendar(ctx, userID, cal)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Change calendar(%v)", cal.ID))
	}

	return err
}

func (s *Service) changeCalendar(ctx context.Context, userID string, cal model.Calendar) error {
	if cal.ID == "" {
		return cerror.NewInvalidContentError(
			nil,
			"id is empty",
		)
	}

	if cal.Name == "" {
		return cerror.NewInvalidContentError(
			nil,
			"name is empty",
		)
	}

	in := false
	for _, uid := range cal.Shares {
		if userID == uid {
			in = true
			break
		}
	}
	if !in {
		return cerror.NewInvalidContentError(
			nil,
			"owner is not in shares",
		)
	}

	c, err := s.repo.Calendar().Find(ctx, cal.ID)
	if err != nil {
		return err
	}

	if userID != c.UserID {
		return cerror.NewAuthorizationError(
			nil,
			fmt.Sprintf("user(%v) does not permit to change calendar(%v)", userID, cal.ID),
		)
	}

	return s.repo.Calendar().Update(ctx, newCalendarData(cal))
}
