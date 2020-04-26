package auth

import (
	"context"
	"errors"
	"fmt"
	"unicode"

	"github.com/google/uuid"
	"github.com/x-color/calendar/model/auth"
	cctx "github.com/x-color/calendar/model/ctx"
	cerror "github.com/x-color/calendar/model/error"
	"golang.org/x/crypto/bcrypt"
)

type Repogitory interface {
	User() UserRepogitory
}

type UserRepogitory interface {
	FindByName(ctx context.Context, name string) (auth.User, error)
	Create(ctx context.Context, user auth.User) error
}

type Logger interface {
	Info(id, msg string)
	Error(id, msg string)
}

type Service struct {
	repo Repogitory
	log  Logger
}

func NewService(repo Repogitory, log Logger) Service {
	return Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) Signup(ctx context.Context, name, password string) (auth.User, error) {
	user, err := s.signup(ctx, name, password)

	reqID := ctx.Value(cctx.ReqIDKey).(string)
	if err != nil {
		s.log.Error(reqID, err.Error())
	} else {
		s.log.Info(reqID, fmt.Sprintf("Sign up user(%v)", reqID))
	}

	return user, err
}

func (s *Service) signup(ctx context.Context, name, password string) (auth.User, error) {
	if err := validateSigninInfo(name, password); err != nil {
		return auth.User{}, err
	}

	_, err := s.repo.User().FindByName(ctx, name)
	if err == nil {
		return auth.User{}, cerror.NewDuplicationError(
			nil,
			fmt.Sprintf("user(%v) already exists", name),
		)
	}
	if !errors.Is(err, cerror.ErrNotFound) {
		return auth.User{}, err
	}

	hash, err := passwordHash(password)
	if err != nil {
		return auth.User{}, cerror.NewInternalError(
			err,
			"failed to hash password",
		)
	}

	user := auth.NewUser(name, hash)
	err = s.repo.User().Create(ctx, user)
	if err != nil {
		return auth.User{}, err
	}
	return user, nil
}

func (s *Service) Signin(ctx context.Context, name, password string) (string, error) {
	if err := validateSigninInfo(name, password); err != nil {
		return "", err
	}

	user, err := s.repo.User().FindByName(ctx, name)
	if errors.Is(err, cerror.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	if err := verifyPassword(user.Password, password); err != nil {
		return "", nil
	}

	session := uuid.New().String()
	// TODO: Store session

	return session, nil
}

func validateSigninInfo(name, password string) error {
	if name == "" {
		return cerror.NewInvalidContentError(
			nil,
			"name is empty",
		)
	}
	if !isValidPassword(password) {
		return cerror.NewInvalidContentError(
			nil,
			"invalid password",
		)
	}
	return nil
}

func isValidPassword(password string) bool {
	hasMinLen := false
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false
	if 7 < len(password) && len(password) < 73 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func verifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
