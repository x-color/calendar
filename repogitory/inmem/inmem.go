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
	planRepo     planRepo
	calUserRepo  calUserRepo
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

func (m *inmem) Plan() cs.PlanRepogitory {
	return &m.planRepo
}

func (m *inmem) CalUser() cs.UserRepogitory {
	return &m.calUserRepo
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
	p := planRepo{
		m:     sync.RWMutex{},
		plans: []calendar.Plan{},
	}
	cu := calUserRepo{
		m:     sync.RWMutex{},
		users: []calendar.User{},
	}
	return inmem{
		userRepo:     u,
		sessionRepo:  s,
		calendarRepo: c,
		planRepo:     p,
		calUserRepo:  cu,
	}
}
