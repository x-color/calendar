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
