package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/x-color/calendar/calendar/model"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
)

func (s *Service) RegisterUser(ctx context.Context, userID string) (model.User, error) {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	user, err := s.registerUser(ctx, userID)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Register user(%v)", user.ID))
	}

	return user, err
}

func (s *Service) registerUser(ctx context.Context, id string) (model.User, error) {
	if id == "" {
		return model.User{}, cerror.NewInvalidContentError(
			nil,
			"id is empty",
		)
	}

	user := model.NewUser(id)
	err := s.repo.User().Create(ctx, newUserData(user))
	if err != nil && !errors.Is(err, cerror.ErrDuplication) {
		return model.User{}, err
	}

	return user, nil
}

func (s *Service) CheckRegistration(ctx context.Context, userID string) error {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

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

	_, err := s.repo.User().Find(ctx, id)
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
