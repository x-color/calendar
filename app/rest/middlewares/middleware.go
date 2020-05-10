package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	as "github.com/x-color/calendar/auth/service"
	"github.com/x-color/calendar/logging"
	cctx "github.com/x-color/calendar/model/ctx"
)

func ReqIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), cctx.ReqIDKey, uuid.New().String())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(logger logging.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Context().Value(cctx.ReqIDKey).(string)
			logger.Uniq(reqID).Info(r.RequestURI)
			next.ServeHTTP(w, r)
		})
	}
}

func ResponseHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func AuthorizationMiddleware(service as.Service) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err := service.Authorize(r.Context(), cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), cctx.UserIDKey, userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
