package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/x-color/calendar/calendar/service"
	cerror "github.com/x-color/calendar/model/error"
)

type userRepo struct {
	db *sql.DB
	tx *sql.Tx
}

func (r *userRepo) Find(ctx context.Context, id string) (service.UserData, error) {
	const query = "SELECT id FROM calendar.users WHERE id = $1"

	user := service.UserData{}
	var err error
	if r.tx != nil {
		err = r.tx.QueryRow(query, id).Scan(&user.ID)
	} else {
		err = r.db.QueryRow(query, id).Scan(&user.ID)
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return user, cerror.NewNotFoundError(
			err,
			fmt.Sprintf("not found a user(%v)", id),
		)
	case err != nil:
		return user, cerror.NewInternalError(
			err,
			"failed to scan query result",
		)
	}

	return user, nil
}

func (r *userRepo) Create(ctx context.Context, user service.UserData) error {
	const query = "INSERT INTO calendar.users (id) VALUES ($1)"

	var err error
	if r.tx != nil {
		_, err = r.tx.Exec(query, user.ID)
	} else {
		_, err = r.db.Exec(query, user.ID)
	}
	if err != nil {
		return cerror.NewInternalError(
			err,
			"failed to query",
		)
	}
	return nil
}
