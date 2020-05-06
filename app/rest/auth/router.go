package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(r *mux.Router, e AuthEndpoint) {
	r.HandleFunc("/signup", e.SignupHandler).Methods(http.MethodPost)
	r.HandleFunc("/signin", e.SigninHandler)
	r.HandleFunc("/signout", e.SignoutHandler)
}
