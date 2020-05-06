package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/x-color/calendar/auth/service"
	cerror "github.com/x-color/calendar/model/error"
)

type sessionRepo struct {
	m        sync.RWMutex
	sessions []service.SessionData
}

func (r *sessionRepo) Find(ctx context.Context, id string) (service.SessionData, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	for _, s := range r.sessions {
		if id == s.ID {
			return s, nil
		}
	}

	return service.SessionData{}, cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found session(%v)", id),
	)
}

func (r *sessionRepo) FindByUserID(ctx context.Context, userID string) (service.SessionData, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	for _, s := range r.sessions {
		if userID == s.UserID {
			return s, nil
		}
	}

	return service.SessionData{}, cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found user(%v) session", userID),
	)
}

func (r *sessionRepo) Create(ctx context.Context, session service.SessionData) error {
	r.m.RLock()
	for _, s := range r.sessions {
		if session.ID == s.ID {
			r.m.RUnlock()
			return cerror.NewDuplicationError(
				nil,
				fmt.Sprintf("same key(%v)", session.ID),
			)
		}
	}
	r.m.RUnlock()
	r.m.Lock()
	r.sessions = append(r.sessions, session)
	r.m.Unlock()
	return nil
}

func (r *sessionRepo) Delete(ctx context.Context, id string) error {
	r.m.Lock()
	defer r.m.Unlock()
	for i, s := range r.sessions {
		if id == s.ID {
			r.sessions = append(r.sessions[:i], r.sessions[i+1:]...)
			return nil
		}
	}
	return cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found session(%v)", id),
	)
}
