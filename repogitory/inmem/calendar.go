package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/x-color/calendar/model/calendar"
	cerror "github.com/x-color/calendar/model/error"
)

type calendarRepo struct {
	m         sync.RWMutex
	calendars []calendar.Calendar
}

func (r *calendarRepo) Create(ctx context.Context, cal calendar.Calendar) error {
	r.m.RLock()
	for _, c := range r.calendars {
		if c.ID == cal.ID {
			r.m.RUnlock()
			return cerror.NewDuplicationError(
				nil,
				fmt.Sprintf("same key(%v)", cal.ID),
			)
		}
	}
	r.m.RUnlock()
	r.m.Lock()
	r.calendars = append(r.calendars, cal)
	r.m.Unlock()
	return nil
}
