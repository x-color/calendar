package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	cerror "github.com/x-color/calendar/model/error"
	"github.com/x-color/calendar/service/auth"
)

type AuthEndpoint struct {
	service auth.Service
}

func (e *AuthEndpoint) signupHandler(w http.ResponseWriter, r *http.Request) {
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

func (e *AuthEndpoint) signinHandler(w http.ResponseWriter, r *http.Request) {
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
		Secure:   false, // It should be 'true' if app is not sample.
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(msgContent{Msg: "signin"})
}

func (e *AuthEndpoint) signoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sign out!\n"))
}
