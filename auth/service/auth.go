package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/x-color/calendar/auth/model"
	"github.com/x-color/calendar/logging"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
)

type Repogitory interface {
	User() UserRepogitory
	Session() SessionRepogitory
}

type UserRepogitory interface {
	FindByName(ctx context.Context, name string) (model.User, error)
	Create(ctx context.Context, user model.User) error
}

type SessionRepogitory interface {
	Find(ctx context.Context, id string) (model.Session, error)
	FindByUserID(ctx context.Context, userID string) (model.Session, error)
	Create(ctx context.Context, session model.Session) error
	Delete(ctx context.Context, id string) error
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

func (s *Service) Signup(ctx context.Context, name, password string) (model.User, error) {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	user, err := s.signup(ctx, name, password)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Sign up user(%v)", name))
	}

	return user, err
}

func (s *Service) signup(ctx context.Context, name, password string) (model.User, error) {
	if err := validateSigninInfo(name, password); err != nil {
		return model.User{}, err
	}

	_, err := s.repo.User().FindByName(ctx, name)
	if err == nil {
		return model.User{}, cerror.NewDuplicationError(
			nil,
			fmt.Sprintf("user(%v) already exists", name),
		)
	}
	if !errors.Is(err, cerror.ErrNotFound) {
		return model.User{}, err
	}

	hash, err := passwordHash(password)
	if err != nil {
		return model.User{}, cerror.NewInternalError(
			err,
			"failed to hash password",
		)
	}

	user := model.NewUser(name, hash)
	err = s.repo.User().Create(ctx, user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *Service) Signin(ctx context.Context, name, password string) (model.Session, error) {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	session, err := s.signin(ctx, name, password)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Sign in user(%v)", name))
	}

	return session, err
}

func (s *Service) signin(ctx context.Context, name, password string) (model.Session, error) {
	if err := validateSigninInfo(name, password); err != nil {
		return model.Session{}, err
	}

	user, err := s.repo.User().FindByName(ctx, name)
	if errors.Is(err, cerror.ErrNotFound) {
		return model.Session{}, cerror.NewAuthorizationError(
			err,
			"user not found",
		)
	} else if err != nil {
		return model.Session{}, err
	}

	if err := verifyPassword(user.Password, password); err != nil {
		return model.Session{}, cerror.NewAuthorizationError(
			err,
			"password is not correct",
		)
	}

	oldSession, err := s.repo.Session().FindByUserID(ctx, user.ID)
	if err == nil {
		err = s.repo.Session().Delete(ctx, oldSession.ID)
		if err != nil {
			s.log.Error(err.Error())
		} else {
			s.log.Info(fmt.Sprintf("use already existed session(%v)", oldSession.ID))
		}
	} else if !errors.Is(err, cerror.ErrNotFound) {
		return model.Session{}, err
	}

	session := model.NewSession(user.ID, time.Now().AddDate(0, 1, 0))
	err = s.repo.Session().Create(ctx, session)
	if err != nil {
		return model.Session{}, err
	}

	return session, nil
}

func (s *Service) Signout(ctx context.Context, id string) error {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	err := s.signout(ctx, id)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Sign out session(%v)", id))
	}

	return err
}

func (s *Service) signout(ctx context.Context, sessionID string) error {
	err := s.repo.Session().Delete(ctx, sessionID)
	if errors.Is(err, cerror.ErrNotFound) {
		return cerror.NewAuthorizationError(
			err,
			fmt.Sprintf("invalid session id(%v)", sessionID),
		)
	} else if err != nil {
		return err
	}

	return nil
}

func (s *Service) Authorize(ctx context.Context, id string) (string, error) {
	reqID := ctx.Value(cctx.ReqIDKey).(string)
	s.log = s.log.Uniq(reqID)

	userID, err := s.authorize(ctx, id)

	if err != nil {
		s.log.Error(err.Error())
	} else {
		s.log.Info(fmt.Sprintf("Authorization user(%v)", userID))
	}

	return userID, err
}

func (s *Service) authorize(ctx context.Context, sessionID string) (string, error) {
	session, err := s.repo.Session().Find(ctx, sessionID)
	if errors.Is(err, cerror.ErrNotFound) {
		return "", cerror.NewAuthorizationError(
			err,
			fmt.Sprintf("invalid session id(%v)", sessionID),
		)
	} else if err != nil {
		return "", err
	}

	if time.Now().After(session.Expires) {
		return "", cerror.NewAuthorizationError(
			nil,
			fmt.Sprintf("session(%v) is already expired", session.ID),
		)
	}

	return session.UserID, nil
}
