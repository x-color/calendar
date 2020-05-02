package inmem

import (
	"sync"

	"github.com/x-color/calendar/model/auth"
	"github.com/x-color/calendar/model/calendar"
	as "github.com/x-color/calendar/service/auth"
	cs "github.com/x-color/calendar/service/calendar"
)

type inmem struct {
	userRepo     userRepo
	sessionRepo  sessionRepo
	calendarRepo calendarRepo
}

func (m *inmem) User() as.UserRepogitory {
	return &m.userRepo
}

func (m *inmem) Session() as.SessionRepogitory {
	return &m.sessionRepo
}

func (m *inmem) Calendar() cs.CalendarRepogitory {
	return &m.calendarRepo
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
	c := calendarRepo{
		m:         sync.RWMutex{},
		calendars: []calendar.Calendar{},
	}
	return inmem{
		userRepo:     u,
		sessionRepo:  s,
		calendarRepo: c,
	}
}
