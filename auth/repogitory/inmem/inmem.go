package inmem

import (
	"sync"

	"github.com/x-color/calendar/auth/model"
	"github.com/x-color/calendar/auth/service"
)

type inmem struct {
	userRepo    userRepo
	sessionRepo sessionRepo
}

func (m *inmem) User() service.UserRepogitory {
	return &m.userRepo
}

func (m *inmem) Session() service.SessionRepogitory {
	return &m.sessionRepo
}

func NewRepogitory() inmem {
	u := userRepo{
		m:     sync.RWMutex{},
		users: []model.User{},
	}
	s := sessionRepo{
		m:        sync.RWMutex{},
		sessions: []model.Session{},
	}
	return inmem{
		userRepo:    u,
		sessionRepo: s,
	}
}
