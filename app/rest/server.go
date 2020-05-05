package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	ase "github.com/x-color/calendar/app/rest/auth"
	cse "github.com/x-color/calendar/app/rest/calendar"
	as "github.com/x-color/calendar/auth/service"
	cs "github.com/x-color/calendar/calendar/service"
	"github.com/x-color/calendar/logging"
)

func StartServer(authService as.Service, calService cs.Service, l logging.Logger) {
	r := newRouter(authService, calService, l)
	http.ListenAndServe(":8080", r)
}

func newRouter(authService as.Service, calService cs.Service, l logging.Logger) *mux.Router {
	ae := ase.NewAuthEndpoint(authService)
	ce := cse.NewCalEndpoint(calService)
	pe := cse.NewPlanEndpoint(calService)
	ue := cse.NewUserEndpoint(calService)

	r := mux.NewRouter()
	r.NotFoundHandler = http.NotFoundHandler()
	r.Use(reqIDMiddleware)
	r.Use(loggingMiddleware(l))
	r.Use(responseHeaderMiddleware)

	ar := r.PathPrefix("/auth").Subrouter()
	ar.HandleFunc("/signup", ae.SignupHandler).Methods(http.MethodPost)
	ar.HandleFunc("/signin", ae.SigninHandler)
	ar.HandleFunc("/signout", ae.SignoutHandler)

	ur := r.PathPrefix("/register").Subrouter()
	ur.Use(authorizationMiddleware(authService))
	ur.HandleFunc("", ue.RegisterHandler).Methods(http.MethodPost)

	cr := r.PathPrefix("/calendars").Subrouter()
	cr.Use(authorizationMiddleware(authService))
	cr.Use(userCheckerMiddleware(calService))
	cr.HandleFunc("", ce.GetCalendarsHandler).Methods(http.MethodGet)
	cr.HandleFunc("", ce.MakeCalendarHandler).Methods(http.MethodPost)
	cr.HandleFunc("/{id}", ce.RemoveCalendarHandler).Methods(http.MethodDelete)
	cr.HandleFunc("/{id}", ce.ChangeCalendarHandler).Methods(http.MethodPatch)

	pr := r.PathPrefix("/plans").Subrouter()
	pr.Use(authorizationMiddleware(authService))
	cr.Use(userCheckerMiddleware(calService))
	pr.HandleFunc("", pe.ScheduleHandler).Methods(http.MethodPost)
	pr.HandleFunc("/{id}", pe.UnsheduleHandler).Methods(http.MethodDelete)
	pr.HandleFunc("/{id}", pe.ResheduleHandler).Methods(http.MethodPatch)

	return r
}
