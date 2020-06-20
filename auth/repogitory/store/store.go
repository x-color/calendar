package store

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/x-color/calendar/auth/service"
)

type rds struct {
	userRepo    userRepo
	sessionRepo sessionRepo
}

func (m *rds) User() service.UserRepogitory {
	return &m.userRepo
}

func (m *rds) Session() service.SessionRepogitory {
	return &m.sessionRepo
}

func NewRepogitory(pdb *sql.DB, rdb *redis.Client) rds {
	u := userRepo{
		db: pdb,
	}
	s := sessionRepo{
		rdb: rdb,
	}
	return rds{
		userRepo:    u,
		sessionRepo: s,
	}
}
