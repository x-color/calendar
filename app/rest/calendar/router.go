package calendar

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-color/calendar/app/rest/middlewares"
	as "github.com/x-color/calendar/auth/service"
	cs "github.com/x-color/calendar/calendar/service"
	cctx "github.com/x-color/calendar/model/ctx"
)

func NewCalendarRouter(r *mux.Router, calService cs.Service, authService as.Service) {
	e := calEndpoint{calService}
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.ResponseHeaderMiddleware)
	r.Use(middlewares.AuthorizationMiddleware(authService))
	r.Use(userCheckerMiddleware(calService))
	r.HandleFunc("", e.GetCalendarsHandler).Methods(http.MethodGet)
	r.HandleFunc("", e.MakeCalendarHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}", e.RemoveCalendarHandler).Methods(http.MethodDelete)
	r.HandleFunc("/{id}", e.ChangeCalendarHandler).Methods(http.MethodPatch)
}

func NewPlanRouter(r *mux.Router, calService cs.Service, authService as.Service) {
	e := planEndpoint{calService}
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.ResponseHeaderMiddleware)
	r.Use(middlewares.AuthorizationMiddleware(authService))
	r.Use(userCheckerMiddleware(calService))
	r.HandleFunc("", e.ScheduleHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}", e.UnsheduleHandler).Methods(http.MethodDelete)
	r.HandleFunc("/{id}", e.ResheduleHandler).Methods(http.MethodPatch)
}

func NewUserRouter(r *mux.Router, calService cs.Service, authService as.Service) {
	e := userEndpoint{calService}
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.ResponseHeaderMiddleware)
	r.Use(middlewares.AuthorizationMiddleware(authService))
	r.HandleFunc("", e.RegisterHandler).Methods(http.MethodPost)
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
