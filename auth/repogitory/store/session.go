package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/x-color/calendar/auth/service"
	cerror "github.com/x-color/calendar/model/error"
)

type sessionRepo struct {
	rdb *redis.Client
}

func (r *sessionRepo) Find(ctx context.Context, id string) (service.SessionData, error) {
	userID, err := r.rdb.Get(ctx, id).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return service.SessionData{}, cerror.NewNotFoundError(
			err,
			fmt.Sprintf("not found a session(%v)", id),
		)
	case err != nil:
		return service.SessionData{}, cerror.NewInternalError(
			err,
			"failed to get session",
		)
	}

	return service.SessionData{id, userID, time.Now().Add(time.Hour).Unix()}, nil
}

func (r *sessionRepo) Create(ctx context.Context, session service.SessionData) error {
	duration := time.Unix(session.Expires, 0).Sub(time.Now())
	set, err := r.rdb.SetNX(ctx, session.ID, session.UserID, duration).Result()
	switch {
	case err != nil:
		return cerror.NewInternalError(
			err,
			"failed to check same session already exists",
		)
	case !set:
		return cerror.NewDuplicationError(
			nil,
			fmt.Sprintf("same key(%v)", session.ID),
		)
	}

	return nil
}

func (r *sessionRepo) Delete(ctx context.Context, id string) error {
	n, err := r.rdb.Del(ctx, id).Result()
	switch {
	case n == 0:
		return cerror.NewNotFoundError(
			err,
			fmt.Sprintf("not found session(%v)", id),
		)
	case err != nil:
		return cerror.NewInternalError(
			err,
			"failed to delete session",
		)
	}
	return nil
}
