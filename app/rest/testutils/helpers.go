package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	ar "github.com/x-color/calendar/auth/repogitory/store"
	as "github.com/x-color/calendar/auth/service"
	cr "github.com/x-color/calendar/calendar/repogitory/store"
	cs "github.com/x-color/calendar/calendar/service"
	"github.com/x-color/calendar/logging"
)

var (
	pdb *sql.DB       = nil
	rdb *redis.Client = nil
)

func IgnoreKey(key string) cmp.Option {
	return cmpopts.IgnoreMapEntries(func(k string, t interface{}) bool {
		return k == key
	})
}

func connectDB() (*sql.DB, *redis.Client) {
	if pdb == nil || rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			"localhost", 5432, "testuser", "password", "calendar")

		var err error
		pdb, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		if err = pdb.Ping(); err != nil {
			panic(err)
		}
	}

	return pdb, rdb
}

func NewAuthRepo() as.Repogitory {
	pdb, rdb := connectDB()
	_, err := pdb.Exec("DELETE FROM auth.users")
	if err != nil {
		panic(err)
	}
	r := ar.NewRepogitory(pdb, rdb)
	return &r
}

func NewCalRepo() cs.Repogitory {
	db, _ := connectDB()

	_, err := pdb.Exec("DELETE FROM calendar.plan_shares")
	if err != nil {
		panic(err)
	}
	_, err = pdb.Exec("DELETE FROM calendar.plans")
	if err != nil {
		panic(err)
	}
	_, err = pdb.Exec("DELETE FROM calendar.calendar_shares")
	if err != nil {
		panic(err)
	}
	_, err = pdb.Exec("DELETE FROM calendar.calendars")
	if err != nil {
		panic(err)
	}
	_, err = pdb.Exec("DELETE FROM calendar.users")
	if err != nil {
		panic(err)
	}

	r := cr.NewRepogitory(db)
	return &r
}

func NewLogger() logging.Logger {
	l := logging.NewLogger(ioutil.Discard)
	return &l
}

func DummyCalService() cs.Service {
	return cs.Service{}
}

func MakeSession(authRepo as.Repogitory) (string, string) {
	userID := uuid.New().String()
	sessionID := uuid.New().String()
	authRepo.User().Create(context.Background(), as.UserData{
		ID:   userID,
		Name: "Alice",
	})
	authRepo.Session().Create(context.Background(), as.SessionData{
		ID:      sessionID,
		UserID:  userID,
		Expires: time.Now().Add(time.Hour).Unix(),
	})
	return userID, sessionID
}
