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

	l := logging.NewLogger(os.Stdout)
	ar := authStore.NewRepogitory(pdb, rdb)
	cr := calStore.NewRepogitory(pdb)
	a := as.NewService(&ar, &l)
	c := cs.NewService(&cr, &l)
	rest.StartServer(a, c, &l, os.Getenv("PORT"))
}
