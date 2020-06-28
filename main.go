package main

import (
	"context"
	"database/sql"
	"log"
	"net/url"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/x-color/calendar/app/rest"
	authStore "github.com/x-color/calendar/auth/repogitory/store"
	as "github.com/x-color/calendar/auth/service"
	calStore "github.com/x-color/calendar/calendar/repogitory/store"
	cs "github.com/x-color/calendar/calendar/service"
	"github.com/x-color/calendar/logging"
)

func main() {
	redisURL, err := url.ParseRequestURI(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	pwd, ok := redisURL.User.Password()
	if !ok {
		log.Fatalln("REDIS_URL does not have password")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisURL.Host,
		Password: pwd,
	})
	defer rdb.Close()
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalln(err)
	}

	pdb, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	defer pdb.Close()
	if err = pdb.Ping(); err != nil {
		log.Fatalln(err)
	}

	if err := createTables(pdb); err != nil {
		log.Fatalln(err)
	}

	l := logging.NewLogger(os.Stdout)
	ar := authStore.NewRepogitory(pdb, rdb)
	cr := calStore.NewRepogitory(pdb)
	a := as.NewService(&ar, &l)
	c := cs.NewService(&cr, &l)
	rest.StartServer(a, c, &l, os.Getenv("PORT"))
}

func createTables(db *sql.DB) error {
	_, err := db.Exec("CREATE SCHEMA IF NOT EXISTS auth")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS calendar")
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS auth.users (
		id CHAR(36) PRIMARY KEY,
		name VARCHAR(64) NOT NULL,
		password VARCHAR(72) NOT NULL
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS calendar.users (
		id CHAR(36) PRIMARY KEY,
		FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS calendar.calendars (
		id CHAR(36) PRIMARY KEY,
		userid CHAR(36),
		name NAME NOT NULL,
		color VARCHAR(20) NOT NULL,
		FOREIGN KEY (userid) REFERENCES calendar.users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS calendar.calendar_shares (
		userid CHAR(36),
		calendarid CHAR(36),
		PRIMARY KEY(userid, calendarid),
		FOREIGN KEY (userid) REFERENCES calendar.users(id) ON DELETE CASCADE,
		FOREIGN KEY (calendarid) REFERENCES calendar.calendars(id) ON DELETE CASCADE
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS calendar.plans (
		id CHAR(36) PRIMARY KEY,
		userid CHAR(36),
		calendarid CHAR(36),
		name NAME NOT NULL,
		memo VARCHAR(400),
		color VARCHAR(20) NOT NULL,
		private BOOLEAN NOT NULL,
		isallday BOOLEAN NOT NULL,
		begintime BIGINT NOT NULL,
		endtime BIGINT NOT NULL,
		FOREIGN KEY (userid) REFERENCES calendar.users(id) ON DELETE CASCADE,
		FOREIGN KEY (calendarid) REFERENCES calendar.calendars(id) ON DELETE CASCADE
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS calendar.plan_shares (
		calendarid CHAR(36),
		planid CHAR(36),
		PRIMARY KEY(calendarid, planid),
		FOREIGN KEY (calendarid) REFERENCES calendar.calendars(id) ON DELETE CASCADE,
		FOREIGN KEY (planid) REFERENCES calendar.plans(id) ON DELETE CASCADE
	)`)
	return err
}
