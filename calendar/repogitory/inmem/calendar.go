package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/x-color/calendar/calendar/service"
	cerror "github.com/x-color/calendar/model/error"
	"github.com/x-color/slice/strs"
)

type calendarRepo struct {
	m         sync.RWMutex
	calendars []service.CalendarData
}

func (r *calendarRepo) Find(ctx context.Context, id string) (service.CalendarData, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	for _, c := range r.calendars {
		if id == c.ID {
			return c, nil
		}
	}
	return service.CalendarData{}, cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found calendar(%v)", id),
	)
}

func (r *calendarRepo) FindByUserID(ctx context.Context, userID string) ([]service.CalendarData, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	cals := []service.CalendarData{}
	for _, c := range r.calendars {
		if strs.Contains(c.Shares, userID) {
			cals = append(cals, c)
		}
	}

	return cals, nil
}

func (r *calendarRepo) Create(ctx context.Context, cal service.CalendarData) error {
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

func (r *calendarRepo) Delete(ctx context.Context, id string) error {
	r.m.Lock()
	defer r.m.Unlock()
	for i, c := range r.calendars {
		if id == c.ID {
			if i == len(r.calendars)-1 {
				r.calendars = r.calendars[:i]
			} else {
				r.calendars = append(r.calendars[:i], r.calendars[i+1:]...)
			}
			return nil
		}
	}
	return cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found calendar(%v)", id),
	)
}

func (r *calendarRepo) Update(ctx context.Context, cal service.CalendarData) error {
	r.m.Lock()
	defer r.m.Unlock()
	for i, c := range r.calendars {
		if cal.ID == c.ID {
			r.calendars[i] = cal
			return nil
		}
	}
	return cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found calendar(%v)", cal.ID),
	)
}
