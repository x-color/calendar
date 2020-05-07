package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/x-color/calendar/calendar/service"
	cerror "github.com/x-color/calendar/model/error"
)

type userRepo struct {
	m     sync.RWMutex
	users []service.UserData
}

func (r *userRepo) Find(ctx context.Context, id string) (service.UserData, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	for _, u := range r.users {
		if id == u.ID {
			return u, nil
		}
	}
	return service.UserData{}, cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found user(%v)", id),
	)
}

func (r *userRepo) Create(ctx context.Context, user service.UserData) error {
	r.m.RLock()
	for _, c := range r.users {
		if c.ID == user.ID {
			r.m.RUnlock()
			return cerror.NewDuplicationError(
				nil,
				fmt.Sprintf("same key(%v)", user.ID),
			)
		}
	}
	r.m.RUnlock()
	r.m.Lock()
	r.users = append(r.users, user)
	r.m.Unlock()
	return nil
}
