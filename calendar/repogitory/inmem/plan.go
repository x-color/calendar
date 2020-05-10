package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/x-color/calendar/calendar/service"
	cerror "github.com/x-color/calendar/model/error"
	"github.com/x-color/slice/strs"
)

type planRepo struct {
	m     sync.RWMutex
	plans []service.PlanData
}

func (r *planRepo) Find(ctx context.Context, id string) (service.PlanData, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	for _, p := range r.plans {
		if id == p.ID {
			return p, nil
		}
	}
	return service.PlanData{}, cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found plan(%v)", id),
	)
}

func (r *planRepo) FindByCalendarID(ctx context.Context, calID string) ([]service.PlanData, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	plans := []service.PlanData{}
	for _, p := range r.plans {
		if strs.Contains(p.Shares, calID) {
			plans = append(plans, p)
		}
	}

	return plans, nil
}

func (r *planRepo) Create(ctx context.Context, plan service.PlanData) error {
	r.m.RLock()
	for _, c := range r.plans {
		if c.ID == plan.ID {
			r.m.RUnlock()
			return cerror.NewDuplicationError(
				nil,
				fmt.Sprintf("same key(%v)", plan.ID),
			)
		}
	}
	r.m.RUnlock()
	r.m.Lock()
	r.plans = append(r.plans, plan)
	r.m.Unlock()
	return nil
}

func (r *planRepo) Delete(ctx context.Context, id string) error {
	r.m.Lock()
	defer r.m.Unlock()
	for i, p := range r.plans {
		if id == p.ID {
			if i == len(r.plans)-1 {
				r.plans = r.plans[:i]
			} else {
				r.plans = append(r.plans[:i], r.plans[i+1:]...)
			}
			return nil
		}
	}
	return cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found plans(%v)", id),
	)
}

func (r *planRepo) Update(ctx context.Context, plan service.PlanData) error {
	r.m.Lock()
	defer r.m.Unlock()
	for i, c := range r.plans {
		if plan.ID == c.ID {
			r.plans[i] = plan
			return nil
		}
	}
	return cerror.NewNotFoundError(
		nil,
		fmt.Sprintf("not found plan(%v)", plan.ID),
	)
}
