package calendar

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-color/calendar/app/rest/middlewares"
	as "github.com/x-color/calendar/auth/service"
	"github.com/x-color/calendar/calendar/service"
	cs "github.com/x-color/calendar/calendar/service"
	cctx "github.com/x-color/calendar/model/ctx"
)

type userEndpoint struct {
	service service.Service
}

func (e *userEndpoint) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cctx.UserIDKey).(string)
	_, err := e.service.RegisterUser(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewUserRouter(r *mux.Router, calService cs.Service, authService as.Service) {
	e := userEndpoint{calService}
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.ResponseHeaderMiddleware)
	r.Use(middlewares.AuthorizationMiddleware(authService))
	r.HandleFunc("", e.RegisterHandler).Methods(http.MethodPost)
}
