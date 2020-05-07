package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x-color/calendar/app/rest/middlewares"
	"github.com/x-color/calendar/auth/service"
	cerror "github.com/x-color/calendar/model/error"
)

type userContent struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type msgContent struct {
	Msg string `json:"message"`
}

type authEndpoint struct {
	service service.Service
}

func (e *authEndpoint) SignupHandler(w http.ResponseWriter, r *http.Request) {
	req := userContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	}

	user, err := e.service.Signup(r.Context(), req.Name, req.Password)
	if errors.Is(err, cerror.ErrInvalidContent) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	} else if errors.Is(err, cerror.ErrDuplication) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(msgContent{"user already exist"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msgContent{"internal server error"})
		return
	}

	json.NewEncoder(w).Encode(userContent{Name: user.Name})
}

func (e *authEndpoint) SigninHandler(w http.ResponseWriter, r *http.Request) {
	req := userContent{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	}

	session, err := e.service.Signin(r.Context(), req.Name, req.Password)
	if errors.Is(err, cerror.ErrInvalidContent) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msgContent{"bad contents"})
		return
	} else if errors.Is(err, cerror.ErrAuthorization) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(msgContent{"signin failed"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msgContent{"internal server error"})
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.Expires,
		Path:     "/",
		Secure:   false, // It should be 'true' if app is not sample.
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(msgContent{Msg: "signin"})
}

func (e *authEndpoint) SignoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(msgContent{"signout failed"})
		return
	}

	err = e.service.Signout(r.Context(), cookie.Value)
	if errors.Is(err, cerror.ErrAuthorization) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(msgContent{"signout failed"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msgContent{"internal server error"})
		return
	}

	json.NewEncoder(w).Encode(msgContent{Msg: "signout"})
}

func NewRouter(r *mux.Router, s service.Service) {
	e := authEndpoint{s}
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.ResponseHeaderMiddleware)
	r.HandleFunc("/signup", e.SignupHandler).Methods(http.MethodPost)
	r.HandleFunc("/signin", e.SigninHandler)
	r.HandleFunc("/signout", e.SignoutHandler).Methods(http.MethodPost)
}
