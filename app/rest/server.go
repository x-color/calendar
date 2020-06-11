package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	ase "github.com/x-color/calendar/app/rest/auth"
	cse "github.com/x-color/calendar/app/rest/calendar"
	"github.com/x-color/calendar/app/rest/middlewares"
	as "github.com/x-color/calendar/auth/service"
	cs "github.com/x-color/calendar/calendar/service"
	"github.com/x-color/calendar/logging"
)

func StartServer(authService as.Service, calService cs.Service, l logging.Logger) {
	r := newRouter(authService, calService, l)
	http.ListenAndServe(":8080", r)
}

func newRouter(authService as.Service, calService cs.Service, l logging.Logger) *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = http.NotFoundHandler()
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.LoggingMiddleware(l))

	apiRouter := r.PathPrefix("/api").Subrouter()

	ar := apiRouter.PathPrefix("/auth").Subrouter()
	ase.NewRouter(ar, authService)

	ur := apiRouter.PathPrefix("/register").Subrouter()
	cse.NewUserRouter(ur, calService, authService)

	cr := apiRouter.PathPrefix("/calendars").Subrouter()
	cse.NewCalendarRouter(cr, calService, authService)

	pr := apiRouter.PathPrefix("/plans").Subrouter()
	cse.NewPlanRouter(pr, calService, authService)

	return r
}
