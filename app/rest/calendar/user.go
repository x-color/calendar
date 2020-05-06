package calendar

import (
	"encoding/json"
	"net/http"

	"github.com/x-color/calendar/calendar/service"
	cctx "github.com/x-color/calendar/model/ctx"
)

type userEndpoint struct {
	service service.Service
}

func (e *userEndpoint) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cctx.UserIDKey).(string)
	_, err := e.service.RegisterUser(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msgContent{"internal server error"})
		return
	}

	json.NewEncoder(w).Encode(msgContent{Msg: "register"})
}
