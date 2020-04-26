package main

import (
	"os"

	"github.com/x-color/calendar/app/rest"
	"github.com/x-color/calendar/logging"
	"github.com/x-color/calendar/repogitory/inmem"
	"github.com/x-color/calendar/service/auth"
)

func main() {
	l := logging.NewLogger(os.Stdout)
	r := inmem.NewRepogitory()
	s := auth.NewService(&r, &l)
	rest.StartServer(s, &l)
}
