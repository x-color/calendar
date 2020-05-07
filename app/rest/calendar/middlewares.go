package calendar

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	cs "github.com/x-color/calendar/calendar/service"
	cctx "github.com/x-color/calendar/model/ctx"
)

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
