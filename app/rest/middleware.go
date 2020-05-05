package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	as "github.com/x-color/calendar/auth/service"
	cs "github.com/x-color/calendar/calendar/service"
	"github.com/x-color/calendar/logging"
	cctx "github.com/x-color/calendar/model/ctx"
)

type msgContent struct {
	Msg string `json:"message"`
}

func reqIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), cctx.ReqIDKey, uuid.New().String())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(logger logging.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Context().Value(cctx.ReqIDKey).(string)
			logger.Uniq(reqID).Info(r.RequestURI)
			next.ServeHTTP(w, r)
		})
	}
}

func responseHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func authorizationMiddleware(service as.Service) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(msgContent{"unauthorization"})
				return
			}

			userID, err := service.Authorize(r.Context(), cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(msgContent{"unauthorization"})
				return
			}
			ctx := context.WithValue(r.Context(), cctx.UserIDKey, userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
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
