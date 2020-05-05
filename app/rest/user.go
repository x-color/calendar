package rest

import (
	"encoding/json"
	"net/http"

	cctx "github.com/x-color/calendar/model/ctx"
	cs "github.com/x-color/calendar/service/calendar"
)

type UserEndpoint struct {
	service cs.Service
}

func (e *UserEndpoint) registerHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cctx.UserIDKey).(string)
	_, err := e.service.RegisterUser(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msgContent{"internal server error"})
		return
	}

	json.NewEncoder(w).Encode(msgContent{Msg: "register"})
}
