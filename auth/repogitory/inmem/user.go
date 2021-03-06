package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/x-color/calendar/auth/service"
	cerror "github.com/x-color/calendar/model/error"
)

type userRepo struct {
	m     sync.RWMutex
	users []service.UserData
}

func (r *userRepo) FindByName(ctx context.Context, name string) (service.UserData, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	for _, u := range r.users {
		if name == u.Name {
			return u, nil
		}
	}

	return service.UserData{}, cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found name(%v)", name),
	)
}

func (r *userRepo) Create(ctx context.Context, user service.UserData) error {
	r.m.RLock()
	for _, u := range r.users {
		if user.ID == u.ID {
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
