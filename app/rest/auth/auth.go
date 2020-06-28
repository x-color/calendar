package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/x-color/calendar/app/rest/middlewares"
	"github.com/x-color/calendar/auth/service"
	cerror "github.com/x-color/calendar/model/error"
)

var secure = len(os.Getenv("SSL_DISABLE")) == 0

type userContent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type authEndpoint struct {
	service service.Service
}

func (e *authEndpoint) SignupHandler(w http.ResponseWriter, r *http.Request) {
	req := userContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := e.service.Signup(r.Context(), req.Name, req.Password)
	if errors.Is(err, cerror.ErrInvalidContent) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if errors.Is(err, cerror.ErrDuplication) {
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (e *authEndpoint) SigninHandler(w http.ResponseWriter, r *http.Request) {
	req := userContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := e.service.Signin(r.Context(), req.Name, req.Password)
	if errors.Is(err, cerror.ErrInvalidContent) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if errors.Is(err, cerror.ErrAuthorization) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.Expires,
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(userContent{
		ID: session.UserID,
	})
}

func (e *authEndpoint) SignoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = e.service.Signout(r.Context(), cookie.Value)
	if errors.Is(err, cerror.ErrAuthorization) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewRouter(r *mux.Router, s service.Service) {
	e := authEndpoint{s}
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.ResponseHeaderMiddleware)
	r.HandleFunc("/signup", e.SignupHandler).Methods(http.MethodPost)
	r.HandleFunc("/signin", e.SigninHandler)
	r.HandleFunc("/signout", e.SignoutHandler).Methods(http.MethodPost)
}
