package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-color/calendar/service"
	"github.com/x-color/calendar/service/auth"
	"github.com/x-color/calendar/service/calendar"
)

type userContent struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type msgContent struct {
	Msg string `json:"message"`
}

func StartServer(as auth.Service, cs calendar.Service, l service.Logger) {
	r := newRouter(as, cs, l)
	http.ListenAndServe(":8080", r)
}

func newRouter(as auth.Service, cs calendar.Service, l service.Logger) *mux.Router {
	ae := AuthEndpoint{as}
	ce := CalEndpoint{cs}
	pe := PlanEndpoint{cs}
	ue := UserEndpoint{cs}

	r := mux.NewRouter()
	r.NotFoundHandler = http.NotFoundHandler()
	r.Use(reqIDMiddleware)
	r.Use(loggingMiddleware(l))
	r.Use(responseHeaderMiddleware)

	ar := r.PathPrefix("/auth").Subrouter()
	ar.HandleFunc("/signup", ae.signupHandler).Methods(http.MethodPost)
	ar.HandleFunc("/signin", ae.signinHandler)
	ar.HandleFunc("/signout", ae.signoutHandler)

	ur := r.PathPrefix("/register").Subrouter()
	ur.Use(authorizationMiddleware(as))
	ur.HandleFunc("", ue.registerHandler).Methods(http.MethodPost)

	cr := r.PathPrefix("/calendars").Subrouter()
	cr.Use(authorizationMiddleware(as))
	cr.HandleFunc("", ce.getCalendarsHandler).Methods(http.MethodGet)
	cr.HandleFunc("", ce.makeCalendarHandler).Methods(http.MethodPost)
	cr.HandleFunc("/{id}", ce.removeCalendarHandler).Methods(http.MethodDelete)
	cr.HandleFunc("/{id}", ce.changeCalendarHandler).Methods(http.MethodPatch)

	pr := r.PathPrefix("/plans").Subrouter()
	pr.Use(authorizationMiddleware(as))
	pr.HandleFunc("", pe.scheduleHandler).Methods(http.MethodPost)
	pr.HandleFunc("/{id}", pe.unsheduleHandler).Methods(http.MethodDelete)
	pr.HandleFunc("/{id}", pe.resheduleHandler).Methods(http.MethodPatch)

	return r
}
