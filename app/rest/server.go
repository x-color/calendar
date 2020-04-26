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

func newRouter(s auth.Service, l auth.Logger) *mux.Router {
	e := AuthEndpoint{s}

	r := mux.NewRouter()
	r.NotFoundHandler = http.NotFoundHandler()
	r.Use(reqIDMiddleware)
	r.Use(loggingMiddleware(l))
	r.Use(responseHeaderMiddleware)

	sr := r.PathPrefix("/auth").Subrouter()
	sr.HandleFunc("/signup", e.signupHandler).Methods(http.MethodPost)
	sr.HandleFunc("/signin", e.signinHandler)
	sr.HandleFunc("/signout", e.signoutHandler)

	return r
}
