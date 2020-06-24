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

type planRepo struct {
	db *sql.DB
	tx *sql.Tx
}

func (r *planRepo) Find(ctx context.Context, id string) (service.PlanData, error) {
	const query = `
		SELECT plans.id, plans.userid, plans.calendarid, plans.name, plans.memo,
			   plans.color, plans.private, plans.isallday, plans.begintime, plans.endtime, shares.calendarid
		FROM calendar.plans plans
		INNER JOIN calendar.plan_shares shares
		ON plans.id = shares.planid 
		WHERE plans.id = $1
	`

	var rows *sql.Rows
	var err error
	if r.tx != nil {
		rows, err = r.tx.Query(query, id)
	} else {
		rows, err = r.db.Query(query, id)
	}
	defer rows.Close()

	var calID string
	plan := service.PlanData{Shares: []string{}}
	for rows.Next() {
		err := rows.Scan(&plan.ID, &plan.UserID, &plan.CalendarID, &plan.Name, &plan.Memo,
			&plan.Color, &plan.Private, &plan.IsAllDay, &plan.Begin, &plan.End, &calID)
		if err != nil {
			return plan, cerror.NewInternalError(
				err,
				"failed to scan query result",
			)
		}
		plan.Shares = append(plan.Shares, calID)
	}

	err = rows.Err()
	switch {
	case errors.Is(err, sql.ErrNoRows) || len(plan.ID) == 0:
		return plan, cerror.NewNotFoundError(
			err,
			fmt.Sprintf("not found plan(%v)", id),
		)
	case err != nil:
		return plan, cerror.NewInternalError(
			err,
			"failed to scan query result",
		)
	}

	return plan, nil
}

func (r *planRepo) FindByCalendarID(ctx context.Context, calID string) ([]service.PlanData, error) {
	const query = `
		SELECT plans.id, plans.userid, plans.calendarid, plans.name, plans.memo,
			   plans.color, plans.private, plans.isallday, plans.begintime,  plans.endtime, shares.calendarid
		FROM calendar.plans plans
		JOIN calendar.plan_shares shares
		ON plans.id = shares.planid
		WHERE plans.id IN (
			SELECT planid
			FROM calendar.plan_shares
			WHERE calendarid = $1
		)
		ORDER BY plans.id
	`

	var rows *sql.Rows
	var err error
	if r.tx != nil {
		rows, err = r.tx.Query(query, calID)
	} else {
		rows, err = r.db.Query(query, calID)
	}
	defer rows.Close()

	plans := []service.PlanData{}

	var id string
	var plan, newPlan service.PlanData
	for rows.Next() {
		err := rows.Scan(&newPlan.ID, &newPlan.UserID, &newPlan.CalendarID, &newPlan.Name, &newPlan.Memo,
			&newPlan.Color, &newPlan.Private, &newPlan.IsAllDay, &newPlan.Begin, &newPlan.End, &id)
		if err != nil {
			return nil, cerror.NewInternalError(
				err,
				"failed to scan query result",
			)
		}

		if newPlan.ID == plan.ID {
			plan.Shares = append(plan.Shares, id)
		} else {
			plans = append(plans, plan)
			plan = service.PlanData{
				ID:         newPlan.ID,
				UserID:     newPlan.UserID,
				CalendarID: newPlan.CalendarID,
				Name:       newPlan.Name,
				Memo:       newPlan.Memo,
				Color:      newPlan.Color,
				Private:    newPlan.Private,
				IsAllDay:   newPlan.IsAllDay,
				Begin:      newPlan.Begin,
				End:        newPlan.End,
				Shares:     []string{id},
			}
		}
	}
	plans = append(plans, plan)

	err = rows.Err()
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, cerror.NewNotFoundError(
			err,
			fmt.Sprintf("not found plans of calendar(%v)", calID),
		)
	case err != nil:
		return nil, cerror.NewInternalError(
			err,
			"failed to scan query result",
		)
	}

	// Remove the first data of plans. It is empty.
	return plans[1:], nil
}

func (r *planRepo) Create(ctx context.Context, plan service.PlanData) error {
	var err error
	if r.tx == nil {
		err = r.transaction(func() error {
			return r.create(ctx, plan)
		})
	} else {
		err = r.create(ctx, plan)
	}

	if err != nil {
		return cerror.NewInternalError(
			err,
			"failed to create plan",
		)
	}
	return nil
}

func (r *planRepo) create(ctx context.Context, plan service.PlanData) error {
	const insPlanQuery = `
		INSERT INTO calendar.plans (id, userid, calendarid, name, memo, color, private, isallday, begintime, endtime)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.tx.Exec(insPlanQuery, plan.ID, plan.UserID, plan.CalendarID, plan.Name, plan.Memo,
		plan.Color, plan.Private, plan.IsAllDay, plan.Begin, plan.End)
	if err != nil {
		return err
	}

	const insSharesQuery = "INSERT INTO calendar.plan_shares (calendarid, planid) VALUES ($1, $2)"
	stmt, err := r.tx.Prepare(insSharesQuery)
	if err != nil {
		return err
	}
	for _, calID := range plan.Shares {
		_, err := stmt.Exec(calID, plan.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *planRepo) Delete(ctx context.Context, id string) error {
	var err error
	if r.tx == nil {
		err = r.transaction(func() error {
			return r.delete(ctx, id)
		})
	} else {
		err = r.delete(ctx, id)
	}

	if err != nil {
		return cerror.NewInternalError(
			err,
			"failed to delete plan",
		)
	}
	return nil
}

func (r *planRepo) delete(ctx context.Context, id string) error {
	const query = "DELETE FROM calendar.plans WHERE id = $1"
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
			fmt.Sprintf("not found plan(%v)", id),
		)
	}

	return nil
}

func (r *planRepo) Update(ctx context.Context, plan service.PlanData) error {
	var err error
	if r.tx == nil {
		err = r.transaction(func() error {
			return r.update(ctx, plan)
		})
	} else {
		err = r.update(ctx, plan)
	}

	if err != nil {
		return cerror.NewInternalError(
			err,
			"failed to update plan",
		)
	}
	return nil
}

func (r *planRepo) update(ctx context.Context, plan service.PlanData) error {
	const query = "SELECT calendarid FROM calendar.plan_shares WHERE planid = $1"

	rows, err := r.tx.Query(query, plan.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	calIDs := []string{}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return cerror.NewInternalError(
				err,
				"failed to scan query result",
			)
		}
		calIDs = append(calIDs, id)
	}

	err = rows.Err()
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return cerror.NewNotFoundError(
			err,
			fmt.Sprintf("not found plan(%v)", plan.ID),
		)
	case err != nil:
		return cerror.NewInternalError(
			err,
			"failed to scan query result",
		)
	}

	if delCalIDs := strs.Sub(calIDs, plan.Shares); len(delCalIDs) > 0 {
		l := make([]string, len(delCalIDs))
		for i := 0; i < len(delCalIDs); i++ {
			l[i] = fmt.Sprintf("$%v", i+2)
		}
		delSharesQuery := `
			DELETE FROM calendar.plan_shares
			WHERE planid = $1 AND calendarid = ANY($2)
		`
		_, err = r.tx.Exec(delSharesQuery, plan.ID, pq.Array(delCalIDs))
		if err != nil {
			return err
		}
	}

	addCalIDs := strs.Sub(plan.Shares, calIDs)
	addSharesQuery := "INSERT INTO calendar.plan_shares (planid, calendarid) VALUES ($1, $2)"
	for _, id := range addCalIDs {
		_, err := r.tx.Exec(addSharesQuery, plan.ID, id)
		if err != nil {
			return err
		}
	}

	const updateCalQuery = `
		UPDATE calendar.plans
		SET name = $1, memo = $2, color = $3, private = $4, isallday = $5, begintime = $6, endtime = $7
		WHERE id = $8
	`
	_, err = r.tx.Exec(updateCalQuery, plan.Name, plan.Memo, plan.Color, plan.Private,
		plan.IsAllDay, plan.Begin, plan.End, plan.ID)
	return err
}

func (r *planRepo) transaction(f func() error) error {
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

func (r *planRepo) beginTx() error {
	tx, err := r.db.Begin()
	if err == nil {
		r.tx = tx
	}
	return err
}

func (r *planRepo) commitTx() error {
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

func (r *planRepo) rollbackTx() error {
	err := r.tx.Rollback()
	r.tx = nil
	return err
}
