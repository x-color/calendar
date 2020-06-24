package main

import (
	"context"
	"database/sql"
	"fmt"
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
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "testuser", "password", "calendar")
	pdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer pdb.Close()
	if err = pdb.Ping(); err != nil {
		panic(err)
	}

	l := logging.NewLogger(os.Stdout)
	ar := authStore.NewRepogitory(pdb, rdb)
	cr := calStore.NewRepogitory(pdb)
	a := as.NewService(&ar, &l)
	c := cs.NewService(&cr, &l)
	rest.StartServer(a, c, &l)
}
