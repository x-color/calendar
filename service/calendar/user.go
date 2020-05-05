package calendar

import (
	"context"
	"errors"
	"fmt"

	"github.com/x-color/calendar/model/calendar"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
)

func (s *Service) RegisterUser(ctx context.Context) (calendar.User, error) {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	userID := ctx.Value(cctx.UserIDKey).(string)
	user, err := s.registerUser(ctx, userID)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Register user(%v)", user.ID))
	}

	return user, err
}

func (s *Service) registerUser(ctx context.Context, id string) (calendar.User, error) {
	if id == "" {
		return calendar.User{}, cerror.NewInvalidContentError(
			nil,
			"id is empty",
		)
	}

	user := calendar.NewUser(id)
	err := s.repo.CalUser().Create(ctx, user)
	if err != nil && !errors.Is(err, cerror.ErrDuplication) {
		return calendar.User{}, err
	}

	return user, nil
}

func (s *Service) CheckRegistration(ctx context.Context) error {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	userID := ctx.Value(cctx.UserIDKey).(string)
	err := s.checkRegistration(ctx, userID)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Registered user(%v)", userID))
	}

	return err
}

func (s *Service) checkRegistration(ctx context.Context, id string) error {
	if id == "" {
		return cerror.NewInvalidContentError(
			nil,
			"id is empty",
		)
	}

	_, err := s.repo.CalUser().Find(ctx, id)
	if errors.Is(err, cerror.ErrNotFound) {
		return cerror.NewAuthorizationError(
			err,
			fmt.Sprintf("user(%v) is not registerd", id),
		)
	} else if err != nil {
		return err
	}

	return nil
}
