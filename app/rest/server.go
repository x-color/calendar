package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-color/calendar/service/auth"
)

type userContent struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type msgContent struct {
	Msg string `json:"message"`
}

func StartServer(s auth.Service, l auth.Logger) {
	r := newRouter(s, l)
	http.ListenAndServe(":8080", r)
}

func newRouter(as auth.Service, l auth.Logger) *mux.Router {
	ae := AuthEndpoint{as}
	ce := CalEndpoint{
		// TODO: calendar.Service
	}

	r := mux.NewRouter()
	r.NotFoundHandler = http.NotFoundHandler()
	r.Use(reqIDMiddleware)
	r.Use(loggingMiddleware(l))
	r.Use(responseHeaderMiddleware)

	ar := r.PathPrefix("/auth").Subrouter()
	ar.HandleFunc("/signup", ae.signupHandler).Methods(http.MethodPost)
	ar.HandleFunc("/signin", ae.signinHandler)
	ar.HandleFunc("/signout", ae.signoutHandler)

	cr := r.PathPrefix("/calendars").Subrouter()
	cr.Use(authorizationMiddleware(as))
	cr.HandleFunc("", ce.getCalendarsHandler).Methods(http.MethodGet)
	cr.HandleFunc("", ce.makeCalendarHandler).Methods(http.MethodPost)
	cr.HandleFunc("/{id}", ce.removeCalendarHandler).Methods(http.MethodDelete)
	cr.HandleFunc("/{id}", ce.changeCalendarHandler).Methods(http.MethodPost)

	return r
}
