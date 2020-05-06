package inmem

import (
	"sync"

	"github.com/x-color/calendar/calendar/service"
)

type inmem struct {
	calendarRepo calendarRepo
	planRepo     planRepo
	calUserRepo  calUserRepo
}

func (m *inmem) Calendar() service.CalendarRepogitory {
	return &m.calendarRepo
}

func (m *inmem) Plan() service.PlanRepogitory {
	return &m.planRepo
}

func (m *inmem) CalUser() service.UserRepogitory {
	return &m.calUserRepo
}

func NewRepogitory() inmem {
	c := calendarRepo{
		m:         sync.RWMutex{},
		calendars: []service.CalendarData{},
	}
	p := planRepo{
		m:     sync.RWMutex{},
		plans: []service.PlanData{},
	}
	cu := calUserRepo{
		m:     sync.RWMutex{},
		users: []service.UserData{},
	}
	return inmem{
		calendarRepo: c,
		planRepo:     p,
		calUserRepo:  cu,
	}
}
