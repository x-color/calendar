package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-color/calendar/app/rest/middlewares"
	"github.com/x-color/calendar/auth/service"
)

func NewRouter(r *mux.Router, s service.Service) {
	e := authEndpoint{s}
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.ResponseHeaderMiddleware)
	r.HandleFunc("/signup", e.SignupHandler).Methods(http.MethodPost)
	r.HandleFunc("/signin", e.SigninHandler)
	r.HandleFunc("/signout", e.SignoutHandler).Methods(http.MethodPost)
}
