package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/x-color/calendar/calendar/service"
	cerror "github.com/x-color/calendar/model/error"
	"github.com/x-color/slice/strs"
)

type calendarRepo struct {
	db *sql.DB
	tx *sql.Tx
}

func (r *calendarRepo) Find(ctx context.Context, id string) (service.CalendarData, error) {
	const query = `
		SELECT cals.id, cals.userid, cals.name, cals.color, shares.userid
		FROM calendar.calendars cals
		INNER JOIN calendar.calendar_shares shares
		ON cals.id = shares.calendarid 
		WHERE cals.id = $1
	`

	var rows *sql.Rows
	var err error
	if r.tx != nil {
		rows, err = r.tx.Query(query, id)
	} else {
		rows, err = r.db.Query(query, id)
	}
	defer rows.Close()

	var userID string
	calendar := service.CalendarData{Shares: []string{}}
	for rows.Next() {
		err := rows.Scan(&calendar.ID, &calendar.UserID, &calendar.Name, &calendar.Color, &userID)
		if err != nil {
			return calendar, cerror.NewInternalError(
				err,
				"failed to scan query result",
			)
		}
		calendar.Shares = append(calendar.Shares, userID)
	}

	err = rows.Err()
	switch {
	case errors.Is(err, sql.ErrNoRows) || len(calendar.ID) == 0:
		return calendar, cerror.NewNotFoundError(
			err,
			fmt.Sprintf("not found calendar(%v)", id),
		)
	case err != nil:
		return calendar, cerror.NewInternalError(
			err,
			"failed to scan query result",
		)
	}

	return calendar, nil
}

func (r *calendarRepo) FindByUserID(ctx context.Context, userID string) ([]service.CalendarData, error) {
	const query = `
		SELECT cals.id, cals.userid, cals.name, cals.color, shares.userid
		FROM calendar.calendars cals
		JOIN calendar.calendar_shares shares
		ON cals.id = shares.calendarid
		WHERE cals.id IN (
			SELECT calendarid
			FROM calendar.calendar_shares
			WHERE userid = $1
		)
		ORDER BY cals.id
	`

	var rows *sql.Rows
	var err error
	if r.tx != nil {
		rows, err = r.tx.Query(query, userID)
	} else {
		rows, err = r.db.Query(query, userID)
	}
	defer rows.Close()

	calendars := []service.CalendarData{}

	var id string
	var cal, newCal service.CalendarData
	for rows.Next() {
		err := rows.Scan(&newCal.ID, &newCal.UserID, &newCal.Name, &newCal.Color, &id)
		if err != nil {
			return nil, cerror.NewInternalError(
				err,
				"failed to scan query result",
			)
		}

		if newCal.ID == cal.ID {
			cal.Shares = append(cal.Shares, id)
		} else {
			calendars = append(calendars, cal)
			cal = service.CalendarData{
				ID:     newCal.ID,
				UserID: newCal.UserID,
				Name:   newCal.Name,
				Color:  newCal.Color,
				Shares: []string{id},
			}
		}
	}
	calendars = append(calendars, cal)

	err = rows.Err()
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, cerror.NewNotFoundError(
			err,
			fmt.Sprintf("not found calendars for user(%v)", id),
		)
	case err != nil:
		return nil, cerror.NewInternalError(
			err,
			"failed to scan query result",
		)
	}

	// Remove the first data of calnedars. It is empty.
	return calendars[1:], nil
}

func (r *calendarRepo) Create(ctx context.Context, cal service.CalendarData) error {
	var err error
	if r.tx == nil {
		err = r.transaction(func() error {
			return r.create(ctx, cal)
		})
	} else {
		err = r.create(ctx, cal)
	}

	if err != nil {
		return cerror.NewInternalError(
			err,
			"failed to create calendar",
		)
	}
	return nil
}

func (r *calendarRepo) create(ctx context.Context, cal service.CalendarData) error {
	const insCalQuery = "INSERT INTO calendar.calendars (id, userid, name, color) VALUES ($1, $2, $3, $4)"
	_, err := r.tx.Exec(insCalQuery, cal.ID, cal.UserID, cal.Name, cal.Color)
	if err != nil {
		return err
	}

	const insSharesQuery = "INSERT INTO calendar.calendar_shares (userid, calendarid) VALUES ($1, $2)"
	stmt, err := r.tx.Prepare(insSharesQuery)
	if err != nil {
		return err
	}
	for _, userID := range cal.Shares {
		_, err := stmt.Exec(userID, cal.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *calendarRepo) Delete(ctx context.Context, id string) error {
	var err error
	if r.tx == nil {
		err = r.transaction(func() error {
			return r.delete(ctx, id)
		})
	} else {
		err = r.delete(ctx, id)
	}
	switch {
	case errors.Is(err, cerror.ErrNotFound):
		return err
	case err != nil:
		return cerror.NewInternalError(
			err,
			"failed to delete calendar",
		)
	}
	return nil
}

func (r *calendarRepo) delete(ctx context.Context, id string) error {
	const query = "DELETE FROM calendar.calendars WHERE id = $1"
	res, err := r.tx.Exec(query, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return cerror.NewNotFoundError(
			nil,
			fmt.Sprintf("not found calendar(%v)", id),
		)
	}

	return nil
}

func (r *calendarRepo) Update(ctx context.Context, cal service.CalendarData) error {
	var err error
	if r.tx == nil {
		err = r.transaction(func() error {
			return r.update(ctx, cal)
		})
	} else {
		err = r.update(ctx, cal)
	}

	switch {
	case errors.Is(err, cerror.ErrNotFound):
		return err
	case err != nil:
		return cerror.NewInternalError(
			err,
			"failed to update calendar",
		)
	}

	return nil
}

func (r *calendarRepo) update(ctx context.Context, cal service.CalendarData) error {
	const query = "SELECT userid FROM calendar.calendar_shares WHERE calendarid = $1"

	rows, err := r.tx.Query(query, cal.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	userIDs := []string{}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return cerror.NewInternalError(
				err,
				"failed to scan query result",
			)
		}
		userIDs = append(userIDs, id)
	}

	err = rows.Err()
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return cerror.NewNotFoundError(
			err,
			fmt.Sprintf("not found calendar(%v)", cal.ID),
		)
	case err != nil:
		return cerror.NewInternalError(
			err,
			"failed to scan query result",
		)
	}

	if delUserIDs := strs.Sub(userIDs, cal.Shares); len(delUserIDs) > 0 {
		l := make([]string, len(delUserIDs))
		for i := 0; i < len(delUserIDs); i++ {
			l[i] = fmt.Sprintf("$%v", i+2)
		}
		delSharesQuery := `
			DELETE FROM calendar.calendar_shares
			WHERE calendarid = $1 AND userid = ANY($2)
		`
		_, err = r.tx.Exec(delSharesQuery, cal.ID, pq.Array(delUserIDs))
		if err != nil {
			return err
		}
	}

	addUserIDs := strs.Sub(cal.Shares, userIDs)
	addSharesQuery := "INSERT INTO calendar.calendar_shares (calendarid, userid) VALUES ($1, $2)"
	for _, id := range addUserIDs {
		_, err := r.tx.Exec(addSharesQuery, cal.ID, id)
		if err != nil {
			return err
		}
	}

	const updateCalQuery = `
		UPDATE calendar.calendars
		SET name = $1, color = $2
		WHERE id = $3
	`
	_, err = r.tx.Exec(updateCalQuery, cal.Name, cal.Color, cal.ID)
	return err
}

func (r *calendarRepo) transaction(f func() error) error {
	err := r.beginTx()
	if err != nil {
		return err
	}

	err = f()
	if err != nil {
		if rerr := r.rollbackTx(); rerr != nil {
			return rerr
		}
		return err
	}

	return r.commitTx()
}

func (r *calendarRepo) beginTx() error {
	tx, err := r.db.Begin()
	if err == nil {
		r.tx = tx
	}
	return err
}

func (r *calendarRepo) commitTx() error {
	err := r.tx.Commit()
	if err != nil {
		if rerr := r.tx.Rollback(); rerr != nil {
			err = cerror.NewInternalError(
				rerr,
				"failed to commit and rollback",
			)
		}
	}
	r.tx = nil
	return err
}

func (r *calendarRepo) rollbackTx() error {
	err := r.tx.Rollback()
	r.tx = nil
	return err
}
