package store

import (
	"database/sql"

	"github.com/x-color/calendar/calendar/service"
	cerror "github.com/x-color/calendar/model/error"
)

type store struct {
	db           *sql.DB
	tx           *sql.Tx
	calendarRepo calendarRepo
	planRepo     planRepo
	userRepo     userRepo
}

func (m *store) Calendar() service.CalendarRepogitory {
	m.calendarRepo.tx = m.tx
	return &m.calendarRepo
}

func (m *store) Plan() service.PlanRepogitory {
	m.planRepo.tx = m.tx
	return &m.planRepo
}

func (m *store) User() service.UserRepogitory {
	m.userRepo.tx = m.tx
	return &m.userRepo
}

func (m *store) BeginTX() error {
	tx, err := m.db.Begin()
	if err != nil {
		return cerror.NewInternalError(
			err,
			"failed to begin transaction",
		)
	}
	m.tx = tx
	return nil
}

func (m *store) Commit() error {
	if m.tx == nil {
		return nil
	}
	err := m.tx.Commit()
	if err != nil {
		err = m.Rollback()
		return cerror.NewInternalError(
			err,
			"failed to commit transaction",
		)
	}
	m.tx = nil
	return nil
}

func (m *store) Rollback() error {
	if m.tx == nil {
		return nil
	}
	err := m.tx.Rollback()
	m.tx = nil
	if err != nil {
		cerror.NewInternalError(
			err,
			"failed to rollback transaction",
		)
	}
	return nil
}

func NewRepogitory(db *sql.DB) store {
	c := calendarRepo{
		db: db,
	}
	p := planRepo{
		db: db,
	}
	u := userRepo{
		db: db,
	}
	return store{
		calendarRepo: c,
		planRepo:     p,
		userRepo:     u,
	}
}
