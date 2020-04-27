package inmem

import (
	"sync"

	"github.com/x-color/calendar/model/auth"
	service "github.com/x-color/calendar/service/auth"
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
		users: []auth.User{},
	}
	s := sessionRepo{
		m:        sync.RWMutex{},
		sessions: []auth.Session{},
	}
	return inmem{
		userRepo:    u,
		sessionRepo: s,
	}
}
