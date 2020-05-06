package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	ase "github.com/x-color/calendar/app/rest/auth"
	cse "github.com/x-color/calendar/app/rest/calendar"
	"github.com/x-color/calendar/app/rest/middlewares"
	as "github.com/x-color/calendar/auth/service"
	cs "github.com/x-color/calendar/calendar/service"
	"github.com/x-color/calendar/logging"
	cctx "github.com/x-color/calendar/model/ctx"
)

type msgContent struct {
	Msg string `json:"message"`
}

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
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.LoggingMiddleware(l))
	r.Use(middlewares.ResponseHeaderMiddleware)

	ar := r.PathPrefix("/auth").Subrouter()
	ase.NewRouter(ar, ae)

	ur := r.PathPrefix("/register").Subrouter()
	ur.Use(middlewares.AuthorizationMiddleware(authService))
	ur.HandleFunc("", ue.RegisterHandler).Methods(http.MethodPost)

	cr := r.PathPrefix("/calendars").Subrouter()
	cr.Use(middlewares.AuthorizationMiddleware(authService))
	cr.Use(userCheckerMiddleware(calService))
	cr.HandleFunc("", ce.GetCalendarsHandler).Methods(http.MethodGet)
	cr.HandleFunc("", ce.MakeCalendarHandler).Methods(http.MethodPost)
	cr.HandleFunc("/{id}", ce.RemoveCalendarHandler).Methods(http.MethodDelete)
	cr.HandleFunc("/{id}", ce.ChangeCalendarHandler).Methods(http.MethodPatch)

	pr := r.PathPrefix("/plans").Subrouter()
	pr.Use(middlewares.AuthorizationMiddleware(authService))
	cr.Use(userCheckerMiddleware(calService))
	pr.HandleFunc("", pe.ScheduleHandler).Methods(http.MethodPost)
	pr.HandleFunc("/{id}", pe.UnsheduleHandler).Methods(http.MethodDelete)
	pr.HandleFunc("/{id}", pe.ResheduleHandler).Methods(http.MethodPatch)

	return r
}

func userCheckerMiddleware(service cs.Service) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(cctx.UserIDKey).(string)
			err := service.CheckRegistration(r.Context(), userID)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(msgContent{"forbidden"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
