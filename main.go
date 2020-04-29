package main

import (
	"os"

	"github.com/x-color/calendar/app/rest"
	"github.com/x-color/calendar/logging"
	"github.com/x-color/calendar/repogitory/inmem"
	"github.com/x-color/calendar/service/auth"
	"github.com/x-color/calendar/service/calendar"
)

func main() {
	l := logging.NewLogger(os.Stdout)
	r := inmem.NewRepogitory()
	as := auth.NewService(&r, &l)
	cs := calendar.NewService(&r, &l)
	rest.StartServer(as, cs, &l)
}
