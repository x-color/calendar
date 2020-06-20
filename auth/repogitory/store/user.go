package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/x-color/calendar/auth/service"
	cerror "github.com/x-color/calendar/model/error"
)

type userRepo struct {
	db *sql.DB
}

func (r *userRepo) FindByName(ctx context.Context, name string) (service.UserData, error) {
	stmt, err := r.db.Prepare("SELECT id, name, password FROM auth.users WHERE name = $1")
	if err != nil {
		return service.UserData{}, cerror.NewInternalError(
			err,
			"failed to build prepare statement",
		)
	}
	defer stmt.Close()

	user := service.UserData{}

	err = stmt.QueryRow(name).Scan(&user.ID, &user.Name, &user.Password)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return user, cerror.NewNotFoundError(
			err,
			fmt.Sprintf("not found a user has the name(%v)", name),
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
	stmt, err := r.db.Prepare("INSERT INTO auth.users (id, name, password) VALUES ($1, $2, $3)")
	if err != nil {
		return cerror.NewInternalError(
			err,
			"failed to build prepare statement",
		)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.Name, user.Password)
	if err != nil {
		return cerror.NewInternalError(
			err,
			"failed to query",
		)
	}
	return nil
}
