package main

import (
	"os"

	"github.com/x-color/calendar/app/rest"
	authInmem "github.com/x-color/calendar/auth/repogitory/inmem"
	as "github.com/x-color/calendar/auth/service"
	calInmem "github.com/x-color/calendar/calendar/repogitory/inmem"
	cs "github.com/x-color/calendar/calendar/service"
	"github.com/x-color/calendar/logging"
)

func main() {
	l := logging.NewLogger(os.Stdout)
	ar := authInmem.NewRepogitory()
	cr := calInmem.NewRepogitory()
	a := as.NewService(&ar, &l)
	c := cs.NewService(&cr, &l)
	rest.StartServer(a, c, &l)
}
