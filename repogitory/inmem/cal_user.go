package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/x-color/calendar/model/calendar"
	cerror "github.com/x-color/calendar/model/error"
)

type calUserRepo struct {
	m     sync.RWMutex
	users []calendar.User
}

func (r *calUserRepo) Find(ctx context.Context, id string) (calendar.User, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	for _, u := range r.users {
		if id == u.ID {
			return u, nil
		}
	}
	return calendar.User{}, cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found user(%v)", id),
	)
}

func (r *calUserRepo) Create(ctx context.Context, user calendar.User) error {
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
